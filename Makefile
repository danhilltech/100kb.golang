RUST_SRC = $(shell find . -type f -name '*.rs' -not -path "./target/*" -o -name '*.proto' -not -path "./target/*")
GO_PROTO_SRC = $(shell find ./lib -type f -name '*.proto' -not -path "./target/*")
REV := $(shell git rev-parse HEAD)

lib/libgobert.so: $(RUST_SRC)
	cargo build --release
	mkdir -p lib
	@cp target/release/libgobert.so lib/libgobert.so

# lib/gobert-cbindgen.h: $(RUST_SRC)
# 	@cd lib/gobert && cbindgen . --lang c -o ../gobert-cbindgen.h

pkg/ai/keywords.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/gobert/src/keywords.proto
pkg/ai/sentence_embedding.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/gobert/src/sentence_embedding.proto
pkg/serialize/article.pb.go: pkg/serialize/article.proto
	@protoc -I=. --go_out=. ./pkg/serialize/article.proto

.PHONY: build
build: lib/libgobert.so lib/gobert-cbindgen.h pkg/ai/keywords.pb.go pkg/ai/sentence_embedding.pb.go pkg/serialize/article.pb.go
	go build -ldflags="-r $(ROOT_DIR)lib" -buildvcs=false

.PHONY: clean
clean:	
	@go clean
	@rm -f 100kb.golang
	@rm -rf target/*

.PHONY: godefs
godefs:
	@go tool cgo -godefs ai.go

.PHONY: dockerbuild
dockerbuild:
	docker build --tag '100kb.golang' .

.PHONY: dockerbuild
dockerterm:
	docker run --gpus all --rm --mount type=bind,source=./dbs,target=/dbs -it 100kb.golang zsh 

.PHONY: run
run:
	docker run --dns=1.1.1.1 --gpus all --mount type=bind,source=./dbs,target=/dbs 100kb.golang -mode=index -http-workers=50 -http-chunk-size=500 -hn-fetch-size=1000000 > log.txt 2>&1

.PHONY: output
output:
	docker run --dns=1.1.1.1 --gpus all --mount type=bind,source=./dbs,target=/dbs 100kb.golang -mode=output > log.txt 2>&1

.PHONY: debug
debug:
	docker run --dns=1.1.1.1 -p 6060:6060 --gpus all --mount type=bind,source=./dbs,target=/dbs 100kb.golang -http-workers=50 -http-chunk-size=500 -hn-fetch-size=1000000 -debug=true > log.txt 2>&1


.PHONY: transfer
transfer:
	rsync -havP --stats dan@192.168.1.3:~/100kb.golang/dbs/output.db .