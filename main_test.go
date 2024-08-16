package main

import (
	"net/http"
	"net/http/httptest"
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
	testRequests := []struct {
		name string
		method string
		expectedStatus int
	}{
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
			getAndPostHandler(w, req)

			if w.Code != testReq.expectedStatus {
				t.Errorf("Expected status code %d, got %d", testReq.expectedStatus, w.Code)
			}
		})
	}
}
