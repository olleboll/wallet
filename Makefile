build:
	docker build . -t wallet

up: up-deps
	sleep 2
	docker-compose up -d api-wallet

up-deps:
	docker-compose up -d db-wallet

down:
	docker-compose down

test: build up
	sleep 1
	go test ./tests
	docker-compose down