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

mkdir -p traces

stop_network() {
  local IP=$1
  echo "stopping network $IP -----------------------"
  ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@"$IP" "killall celestia-appd &"
}

# Loop through the IPs and run the download logs function
for IP in $DROPLET_IPS; do
  stop_network "$IP"
done
