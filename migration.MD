# create a new migration
# example: migrate create -ext sql -dir ./migration -seq test
migrate create -ext sql -dir <your_migration_folder_path> -seq <migration_name>

# run migration up
migrate -database "postgres://user:password@localhost:5432/database_name?sslmode=disable" -path ./migration up

# run migration down
migrate -database "postgres://user:password@localhost:5432/database_name?sslmode=disable" -path ./migration down


