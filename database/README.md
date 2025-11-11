# Database Package

This package handles PostgreSQL database connections and initialization for the shodan-backend application.

## Configuration

The database connection uses PostgreSQL with the following configuration:
- Connection string: `postgresql://shodan_user:RR3VRkIbn3DF34A0PDFI5IeEu0GqjSUR@dpg-d49322je5dus73ch5ej0-a/shodan`
- The connection can be overridden using the `DATABASE_URL` environment variable

## Usage

```go
import "shodan-backend/database"

db, err := database.InitDB()
if err != nil {
    log.Fatal(err)
}
defer database.CloseDB(db)
```

## Features

- Auto-migration for models
- Connection health checking
- Graceful connection closing
- Environment variable support for connection string override