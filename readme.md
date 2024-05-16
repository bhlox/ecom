![code coverage badge](https://github.com/bhlox/ecom/actions/workflows/tests.yml/badge.svg)

# Minimal E-commerce Setup

A backend REST api implementation for an e-commerce platform, written in GoLang. It provides essential functionalities for user sign-ups, authentication using JSON Web Tokens (JWT), checkout processes, and retrieving orders. This project is more of a personal experience to understand core concepts of Go, particularly setting up an HTTP application.

## Tools Used

This project utilizes several tools to streamline development and testing processes. Here's a brief overview of some of the key tools used:

- [air](https://github.com/cosmtrek/air) - development live reload
- [sqlc](https://sqlc.io/) - somewhat sql query generator / ORM
- [Goose](https://github.com/pressly/goose) - manage db schemas
- [DockerTest](https://github.com/ory/dockertest)

### Prerequisites

- Go
- Docker

### Setting Up the Environment

Provide the following variables to your `.env`. the example variables are provided in the `.env.example` file.

### RUNNING THE PROJECT

Side note: if you have `task` installed on your machine, a `taskfile.yml` is provided containing the commands to run and test the project.

```
    go run cmd/main.go
```

### RUNNING TESTS

```
    go test -v ./...
```

### DOCKER

```
    docker build -t image-name .
    docker run --env-file=.env -p 8080:8080 image-name
```

-This is currently running on GCP
