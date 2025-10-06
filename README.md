# Verve API

## Database Migrations

This project uses SQL migrations to manage database schema changes. Migrations are located in the `/migrations` directory.

### Prerequisites

1. Install golang-migrate:
```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### Running Migrations

1. Set up your database connection string as an environment variable:
```bash
export VERVE_DB_URL="postgres://username:password@localhost:5432/verve?sslmode=disable"
```

2. Run all pending migrations:
```bash
migrate -database "${VERVE_DB_URL}" -path migrations up
```

3. Roll back the last migration:
```bash
migrate -database "${VERVE_DB_URL}" -path migrations down 1
```

4. Roll back all migrations:
```bash
migrate -database "${VERVE_DB_URL}" -path migrations down
```

### Creating New Migrations

To create a new migration:

```bash
migrate create -ext sql -dir migrations -seq migration_name
```

This will create two files:
- `XXXXXX_migration_name.up.sql`: Contains the changes to apply
- `XXXXXX_migration_name.down.sql`: Contains the commands to roll back the changes

## API Documentation

This project uses Swagger for API documentation. The documentation is automatically generated from code annotations.

### Generating Swagger Documentation

1. Install swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate the documentation:
```bash
swag init -g cmd/server/main.go --output docs
```

### Viewing API Documentation

Once the server is running, you can access the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

### Swagger Annotations

Example of documenting an endpoint:

```go
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User details"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func CreateUserHandler() {}
```

Common annotations:
- `@Summary`: A short summary of what the endpoint does
- `@Description`: A detailed description of the endpoint
- `@Tags`: Grouping for the endpoint in Swagger UI
- `@Accept`: Accepted request content types
- `@Produce`: Response content types
- `@Param`: Parameters the endpoint accepts
- `@Success`: Successful response
- `@Failure`: Error responses
- `@Router`: The endpoint's path and method

Remember to regenerate the documentation whenever you make changes to the annotations:
```bash
swag init -g cmd/server/main.go --output docs
```