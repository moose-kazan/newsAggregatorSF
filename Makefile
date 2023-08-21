.PHONY:

test:
	cd gateway && make test
	cd serviceComments && make test
	cd serviceNews && make test
	cd serviceModerate && make test

build:
	docker-compose build

run: build
	docker-compose up

stop:
	docker-compose down
