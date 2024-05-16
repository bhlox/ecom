# syntax=docker/dockerfile:1

FROM golang:1.22.2 AS build-stage
  WORKDIR /app

  COPY go.mod go.sum ./
  RUN go mod download

  COPY . .

  RUN apt-get update && apt-get install -y ca-certificates
  RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

  # Run the tests in the container
# FROM build-stage AS run-test-stage
#   RUN go test -v ./...

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
  WORKDIR /

  COPY --from=build-stage /api /api

  COPY --from=build-stage /etc/ssl/certs /etc/ssl/certs

  EXPOSE 8080

  ENTRYPOINT ["/api"]