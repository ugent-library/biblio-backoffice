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
BIBLIO_BACKEND_SESSION_SECRET=mysecret go run main.go server start
```

The following variables must be set:

```bash
BIBLIO_BACKEND_BASE_URL
BIBLIO_BACKEND_SESSION_SECRET
BIBLIO_BACKEND_FILE_DIR
BIBLIO_BACKEND_FRONTEND_URL
BIBLIO_BACKEND_FRONTEND_USERNAME
BIBLIO_BACKEND_FRONTEND_PASSWORD
BIBLIO_BACKEND_OIDC_URL
BIBLIO_BACKEND_OIDC_CLIENT_ID
BIBLIO_BACKEND_OIDC_CLIENT_SECRET
BIBLIO_BACKEND_CSRF_NAME
BIBLIO_BACKEND_CSRF_SECRET
BIBLIO_BACKEND_ORCID_CLIENT_ID
BIBLIO_BACKEND_ORCID_CLIENT_SECRET
BIBLIO_BACKEND_ORCID_SANDBOX
BIBLIO_BACKEND_PG_CONN (default: postgres://localhost:5432/biblio_backend?sslmode=disable)
BIBLIO_BACKEND_PUBLICATION_INDEX (default: biblio_backend_publications)
BIBLIO_BACKEND_DATASET_INDEX (default: biblio_backend_datasets)
BIBLIO_BACKEND_ES6_URL (default: http://localhost:9200)
```

The following variables may be set:

```bash
BIBLIO_BACKEND_MODE (default: production)
BIBLIO_BACKEND_PORT (default: 3000)
BIBLIO_BACKEND_HOST (default: localhost)
BIBLIO_BACKEND_SESSION_NAME (default: biblio-backoffice)
BIBLIO_BACKEND_SESSION_MAX_AGE (default: 86400 * 30 // 30 days)
BIBLIO_BACKEND_HDL_SRV_ENABLED (default: false)
BIBLIO_BACKEND_HDL_SRV_URL
BIBLIO_BACKEND_HDL_SRV_PREFIX
BIBLIO_BACKEND_HDL_SRV_USERNAME
BIBLIO_BACKEND_HDL_SRV_PASSWORD
```

For the gRPC server and client:

```bash
BIBLIO_BACKEND_HOST
BIBLIO_BACKEND_PORT
BIBLIO_BACKEND_USERNAME
BIBLIO_BACKEND_PASSWORD
BIBLIO_BACKEND_INSECURE (default: false)
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
export BIBLIO_BACKEND_MODE="development"
export BIBLIO_BACKEND_BASE_URL="http://localhost:3001"
...
```

Alternatively, adapt this command in your `.reflex.conf` to suit your needs.
```
'source .env && go run main.go server start --host localhost --port 3001'
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
PGDATABASE=biblio_backend tern migrate
```

More info [here](https://github.com/jackc/tern).

List of PG env variables [here](https://www.postgresql.org/docs/current/libpq-envars.html).
