package sessionrepo

import (
	"net/http"
	"sync"

	interr "github.com/AlexNov03/AuthService/errors/internalerr"
)

// **************************************************************
// *												   			*
// *     Mock БД на структурах для тестирования        			*
// *  Так делать нехорошо, скоро сессии будут лежать в Reddis   *
// *                                                   			*
// **************************************************************

type SessionRepo struct {
	storage map[string]string
	mu      sync.RWMutex
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		storage: make(map[string]string),
		mu:      sync.RWMutex{},
	}
}

func (sr *SessionRepo) AddSession(sessionID string, userID string) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.storage[sessionID] = userID
}

func (sr *SessionRepo) GetSessionUser(sessionID string) (string, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	if val, ok := sr.storage[sessionID]; ok {
		return val, nil
	}
	return "", interr.NewInternalError(http.StatusNotFound, "repository: no user has this session ID")
}

func (sr *SessionRepo) DeleteSession(sessionID string) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	delete(sr.storage, sessionID)
}
