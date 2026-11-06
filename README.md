# StockFlow

Inventory & Sales Management System

## Tech Stack
- Backend: Go with go fiber
- Frontend: React + TypeScript
- Database: PostgreSQL

## Run

### Database with migration and seeder only
- copy .env from .env.example
    ```bash
    cp .env.example .env
    ```

- update content .env
- run docker compose for database, migrate and seed
    ```bash
    docker compose up -d db migrate seed
    ```
### Backend only
- copy .env from .env.example in backend folder
    ```bash
    cp .env.example .env
    ```
- update content .env in backend folder
- run backend after running database with migrate and seed
    ```bash
    docker compose up -d backend
    ```