package server

import (
	"fmt"
	gmux "github.com/gorilla/mux"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/mdev5000/qlik_message/messages"
	msgs "github.com/mdev5000/qlik_message/server/messages"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
)

type Services struct {
	Log             *logging.Logger
	MessagesService *messages.Service
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
		case "GET", "POST", "PUT", "DELETE":
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

	messageHandler := msgs.NewHandler(svc.Log, svc.MessagesService)
	messages := mux.PathPrefix("/messages").Subrouter()
	messages.HandleFunc("", messageHandler.Create).Methods("POST")
	messages.HandleFunc("", messageHandler.List).Methods("GET")
	messages.HandleFunc("", acceptsHandler("GET", "POST"))

	message := messages.HandleFunc("/{id}", messageHandler.Read).Subrouter()
	message.HandleFunc("", messageHandler.Read).Methods("GET")
	message.HandleFunc("", messageHandler.Update).Methods("PUT")
	message.HandleFunc("", messageHandler.Delete).Methods("DELETE")
	message.HandleFunc("", acceptsHandler("DELETE", "GET", "PUT"))

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	if cfg.LogRequest {
		n.Use(negroni.NewLogger())
	}
	n.UseHandler(mux)

	return n, nil
}

func acceptsHandler(methods ...string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		if _, err := fmt.Fprintf(w, "Allow: %s", strings.Join(methods, ", ")); err != nil {
			handleWriteError(err)
		}
	}
}

func handleWriteError(err error) {

}
