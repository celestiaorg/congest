#!/bin/bash

# Expand the SSH key path
SSH_KEY="$HOME/.ssh/id_rsa"

# Fetch the IP addresses from Pulumi stack outputs
STACK_OUTPUT=$(pulumi stack output -j)
DROPLET_IPS=$(echo "$STACK_OUTPUT" | jq -r '.[]')

# Variables
USER="root"
TMUX_SESSION_NAME="txsim_session"
COMMAND="txsim --blob 1 --blob-amounts 1 --blob-sizes 500000-1000000 --key-path .celestia-app --grpc-endpoint localhost:9090 --feegrant"

for IP in $DROPLET_IPS; do
    # SSH into the remote server and start a tmux session
    ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i ${SSH_KEY} ${USER}@${IP} << EOF
tmux new-session -d -s ${TMUX_SESSION_NAME}
tmux send-keys -t ${TMUX_SESSION_NAME} "${COMMAND}" C-m
EOF
done
