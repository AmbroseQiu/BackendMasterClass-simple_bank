postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:PAxG9FR0e58AJSFJpmgr@simplebank.cwsfirei38lc.ap-northeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:PAxG9FR0e58AJSFJpmgr@simplebank.cwsfirei38lc.ap-northeast-1.rds.amazonaws.com:5432/simple_bank" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:PAxG9FR0e58AJSFJpmgr@simplebank.cwsfirei38lc.ap-northeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:PAxG9FR0e58AJSFJpmgr@simplebank.cwsfirei38lc.ap-northeast-1.rds.amazonaws.com:5432/simple_bank" -verbose down 1

server:
	go run main.go
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
mock:
	mockgen -package mockdb -destination=./db/mock/store.go github.com/backendmaster/simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migratieup migratiedown migratieup1 migratiedown1 sqlc test server