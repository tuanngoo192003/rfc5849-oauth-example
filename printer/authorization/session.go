package authorization

import "sync"

// Session storage (in-memory for demo)
type Session struct {
	mu           sync.RWMutex
	tempToken    string
	tempSecret   string
	accessToken  string
	accessSecret string
	photos       []Photo
	logs         []string
}

type Photo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Owner string `json:"owner"`
}

var session = &Session{
	logs: []string{},
}

func GetSession() *Session {
	return session
}

func (s *Session) GetMu() *sync.RWMutex {
	return &s.mu
}

func (s *Session) GetTempToken() string {
	return s.tempToken
}

func (s *Session) SetTempToken(tempToken string) {
	s.tempToken = tempToken
}

func (s *Session) GetTempSecret() string {
	return s.tempSecret
}

func (s *Session) SetTempSecret(tempSecret string) {
	s.tempSecret = tempSecret
}

func (s *Session) GetAccessToken() string {
	return s.accessToken
}

func (s *Session) SetAccessToken(accessToken string) {
	s.accessToken = accessToken
}

func (s *Session) GetAccessSecret() string {
	return s.accessSecret
}

func (s *Session) SetAccessSecret(accessSecret string) {
	s.accessSecret = accessSecret
}

func (s *Session) GetPhotos() []Photo {
	return s.photos
}

func (s *Session) SetPhotos(photos []Photo) {
	s.photos = append(s.photos, photos...)
}

func (s *Session) GetLogs() []string {
	return s.logs
}

func (s *Session) SetLogs(logs []string) {
	s.logs = append(s.logs, logs...)
}
