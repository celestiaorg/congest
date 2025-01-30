#!/bin/bash

# Set default SSH key location
DEFAULT_SSH_KEY="~/.ssh/id_ed25519"
# Allow user to override the SSH key location
SSH_KEY=${SSH_KEY:-$DEFAULT_SSH_KEY}

# Hardcoded path to the JSON file containing node info
NODES_JSON="/home/evan/go/src/github.com/celestiaorg/congest/payload/validators.json"

# Ensure required arguments are provided
if [ $# -ne 2 ]; then
  echo "Usage: source script.sh <trace_file_name.jsonl> <path/to/destination>"
fi

REMOTE_FILE=$1
DEST_PATH=$2

# Validate the JSON file exists
if [ ! -f "$NODES_JSON" ]; then
  echo "Node information JSON file not found at $NODES_JSON."
  exit 1
fi

# Create the destination directory if it doesn't exist
mkdir -p "$DEST_PATH"

# Function to download a file from the remote server
download_file() {
  local NODE_NAME=$1
  local IP=$2
  echo "Downloading $REMOTE_FILE from $NODE_NAME ($IP) -----------------------"
  scp -i "$SSH_KEY" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "root@$IP:/root/.celestia-app/data/traces/$REMOTE_FILE" "$DEST_PATH/$NODE_NAME-$REMOTE_FILE"
  
  if [ $? -eq 0 ]; then
    echo "Download successful from $NODE_NAME ($IP)."
  else
    echo "Download failed from $NODE_NAME ($IP)."
  fi
}

# Parse the JSON and iterate over nodes
jq -r 'to_entries[] | "\(.key) \(.value.ip)"' "$NODES_JSON" | while read -r NODE_NAME IP; do
  download_file "$NODE_NAME" "$IP" &
done

# Wait for all background processes to finish
wait
