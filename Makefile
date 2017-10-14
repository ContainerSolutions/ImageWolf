# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash
GO := go
BINARY = registry-x86_64
GOARCH = amd64
LDFLAGS= -ldflags '-extldflags "-static"'
BUILDARGS = -v ${LDFLAGS} -o ./bin/${BINARY} ./src
LOG= ImageWolf.log
DEPENDENCIES := 	github.com/anacrolix/torrent \
	github.com/anacrolix/utp \
	github.com/docker/distribution/notifications \
	github.com/gorilla/mux

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
	+ GOOS=linux CGO_ENABLED=0 GOARCH=${GOARCH} go build ${BUILDARGS}

.PHONY: deps
deps:
	+ go get -v $(DEPENDENCIES)

.PHONY: env
env:
						go env

#	EOF!
