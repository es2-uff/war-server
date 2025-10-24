.PHONY: docker-dev docker-hot dev build docker-down docker-clean docker-fclean

# Load environment variables
-include .env
export

# (DOCKER STATIC BUILD)
docker-dev:
	docker-compose --profile dev up -d

# (docker hot reload)
docker-hot:
	docker-compose --profile hot up -d

docker-down:
	docker-compose --profile hot down
	docker-compose --profile dev down

docker-clean: docker-down
	@if [ -n "$$(docker images -q es2_server:latest)" ]; then \
		docker rmi es2_server:latest; \
	fi

docker-fclean: docker-clean
	@if [ -n "$$(docker images -q redis:7-alpine)" ]; then \
		docker rmi redis:7-alpine; \
	fi
	docker volume prune -f
	docker network prune -f


# (HOST MACHINE HOT RELOAD IF EXISTS)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Running without hot reload..."; \
		go run cmd/api/main.go; \
	fi
