#!/bin/bash

SSH_KEY="~/.ssh/id_rsa"

# Fetch the IP addresses from Pulumi stack outputs
STACK_OUTPUT=$(pulumi stack output -j)
DROPLET_IPS=$(echo "$STACK_OUTPUT" | jq -r '.[]')

for IP in $DROPLET_IPS; do
  echo "booting $IP -----------------------------------------------------"
  ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "root@$IP" \
  "tmux new -d -s init_install 'source /root/payload/init_install.sh'"
done