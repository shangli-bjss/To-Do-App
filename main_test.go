package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	go main()
	time.Sleep(1 * time.Second)
	resp, err := http.Get("http://localhost:8080/todos")
	if err != nil {
		t.Errorf("Server failed to start %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Sever failed to start %v", resp.Status)
	}

	if err := http.ListenAndServe(":8080", nil); err == nil {
		t.Errorf("Server started with port already in use")
	}
}

func TestGetAndPostHandler(t *testing.T) {
	type request struct {
		name string
		method string
		expectedStatus int
	}
	testRequests := []request{
		{"GET request", http.MethodGet, http.StatusOK},
		{"POST request", http.MethodPost, http.StatusCreated},
		{"PUT request", http.MethodPut, http.StatusMethodNotAllowed},
		{"DELETE request", http.MethodDelete, http.StatusMethodNotAllowed},
	}
	
	for _, testReq := range(testRequests) {
		t.Run(testReq.name, func(t *testing.T) {
			req, err := http.NewRequest(testReq.method, "/todos", nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			if req.Method == http.MethodPost {
				req.Body = io.NopCloser(strings.NewReader(`{"title":"test","completed":false}`))
			}

			getAndPostHandler(w, req)

			if w.Code != testReq.expectedStatus {
				t.Errorf("Expected status code %d, got %d", testReq.expectedStatus, w.Code)
			}
		})
	}
}

func TestPutAndDeleteHandler(t *testing.T) {
	todoList = []ToDo{}

	todoList = append(todoList, ToDo{
		Id:        "test-id",
		Title:     "Test Todo",
		Completed: new(bool),
	})

	tests := []struct {
		name       string
		method     string
		id         string
		statusCode int
	}{
		{"PUT request", http.MethodPut, "test-id", http.StatusOK},
		{"DELETE request", http.MethodDelete, "test-id", http.StatusOK},
		{"GET request", http.MethodGet, "", http.StatusMethodNotAllowed},
		{"POST request", http.MethodPost, "", http.StatusMethodNotAllowed},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, "/todos/"+test.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			if req.Method == http.MethodPut {
				req.Body = io.NopCloser(strings.NewReader(`{"title":"test","completed":false}`))
			}
			putAndDeleteHandler(w, req)

			if w.Code != test.statusCode {
				t.Errorf("expected status code %d, got %d", test.statusCode, w.Code)
			}
		})
	}
}

func TestGetTodo(t *testing.T) {
	todoList = []ToDo{}

	todoList = append(todoList, ToDo{
		Id:        "test-id",
		Title:     "Test Todo",
		Completed: new(bool),
	})

	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(getAndPostHandler)

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
		t.Errorf("expected 1 todo but got %d", len(response))
	}

	if response[0].Id != "test-id" {
		t.Errorf("expectd id is 'test-id', but got %s", response[0].Id)
	}
}

func TestPostTodo(t *testing.T) {
	todoList = []ToDo{}

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
	handler := http.HandlerFunc(getAndPostHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	if len(todoList) != 1 {
		t.Errorf("expected got 1 todo item, but got %d", len(todoList))
	}

	if todoList[0].Title != "new todo" {
		t.Errorf("expected new todo title is 'new todo', but got %s", todoList[0].Title)
	}
}

func TestPutTodo(t *testing.T) {
	todoList = []ToDo{}

	testid := "test-id"
	todoList = append(todoList, ToDo{
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
	handler := http.HandlerFunc(putAndDeleteHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if todoList[0].Title != "updated todo" {
		t.Errorf("expected todo title is 'updated todo', but got %s", todoList[0].Title)
	}
}

func TestDeleteTodo(t *testing.T) {
	todoList = []ToDo{}

	testid := "test-id"
	todoList = append(todoList, ToDo{
		Id: testid,
		Title: "test todo",
		Completed: new(bool),
	})

	req, err := http.NewRequest("DELETE", "/todos/"+testid, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(putAndDeleteHandler)
	
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if len(todoList) != 0 {
		t.Errorf("expected todolist to be empty, but got %d items", len(todoList))
	}
}