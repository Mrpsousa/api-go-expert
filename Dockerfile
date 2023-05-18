FROM golang:latest

WORKDIR /app

COPY . .

EXPOSE 8000

RUN go mod download

CMD ["go", "run", "cmd/server/main.go"]