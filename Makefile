GOOSE:=goose -env .env
.PHONY: sql migrate-up migrate-down migrate-status

sql:
	sqlc generate -f ./pkg/database/sqlc.yaml

migrate-up:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-status:
	$(GOOSE) status

