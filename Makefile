SHELL := /bin/bash

install:
	go get -u github.com/kardianos/govendor
	rm -rf vendor
	govendor init
	govendor add +external
	go build ./...
	go install ./...
	devops-tool version