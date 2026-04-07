FROM golang:1.25.7-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bench-app ./bench/

CMD ["./bench-app"]
