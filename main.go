package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"OpnLaaS.cyber.unh.edu/database"
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

	if !database.Connect() {
		return
	}

	metadataJSON, _ := json.Marshal(map[string]interface{}{
		"name":                 lib.Config.LabName,
		"organization":         lib.Config.LabOrg,
		"contact":              lib.Config.LabContact,
		"emailDomainWhiteList": lib.Config.EmailDomainWhiteList,
	})

	http.HandleFunc("/metadata.json", func(w http.ResponseWriter, r *http.Request) {
		withCors(w, r)
		w.Header().Set("Content-Type", "application/json")
		w.Write(metadataJSON)
	})

	// Create
	http.HandleFunc("/api/user/create", func(w http.ResponseWriter, r *http.Request) {
		withCors(w, r)

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		obj := struct {
			Email    string `json:"email"`
			First    string `json:"firstName"`
			Last     string `json:"lastName"`
			Password string `json:"password"`
		}{}

		err := json.Unmarshal(body, &obj)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, statusCode := database.CreateUser(strings.ToLower(obj.Email), obj.First, obj.Last, database.HashPassword(obj.Password))

		if statusCode != nil {
			switch statusCode {
			case database.ErrUserExists:
				w.WriteHeader(http.StatusIMUsed)
			case database.ErrBadData:
				w.WriteHeader(http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    strings.ToLower(obj.Email),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    TokenFor(strings.ToLower(obj.Email)),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write(user.JSON())

		lib.Log.Basic(fmt.Sprintf("User %s created", obj.Email))
	})

	// Login
	http.HandleFunc("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		withCors(w, r)
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		obj := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		err := json.Unmarshal(body, &obj)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := database.GetUser(strings.ToLower(obj.Email))

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user.PasswordHash != database.HashPassword(obj.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    strings.ToLower(obj.Email),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    TokenFor(strings.ToLower(obj.Email)),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		w.WriteHeader(http.StatusOK)

		lib.Log.Basic(fmt.Sprintf("User %s logged in", obj.Email))
	})

	// Me
	http.HandleFunc("/api/user/me", func(w http.ResponseWriter, r *http.Request) {
		withCors(w, r)
		if !withAuth(w, r) {
			return
		}

		email, _ := r.Cookie("email")
		user, err := database.GetUser(email.Value)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(user.JSON())
	})

	// Logout
	http.HandleFunc("/api/user/logout", func(w http.ResponseWriter, r *http.Request) {
		withCors(w, r)

		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    "",
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		w.WriteHeader(http.StatusOK)

		if email, err := r.Cookie("email"); err == nil {
			lib.Log.Basic(fmt.Sprintf("User %s logged out", email.Value))
		} else {
			lib.Log.Basic("Unknown user logged out")
		}
	})

	// Delete user (self only)
	http.HandleFunc("/api/user/delete", func(w http.ResponseWriter, r *http.Request) {
		withCors(w, r)

		if !withAuth(w, r) {
			return
		}

		if r.Method != "DELETE" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		email, _ := r.Cookie("email")
		err := database.DeleteUser(email.Value)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "email",
			Value:    "",
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})

		w.WriteHeader(http.StatusOK)

		lib.Log.Basic(fmt.Sprintf("User %s deleted", email.Value))
	})

	lib.Log.Status(fmt.Sprintf("Server started on port %d", lib.Config.Port))
	var at string = fmt.Sprintf("%s:%d", lib.Config.Host, lib.Config.Port)

	if lib.Config.TlsDir != "" {
		http.ListenAndServeTLS(at, lib.Config.TlsDir+"/fullchain.pem", lib.Config.TlsDir+"/privkey.pem", nil)
	} else {
		http.ListenAndServe(at, nil)
	}
}
