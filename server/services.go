package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"todoapp/models"

	"github.com/google/uuid"
)

type ToDo = models.ToDo

var TodoList = make([]ToDo, 0)
var todoListMutex sync.Mutex

func getTodo(w http.ResponseWriter){
	todoListMutex.Lock()
	defer todoListMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TodoList)
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

	TodoList = append(TodoList, newTodo)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New To Do item added successfully")
}

func putTodo(w http.ResponseWriter, req *http.Request, id string){
	todoListMutex.Lock()
	defer todoListMutex.Unlock()

	var todoToUpdate *ToDo

	for i, todo := range(TodoList) {
		if todo.Id == id{
			todoToUpdate = &TodoList[i]
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

		for i, todo := range TodoList {
			if todo.Id == id {
				indexToDelete = i
				break
			}
		}

		if indexToDelete == -1 || len(TodoList) < 1  {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		TodoList = slices.Delete(TodoList, indexToDelete, indexToDelete+1)

		fmt.Fprintf(w, "Item deleted successfully")
}
