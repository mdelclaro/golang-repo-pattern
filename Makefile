test: 
	go test ./src/api/handlers/...

build-dependencies:
	docker compose build
	
run:
	docker compose up