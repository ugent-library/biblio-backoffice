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

- A PostgreSQL database
- An ElasticSearch index
- An OpenID Connect endpoint (e.g. [Keycloak](https://www.keycloak.org/))
- Disk storage

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
BIBLIO_BACKOFFICE_PG_CONN (e.g.: postgres://localhost:5432/biblio_backoffice?sslmode=disable)
BIBLIO_BACKOFFICE_PUBLICATION_INDEX (e.g.: biblio_backoffice_publications)
BIBLIO_BACKOFFICE_DATASET_INDEX (e.g.: biblio_backoffice_datasets)
BIBLIO_BACKOFFICE_ES6_URL (e.g.: http://localhost:9200)
BIBLIO_BACKOFFICE_MONGODB_URL (e.g.: mongodb://localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000)
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
BIBLIO_BACKOFFICE_HDL_SRV_URL (e.g.: http://localhost:4000/handles)
BIBLIO_BACKOFFICE_HDL_SRV_PREFIX (e.g. 1854)
BIBLIO_BACKOFFICE_HDL_SRV_USERNAME
BIBLIO_BACKOFFICE_HDL_SRV_PASSWORD
BIBLIO_BACKOFFICE_OAI_API_URL
BIBLIO_BACKOFFICE_OAI_API_KEY
BIBLIO_BACKOFFICE_TIMEZONE (default: Europe/Brussels)
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

This project uses [wgo](https://github.com/bokwoon95/wgo) to watch for file
changes and recompile the application and assets:

```bash
cd biblio-backoffice
make dev
# or
npm run dev
```

### Running tests

This project contains Cypress integration tests.

To run the tests in the CLI run:

```bash
npm test
```

If you want to run a recorded (CLI) run to the [Cypress Dashboard](https://cloud.cypress.io/projects/mjg74d/runs), first copy the file `cypress/cypress-record.sh.example` to `cypress/cypress-record.sh` and replace the record key with the secret key from Cypress. Then run:

```bash
./cypress/cypress-record.sh
```

To run the test in the GUI run:

```bash
npm run cypress:open
```

## SASS/SCSS & asset compilation

Install node dependencies:

```bash
npm install
```

Assets will be recompiled automatically if you use [wgo](https://github.com/bokwoon95/wgo) (see above).

Build assets manually:

```bash
node esbuild.mjs
# or
npm run build:assets
```

esbuild [documentation](https://esbuild.github.io/).

## Database migrations

```
go install github.com/jackc/tern@latest
cd etc/pg/migrations
PGDATABASE=biblio_backoffice tern migrate
```

More info [here](https://github.com/jackc/tern).

List of PG env variables [here](https://www.postgresql.org/docs/current/libpq-envars.html).

## Dev Containers

This project supports [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers). Following these steps
will auto setup a containerized development environment for this project. In VS Code, you will be able to start a terminal
that logs into a Docker container. This will allow you to write and interact with the code inside a self-contained sandbox.

**Installing the Dev Containers extension**

1. Open VS Code.
2. Go to the [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension page.
3. Click the `install` button to install the extension in VS Code.

**Add configuration**

1. `cp .devcontainer.env.example .devcontainer.env`

**Open in Dev Containers**

1. Open the project directory in VS Code.
2. Click on the green "Open a remote window" button in the lower left window corner.
3. Choose "reopen in container" from the popup menu.
4. The green button should now read "Dev Container: App name" when successfully opened.
5. Open a new terminal in VS Code from the `Terminal` menu link.

You are now logged into the dev container and ready to develop code, write code, push to git or execute commands.

**Run the project**

1. Open a new terminal in VS Code from the `Terminal` menu link.
2. Execute this command `make dev`.
3. Once the application has started, VS Code will show a popup with a link that opens the project in your browser.

**Networking**

The application and its dependencies run on these ports:

| Application       | Port |
| ----------------- | ---- |
| Biblio Backoffice | 3001 |
| Mock OIDC         | 3002 |
| DB Application    | 3051 |
| Elastic Search    | 3061 |
