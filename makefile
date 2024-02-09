run:
	sh ./scripts/run.sh

down:
	docker-compose down
.PHONY: create-network run