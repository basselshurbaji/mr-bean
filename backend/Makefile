.PHONY: build test migrate clean

build:
	@bash scripts/build.sh

run:
	@go run ./cmd/server

test:
	@bash scripts/build.sh
	@bash scripts/test.sh

migrate:
	@bash scripts/migrate.sh

clean:
	@bash scripts/clean.sh