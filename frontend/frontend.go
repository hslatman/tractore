// Service frontend serves the frontend for development purposes.
package frontend

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"encore.dev/rlog"
	"github.com/gorilla/mux"
)

var (
	//go:embed ui
	ui embed.FS
)

var secrets struct {
	MercureUsername string // the Mercure username
	MercurePassword string // the Mercure password
}

//encore:api private
func Dummy(context.Context) error { return nil }

// Serve serves the frontend for development using a raw endpoint.
// Learn more: https://encore.dev/docs/primitives/services-and-apis#raw-endpoints
// For production use we recommend deploying the frontend
// using Vercel, Netlify, or similar.
//
//encore:api public raw path=/frontend/*path
func Serve(w http.ResponseWriter, req *http.Request) {
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

	config.Webroot = "/frontend/"
	config.Version = "dev"

	VerifyConfig()

	serverRoot, err := fs.Sub(ui, "ui")
	if err != nil {
		rlog.Error("failed getting server root", "error", err)

		w.WriteHeader(500)
		w.Write(nil)
		return
	}

	r := apiRoutes()

	// redirect to webroot if no trailing slash
	if config.Webroot != "/" {
		redirect := strings.TrimRight(config.Webroot, "/")
		r.HandleFunc(redirect, middleWareFunc(addSlashToWebroot)).Methods("GET")
	}

	h := mux.NewRouter()
	h.HandleFunc(config.Webroot+"index.html", middleWareFunc(index))
	h.HandleFunc(config.Webroot, middleWareFunc(index))
	h.PathPrefix(config.Webroot + "dist/").Handler(middlewareHandler(http.StripPrefix(config.Webroot, http.FileServer(http.FS(serverRoot)))))
	h.PathPrefix(config.Webroot + "api/").Handler(middlewareHandler(http.StripPrefix(config.Webroot, http.FileServer(http.FS(serverRoot)))))
	h.Path(config.Webroot + "favicon.ico").Handler(middlewareHandler(http.StripPrefix(config.Webroot, http.FileServer(http.FS(serverRoot)))))
	h.Path(config.Webroot + "favicon.svg").Handler(middlewareHandler(http.StripPrefix(config.Webroot, http.FileServer(http.FS(serverRoot)))))
	h.Path(config.Webroot + "mailpit.svg").Handler(middlewareHandler(http.StripPrefix(config.Webroot, http.FileServer(http.FS(serverRoot)))))
	h.Path(config.Webroot + "notification.png").Handler(middlewareHandler(http.StripPrefix(config.Webroot, http.FileServer(http.FS(serverRoot)))))
	h.PathPrefix(config.Webroot + "view/").Handler(middleWareFunc(index)).Methods("GET")
	h.Path(config.Webroot + "search").Handler(middleWareFunc(index)).Methods("GET")

	h.ServeHTTP(w, req)
}

// Just returns the default HTML template
func index(w http.ResponseWriter, _ *http.Request) {

	var h = `<!DOCTYPE html>
<html lang="en" class="h-100">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
	<meta name="referrer" content="no-referrer">
	<meta name="robots" content="noindex, nofollow, noarchive">
	<link rel="icon" href="{{ .Webroot }}favicon.svg">
	<title>Tractore</title>
	<link rel=stylesheet href="{{ .Webroot }}dist/app.css?{{ .Version }}">
</head>

<body class="h-100">
	<div class="container-fluid h-100 d-flex flex-column" id="app" data-webroot="{{ .Webroot }}" data-version="{{ .Version }}">
		<noscript>You require JavaScript to use this app.</noscript>
	</div>

	<script src="{{ .Webroot }}dist/app.js?{{ .Version }}"></script>
</body>

</html>`

	t, err := template.New("index").Parse(h)
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
		return
	}

	data := struct {
		Webroot string
		Version string
	}{
		Webroot: config.Webroot,
		Version: config.Version,
	}

	buff := new(bytes.Buffer)

	err = t.Execute(buff, data)
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
		return
	}

	buff.Bytes()

	w.Header().Add("Content-Type", "text/html")
	_, _ = w.Write(buff.Bytes())
}

func apiRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(config.Webroot+"api/v1/webui", middleWareFunc(webUIConfig)).Methods("GET")

	return r
}

// Response includes global web UI settings
type webUIConfiguration struct {
	// Message Relay information
	MessageRelay struct {
		// Whether message relaying (release) is enabled
		Enabled bool
		// The configured SMTP server address
		SMTPServer string
		// Enforced Return-Path (if set) for relay bounces
		ReturnPath string
		// Allowlist of accepted recipients
		RecipientAllowlist string
	}

	// Whether the HTML check has been globally disabled
	DisableHTMLCheck bool

	// Whether SpamAssassin is enabled
	SpamAssassin bool

	// Whether messages with duplicate IDs are ignored
	DuplicatesIgnored bool
}

// WebUIConfig returns configuration settings for the web UI.
func webUIConfig(w http.ResponseWriter, _ *http.Request) {
	conf := webUIConfiguration{}

	conf.DisableHTMLCheck = false
	conf.SpamAssassin = false
	conf.DuplicatesIgnored = false

	bytes, _ := json.Marshal(conf)

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// MiddleWareFunc http middleware adds optional basic authentication
// and gzip compression.
func middleWareFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

// MiddlewareHandler http middleware adds optional basic authentication
// and gzip compression
func middlewareHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		h.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

// Redirect to webroot
func addSlashToWebroot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.Webroot, http.StatusFound)
}

var config struct {
	Webroot                string
	Version                string
	ContentSecurityPolicy  string
	BlockRemoteCSSAndFonts bool
}

// VerifyConfig wil do some basic checking
func VerifyConfig() error {
	cssFontRestriction := "*"
	if config.BlockRemoteCSSAndFonts {
		cssFontRestriction = "'self'"
	}

	config.ContentSecurityPolicy = fmt.Sprintf("default-src 'self'; script-src 'self'; style-src %s 'unsafe-inline'; frame-src 'self'; img-src * data: blob:; font-src %s data:; media-src 'self'; connect-src 'self' ws: wss:; object-src 'none'; base-uri 'self';",
		cssFontRestriction, cssFontRestriction,
	)

	return nil
}
