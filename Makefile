install-dev:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.7.2

lint:
	golangci-lint run 

lint-fix:
	golangci-lint fmt
	golangci-lint run --fix
