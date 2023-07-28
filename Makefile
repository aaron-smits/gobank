build: 
	@go build -o bin/gobank

run: build
	docker compose up -d
ifdef attach
	docker compose up
endif

stop:

ifdef clean
	docker compose down -v
	docker volume prune
endif
	docker compose down