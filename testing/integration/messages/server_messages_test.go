package messages

import (
	"bytes"
	"encoding/json"
	"fmt"
	msgs "github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/server/handler"
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

func TestMessages_canCreateAndGetAMessage(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	// Create the message.
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, requestString(t, "POST", "/messages", `{"message": "my message"}`))
	require.Equal(t, http.StatusCreated, rr.Code)
	loc := rr.Header().Get("Location")
	require.True(t, regexp.MustCompile("^/messages/[0-9]+$").MatchString(loc))

	id := messageIdFromLocation(t, loc)
	msg, err := svc.MessagesService.Read(id)
	require.NoError(t, err)

	// Then retrieve it.
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, requestEmpty(t, "GET", loc))
	requireJsonOk(t, rr2)
	var m messages.MessageResponseJSON
	require.NoError(t, json.Unmarshal(rr2.Body.Bytes(), &m))

	// Check the json.
	require.Equal(t, msg.Id, m.Id)
	require.Equal(t, msg.Version, m.Version)
	require.Equal(t, "my message", m.Message)
	require.True(t, msg.CreatedAt.Equal(*m.CreatedAt))
	require.True(t, msg.UpdatedAt.Equal(*m.UpdatedAt))

	// Check the headers.
	require.Equal(t, handler.LastModifiedFormat(msg.UpdatedAt), rr2.Header().Get("Last-Modified"))
	require.Equal(t, `"1"`, rr2.Header().Get("ETag"))
}

func messageIdFromLocation(t *testing.T, uri string) msgs.MessageId {
	var id msgs.MessageId
	_, err := fmt.Fscanf(bytes.NewBufferString(uri), "/messages/%d", &id)
	require.NoError(t, err)
	return id
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
			requireJson(t, rr)
			require.Equal(t, c.rs, rr.Body.String())
		})
	}
}

func TestMessage_canDeleteMessageAnd404WhenMessageDoesNotExist(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	id, err := svc.MessagesService.Create(msgs.ModifyMessage{Message: "my message"})
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, requestEmpty(t, "DELETE", uris.Message(id)))
	require.Equal(t, http.StatusOK, rr.Code)

	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, requestEmpty(t, "GET", uris.Message(id)))
	require.Equal(t, http.StatusNotFound, rr2.Code)
}

func TestMessage_canUpdateMessage(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	id, err := svc.MessagesService.Create(msgs.ModifyMessage{Message: "first message"})
	require.NoError(t, err)

	// Update the message.
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, requestString(t, "PUT", uris.Message(id), `{"message": "new message"}`))
	require.Equal(t, http.StatusOK, rr.Code)

	// Then retrieve it.
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, requestEmpty(t, "GET", uris.Message(id)))
	require.Equal(t, http.StatusOK, rr2.Code)
	var m messages.MessageResponseJSON
	require.NoError(t, json.Unmarshal(rr2.Body.Bytes(), &m))
	require.Equal(t, "new message", m.Message)
}

func TestMessage_canListMessages(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	id1, err := svc.MessagesService.Create(msgs.ModifyMessage{Message: "first message"})
	require.NoError(t, err)
	_, err = svc.MessagesService.Create(msgs.ModifyMessage{Message: "second message"})
	require.NoError(t, err)
	_, err = svc.MessagesService.Create(msgs.ModifyMessage{Message: "atttta"})
	require.NoError(t, err)
	_, err = svc.MessagesService.Create(msgs.ModifyMessage{Message: "last message"})
	require.NoError(t, err)

	t.Run("with fields, pageSize, and pageStartIndex", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, requestEmpty(t, "GET", "/messages?fields=message,isPalindrome&pageSize=2&pageStartIndex=1"))
		requireJsonOk(t, rr)
		expected := fmt.Sprintf(`{"messages":[` +
			`{"message":"atttta","isPalindrome":true},` +
			`{"message":"last message","isPalindrome":false}` +
			"]}\n")
		require.Equal(t, expected, rr.Body.String())
	})

	t.Run("show all fields when non specified in query", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, requestEmpty(t, "GET", "/messages?pageSize=1"))
		requireJsonOk(t, rr)
		var data messages.MessageListResponseJSON
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &data))
		require.Len(t, data.Messages, 1)
		require.Equal(t, id1, data.Messages[0].Id)
		require.Equal(t, 1, data.Messages[0].Version)
		require.NotNil(t, data.Messages[0].CreatedAt)
		require.NotNil(t, data.Messages[0].UpdatedAt)
		require.Equal(t, "first message", data.Messages[0].Message)
	})
}

func TestMessage_whenListingMessages_errors(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	t.Run("error when invalid field name", func(t *testing.T) {
		rr := httptest.NewRecorder()
		serve(t, db, rr, requestEmpty(t, "GET", "/messages?fields=id,notAField,version"))
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"errors\":[{\"error\":\"invalid messages fields: notAField\"}]}\n", rr.Body.String())
	})

	t.Run("error when invalid page limit", func(t *testing.T) {
		rr := httptest.NewRecorder()
		serve(t, db, rr, requestEmpty(t, "GET", "/messages?pageSize=badSize"))
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"errors\":[\"invalid pageSize value\"]}\n", rr.Body.String())
	})

	t.Run("error when invalid page start index", func(t *testing.T) {
		rr := httptest.NewRecorder()
		serve(t, db, rr, requestEmpty(t, "GET", "/messages?pageStartIndex=badIndex"))
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"errors\":[\"invalid pageStartIndex value\"]}\n", rr.Body.String())
	})
}

func requireJsonOk(t *testing.T, rr *httptest.ResponseRecorder) {
	require.Equal(t, http.StatusOK, rr.Code)
	requireJson(t, rr)
}

func requireJson(t *testing.T, rr *httptest.ResponseRecorder) {
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}
