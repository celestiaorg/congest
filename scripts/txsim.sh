#!/bin/bash

run_as_user() {
  sleep 180 && txsim --blob 1 --blob-amounts 4 --blob-sizes 500000-1000000 --key-path .celestia-app --grpc-endpoint localhost:9090 --feegrant
}


run_as_user &

