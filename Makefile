migrate:
	dbmate dump

.PHONY: db-up
db-up:
	sudo docker compose up -d postgres_db db_migrations

delete-db:
	sudo docker compose down -v

.PHONY: dev
dev: db-up
	templ generate
	sqlc generate
	DB_CONTAINER_URL="postgres://adwaith:req3110@localhost:5432/reqdb?sslmode=disable" PORT=8080 air

.PHONY: vet
vet:
	go vet ./...

run: vet
	templ generate
	sqlc generate 
	dbmate dump
	sudo docker compose up --build
compile:
	templ generate
	sqlc generate 

