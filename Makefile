# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash
GO := go
BINARY = registry-x86_64
GOARCH = amd64
VENDOR=${GOPATH}/src/github.com/ArangoGutierrez/ImageWolf/vendor
LDFLAGS= -ldflags '-extldflags "-static"'
BUILDARGS = -v ${LDFLAGS} -o ./bin/${BINARY} ./src
LOG= log
DEPENDENCIES := 	github.com/anacrolix/torrent \
	github.com/anacrolix/utp \
	github.com/docker/distribution/notifications

.PHONY: clean
clean:
		@if [ -f bin/${BINARY} ] ; then rm bin/${BINARY} ; fi
		@if [ -f ${LOG} ] ; then rm ${LOG} ; fi

.PHONY: build
build:
	@if [ -f ./bin/${BINARY} ] ; then rm ./bin/${BINARY} ; fi
	+ GOOS=linux CGO_ENABLED=0 GOARCH=${GOARCH} go build ${BUILDARGS}

.PHONY: test
test:
	@if [ -f ./bin/${BINARY} ] ; then rm ./bin/${BINARY} ; fi
	+ GOOS=linux CGO_ENABLED=0 GOARCH=${GOARCH} GOPATH=${VENDOR} go build ${BUILDARGS}

.PHONY: deps
deps:
	+	GOPATH=${VENDOR} go get -v $(DEPENDENCIES)

.PHONY: env
env:
						go env

#	EOF!
