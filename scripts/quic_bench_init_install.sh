#!/bin/bash

ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
  wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz && export PATH=$PATH:/usr/local/go/bin && go version
elif [[ "$ARCH" == "aarch64" ]]; then
    wget https://go.dev/dl/go1.23.2.linux-arm64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.2.linux-arm64.tar.gz && export PATH=$PATH:/usr/local/go/bin && go version
else
    exit 1
fi

apt install build-essential jq git htop vim nethogs --yes -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold"

export PATH=$PATH:~/go/bin:/usr/local/go/bin
echo 'export PATH=$PATH:~/go/bin:/usr/local/go/bin' >> ~/.bashrc

# increase udp buffers
sudo sysctl -w net.core.rmem_max=16777216
sudo sysctl -w net.core.wmem_max=16777216
sudo sysctl -w net.core.rmem_default=8388608
sudo sysctl -w net.core.wmem_default=8388608
sudo sysctl -w net.ipv4.udp_mem="8388608 8388608 16777216"
sudo sysctl -w net.ipv4.udp_rmem_min=1638400
sudo sysctl -w net.ipv4.udp_wmem_min=1638400

git clone https://github.com/rach-id/quic-bench
cd quic-bench
go mod tidy
go build
./quic-bench -listen 0.0.0.0:4242 -peersFile /root/payload/validators.json
