package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	userStore    = NewUserStore()
	sessionstore = NewSessionStore()
	webAuthn     *webauthn.WebAuthn
)

func main() {
	var err error
	slog.SetLogLoggerLevel(slog.LevelDebug)
	wconfig := &webauthn.Config{
		RPDisplayName: "Gord's Webauthn",
		RPID:          "webauthn.local",
		RPOrigins:     []string{"https:login.webauthn.local"},
	}

	webAuthn, err = webauthn.New(wconfig)
	if err != nil {
		slog.Error("could not create new webauthn", "error", err)
		return
	}

	fmt.Printf("Webauthn: %#v\n", webAuthn)

	mux := http.NewServeMux()
	mux.HandleFunc("/register/{username}", BeginRegistration)

	mux.Handle("/", http.FileServer(http.Dir("./backend")))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func writeResponse(w http.ResponseWriter, msg any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		slog.Error("Could not marshal msg", "msg", msg)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	slog.Debug("BeginRegistration Start")
	defer slog.Debug("BeginRegistration end")
	username := r.PathValue("username")
	if username == "" {
		slog.Error("Invalid User Name")
		http.Error(w, "Invalid user name", http.StatusBadRequest)
		return
	}

	user, err := userStore.GetUser(username)
	if err != nil {
		if err != ErrorUserNotFound {
			slog.Error("unable to get user", "username", username, "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		slog.Info("user not found, creating new user", "username", username)
		user = NewUser(username)
		userStore.AddUser(user)
	}
	options, sessionData, err := webAuthn.BeginRegistration(user)
	if err != nil {
		slog.Error("could not begin registration", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	sessionId := sessionstore.StartSession(sessionData)
	cookie := &http.Cookie{
		Name:  "registration",
		Value: sessionId,
		Path:  "/",
	}

	http.SetCookie(w, cookie)

	writeResponse(w, options, http.StatusOK)
}
