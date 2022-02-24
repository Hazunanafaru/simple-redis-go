# Simple Redis with Go

## Description
Try to implement simple Redis command like SET and GET with Go programming language.

## Requirements
1. Go installed
2. Docker

## Usage
1. Pull redis docker image from docker hub with this command

    `docker pull redis:6.2-alpine`
2. Run the image in container

    `docker run -d -p 6379:6379 redis:6.2-alpine`
3. Get needed go external package
    
    `go mod tidy`
4. Run the `main.go` file

    `go run main.go`
4. See the result

## References
[1] ["Redis Docker Hub"](https://hub.docker.com/_/redis)