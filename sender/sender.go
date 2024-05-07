// Service sender sends processed emails
package sender

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"encore.dev/metrics"
	"encore.dev/pubsub"
	"encore.dev/rlog"
	"github.com/jhillyerd/enmime"
	"github.com/mailgun/mailgun-go/v4"

	"encore.app/mercure"
	"encore.app/pkg/events"
	"encore.app/pkg/xerrs"
	"encore.app/processor"
)

var EmailsSent = metrics.NewCounter[uint64]("emails_sent", metrics.CounterConfig{})

//encore:service
type Service struct {
}

var secrets struct {
	MailgunDomain string // the Mailgun domain
	MailgunAPIKey string // the Mailgun API key
}

// Send sends an email using Mailgun
func (s *Service) Send(ctx context.Context, im *events.OutgoingEmail) error {
	raw, err := base64.StdEncoding.DecodeString(im.Raw)
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed to decode raw mail: %w", err))
	}

	env, err := enmime.ReadEnvelope(bytes.NewBuffer(raw))
	if err != nil {
		return xerrs.Internal(fmt.Errorf("failed to read envelope: %w", err))
	}

	subject := env.GetHeader("Subject")
	text := env.Text

	if secrets.MailgunDomain != "" && secrets.MailgunAPIKey != "" {
		mg := mailgun.NewMailgun(secrets.MailgunDomain, secrets.MailgunAPIKey)
		m := mg.NewMessage(
			im.From,
			subject,
			text,
			im.To,
		)
		m.SetHtml(env.HTML)
		_, id, err := mg.Send(ctx, m)
		if err != nil {
			return xerrs.Internal(fmt.Errorf("failed to send email: %w", err))
		}
		rlog.Debug("sent email", "to", im.To, "from", im.From, "mailgun_id", id)
	} else {
		rlog.Warn("skipped sending mail because Mailgun secrets are not set")
	}

	e := &events.MercureMessage{
		ID:    im.ID,
		From:  im.From,
		To:    im.To,
		State: "sent",
		Raw:   im.Raw,
	}

	if err := mercure.PublishSent(ctx, e); err != nil {
		return xerrs.Internal(fmt.Errorf("failed publishing to Mercure: %w", err))
	}

	EmailsSent.Increment()

	return nil
}

var _ = pubsub.NewSubscription(
	processor.OutgoingEmailsTopic, "send",
	pubsub.SubscriptionConfig[*events.OutgoingEmail]{
		Handler: pubsub.MethodHandler((*Service).Send),
	},
)
