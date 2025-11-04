# TP Grupal Desarrollo de Software 2025 - Grupo 1

- **frontend** → React + Vite  
- **api** → Go

## Crear archivo .env en /api
```bash
# Postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=12345678
DB_NAME=proyecto-desarrollo
DB_SSLMODE=disable

# puerto
APP_PORT=8080

# (development/production)
APP_ENV=development

# Generar documentacion
GENERATE_DOCS=0

#JWT
JWT_SECRET=asd23wfsfsdfsdfsdfsdfsdf
JWT_TTL_MIN=30

# Datos ADMIN
ADMIN_EMAIL=admin@quepc.com
ADMIN_PASSWORD=12345678
ADMIN_NAME=Administrador General

````

## Frontend
```bash
cd frontend/
npm install 
npm run dev
````

## API

```bash
cd api/
go run .
```
