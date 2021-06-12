
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

# Run the staticcheck static analysis tool (https://staticcheck.io/)
staticcheck:
	staticcheck ./...

# Run the go vet command
vet:
	go vet ./...



# TESTING
# -------------------------------------------------------

# Run all tests for the project
test:
	go test ./...

FORCE: