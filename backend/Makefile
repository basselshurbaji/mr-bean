BINARY := bin/server

.PHONY: build test migrate clean

build:
	@zsh scripts/build.sh
	@go build -o $(BINARY) ./cmd/server

test:
	@zsh scripts/test.sh

migrate:
	@zsh scripts/migrate.sh

clean:
	@rm -rf bin/
