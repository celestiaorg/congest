#!/bin/bash

tmux kill-session -t app
rm -rf /root/payload
rm -rf /root/payload.tar.gz
rm -rf /root/.celestia-app
rm -rf /root/celestia-app
rm -rf /root/logs
rm -rf /root/go/bin/*
rm -rf /root/quic-bench
