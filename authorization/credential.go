package authorization

import (
	"sync"
	"time"
)

// Storage for OAuth tokens and temporary credentials
type OAuthStore struct {
	mu                    sync.RWMutex
	tempCredentials       map[string]*TempCredential
	accessTokens          map[string]*AccessToken
	authorizedCredentials map[string]string // maps temp token to verifier
}

func (o *OAuthStore) GetMu() *sync.RWMutex {
	return &o.mu
}

func (o *OAuthStore) GetTemporaryCredentials() map[string]*TempCredential {
	return o.tempCredentials
}

func (o *OAuthStore) GetAccessToken() map[string]*AccessToken {
	return o.accessTokens
}

func (o *OAuthStore) GetAuthorizedCredentials() map[string]string {
	return o.authorizedCredentials
}

type TempCredential struct {
	Token      string
	Secret     string
	Callback   string
	Timestamp  time.Time
	Authorized bool
	Username   string
}

type AccessToken struct {
	Token    string
	Secret   string
	Username string
}

var store = &OAuthStore{
	tempCredentials:       make(map[string]*TempCredential),
	accessTokens:          make(map[string]*AccessToken),
	authorizedCredentials: make(map[string]string),
}

func GetStore() *OAuthStore {
	return store
}
