FROM golang:alpine AS builder

WORKDIR /app

ADD go.mod go.sum /app/

RUN go mod download 

COPY . . 

RUN go build -o myapp cmd/main.go

FROM alpine

WORKDIR /app

COPY .env .env

COPY --from=builder /app/myapp .

EXPOSE 8080

CMD ["./myapp"]

