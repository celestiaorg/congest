#!/bin/bash

tmux new -d -s tacosim "txsim --blob 1 --blob-amounts 1 --blob-sizes 500000-1000000 --key-path .celestia-app --grpc-endpoint localhost:9090 --feegrant"
