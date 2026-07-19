migrate:
	dbmate dump

delete:
	sudo docker compose down -v



run:
	templ generate
	sqlc generate 
	sudo docker compose up --build
compile:
	templ generate
	sqlc generate 

