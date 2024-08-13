package main

import "sync"

func createTodoList(todos ...ToDo) []ToDo {
	return todos
}

func printTodo() {
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