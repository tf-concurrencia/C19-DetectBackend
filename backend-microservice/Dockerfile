FROM golang:alpine AS build

WORKDIR /app
COPY . .
RUN go build ./backend-service.go

ENTRYPOINT ["/app/backend-service"][]
