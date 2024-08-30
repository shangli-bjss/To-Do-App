package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testId = "testId"

func TestGetTodo(t *testing.T) {
	TodoList = []ToDo{}

	TodoList = append(TodoList, ToDo{
		Id: testId,
		Title: "Test Todo",
		Completed: new(bool),
	}, ToDo{
		Id: "123",
		Title: "read book",
		Completed: new(bool),
	})

	w := httptest.NewRecorder()

	getTodos(w)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var response []ToDo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 2 {
		t.Errorf("expected 1 todo but got %d", len(response))
	}

	if response[0].Id != testId {
		t.Errorf("expectd id is 'test-id', but got %s", response[0].Id)
	}
}

func TestGetTodoById(t *testing.T) {
	TodoList = []ToDo{}

	TodoList = append(TodoList, ToDo{
		Id: testId,
		Title: "Test Todo",
		Completed: new(bool),
	}, ToDo{
		Id: "123",
		Title: "read book",
		Completed: new(bool),
	})


	w := httptest.NewRecorder()
	getTodoById(w, testId)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var response ToDo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Id != testId {
		t.Errorf("expectd id is 'test-id', but got %s", response.Id)
	}
}

func TestPostTodo(t *testing.T) {
	TodoList = []ToDo{}

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
	postTodo(w, req)

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
	TodoList = []ToDo{}

	newTodo := ToDo{
		Title: "",
		Completed: nil,
	}

	body, _ := json.Marshal(newTodo)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	postTodo(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	if len(TodoList) != 0 {
		t.Errorf("expected got 0 todo item, but got %d", len(TodoList))
	}
}

func TestPutTodo(t *testing.T) {
	TodoList = []ToDo{}

	TodoList = append(TodoList, ToDo{
		Id: testId,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := ToDo{
		Id: testId,
		Title: "updated todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+testId, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	putTodo(w, req, testId)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if TodoList[0].Title != "updated todo" {
		t.Errorf("expected todo title is 'updated todo', but got %s", TodoList[0].Title)
	}
}

func TestPutTodoInvalidId(t *testing.T) {
	TodoList = []ToDo{}
	invalidId := "invalidId"

	TodoList = append(TodoList, ToDo{
		Id: testId,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := ToDo{
		Title: "updated todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+invalidId, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	putTodo(w, req, invalidId)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, but got %d", http.StatusNotFound, w.Code)
	}
}

func TestPutTodoWithId(t *testing.T) {
	TodoList = []ToDo{}
	newId := "newId"

	TodoList = append(TodoList, ToDo{
		Id: testId,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := ToDo{
		Id: newId,
		Title: "update todo",
		Completed: new(bool),
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+testId, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	putTodo(w, req, testId)

	if status := w.Code; status != http.StatusForbidden {
		t.Errorf("expected status code %d, but got %d", http.StatusForbidden, w.Code)
	}
}

func TestPutTodoInvalidJSONBody(t *testing.T) {
	TodoList = []ToDo{}

	TodoList = append(TodoList, ToDo{
		Id: testId,
		Title: "test todo",
		Completed: new(bool),
	})

	updatedTodo := ToDo{
		Id: testId,
		Title: "",
		Completed: nil,
	}

	body, _ := json.Marshal(updatedTodo)
	req, err := http.NewRequest("PUT", "/todos/"+testId, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	putTodo(w, req, testId)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteTodo(t *testing.T) {
	TodoList = []ToDo{}

	TodoList = append(TodoList, ToDo{
		Id: testId,
		Title: "test todo",
		Completed: new(bool),
	})

	w := httptest.NewRecorder()
	deleteTodo(w, testId)
	
	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if len(TodoList) != 0 {
		t.Errorf("expected todolist to be empty, but got %d items", len(TodoList))
	}
}