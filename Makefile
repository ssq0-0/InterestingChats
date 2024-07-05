BINARY_NAME=Chats
CMD_DIR=./cmd/server

.PHONY: all build run deps clean

all: build

build:
	go build -o ${BINARY_NAME} ${CMD_DIR}

run: 
	go run ${CMD_DIR}

deps:
	go mod tidy

clean:
	rm -f ${BINARY_NAME}
