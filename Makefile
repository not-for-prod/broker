linter:
	golangci-lint --config .golangci.yaml run

fmt:
	golangci-lint --config .golangci.yaml fmt