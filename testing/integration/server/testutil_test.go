package server

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/mdev5000/messageappdemo/approot"
	"github.com/mdev5000/messageappdemo/logging"
	"github.com/mdev5000/messageappdemo/postgres"
	"github.com/mdev5000/messageappdemo/server"
	"github.com/stretchr/testify/require"
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
	h, err := server.Handler(server.Services{Log: logging.NoLog()}, server.Config{LogRequest: false})
	require.NoError(t, err)
	return h
}

func noDbServe(t *testing.T, w http.ResponseWriter, r *http.Request) {
	noDbHandler(t).ServeHTTP(w, r)
}

func handlerWithDb(t *testing.T, db *postgres.DB) (http.Handler, *approot.Services) {
	log := logging.NoLog()
	svcs := approot.Setup(db, log)
	svch := server.Services{
		Log:             svcs.Log,
		MessagesService: svcs.MessagesService,
	}
	h, err := server.Handler(svch, server.Config{LogRequest: false})
	require.NoError(t, err)
	return h, svcs
}

func serve(t *testing.T, db *postgres.DB, w http.ResponseWriter, r *http.Request) {
	h, _ := handlerWithDb(t, db)
	h.ServeHTTP(w, r)
}

func requestEmpty(t *testing.T, method, url string) *http.Request {
	return request(t, method, url, nil)
}

func request(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req
}

func requestString(t *testing.T, method, url string, body string) *http.Request {
	b := bytes.NewBufferString(body)
	return request(t, method, url, b)
}
