NAME = configurator-api

.PHONY: build run dev clean

build:
	go build -o $(NAME) .

run:
	./$(NAME)

dev:
	go run .

clean:
	rm -f $(NAME)
