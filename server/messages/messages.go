package messages

import (
	"fmt"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/server/handler"
	"github.com/mdev5000/qlik_message/server/uris"
	"net/http"
)

type createdJson struct {
	Message string `json:"message"`
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
	var resp createdJson
	if err := handler.DecodeJsonOrError(w, r, &resp); err != nil {
		return
	}

	id, err := h.messagesSvc.Create(messages.CreateMessage{
		Message: resp.Message,
	})
	if err != nil {
		handler.SendErrorResponse(h.log, op, err, w)
		return
	}

	w.Header().Set("Location", uris.ReadMessage(id))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	fmt.Println("read")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete")
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("list")
}
