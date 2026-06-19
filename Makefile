INFRA = -f docker-compose.infra.yml
APP   = -f docker-compose.yml
DEV	  = -f docker-compose.dev.yml
CLIENT = -f docker-compose.client.yml

# For the development purposes
dev.up:
	docker compose --env-file .env.docker $(INFRA) $(DEV) up -d

dev.down:
	docker compose $(INFRA) $(DEV) down

# For Docker and actual deployment
prod.up:
	docker compose --env-file .env.docker $(APP) $(INFRA) up -d

prod.down:
	docker compose $(APP) $(INFRA) down -v

prog.logs:
	docker compose $(APP) $(INFRA) logs -f $(SERVICE)

# For client
client.up:
	docker compose $(CLIENT) up -d

client.down:
	docker compose $(CLIENT) down