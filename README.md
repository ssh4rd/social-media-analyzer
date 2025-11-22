# social-media-analyzer

A Go web application with a custom HTTP router, PostgreSQL database integration, and Docker containerization.

## Features

- **Custom HTTP Router**: Lightweight router with support for:
  - Path parameters (`:param`)
  - Catch-all routes (`*wildcard`)
  - Middleware support
  - Method-based routing (GET, POST, etc.)
- **PostgreSQL Database**: Fully containerized database with persistent storage
- **Template Rendering**: HTML templates for server-side rendering
- **Hot Reload**: Development environment with Air for automatic reloading
- **Docker Support**: Complete containerization with Docker Compose

## Project Structure

```
social-media-analyzer/
├── cmd/
│   └── app/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration loader
│   ├── db/
│   │   ├── db.go            # Database package
│   │   └── migrations/       # Database migrations
│   ├── dto/                  # Data Transfer Objects
│   ├── models/               # Domain models
│   ├── repo/                 # Repository layer
│   ├── service/              # Business logic
│   ├── transport/
│   │   └── http/
│   │       ├── controller/   # HTTP handlers
│   │       └── router/       # Custom router implementation
│   └── util/                 # Utility functions
├── web/
│   ├── static/               # Static assets (CSS, JS, images)
│   └── templates/            # HTML templates
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile                # Application Dockerfile
├── .env                      # Environment variables (local, not in git)
└── .env.example              # Environment variables template
```

## Prerequisites

- **Go 1.24.0** or higher
- **Docker** and **Docker Compose**
- **Air** (optional, for development hot reload)

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd social-media-analyzer
```

### 2. Configure Environment Variables

Copy the example environment file and adjust as needed:

```bash
cp .env.example .env
```

Default configuration:
```env
# Application
PORT=3000

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=social-media-analyzer
```

### 3. Run with Docker Compose

Build and start all services (application + PostgreSQL):

```bash
docker-compose up --build
```

The application will be available at `http://localhost:3000`

### 4. Stop the Services

```bash
docker-compose down
```

To remove volumes (database data):

```bash
docker-compose down -v
```

## Data Loading

### Loading Data from CSV

The application includes a data loader that populates the database with employee travel data from a CSV file.

#### Prerequisites

- Docker and Docker Compose running
- CSV file at `datasets/employee_travel_data.csv`

#### Running the Loader

Once the application is running with `docker-compose up`, load the data:

```bash
docker-compose exec app ./loader -file datasets/employee_travel_data.csv
```

The loader will:
- Parse the CSV file
- Create employee records
- Create business trip records
- Link employees to trips with expense information

To verify the data was loaded:

```bash
# Check table counts
docker-compose exec postgres psql -U postgres -d social-media-analyzer -c "\
SELECT 'employees' as table_name, COUNT(*) as row_count FROM employees UNION ALL \
SELECT 'business_trips', COUNT(*) FROM business_trips UNION ALL \
SELECT 'assignment_to_trips', COUNT(*) FROM assignment_to_trips;"
```

## Development

### Local Development with Air

For hot reload during development:

1. Install Air:
```bash
go install github.com/cosmtrek/air@latest
```

2. Run Air:
```bash
air
```

Configuration is in `.air.toml`.

### Local Development without Docker

1. Start PostgreSQL locally or use a remote instance

2. Update `.env` with your database connection:
```env
DB_HOST=localhost
DB_PORT=5432
```

3. Run the application:
```bash
go run cmd/app/main.go
```

## Configuration

The application uses environment variables for configuration, loaded through `internal/config/config.go`:

- `PORT` - HTTP server port (default: 3000)
- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: social-media-analyzer)

## API Endpoints

- `GET /` - Main page
- `GET /employee/:id` - Get employee by ID
- `/static/*` - Static file server

## Database

### Migrations

Database migrations are located in `internal/db/migrations/` and are automatically executed on container startup.

### PostgreSQL Container

- **Port**: 5432 (mapped to host)
- **Data**: Persisted in Docker volume `postgres_data`
- **Health Check**: Ensures database is ready before app starts

## Docker Services

### app
- Built from local Dockerfile
- Depends on PostgreSQL
- Port 3000 exposed
- Volume mount for web assets (hot reload in development)

### postgres
- PostgreSQL 16 Alpine
- Persistent data volume
- Health check enabled
- Auto-loads migration scripts

## Building

### Build the Application

```bash
go build -o main ./cmd/app
```

### Build Docker Image

```bash
docker build -t social-media-analyzer .
```

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
