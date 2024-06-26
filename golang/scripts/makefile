
WORKDIR ?= $(shell pwd)/..
DOT_ENV ?= ${WORKDIR}/internal/config/.env
include ${DOT_ENV}
export $(shell sed 's/=.*//' ${DOT_ENV})

default:
	go build -mod=vendor -o ${WORKDIR}/bin/main ${WORKDIR}/cmd/main.go
	${WORKDIR}/bin/main