package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

const maxJSONBodyBytes = 1 << 20 // 1 MB

type APIError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIError{Error: message})
}

func writeInternalError(w http.ResponseWriter, message string, err error) {
	if err != nil {
		slog.Error(message, "error", err)
	}
	writeError(w, http.StatusInternalServerError, message)
}

// decodeJSON safely decodes a bounded JSON request body, rejecting unknown fields and trailing data.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "corpo da requisição excede o limite permitido")
			return false
		}
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return false
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return false
	}

	return true
}
