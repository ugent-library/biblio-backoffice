# Installation

These installation instructions will guide you to set up this project on MacOS for local development.

## Prerequisites

These dependencies are required:

* PostgreSQL 
* ElasticSearch 6.8
* An OIDC provider (e.g. Keycloak)
* MongoDB

For development, you will also need:

* Node 14
* Go
* Python 2 (for libsass)

For importing data into the authority backend:

* Benthos

## Preparation

### Postgres

On Macos:

```
brew install postgresql@15
brew services start postgresql@15
```

After installing PostgreSQL, make the postgresql applications available:

```
echo 'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
# check if you can access postgres
psql postgres
```

Create a new PostgreSQL database and user:

```
createuser biblio --interactive --pwprompt --createdb
createdb biblio -O biblio -T template0 -l en_US.UTF-8 -E UTF8
```

Grant all privileges to the `biblio` user on the `biblio` database:

```
psql postgres
# Next execute this SQL query
GRANT ALL PRIVILEGES ON DATABASE biblio TO biblio;
# Check if you can see the databse, its grants:
\l
\q
```

### Elasticsearch

On Macos:

```
brew install elasticsearch@6
brew services start elasticsearch@6
```

**Note**: Elasticsearch@6 is disabled since 06/2023. You will need to edit the formula and 
change this line:

```
disable! date: "2023-06-19", because: :unsupported
```

to

```
deprecate! date: "2022-02-10", because: :unsupported
```

After installing ElasticSearch, you will need to install the `analyzer-icu` plugin:

```
echo 'export PATH="/opt/homebrew/opt/elasticsearch@6/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
elasticsearch-plugin install analysis-icu
```

Finally, restart the elasticsearch server to make sure the analysis-icu plugin is loaded:

```
brew services restart elasticsearch@6
```

### Mongodb

Refer to https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-os-x/ for installation instructions:

On Macos:

```
brew tap mongodb/brew
brew update
brew install mongodb-community@6.0
brew services start mongodb-community@6.0
```

### Benthos

[Benthos](https://benthos.dev) is a declarative data streaming service. It's a data-processing toolkit.

On Macos:

```
brew install benthos
```

### Go

On Macos

```
brew install go
```

Next, add the go BIN folder to your PATH variable:

```
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```

### NodeJS

On Macos

```
brew install nvm
```

Add this to the end of your `~/.zshrc` file:

```
export NVM_DIR="$HOME/.nvm"
[ -s "/opt/homebrew/opt/nvm/nvm.sh" ] && \. "/opt/homebrew/opt/nvm/nvm.sh"  # This loads nvm
[ -s "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm" ] && \. "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm"  # This loads nvm bash_completion
```

You should be able to run `nvm` after sourcing the `.zshrc` file:

```
source ~/.zshrc
nvm
```

Next, install node v14.17.1:

```
# This can take long while
nvm install v14.17.1
```

The `biblio-backoffice` directory has an `.nvmrc` file which specifies the version number of NodeJS required to run all the frontend tools required for frontend development. Before you start developing, do this:

```
cd biblio-backoffice
nvm use
# output
Found '/Users/netsensei/Workspace/biblio-backoffice/.nvmrc' with version <v14.17.1>
Now using node v14.17.1 (npm v6.14.13)
```

You're now all set.

## Step 1: Setting up a basic development environment

Clone the repository to your local machine:

```bash
git clone git@github.com:ugent-library/biblio-backoffice.git
cd biblio-backoffice
```

Next, copy these configuration files:

```
cp .env.example .env
cp .reflex.example.conf .reflex.conf
```

Configure the `.env` file accordingly:

```
# Run in development mode
export BIBLIO_BACKOFFICE_MODE="development"
# Base URL of the installation
export BIBLIO_BACKOFFICE_BASE_URL="http://localhost:3001"
# Credentials for the gRPC service
export BIBLIO_BACKOFFICE_ADMIN_USERNAME="admin"
export BIBLIO_BACKOFFICE_ADMIN_PASSWORD="admin"
export BIBLIO_BACKOFFICE_CURATOR_USERNAME="curator"
export BIBLIO_BACKOFFICE_CURATOR_PASSWORD="curator"
# Port on which the gRPC service listens
export BIBLIO_BACKOFFICE_PORT="3002"
# Set to a path where any uploaded files will reside:
export BIBLIO_BACKOFFICE_FILE_DIR=""
# session secret
export BIBLIO_BACKOFFICE_SESSION_SECRET="secret"
# CSRF secret
export BIBLIO_BACKOFFICE_CSRF_SECRET="secret"
# Connection to a biblio frontend installation 
export BIBLIO_BACKOFFICE_FRONTEND_URL=
export BIBLIO_BACKOFFICE_FRONTEND_USERNAME=
export BIBLIO_BACKOFFICE_FRONTEND_PASSWORD=
# ORCID credentials (Sandbox credentials)
export BIBLIO_BACKOFFICE_ORCID_CLIENT_ID=
export BIBLIO_BACKOFFICE_ORCID_CLIENT_SECRET=
export BIBLIO_BACKOFFICE_ORCID_SANDBOX="true"
# Connection to an OIDC provider
export BIBLIO_BACKOFFICE_OIDC_URL=
export BIBLIO_BACKOFFICE_OIDC_CLIENT_ID=
export BIBLIO_BACKOFFICE_OIDC_CLIENT_SECRET=
# Connection to PostgreSQL
export BIBLIO_BACKOFFICE_PG_CONN="postgres://biblio:biblio@localhost:5432/biblio"
# Connection to ElasticSearch
export BIBLIO_BACKOFFICE_ES6_URL="http://localhost:9200"
export BIBLIO_BACKOFFICE_DATASET_INDEX="biblio_datsets"
export BIBLIO_BACKOFFICE_PUBLICATION_INDEX="biblio_publications"
# Connection to MongoDB
export BIBLIO_BACKOFFICE_MONGODB_URL="mongodb://localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
# Connection to the ElasticSearch of Biblio Frontend
# For local development, use the same ElasticSearch instance as the one used by Biblio Backoffice
export BIBLIO_BACKOFFICE_FRONTEND_ES6_URL="http://localhost:9200"
```

Next, create the database tables by running the migrations:

```
cp .tern.example.conf tern.conf
```

Then open `tern.conf` and change the `port` variable to a value that matches the port on which your
postgresql server listens (default, this should be: 5432).

Now run the migrations using the `tern` program:

````
go install github.com/jackc/tern@latest
tern migrate --migrations etc/snapstore/migrations
```

Next, initialize the ElasticSearch indices:

```
go run main.go reset
```

Next, use `npm` to download all dependencies for frontend development:

```
npm install
```

## Step 2: Boot the application, start developing

Development involves re-compiling and running the codebase via `go run`. Repeatedly executing this
command is cumbersome. The `reflex` utility solves this. It watches the codebase for any changes to
files and will re-compile and restart the application.

Install reflex: 

```
go install github.com/cespare/reflex@latest
```

Starting the application:

```
reflex -d none -c .reflex.conf
```

If the configuration in `.env` contains no errors, you should see this output:

```
Starting service
2023-04-20T11:28:23.537+0200    INFO    commands/server.go:151  starting server at localhost:3001
```

Open your browser and navigate to `http://localhost:3001`. You should see the application.

If you want to edit the JavaScript or SCSS code, you will need to run the `mix` utility in a 
seperate terminal session. Laravel Mix will watch the assets directory and rebuild the assets
when it detects changes to the files.

```
npx mix
```

## Step 3: setup the authority backend

The authority backend embedded in the bibli-backoffice application serves two purposes:

* Look up of profile information used for authorizing users after OIDC authentication during login.
* A historic database of projects, organizations (affiliations) and people used to enrich records
  (i.e. Contributor fields, Project fields, Departments,...).

The authority backend consists of:

* A mongodb authority database containing these collections:
  * person - profile information
  * project - project information
* Several ElasticSearch indices which power the "suggest" feature for people, projects and departments:
  * person
  * project
  * organization

Without an initialized authority backend, users can't log in, or use the "suggest" feature.

### Initializing the mongoDB database with minimal data (optional)

Create a new file called person.json with a single record. Fill out the fields with your own account info:

```
{
  "_id": "00000000-0000-0000-0000-000000000001",
  "date_updated": "2022-05-06T22:02:50Z",
  "first_name": "YOUR FIRST NAME",
  "ugent_username": "UGENT USERNAME",
  "ugent_id": [
    "UGENT ID"
  ],
  "email": "UGENT EMAIL",
  "roles": [
    "biblio-admin"
  ],
  "ids":  [
    "00000000-0000-0000-0000-000000000001"
  ]
  "date_created": "2022-05-06T22:02:50Z",
  "last_name": "YOUR FIRST NAME",
  "active": 1,
  "full_name": "YOUR FULL NAME",
  "department": [
    {
      "_id": "CA20"
    }
  ]
}
```

Then, import that file in mongo:

```
mongoimport --uri "mongodb://localhost:27017/authority?directConnection=true&serverSelectionTimeoutMS=2000" --collection person person.json
```

This will automatically create a new database called `authority` with a `person` collection.

**You should now be able to log in using your own UGent account credentials via the OIDC provider.**

### Initializing the ElasticSearch indices

Create these ElasticSearch indices:

```
curl -X PUT http://localhost:9200/biblio_person -H 'Content-Type: application/json' -d @etc/es6/person.json
curl -X PUT http://localhost:9200/biblio_project -H 'Content-Type: application/json' -d @etc/es6/project.json
curl -X PUT http://localhost:9200/biblio_organization -H 'Content-Type: application/json' -d @etc/es6/organization.json
```

While empty, you need these for the "suggest contributor, project, department" functionality to work without throwing errors.

### Importing a complete dataset in the authority backend

You will need to obtain a full copy from the authority backend on either test or production environment:

* MongoDB dump of the authority > person collection: `authority_person.json`
* MongoDB dump of the authority > project collection: `authority_project.json`
* ES dump of the organizations index: `es_organizations.json`

Importing the MongoDB data dumps:

```
mongoimport --uri "mongodb://localhost:27017/authority?directConnection=true&serverSelectionTimeoutMS=2000" --collection person authority_person.json
mongoimport --uri "mongodb://localhost:27017/authority?directConnection=true&serverSelectionTimeoutMS=2000" --collection project authority_project.json
```

Next, load the data from the dumps into ElasticSearch using [Benthos](https://benthos.dev).

Create a new file called `benthos.yml` and add the configuration below:

```
input:
  label: ""
  stdin:
    codec: lines
    max_buffer: 1000000
pipeline:
  threads: -1
  processors:
    - mapping: |
        root = this
        root.id = this._id
        root._id = deleted()
        root.active = deleted()
output:
  label: ""
  elasticsearch:
    urls: [ "http://localhost:9200" ]
    index: "biblio_project"
    id: ${! json("id") }
    type: "project"
    max_in_flight: 64
    batching:
      count: 0
      byte_size: 0
      period: ""
      check: ""
```

This file will allow us to import `authority_project.json` file into the `biblio_project` index.

Run this command to import projects:

```
cat authority_project.json | benthos -c benthos.yml
```

Now, change the YAML file:

```
index: "biblio_person"
type: "person"
```

Run this command to import persons:

```
cat authority_person.json | benthos -c benthos.yml
```

## Step 4: import publications / datasets

### Starting gRPC API


First, start the API service in a separate terminal window/tab, don't close the terminal tab where the backoffice service is running!

```
source .env
go run main.go api start
```

Output:

```
2023-04-21T13:48:33.744+0200    INFO    commands/api.go:59      Listening at localhost:3002
```

The gRPC endpoint is now reachable at `localhost:3002`. Typically, gRPC endpoints are protected via 
SSL/TLS, but for local development, this isn't used.

The endpoint is secured via Basic Authorization. The username & password for the admin & curate users is configured in the `.env`. file via these variables:

```
export BIBLIO_BACKOFFICE_ADMIN_USERNAME="admin"
export BIBLIO_BACKOFFICE_ADMIN_PASSWORD="admin"
export BIBLIO_BACKOFFICE_CURATOR_USERNAME="curator"
export BIBLIO_BACKOFFICE_CURATOR_PASSWORD="curator"
```

### Testing the gRPC service

To test if you can connect with the gRPC service, use the `grpcurl` command. This is a version of `curl` but specifically for gRPC instead of HTTP. You will need to open a new terminal window, don't close the window which has the gRPC service running!

First, install `grpcurl`

```
brew install grpcurl
```

Next, test the connection:

```
echo -n "admin:admin" | base64
# copy the output and replace CREDS in the command below
grpcurl -v -expand-headers -H 'Authorization:Basic CREDS' -plaintext localhost:3002 list biblio.v1.Biblio
# output should be a list of available gRPC methods:
biblio.v1.Biblio.AddDatasets
biblio.v1.Biblio.AddFile
biblio.v1.Biblio.AddPublications
...
```

### Using the gRPC biblio client

The gRPC client code lives in a separate `client` directory in the `biblio-backoffice` project.
It's a separate program with its own `main.go` file.

Create a new `.env` file in the `client` directory:

```
cd client
touch .env
```

Next, edit the file and add these exported variables:

```
export BIBLIO_BACKOFFICE_HOST="localhost"
export BIBLIO_BACKOFFICE_PORT=3002
export BIBLIO_BACKOFFICE_USERNAME="admin"
export BIBLIO_BACKOFFICE_PASSWORD="admin"
export BIBLIO_BACKOFFICE_INSECURE="true"
```

Let's try the client:

```
# Load the configured variables
source .env
# Try to fetch all publications
go run main.go publication get-all
```

At this point, the `publication get-all` command isn't going to return anything since there's no record yet.

### Importing publications and records

You will need to obtain a dump of publications / datasets from either the test or production environment. These are 2 files containing records in JSONL format (JSON Lines, each line contains a separate record).

Let's import publications:

```
cd client
source .env
go run main.go publication import < location-to-data/publications.jsonl
```

You should see this output:

```
stored and index publication XXX at line XXX
```

Same for datasets:

```
go run main.go dataset import < location-to-data/datasets.jsonl
```

