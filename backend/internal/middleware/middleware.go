package middleware

import (
	"log"
	"net/http"
)

// HTTPError allows handlers to specify HTTP error codes and messages.
type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

// AppHandler is a custom handler type that returns an error.
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// ErrorHandler é um middleware que captura erros e os trata de forma global
func ErrorHandler(handler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Erro recuperado (Panic): %v", err)
				http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			}
		}()

		if err := handler(w, r); err != nil {
			log.Printf("Erro na requisição: %v", err)

			// Check if the error is a specific HTTP error. If not, return a generic 500.
			if httpErr, ok := err.(HTTPError); ok {
				http.Error(w, httpErr.Message, httpErr.Code)
			} else {
				http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			}
		}
	}
}
