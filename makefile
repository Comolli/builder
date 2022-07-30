GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

.PHONY: coverage-report
# coverage-report
coverage-report: 
	go test -race -covermode=atomic -coverprofile=coverage.out
	go tool cover -html=coverage.out	
