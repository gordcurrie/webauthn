package main

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
)

var ErrorUserNotFound = errors.New("user does not exist")

type UserStore struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]User),
	}
}

func (u *UserStore) GetUser(username string) (User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	user, ok := u.users[username]
	if !ok {
		return User{}, ErrorUserNotFound
	}
	return user, nil
}

func (u *UserStore) AddUser(user User) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users[user.Name] = user
}

type User struct {
	Id          int
	Name        string
	Credentials []webauthn.Credential
}

func NewUser(username string) User {
	return User{
		Id:   rand.Int(),
		Name: username,
	}
}

func (u User) WebAuthnID() []byte {
	id := strconv.Itoa(u.Id)
	return []byte(id)
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.Name
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return nil
}
