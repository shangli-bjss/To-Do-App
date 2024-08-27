package server

import (
	"net/http"
)

func registerRoutes() {
    http.HandleFunc("/todos", getAndPostHandler)
    http.HandleFunc("/todos/", putAndDeleteHandler)
}
