migrate:
	dbmate dump

delete:
	sudo docker compose down -v

.PHONY: vet
vet:
	go vet ./...

run: vet
	templ generate
	sqlc generate 
	sudo docker compose up --build
compile:
	templ generate
	sqlc generate 

