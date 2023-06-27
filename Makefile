mup:
	migrate -path db/migration -database "postgresql://postgres:12345678@localhost:5432/postgres?sslmode=disable" -verbose up
mdown:
	migrate -path db/migration -database "postgresql://postgres:12345678@vienan-prod.cgysah7oi4oj.ap-southeast-1.rds.amazonaws.com:5432/postgres" -verbose down
mforce:
	migrate -path db/migration -database "postgresql://postgres:12345678@vienan-prod.cgysah7oi4oj.ap-southeast-1.rds.amazonaws.com:5432/postgres" -verbose force 1
	 
sqlc:
	docker run --rm -v ".://src" -w //src kjconroy/sqlc generate 

test:
	go test -v -cover ./...
	
.PHONY: mup mdown mforce sqlc test