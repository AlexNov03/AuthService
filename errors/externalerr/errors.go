package externalerr

import (
	"encoding/json"
	"net/http"

	"github.com/AlexNov03/AuthService/errors/internalerr"
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

func ProcessInternalError(w http.ResponseWriter, error *internalerr.InternalError) {
	w.WriteHeader(error.Code)
	json.NewEncoder(w).Encode(models.ExternalError{Error: error.Message})
}
