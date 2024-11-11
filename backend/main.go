package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	Id          int
	Name        string
	Credentials []webauthn.Credential
}

var UserStore = make(map[int]User)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	wconfig := &webauthn.Config{
		RPDisplayName: "Test Webauthn",
		RPID:          "webauthn.local",
		RPOrigins:     []string{"https:login.webauthn.local"},
	}

	webauthn, err := webauthn.New(wconfig)
	if err != nil {
		slog.Error("could not create new webauthn", "error", err)
		return
	}

	fmt.Printf("Webauthn: %#v\n", webauthn)

	mux := http.NewServeMux()
	mux.HandleFunc("/register/{id}", BeginRegistration)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func GetUser(id int) (User, error) {
	user, ok := UserStore[id]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	slog.Debug("BeginRegistraion")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("Invalid User ID", "id", r.PathValue("id"))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	slog.Info("Begin registratin", "id", id, "user", user)
}
