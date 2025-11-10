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
# TP Grupal - Desarrollo de Software (2025) - Grupo 1

Proyecto full-stack que reúne una API en Go y un frontend en React (Vite). El objetivo es gestionar componentes y softwares relacionados a carreras y recomendaciones.

Stack principal
- Frontend: React + Vite + TypeScript
- Backend: Go (chi, gorm) + PostgreSQL

Índice
1. Descripción
2. Requisitos
3. Instalación
	 - Configurar variables de entorno
	 - Iniciar backend
	 - Iniciar frontend
4. Uso (rápido)
5. Estructura de assets e imágenes
6. Troubleshooting (errores comunes)
7. Licencia
8. Insignias

## 1. Descripción

Este repositorio contiene una aplicación para administrar componentes, softwares y recomendaciones académicas/técnicas. El panel de administración permite crear, editar y listar entidades (softwares, componentes, carreras, etc.). La API sirve endpoints REST bajo `/v1` y expone assets estáticos bajo `/static`.

## 2. Requisitos

- Go (1.20+ recomendado, ver `go.mod`)
- Node.js (16+ recomendado) y npm
- PostgreSQL (servidor local o remoto)
- psql (cliente, opcional para tareas de migración/seed)

## 3. Instalación

3.1 Clonar el repositorio

```powershell
git clone <repo-url>
cd 2025-UTN-TPGrupal-Grupo1
```

3.2 Configurar variables de entorno (backend)

Crear un archivo `.env` en la carpeta `api/` con al menos las variables de conexión a PostgreSQL y otras opcionales. Ejemplo mínimo:

```powershell
# api/.env (ejemplo mínimo)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=your_db_name
DB_SSLMODE=disable
APP_PORT=8080
APP_ENV=development
```

Opcionales (JWT, ADMIN seed, assets base url):

```powershell
GENERATE_DOCS=0
JWT_SECRET=replace_with_secure_secret
JWT_TTL_MIN=30
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=changeme
API_ASSETS_BASE_URL=
```

3.3 Crear la base de datos (si no existe)

Use `psql` o una herramienta gráfica. Ejemplo con `psql` (PowerShell):

```powershell
& "C:\Program Files\PostgreSQL\17\bin\psql.exe" -U postgres -c "CREATE DATABASE your_db_name;"
```

3.4 Instalar dependencias y ejecutar

Backend (desde la carpeta `api`):

```powershell
cd api
go mod download
go run .
# o para desarrollo: go run main.go
```

Frontend (desde la carpeta `frontend`):

```powershell
cd frontend
npm install
npm run dev
```

Por defecto el frontend usa Vite y arrancará en `http://localhost:5173` (si ese puerto está en uso, Vite elegirá otro puerto libre). La API corre en `http://localhost:8080` según `.env`.

## 4. Uso (rápido)

- Abrir el frontend en el navegador (p. ej. `http://localhost:5173` o el puerto indicado por Vite).
- Iniciar sesión en el panel de administración con las credenciales de admin si las creó en `.env` (o usar el seed de admin si existe).
- Crear entidades (softwares, componentes, carreras) desde el Panel de Admin.

API endpoints principales (ejemplos):

- Listar carreras: `GET /v1/carreras`
- Crear software: `POST /v1/softwares` (payload en JSON)
- Assets estáticos: `GET /static/images/<categoria>/<archivo>`

> Nota: la API espera JSON para la mayoría de endpoints. Si la funcionalidad de subida de imágenes no está configurada completamente en el backend, use el campo `image_url` para indicar una URL pública en lugar de subir archivos.

## 5. Estructura de assets e imágenes

Estructura utilizada en el proyecto (ubicación en el servidor):

```
api/static/
	images/
		carreras/
		juegos/
		softwares/
```

- Cuando se almacena la ruta en la base de datos, suele guardarse la ruta relativa dentro de `static`, por ejemplo: `images/carreras/imagen.png`.
- El servidor expone estos archivos bajo `/static/...`. Si configura `API_ASSETS_BASE_URL`, las respuestas generadas usarán esa base.

## 6. Troubleshooting (errores comunes)

- Error: `API Error 400: {"error":"JSON inválido o vacío"}`
	- Causa frecuente: el frontend envía `FormData` (multipart) mientras que el backend está esperando JSON en el body.
	- Soluciones:
		- Opción rápida: usar el campo `image_url` en lugar de adjuntar una imagen, enviar sólo JSON.
		- Implementación completa: actualizar el backend para leer `multipart/form-data` y procesar el campo `data` (JSON) y `imagen` (file). Esto implica usar `r.ParseMultipartForm(...)`, `r.FormValue("data")` y `r.FormFile("imagen")` antes de decodificar el JSON.

- Error: `FATAL: no existe la base de datos "<name>"` al iniciar el servidor
	- Asegúrese de que la base de datos existe y coincide con `DB_NAME` en `.env`. Cree la base de datos con `psql` si es necesario.

- Error de CORS entre frontend y backend
	- Verifique que el backend permite el origen donde corre el frontend (por ejemplo `http://localhost:5173`). Ajuste la lista de AllowedOrigins en el middleware CORS del servidor.

- Problema: puerto 5173 en uso (Vite)
	- Vite intentará otro puerto por defecto; revise la URL que imprime en la consola al iniciar.

Si necesita ayuda con un error específico, incluya el mensaje de error completo y el contexto (qué acción realizaba, payload enviado, y logs del backend si los hay).


## 7. Insignias

- Go version: 1.25.1
- Node version: 22.19.0
- Npm version: 10.9.3


