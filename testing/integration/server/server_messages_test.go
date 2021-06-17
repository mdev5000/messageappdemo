package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	msgs "github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/server"
	"github.com/mdev5000/qlik_message/server/handler"
	"github.com/mdev5000/qlik_message/server/messages"
	"github.com/mdev5000/qlik_message/server/uris"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// POST - /messages, GET - /messages/{id}
// --------------------------------------------

func TestMessages_canCreateAndGetAMessage(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	// Create the message.
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, requestString(t, "POST", "/messages", `{"message": "my message"}`))
	require.Equal(t, http.StatusCreated, rr.Code)
	loc := rr.Header().Get("Location")
	require.Equal(t, `"1"`, rr.Header().Get("ETag"))
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

// DELETE - /messages/{id}
// --------------------------------------------

func TestMessage_canDeleteMessageAnd404WhenMessageDoesNotExist(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	id, err := svc.MessagesService.Create(msgs.ModifyMessage{Message: "my message"})
	require.NoError(t, err)

	t.Run("can delete a message and get a 404 upon trying to request again", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, requestEmpty(t, "DELETE", uris.Message(id)))
		require.Equal(t, http.StatusOK, rr.Code)

		rr2 := httptest.NewRecorder()
		h.ServeHTTP(rr2, requestEmpty(t, "GET", uris.Message(id)))
		require.Equal(t, http.StatusNotFound, rr2.Code)
	})

	t.Run("delete returns 200 for non-existent message", func(t *testing.T) {
		// See https://stackoverflow.com/questions/6474223/should-deleting-a-non-existent-resource-result-in-a-404-in-restful-rails
		// for why this is the case.
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, requestEmpty(t, "DELETE", uris.Message(5)))
		require.Equal(t, http.StatusOK, rr.Code)
	})
}

// PUT - /messages/{id}
// --------------------------------------------

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
	require.Equal(t, `"2"`, rr.Header().Get("ETag"))
	var m messages.MessageResponseJSON
	require.NoError(t, json.Unmarshal(rr2.Body.Bytes(), &m))
	require.Equal(t, "new message", m.Message)
}

// GET - /messages
// --------------------------------------------

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

	t.Run("show all messages by default", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, requestEmpty(t, "GET", "/messages?fields=message"))
		requireJsonOk(t, rr)
		expected := fmt.Sprintf(`{"messages":[` +
			`{"message":"first message"},` +
			`{"message":"second message"},` +
			`{"message":"atttta"},` +
			`{"message":"last message"}` +
			`]}`)
		require.Equal(t, expected, rr.Body.String())
	})

	t.Run("with fields, pageSize, and pageStartIndex", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, requestEmpty(t, "GET", "/messages?fields=message,isPalindrome&pageSize=2&pageStartIndex=1"))
		requireJsonOk(t, rr)
		expected := fmt.Sprintf(`{"messages":[` +
			`{"message":"atttta","isPalindrome":true},` +
			`{"message":"last message","isPalindrome":false}` +
			`]}`)
		require.Equal(t, expected, rr.Body.String())
	})

	t.Run("show all fields when no field filter specified in query", func(t *testing.T) {
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
		require.Equal(t, "{\"errors\":[{\"error\":\"invalid messages fields: notAField\"}]}", rr.Body.String())
	})

	t.Run("error when invalid page limit", func(t *testing.T) {
		rr := httptest.NewRecorder()
		serve(t, db, rr, requestEmpty(t, "GET", "/messages?pageSize=badSize"))
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"errors\":[\"invalid pageSize value\"]}", rr.Body.String())
	})

	t.Run("error when invalid page start index", func(t *testing.T) {
		rr := httptest.NewRecorder()
		serve(t, db, rr, requestEmpty(t, "GET", "/messages?pageStartIndex=badIndex"))
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"errors\":[\"invalid pageStartIndex value\"]}", rr.Body.String())
	})
}

// * - /messages/{id}
// --------------------------------------------

func TestMessage_errorOnBadId(t *testing.T) {
	for _, method := range []string{"GET", "PUT", "DELETE"} {
		t.Run(fmt.Sprintf("error when id invalid for %s", method), func(t *testing.T) {
			rr := httptest.NewRecorder()
			noDbServe(t, rr, requestEmpty(t, "GET", "/messages/duck"))
			require.Equal(t, http.StatusBadRequest, rr.Code)
			requireJson(t, rr)
			require.Equal(t, "{\"errors\":[{\"error\":\"invalid message id\"}]}", rr.Body.String())
		})
	}
}

func TestMessage_404OnNotFound(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, _ := handlerWithDb(t, db)

	cases := []struct {
		method string
		body   string
	}{
		{"GET", ""},
		{"HEAD", ""},
		{"PUT", `{"message":"my message"}`},
	}

	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, requestString(t, c.method, uris.Message(5), c.body))
			require.Equal(t, http.StatusNotFound, rr.Code)
		})
	}
}

