package server

import (
	"bytes"
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

var allMethods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"OPTIONS",
	"TRACE",
	"PATCH",
}

func allMethodsExcept(methods ...string) []string {
	numMethods := len(allMethods) - len(methods)
	if numMethods < 1 {
		panic("invalid methods list at least one existing method must not exist in the list")
	}
	out := make([]string, numMethods)
	i := 0
Top:
	for _, method := range allMethods {
		for _, exclude := range methods {
			if method == exclude {
				continue Top
			}
		}
		out[i] = method
		i++
	}
	return out
}

func noDbHandler(t *testing.T) http.Handler {
	h, err := Handler(Services{}, Config{LogRequest: false})
	require.NoError(t, err)
	return h
}

func noDbServe(t *testing.T, w http.ResponseWriter, r *http.Request) {
	noDbHandler(t).ServeHTTP(w, r)
}

func handlerWithDb(t *testing.T, db *postgres.DB) http.Handler {
	log := logging.NoLog()
	messagesRepo := data.NewMessageRepository(db)
	h, err := Handler(Services{
		Log:             log,
		MessagesService: messages.NewService(log, messagesRepo),
	}, Config{LogRequest: false})
	require.NoError(t, err)
	return h
}

//func serve(t *testing.T, db *postgres.DB, w http.ResponseWriter, r *http.Request) {
//	handlerWithDb(t, db).ServeHTTP(w, r)
//}

func requestEmpty(t *testing.T, method, url string) *http.Request {
	return request(t, method, url, nil)
}

func request(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func requestString(t *testing.T, method, url string, body string) *http.Request {
	b := bytes.NewBufferString(body)
	return request(t, method, url, b)
}
