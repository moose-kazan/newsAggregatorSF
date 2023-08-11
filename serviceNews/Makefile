.PHOMY: build

clean:
	rm -rfv build

build: deps
	mkdir -p build
	go build -o build/serviceNews

deps:
	go get github.com/gorilla/mux
	go get github.com/jackc/pgx/v4/pgxpool
