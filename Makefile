PKG=github.com/ocramh/challenge2

.PHONY: install
install:
	go install ${PKG}/cmd/provider

.PHONY: test
test:
	@go test ./... -coverprofile=cover.out