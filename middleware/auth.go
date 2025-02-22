package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/AlexNov03/AuthService/errors/externalerr"
)

type SessionUsecase interface {
	AddSession(sessionID string, userID string)
	GetSessionUser(sessionID string) (string, error)
	DeleteSession(sessionID string)
}

type Middleware struct {
	uc SessionUsecase
}

func NewMiddleware(usecase SessionUsecase) *Middleware {
	return &Middleware{uc: usecase}
}

func (m *Middleware) RequireAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if errors.Is(err, http.ErrNoCookie) {
			externalerr.ProcessUnauthorizedError(w, "method require authorization")
			return
		}

		userID, err := m.uc.GetSessionUser(cookie.Value)

		if err != nil {
			externalerr.ProcessError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)

	})
}
