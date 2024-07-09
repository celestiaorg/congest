#!/bin/bash

DIRECTORY_TO_TRANSFER="./payload"
ARCHIVE_NAME="payload.tar.gz"
SSH_KEY="~/.ssh/id_rsa"

# Fetch the IP addresses from Pulumi stack outputs
pulumi stack output -j > ./payload/ips.json
STACK_OUTPUT=$(pulumi stack output -j)
echo $STACK_OUTPUT
DROPLET_IPS=$(echo "$STACK_OUTPUT" | jq -r '.[]')

# Compress the directory
echo "Compressing the directory $DIRECTORY_TO_TRANSFER..."
tar -czf "$ARCHIVE_NAME" -C "$(dirname "$DIRECTORY_TO_TRANSFER")" "$(basename "$DIRECTORY_TO_TRANSFER")"

for IP in $DROPLET_IPS; do
  echo "Transferring files to $IP -----------------------------------------------------"
  scp -i "$SSH_KEY" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "$ARCHIVE_NAME" "root@$IP:/root/"

  # Uncompress the directory on the remote node
  echo "Uncompressing the directory on $IP..."
  ssh -i "$SSH_KEY" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "root@$IP" "tar -xzf /root/$ARCHIVE_NAME -C /root/"
done

# Cleanup local archive
rm "$ARCHIVE_NAME"