# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash
GO := go

BINARY = registry-x86_64
GOARCH = amd64
LDFLAGS= -ldflags '-extldflags "-static"'
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
	@if [ -f ./bin/${BINARY} ] ; then rm ./bin/${BINARY}-* ; fi
	+ GOOS=linux CGO_ENABLED=0 GOARCH=${GOARCH} go build ${LDFLAGS} -o ./bin/${BINARY} ./src

.PHONY: test
test:
	@if [ -f ./bin/${BINARY} ] ; then rm ./bin/${BINARY} ; fi
	+ GOOS=linux CGO_ENABLED=0 GOARCH=${GOARCH} go build ${LDFLAGS} -o ./bin/${BINARY} ./src
	@if [ -f ./bin/${BINARY} ] ;	then return 0; else	return 1;	fi

.PHONY: deps
deps:
			go get $(DEPENDENCIES)

.PHONY: env
env:
						go env

#	EOF!
