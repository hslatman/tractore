package smtpd

import (
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct {
	c chan Mail
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		c: bkd.c,
	}, nil
}

// A Session is returned after successful login.
type Session struct {
	c    chan Mail
	from string
	to   string
}

// AuthMechanisms returns a slice of available auth mechanisms; only PLAIN is
// supported in this example.
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators.
func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != "username" || password != "password" {
			return errors.New("invalid username or password")
		}
		return nil
	}), nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	slog.Info("incoming mail", "from", s.from, "to", s.to)

	s.c <- Mail{
		From: s.from,
		To:   s.to,
		Raw:  b,
	}

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

var _ smtp.Backend = (*Backend)(nil)
var _ smtp.Session = (*Session)(nil)

type Mail struct {
	From string
	To   string
	Raw  []byte
}

// New returns a new [smtp.Server] and a channel that returns [Mail]
// structs. After starting the server, any email sent to it will be
// sent to the channel, so that it can be further processed. The server
// (currently) does not require authentication or encryption, so don't
// run this in an untrusted environment.
func New() (*smtp.Server, chan Mail) {
	c := make(chan Mail, 100)
	be := &Backend{
		c: c,
	}

	s := smtp.NewServer(be)

	s.Addr = "localhost:1025"
	s.Domain = "localhost"
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	return s, c
}
