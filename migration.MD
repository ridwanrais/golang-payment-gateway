## Create a new migration

To create a new migration, use the following command:

	$ migrate create -ext sql -dir ./migration -seq <migration_name>
    
Replace `<your_migration_folder_path>` with the path where you want to store your migration files, and `<migration_name>` with the desired name for the migration.


## Run migration up

To apply the migrations to your PostgreSQL database, use the following command:

    $ migrate -database "postgres://user:password@localhost:5432/database_name?sslmode=disable" -path ./migration up


Replace the connection string in the `-database` flag with your PostgreSQL connection details.


## Run migration down (CAUTION: MAKE SURE YOU KNOW WHAT YOU ARE DOING BEFORE PROCEEDING)

To rollback the migrations, use the following command:

    $ migrate -database "postgres://user:password@localhost:5432/database_name?sslmode=disable" -path ./migration down
    
Replace the connection string in the `-database` flag with your PostgreSQL connection details.




