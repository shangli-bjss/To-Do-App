package main

import (
	"encoding/json"
	"fmt"
	"os"
)


type ToDo struct {
	Item string `json:"item"`
	Status string `json:"status"`
}

var pl = fmt.Println

func createTodoList(todos ...ToDo) []ToDo {
	return todos
}

func writeToJsonFile(filename string, todoList []ToDo) error {
	file, error := os.Create(filename)
	if error != nil {
		return error
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")

	err := encoder.Encode(todoList)
	if err != nil {
		return err
	}

	return nil
}


func readFromJsonFile(filename string) ([]ToDo, error) {
	var todoList []ToDo
	file, error := os.Open(filename)
	if error != nil {
		return nil, error
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err := decoder.Decode(&todoList)
	if err != nil {
		return nil, err
	}

	return todoList, nil
}

func displayTodoList(todoList []ToDo) {
	for i, todo := range todoList {
		pl(i+1, ": ", todo.Item, " - ", todo.Status)
	}
}

func main() {
	todoList, err := readFromJsonFile("todo_list.json")
	if err != nil {
		pl("Error reading from JSON file", err)
		return
	}

	displayTodoList(todoList)
}


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