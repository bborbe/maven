install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/maven_repo_server/maven_repo_server.go
test:
	GO15VENDOREXPERIMENT=1 go test `glide novendor`
check:
	golint ./...
	errcheck -ignore '(Close|Write)' ./...
run:
	maven_repo_server \
	-loglevel=INFO \
	-port=8080 \
	-root=/tmp
open:
	open http://localhost:7777/storage/
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf vendor var
