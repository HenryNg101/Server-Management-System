COMPOSE_INFRA  = -f docker-compose.infra.yml
COMPOSE_APP    = -f docker-compose.yml
COMPOSE_DEV	   = -f docker-compose.dev.yml
COMPOSE_CLIENT = -f docker-compose.client.yml
ENV_FILE       = .env.docker

.PHONY: env_init dev.build dev.up dev.down prod.build prod.up prod.down prod.logs client.build client.up client.down

# Detect if the OS is Windows or Linux/WSL, then grab the IP address of the Linux WSL instance dynamically. 
# This is to bypass the weird networking issues from Windows Hyper-V, where they randomly reserve ports (Even if not use) on localhost
ifeq ($(OS),Windows_NT)
    HOST := $(shell wsl -e bash -c 'hostname -I | awk "{print $$1}"')
else ifeq ($(shell uname -s),Linux)
	# If running natively inside Linux bash
	HOST := $(shell hostname -I | awk '{print $$1}')
# Temporarily comment out for Mac, because the port reservation problem doesn't quite exist on MacOS's network virtualization for Docker yet
# else
#     # If running natively inside MacOS cmdline, use the hostname command to get the IP address
#     HOST := $(shell hostname)
endif

# Safe default fallback if the above evaluation returns empty
HOST ?= 127.0.0.1

# Clean up old HOST definitions on env files and append the fresh current IP into the file
env_init:
	$(eval export ENV_FILE=$(ENV_FILE))
	@echo "Updating $(ENV_FILE) with active IP: $(HOST)"
ifeq ($(OS),Windows_NT)
	@powershell -Command "if (Test-Path $(ENV_FILE)) { (Get-Content $(ENV_FILE)) -notmatch '^ *HOST=' | Set-Content $(ENV_FILE) }"
	@powershell -Command "Add-Content $(ENV_FILE) 'HOST=$(HOST)'"
else
	@sed -i '/^ *HOST=/d' $(ENV_FILE) 2>/dev/null || true
	@echo "HOST=$(HOST)" >> $(ENV_FILE)
endif

# ==========================================
# For actual deployment
# ==========================================
prod.build: env_init
#	Build cleanly without crashing resources, using sequential cache-warming by building one service first, to allow Docker caching even working at all
#	This is to avoid the issue of Docker caching not working properly when building multiple services at once, as Docker builds services in parallel and cache is not utilized effectively right in the first build, which can lead to unnecessary rebuilds of GBs of data, longer build times, and in worst cases, thrashings.
	docker compose --env-file $(ENV_FILE) $(COMPOSE_INFRA) $(COMPOSE_APP) build api
	docker compose --env-file $(ENV_FILE) $(COMPOSE_INFRA) $(COMPOSE_APP) build

prod.up: prod.build
	docker compose --env-file $(ENV_FILE) $(COMPOSE_INFRA) $(COMPOSE_APP) up -d --no-build
	@echo "========================================="
	@echo "Compose services are available! Access them through: http://$(HOST):<port>"
	@echo "========================================="

prod.down:
	docker compose $(COMPOSE_INFRA) $(COMPOSE_APP) down -v

prod.logs:
	docker compose $(COMPOSE_INFRA) $(COMPOSE_APP) logs -f $(SERVICE)

# ==========================================
# For development. The only difference from this to production is that, in this one, ports are exposed for debugging purposes. In production, ports for infrastructure are not exposed to the outside world for security reasons.
# ==========================================
dev.build: env_init
	docker compose --env-file $(ENV_FILE) $(COMPOSE_INFRA) $(COMPOSE_DEV) $(COMPOSE_APP) build api
	docker compose --env-file $(ENV_FILE) $(COMPOSE_INFRA) $(COMPOSE_DEV) $(COMPOSE_APP) build

dev.up: dev.build
	docker compose --env-file $(ENV_FILE) $(COMPOSE_INFRA) $(COMPOSE_DEV) $(COMPOSE_APP) up -d --no-build
	@echo "========================================="
	@echo "Compose services are available! Access them through: http://$(HOST):<port>"
	@echo "========================================="

dev.down:
	docker compose --env-file $(ENV_FILE) $(COMPOSE_INFRA) $(COMPOSE_DEV) $(COMPOSE_APP) down

# ==========================================
# For testing monitoring agent purposes
# ==========================================
client.build: env_init
	docker compose --env-file $(ENV_FILE) $(COMPOSE_CLIENT) build

client.up: client.build
	docker compose --env-file $(ENV_FILE) $(COMPOSE_CLIENT) up -d --no-build
	@echo "========================================="
	@echo "Docker compose services are available! Access them through: http://$(HOST):<port>"
	@echo "========================================="

client.down:
	docker compose $(COMPOSE_CLIENT) down
