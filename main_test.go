package main

import (
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

type request struct {
	name string
	method string
	expectedStatus int
}


func TestGetAndPostHandler(t *testing.T) {
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
				req.Body = io.NopCloser(strings.NewReader(`{"item":"test","status":"Pending"}`))
			}
			getAndPostHandler(w, req)

			if w.Code != testReq.expectedStatus {
				t.Errorf("Expected status code %d, got %d", testReq.expectedStatus, w.Code)
			}
		})
	}
}

func TestPutAndDeleteHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         string
		statusCode int
	}{
		{"PUT request", http.MethodPut, "123", http.StatusNotFound}, // TODO: not found
		{"DELETE request", http.MethodDelete, "123", http.StatusNotFound},
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
				req.Body = io.NopCloser(strings.NewReader(`{"item":"test","status":"Pending"}`))
			}
			putAndDeleteHandler(w, req)

			if w.Code != test.statusCode {
				t.Errorf("expected status code %d, got %d", test.statusCode, w.Code)
			}
		})
	}
}