package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTodosHandler(t *testing.T) {
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

			todosHandler(w, req)

			if w.Code != testReq.expectedStatus {
				t.Errorf("Expected status code %d, got %d", testReq.expectedStatus, w.Code)
			}
		})
	}
}

func TestTodosByIdHandle(t *testing.T) {
	TodoList = []Todo{}
	mockid := "mockid"

	TodoList = append(TodoList, Todo{
		Id:        mockid,
		Title:     "Test Todo",
		Completed: new(bool),
	})

	tests := []struct {
		name       string
		method     string
		id         string
		expectedStatus int
	}{
		{"GET request", http.MethodGet, mockid, http.StatusOK},
		{"POST request", http.MethodPost, mockid, http.StatusMethodNotAllowed},
		{"PUT request invalid id", http.MethodPut, "", http.StatusBadRequest},
		{"DELETE request invalid id", http.MethodDelete, "", http.StatusBadRequest},
		{"PUT request", http.MethodPut, mockid, http.StatusOK},
		{"DELETE request", http.MethodDelete, mockid, http.StatusOK},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, "/todos/" + test.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			if req.Method == http.MethodPut {
				req.Body = io.NopCloser(strings.NewReader(`{"title":"updated test","completed":true}`))
			}

			todosByIdHandler(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("expected status code %d, got %d", test.expectedStatus, w.Code)
			}
		})
	}
}


