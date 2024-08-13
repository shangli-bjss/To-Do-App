package main

import (
	"fmt"
	"sync"
)


type ToDo struct {
	Item string `json:"item"`
	Status string `json:"status"`
}

var pl = fmt.Println

func createTodoList(todos ...ToDo) []ToDo {
	return todos
}

// func displayTodoList(todoList []ToDo) {
// 	for i, todo := range todoList {
// 		pl(i+1, ": ", todo.Item, " - ", todo.Status)
// 	}
// }

func main() {
	todoList := createTodoList(
	ToDo{"Buy groceries", "Pending"},
	ToDo{"Clean the house", "In Progress"},
	ToDo{"Write blog post", "Completed"},
	ToDo{"Finish project", "Pending"},
	ToDo{"Call mom", "Completed"},
	ToDo{"Attend meeting", "Pending"},
	ToDo{"Read book", "In Progress"},
	ToDo{"Exercise", "Pending"},
	ToDo{"Plan vacation", "Completed"},
	ToDo{"Pay bills", "In Progress"},
)
	var wg sync.WaitGroup
	var mu sync.Mutex

	done := make(chan bool)

	wg.Add(2)

	// print to do items
	go func() {
		defer wg.Done()
		for i:=0; i<len(todoList); i++ {
			mu.Lock()
			pl("To Do: ", todoList[i].Item)
			mu.Unlock()
			done <- true
			<-done // wait for another routine to finish
		}
	}()

	// print the status
	go func() {
		defer wg.Done()
		for i:=0; i<len(todoList); i++ {
			<-done
			mu.Lock()
			pl("Status: ", todoList[i].Status)
			mu.Unlock()
			done <- true
		}
	}()

	wg.Wait()
	close(done)

	pl("All routines completed")
}

// func main() {
// 	todoList, err := readFromJsonFile("todo_list.json")
// 	if err != nil {
// 		pl("Error reading from JSON file", err)
// 		return
// 	}

// 	displayTodoList(todoList)
// }


// func main() {
// 	todoList := createTodoList(
// 	ToDo{"Buy groceries", "Pending"},
// 	ToDo{"Clean the house", "In Progress"},
// 	ToDo{"Write blog post", "Completed"},
// 	ToDo{"Finish project", "Pending"},
// 	ToDo{"Call mom", "Completed"},
// 	ToDo{"Attend meeting", "Pending"},
// 	ToDo{"Read book", "In Progress"},
// 	ToDo{"Exercise", "Pending"},
// 	ToDo{"Plan vacation", "Completed"},
// 	ToDo{"Pay bills", "In Progress"},
// )

// err := writeToJsonFile("todo_list.json", todoList)
// if err != nil {
// 	pl("Error writing to JSON file", err)
// 	return
// }

// pl("ToDo list has been written to todo_list.json successfully")

// }