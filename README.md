# Villainrsty E-commerce Server

RESTful e-commerce API server built with Go (Chi router), JWT authentication, Casbin for role-based access control, and PostgreSQL.

## Tech Stack

- **Language:** Go 1.25+
- **Router:** Chi v5
- **Database:** PostgreSQL (pgx driver)
- **Authentication:** JWT (golang-jwt)
- **Authorization:** Casbin v3 (RBAC)
- **Validation:** go-playground/validator

## Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL
- Air (for hot reload, optional)

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/villainrsty-ecommerce-server.git
cd villainrsty-ecommerce-server
```

2. Install dependencies

```bash
go mod download
```

3. Setup environment variables

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run database migrations

```bash
# Add your migration command here
```

5. Run the server

```bash
go run cmd/api/main.go
```

Or with Air (hot reload):

```bash
air
```

## API Endpoints

Server runs on `http://localhost:3000`

## License

MIT
