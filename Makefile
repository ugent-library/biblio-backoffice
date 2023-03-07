export GO111MODULE=on

-include .env
BIBLIO_BACKOFFICE_PG_CONN=postgres://biblio:biblio@localhost:2345/biblio
BIBLIO_BACKOFFICE_ES6_URL=http://localhost:9400
BIBLIO_BACKOFFICE_DATASET_INDEX=biblio_datasets
BIBLIO_BACKOFFICE_PUBLICATION_INDEX=biblio_publications
BIBLIO_BACKOFFICE_PORT=3999
BIBLIO_BACKOFFICE_MONGODB_URL=mongodb://localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000
BIBLIO_BACKOFFICE_FRONTEND_ES6_URL=http://localhost:9200
export

setup-test-env: init-test-env create-db-tables create-es-indices

tear-test-env:
	docker-compose down

run-tests:
	@cd client && go clean -testcache ./... && go test ./... -v -cover

coverage:
	@cd client && go clean -testcache ./... && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out

init-test-env:
	docker-compose up -d
	@echo "> Waiting 45 seconds for ElasticSearch to fully boot..."
	@sleep 45
	curl -XPUT -H "Content-Type: application/json" http://localhost:9400/_cluster/settings -d '{ "transient": { "cluster.routing.allocation.disk.threshold_enabled": false } }'

create-es-indices:
	go run main.go publication reindex
	go run main.go dataset reindex
	curl -XPUT -H "Content-Type: application/json" http://localhost:9400/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}'

create-db-tables:
	go install github.com/jackc/tern@latest
	tern migrate --migrations etc/snapstore/migrations

api-server:
	go run main.go api start


# Declare all targets as "PHONY", see https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html.
MAKEFLAGS += --always-make
.PHONY: init-test-env create-es-indices create-db-tables api-server setup-test-env tear-test-env run-test-grpc-client