#  REST API with JWT Authentication and Todo Management

This is a simple RESTful API implemented in Golang that provides authentication functionality along with CRUD operations for managing todos. It uses SQLlite as the database and JWT (JSON Web Tokens) for authentication.

## Features

* Create a new todo
* Get all todos
* Get a specific todo by ID
* Update a todo
* Delete a todo
* Signup User
* Signin User

## Installation

1. Ensure you have Go installed on your system.
2. Clone the repository: 
```bash
git clone https://github.com/akhand3108/rest-api-golang.git
```
3. Change to the project directory: 
```bash
cd rest-api-golang
```
4. Install the required dependencies: 
```bash
go mod download
```

5. Set up the environment variables by creating a `.env` file in the project root directory. Use the provided `.env.template` file as a reference.
6. Build the app: 
```bash
go build .
```
7. Run the App
```bash
./restgo
```

## Usage

Access the server at `http://localhost:8080` in your web browser or using an API testing tool like cURL or Postman.

## API Endpoints

The following API endpoints are available:

* `GET /todos`: Get all todos
* `POST /todos`: Create a new todo
* `GET /todos/{id}`: Get a specific todo by ID
* `PUT /todos/{id}`: Update a todo
* `DELETE /todos/{id}`: Delete a todo
* `POST /signup`: Create a new user account
* `POST /signin`: Sign in to an existing user account


## API Documentation

The following endpoints are available in the API:

### Sign Up

Create a new user account.

- **URL**: `/signup`
- **Method**: `POST`

### Sign In

Sign in to an existing user account.

- **URL**: `/signin`
- **Method**: `POST`

### Get All Todos

Retrieve all todos for the authenticated user.

- **URL**: `/todos`
- **Method**: `GET`

### Create Todo

Create a new todo for the authenticated user.

- **URL**: `/todos`
- **Method**: `POST`

### Get Todo by ID

Retrieve a specific todo by its ID for the authenticated user.

- **URL**: `/todos/{id}`
- **Method**: `GET`

### Update Todo by ID

Update a specific todo by its ID for the authenticated user.

- **URL**: `/todos/{id}`
- **Method**: `PUT`

### Delete Todo by ID

Delete a specific todo by its ID for the authenticated user.

- **URL**: `/todos/{id}`
- **Method**: `DELETE`

Note: The authenticated routes require a bearer token in the request headers. Please include the token in the `Authorization` header as `Bearer <token>`.



## Dependencies

The app depends on the following libraries:

* `chi`: A web framework for Golang.
* `jwt-go`: JWT in Golang

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).




