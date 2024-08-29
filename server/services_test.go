package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTodo(t *testing.T) {
	TodoList = []Todo{}

	TodoList = append(TodoList, Todo{
		Id: "test-id",
		Title: "Test Todo",
		Completed: new(bool),
	}, Todo{
		Id: "123",
		Title: "read book",
		Completed: new(bool),
	})

	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var response []Todo
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 2 {
		fmt.Println("todolist on test", TodoList)
		t.Errorf("expected 1 todo but got %d", len(response))
	}

	if response[0].Id != "test-id" {
		t.Errorf("expectd id is 'test-id', but got %s", response[0].Id)
	}
}

func TestGetTodoById(t *testing.T) {
	TodoList = []Todo{}

	TodoList = append(TodoList, Todo{
		Id: "test-id",
		Title: "Test Todo",
		Completed: new(bool),
	}, Todo{
		Id: "123",
		Title: "read book",
		Completed: new(bool),
	})

	req, err := http.NewRequest("GET", "/todos/test-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosByIdHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var response Todo
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Id != "test-id" {
		t.Errorf("expectd id is 'test-id', but got %s", response.Id)
	}
}

func TestPostTodo(t *testing.T) {
	TodoList = []Todo{}

	newTodo := Todo{
		Title: "new todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(newTodo)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	if len(TodoList) != 1 {
		t.Errorf("expected got 1 todo item, but got %d", len(TodoList))
	}

	if TodoList[0].Title != "new todo" {
		t.Errorf("expected new todo title is 'new todo', but got %s", TodoList[0].Title)
	}
}

func TestPostTodoInvalidJSONBody(t *testing.T) {
	TodoList = []Todo{}

	newTodo := Todo{
		Title: "",
		Completed: nil,
	}

	body, _ := json.Marshal(newTodo)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	if len(TodoList) != 0 {
		t.Errorf("expected got 0 todo item, but got %d", len(TodoList))
	}
}

func TestPutTodo(t *testing.T) {
	TodoList = []Todo{}

	testid := "test-id"
	TodoList = append(TodoList, Todo{
		Id: testid,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := Todo{
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
	handler := http.HandlerFunc(todosByIdHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if TodoList[0].Title != "updated todo" {
		t.Errorf("expected todo title is 'updated todo', but got %s", TodoList[0].Title)
	}
}

func TestPutTodoInvalidId(t *testing.T) {
	TodoList = []Todo{}
	invalidId := "invalidId"

	TodoList = append(TodoList, Todo{
		Id: "testid",
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := Todo{
		Id: invalidId,
		Title: "updated todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+invalidId, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosByIdHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, but got %d", http.StatusNotFound, w.Code)
	}
}

func TestPutTodoWithId(t *testing.T) {
	TodoList = []Todo{}
	testId := "testid"

	TodoList = append(TodoList, Todo{
		Id: testId,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := Todo{
		Id: "newid",
		Title: "update todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+testId, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosByIdHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusForbidden {
		t.Errorf("expected status code %d, but got %d", http.StatusForbidden, w.Code)
	}
}

func TestPutTodoInvalidJSONBody(t *testing.T) {
	TodoList = []Todo{}

	testid := "test-id"
	TodoList = append(TodoList, Todo{
		Id: testid,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := Todo{
		Id: testid,
		Title: "",
		Completed: nil,
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+testid, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosByIdHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteTodo(t *testing.T) {
	TodoList = []Todo{}

	testid := "test-id"
	TodoList = append(TodoList, Todo{
		Id: testid,
		Title: "test todo",
		Completed: new(bool),
	})

	req, err := http.NewRequest("DELETE", "/todos/"+testid, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(todosByIdHandler)
	
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if len(TodoList) != 0 {
		t.Errorf("expected todolist to be empty, but got %d items", len(TodoList))
	}
}