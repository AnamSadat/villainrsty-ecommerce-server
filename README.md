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

## Casbin RBAC Quick Start

Project ini sudah include RBAC Casbin untuk route yang merepresentasikan alur real project:

- Products: `GET /products`, `POST /products`, `PATCH /products/{id}`, `DELETE /products/{id}`
- Orders: `GET /orders`, `POST /orders`
- Admin Users: `GET /admin/users`, `POST /admin/users`, `DELETE /admin/users/{id}`

### 1 Konfigurasi env

Tambahkan env berikut:

```env
CASBIN_MODEL_PATH=configs/casbin_model.conf
CASBIN_POLICY_PATH=configs/casbin_policy.csv

# Mapping role berbasis email
RBAC_ADMIN_EMAILS=admin@mail.com
RBAC_STAFF_EMAILS=staff@mail.com,ops@mail.com
```

### 2 Login dan akses endpoint RBAC

1. Login pakai endpoint auth untuk dapat access token.
2. Kirim header `Authorization: Bearer <access_token>`.
3. Panggil endpoint `/products`, `/orders`, `/admin/users` sesuai role.

Jika role tidak punya akses -> response `403 FORBIDDEN`.

## License

MIT
