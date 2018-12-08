.PHONY: all help test build

all: build

deps:
	@docker run --rm -it -v "${GOPATH}"/src:/go/src \
			-w /go/src/github.com/pottava/ecr-creds \
			supinf/go-dep:0.5 ensure

test:
	@docker run --rm -it -v "${GOPATH}"/src:/go/src \
			-w /go/src/github.com/pottava/ecr-creds \
			supinf/golangci-lint:1.12 \
			run --config .golangci.yml
	@docker run --rm -it -v "${GOPATH}"/src:/go/src \
			-w /go/src/github.com/pottava/ecr-creds \
			--entrypoint go supinf/golangci-lint:1.12 \
			test -vet off `go list ./...`

build:
	@docker run --rm -it -v "${GOPATH}"/src:/go/src \
			-w /go/src/github.com/pottava/ecr-creds \
			pottava/gox:go1.11 --osarch "linux/amd64 darwin/amd64 windows/amd64" \
			-ldflags "-s -w" -output "dist/{{.OS}}_{{.Arch}}"
