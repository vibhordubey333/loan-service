FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

EXPOSE 8080

CMD ["./main"]