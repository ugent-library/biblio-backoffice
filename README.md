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
BIBLIO_BACKEND_SESSION_SECRET
BIBLIO_BACKEND_CSRF_SECRET
BIBLIO_BACKEND_LIBRECAT_URL
BIBLIO_BACKEND_LIBRECAT_USERNAME
BIBLIO_BACKEND_LIBRECAT_PASSWORD
BIBLIO_BACKEND_ORCID_CLIENT_ID
BIBLIO_BACKEND_ORCID_CLIENT_SECRET
BIBLIO_BACKEND_ORCID_CLIENT_SANDBOX
BIBLIO_BACKEND_OIDC_URL
BIBLIO_BACKEND_OIDC_CLIENT_ID
BIBLIO_BACKEND_OIDC_CLIENT_SECRET

```
The following variables may be set:

```
BIBLIO_BACKEND_HOST
BIBLIO_BACKEND_PORT
BIBLIO_BACKEND_MODE
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
