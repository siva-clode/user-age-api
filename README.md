# User Age API (GoFiber + sqlc + Postgres)

## Features
- Create/Get/Update/Delete user with `name` and `dob`
- `age` is calculated dynamically when returning user(s)
- Uses sqlc for DB access
- Input validation with go-playground/validator
- Zap logging
- Middleware: request-id and request logging
- Optional Docker compose for Postgres

## Setup (local)

1. Install dependencies:
   - Go >=1.20
   - sqlc (https://sqlc.dev)
   - PostgreSQL (or use Docker Compose below)

2. Run Postgres (option A: Docker)
   ```bash
   docker-compose up -d
