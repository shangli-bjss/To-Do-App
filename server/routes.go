package server

import (
	"net/http"
)

func registerRoutes() {
    http.HandleFunc("/todos", todosHandler)
    http.HandleFunc("/todos/", todosByIdHandler)
}
