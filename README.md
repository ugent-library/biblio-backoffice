# biblio-backoffice

The backoffice application for [Biblio](https://biblio.ugent.be), the Ghent University
Academic Bibliography.

## Introduction

The application consists of three components. A Web server which serves the Web UI, a
[grpc](https://grpc.io/) server and a grpc CLI client. The Web UI is intended for
librarians and researchers. The grpc server/client is geared towards data management
by system administrators and data curators.

## Prerequisites

The application stores data in several stores. You will need:

* A PostgreSQL database
* An ElasticSearch index
* An OpenID Connect endpoint (e.g. [Keycloak](https://www.keycloak.org/))
* Disk storage

## Quickstart Web UI

### Manual

```bash
git clone git@github.com:ugent-library/biblio-backoffice.git
go run main.go
```

Alternatively, build the binary:

```bash
git clone git@github.com:ugent-library/biblio-backoffice.git
cd biblio-backoffice
go build -o /tmp/biblio-backoffice
./tmp/biblio-backoffice
# Or cross compile e.g.
GOOS=linux GOARCH=amd64 go build -o /tmp/biblio-backoffice.linux.amd64
```

Starting the Web UI server:

```bash
go run main.go start server
```

The server will exit with an error if it can't connect to a PostgreSQL,
ElasticSearch or an OpenID Connect endpoint.

Refer to the configuration section to learn how to configure the Web UI
server.

### Docker

```bash
docker run --name some-backoffice -d ugentlib/biblio-backoffice:dev
```

Use `docker-compose` to manage all arguments (ports, volumes, env variables,...)

## GRPC Server & Client

### Server

The gRPC protocol requires TLS support. The server itself currently doesn't
have TLS support built-in. You will need to rely on a proxy which provides
TLS and funnels traffic to the gRPC server. [Traefik v2](https://doc.traefik.io/traefik/user-guides/grpc/)
offers this via Let's Encrypt. The gRPC server shouldn't run under a subpath,
but directly on the root path of a FQDN separate from the domainname used by the
Web UI e.g. grpc.myapplication.tld.

Starting the GRPC server:

```bash
go run main.go api start
# or
./tmp/biblio-backoffice api start
```

The server will exit with an error if it can't connect to a PostgreSQL,
ElasticSearch or an OpenID Connect endpoint.

Refer to the configuration section to learn how to configure the gRPC server.

### Client

```bash
git clone git@github.com:ugent-library/biblio-backoffice.git
cd biblio-backoffice/client
go run main.go --host <grpc.host.tld> --port <port> --username <username> --password <password>
```

Or build & run as a (distributable) binary:

```bash
git clone git@github.com:ugent-library/biblio-backoffice.git
cd biblio-backoffice
go build -o /tmp/biblio-backoffice-client client/main.go
./tmp/biblio-backoffice-client
# Or cross compile e.g.
GOOS=linux GOARCH=amd64 go build -o /tmp/biblio-backoffice-client.linux.amd64 client/main.go
```

The gRPC client is intended to be used from any machine (e.g. a data
curator's local machine), not just the same server where the gRPC server is
running. This allows for anyone with credentials to manage data remotely.


## Configuration

Configuration can be passed as an argument:

```
go run main.go server start --session-secret mysecret
```

Or as an environment variables:

```
BIBLIO_BACKOFFICE_SESSION_SECRET=mysecret go run main.go server start
```

The following variables must be set:

```bash
BIBLIO_BACKOFFICE_BASE_URL
BIBLIO_BACKOFFICE_SESSION_SECRET
BIBLIO_BACKOFFICE_FILE_DIR
BIBLIO_BACKOFFICE_FRONTEND_URL
BIBLIO_BACKOFFICE_FRONTEND_USERNAME
BIBLIO_BACKOFFICE_FRONTEND_PASSWORD
BIBLIO_BACKOFFICE_OIDC_URL
BIBLIO_BACKOFFICE_OIDC_CLIENT_ID
BIBLIO_BACKOFFICE_OIDC_CLIENT_SECRET
BIBLIO_BACKOFFICE_CSRF_NAME
BIBLIO_BACKOFFICE_CSRF_SECRET
BIBLIO_BACKOFFICE_ORCID_CLIENT_ID
BIBLIO_BACKOFFICE_ORCID_CLIENT_SECRET
BIBLIO_BACKOFFICE_ORCID_SANDBOX
BIBLIO_BACKOFFICE_PG_CONN (default: postgres://localhost:5432/biblio_backoffice?sslmode=disable)
BIBLIO_BACKOFFICE_PUBLICATION_INDEX (default: biblio_backoffice_publications)
BIBLIO_BACKOFFICE_DATASET_INDEX (default: biblio_backoffice_datasets)
BIBLIO_BACKOFFICE_ES6_URL (default: http://localhost:9200)
BIBLIO_BACKOFFICE_MONGODB_URL (default: mongodb://localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000)
BIBLIO_BACKOFFICE_FRONTEND_ES6_URL (default: http://localhost: 9200)
```

The following variables may be set:

```bash
BIBLIO_BACKOFFICE_MODE (default: production)
BIBLIO_BACKOFFICE_PORT (default: 3000)
BIBLIO_BACKOFFICE_HOST (default: localhost)
BIBLIO_BACKOFFICE_SESSION_NAME (default: biblio-backoffice)
BIBLIO_BACKOFFICE_SESSION_MAX_AGE (default: 86400 * 30 // 30 days)
BIBLIO_BACKOFFICE_HDL_SRV_ENABLED (default: false)
BIBLIO_BACKOFFICE_HDL_SRV_URL
BIBLIO_BACKOFFICE_HDL_SRV_PREFIX
BIBLIO_BACKOFFICE_HDL_SRV_USERNAME
BIBLIO_BACKOFFICE_HDL_SRV_PASSWORD
BIBLIO_BACKOFFICE_OAI_API_URL
BIBLIO_BACKOFFICE_OAI_API_KEY
```

For the gRPC server:

```bash
BIBLIO_BACKOFFICE_HOST
BIBLIO_BACKOFFICE_PORT
BIBLIO_BACKOFFICE_ADMIN_USERNAME
BIBLIO_BACKOFFICE_ADMIN_PASSWORD
BIBLIO_BACKOFFICE_CURATOR_USERNAME
BIBLIO_BACKOFFICE_CURATOR_PASSWORD
BIBLIO_BACKOFFICE_INSECURE (default: false)
BIBLIO_BACKOFFICE_TIMEOUT (default: 5s)
```

And the gRPC client:

```bash
BIBLIO_BACKOFFICE_HOST
BIBLIO_BACKOFFICE_PORT
BIBLIO_BACKOFFICE_USERNAME
BIBLIO_BACKOFFICE_PASSWORD
BIBLIO_BACKOFFICE_INSECURE (default: false)
```

## Development

This project uses [reflex](https://github.com/cespare/reflex) to watch for file
changes and recompile the application and assets:

```
cd biblio-backoffice
go install github.com/cespare/reflex@latest
cp .reflex.example.conf .reflex.conf
reflex -d none -c .reflex.conf
```

Refer to `.reflex.example.conf`. The command assumes the existence of a `.env` file
which exports all relevant environment variables:

```bash
export BIBLIO_BACKOFFICE_MODE="development"
export BIBLIO_BACKOFFICE_BASE_URL="http://localhost:3001"
...
```

Alternatively, adapt this command in your `.reflex.conf` to suit your needs.
```
'source .env && go run main.go server start --host localhost --port 3001'
```

### Running tests

This project has integration tests for the GRPC client/server.

Setting up the test environment:

```bash
# Set up a docker-based isolated, testing environment
make setup-test-env
# Run the GRPC API server
make api-server
```

In a separate terminal, you can now run the integration tests:

```bash
make run-tests
# or alternatively
go test client/... -v
```

Tearing down the test environment:

```bash
make tear-test-env
```

## SASS/SCSS & asset compilation

Install node dependencies:

```bash
npm install
```

Assets will be recompiled automatically if you use [reflex](https://github.com/cespare/reflex) (see above).

Build assets manually:

```
npx mix
```

Build production assets manually:

```
npx mix --production
```

Laravel Mix [documentation](https://laravel.com/docs/8.x).

## Database migrations

```
go install github.com/jackc/tern@latest
cd etc/pg/migrations
PGDATABASE=biblio_backoffice tern migrate
```

More info [here](https://github.com/jackc/tern).

List of PG env variables [here](https://www.postgresql.org/docs/current/libpq-envars.html).
