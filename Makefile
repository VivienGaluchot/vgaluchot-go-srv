# build tool

TMP := tmp/build

all:build

build:
	mkdir -p $(TMP)
	go build -v -o $(TMP)/vgaluchot-go-srv ./cmd/vgaluchot-go-srv

run: build
	cd web ; ../$(TMP)/vgaluchot-go-srv

clean:
	rm -rf $(TMP)