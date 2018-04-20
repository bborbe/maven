
all: test install run

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install cmd/maven-repo-server/*.go

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

vet:
	go tool vet .
	go tool vet --shadow .

lint:
	golint -min_confidence 1 ./...

errcheck:
	errcheck -ignore '(Close|Write)' ./...

check: lint vet errcheck

run:
	maven_repo_server \
	-logtostderr \
	-v=2 \
	-port=8080 \
	-root=/tmp

open:
	open http://localhost:7777/storage/

format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install

clean:
	rm -rf vendor var
