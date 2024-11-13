#!/bin/bash

# Set default SSH key location
DEFAULT_SSH_KEY="~/.ssh/do"

# Allow user to override the SSH key location
SSH_KEY=${SSH_KEY:-$DEFAULT_SSH_KEY}


# Rest of your script goes here
TIMEOUT=120

# Fetch the IP addresses from Pulumi stack outputs
STACK_OUTPUT=$(pulumi stack output -j)
DROPLET_IPS=$(echo "$STACK_OUTPUT" | jq -r '.[]')

# Function to boot a node
boot_node() {
  local IP=$1
  {
    echo "Booting $IP -----------------------------------------------------"
    ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "root@$IP" \
    "tmux new -d -s app 'source /root/payload/init_install.sh'"
    echo "Boot complete for $IP"
  } &

  PID=$!
  (sleep $TIMEOUT && kill -HUP $PID) 2>/dev/null &

  if wait $PID 2>/dev/null; then
    echo "$IP: Boot finished within timeout"
  else
    echo "$IP: Boot operation timed out"
  fi
}

# Loop through the IPs and run the boot command in parallel
for IP in $DROPLET_IPS; do
  boot_node "$IP"
done
