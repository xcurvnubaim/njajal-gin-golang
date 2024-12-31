# Migrate database
.PHONY: migrate
migrate:
	go run ./cmd/migration/main.go migrate up

# Migrate Fresh database
.PHONY: migrate-fresh
migrate-fresh:
	go run ./cmd/migration/main.go migrate:fresh

# Rollback database
.PHONY: rollback
rollback:
	go run ./cmd/migration/main.go migrate down

# Create migration
.PHONY: create-migration
create-migration:
	go run ./cmd/migration/main.go create migration ${name}