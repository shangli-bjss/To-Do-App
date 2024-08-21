package server

import (
	"net/http"
	"testing"
)

func TestMain(t *testing.T) {
	go StartServer()
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

