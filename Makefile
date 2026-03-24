NAME = configurator-api

.PHONY: build run dev clean

build:
	go build -o $(NAME) .

run:
	build ./$(NAME)

dev:
	go run .

clean:
	rm $(APP_NAME)
