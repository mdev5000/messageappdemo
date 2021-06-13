package main

import (
	"github.com/mdev5000/qlik_message/server"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	handler, err := server.Handler(server.Config{
		LogRequest: true,
	})
	if err != nil {
		return err
	}
	return http.ListenAndServe("localhost:8000", handler)
}
