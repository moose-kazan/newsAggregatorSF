.PHOMY: build

clean:
	rm -rfv build

build: deps
	mkdir -p build
	go build -o build/serviceNews
	cp run.sh build/run.sh

deps:
	go get github.com/gorilla/mux
	go get github.com/jackc/pgx/v4/pgxpool

test:
	@echo Prepare empty DB
	docker run \
		-p 15432:5432 \
		--rm \
		--name postgre-server-test \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=test \
		-d postgres
	sleep 3
	@echo Load schema
	psql postgresql://postgres:postgres@127.0.0.1:15432/test < db_schema.sql

	@echo Run tests
	go test -v -covermode=count './...'

	@echo Cleanup
	docker stop postgre-server-test
