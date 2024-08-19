.PHONY: up down deploy collect-traces

collect-traces:
	. ./collect_traces.sh
	sleep 120

# up wraps the pulumi up command to automate calling the scripts that initialize the network
up:
	set -e; \
	pulumi up --yes --continue-on-error && \
	. ./send_payload.sh && \
	. ./boot_network.sh

# down wraps the pulumi down command to also clean up local objects
down:
	collect-traces
	pulumi down --yes
	. ./scripts/clean.sh
	pulumi refresh --yes

# deploy wraps the up command by exporting the experiment and chain id as env vars before hand
deploy:
	$(MAKE) EXPERIMENT_NAME=$(word 2,$(MAKECMDGOALS)) EXPERIMENT_CHAIN_ID=$(word 3,$(MAKECMDGOALS)) up


%:
	@:
