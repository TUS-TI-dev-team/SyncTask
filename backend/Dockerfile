FROM golang:1.26-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go install github.com/air-verse/air@latest

CMD ["air"]
