.PHONY: build
build: tidy
	@go build -o build/casechan *.go


.PHONY: tidy
tidy:
	@go mod tidy
