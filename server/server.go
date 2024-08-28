package server

import (
	"fmt"
	"log"
	"net/http"
	"todoapp/store"
)

func StartServer() {
    var err error
    st, err = store.NewTodoStore("./store/todos.db")
    if err != nil {
        log.Fatal(err)
    }
    defer st.Close()

    registerRoutes()
    fmt.Println("Starting server on port: 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v\n", err)
    }
}
