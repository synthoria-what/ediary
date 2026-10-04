package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"synthori/ediary/m/models"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(&data); err != nil {
		fmt.Printf("error when encoding data")
	}
}

func WriteSuccess[T any](w http.ResponseWriter, status int, message string, data T) {
	now := time.Now().UTC()

	response := models.Response[T]{
		Message:    message,
		StatusCode: status,
		Timestamp:  now,
		Data:       &data,
	}

	WriteJSON(w, status, response)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	now := time.Now().UTC()

	response := models.Response[string]{
		Message:    err.Error(),
		StatusCode: status,
		Timestamp:  now,
		Error: &models.APIError{
			Code:    http.StatusText(status),
			Details: err.Error(),
		},
	}

	WriteJSON(w, status, response)
}
