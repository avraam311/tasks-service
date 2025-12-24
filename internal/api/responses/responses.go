package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ErrCodeNotFound   = "NOT_FOUND"
	ErrInternalServer = "INTERNAL ERROR"
	ErrInvalidJSON    = "INVALID JSON"
)

type Success struct {
	Result interface{} `json:"result"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("responses/responses.go - failed to send json")
	}
	return nil
}

func ResponseOK(w http.ResponseWriter, result interface{}) error {
	err := WriteJSON(w, http.StatusOK, Success{Result: result})
	if err != nil {
		return err
	}
	return nil
}

func ResponseCreated(w http.ResponseWriter, result interface{}) error {
	err := WriteJSON(w, http.StatusCreated, Success{Result: result})
	if err != nil {
		return err
	}
	return nil
}

func ResponseError(w http.ResponseWriter, code string, message string, statusCode int) error {
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message

	err := WriteJSON(w, statusCode, resp)
	if err != nil {
		return err
	}
	return nil
}
