package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
)


type ToDo struct {
	Item string `json:"item"`
	Status string `json:"status"`
}

var todoList = make([]ToDo, 0)

func main() {
	http.HandleFunc("/todos", getAndPostHandler)
	http.HandleFunc("/todos/", putAndDeleteHandler)

	pl("Starting server on port: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}

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
	idStr := strings.TrimPrefix(req.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > len(todoList) {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}

	switch req.Method {
	case http.MethodPut:
		putTodo(w, req, id-1)
	case http.MethodDelete:
		deleteTodo(w, id-1)
	default:
		http.Error(w, "Request Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTodo(w http.ResponseWriter){
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todoList)
}

func postTodo(w http.ResponseWriter, req *http.Request){
	var newTodo ToDo
	if err := json.NewDecoder(req.Body).Decode(&newTodo); err != nil { // Decoder read the JSON and store the value to newTodo
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(newTodo.Status)>0 && !isStatusValid(newTodo.Status) {
		http.Error(w, "Bad Request - The status should only be Pending/In Progress/Compeleted", http.StatusBadRequest)
		return
	}

	todoList = append(todoList, newTodo)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New To Do item added successfully")
}

func putTodo(w http.ResponseWriter, req *http.Request, index int){
	var updatedTodo ToDo
	if err := json.NewDecoder(req.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(updatedTodo.Status)>0 && !isStatusValid(updatedTodo.Status) {
		http.Error(w, "Bad Request - The status should only be Pending/In Progress/Compeleted", http.StatusBadRequest)
		return
	}

	if !isIndexValid(index){
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	todoList[index] = updatedTodo
	fmt.Fprintf(w, "Item updated successfully")
}

func deleteTodo(w http.ResponseWriter, index int){
		if !isIndexValid(index) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		todoList = slices.Delete(todoList, index, index+1)
		fmt.Fprintf(w, "Item deleted successfully")
}