package main

import (
	"fmt"
	"slices"
)

var pl = fmt.Println

func displayTodoList() {
	if len(todoList) < 1 {
		pl("Todo list is empty\n")
		return
	}
	for i, todo := range todoList {
		pl(i+1, ": ", todo.Item, " - ", todo.Status)
	}
}

func addTodoList(item, status string) {
	todoList = append(todoList, ToDo{Item: item, Status: status})
	pl("To Do item added succesfully\n")
}

func isIndexValid(index int) bool {
	if index < 0 || index > len(todoList) {
		return false
	}
	return true
}

func isStatusValid(status string) bool {
	if status != "Pending" && status != "In Progress" && status !="Completed" {
		return false
	}
	return true
}

func updateTodoItem(index int, field, value string) {
	if !isIndexValid(index) {
		pl("The index is invalid")
		return
	}
	switch field {
	case "item":
		todoList[index].Item = value
	case "status":
		todoList[index].Status = value
	default:
		pl("Invalid field")
		return
	}
	fmt.Printf("Item updated successfully\n")
}

func deleteTodoItem(index int){
	if !isIndexValid(index) {
		pl("The index is invalid")
		return
	}
	todoList = slices.Delete(todoList, index, index+1)
	fmt.Printf("Item deleted successfully\n")
}