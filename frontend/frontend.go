// Service frontend serves the frontend for development purposes.
package frontend

import (
	"crypto/sha256"
	"crypto/subtle"
	"embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed dist
	dist embed.FS

	assets, _ = fs.Sub(dist, "dist")
	handler   = http.StripPrefix("/frontend/", http.FileServer(http.FS(assets)))
)

var secrets struct {
	MercureUsername string // the Mercure username
	MercurePassword string // the Mercure password
}

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

	handler.ServeHTTP(w, req)
}
