# REST API FOR TODO 

This is a simple Todo App implemented in Go using the Chi router.

## Features

* Create a new todo
* Get all todos
* Get a specific todo by ID
* Update a todo
* Delete a todo

## Installation

1. Ensure you have Go installed on your system.
2. Clone the repository: `git clone https://github.com/akhand3108/rest-api-golang.git`
3. Change to the project directory: `cd rest-api-golang`
4. Install the required dependencies: `go mod download`

## Usage

1. Start the server: `go run main.go`
2. Access the server at `http://localhost:8080` in your web browser or using an API testing tool like cURL or Postman.

## API Endpoints

* `GET /todos` - Get all todos
* `POST /todos` - Create a new todo
* `GET /todos/{id}` - Get a specific todo by ID
* `PUT /todos/{id}` - Update a todo
* `DELETE /todos/{id}` - Delete a todo

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
