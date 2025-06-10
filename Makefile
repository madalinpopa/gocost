.PHONY: test staticcheck check build run

.DEFAULT_GOAL := test

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
#   Run GO commands
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
build:
	go build -o bin/gocost ./cmd/gocost

# Cross-platform builds
build/linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/gocost-linux-amd64 ./cmd/gocost

build/windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/gocost-windows-amd64.exe ./cmd/gocost

build/darwin:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/gocost-darwin-arm64 ./cmd/gocost

build/all: build/linux build/windows build/darwin

run:
	go run ./cmd/gocost

# run go tests
test:
	go test ./internal...

# check for data race conditions
test/race:
	go -race test ./internal...

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
#   Run code format and code style commands
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

# run go vet tool
vet:
	go vet ./internal...

# run staticcheck tool
staticcheck:
	go tool staticcheck ./internal...

# run all tools
check: vet staticcheck


# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
#   Run release commands
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
release/patch:
		@if [ $$(git tag | wc -l) -eq 0 ]; then \
    		NEW_TAG="v0.0.1"; \
    	else \
    		LATEST_TAG=$$(git describe --tags `git rev-list --tags --max-count=1`); \
    		MAJOR=$$(echo $$LATEST_TAG | cut -d. -f1 | tr -d 'v'); \
    		MINOR=$$(echo $$LATEST_TAG | cut -d. -f2); \
    		PATCH=$$(echo $$LATEST_TAG | cut -d. -f3); \
    		NEW_PATCH=$$((PATCH + 1)); \
    		NEW_TAG="v$$MAJOR.$$MINOR.$$NEW_PATCH"; \
    	fi; \
    	git tag -a $$NEW_TAG -m "Release $$NEW_TAG" && \
    	echo "Created new tag: $$NEW_TAG"


release/minor:
		@if [ $$(git tag | wc -l) -eq 0 ]; then \
    		NEW_TAG="v0.1.0"; \
    	else \
    		LATEST_TAG=$$(git describe --tags `git rev-list --tags --max-count=1`); \
    		MAJOR=$$(echo $$LATEST_TAG | cut -d. -f1 | tr -d 'v'); \
    		MINOR=$$(echo $$LATEST_TAG | cut -d. -f2); \
    		NEW_MINOR=$$((MINOR + 1)); \
    		NEW_TAG="v$$MAJOR.$$NEW_MINOR.0"; \
    	fi; \
    	git tag -a $$NEW_TAG -m "Release $$NEW_TAG" && \
    	echo "Created new tag: $$NEW_TAG"

release/major:
		@if [ $$(git tag | wc -l) -eq 0 ]; then \
    		NEW_TAG="v1.0.0"; \
    	else \
    		LATEST_TAG=$$(git describe --tags `git rev-list --tags --max-count=1`); \
    		MAJOR=$$(echo $$LATEST_TAG | cut -d. -f1 | tr -d 'v'); \
    		NEW_MAJOR=$$((MAJOR + 1)); \
    		NEW_TAG="v$$NEW_MAJOR.0.0"; \
    	fi; \
    	git tag -a $$NEW_TAG -m "Release $$NEW_TAG" && \
    	echo "Created new tag: $$NEW_TAG"
