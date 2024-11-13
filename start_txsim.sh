#!/bin/bash

# Expand the SSH key path
# Set default SSH key location
DEFAULT_SSH_KEY="~/.ssh/do"

# Allow user to override the SSH key location
SSH_KEY=${SSH_KEY:-$DEFAULT_SSH_KEY}

TIMEOUT=60

# Fetch the IP addresses from Pulumi stack outputs
STACK_OUTPUT=$(pulumi stack output -j)
DROPLET_IPS=$(echo "$STACK_OUTPUT" | jq -r '.[]')

# Variables
USER="root"
TMUX_SESSION_NAME="txsim"
COMMAND="./go/bin/txsim .celestia-app/keyring-test --blob 3 --blob-amounts 1 --blob-sizes 100000-200001 --key-path .celestia-app --grpc-endpoint localhost:9090 --feegrant"
# COMMAND="tmux send-keys -t app 'export SEEN_LIMIT=83' C-m"

# Function to start tmux session on a remote server
start_tmux_session() {
  local IP=$1
  {
    echo "Starting tmux session on $IP -----------------------------------------------------"
    ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i ${SSH_KEY} ${USER}@${IP} << EOF
tmux new-session -d -s ${TMUX_SESSION_NAME}
tmux send-keys -t ${TMUX_SESSION_NAME} "${COMMAND}" C-m
EOF
    echo "Tmux session started on $IP"
  } &

  PID=$!
  (sleep $TIMEOUT && kill -HUP $PID) 2>/dev/null &

  if wait $PID 2>/dev/null; then
    echo "$IP: Tmux session started within timeout"
  else
    echo "$IP: Operation timed out"
  fi
}

MAX_LOOPS=100
COUNTER=0

for IP in $DROPLET_IPS; do
  start_tmux_session "$IP" &

  # Increment the counter
  COUNTER=$((COUNTER + 1))

  # Check if the counter has reached the max number of loops
  if [ "$COUNTER" -ge "$MAX_LOOPS" ]; then
    break
  fi
done

# Wait for all background processes to finish
wait
