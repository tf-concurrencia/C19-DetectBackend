FROM golang:alpine AS build

WORKDIR /app
COPY . .
RUN go build ./cmd/main.go

ENTRYPOINT [ "/app/main" ][]