all: build test

gen:
	go fmt ./...

build:
	sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o marathon-pipeline'

install: all do_install

do_install:
	go install

test:
	go test -cover ./...
	go vet ./...

docker:
	docker run --rm -it -e "GOPATH=/go" -v "${PWD}:/go/src/github.com/kpacha/marathon-pipeline" -w /go/src/github.com/kpacha/marathon-pipeline golang:1.5.1 make
	docker build -f Dockerfile-min -t kpacha/marathon-pipeline:latest-min .