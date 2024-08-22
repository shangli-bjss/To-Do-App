# To-Do App

A simple To-Do App built using Go.

## Features

* List all To-Do items
* Create new To-Do items
* Update existing To-Do items
* Delete To-Do items


## API Endpoints

* `GET /todos`: Get all To-Do items
* `POST /todos`: Create a new To-Do item
* `PUT /todos/{id}`: Update a To-Do item
* `DELETE /todos/{id}`: Delete a To-Do item

## Running the App

To run the app, navigate to the project directory and execute the following command:

```bash
go run main.go
```

This will start the server, and you can access the API endpoints using a tool like curl or a REST client (e.g. Postman).