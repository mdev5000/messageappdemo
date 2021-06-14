package server

import (
	"encoding/json"
	msgs "github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/server/messages"
	"github.com/mdev5000/qlik_message/server/uris"
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

	h, _ := handlerWithDb(t, db)

	// Create the message.
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, requestString(t, "POST", "/messages", `{"message": "my message"}`))
	require.Equal(t, http.StatusCreated, rr.Code)
	loc := rr.Header().Get("Location")
	require.True(t, regexp.MustCompile("^/messages/[0-9]+$").MatchString(loc))

	// Then retrieve it
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, requestEmpty(t, "GET", loc))
	require.Equal(t, http.StatusOK, rr2.Code)
	var m messages.MessageResponseJSON
	require.NoError(t, json.Unmarshal(rr2.Body.Bytes(), &m))
	require.Equal(t, "my message", m.Message)
}

func TestMessages_returnsErrorWhen(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		rq   string
		rs   string
	}{
		{
			"empty message",
			`{"message": ""}`,
			`{"errors":[{"field":"message","error":"Message field cannot be blank."}]}` + "\n",
		},
		{
			"invalid json",
			`{{`,
			`{"errors":[{"error":"invalid json"}]}` + "\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			noDbServe(t, rr, requestString(t, "POST", "/messages", c.rq))
			require.Equal(t, http.StatusBadRequest, rr.Code)
			require.Equal(t, c.rs, rr.Body.String())
		})
	}
}

func TestMessage_canDeleteMessageAnd404WhenMessageDoesNotExist(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	id, err := svc.MessagesService.Create(msgs.CreateMessage{Message: "my message"})
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, requestEmpty(t, "DELETE", uris.Message(id)))
	require.Equal(t, http.StatusOK, rr.Code)

	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, requestEmpty(t, "GET", uris.Message(id)))
	require.Equal(t, http.StatusNotFound, rr2.Code)
}
