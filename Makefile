INFRA  = -f docker-compose.infra.yml
APP    = -f docker-compose.yml
DEV	   = -f docker-compose.dev.yml
CLIENT = -f docker-compose.client.yml

.PHONY: dev.build dev.up dev.down prod.build prod.up prod.down prod.logs client.build client.up client.down

# Detect if the OS is Windows or Linux/WSL
ifeq ($(OS),Windows_NT)
    # If running on Windows PowerShell/CMD, reach inside WSL to get the exact Linux IP
    HOST := $(shell wsl -e bash -c 'hostname -I | awk "{print $$1}"')
else
    # If running natively inside Linux bash
    HOST := $(shell hostname -I | awk '{print $$1}')
endif

# Safe default fallback if the above evaluation returns empty
HOST ?= 127.0.0.1

# Reusable macro to build cleanly without crashing resources.
# Note: Variables in a Makefile target recipe must be prefixed natively on the same line.
define safe_build
	HOST=$(HOST) docker compose --env-file .env.docker $(1) build api
	HOST=$(HOST) docker compose --env-file .env.docker $(1) build
endef

# ==========================================
# DEVELOPMENT PURPOSES
# ==========================================
dev.build:
#	$(call safe_build, $(INFRA) $(DEV))
	HOST=$(HOST) docker compose --env-file .env.docker $(INFRA) $(DEV) build

dev.up: dev.build
	HOST=$(HOST) docker compose --env-file .env.docker $(INFRA) $(DEV) up -d --no-build
	@echo "========================================="
	@echo "🚀 Dev services are available! Access them through: http://$(HOST):<port>"
	@echo "========================================="

dev.down:
	docker compose $(INFRA) $(DEV) down

# ==========================================
# PRODUCTION / ACTUAL DEPLOYMENT
# ==========================================
prod.build:
	$(call safe_build, $(APP) $(INFRA))

prod.up: prod.build
	HOST=$(HOST) docker compose --env-file .env.docker $(APP) $(INFRA) up -d --no-build
	@echo "========================================="
	@echo "🚀 Prod services are available! Access them through: http://$(HOST):<port>"
	@echo "========================================="

prod.down:
	docker compose $(APP) $(INFRA) down -v

prod.logs:
	docker compose $(APP) $(INFRA) logs -f $(SERVICE)

# ==========================================
# CLIENT TESTING PURPOSES
# ==========================================
client.build:
	HOST=$(HOST) docker compose --env-file .env.docker $(CLIENT) build

client.up: client.build
	HOST=$(HOST) docker compose --env-file .env.docker $(CLIENT) up -d --no-build
	@echo "========================================="
	@echo "🚀 Client services are available! Access them through: http://$(HOST):<port>"
	@echo "========================================="

client.down:
	docker compose $(CLIENT) down
