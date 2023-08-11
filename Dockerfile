# Use an official Go runtime as a parent image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /usr/src/shopit

RUN go install github.com/pilu/fresh@latest

COPY . .

RUN go mod tidy