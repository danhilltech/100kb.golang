RUST_SRC = $(shell find . -type f -name '*.rs' -not -path "./target/*")
REV := $(shell git rev-parse HEAD)

lib/libgobert.so: $(RUST_SRC)
	cargo build --release
	@cp target/release/libgobert.so lib/libgobert.so

lib/gobert-cbindgen.h: $(RUST_SRC)
	@cd lib/gobert && cbindgen . --lang c -o ../gobert-cbindgen.h

.PHONY: build
build: lib/libgobert.so lib/gobert-cbindgen.h
	go build -ldflags="-r $(ROOT_DIR)lib" -buildvcs=false

.PHONY: clean
	go clean
	@rm -f 100kb.golang
