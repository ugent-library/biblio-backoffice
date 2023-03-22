export GO111MODULE=on

-include .env
BIBLIO_BACKOFFICE_USERNAME=admin
BIBLIO_BACKOFFICE_PASSWORD=admin
BIBLIO_BACKOFFICE_PG_CONN=postgres://biblio:biblio@localhost:2345/biblio
BIBLIO_BACKOFFICE_ES6_URL=http://localhost:9400
BIBLIO_BACKOFFICE_DATASET_INDEX=biblio_datasets
BIBLIO_BACKOFFICE_PUBLICATION_INDEX=biblio_publications
BIBLIO_BACKOFFICE_PORT=3999
BIBLIO_BACKOFFICE_MONGODB_URL=mongodb://localhost:27020/?directConnection=true&serverSelectionTimeoutMS=2000
BIBLIO_BACKOFFICE_FRONTEND_ES6_URL=http://localhost:9400
export

setup-test-env: init-test-env create-db-tables create-es-indices create-mongo-db

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
	curl -X PUT -H "Content-Type: application/json" http://localhost:9400/_cluster/settings -d '{ "transient": { "cluster.routing.allocation.disk.threshold_enabled": false } }'

create-es-indices:
	go run main.go publication reindex
	go run main.go dataset reindex
	curl -X PUT http://localhost:9400/biblio_person -H 'Content-Type: application/json' -d @etc/es6/person.json
	curl -X PUT http://localhost:9400/biblio_project -H 'Content-Type: application/json' -d @etc/es6/project.json
	curl -X PUT http://localhost:9400/biblio_organization -H 'Content-Type: application/json' -d @etc/es6/organization.json
	curl -X PUT -H "Content-Type: application/json" http://localhost:9400/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}'
	curl -X POST http://localhost:9400/biblio_person/person/_bulk -H "Content-Type: application/x-ndjson" --data-binary @etc/fixtures/person.es6.json
	curl -X POST http://localhost:9400/biblio_organization/organization/_bulk -H "Content-Type: application/x-ndjson" --data-binary @etc/fixtures/organization.es6.json

create-mongo-db:
	mongoimport --uri "mongodb://localhost:27020/authority?directConnection=true&serverSelectionTimeoutMS=2000" --collection person  ./etc/fixtures/person.json

create-db-tables:
	go install github.com/jackc/tern@latest
	tern migrate --migrations etc/snapstore/migrations

api-server:
	go run main.go api start


# Declare all targets as "PHONY", see https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html.
MAKEFLAGS += --always-make
.PHONY: init-test-env create-es-indices create-db-tables api-server setup-test-env tear-test-env run-test-grpc-client