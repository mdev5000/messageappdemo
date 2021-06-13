package server

import (
	"encoding/json"
	"github.com/mdev5000/qlik_message/server/messages"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestMessages_405WhenMethodIsNotSupported(t *testing.T) {
	t.Parallel()
	h := noDbHandler(t)
	for _, method := range allMethodsExcept("GET", "POST", "OPTIONS") {
		t.Run("method "+method, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, requestEmpty(t, method, "/messages"))
			require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
			require.Equal(t, "Allow: GET, POST", rr.Body.String())
		})
	}
}

func TestMessages_canListOptions(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	noDbServe(t, rr, requestEmpty(t, "OPTIONS", "/messages"))
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "Allow: GET, POST", rr.Body.String())
}

func TestMessages_canCreateAMessage(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	hander := handlerWithDb(t, db)

	// Create the message.
	rr := httptest.NewRecorder()
	hander.ServeHTTP(rr, requestString(t, "POST", "/messages", `{"message": "my message"}`))
	require.Equal(t, http.StatusCreated, rr.Code)
	loc := rr.Header().Get("Location")
	require.True(t, regexp.MustCompile("^/messages/[0-9]+$").MatchString(loc))

	// Then retrieve it
	rr2 := httptest.NewRecorder()
	hander.ServeHTTP(rr2, requestEmpty(t, "GET", loc))
	require.Equal(t, http.StatusOK, rr2.Code)
	var m messages.MessageResponseJSON
	require.NoError(t, json.Unmarshal(rr2.Body.Bytes(), &m))
	require.Equal(t, "my message", m.Message)
}
