.PHOMY: build

clean:
	rm -rfv build

build:
	mkdir -p build
	go build -o build/gateway

deps:
	go get github.com/gorilla/mux

test:
	@echo Run tests
	go test -v -covermode=count './...'
