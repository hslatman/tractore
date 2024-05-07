// Service receiver receives emails
package receiver

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/metrics"
	"encore.dev/pubsub"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jhillyerd/enmime"
	"x.encore.dev/infra/pubsub/outbox"

	"encore.app/mercure"
	"encore.app/pkg/events"
	"encore.app/pkg/xerrs"
	"encore.app/receiver/ent"
)

//encore:service
type Service struct {
	ent *ent.Client
}

var IncomingEmailsTopic = pubsub.NewTopic[*events.IncomingEmail]("incoming", pubsub.TopicConfig{DeliveryGuarantee: pubsub.ExactlyOnce})

var incomingEmailsRef = pubsub.TopicRef[pubsub.Publisher[*events.IncomingEmail]](IncomingEmailsTopic)

var EmailsReceived = metrics.NewCounter[uint64]("emails_received", metrics.CounterConfig{})

func initService() (*Service, error) {
	if secrets.MailAPIToken == "" {
		return nil, errors.New("secret mail API token not configured")
	}

	// Initialize the relay to poll from our database.
	relay := outbox.NewRelay(outbox.SQLDBStore(db))
	outbox.RegisterTopic(relay, incomingEmailsRef)
	rlog.Info("starting polling for messages")
	go relay.PollForMessages(context.Background(), -1) // TODO: shutdown nicely

	driver := entsql.OpenDB(dialect.Postgres, db.Stdlib())
	entClient := ent.NewClient(ent.Driver(driver))

	return &Service{ent: entClient}, nil
}

func (s *Service) Shutdown(force context.Context) {
	rlog.Info("shutting down mail service")
}

type IngestRequest struct {
	To   string `json:"to"`
	From string `json:"from"`
	Raw  string `json:"raw"`
}

//encore:api auth
func (s *Service) Ingest(ctx context.Context, ir *IngestRequest) error {
	tx, err := s.ent.BeginTx(ctx, nil)
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed starting DB transaction: %w", err))
	}
	defer tx.Rollback()

	m, err := tx.IncomingMail.
		Create().
		SetFrom(ir.From).
		SetTo(ir.To).
		SetRaw([]byte(ir.Raw)).
		Save(ctx)
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed to save incoming mail: %w", err))
	}

	e := &events.IncomingEmail{
		ID:   m.ID,
		To:   ir.To,
		From: ir.From,
		Raw:  ir.Raw,
	}

	eb, err := json.Marshal(e)
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed marshaling JSON: %w", err))
	}

	_, err = tx.Outbox.
		Create().
		SetTopic(IncomingEmailsTopic.Meta().Name).
		SetInsertedAt(time.Now()).
		SetData(eb).
		Save(ctx)
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed creating outbox item: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return xerrs.Internal(fmt.Errorf("failed committing transaction: %w", err))
	}

	b, err := base64.StdEncoding.DecodeString(ir.Raw)
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed to decode raw mail: %w", err))
	}
	env, err := enmime.ReadEnvelope(bytes.NewBuffer(b))
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed to read envelope: %w", err))
	}
	date, err := env.Date()
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed to read envelope: %w", err))
	}

	var from *mail.Address
	fromData := addressToSlice(env, "From")
	if len(fromData) > 0 {
		from = fromData[0]
	} else if env.GetHeader("From") != "" {
		from = &mail.Address{Name: env.GetHeader("From")}
	}
	messageID := strings.Trim(env.GetHeader("Message-ID"), "<>")

	if err := mercure.PublishIncoming(ctx, &events.MercureMessage{
		ID:    m.ID,
		To:    ir.To,
		From:  ir.From,
		State: "received",
		Raw:   ir.Raw,
		Data: events.Data{
			ID:        strconv.Itoa(m.ID),
			MessageID: messageID,
			From:      from,
			To:        addressToSlice(env, "To"),
			Cc:        addressToSlice(env, "Cc"),
			Bcc:       addressToSlice(env, "Bcc"),
			ReplyTo:   addressToSlice(env, "Reply-To"),
			Date:      time.UnixMilli(date.UnixMilli()),
			Subject:   env.GetHeader("Subject"),
			Text:      env.Text,
			HTML:      env.HTML,
			Tags:      []string{"received"},
			Size:      len(b),
			Created:   time.Now(),
		},
	}); err != nil {
		return xerrs.Internal(fmt.Errorf("failed publishing to Mercure: %w", err))
	}

	EmailsReceived.Increment()

	return nil
}

// Return a header field as a []*mail.Address, or "null" is not found/empty
func addressToSlice(env *enmime.Envelope, key string) []*mail.Address {
	data, err := env.AddressList(key)
	if err != nil || data == nil {
		return []*mail.Address{}
	}

	return data
}

var secrets struct {
	MailAPIToken string // static API token for authentication to mail API
}

// Data can be named whatever you prefer (but must be exported).
type Data struct {
	Username string
}

// TODO: add (opportunistic) mTLS auth? I.e. TLS handshake is already
// finished, but can we still verify the client certs, if provided?

//encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, *Data, error) {
	// Validate the token and look up the user id and user data,
	// for example by calling Firebase Auth.
	if subtle.ConstantTimeCompare([]byte(token), []byte(secrets.MailAPIToken)) != 1 {
		return "", nil, errors.New("invalid API token")
	}

	return "testuser", &Data{Username: "testuser"}, nil
}

var db = sqldb.NewDatabase("receiver", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
