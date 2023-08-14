# Use an official Go runtime as a parent image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /usr/src/backend

RUN go install github.com/pilu/fresh@latest

COPY . .

# Update deps
RUN go mod tidy

# Build executeable
RUN GOOS=linux CGO_ENABLED=0 go build -o /main

CMD [ "/main" ]