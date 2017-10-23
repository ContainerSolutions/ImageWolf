# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash
GO := go
BINARY = ImageWolf
TESTBINARY = healt
GOARCH = amd64
LDFLAGS= -ldflags '-extldflags "-static"'
BUILDARGS = -v ${LDFLAGS} -o ./bin/${BINARY} ./src
TESTBUILDARGS = -v ${LDFLAGS} -o ./bin/${TESTBINARY} ./test
LOG= ImageWolf.log
DEPENDENCIES := 	github.com/anacrolix/torrent \
	github.com/anacrolix/utp \
	github.com/docker/distribution/notifications \
	github.com/gorilla/mux \
	github.com/anacrolix/torrent/bencode \
	github.com/anacrolix/torrent/metainfo \
	github.com/docker/distribution/notifications \
	github.com/docker/docker/api/types \
	github.com/docker/docker/client \
	golang.org/x/net/context

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
	@if [ -f ./bin/${BINARY} ] ; then rm ./bin/${BINARY} ; fi
	+	GOOS=linux CGO_ENABLED=0 GOARCH=${GOARCH} go build ${BUILDARGS}

.PHONY: deps
deps:
	+	go get -v -d $(DEPENDENCIES)

.PHONY: env
env:
	+	go env

#	EOF!
