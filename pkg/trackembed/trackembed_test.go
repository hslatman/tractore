package trackembed

import (
	"testing"

	"encore.app/internal/testutils"
	"github.com/jhillyerd/enmime"
	"github.com/stretchr/testify/assert"
)

func TestPixel(t *testing.T) {
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
			name: "ok/text",
			args: args{
				raw:         testutils.EncodeMail(t, textMail),
				trackingURL: "http://test.example.com/t?id=text",
			},
		},
		{
			name: "ok/html",
			args: args{
				raw:         testutils.EncodeMail(t, htmlMail),
				trackingURL: "http://test.example.com/t?id=html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Pixel(tt.args.raw, tt.args.trackingURL)
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
