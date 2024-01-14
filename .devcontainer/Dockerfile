FROM debian:latest

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10

COPY --from=golang:1.21-alpine /usr/local/go/ /usr/local/go/
COPY --from=bitnami/pytorch:2.1.2 /opt/bitnami/python/lib/python3.11/site-packages/torch/ /opt/libtorch/

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV TORCH_HOME=/opt/libtorch
ENV LIBTORCH=/opt/libtorch
ENV LD_LIBRARY_PATH=${LIBTORCH}/lib:$LD_LIBRARY_PATH
ENV DEBIAN_FRONTEND=noninteractive
ENV LIBTORCH_URL=https://download.pytorch.org/libtorch/cpu/libtorch-cxx11-abi-shared-with-deps-2.1.2%2Bcpu.zip
# ENV LIBTORCH_URL=https://download.pytorch.org/libtorch/cpu/libtorch-macos-2.1.2.zip
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
ENV PATH="/home/vscode/.cargo/bin:/usr/local/go/bin:/go/bin:${PATH}"
ENV LIBTORCH_BYPASS_VERSION_CHECK="true"

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
  && apt-get update -y --no-install-recommends \
  && apt-get -y install --no-install-recommends bash git wget curl tzdata build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev cmake unzip zsh ca-certificates sudo apt-transport-https nano zip openssh-client apt-utils pkg-config clang python3 python3-pip ninja-build gcc python3-dev protobuf-compiler \
  && apt-get autoremove -y

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
  echo "${USERNAME} ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers && \
  chmod 0440 /etc/sudoers && \
  chmod g+w /etc/passwd 


# RUN cd /tmp/ && \
#     wget https://www.python.org/ftp/python/3.12.1/Python-3.12.1.tgz  && \
#     tar -xzvf Python-3.12.1.tgz  && \
#     cd Python-3.12.1/ && \
#     ./configure --enable-optimizations && make -j `nproc` && make altinstall && \
#     ln -s /usr/local/bin/python3.12 /usr/local/bin/python


WORKDIR /opt

# x86
# RUN curl -fsSL --insecure -o libtorch.zip  $LIBTORCH_URL \
#   && unzip -q libtorch.zip \
#   && rm libtorch.zip

# RUN git clone --recursive https://github.com/pytorch/pytorch -b v2.0.1 --depth 1
# RUN cd /opt/pytorch && \
#   pip3 install -r requirements.txt --break-system-packages
# RUN cd /opt/pytorch && \
#   python3 tools/build_libtorch.py

COPY setup.sh /tmp/setup.sh
RUN bash /tmp/setup.sh \
  && apt-get clean -y

RUN chsh -s /usr/bin/zsh vscode

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

ENV TERM xterm
# set the zsh theme
ENV ZSH_THEME agnoster

# run the installation script  
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true

# RUN chsh -s $(which zsh)
RUN cargo install --force cbindgen
