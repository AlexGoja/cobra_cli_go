SHELL := /bin/bash

install:
	/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
	brew install make
	go get -u github.com/kardianos/govendor
	rm -rf vendor
	govendor init
	govendor add +external
	go build ./...
	go install ./...
	devops-tool version