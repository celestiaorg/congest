#!/bin/bash

ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
  wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz && export PATH=$PATH:/usr/local/go/bin && go version
elif [[ "$ARCH" == "aarch64" ]]; then
    wget https://go.dev/dl/go1.23.2.linux-arm64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.2.linux-arm64.tar.gz && export PATH=$PATH:/usr/local/go/bin && go version
else
    exit 1
fi

apt install -y build-essential jq git htop vim nethogs

export PATH=$PATH:~/go/bin:/usr/local/go/bin
echo 'export PATH=$PATH:~/go/bin:/usr/local/go/bin' >> ~/.bashrc

git clone https://github.com/rach-id/quic-bench
cd quic-bench
go mod tidy
go build
./quic-bench -listen 0.0.0.0:4242 -peersFile /root/payload/validators.json
