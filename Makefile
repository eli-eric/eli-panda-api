# Main go commands
SWAG := go run github.com/swaggo/swag/cmd/swag@v1.16.3

swagger:
	$(SWAG) init -g server.go
	cp -r ./docs/swagger.yaml ./open-api-specification/panda-api.yaml

run: swagger
	go run server.go	

build: swagger
	go build -v -ldflags "-s -w"

install:
	go mod download && go mod verify

db-local-up:
	docker compose -f docker-compose-local.yml up -d panda-dev-neo4j

db-local-down:
	docker compose -f docker-compose-local.yml stop panda-dev-neo4j

db-local-status:
	docker compose -f docker-compose-local.yml ps panda-dev-neo4j

# Read-only. Reports duplicate researcher records and stale current ResearcherIDs,
# both of which produce a wrong RIV delivery. Override NEO4J_PASSWORD as needed.
NEO4J_PASSWORD ?= elipanda2022
researcher-audit:
	docker exec -i panda-dev-neo4j cypher-shell -u neo4j -p '$(NEO4J_PASSWORD)' \
		--format plain < db/neo4j/reports/researcher-identity-audit.cypher

tunel:
	ssh -L 7472:127.0.0.1:7472 -L 7682:127.0.0.1:7682 -L 7471:127.0.0.1:7471 -L 7681:127.0.0.1:7681 -L 7470:127.0.0.1:7470 -L 7680:127.0.0.1:7680 -L 9000:127.0.0.1:9000 -L 9090:127.0.0.1:9090 panda@panda.eli-laser.eu

test:
	go test ./...
