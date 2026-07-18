GOOSE:=goose -env .env
.PHONY: sql migrate

sql:
	sqlc generate -f ./pkg/database/sqlc.yaml

migrate:
	$(GOOSE) up

