run:
	templ generate
	sqlc generate 
	go run .
compile:
	templ generate
	sqlc generate 

