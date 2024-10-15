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

# Build the application
.PHONE: build
build:
	go build -o main ./cmd/api/main.go

# Run the application
.PHONE: run
run:
	./main