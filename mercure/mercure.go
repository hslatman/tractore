package mercure

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"encore.dev/rlog"
	"github.com/dunglas/mercure"
	"github.com/golang-jwt/jwt/v5"

	"encore.app/pkg/events"
)

//encore:service
type Service struct {
	hub   *mercure.Hub
	token string
}

var singleton *Service

// claims contains Mercure's JWT claims.
type claims struct {
	Mercure mercureClaim `json:"mercure"`
	jwt.RegisteredClaims
}

type mercureClaim struct {
	Publish   []string    `json:"publish"`
	Subscribe []string    `json:"subscribe"`
	Payload   interface{} `json:"payload"`
}

func initService() (*Service, error) {
	if singleton != nil {
		return singleton, nil
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed creating Ed25519 key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims{
		Mercure: mercureClaim{
			Publish: []string{"*"}, // allow internal publisher to publish on all topics
		},
	})

	signed, err := token.SignedString(priv)
	if err != nil {
		return nil, fmt.Errorf("failed signing token: %w", err)
	}

	hub, err := mercure.NewHub(
		mercure.WithAnonymous(), // TODO: maybe only let authorized users subscribe, after removing HTTP basic auth
		mercure.WithDebug(),
		mercure.WithDemo(),
		mercure.WithPublisherJWTKeyFunc(func(_ *jwt.Token) (interface{}, error) {
			return pub, nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating mercure hub: %w", err)
	}

	singleton = &Service{
		hub:   hub,
		token: signed,
	}

	return singleton, nil
}

var secrets struct {
	MercureUsername string // the Mercure username
	MercurePassword string // the Mercure password
}

//encore:api public raw path=/.well-known/mercure
func (s *Service) Serve(w http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	usernameHash := sha256.Sum256([]byte(username))
	passwordHash := sha256.Sum256([]byte(password))
	expectedUsernameHash := sha256.Sum256([]byte(secrets.MercureUsername))
	expectedPasswordHash := sha256.Sum256([]byte(secrets.MercurePassword))

	usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
	passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

	if !usernameMatch || !passwordMatch {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	s.hub.ServeHTTP(w, req)
}

//encore:api private
func PublishIncoming(ctx context.Context, m *events.IncomingEmail) error {
	if singleton == nil {
		return errors.New("mercure singleton instance not available")
	}

	return singleton.publish(ctx, m)
}

//encore:api private
func PublishOutgoing(ctx context.Context, m *events.OutgoingEmail) error {
	if singleton == nil {
		return errors.New("mercure singleton instance not available")
	}

	return singleton.publish(ctx, m)
}

//encore:api private
func PublishSent(ctx context.Context, m *events.SentEmail) error {
	if singleton == nil {
		return errors.New("mercure singleton instance not available")
	}

	return singleton.publish(ctx, m)
}

//encore:api private
func PublishTrack(ctx context.Context, m *events.TrackEmail) error {
	if singleton == nil {
		return errors.New("mercure singleton instance not available")
	}

	return singleton.publish(ctx, m)
}

func (s *Service) publish(ctx context.Context, m any) error {
	// TODO: filter out which requests to serve to whom? E.g. only to to/from certain user/email topic?
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed marshaling email: %w", err)
	}

	data := url.Values{}
	data.Set("topic", "mail")
	data.Set("data", string(b))
	enc := data.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", "/.well-known/mercure", strings.NewReader(enc))
	if err != nil {
		return fmt.Errorf("failed creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	rw := httptest.NewRecorder()

	rlog.Debug("publishing message", "enc", enc)

	s.hub.PublishHandler(rw, req)

	if rw.Code >= 300 {
		b, _ = io.ReadAll(rw.Body)
		rlog.Error("failed sending message", "error", string(b))
		return fmt.Errorf("publishing failed with code %d", rw.Code)
	}

	rlog.Info("finished publish", "code", rw.Code)

	return nil
}
