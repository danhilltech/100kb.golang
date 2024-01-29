RUST_SRC = $(shell find . -type f -name '*.rs' -not -path "./target/*" -o -name '*.proto' -not -path "./target/*")
GO_PROTO_SRC = $(shell find ./lib -type f -name '*.proto' -not -path "./target/*")
REV := $(shell git rev-parse HEAD)

lib/libgobert.so: $(RUST_SRC)
	rm target/release/libgobert.so || true
	cargo build --release
	mkdir -p lib
	@cp target/release/libgobert.so lib/libgobert.so

lib/libgoadblock.so: $(RUST_SRC)
	rm target/release/libgoadblock.so || true
	cargo build --release
	mkdir -p lib
	@cp target/release/libgoadblock.so lib/libgoadblock.so

# lib/gobert-cbindgen.h: $(RUST_SRC)
# 	@cd lib/gobert && cbindgen . --lang c -o ../gobert-cbindgen.h

pkg/ai/keywords.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/gobert/src/keywords.proto
pkg/ai/sentence_embedding.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/gobert/src/sentence_embedding.proto
pkg/serialize/article.pb.go: pkg/serialize/article.proto
	@protoc -I=. --go_out=. ./pkg/serialize/article.proto
pkg/parsing/adblock.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/goadblock/src/adblock.proto

.PHONY: build
build: lib/libgobert.so lib/libgoadblock.so lib/gobert-cbindgen.h pkg/ai/keywords.pb.go pkg/ai/sentence_embedding.pb.go pkg/serialize/article.pb.go pkg/parsing/adblock.pb.go
	go build -ldflags="-r $(ROOT_DIR)lib" -buildvcs=false

.PHONY: clean
clean:	
	@go clean
	@rm -f 100kb.golang
	@rm -rf target/*
	@rm -rf .cache/*
	@rm -rf dbs/output*

.PHONY: godefs
godefs:
	@go tool cgo -godefs ai.go

.PHONY: dockerbuild
dockerbuild:
	docker build --tag '100kb.golang' .

.PHONY: dockerbuild
dockerterm:
	docker run --gpus all --rm --mount type=bind,source=./dbs,target=/dbs --mount type=bind,source=./.cache,target=/cache -it --entrypoint zsh  100kb.golang

.PHONY: index
index:
	docker run --dns=1.1.1.1 --gpus all --mount type=bind,source=./dbs,target=/dbs --mount type=bind,source=./.cache,target=/cache 100kb.golang -mode=index -http-chunk-size=500 -hn-fetch-size=1000000 --cache-dir=/cache > log.txt 2>&1

.PHONY: meta
meta:
	docker run --dns=1.1.1.1 --gpus all --mount type=bind,source=./dbs,target=/dbs --mount type=bind,source=./.cache,target=/cache 100kb.golang -mode=meta -util=0.8 -http-chunk-size=500 -hn-fetch-size=1000000 --cache-dir=/cache > log.txt 2>&1



.PHONY: output
output:
	docker run --dns=1.1.1.1 --gpus all --mount type=bind,source=./dbs,target=/dbs --mount type=bind,source=./.cache,target=/cache --mount type=bind,source=./output,target=/app/output -p 8080:8080 100kb.golang -mode=output > log.txt 2>&1


.PHONY: transfer
transfer:
	rsync -havP --stats dan@192.168.1.3:~/100kb.golang/dbs/output.db .