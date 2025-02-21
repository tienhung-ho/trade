# Makefile
include .env
.PHONY: up down ps logs mysql-cli redis-cli clean createmigrate

# Docker compose commands
up:
	docker-compose up -d

down:
	docker-compose down

ps:
	docker-compose ps

logs:
	docker-compose logs -f

# Database management
createmigrate:
	migrate create -ext sql -dir migrations/mysql -seq init_schema

mysql-cli:
	docker exec -it trade-mysql mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) $(MYSQL_DATABASE)

redis-cli:
	docker exec -it trade-redis redis-cli

# Database operations
createdb:
	docker exec -it trade-mysql mysql -u root -p$(MYSQL_ROOT_PASSWORD) -e "CREATE DATABASE IF NOT EXISTS $(MYSQL_DATABASE);"

dropdb:
	docker exec -it trade-mysql mysql -u root -p$(MYSQL_ROOT_PASSWORD) -e "DROP DATABASE IF EXISTS $(MYSQL_DATABASE);"

# Migration commands
migrateup:
	migrate -path migrations/mysql -database "$(MYSQL_URL)" -verbose up

migratedown:
	migrate -path migrations/mysql -database "$(MYSQL_URL)" -verbose down
# Cleanup
clean:
	docker-compose down -v
