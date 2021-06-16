
# PRE COMMIT
# -------------------------------------------------------

# This should be run prior to any commits, runs the various tools that should pass before committing code.
precommit: fmt test.race.nocache cover lint



# DEPENDENCIES
# -------------------------------------------------------

# Install all required dependencies that can reasonably be installed (ex. this will not install Docker).
deps:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/godoc@latest



# BUILDING
# -------------------------------------------------------

# Build the application.
build:
	@rm -rf _build
	@mkdir -p _build
	go build -o _build/messageapp ./cmd/messageapp/messageapp.go
	@echo "App can be found at: _build/messageapp"



# LINTING + FORMATING
# -------------------------------------------------------

# Run go code formatting on the code base.
fmt:
	go fmt ./...

# Run all available linting and static analysis tools.
lint: staticcheck vet

# Run the staticcheck static analysis tool (https://staticcheck.io/).
staticcheck:
	staticcheck ./...

# Run the go vet command.
vet:
	go vet ./...



# TESTING + CODE COVERAGE
# -------------------------------------------------------

# Same as test.race but clears the test cache first.
test.race.nocache:
	go clean -testcache && go test -race ./...

# Run all tests with the race detector enabled.
test.race:
	go test -race ./...

# Run all tests for the project.
test:
	go test ./...

# Run only test that are not related to the database. Database test can be slow to run this is helpful when you are
# mostly running unit and fast integration tests.
test.nodb:
	NODB=1 go test ./...

# Run property based tests, note this can take a while.
test.prop:
	go test --tags=propertyTests ./messages

# Print code coverage.
cover:
	go test -cover ./...

# Generate and view code coverage report.
cover.report:
	@mkdir -p _tmp
	@go test -coverprofile _tmp/coverage.out ./...
	@go tool cover -html _tmp/coverage.out



# DOCUMENTATION
# -------------------------------------------------------

# Run Godoc documentation server.
docs:
	@echo "View at http://localhost:3000/pkg/github.com/mdev5000/qlik_message"
	godoc -http=:3000

# Generate the API documentation using openapi.
docs.api.gen:
	@rm -rf _docs/api

	docker run --rm \
      -u $$(id -u ${USER}):$$(id -g ${USER}) \
      -v ${PWD}:/local openapitools/openapi-generator-cli generate \
      -i /local/_openapi/messages.json \
      -g markdown \
      -o /local/_docs/api

	docker run --rm \
      -u $$(id -u ${USER}):$$(id -g ${USER}) \
      -v ${PWD}:/local openapitools/openapi-generator-cli generate \
      -i /local/_openapi/messages.json \
      -g html2 \
      -o /local/_docs/api/html

# MISC
# -------------------------------------------------------

FORCE: