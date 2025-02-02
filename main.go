package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"OpnLaaS.cyber.unh.edu/lib"
)

type Token struct {
	Username, Token string
	Expires         time.Time
}

var tokens map[string]*Token = make(map[string]*Token)

func RandomString(length int) string {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", bytes)
}

func TokenFor(email string) string {
	if token, ok := tokens[email]; ok {
		if token.Expires.After(time.Now()) {
			return token.Token
		}
	}

	tokens[email] = &Token{
		Username: email,
		Token:    RandomString(32),
		Expires:  time.Now().Add(time.Minute * 10),
	}

	return tokens[email].Token
}

func withAuth(w http.ResponseWriter, r *http.Request) bool {
	email, err := r.Cookie("email")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	token, err := r.Cookie("token")

	if err != nil || token.Value != TokenFor(strings.ToLower(email.Value)) {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func withCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
}

func main() {
	if err := lib.InitEnv(); err != nil {
		lib.Log.Error("Could not initialize environment: " + err.Error())
		return
	} else {
		lib.Log.Status("Successfully initialized environment")
	}

	// Basic http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		withCors(w, r)
		http.ServeFile(w, r, "index.html")
	})

	lib.Log.Status(fmt.Sprintf("Server started on port %d", lib.Config.Port))
	var at string = fmt.Sprintf("%s:%d", lib.Config.Host, lib.Config.Port)

	if lib.Config.TlsDir != "" {
		http.ListenAndServeTLS(at, lib.Config.TlsDir+"/fullchain.pem", lib.Config.TlsDir+"/privkey.pem", nil)
	} else {
		http.ListenAndServe(at, nil)
	}
}
