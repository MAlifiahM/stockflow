# StockFlow

Inventory & Sales Management System

## Tech Stack
- Backend: Go with go fiber
- Frontend: React + TypeScript
- Database: PostgreSQL

## Run

### All
- copy .env in backend and fill the value
    ```bash
    cd backend
    cp .env.example .env
    ```
- copy .env from .env.example and fill value for database same with backend
    ```bash
    cp .env.example .env
    ```
- copy .env from .env.example from frontend and fill url with backend url
    ```bash
    cp .env.example .env
    ```
- run all frontend and backend with migrations and seeders
    ```bash
    docker compose up -d
    ```

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
### Frontend only
- copy .env from .env.example in frontend folder
    ```bash
    cp .env.example .env
    ```
- update url for backend using url backend + port
- run frontend only
    ```bash
    docker compose up -d frontend
    ```