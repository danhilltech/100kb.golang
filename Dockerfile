FROM nvidia/cuda:11.8.0-devel-ubuntu22.04

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10

COPY --from=golang:1.21-alpine /usr/local/go/ /usr/local/go/

ARG USERNAME=builder
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV TORCH_HOME=/usr/local/lib/libtorch
ENV LIBTORCH=/usr/local/lib/libtorch
ENV LD_LIBRARY_PATH=${LIBTORCH}/lib:$LD_LIBRARY_PATH
ENV DEBIAN_FRONTEND=noninteractive
ENV LIBTORCH_URL=https://download.pytorch.org/libtorch/cu118/libtorch-cxx11-abi-shared-with-deps-2.1.2%2Bcu118.zip
# ENV LIBTORCH_URL=https://download.pytorch.org/libtorch/cpu/libtorch-macos-2.1.2.zip
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
ENV PATH="/home/builder/.cargo/bin:/usr/local/go/bin:/go/bin:${PATH}"
ENV LIBTORCH_BYPASS_VERSION_CHECK="true"
ENV LD_LIBRARY_PATH="/usr/local/cuda/lib64:/usr/local/lib:${LD_LIBRARY_PATH}"
ENV DEBIAN_FRONTEND=noninteractive
ENV NVIDIA_VISIBLE_DEVICES all
ENV NVIDIA_DRIVER_CAPABILITIES all

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get update -y --no-install-recommends \
    && apt-get -y install --no-install-recommends bash git wget curl tzdata build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev cmake unzip zsh ca-certificates sudo apt-transport-https nano zip openssh-client apt-utils pkg-config gcc protobuf-compiler \
    && apt-get autoremove -y

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


WORKDIR /opt

RUN chown -R builder:builder /opt
RUN chmod -R 755 /opt
RUN mkdir -p /go && chmod -R 777 /go
RUN git config --global --add safe.directory /opt

USER builder


RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

ENV TERM xterm
# set the zsh theme
ENV ZSH_THEME agnoster

# run the installation script  
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true

COPY ./lib ./lib
COPY ./pkg ./pkg
COPY ./views ./views
COPY ./*.go .
COPY ./Cargo.* .
COPY ./go.mod .
COPY ./go.sum .
COPY ./Makefile .
USER root
RUN chown -R builder:builder /opt
RUN chmod -R 755 /opt
USER builder

RUN go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

# # Build the Go app
RUN --mount=type=cache,id=rustcache,target=/usr/local/cargo/registry,uid=1000,gid=1000 \
    --mount=type=cache,id=gomod,target=/go/pkg/mod,uid=1000,gid=1000 \
    --mount=type=cache,id=gobuild,target=/home/builder/.cache/go-build,uid=1000,gid=1000 \
    --mount=type=cache,id=gobuildtmp,target=/tmp/go-build,uid=1000,gid=1000 \
    make build

# # # Run the executable
ENTRYPOINT ["./100kb.golang"]