# Qlik message service

A REST service that allows clients to create and modify messages 
and request information about the messages.

---

## Testing

Test can be run via:

```bash
# run all tests
make test

# run all tests with race detector enabled
make test.race
```

You can also test individual packages with the standard Go
tooling.

```bash
go test ./data # run tests for the data package
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
ok  	github.com/mdev5000/qlik_message/data	(cached)
?   	github.com/mdev5000/qlik_message/postgres	[no test files]
go test -cover ./...
ok  	github.com/mdev5000/qlik_message/data	(cached)	coverage: 80.5% of statements
?   	github.com/mdev5000/qlik_message/postgres	[no test files]
staticcheck ./...
go vet ./...
```