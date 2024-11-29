#!/bin/bash

CELESTIA_APP_COMMIT="323e54a6b536e3ef50d1237d2e8a0fc331d841ed"
CELES_HOME=".celestia-app"
MONIKER="validator"
ARCHIVE_NAME="payload.tar.gz"
export CHAIN_ID=another-test
export DEBIAN_FRONTEND=noninteractive

apt update -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold"

export DEBIAN_FRONTEND=noninteractive
apt-get update -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold"
apt-get install git build-essential curl jq chrony snapd nethogs --yes -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold"

systemctl enable chrony
systemctl start chrony

#export GOOGLE_APPLICATION_CREDENTIALS="/root/payload/congest-remote-key-gbq.json"

# Ensure the script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please run with sudo or as root."
  exit 1
fi

# Load the BBR module
echo "Loading BBR module..."
modprobe tcp_bbr

# Verify if the BBR module is loaded
if lsmod | grep -q "tcp_bbr"; then
  echo "BBR module loaded successfully."
else
  echo "Failed to load BBR module."
  exit 1
fi

# Add BBR to the list of available congestion control algorithms
echo "Updating sysctl settings..."
sysctl -w net.core.default_qdisc=fq
sysctl -w net.ipv4.tcp_congestion_control=bbr

# Enable MPTCP
sysctl -w net.mptcp.enabled=1
# Set the path manager to ndiffports
sysctl -w net.mptcp.mptcp_path_manager=ndiffports

# Specify the number of subflows
SUBFLOWS=16
sysctl -w net.mptcp.mptcp_ndiffports=$SUBFLOWS

# Make the changes persistent across reboots
echo "Making changes persistent..."
echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf

#Verify the current TCP congestion control algorithm
current_algo=$(sysctl net.ipv4.tcp_congestion_control | awk '{print $3}')
if [ "$current_algo" == "bbr" ]; then
  echo "Successfully switched to BBR congestion control algorithm."
else
  echo "Failed to switch to BBR. Current algorithm is $current_algo."
  exit 1
fi

echo "Script completed successfully."

tar -xzf /root/$ARCHIVE_NAME -C /root/

#source ./vars.sh

# increase udp buffers
sudo sysctl -w net.core.rmem_max=16777216
sudo sysctl -w net.core.wmem_max=16777216
sudo sysctl -w net.core.rmem_default=8388608
sudo sysctl -w net.core.wmem_default=8388608
sudo sysctl -w net.ipv4.udp_mem="8388608 8388608 16777216"
sudo sysctl -w net.ipv4.udp_rmem_min=1638400
sudo sysctl -w net.ipv4.udp_wmem_min=1638400

sudo snap install go --channel=1.22/stable --classic

echo 'export GOPATH="$HOME/go"' >> ~/.profile
echo 'export GOBIN="$GOPATH/bin"' >> ~/.profile
echo 'export PATH="$GOBIN:$PATH"' >> ~/.profile
source ~/.profile

cd $HOME
git clone https://github.com/celestiaorg/celestia-app
cd celestia-app
git checkout $CELESTIA_APP_COMMIT
make install
make txsim-install

celestia-appd config chain-id $CHAIN_ID

celestia-appd init --chain-id=$CHAIN_ID --home $CELES_HOME $MONIKER

cd $HOME

# Get the hostname
hostname=$(hostname)

# Parse the first part of the hostname
parsed_hostname=$(echo $hostname | awk -F'-' '{print $1 "-" $2}')

mv payload/$parsed_hostname/node_key.json $HOME/$CELES_HOME/config/node_key.json

mv payload/$parsed_hostname/priv_validator_key.json $HOME/$CELES_HOME/config/priv_validator_key.json

mv payload/$parsed_hostname/priv_validator_state.json $HOME/$CELES_HOME/data/priv_validator_state.json

cp payload/genesis.json $HOME/$CELES_HOME/config/genesis.json

cp payload/addrbook.json $HOME/$CELES_HOME/config/addrbook.json

mv payload/$parsed_hostname/app.toml $HOME/$CELES_HOME/config/app.toml

mv payload/$parsed_hostname/config.toml $HOME/$CELES_HOME/config/config.toml

cp -r payload/$parsed_hostname/keyring-test $HOME/$CELES_HOME

# run txsim script which starts a sleep timer and txsim in a different tmux session
#source payload/txsim.sh


# Get the hostname of the machine
HOSTNAME=$(hostname)

# Base command
COMMAND="celestia-appd start"

# Define log file path
LOG_FILE="/root/logs"

# Check if the hostname matches the specific value
if [[ "$HOSTNAME" == "validator-2" ]]; then
    # If it matches, don't add the log level flag
    echo "Starting celestia-appd without log level flag on $HOSTNAME"
else
    # Otherwise, add the log level flag
    COMMAND+=" --log_level=\"error\""
fi

for i in {1..10}; do
    eval $COMMAND 2>&1 | tee -a "$LOG_FILE"
    sleep 10
done
