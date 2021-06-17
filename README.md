# Message API Demo

A REST service that allows clients to create and modify messages 
and request information about the messages.

## Architecture + Docs

For details on the architecture please see [here](./_docs/arch/README.md).

For Go documentation run:

```bash
make docs
```

---

## Setup

The following are required:
- `Go 1.16`
- `Docker Compose` + `Docker`
- `Make`

To setup the remaining dependencies run:

```bash
make deps
```

---

## Building + Running

### Building

To build the application run:

```bash
make build
```

### Running

To run specify the database to connect to, ex:

```bash
HOST=localhost PORT=8000 DATABASE_URL=postgresql://user:password@host messageapp
```

You can run db migrations prior to running the app via:

```bash
MIGRATE=1 ... messageapp
```

You can also run the server with TLS:

```bash
KEY=key.pem CERT=certificate.pem ... messageapp -tls
```

Full local example:

```bash
make build
docker-compose up -d
MIGRATE=1 HOST=localhost PORT=10443 \
KEY=_examples/cert/server.key CERT=_examples/cert/server.cert \
DATABASE_URL=postgresql://postgres:postgres@localhost?sslmode=disable \
./_build/messageapidemo -tls
```

---

## Development

To run the development server run the following:

```bash
docker-compose up -d
go run ./cmd/devserver/devserver.go

# or run the following for more options:
go run ./cmd/devserver/devserver.go -h
```

To tear the dev database down run:

```bash
docker-compose down
```

---

## Testing

Tests can be run via:

```bash
# run all tests
make test

# run all tests with race detector enabled
make test.race

# run all tests with race detector enabled and clear the test cache first.
make test.race.nocache
```

You can also test individual packages with the standard Go
tooling.

```bash
go test ./data # run tests for the data package
```

You can include property based tests by adding the `propertyTests` build tag and specify the number of cases via 
`NUM_PROP_TESTS`, ex:

```bash
NUM_PROP_TESTS=20000 go test --tags=propertyTests ./messages
```

### Coverage

You can check coverage via:

```bash
make cover
```

Or view a coverage report via:

```bash
make cover.report
```

---

## API Documentation

You can view the API documentation [here](./_docs/api/README.md) or for a nicer experience view it in html at 
`./_docs/api/html/index.html`

The API documentation is written in the `openapi` json specification generated using `openapi-generator`.

### API Overview

```bash
# list messages
curl http://localhost:8000/messages
curl http://localhost:8000/messages?fields=id,message&pageSize=20&pageStartIndex=2

# view message
curl http://localhost:8000/messages/5

# create message
curl -X POST http://localhost:8000/messages --data '{"message":"first"}' \
-H 'Content-Type: application/json; charset=UTF-8'

# update message
curl -X PUT http://localhost:8000/messages/1 --data '{"message":"new message"}' \
-H 'Content-Type: application/json; charset=UTF-8'

# delete message
curl -X DELETE http://localhost:8000/messages/1
```

### Update API documentation

To update the API documentation after changing the `/_openapi/messages.json` file, run the following:

```bash
make docs.api.gen
```

---

## Committing code

Before committing code you should run the `precommit` make
command. The command formats code, runs the tests with the
race detector enabled, runs linting, and reports on test
coverage.

All tests should pass and no linting errors should be present, ex.:

```
> make precommit
go fmt ./...
go test -race ./...
ok  	github.com/mdev5000/messageappdemo/data	(cached)
?   	github.com/mdev5000/messageappdemo/postgres	[no test files]
go test -cover ./...
ok  	github.com/mdev5000/messageappdemo/data	(cached)	coverage: 80.5% of statements
?   	github.com/mdev5000/messageappdemo/postgres	[no test files]
staticcheck ./...
go vet ./...
```