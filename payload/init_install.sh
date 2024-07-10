#!/bin/bash

CELES_HOME=".celestia-app"
MONIKER="validator"

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

celestia-appd start