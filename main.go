package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)


type ToDo struct {
	Item string `json:"item"`
	Status string `json:"status"`
}

var todoList []ToDo

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
	pl("To Do item added succesfully")
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


func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		pl("---Welcome to the BEST To Do App---")
		pl("1. View To Do List")
		pl("2. Add To Do Item")
		pl("3. Update To Do Item")
		pl("4. Delete To Do Item")
		pl("5. Exit")
		pl("Enter your choice: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			pl("Invalid input, please enter a number.")
			continue
		}

		switch choice {
		case 1:
			displayTodoList()

		case 2:
			pl("Enter the new To Do item: ")
			item, _ := reader.ReadString('\n')
			item = strings.TrimSpace(item)

			pl("Enter the status (Pending/In Progress/Completed): ")
			status, _ := reader.ReadString('\n')
			status = strings.TrimSpace(status)
			if !isStatusValid(status) {
				pl("Invaild status - Please enter 'Pending/In Progress/Completed' only\n")
				continue
			}

			addTodoList(item, status)

		case 3:
			displayTodoList()
			pl("Enter the item you want to update: ")
			indexStr, _ := reader.ReadString('\n')
			indexStr = strings.TrimSpace(indexStr)

			index, err := strconv.Atoi(indexStr)
			if err != nil || !isIndexValid(index) {
				pl("Invalid input, please enter a vaild number")
				continue
			}

			pl("1. Updated the name of the item")
			pl("2. Updated the status of the item")
			pl("Enter your choice: ")

			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)
			choice, err := strconv.Atoi(choiceStr)
			if err != nil {
				pl("Invalid input, please enter a number.")
				continue
			}

			switch choice {
			case 1:
				pl("Enter the new name of the item: ")
				item, _ := reader.ReadString('\n')
				item = strings.TrimSpace(item)
				updateTodoItem(index-1, "item", item)
			case 2:
				pl("Enter the new status of the item (Pending/In Progress/Completed): ")
				status, _ := reader.ReadString('\n')
				status = strings.TrimSpace(status)
				if !isStatusValid(status) {
					pl("Invaild status - Please enter 'Pending/In Progress/Completed' only")
					continue
				}
				updateTodoItem(index-1, "status", status)
			default:
				pl("Invalid choice - please enter the number 1 or 2")
			}

		case 4:
			displayTodoList()
			pl("Enter the item you want to delete: ")
			indexStr, _ := reader.ReadString('\n')
			indexStr = strings.TrimSpace(indexStr)

			index, err := strconv.Atoi(indexStr)
			if err != nil || !isIndexValid(index) {
				pl("Invalid input, please enter a vaild number")
				continue
			}

			deleteTodoItem(index-1)

		case 5:
			return

		default:
			pl("Invalid choice, please enter a number between 1 and 5.")
		}
	}
}