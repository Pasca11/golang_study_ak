FROM golang:1.21-alpine

WORKDIR /usr/local/src

COPY ./go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

CMD ["./app"]



