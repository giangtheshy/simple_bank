mup:
	migrate -path db/migration -database "postgresql://postgres:12345678@vienan-prod.cgysah7oi4oj.ap-southeast-1.rds.amazonaws.com:5432/postgres?sslmode=disable" -verbose up
mdown:
	migrate -path db/migration -database "postgresql://postgres:12345678@vienan-prod.cgysah7oi4oj.ap-southeast-1.rds.amazonaws.com:5432/postgres" -verbose down
mforce:
	migrate -path db/migration -database "postgresql://postgres:12345678@vienan-prod.cgysah7oi4oj.ap-southeast-1.rds.amazonaws.com:5432/postgres" -verbose force 1
migrateup-github:
	migrate -path db/migration -database "postgresql://postgres:12345678@localhost:5432/postgres?sslmode=disable" -verbose up
	 
sqlc:
	docker run --rm -v ".://src" -w //src kjconroy/sqlc generate 

test:
	go test -v -cover ./...
	
server:
	gin -p 8081 -i run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/giangtheshy/simple_bank/db/sqlc Store


.PHONY: mup mdown mforce sqlc test server mock