package server

import (
    "fmt"
    "log"
    "net/http"
)

func StartServer() {
    registerRoutes()
    fmt.Println("Starting server on port: 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v\n", err)
    }
}
