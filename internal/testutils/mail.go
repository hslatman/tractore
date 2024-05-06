package testutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jhillyerd/enmime"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func EncodeMail(t *testing.T, b enmime.MailBuilder) []byte {
	t.Helper()
	p, err := b.Build()
	require.NoError(t, err)

	var w bytes.Buffer
	require.NoError(t, p.Encode(&w))

	return w.Bytes()
}

func ExtractTrackingURL(t *testing.T, raw []byte) string {
	t.Helper()
	env, err := enmime.ReadEnvelope(bytes.NewBuffer(raw))
	require.NoError(t, err)

	doc, err := html.Parse(strings.NewReader(env.HTML))
	require.NoError(t, err)

	var findTrackingURL func(*html.Node) string
	findTrackingURL = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					return a.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if trackingURL := findTrackingURL(c); trackingURL != "" {
				return trackingURL
			}
		}

		return ""
	}

	return findTrackingURL(doc)
}
