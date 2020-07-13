FROM golang:latest
WORKDIR /app
ADD . /app
CMD go run main/search.go