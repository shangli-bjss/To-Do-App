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
var todoListMutex = &sync.RWMutex{}

func getTodos(w http.ResponseWriter){
	todoListMutex.RLock()
	defer todoListMutex.RUnlock()

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

func getTodoById(w http.ResponseWriter, id string){
	todoListMutex.RLock()
	defer todoListMutex.RUnlock()

	var findTodoById *ToDo

	for i, todo := range(TodoList) {
		if todo.Id == id {
			findTodoById = &TodoList[i]
			break
		}
	}

	if findTodoById == nil {
		http.Error(w, "Item not found - Please check the Id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(findTodoById)
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

	var updatedTodo ToDo = *todoToUpdate

	if err := json.NewDecoder(req.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if updatedTodo.Id != todoToUpdate.Id {
		http.Error(w, "Forbidden - the id can't be updated", http.StatusForbidden)
		return
	}

	if updatedTodo.Title == "" || updatedTodo.Completed == nil {
		http.Error(w, "Invalid JSON body - the title and completed can not be empty", http.StatusBadRequest)
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
			http.Error(w, "Item not found - Please check the Id", http.StatusNotFound)
			return
		}

		TodoList = slices.Delete(TodoList, indexToDelete, indexToDelete+1)

		fmt.Fprintf(w, "Item deleted successfully")
}
