ENV ?= stg
GO_CMD ?= go

run:
	$(GO_CMD) run cmd/main.go --env=$(ENV)

test:
	$(GO_CMD) test ./... -v

test-class:
	$(GO_CMD) test ./internal/class/... -v

test-lesson:
	$(GO_CMD) test ./internal/lesson/... -v

test-user:
	$(GO_CMD) test ./internal/user/... -v

test-set:
	$(GO_CMD) test ./internal/set/... -v

test-question:
	$(GO_CMD) test ./internal/question/... -v

fmt:
	$(GO_CMD) fmt ./...

clean:
	$(GO_CMD) clean -testcache

lint:
	golangci-lint run ./...

ci: fmt lint test
