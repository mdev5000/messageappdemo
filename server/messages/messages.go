package messages

import (
	"github.com/gorilla/mux"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/server/handler"
	"github.com/mdev5000/qlik_message/server/uris"
	"net/http"
	"strconv"
	"time"
)

type modifyMessageJSON struct {
	Message string `json:"message"`
}

func (m *modifyMessageJSON) toModifyMessage() messages.ModifyMessage {
	return messages.ModifyMessage{
		Message: m.Message,
	}
}

type MessageListResponseJSON struct {
	Messages []MessageResponseJSON `json:"messages,omitempty"`
}

type MessageResponseJSON struct {
	Id           *messages.MessageId      `json:"id,omitempty"`
	Version      *messages.MessageVersion `json:"version,omitempty"`
	CreatedAt    *time.Time               `json:"created_at,omitempty"`
	UpdatedAt    *time.Time               `json:"updated_at,omitempty"`
	Message      string                   `json:"message,omitempty"`
	IsPalindrome bool                     `json:"isPalindrome"`
}

type Handler struct {
	log         *logging.Logger
	messagesSvc *messages.Service
}

func NewHandler(log *logging.Logger, messageSvc *messages.Service) *Handler {
	return &Handler{
		log:         log,
		messagesSvc: messageSvc,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "MessagesHandler.Create"

	var resp modifyMessageJSON
	if !handler.DecodeJsonOrError(h.log, op, w, r, &resp) {
		return
	}

	id, err := h.messagesSvc.Create(resp.toModifyMessage())
	if err != nil {
		handler.SendErrorResponse(h.log, op, w, err)
		return
	}

	w.Header().Set("Location", uris.Message(id))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	const op = "MessagesHandler.Read"

	id, ok := h.readIdFromUri(op, w, r)
	if !ok {
		return
	}

	message, err := h.messagesSvc.Read(id)
	if err != nil {
		handler.SendErrorResponse(h.log, op, w, err)
		return
	}

	handler.EncodeJsonOrError(op, h.log, w, messageToJsonValue(message))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "MessagesHandler.Update"

	id, ok := h.readIdFromUri(op, w, r)
	if !ok {
		return
	}

	var resp modifyMessageJSON
	if !handler.DecodeJsonOrError(h.log, op, w, r, &resp) {
		return
	}

	_, err := h.messagesSvc.Update(id, resp.toModifyMessage())
	if err != nil {
		handler.SendErrorResponse(h.log, op, w, err)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "MessagesHandler.Delete"

	id, ok := h.readIdFromUri(op, w, r)
	if !ok {
		return
	}

	if err := h.messagesSvc.Delete(id); err != nil {
		handler.SendErrorResponse(h.log, op, w, err)
		return
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	const op = "MessagesHandler.List"

	fields, limit, offset, err := handler.GetQueryParams(op, r)
	if err != nil {
		handler.SendErrorResponse(h.log, op, w, err)
		return
	}

	msgs, err := h.messagesSvc.List(messages.MessageQuery{
		Fields: fields,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		handler.SendErrorResponse(h.log, op, w, err)
		return
	}

	out := make([]MessageResponseJSON, len(msgs))
	for i, msg := range msgs {
		out[i] = messageToJsonValue(msg)
	}

	handler.EncodeJsonOrError(op, h.log, w, MessageListResponseJSON{Messages: out})
}

func (h *Handler) readIdFromUri(op string, w http.ResponseWriter, r *http.Request) (messages.MessageId, bool) {
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		appErr := apperrors.Error{Op: op}
		appErr.AddResponse(apperrors.ErrorResponse("invalid message id"))
		handler.SendErrorResponse(h.log, op, w, &appErr)
		return 0, false
	}
	return messages.MessageId(id), true
}

func messageToJsonValue(message *messages.Message) MessageResponseJSON {
	emptyTime := time.Time{}
	mr := MessageResponseJSON{
		Message:      message.Message,
		IsPalindrome: message.IsPalindrome,
	}
	if message.Id != 0 {
		mr.Id = &message.Id
	}
	if message.Version != 0 {
		mr.Version = &message.Version
	}
	if !emptyTime.Equal(message.CreatedAt) {
		mr.CreatedAt = &message.CreatedAt
	}
	if !emptyTime.Equal(message.UpdatedAt) {
		mr.UpdatedAt = &message.UpdatedAt
	}
	return mr
}
