package server

import (
	"fmt"
	gmux "github.com/gorilla/mux"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/logging"
	msgs "github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/server/handler"
	msgh "github.com/mdev5000/qlik_message/server/messages"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
)

type Services struct {
	Log             *logging.Logger
	MessagesService *msgs.Service
}

type Config struct {
	LogRequest bool
}

func standardServiceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This service only accepts json, so indicate to the client.
		w.Header().Set("Accept", "application/json")

		// Ensure the content type is correctly set to json
		switch r.Method {
		case "POST", "PUT":
			if r.Header.Get("Content-Type") != "application/json" {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

func Handler(svc Services, cfg Config) (http.Handler, error) {
	mux := gmux.NewRouter()
	mux.Use(standardServiceMiddleware)

	messageHandler := msgh.NewHandler(svc.Log, svc.MessagesService)
	messages := mux.PathPrefix("/messages").Subrouter()
	messages.HandleFunc("", messageHandler.Create).Methods("POST")
	messages.HandleFunc("", messageHandler.List).Methods("GET")
	messages.HandleFunc("", acceptsHandler(svc.Log, "GET", "POST"))

	message := messages.HandleFunc("/{id}", messageHandler.Read).Subrouter()
	message.HandleFunc("", messageHandler.Read).Methods("GET")
	message.HandleFunc("", messageHandler.Update).Methods("PUT")
	message.HandleFunc("", messageHandler.Delete).Methods("DELETE")
	message.HandleFunc("", acceptsHandler(svc.Log, "DELETE", "GET", "PUT"))

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	if cfg.LogRequest {
		n.Use(negroni.NewLogger())
	}
	n.UseHandler(mux)

	return n, nil
}

func acceptsHandler(log *logging.Logger, methods ...string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "server.acceptsHandler"
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		if _, err := fmt.Fprintf(w, "Allow: %s", strings.Join(methods, ", ")); err != nil {
			handler.SendErrorResponse(log, op, w, &apperrors.Error{
				EType: apperrors.ETInternal,
				Op:    op,
				Err:   err,
				Stack: errors.WithStack(err),
			})
			return
		}
	}
}
