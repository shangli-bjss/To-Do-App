package main

import (
	"encoding/json"
	"os"
)

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