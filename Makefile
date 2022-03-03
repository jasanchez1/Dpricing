DB_URL = postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}
DB_NAME = dpricing

create-db:
	psql $(DB_URL) -w -c "create database $(DB_NAME)"

drop-db:
	psql $(DB_URL) -w -c "drop database $(DB_NAME)"
migrations-up:
	migrate -database $(DB_URL)/$(DB_NAME)?sslmode=disable -path db/migrations up