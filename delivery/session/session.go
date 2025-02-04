package sessiondelivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AlexNov03/AuthService/errors/externalerr"
	interr "github.com/AlexNov03/AuthService/errors/internalerr"
	"github.com/AlexNov03/AuthService/models"
	"github.com/google/uuid"
)

type SessionUsecase interface {
	AddSession(sessionID string, userID string)
	GetSessionUser(sessionID string) (string, error)
	DeleteSession(sessionID string)
}

type UserUsecase interface {
	GetUserInfo(userID string) (*models.UserInfo, error)
	AddUserInfo(userInfo *models.UserRegData) string
	FindUser(loginData *models.UserLoginData) (*models.UserInfo, error)
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
		externalerr.ProcessBadRequestError(w, "login input data is not valid")
		return
	}

	userInfo, err := sd.uu.FindUser(loginData)

	var internalError *interr.InternalError
	if err != nil {
		if ok := errors.As(err, &internalError); ok {
			externalerr.ProcessInternalError(w, internalError)
			return
		}
		externalerr.ProcessInternalServerError(w, "unknown internal server error")
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
		externalerr.ProcessUnauthorizedError(w, "logout imposiible: user unauthorized")
		return
	}

	sd.su.DeleteSession(cookie.Value)

	cookie.Expires = time.Now().Add(-time.Hour)

	http.SetCookie(w, cookie)

}

func (sd *SessionDelivery) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	regData := &models.UserRegData{}

	err := json.NewDecoder(r.Body).Decode(regData)

	if err != nil {
		externalerr.ProcessBadRequestError(w, "signup input data is not valid")
		return
	}

	_, err = sd.uu.FindUser(&models.UserLoginData{Email: regData.Email, Password: regData.Password})

	if err == nil {
		externalerr.ProcessAlreadyExistsError(w, "this user already exists")
		return
	}

	var internalError *interr.InternalError

	if ok := errors.As(err, &internalError); ok {
		if internalError.Code != http.StatusNotFound {
			externalerr.ProcessInternalError(w, internalError)
			return
		}
	} else {
		externalerr.ProcessInternalServerError(w, "unknown internal server error")
		return
	}

	userID := sd.uu.AddUserInfo(regData)

	sessionID := uuid.NewString()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(time.Hour * 12),
		HttpOnly: true,
		Secure:   true,
	}

	sd.su.AddSession(sessionID, userID)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
