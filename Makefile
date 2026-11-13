GO ?= go
GOFMT ?= gofmt

.PHONY: check-go build-cli
check-go:
	@test -z "$$($(GOFMT) -l cmd/storm/*.go cmd/storm-package/*.go internal/cli/*.go internal/release/*.go test/acceptance/*.go)"
	$(GO) vet ./...
	@mkdir -p .artifacts
	$(GO) test -race -count=1 -covermode=atomic -coverpkg=./... -coverprofile=.artifacts/coverage.out ./...
	$(GO) tool cover -func=.artifacts/coverage.out
	@$(GO) tool cover -func=.artifacts/coverage.out | awk '/^total:/ { if ($$3+0 < 97) exit 1 }'
	$(GO) build ./...
	$(GO) mod verify
	$(GO) mod tidy -diff

build-cli:
	@mkdir -p .artifacts
	$(GO) build -trimpath -o .artifacts/storm ./cmd/storm
