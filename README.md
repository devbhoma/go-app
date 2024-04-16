# go-app

This project is built in Golang following a microservice architecture. The folder structure is organized as follows:

```
go-app
├── cmd
│   ├── admin
│   ├── apiserver
│   ├── cli
│   └── wsserver
├── config
├── internal
│   ├── authorization
│   ├── database
│   ├── endpoints
│   ├── httpserver
│   ├── store
│   └── utils
└── main.go
```

# Directory Explanations:

`cmd`: This folder contains microservice modules and HTTP server handlers for the REST API and WebSocket server.

`admin`: Used for PostgreSQL database migrations.

`apiserver`: Implements REST API server functionality.

`cli`: command-line setup using the Cobra library.

`wsserver`: Handles WebSocket server setup.

`config`: Holds project configuration file data, synchronized with a .env file.

`internal`: Contains common components dynamically utilized by cmd microservices.

`authorization`: Provides common handlers for the authentication system.

`database`: Manages database operations.

`endpoints`: Holds business logic code for all modules, facilitating centralized handling.

`httpserver`: Utilizes the Gin library to create an API server, used in cmd `apiserver` & `wsserver`.

`store`: Manages database collections and entities.