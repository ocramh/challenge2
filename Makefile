PKG=github.com/ocramh/challenge2

.PHONY: install
install:
	go install ${PKG}/cmd/challenge2

.PHONY: test
test:
	@go test ./... -coverprofile=cover.out