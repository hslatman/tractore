package wasm

import (
	_ "embed"
	"testing"

	"github.com/jhillyerd/enmime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"encore.app/internal/testutils"
)

func TestPlugin_EmbedPixel(t *testing.T) {
	textMail := enmime.Builder().
		From("test", "test@example.com").
		To("to", "to@example.com").
		Subject("Testing").
		Text([]byte("testing ..."))

	htmlMail := enmime.Builder().
		From("test", "test@example.com").
		To("to", "to@example.com").
		Subject("Testing").
		Text([]byte("testing with HTML ...")).
		HTML([]byte("<html><body><p>testing with HTML ...</p></body></html>"))

	type args struct {
		raw         []byte
		trackingURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "text",
			args: args{
				raw:         testutils.EncodeMail(t, textMail),
				trackingURL: "http://test.example.com/t?id=text",
			},
		},
		{
			name: "html",
			args: args{
				raw:         testutils.EncodeMail(t, htmlMail),
				trackingURL: "http://test.example.com/t?id=html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New()
			require.NoError(t, err)

			got, err := p.EmbedPixel(tt.args.raw, tt.args.trackingURL)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
				return
			}

			trackingURL := testutils.ExtractTrackingURL(t, got)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.trackingURL, trackingURL)
		})
	}
}
