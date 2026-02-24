FROM golang:1.25.7-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o router-app ./router/main.go

EXPOSE 7999

CMD ["./router-app"]