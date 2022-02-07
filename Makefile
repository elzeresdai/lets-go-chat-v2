migrateup:
	migrate -path internal/db/migration -database "postgresql://root:root@localhost:5432/lets_go_chat?sslmode=disable" -verbose up
migratedown:
	migrate -path internal/db/migration -database "postgresql://root:root@localhost:5432/lets_go_chat?sslmode=disable" -verbose down