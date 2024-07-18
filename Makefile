.PHONY: up down

up:
	pulumi up --yes && source send_payload.sh && source boot_network.sh

down:
	pulumi down --yes
