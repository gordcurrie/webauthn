package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
)

type sessionStore struct {
	sessions map[string]*webauthn.SessionData
	mu       sync.RWMutex
}

func NewSessionStore() *sessionStore {
	return &sessionStore{sessions: make(map[string]*webauthn.SessionData)}
}

func (store *sessionStore) GetSession(id string) (*webauthn.SessionData, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	session, ok := store.sessions[id]
	if !ok {
		return nil, fmt.Errorf("error getting session:  %s", id)
	}

	return session, nil
}

func (store *sessionStore) StartSession(data *webauthn.SessionData) string {
	store.mu.Lock()
	defer store.mu.Unlock()

	id := strconv.Itoa(rand.Int())

	store.sessions[id] = data

	return id
}
