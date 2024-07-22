.PHONY: deploy destroy test

deploy:
	set -e; \
	pulumi up --yes && \
	bash send_payload.sh && \
	bash boot_network.sh

destroy:
	pulumi down --yes
	bash clean.sh

test:
	pulumi config set test $(word 2,$(MAKECMDGOALS))
	$(MAKE) deploy

%:
	@:
