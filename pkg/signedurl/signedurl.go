package signedurl

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/leg100/surl"
)

var secrets struct {
	SignedURLSecret string // secret key for signing and verifying URLs
}

var initOnce sync.Once
var signer *surl.Signer
var initErr error

const lifespan = 365 * 24 * time.Hour // a year

func New(u string) (string, error) {
	initOnce.Do(func() {
		if secrets.SignedURLSecret == "" {
			initErr = errors.New("no URL signer secret available")
		}
		signer = surl.New([]byte(secrets.SignedURLSecret))
	})
	if initErr != nil {
		return "", fmt.Errorf("failed initializing URL signer: %w", initErr)
	}

	if signer == nil {
		return "", errors.New("no URL signer available")
	}

	// Create a signed URL that expires in one hour.
	signed, err := signer.Sign(u, lifespan)
	if err != nil {
		return "", fmt.Errorf("failed creating signed URL: %w", err)
	}

	return signed, nil
}

func Verify(u string) error {
	initOnce.Do(func() {
		if secrets.SignedURLSecret == "" {
			initErr = errors.New("no URL verifier secret available")
		}
		signer = surl.New([]byte(secrets.SignedURLSecret))
	})
	if initErr != nil {
		return fmt.Errorf("failed initializing URL signer: %w", initErr)
	}

	if signer == nil {
		return errors.New("no URL verifier available")
	}

	if err := signer.Verify(u); err != nil {
		return fmt.Errorf("failed verifying URL signature: %w", err)
	}

	return nil
}
