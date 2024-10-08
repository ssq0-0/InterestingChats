FROM --platform=linux/amd64 golang:1.22.2

WORKDIR /app

RUN apt-get update && apt-get install -y \
    build-essential \
    g++ \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY config.json ./

RUN go build -o DB ./cmd/main.go

CMD ["./DB"]