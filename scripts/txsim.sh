#!/bin/bash

# Variables
USER="root"
TMUX_SESSION_NAME="txsim_session_2"
COMMAND="sleep 180 && txsim --blob 1 --blob-amounts 1 --blob-sizes 500000-1000000 --key-path .celestia-app --grpc-endpoint localhost:9090 --feegrant"

# Function to start tmux session on a remote server
start_tmux_session() {
  tmux new-session -d -s ${TMUX_SESSION_NAME}
  tmux send-keys -t ${TMUX_SESSION_NAME} "${COMMAND}" C-m
}

# Optional: Run as a specific user
run_as_user() {
  su - ${USER} -c "$(declare -f start_tmux_session); start_tmux_session"
}


run_as_user &

wait
