package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
)


type ToDo struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Completed *bool `json:"completed"`
}

var todoList = make([]ToDo, 0)
var pl = fmt.Println
var todoListMutex sync.Mutex

func main() {
	http.HandleFunc("/todos", getAndPostHandler)
	http.HandleFunc("/todos/", putAndDeleteHandler)

	pl("Starting server on port: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Server failed: %v\n", err)
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
	id := strings.TrimPrefix(req.URL.Path, "/todos/")
	switch req.Method {
	case http.MethodPut:
		putTodo(w, req, id)
	case http.MethodDelete:
		deleteTodo(w, id)
	default:
		http.Error(w, "Request Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTodo(w http.ResponseWriter){
	todoListMutex.Lock()
	defer todoListMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

func postTodo(w http.ResponseWriter, req *http.Request){
	todoListMutex.Lock()
	defer todoListMutex.Unlock()
	
	var newTodo ToDo

	if err := json.NewDecoder(req.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if newTodo.Id != "" || newTodo.Title == "" || newTodo.Completed == nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	newTodo.Id = uuid.New().String()

	todoList = append(todoList, newTodo)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New To Do item added successfully")
}

func putTodo(w http.ResponseWriter, req *http.Request, id string){
	todoListMutex.Lock()
	defer todoListMutex.Unlock()

	var todoToUpdate *ToDo

	for i, todo := range(todoList) {
		if todo.Id == id{
			todoToUpdate = &todoList[i]
			break
		}
	}

	if todoToUpdate == nil {
		http.Error(w, "Item not found - Please check the Id", http.StatusNotFound)
		return
	}

	var updatedTodo ToDo

	updatedTodo.Id = id

	if err := json.NewDecoder(req.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if updatedTodo.Title == "" || updatedTodo.Completed == nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	*todoToUpdate = updatedTodo
	fmt.Fprintf(w, "Item updated successfully")
}

func deleteTodo(w http.ResponseWriter, id string){
	todoListMutex.Lock()
	defer todoListMutex.Unlock()

	var indexToDelete int = -1

		for i, todo := range todoList {
			if todo.Id == id {
				indexToDelete = i
				break
			}
		}

		if indexToDelete == -1 || len(todoList) < 1  {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		todoList = slices.Delete(todoList, indexToDelete, indexToDelete+1)

		fmt.Fprintf(w, "Item deleted successfully")
}