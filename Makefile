INFRA  = -f docker-compose.infra.yml
APP    = -f docker-compose.yml
DEV	   = -f docker-compose.dev.yml
CLIENT = -f docker-compose.client.yml

.PHONY: dev.build dev.up dev.down prod.build prod.up prod.down prod.logs client.build client.up client.down

# Reusable macro to build cleanly without crashing resources
# It works by building a single service first, then build all others later
# Why needs to do that ? Because each service spans across the entire codebase, which if is built simultaneously by Docker compose build, it would consume every resources of the Linux VM/host for Docker
# Usage: $(call safe_build, <compose_files>)
define safe_build
	docker compose --env-file .env.docker $(1) build api
	docker compose --env-file .env.docker $(1) build
endef

# ==========================================
# DEVELOPMENT PURPOSES
# ==========================================
dev.build:
	$(call safe_build, $(INFRA) $(DEV))

dev.up: dev.build
	docker compose --env-file .env.docker $(INFRA) $(DEV) up -d --no-build

dev.down:
	docker compose $(INFRA) $(DEV) down

# ==========================================
# PRODUCTION / ACTUAL DEPLOYMENT
# ==========================================
prod.build:
	$(call safe_build, $(APP) $(INFRA))

prod.up: prod.build
	docker compose --env-file .env.docker $(APP) $(INFRA) up -d --no-build

prod.down:
	docker compose $(APP) $(INFRA) down -v

prod.logs:
	docker compose $(APP) $(INFRA) logs -f $(SERVICE)

# ==========================================
# CLIENT TESTING PURPOSES
# ==========================================
client.build:
	docker compose --env-file .env.docker $(CLIENT) build

client.up: client.build
	docker compose --env-file .env.docker $(CLIENT) up -d --no-build

client.down:
	docker compose $(CLIENT) down