package sessiondelivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AlexNov03/AuthService/errors/externalerr"
	"github.com/AlexNov03/AuthService/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SessionUsecase interface {
	AddSession(sessionID string, userID string)
	GetSessionUser(sessionID string) (string, error)
	DeleteSession(sessionID string)
}

type UserUsecase interface {
	SignUp(regData *models.UserRegData) (*models.UserInfo, error)
	Login(loginData *models.UserLoginData) (*models.UserInfo, error)
}

type SessionDelivery struct {
	su SessionUsecase
	uu UserUsecase
}

func NewSessionDelivery(sessionuc SessionUsecase, useruc UserUsecase) *SessionDelivery {
	return &SessionDelivery{
		su: sessionuc,
		uu: useruc,
	}
}

func (sd *SessionDelivery) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	loginData := &models.UserLoginData{}

	err := json.NewDecoder(r.Body).Decode(loginData)

	if err != nil {
		externalerr.ProcessBadRequestError(w, err.Error())
		return
	}

	validator := validator.New(validator.WithRequiredStructEnabled())
	err = validator.Struct(loginData)

	if err != nil {
		externalerr.ProcessBadRequestError(w, err.Error())
		return
	}

	userInfo, err := sd.uu.Login(loginData)
	if err != nil {
		externalerr.ProcessError(w, err)
		return
	}

	sessionID := uuid.NewString()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(time.Hour * 12),
		HttpOnly: true,
		Secure:   true,
	}

	sd.su.AddSession(sessionID, userInfo.ID)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func (sd *SessionDelivery) LogOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")

	if errors.Is(err, http.ErrNoCookie) {
		externalerr.ProcessUnauthorizedError(w, "logout impossible: user unauthorized")
		return
	}

	sd.su.DeleteSession(cookie.Value)

	cookie.Expires = time.Now().Add(-time.Hour)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func (sd *SessionDelivery) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	regData := &models.UserRegData{}

	err := json.NewDecoder(r.Body).Decode(regData)

	if err != nil {
		externalerr.ProcessBadRequestError(w, err.Error())
		return
	}

	validator := validator.New(validator.WithRequiredStructEnabled())
	err = validator.Struct(regData)

	if err != nil {
		externalerr.ProcessBadRequestError(w, err.Error())
		return
	}

	userInfo, err := sd.uu.SignUp(regData)
	if err != nil {
		externalerr.ProcessError(w, err)
		return
	}

	sessionID := uuid.NewString()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(time.Hour * 12),
		HttpOnly: true,
		Secure:   true,
	}

	sd.su.AddSession(sessionID, userInfo.ID)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
