package server

import (
	"fmt"
	gmux "github.com/gorilla/mux"
	msgs "github.com/mdev5000/qlik_message/server/messages"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
)

type Config struct {
	LogRequest bool
}

func Handler(cfg Config) (http.Handler, error) {
	mux := gmux.NewRouter()

	messageHandler := msgs.NewHandler()
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
