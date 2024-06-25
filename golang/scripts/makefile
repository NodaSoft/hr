include ./.env

export $(shell sed 's/=.*//' .env)

default:
	go build -mod=vendor -o ../bin/main ../cmd/main.go
	../bin/main