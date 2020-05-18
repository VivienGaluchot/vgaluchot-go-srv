# build tool

TMP := tmp/build

.PHONY: all build test run clean

all:build

build:
	mkdir -p $(TMP)
	go build -v -o $(TMP)/vgaluchot-go-srv ./cmd/vgaluchot-go-srv

test:
	@export WEB_DIR="../../web" ; go test -v ./cmd/vgaluchot-go-srv
	@export WEB_DIR="../../web" ; go test -v ./deployments/gae
	@export WEB_DIR="../../../../web" ; go test -v ./internal/app/vgaluchot-go-srv/server

run: build
	cd web ; ../$(TMP)/vgaluchot-go-srv

clean:
	rm -rf $(TMP)