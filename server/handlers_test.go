package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	TodoList = []ToDo{}

	TodoList = append(TodoList, ToDo{
		Id:        "test-id",
		Title:     "Test Todo",
		Completed: new(bool),
	})

	tests := []struct {
		name       string
		method     string
		id         string
		expectedStatus int
	}{
		{"PUT request", http.MethodPut, "test-id", http.StatusOK},
		{"DELETE request", http.MethodDelete, "test-id", http.StatusOK},
		{"GET request", http.MethodGet, "test-id", http.StatusMethodNotAllowed},
		{"POST request", http.MethodPost, "test-id", http.StatusMethodNotAllowed},
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

			putAndDeleteHandler(w, req)

			if w.Code != test.expectedStatus {
				fmt.Println(req.URL)
				fmt.Println(TodoList)
				t.Errorf("expected status code %d, got %d", test.expectedStatus, w.Code)
			}
		})
	}
}

