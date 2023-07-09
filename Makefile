.PHONY: lib
lib: go-test
	go build -buildmode=c-shared -ldflags='-s -w' -o go_template/bind/template.so go/lib/template.go

.PHONY: bin
bin: go-test
	CGO_ENABLED=0 go build -ldflags='-s -w' -o bin/go-template go/template.go

.PHONY: go-test
go-test:
	go test ./go/template-functions

.PHONY: py-test
py-test:
	PYTHONPATH="$(PWD):$(PYTHONPATH)" python tests/test_utils.py

.PHONY: py-build
py-build: py-test
	python setup.py bdist_wheel

.PHONY: build-on-host
build-on-host: lib py-test py-build

.PHONY: lib-in-docker
lib-in-docker:
	docker run --platform=linux/amd64 --rm -it -w /app -v ${PWD}:/app -v ${PWD}/.go_cache:/go/pkg golang:bullseye make lib

.PHONY: build-in-docker
build-in-docker: lib-in-docker
	docker run --platform=linux/amd64 -e PYTHONPATH='/app' --rm -it -w /app -v ${PWD}:/app python:bullseye make py-build

.PHONY: clean
clean:
	@rm -rf .pytest_cache/ build/ dist/
	@find . -not -path './.venv*' -path '*/__pycache__*' -delete
	@find . -not -path './.venv*' -path '*/*.egg-info*' -delete
