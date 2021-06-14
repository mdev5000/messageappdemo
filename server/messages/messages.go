package messages

import (
	"fmt"
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

type createJSON struct {
	Message string `json:"message"`
}

type MessageResponseJSON struct {
	Id        messages.MessageId      `json:"id"`
	Version   messages.MessageVersion `json:"version"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	Message   string                  `json:"message"`
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

	var resp createJSON
	if !handler.DecodeJsonOrError(h.log, op, w, r, &resp) {
		return
	}

	id, err := h.messagesSvc.Create(messages.CreateMessage{
		Message: resp.Message,
	})
	if err != nil {
		handler.SendErrorResponse(h.log, op, w, err)
		return
	}

	w.Header().Set("Location", uris.Message(id))
	w.WriteHeader(http.StatusCreated)
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

	handler.EncodeJsonOrError(op, h.log, w, MessageResponseJSON{
		Id:        message.Id,
		Version:   message.Version,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		Message:   message.Message,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update")
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
	fmt.Println("list")
}
