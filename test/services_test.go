package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"todoapp/server"
)

func TestGetTodo(t *testing.T) {
	server.TodoList = []ToDo{}

	server.TodoList = append(server.TodoList, ToDo{
		Id:        "test-id",
		Title:     "Test Todo",
		Completed: new(bool),
	})

	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetAndPostHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var response []ToDo
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 1 {
		fmt.Println("todolist on test", server.TodoList)
		t.Errorf("expected 1 todo but got %d", len(response))
	}

	if response[0].Id != "test-id" {
		t.Errorf("expectd id is 'test-id', but got %s", response[0].Id)
	}
}

func TestPostTodo(t *testing.T) {
	server.TodoList = []ToDo{}

	newTodo := ToDo{
		Title: "new todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(newTodo)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetAndPostHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	if len(server.TodoList) != 1 {
		t.Errorf("expected got 1 todo item, but got %d", len(server.TodoList))
	}

	if server.TodoList[0].Title != "new todo" {
		t.Errorf("expected new todo title is 'new todo', but got %s", server.TodoList[0].Title)
	}
}

func TestPutTodo(t *testing.T) {
	server.TodoList = []ToDo{}

	testid := "test-id"
	server.TodoList = append(server.TodoList, ToDo{
		Id: testid,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := ToDo{
		Id: testid,
		Title: "updated todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+testid, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(server.PutAndDeleteHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if server.TodoList[0].Title != "updated todo" {
		t.Errorf("expected todo title is 'updated todo', but got %s", server.TodoList[0].Title)
	}
}

func TestDeleteTodo(t *testing.T) {
	server.TodoList = []ToDo{}

	testid := "test-id"
	server.TodoList = append(server.TodoList, ToDo{
		Id: testid,
		Title: "test todo",
		Completed: new(bool),
	})

	req, err := http.NewRequest("DELETE", "/todos/"+testid, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(server.PutAndDeleteHandler)
	
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if len(server.TodoList) != 0 {
		t.Errorf("expected todolist to be empty, but got %d items", len(server.TodoList))
	}
}