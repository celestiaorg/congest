#!/bin/bash

# Set default SSH key location
DEFAULT_SSH_KEY="~/.ssh/do"
# Allow user to override the SSH key location
SSH_KEY=${SSH_KEY:-$DEFAULT_SSH_KEY}

# Fetch the IP addresses from Pulumi stack outputs
pulumi stack output -j > ./payload/ips.json
STACK_OUTPUT=$(pulumi stack output -j)
echo $STACK_OUTPUT
DROPLET_IPS=$(echo "$STACK_OUTPUT" | jq -r '.[]')

# Function to transfer and uncompress files on the remote server
remove_data() {
  local IP=$1
  echo "cleaning up $IP -----------------------"
  ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "root@$IP" "bash payload/cleanup_server.sh"
}

# Loop through the IPs and run the transfer and uncompress in parallel
for IP in $DROPLET_IPS; do
  remove_data "$IP" &
done

wait