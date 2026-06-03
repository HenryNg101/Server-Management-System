up:
	docker compose --env-file .env.docker up -d

down:
	docker compose down

logs:
	docker compose logs -f