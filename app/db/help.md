# Migration Commands

**CREATE MIGRATION**
migrate create -ext sql -dir app/db/migrations create_package_table
**RUN MIGRATION**
migrate -path app/db/migrations -database "mysql://root@tcp(127.0.0.1:3306)/lift_gym_db" up
