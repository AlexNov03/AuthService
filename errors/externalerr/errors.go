package externalerr

import (
	"encoding/json"
	"errors"
	"net/http"

	interr "github.com/AlexNov03/AuthService/errors/internalerr"
	"github.com/AlexNov03/AuthService/models"
)

func ProcessInternalServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(models.ExternalError{Error: message})
}

func ProcessBadRequestError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(models.ExternalError{Error: message})
}

func ProcessAlreadyExistsError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(models.ExternalError{Error: message})
}

func ProcessUnauthorizedError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(models.ExternalError{Error: message})
}

func ProcessError(w http.ResponseWriter, err error) {
	var internalError *interr.InternalError
	if err != nil {
		if ok := errors.As(err, &internalError); ok {
			w.WriteHeader(internalError.Code)
			json.NewEncoder(w).Encode(models.ExternalError{Error: internalError.Message})
			return
		}
		ProcessInternalServerError(w, "unknown internal server error")
		return
	}
}
