#
# Simple Makefile for Golang based Projects.
#
PROJECT = wsfn

GIT_GROUP = caltechlibrary

RELEASE_DATE=$(shell date +'%Y-%m-%d')

RELEASE_HASH=$(shell git log --pretty=format:'%h' -n 1)

PROGRAMS = $(shell ls -1 cmd)

MAN_PAGES_1 = $(shell ls -1 *.1.md | sed -E 's/\.1.md/.1/g')

MAN_PAGES_3 = $(shell ls -1 *.3.md | sed -E 's/\.3.md/.3/g')

HTML_PAGES = $(shell find . -type f | grep -E '\.html')

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PACKAGE = $(shell ls -1 *.go | grep -v 'version.go')

SUBPACKAGES = $(shell ls -1 */*.go)

OS = $(shell uname)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

ifneq ($(prefix),)
	PREFIX = $(prefix)
endif

EXT =
ifeq ($(OS), Windows)
	EXT = .exe
endif

DIST_FOLDERS = bin/*

build: version.go $(PROGRAMS) man CITATION.cff about.md installer.sh installer.ps1

version.go: .FORCE
	cmt codemeta.json version.go	

$(PROGRAMS): $(PACKAGE)
	@mkdir -p bin
	go build -o "bin/$@$(EXT)" cmd/$@/*.go
	@./bin/$@ -help >$@.1.md

man: $(MAN_PAGES_1) $(MAN_PAGES_3)

man/man1:
	mkdir -p man/man1

man/man3:
	mkdir -p man/man3

$(MAN_PAGES_1): man/man1 .FORCE
	pandoc $@.md --from markdown --to man -s >man/man1/$@

$(MAN_PAGES_3): man/man3 .FORCE
	pandoc $@.md --from markdown --to man -s >man/man3/$@

CITATION.cff: .FORCE
	cmt codemeta.json CITATION.cff

about.md: .FORCE 
	cmt codemeta.json about.md

installer.sh: .FORCE
	@cmt codemeta.json installer.sh
	@chmod 775 installer.sh
	@git add -f installer.sh

installer.ps1: .FORCE
	@cmt codemeta.json installer.ps1
	@chmod 775 installer.ps1
	@git add -f installer.ps1

clean-website:
	make -f website.mak clean

website: clean-website .FORCE
	make -f website.mak


# NOTE: on macOS you must use "mv" instead of "cp" to avoid problems
install: build man .FORCE
	@if [ ! -d $(PREFIX)/bin ]; then mkdir -p $(PREFIX)/bin; fi
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f ./bin/$$FNAME ]; then mv -v ./bin/$$FNAME $(PREFIX)/bin/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"
	@echo ""
	@if [ ! -d $(PREFIX)/man/man1 ]; then mkdir -p $(PREFIX)/man/man1; fi
	@for MAN_PAGE in $(MAN_PAGES); do cp -v man/man1/$$MAN_PAGE $(PREFIX)/man/man1/;done
	@echo ""
	@echo "Make sure $(PREFIX)/man is in your MANPATH"
	@echo ""

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME ]; then rm -v $(PREFIX)/bin/$$FNAME; fi; done
	@echo "Removing man pages from  $(PREFIX)/man/man1"
	@for MAN_PAGE in $(MAN_PAGES_1); do if [ -f "$(PREFIX)/man/man1/$$MAN_PAGE" ]; then rm "$(PREFIX)/man/man1/$$MAN_PAGE"; fi; done
	@echo "Removing man pages from  $(PREFIX)/man/man3"
	@for MAN_PAGE in $(MAN_PAGES_3); do if [ -f "$(PREFIX)/man/man3/$$MAN_PAGE" ]; then rm "$(PREFIX)/man/man3/$$MAN_PAGE"; fi; done


hash: .FORCE
	git log --pretty=format:'%h' -n 1

check: .FORCE
	for FNAME in $(shell ls -1 *.go); do go fmt $$FNAME; done
	go vet *.go

test: clean build
	go test

clean: 
	-if [ -d bin ]; then rm -fR bin; fi
	-if [ -d dist ]; then rm -fR dist; fi
	-if [ -d testout ]; then rm -fR testout; fi
	-for MAN_PAGE in $(MAN_PAGES); do if [ -f man/man1/$$MAN_PAGE.1 ]; then rm man/man1/$$MAN_PAGE.1; fi;done

status:
	git status

save:
	@if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

refresh:
	git fetch origin
	git pull origin $(BRANCH)

publish: build website save .FORCE
	./publish.bash


dist/Linux-x86_64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env  GOOS=linux GOARCH=amd64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/Linux-aarch64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env  GOOS=linux GOARCH=arm64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-aarch64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/macOS-x86_64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin


dist/macOS-arm64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-arm64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin


dist/Windows-x86_64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o "dist/bin/$${FNAME}.exe" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* 
	@rm -fR dist/bin

dist/Windows-arm64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=arm64 go build -o "dist/bin/$${FNAME}.exe" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-arm64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/RaspberryPiOS-arm7: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-RaspberryPiOS-arm7.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

distribute_docs:
	@mkdir -p dist/
	@cp -v codemeta.json dist/
	@cp -v CITATION.cff dist/
	@cp -v README.md dist/
	@cp -v LICENSE dist/
	@cp -v INSTALL.md dist/
	@cp -vR man dist/

release: save build save distribute_docs dist/Linux-x86_64 dist/Linux-aarch64 dist/macOS-x86_64 dist/macOS-arm64 dist/Windows-x86_64 dist/Windows-arm64 dist/RaspberryPiOS-arm7
	echo "Ready to do ./release.bash"

.FORCE:
