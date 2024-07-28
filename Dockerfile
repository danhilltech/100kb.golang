FROM nvidia/cuda:11.8.0-devel-ubuntu22.04

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10

COPY --from=golang:1.22-alpine /usr/local/go/ /usr/local/go/

ARG USERNAME=builder
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG LIBTORCH_URL=https://download.pytorch.org/libtorch/cu118/libtorch-cxx11-abi-shared-with-deps-2.2.2%2Bcu118.zip
ARG CHROME_URL=https://storage.googleapis.com/chrome-for-testing-public/129.0.6622.0/linux64/chrome-linux64.zip

ENV TORCH_HOME=/usr/local/lib/libtorch
ENV LIBTORCH=/usr/local/lib/libtorch
ENV LD_LIBRARY_PATH=${LIBTORCH}/lib:$LD_LIBRARY_PATH
ENV DEBIAN_FRONTEND=noninteractive

# ENV LIBTORCH_URL=https://download.pytorch.org/libtorch/cpu/libtorch-macos-2.1.2.zip
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
ENV PATH="/home/builder/.cargo/bin:/usr/local/go/bin:/go/bin:${PATH}"
ENV LIBTORCH_BYPASS_VERSION_CHECK="true"
ENV LD_LIBRARY_PATH="/usr/local/cuda/lib64:/usr/local/lib:${LD_LIBRARY_PATH}"
ENV DEBIAN_FRONTEND=noninteractive
ENV NVIDIA_VISIBLE_DEVICES all
ENV NVIDIA_DRIVER_CAPABILITIES all

RUN apt-get update \
    && apt-get install -y software-properties-common \
    && apt-get update -y --no-install-recommends \
    && apt-get -y install --no-install-recommends bash git wget curl tzdata build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev cmake unzip zsh ca-certificates sudo apt-transport-https nano zip openssh-client apt-utils pkg-config gcc protobuf-compiler \
    && apt-get autoremove -y

# CHROME
RUN curl -fsSL -o chrome.zip $CHROME_URL \
    && unzip -q chrome.zip -d /chrome \
    && rm chrome.zip

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    echo "${USERNAME} ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers && \
    chmod 0440 /etc/sudoers && \
    chmod g+w /etc/passwd 

RUN chsh -s /usr/bin/zsh builder

WORKDIR /usr/local/lib

RUN curl -fsSL --insecure -o libtorch.zip $LIBTORCH_URL \
    && unzip -q libtorch.zip \
    && rm libtorch.zip


WORKDIR /app

RUN chown -R builder:builder /app
RUN chmod -R 755 /app
RUN mkdir -p /go && chmod -R 777 /go
RUN git config --global --add safe.directory /app

USER builder

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

ENV TERM xterm
# set the zsh theme
ENV ZSH_THEME agnoster

ENV RUSTBERT_CACHE /app/models
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/app/lib

# run the installation script  
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true

COPY --chown=builder:builder --chmod=755 ./lib/ ./lib/
COPY --chown=builder:builder --chmod=755 ./pkg/ ./pkg/
COPY --chown=builder:builder --chmod=755 ./*.go .
COPY --chown=builder:builder --chmod=755 ./Cargo.* .
COPY --chown=builder:builder --chmod=755 ./go.mod .
COPY --chown=builder:builder --chmod=755 ./go.sum .
COPY --chown=builder:builder --chmod=755 ./Makefile .

# # Build the Go app
RUN --mount=type=cache,id=rustcache,target=/usr/local/cargo/registry,uid=1000,gid=1000 \
    --mount=type=cache,id=rustbuild,target=/app/target,uid=1000,gid=1000 \
    --mount=type=cache,id=gomod,target=/go/pkg/mod,uid=1000,gid=1000 \
    --mount=type=cache,id=gobuild,target=/home/builder/.cache/go-build,uid=1000,gid=1000 \
    --mount=type=cache,id=gobuildtmp,target=/tmp/go-build,uid=1000,gid=1000 \
    make build

# # # Run the executable
ENTRYPOINT ["./100kb.golang"]