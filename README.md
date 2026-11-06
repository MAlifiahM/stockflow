# StockFlow

Inventory & Sales Management System

## Tech Stack
- Backend: Go with go fiber
- Frontend: React + TypeScript
- Database: PostgreSQL

## Run

### Database with migration and seeder only
```bash
docker compose up -d db migrate seed
```
### Backend only
```bash
docker compose up -d backend
```