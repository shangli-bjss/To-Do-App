package server

import (
	"net/http"
)

func RegisterRoutes() {
    http.HandleFunc("/todos", GetAndPostHandler)
    http.HandleFunc("/todos/", PutAndDeleteHandler)
}
