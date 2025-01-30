.PHONY: up down deploy collect-traces

collect-traces:
	. ./collect_traces.sh
	sleep 120

# up wraps the pulumi up command to automate calling the scripts that initialize the network
up:
	set -e; \
	pulumi up --yes --continue-on-error # && \
	# . ./send_payload.sh && \
	# . ./boot_network.sh

# down wraps the pulumi down command to also clean up local objects
down:
	pulumi down --yes
	. ./scripts/clean.sh
	pulumi refresh --yes

# force-down uses the digitalocean's and linode's cli to manually close all
# instances with the "temp" tag. This can be needed in certain scenarios where
# pulumi doesn't register the nodes that are started.
force-down:
	doctl compute droplet list --tag-name temp --format ID --no-header | xargs -n1 doctl compute droplet delete -f
	linode-cli linodes list --tags temp --json 2>/dev/null | jq -r '.[].id' | xargs -n1 linode-cli linodes delete 2>/dev/null

# deploy wraps the up command by exporting the experiment and chain id as env vars before hand
deploy:
	$(MAKE) EXPERIMENT_NAME=$(word 2,$(MAKECMDGOALS)) EXPERIMENT_CHAIN_ID=$(word 3,$(MAKECMDGOALS)) up

# overwrite erases the current network and deploys a new one.

JSON_FILE=./payload/ips.json

ssh:
	@validator_name=$(word 2,$(MAKECMDGOALS)); \
	ip=$(shell jq -r '.["'$$validator_name'"]' $(JSON_FILE)); \
	if [ -z "$$ip" ]; then \
		echo "Validator $$validator_name not found!"; \
	else \
		echo "SSH into $$validator_name ($$ip)"; \
		ssh root@$$ip; \
	fi

%:
	@:
