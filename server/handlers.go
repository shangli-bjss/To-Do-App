package server

import (
	"net/http"
	"strings"
)

func getAndPostHandler(w http.ResponseWriter, req *http.Request){
	switch req.Method{
	case http.MethodGet:
		getTodo(w)
	case http.MethodPost:
		postTodo(w, req)
	default:
		http.Error(w, "Request Method not allowed", http.StatusMethodNotAllowed)
	}
}

func putAndDeleteHandler(w http.ResponseWriter, req *http.Request){
	id := strings.TrimPrefix(req.URL.Path, "/todos/")
	if id == "" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
	}
	switch req.Method {
	case http.MethodPut:
		putTodo(w, req, id)
	case http.MethodDelete:
		deleteTodo(w, id)
	default:
		http.Error(w, "Request Method not allowed", http.StatusMethodNotAllowed)
	}
}

