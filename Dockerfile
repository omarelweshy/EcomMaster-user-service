FROM golang:1.21.4 AS builder


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o user_service ./cmd

EXPOSE 8000

CMD ["./user_service"]
