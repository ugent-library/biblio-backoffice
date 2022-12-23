export GO111MODULE=on

-include .env
BIBLIO_BACKEND_PG_CONN=postgres://biblio:biblio@localhost:2345/biblio
BIBLIO_BACKEND_ES6_URL=http://localhost:9400
BIBLIO_BACKEND_DATASET_INDEX=biblio_datasets
BIBLIO_BACKEND_PUBLICATION_INDEX=biblio_publications
BIBLIO_BACKEND_PORT=3999
export

setup-test-env: init-test-env create-db-tables create-es-indices

tear-test-env:
	docker-compose down

run-tests
	@cd client && go test ./...

init-test-env:
	docker-compose up -d
	@echo "> Waiting 45 seconds for ElasticSearch to fully boot..."
	@sleep 45

create-es-indices:
	go run main.go index publication create
	go run main.go index dataset create

create-db-tables:
	go install github.com/jackc/tern@latest
	tern migrate --migrations etc/snapstore/migrations

api-server:
	go run main.go api start


# Declare all targets as "PHONY", see https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html.
MAKEFLAGS += --always-make
.PHONY: init-test-env create-es-indices create-db-tables api-server setup-test-env tear-test-env run-test-grpc-client