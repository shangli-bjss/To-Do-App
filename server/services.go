package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todoapp/models"
	"todoapp/store"
)

type ToDo = models.ToDo

var st *store.TodoStore

func getTodos(w http.ResponseWriter){
	todos, err := st.GetAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func postTodo(w http.ResponseWriter, req *http.Request){
	var newTodo ToDo

	if err := json.NewDecoder(req.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if newTodo.Id != "" || newTodo.Title == "" || newTodo.Completed == nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err := st.CreateTodo(newTodo.Title, newTodo.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New To Do item added successfully")
}

func getTodoById(w http.ResponseWriter, id string){
	todo, err := st.GetTodo(id)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func putTodo(w http.ResponseWriter, req *http.Request, id string){
	todo, err := st.GetTodo(id)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var updatedTodo ToDo = todo

	if err := json.NewDecoder(req.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if updatedTodo.Id != todo.Id {
		http.Error(w, "Forbidden - the id can't be updated", http.StatusForbidden)
		return
	}

	if updatedTodo.Title == "" || updatedTodo.Completed == nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	err = st.UpdateTodo(id, updatedTodo.Title, *updatedTodo.Completed)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Item updated successfully")
}

func deleteTodo(w http.ResponseWriter, id string){
	err := st.DeleteTodo(id)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
		fmt.Fprintf(w, "Item deleted successfully")
}
