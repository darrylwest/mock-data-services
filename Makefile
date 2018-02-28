export GOPATH:=$(HOME)/.gopath:$(PWD)
VERSION=`cat VERSION`

.PHONY: test build format lint

test:
	@( go vet src/*/*.go )
	@( cd test/unit && go test )
	@( make lint )

version:
	@( echo $(VERSION) )

build:
	@[ -d bin ] || mkdir bin
	go build -o bin/mack-data-proxy src/main.go

install-deps:
	go get -u github.com/shuLhan/go-bindata/...
	go get github.com/go-zoo/bone
	go get github.com/golang/lint/golint
	go get github.com/franela/goblin
	go get -u github.com/darrylwest/go-unique/unique
	go get github.com/boltdb/bolt
	go get github.com/darrylwest/cassava-logger/logger

format:
	( gofmt -s -w src/*.go src/*/*.go test/*/*.go )

lint:
	@( golint src/... && golint test/... )

watch:
	go-watcher --loglevel=5

integration:
	go run test/integration/db-utils.go

edit:
	make format
	vi -O3 src/*/*.go test/*/*.go src/*.go

