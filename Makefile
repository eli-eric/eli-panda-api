# Main go commands
swagger:
	swag init -g server.go
	cp -r ./docs/swagger.yaml ./open-api-specification/panda-api.yaml

run: swagger
	go run server.go	

build: swagger
	go build -v -ldflags "-s -w"

install:
	go mod download && go mod verify

db-local-up:
	docker-compose -f docker/docker-compose-databases-local.yml up -d

db-local-down:
	docker-compose -f docker/docker-compose-databases-local.yml down

tunel:
	ssh -L 7472:127.0.0.1:7472 -L 7682:127.0.0.1:7682 -L 7471:127.0.0.1:7471 -L 7681:127.0.0.1:7681 -L 7470:127.0.0.1:7470 -L 7680:127.0.0.1:7680 -L 9000:127.0.0.1:9000 -L 9090:127.0.0.1:9090 panda@panda.eli-laser.eu

test:
	go test ./...