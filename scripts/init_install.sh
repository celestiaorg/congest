#!/bin/bash

CELES_HOME=".celestia-app"
MONIKER="validator"
ARCHIVE_NAME="payload.tar.gz"

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

# Make the changes persistent across reboots
echo "Making changes persistent..."
echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf

# Verify the current TCP congestion control algorithm
current_algo=$(sysctl net.ipv4.tcp_congestion_control | awk '{print $3}')
if [ "$current_algo" == "bbr" ]; then
  echo "Successfully switched to BBR congestion control algorithm."
else
  echo "Failed to switch to BBR. Current algorithm is $current_algo."
  exit 1
fi

echo "Script completed successfully."

tar -xzf /root/$ARCHIVE_NAME -C /root/

source ./vars.sh

sudo apt install git build-essential ufw curl jq snapd --yes

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
source payload/txsim.sh

celestia-appd start