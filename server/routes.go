package server

import (
	"net/http"
)

func RegisterRoutes() {
    http.HandleFunc("/todos", getAndPostHandler)
    http.HandleFunc("/todos/", putAndDeleteHandler)
}
