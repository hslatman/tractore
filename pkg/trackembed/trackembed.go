package trackembed

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jhillyerd/enmime"
	"golang.org/x/net/html"
)

// Pixel embeds an HTML "img" tag pointing to the tracking URL in the
// (raw) email message. If the raw email message does not contain HTML,
// the existing text content is transformed into HTML and added as an
// additional email part. If the raw email did already contain HTML, the
// "img" tag will be included at the end of the "body" node.
func Pixel(raw []byte, trackingURL string) ([]byte, error) {
	env, err := enmime.ReadEnvelope(bytes.NewBuffer(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to read envelope: %w", err)
	}

	htmlContent := env.HTML
	if htmlContent == "" {
		htmlContent = fmt.Sprintf(`<html><p>%s</p></html>`, env.Text) // TODO: nicer text to HTML formatting?
	}

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed parsing HTML: %w", err)
	}

	img := &html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			{Key: "src", Val: trackingURL},
			{Key: "width", Val: "1"},
			{Key: "height", Val: "1"},
		},
	}

	bodyNode := doc.LastChild.LastChild // root -> html -> body
	bodyNode.AppendChild(img)

	var hw bytes.Buffer
	if err := html.Render(&hw, doc); err != nil {
		return nil, fmt.Errorf("failed rendering new HTML: %w", err)
	}

	newHTMLContent := hw.Bytes()

	from, fromAddress, err := split(env.GetHeader("From"))
	if err != nil {
		return nil, fmt.Errorf("failed getting From header: %w", err)
	}
	to, toAddress, err := split(env.GetHeader("To"))
	if err != nil {
		return nil, fmt.Errorf("failed getting To header: %w", err)
	}

	bldr := enmime.Builder().
		From(from, fromAddress).
		To(to, toAddress).
		Subject(env.GetHeader("Subject")).
		Text([]byte(env.Text)).
		HTML(newHTMLContent)

	p, err := bldr.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build email: %w", err)
	}

	var w bytes.Buffer
	if err := p.Encode(&w); err != nil {
		return nil, fmt.Errorf("failed encoding email: %w", err)
	}

	return w.Bytes(), nil
}

func split(s string) (string, string, error) {
	i1 := strings.Index(s, "<")
	i2 := strings.Index(s, ">")
	if i1 < 2 {
		return s, s, nil
	}
	p1 := s[:i1-1]
	p2 := s[i1+1 : i2]
	return p1, p2, nil
}
