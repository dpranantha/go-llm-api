DEFAULT_GOAL := run

.PHONY: clean fmt vet build run run-air graphql-generate

# Default framework is Gin
FRAMEWORK ?= gin

clean:
	go clean

fmt:
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	@mkdir -p bin
	@for dir in $$(find . -name 'main.go' -exec dirname {} \;); do \
		name=$$(basename $$dir); \
		go build -o bin/$$name $$dir || exit 1; \
	done

run:
	@echo "Running with framework: $(FRAMEWORK)"
	@go run main.go -fiber=$(if $(findstring fiber,$(FRAMEWORK)),true,false)

run-air:
	air

graphql-generate:
	cd back-end/graphql && go get github.com/99designs/gqlgen@v0.17.70 && go run github.com/99designs/gqlgen generate
