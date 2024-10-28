# Start PostgreSQL using docker-compose.db.yml
.PHONY: db-up
db-up:
	docker-compose -f docker-compose.db.yml up -d

# Stop PostgreSQL
.PHONY: db-down
db-down:
	docker-compose -f docker-compose.db.yml down -v
