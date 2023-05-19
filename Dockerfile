FROM golang:latest as builder

WORKDIR /app

COPY . .

EXPOSE 8000

RUN GOOS=linux go build -o server ./cmd/server/

CMD ["./server"]