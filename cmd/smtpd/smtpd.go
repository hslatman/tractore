package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"log/slog"
	"sync"

	"encore.app/pkg/client"
	"encore.app/pkg/smtpd"
	"github.com/jhillyerd/enmime"
)

var token string

func init() {
	flag.StringVar(&token, "token", "", "the API token for the mail service")
}

func main() {
	flag.Parse()

	if token == "" {
		panic(errors.New("'token' is required"))
	}

	c, err := client.New(client.Local, client.WithAuthToken(token))
	if err != nil {
		panic(err)
	}

	s, ch := smtpd.New()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("starting server ...", "addr", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			slog.Error("failed serving SMTP server", "error", err)
		}
	}()

	go func() {
		for m := range ch {
			ctx := context.Background()

			// try to read the enveloped message as basic validation
			_, err := enmime.ReadEnvelope(bytes.NewBuffer(m.Raw))
			if err != nil {
				slog.Error("failed reading raw mail", "error", err)
				continue
			}

			err = c.Receiver.Ingest(ctx, client.ReceiverIngestRequest{
				To:   m.To,
				From: m.From,
				Raw:  base64.StdEncoding.EncodeToString(m.Raw),
			})

			if err != nil {
				slog.Error("failed sending message", "error", err)
			} else {
				slog.Info("sent message")
			}
		}
	}()

	wg.Wait()
}
