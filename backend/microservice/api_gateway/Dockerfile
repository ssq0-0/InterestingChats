FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY config.json .

RUN go build -o AG ./cmd/main.go

CMD ["./AG"]