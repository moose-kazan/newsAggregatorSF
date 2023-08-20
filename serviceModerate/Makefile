.PHOMY: build

clean:
	rm -rfv build

build: deps
	mkdir -p build
	go build -o build/serviceModerate
	cp run.sh build/run.sh

deps:
	go get github.com/gorilla/mux

test:
	@echo Run tests
	go test -v -covermode=count './...'

