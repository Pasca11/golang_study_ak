FROM golang:alpine

WORKDIR /usr/local/src

COPY . .

RUN go mod download

RUN go build -o app .

EXPOSE 4444

CMD ["./app"]