func TestMessages_createOrUpdateReturnsErrorWhen(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		rq   string
		rs   string
	}{
		{
			"empty message",
			`{"message": ""}`,
			`{"errors":[{"field":"message","error":"Message field cannot be blank."}]}`,
		},
		{
			"invalid json",
			`{{`,
			`{"errors":[{"error":"invalid json"}]}`,
		},
		{
			"empty json",
			``,
			`{"errors":[{"error":"invalid json"}]}`,
		},
	}
	methods := []struct {
		method string
		path   string
	}{
		{"POST", "/messages"},
		{"PUT", uris.Message(1)},
	}
	for _, method := range methods {
		for _, c := range cases {
			t.Run(method.method+" - "+c.name, func(t *testing.T) {
				rr := httptest.NewRecorder()
				noDbServe(t, rr, requestString(t, method.method, method.path, c.rq))
				require.Equal(t, http.StatusBadRequest, rr.Code)
				requireJson(t, rr)
				require.Equal(t, c.rs, rr.Body.String())
			})
		}
	}
}

func TestMessages_OptionsMethodAndReturns405WhenMethodIsNotSupported(t *testing.T) {
	t.Parallel()
	cases := []struct {
		uri        string
		badMethods []string
		allow      string
	}{
		{
			"/messages",
			allMethodsExcept("GET", "POST", "OPTIONS", "HEAD"),
			"GET, HEAD, OPTIONS, POST"},
		{
			"/messages/1",
			allMethodsExcept("GET", "PUT", "DELETE", "OPTIONS", "HEAD"),
			"DELETE, GET, HEAD, OPTIONS, PUT"},
	}

	h := noDbHandler(t)
	for _, c := range cases {

		t.Run("OPTIONS returns Allow for "+c.uri, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, requestEmpty(t, "OPTIONS", c.uri))
			require.Equal(t, http.StatusOK, rr.Code)
			require.Equal(t, "Allow: "+c.allow, rr.Body.String())
		})

		for _, method := range c.badMethods {

			t.Run("504 when method "+method+" "+c.uri, func(t *testing.T) {
				rr := httptest.NewRecorder()
				h.ServeHTTP(rr, requestEmpty(t, method, c.uri))
				require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
				require.Equal(t, "Allow: "+c.allow, rr.Body.String())
			})
		}
	}
}

func TestMessages_headRequests(t *testing.T) {
	db, dbClose := acquireDb(t)
	defer dbClose()

	h, svc := handlerWithDb(t, db)

	id, err := svc.MessagesService.Create(msgs.ModifyMessage{Message: "message"})
	require.NoError(t, err)

	cases := []string{
		"/messages",
		uris.Message(id),
	}

	for _, uri := range cases {
		t.Run("HEAD for "+uri, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, requestEmpty(t, "HEAD", uri))
			require.Equal(t, http.StatusOK, rr.Code)
			// Check no body is sent in the HEAD request.
			require.Equal(t, "", rr.Body.String())
		})
	}
}

func TestMessages_406NotAcceptableWhenInvalidContentType(t *testing.T) {
	t.Parallel()
	h := noDbHandler(t)

	cases := []struct {
		method      string
		uri         string
		contentType string
		body        string
	}{
		{"POST", "/messages", "application/xml", "<xml></xml>"},
		{"PUT", uris.Message(5), "application/xml", "<xml></xml>"},
	}
	for _, c := range cases {
		t.Run(c.method+" "+c.uri+" content type: "+c.contentType, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := requestString(t, c.method, c.uri, c.body)
			req.Header.Set("Content-Type", c.contentType)
			h.ServeHTTP(rr, req)
			require.Equal(t, http.StatusNotAcceptable, rr.Code)
			require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Accept"))
		})
	}
}

// * - /messages/{id} (security)
// --------------------------------------------

func TestMessages_secureHeadersApplied(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	// These headers are currently applied to all paths so just check that they are there.
	noDbServe(t, rr, requestEmpty(t, "DELETE", uris.Message(7)))
	require.Equal(t, "deny", rr.Header().Get("X-Frame-Options"))
	require.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "frame-ancestors 'none'", rr.Header().Get("Content-Security-Policy"))
}

func TestMessages_errorWhenRequestBodyIsTooBig(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	// This setting is applied to all paths so just check it exists.
	size := server.MaxBodySize + 1
	body := make([]byte, size)
	// Make sure the start is valid json.
	n := copy(body, `{"message": "long"`)
	copy(body[n:], strings.Repeat(" ", size-n-1))
	body[len(body)-1] = '}'
	noDbServe(t, rr, request(t, "POST", "/messages", bytes.NewBuffer(body)))
	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, `{"errors":[{"error":"request body too large"}]}`, rr.Body.String())
}

// helpers

func requireJsonOk(t *testing.T, rr *httptest.ResponseRecorder) {
	require.Equal(t, http.StatusOK, rr.Code)
	requireJson(t, rr)
}

func requireJson(t *testing.T, rr *httptest.ResponseRecorder) {
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
}
