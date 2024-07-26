.PHONY: up down test

up:
	set -e; \
	pulumi up --yes && \
	. ./send_payload.sh && \
	. ./boot_network.sh && \
	. ./start_txsim.sh

down:
	pulumi down --yes
	. ./scripts/clean.sh

deploy:
	$(MAKE) EXPERIMENT_NAME=$(word 2,$(MAKECMDGOALS)) EXPERIMENT_CHAIN_ID=$(word 3,$(MAKECMDGOALS)) up


%:
	@:
