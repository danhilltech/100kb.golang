RUST_SRC = $(shell find ./lib -type f -name '*.rs' -not -path "./target/*" -o -name '*.proto' -not -path "./target/*")
GO_PROTO_SRC = $(shell find ./lib -type f -name '*.proto' -not -path "./target/*")
GO_SVM_SRC = $(shell find ./pkg/svm -type f -name '*.cpp')
REV := $(shell git rev-parse HEAD)
HAS_CUDA = ${shell command -v nvidia-smi}
DOCKER_TAG = danhilltech/100kb
DOCKER_CORE_ARGS = --cap-add=SYS_ADMIN --dns=1.1.1.1 --mount type=bind,source=./dbs,target=/dbs --mount type=bind,source=./.cache,target=/cache  --mount type=bind,source=./models,target=/app/models --mount type=bind,source=./output,target=/app/output --mount type=bind,source=./train,target=/train -p 9081:8081
DOCKER_RUN_ARGS = -cache-dir=/cache
DOCKER_GPUS = ${shell if command -v nvidia_smi >&/dev/null; then echo "--gpus all"; fi}

GO_BUILD_TAGS :=
CUDA :=
LIBTORCH_URL := https://download.pytorch.org/libtorch/cpu/libtorch-cxx11-abi-shared-with-deps-2.2.2%2Bcpu.zip

ifneq (, ${HAS_CUDA})
CUDA = 1
GO_BUILD_TAGS += cuda
LIBTORCH_URL = https://download.pytorch.org/libtorch/cu118/libtorch-cxx11-abi-shared-with-deps-2.2.2%2Bcu118.zip
endif


.PHONY: debug
debug:
	@echo "🍔🍔🍔🍔🍔"
	@echo "CUDA: ${CUDA}"
	@echo "Build Tags: ${GO_BUILD_TAGS}"
	@echo "Docker GPUs: ${DOCKER_GPUS}"
	@echo "Has Cuda: ${HAS_CUDA}"
ifdef CUDA
	@echo "Will build for cuda"
endif

lib/libgobert.so: $(RUST_SRC)
ifdef CUDA
	@echo "👉 Building libgobert"
	@rm -f target/release/libgobert.so || true
	cargo build --release --package gobert
	@mkdir -p lib
	@cp target/release/libgobert.so lib/libgobert.so
endif

lib/libgoadblock.so: $(RUST_SRC)
	@echo "👉 Building libgoadblock"
	@rm -f target/release/libgoadblock.so || true
	cargo build --release --package goadblock
	@mkdir -p lib
	@cp target/release/libgoadblock.so lib/libgoadblock.so



pkg/ai/keywords.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/gobert/src/keywords.proto
pkg/ai/sentence_embedding.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/gobert/src/sentence_embedding.proto
pkg/ai/zero_shot.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/gobert/src/zero_shot.proto
pkg/serialize/article.pb.go: pkg/serialize/article.proto
	@protoc -I=. --go_out=. ./pkg/serialize/article.proto
pkg/parsing/adblock.pb.go: ${GO_PROTO_SRC}
	@protoc -I=. --go_out=. ./lib/goadblock/src/adblock.proto


.PHONY: build
build: lib/libgobert.so lib/libgoadblock.so pkg/ai/keywords.pb.go pkg/ai/sentence_embedding.pb.go pkg/serialize/article.pb.go pkg/parsing/adblock.pb.go pkg/ai/zero_shot.pb.go
	@echo "👉 Building go binary"
	go build -ldflags="-r $(ROOT_DIR)lib" -tags '$(GO_BUILD_TAGS)'

.PHONY: release
release: clean debug build

.PHONY: clean
clean:
	@echo "🧹 Cleaning"
	@go clean
	@rm -f 100kb.golang
	@rm -rf target
	@rm -rf .cache
	@rm -rf dbs/output*
	@rm -rf lib/*.so
	@rm -rf pkg/svm/libsvm.so
	@mkdir .cache
	@mkdir target

.PHONY: godefs
godefs:
	@go tool cgo -godefs ai.go

.PHONY: dockerbuild
dockerbuild:
	docker build --tag '${DOCKER_TAG}' --build-arg="LIBTORCH_URL=${LIBTORCH_URL}" .

.PHONY: dockerterm
dockerterm:
	docker run ${DOCKER_GPUS} ${DOCKER_CORE_ARGS} --rm  -it --entrypoint zsh ${DOCKER_TAG}

.PHONY: search
search:
	docker run ${DOCKER_GPUS} ${DOCKER_CORE_ARGS} ${DOCKER_TAG} -mode=search -http-chunk-size=200 -hn-fetch-size=1000000 ${DOCKER_RUN_ARGS}

.PHONY: index
index:
	docker run ${DOCKER_GPUS} ${DOCKER_CORE_ARGS} ${DOCKER_TAG} -mode=index -http-chunk-size=200 ${DOCKER_RUN_ARGS}


.PHONY: meta
meta:
	docker run ${DOCKER_GPUS} ${DOCKER_CORE_ARGS} ${DOCKER_TAG} -mode=meta ${DOCKER_RUN_ARGS}

.PHONY: train
train:
	docker run ${DOCKER_GPUS} ${DOCKER_CORE_ARGS} ${DOCKER_TAG} -mode=train -train-dir=/train ${DOCKER_RUN_ARGS}

.PHONY: output
output:
	docker run ${DOCKER_GPUS} ${DOCKER_CORE_ARGS} ${DOCKER_TAG} -mode=output ${DOCKER_RUN_ARGS}

.PHONY: transfer
transfer:
	rsync -havP --stats dan@192.168.1.3:~/100kb.golang/dbs/output.db .