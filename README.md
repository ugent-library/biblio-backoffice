# biblio-backend

## Configuration

Configuration can be passed as an argument:

```
go run main.go server start --session-secret mysecret
```

Or as an env variable:

```
BIBLIO_BACKEND_SESSION_SECRET=mysecret go run main.go server start
```

The following variables must be set:

```
BIBLIO_BACKEND_FILE_DIR
BIBLIO_BACKEND_SESSION_SECRET
BIBLIO_BACKEND_FRONTEND_URL
BIBLIO_BACKEND_FRONTEND_USERNAME
BIBLIO_BACKEND_FRONTEND_PASSWORD
BIBLIO_BACKEND_ORCID_CLIENT_ID
BIBLIO_BACKEND_ORCID_CLIENT_SECRET
BIBLIO_BACKEND_ORCID_CLIENT_SANDBOX
BIBLIO_BACKEND_OIDC_URL
BIBLIO_BACKEND_OIDC_CLIENT_ID
BIBLIO_BACKEND_OIDC_CLIENT_SECRET
BIBLIO_BACKEND_CSRF_SECRET

```
The following variables may be set:

```
BIBLIO_BACKEND_HOST
BIBLIO_BACKEND_PORT
BIBLIO_BACKEND_MODE
BIBLIO_BACKEND_SESSION_NAME
BIBLIO_BACKEND_CSRF_NAME
```

## Assets

Install node dependencies:

```
cd services/webapp/
npm install
```

Build assets:

```
cd services/webapp/
npx mix
```

Watch file changes in development:

```
cd services/webapp/
npx mix watch
```

Build production assets:

```
cd services/webapp/
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

## Start server

Install node dependencies:

```
npm install
```

To start the development server with live reload:

```
npm run dev
```

Run the server directly:

```
go run main.go server start
```

## Build

```
go build -o biblio-backend main.go
```
