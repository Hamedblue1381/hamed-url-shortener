DB_URL=postgresql://root:secret@localhost:5432/hamed_bank?sslmode=disable

postgres:
	sudo docker run --name hamed-postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	sudo docker exec -it hamed-postgres createdb --username=root --owner=root hamed_url_shortener

dropdb:
	sudo docker exec -it hamed-postgres dropdb hamed_url_shortener

server:
	go run main.go

.PHONY: postgres createdb dropdb server
