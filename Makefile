
# PRE COMMIT
# -------------------------------------------------------

# This should be run prior to any commits, runs the various tools that should pass before committing code.
precommit: fmt test.race cover lint



# DEPENDENCIES
# -------------------------------------------------------

# Install all required dependencies that can reasonably be installed (ex. this will not install Docker).
dependencies: dependencies.linting

dependencies.linting:
	go install honnef.co/go/tools/cmd/staticcheck@latest



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

# Run all tests with the race detector enabled.
test.race:
	go test -race ./...

# Run all tests for the project.
test:
	go test ./...

# Print code coverage.
cover:
	go test -cover ./...

# Generate and view code coverage report.
cover.report:
	@mkdir -p _tmp
	@go test -coverprofile _tmp/coverage.out ./...
	@go tool cover -html _tmp/coverage.out



# DEV
# -------------------------------------------------------

# This should be run prior to any commits, runs the various tools that should pass before committing code.
dev.run:
	UID=$${UID} GID=$${GID} docker-compose up -d



# MISC
# -------------------------------------------------------

FORCE: