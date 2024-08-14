FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum main.go ./

RUN go mod download

CMD ["go run main.go"]


