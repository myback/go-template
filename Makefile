PLATFORM ?=

ifeq ($(OS),Windows_NT)
	PLATFORM := win_amd64
else
	UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
    	PLATFORM := manylinux2014-x86_64
    endif
    ifeq ($(UNAME_S),Darwin)
        PLATFORM := macosx-11-0-arm64
    endif
endif

.PHONY: lib
lib: go-test
	cd src/go && go build -buildmode=c-shared -ldflags='-s -w' -o ../go_template/bind/template.so lib/template.go

.PHONY: bin
bin: go-test
	CGO_ENABLED=0 go build -ldflags='-s -w' -o bin/go-template src/go/template.go

.PHONY: go-test
go-test:
	cd src/go && go test ./tmpl-functions

.PHONY: py-test
py-test:
	python -m unittest discover

.PHONY: tox
tox:
	tox -- --plat-name=$(PLATFORM)

.PHONY: dist
dist:
	python -m build -w -C="--global-option=--plat-name=$(PLATFORM)"
