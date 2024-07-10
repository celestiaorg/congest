#!/bin/bash

# Expand the SSH key path
SSH_KEY="$HOME/.ssh/id_rsa"
TIMEOUT=60

# Fetch the IP addresses from Pulumi stack outputs
STACK_OUTPUT=$(pulumi stack output -j)
DROPLET_IPS=$(echo "$STACK_OUTPUT" | jq -r '.[]')

# Variables
USER="root"
TMUX_SESSION_NAME="txsim_session"
COMMAND="txsim --blob 1 --blob-amounts 1 --blob-sizes 500000-1000000 --key-path .celestia-app --grpc-endpoint localhost:9090 --feegrant"

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

# Loop through the IPs and run the start_tmux_session function in parallel
for IP in $DROPLET_IPS; do
  start_tmux_session "$IP" &
done

# Wait for all background processes to finish
wait
