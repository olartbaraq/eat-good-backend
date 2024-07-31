c_m:
	#create a new migration
	@migrate create -ext sql -dir db/migrations -seq $(name)

p_up:
	#create all services listed in docker compose file with docker
	@docker compose up -d

p_down:
	#delete all services listed in docker compose file
	@docker compose down

db_up:
	#create a database from the db server
	@docker exec -it eat_good_test createdb --username=root --owner=root eat_good_db_test
	@docker exec -it eat_good_live createdb --username=root --owner=root eat_good_db

db_down:
	#delete a database from the db server
	@docker exec -it eat_good_test dropdb --username=root eat_good_db_test
	@docker exec -it eat_good_live dropdb --username=root eat_good_db

m_up:
	#run a migration to the database
	@migrate -path db/migrations -database "postgres://root:testing@localhost:5432/eat_good_db_test?sslmode=disable" up
	@migrate -path db/migrations -database "postgres://root:testing@localhost:5433/eat_good_db?sslmode=disable" up

m_down:
	#revert the migration from the database
	@migrate -path db/migrations -database "postgres://root:testing@localhost:5432/eat_good_db_test?sslmode=disable" down
	@migrate -path db/migrations -database "postgres://root:testing@localhost:5433/eat_good_db?sslmode=disable" down

sqlc:
	#generate the sql queries to golang
	sqlc generate

test:
	#run all tests in test directory
	@go test -v -cover ./...