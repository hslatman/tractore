package mercure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/r3labs/sse/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"encore.app/pkg/events"
)

func TestService(t *testing.T) {
	secrets.MercureUsername = "mercure-testuser"
	secrets.MercurePassword = "mercure-testpass"

	svc, err := initService()
	require.NoError(t, err)

	srv := newServer(t, svc)
	defer srv.Close()

	c := sse.NewClient(fmt.Sprintf("%s/.well-known/mercure?topic=mail", srv.URL))
	httpClient := &http.Client{
		Transport: basicAuth(http.DefaultClient.Do),
	}
	c.Connection = httpClient

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // TODO: make this less flaky?
	defer cancel()

	numberOfMessages := 0

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		c.OnDisconnect(func(c *sse.Client) {
			wg.Done()
		})
		err = c.SubscribeWithContext(ctx, "", func(msg *sse.Event) {
			numberOfMessages += 1
			var m = &events.MercureMessage{}
			require.NoError(t, json.Unmarshal(msg.Data, m))
			assert.Equal(t, &events.MercureMessage{
				ID:    1,
				To:    "to@example.com",
				From:  "from@example.com",
				State: "received",
			}, m)
		})
		require.NoError(t, err)
	}()

	// hacky way to allow the client some time to subscribe and
	// get its read loop started. It seems the OnConnect callback
	// doesn't work for this purpose, somehow?
	time.Sleep(time.Second)

	svc.PublishIncoming(ctx, &events.MercureMessage{
		ID:    1,
		To:    "to@example.com",
		From:  "from@example.com",
		State: "received",
	})

	wg.Wait()

	require.Equal(t, 1, numberOfMessages)
}

// roundTripFunc wraps a func to make it into an http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements RoundTripper interface
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func basicAuth(next roundTripFunc) roundTripFunc {
	return roundTripFunc(func(req *http.Request) (*http.Response, error) {
		req.SetBasicAuth(secrets.MercureUsername, secrets.MercurePassword)
		return next(req)
	})
}

func newServer(t *testing.T, svc *Service) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/mercure", svc.Serve)
	server := httptest.NewServer(mux)

	return server
}
