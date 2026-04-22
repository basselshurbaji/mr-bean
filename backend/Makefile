BINARY := bin/server

.PHONY: build test clean

build:
	@zsh scripts/build.sh
	@go build -o $(BINARY) ./cmd/server

test:
	@zsh scripts/test.sh

clean:
	@rm -rf bin/
