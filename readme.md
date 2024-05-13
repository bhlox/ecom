# Minimal E-commerce Setup

A backend REST api implementation for an e-commerce platform, written in GoLang. It provides essential functionalities for user sign-ups, authentication using JSON Web Tokens (JWT), checkout processes, and retrieving orders. This project is more of a personal experience to understand core concepts of Go, particularly setting up an HTTP application.

## Tools Used

This project utilizes several tools to streamline development and testing processes. Here's a brief overview of some of the key tools used:

- [air](https://github.com/cosmtrek/air) - development live reload
- [sqlc](https://sqlc.io/) - somewhat sql query generator / ORM
- [Goose](https://github.com/pressly/goose) - manage db schemas
- [Testcontainers](https://testcontainers.com/)

## Testcontainers

Testcontainers is integrated to set up and tear down database containers for our integration tests. This ensures that each test runs in a clean environment, preventing state leakage between tests and making our tests more reliable and easier to debug.

Currently a bug is still present when executing the test command for this project.
If you are encountering this along the test. Just keep testing till it passses (I know it's stupid).
issue -> [Create Reaper](https://github.com/testcontainers/testcontainers-go/issues/2172)

```
20xx/xx/xx xx:xx:xx 🔥 Reaper obtained from Docker for this test session b7299efdc4dd5a3319a6beafba78c721cc497843403fcde54323ff283352933b
routes_test.go:52: port not found: creating reaper failed: failed to create container
```

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
