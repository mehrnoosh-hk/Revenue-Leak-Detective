.PHONY: test

test:
	@echo "Formatting and testing Go code..."
	@go fmt ./services/api/...
	@go test ./services/api/... -v