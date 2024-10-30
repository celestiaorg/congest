#!/bin/bash
export DEBIAN_FRONTEND=noninteractive

# Check if a session name is passed as the first argument
if [ -z "$1" ]; then
  echo "Usage: $0 <session-name>"
  exit 1
fi

# Kill the specified tmux session
tmux kill-session -t "$1"

# Kill the current tmux session
tmux kill-session