#
# Simple Makefile
#

PROJECT = wsfn

VERSION = $(shell grep -m1 "Version = " $(PROJECT).go | cut -d\` -f 2)

BRANCH = $(shell git branch | grep "* " | cut -d\   -f 2)

PKGASSETS = $(shell which pkgassets)

OS = $(shell uname)

EXT =
ifeq ($(OS),Windows)
	EXT = .exe
endif

build: bin/webserver$(EXT) bin/webaccess$(EXT)

bin/webserver$(EXT): access.go cors.go \
	defaults.go json.go license.go logger.go \
	redirects.go service.go wsfn.go \
	safefilesystem.go \
	cmd/webserver/webserver.go
	go build -o bin/webserver$(EXT) cmd/webserver/webserver.go

bin/webaccess$(EXT): access.go cors.go \
	defaults.go json.go license.go logger.go \
	redirects.go service.go wsfn.go \
	safefilesystem.go \
	cmd/webaccess/webaccess.go
	go build -o bin/webaccess$(EXT) cmd/webaccess/webaccess.go

lint:
	golint access.go
	golint cors.go
	golint cors_test.go
	golint defaults.go
	golint json.go
	golint json_test.go
	golint license.go
	golint logger.go
	golint redirects.go
	golint service.go
	golint wsfn.go
	golint wsfn_test.go
	goling safefilesystem.go
	golint cmd/webserver/webserver.go
	golint cmd/webaccess/webaccess.go

format:
	gofmt -w access.go
	gofmt -w cors.go
	gofmt -w cors_test.go
	gofmt -w defaults.go
	gofmt -w json.go
	gofmt -w json_test.go
	gofmt -w license.go
	gofmt -w logger.go
	gofmt -w redirects.go
	gofmt -w service.go
	gofmt -w wsfn.go
	gofmt -w wsfn_test.go  
	gofmt -w safefilesystem.go
	gofmt -w cmd/webserver/webserver.go
	gofmt -w cmd/webaccess/webaccess.go

test: bin/webserver$(EXT) bin/webaccess$(EXT)
	go test

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi

install: 
	env GOBIN=$(GOPATH)/bin go install cmd/webserver/webserver.go
	env GOBIN=$(GOPATH)/bin go install cmd/webaccess/webaccess.go


dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/webserver cmd/webserver/webserver.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/webaccess cmd/webaccess/webaccess.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/* docs/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/webserver.exe cmd/webserver/webserver.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/webaccess.exe cmd/webaccess/webaccess.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/* docs/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/webserver cmd/webserver/webserver.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/webaccess cmd/webaccess/webaccess.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md bin/* docs/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/webserver cmd/webserver/webserver.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/webaccess cmd/webaccess/webaccess.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/* docs/*
	rm -fR dist/bin

distribute_docs:
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -vR docs dist/

release: clean website wsfn.go cmd/webserver/webserver.go distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7

website:
	./mk_website.py

publish:
	./mk_website.py
	./publish.bash

FORCE:
