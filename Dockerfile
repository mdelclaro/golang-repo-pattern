FROM golang:1.23

WORKDIR /usr/cmd/main

RUN go install github.com/air-verse/air@latest && \
	go install go.uber.org/mock/mockgen@latest && \
	go install gotest.tools/gotestsum@latest

COPY . .
RUN go mod tidy