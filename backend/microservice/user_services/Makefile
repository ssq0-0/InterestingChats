BINARY_NAME=UserService
CMD_DIR=./cmd

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
