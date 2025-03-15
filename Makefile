
build:
	docker-compose --file deployments/docker/docker-compose.yml down
	docker-compose --file deployments/docker/docker-compose.yml up --build

run:
	go run ./cmd/server/main.go

generate:
	go generate api/generate.go