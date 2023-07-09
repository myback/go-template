lib: test
	go build -buildmode=c-shared -ldflags='-s -w' -o go_template/bind/template.so go/lib/template.go

binary: test
	CGO_ENABLED=0 go build -ldflags='-s -w' -o bin/go-template go/template.go

test:
	PYTHONPATH="$(PWD):$(PYTHONPATH)" python3 tests/test_utils.py
	@echo ''
	go test ./go/template-functions
