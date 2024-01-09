#!/usr/bin/env bash

set -e
# Install Go tools that are isImportant && !replacedByGopls based on
# https://github.com/golang/vscode-go/blob/0ff533d408e4eb8ea54ce84d6efa8b2524d62873/src/goToolsInformation.ts
# Exception `dlv-dap` is a copy of github.com/go-delve/delve/cmd/dlv built from the master.
TARGET_GOROOT=${2:-"/usr/local/go"}
TARGET_GOPATH=${3:-"/go"}

GO_TOOLS="\
    golang.org/x/tools/gopls@latest \
    honnef.co/go/tools/cmd/staticcheck@latest \
    golang.org/x/lint/golint@latest \
    github.com/mgechev/revive@latest \
    github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest \
    github.com/ramya-rao-a/go-outline@latest \
    github.com/go-delve/delve/cmd/dlv@latest \
    github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

echo "Installing common Go tools..."
export PATH=${TARGET_GOROOT}/bin:${PATH}
mkdir -p /tmp/gotools /usr/local/etc/vscode-dev-containers ${TARGET_GOPATH}/bin
chmod -R 0777 ${TARGET_GOPATH}
cd /tmp/gotools
export GOPATH=/tmp/gotools
export GOCACHE=/tmp/gotools/cache

# Use go get for versions of go under 1.16
go_install_command=install
if [[ "1.16" > "$(go version | grep -oP 'go\K[0-9]+\.[0-9]+(\.[0-9]+)?')" ]]; then
    export GO111MODULE=on
    go_install_command=get
    echo "Go version < 1.16, using go get."
fi 

(echo "${GO_TOOLS}" | xargs -n 1 go ${go_install_command} -v )2>&1 | tee -a /usr/local/etc/vscode-dev-containers/go.log

# Move Go tools into path and clean up
mv /tmp/gotools/bin/* ${TARGET_GOPATH}/bin/

# install dlv-dap (dlv@master)
go ${go_install_command} -v github.com/go-delve/delve/cmd/dlv@master 2>&1 | tee -a /usr/local/etc/vscode-dev-containers/go.log
mv /tmp/gotools/bin/dlv ${TARGET_GOPATH}/bin/dlv-dap

rm -rf /tmp/gotools