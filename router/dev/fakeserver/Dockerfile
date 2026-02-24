FROM golang:1.25.7-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o fakeserver ./router/dev/fakeserver

EXPOSE 8000

CMD ["./fakeserver"]