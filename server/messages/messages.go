package messages

import (
	"fmt"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Create(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("create")
}

func (h *Handler) Read(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("read")
}

func (h *Handler) Update(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("update")
}

func (h *Handler) Delete(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("delete")
}

func (h *Handler) List(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("list")
}
