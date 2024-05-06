package processor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"encore.dev"
	"encore.dev/beta/errs"
	"encore.dev/metrics"
	"encore.dev/pubsub"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"x.encore.dev/infra/pubsub/outbox"

	"encore.app/mercure"
	"encore.app/pkg/events"
	"encore.app/pkg/wasm"
	"encore.app/processor/ent"
	"encore.app/receiver"
)

var OutgoingEmailsTopic = pubsub.NewTopic[*events.OutgoingEmail]("outgoing", pubsub.TopicConfig{DeliveryGuarantee: pubsub.ExactlyOnce})

var outgoingEmailsRef = pubsub.TopicRef[pubsub.Publisher[*events.OutgoingEmail]](OutgoingEmailsTopic)

var EmailsProcessed = metrics.NewCounter[uint64]("emails_processed", metrics.CounterConfig{})

//encore:service
type Service struct {
	ent  *ent.Client
	wasm *wasm.Plugin
}

func initService() (*Service, error) {
	// Initialize the relay to poll from our database.
	relay := outbox.NewRelay(outbox.SQLDBStore(db))

	outbox.RegisterTopic(relay, outgoingEmailsRef)

	rlog.Info("starting polling for messages")

	go relay.PollForMessages(context.Background(), -1) // TODO: shutdown nicely

	driver := entsql.OpenDB(dialect.Postgres, db.Stdlib())
	entClient := ent.NewClient(ent.Driver(driver))

	plugin, err := wasm.New()
	if err != nil {
		return nil, fmt.Errorf("failed instantiating Wasm plugin: %w", err)
	}

	return &Service{
		ent:  entClient,
		wasm: plugin,
	}, nil
}

func (s *Service) Process(ctx context.Context, im *events.IncomingEmail) error {
	tx, err := s.ent.BeginTx(ctx, nil)
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed to start transaction: %w", err),
			},
		}
	}
	defer tx.Rollback()

	b, err := base64.StdEncoding.DecodeString(im.Raw)
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed to decode raw mail: %w", err),
			},
		}
	}

	m, err := tx.Mail.
		Create().
		SetFrom(im.From).
		SetTo(im.To).
		SetRaw(b).
		SetIncomingMailID(im.ID).
		Save(ctx)
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed to save incoming mail: %w", err),
			},
		}
	}

	if err := tx.Commit(); err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed to commit: %w", err),
			},
		}
	}

	meta := encore.Meta()
	baseURL := meta.APIBaseURL.String()
	trackingURL := fmt.Sprintf("%s/t?id=%d", baseURL, im.ID) // TODO: use a protected ID

	rlog.Info("embedding tracking pixel", "url", trackingURL)

	raw, err := s.wasm.EmbedPixel(m.Raw, trackingURL)
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed embedding pixel: %w", err),
			},
		}
	}

	e := &events.OutgoingEmail{
		ID:   im.ID,
		To:   im.To,
		From: im.From,
		Raw:  base64.StdEncoding.EncodeToString(raw),
	}

	eb, err := json.Marshal(e)
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed to marshal outgoing email: %w", err),
			},
		}
	}

	_, err = s.ent.Outbox.
		Create().
		SetTopic(OutgoingEmailsTopic.Meta().Name).
		SetInsertedAt(time.Now()).
		SetData(eb).
		Save(ctx)
	if err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed to create outbox entry for outgoing email: %w", err),
			},
		}
	}

	// TODO: set some state in the mercure event, so that UI can update it?

	if err := mercure.PublishOutgoing(ctx, e); err != nil {
		return &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": fmt.Errorf("failed publishing update to mercure: %w", err),
			},
		}
	}

	EmailsProcessed.Increment()

	return nil
}

var _ = pubsub.NewSubscription(
	receiver.IncomingEmailsTopic, "process",
	pubsub.SubscriptionConfig[*events.IncomingEmail]{
		Handler: pubsub.MethodHandler((*Service).Process),
	},
)

var db = sqldb.NewDatabase("processor", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
