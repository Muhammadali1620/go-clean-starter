# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./migrations ./migrations
COPY ./docs ./docs
COPY ./env ./env

RUN CGO_ENABLED=0 GOOS=linux go build -o new_project_app ./cmd/app/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

ENV CONTAINER_MODE=1

COPY --from=builder /app/new_project_app .

EXPOSE 8080

CMD ["./new_project_app"]
