include .env

migration:
	goose -dir internal/db/migrations create $(name) sql

migrate:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up

migrate-up-one:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up-by-one

migratedbs:
	make migrate

rollback:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} down

api:
	go run cmd/api/*.go