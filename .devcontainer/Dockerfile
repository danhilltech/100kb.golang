FROM debian:latest

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10

COPY --from=golang:1.23-alpine /usr/local/go/ /usr/local/go/
COPY --from=bitnami/pytorch:2.2.2 /opt/bitnami/python/lib/python3.11/site-packages/torch/ /opt/libtorch/
# COPY --from=bitnami/pytorch:2.2.2 /opt/bitnami/python/lib/python3.11/site-packages/torch.libs/ /opt/libtorch/lib/
# above is needed on arm64
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV TORCH_HOME=/opt/libtorch
ENV LIBTORCH=/opt/libtorch
ENV LD_LIBRARY_PATH=${LIBTORCH}/lib:$LD_LIBRARY_PATH:/workspaces/100kb.golang/lib
ENV DEBIAN_FRONTEND=noninteractive
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
ENV PATH="/home/vscode/.cargo/bin:/usr/local/go/bin:/go/bin:${PATH}"
ENV LIBTORCH_BYPASS_VERSION_CHECK="true"

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
  && apt-get update -y --no-install-recommends \
  && apt-get -y install --no-install-recommends bash git wget curl tzdata build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev cmake unzip zsh ca-certificates sudo apt-transport-https nano zip openssh-client apt-utils pkg-config clang python3 python3-pip ninja-build gcc python3-dev protobuf-compiler chromium \
  && apt-get autoremove -y

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
  echo "${USERNAME} ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers && \
  chmod 0440 /etc/sudoers && \
  chmod g+w /etc/passwd 

WORKDIR /opt

COPY setup.sh /tmp/setup.sh
RUN bash /tmp/setup.sh \
  && apt-get clean -y

RUN chsh -s /usr/bin/zsh vscode
RUN mkdir -p /dbs

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode


RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

ENV TERM xterm
# set the zsh theme
ENV ZSH_THEME agnoster
ENV RUSTBERT_CACHE /workspaces/100kb.golang/models
ENV LIBTORCH_CXX11_ABI 0

# run the installation script  
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true

