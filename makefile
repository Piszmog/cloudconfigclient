# Parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test

# Commands
all: clean vet test build
clean:
	$(GO_CLEAN)
vet:
	$(GO_CMD) vet
test:
	$(GO_TEST) -v ./...
build:
	$(GO_BUILD)