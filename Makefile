.PHONY: build-app build-poller run-app run-poller clean migrate-up migrate-down

all: build-app build-poller

build-app:
	go build -o bin/app cmd/app/main.go

build-poller:
	go build -o bin/poller cmd/poller/main.go

run-app: build-app
	./bin/app

run-poller: build-poller
	./bin/poller

clean:
	rm -rf bin/*

migrate-up:
	goose sqlite3 bigmacindex.db -dir migrations up

migrate-down:
	goose sqlite3 bigmacindex.db -dir migrations down