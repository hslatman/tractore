package signedurl

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/leg100/surl"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	secrets.SignedURLSecret = "test-secret"
	got, err := New(fmt.Sprintf("https://test.example.com/t?id=%d", 1))
	require.NoError(t, err)

	u, err := url.Parse(got)
	require.NoError(t, err)

	require.Equal(t, "1", u.Query().Get("id"))

	signer := surl.New([]byte(secrets.SignedURLSecret))

	err = signer.Verify(got)
	require.NoError(t, err)
}

func TestVerify(t *testing.T) {
	secrets.SignedURLSecret = "test-secret"
	signer := surl.New([]byte(secrets.SignedURLSecret))
	got, err := signer.Sign(fmt.Sprintf("https://test.example.com/t?id=%d", 1), time.Hour)
	require.NoError(t, err)

	u, err := url.Parse(got)
	require.NoError(t, err)

	require.Equal(t, "1", u.Query().Get("id"))
	require.NotEmpty(t, u.Query().Get("signature"))

	require.NoError(t, Verify(got))
}
