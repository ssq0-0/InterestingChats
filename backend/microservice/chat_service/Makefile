.PHONY: up build down

IMAGE_NAME=cs

up:
	docker compose up --build

build:
	docker build -t $(IMAGE_NAME) .

down:
	docker compose down
