# Frontend environment configuration (Vite)

Vite lee variables desde archivos `.env` según el modo. Las variables deben estar prefijadas con `VITE_` para ser accesibles en el código mediante `import.meta.env`.

## Archivos recomendados

- Desarrollo: `frontend/.env.development` (incluido en el repo) con valores cómodos para dev.
- Producción: evita `.env` en el servidor en runtime. Inyecta `VITE_*` en el entorno del proceso de build (CI/CD) antes de `npm run build`. Si necesitas un archivo de ejemplo, usa `frontend/.env.production.example` (incluido).
- No comitees archivos locales o productivos: `.env.local`, `.env.*.local` y opcionalmente `.env.production` están ignorados en `.gitignore`.

## Variables típicas

```ini
# API base URL para el backend (incluye /v1)
VITE_API_BASE_URL=http://localhost:8080/v1
```

## Uso en código

```ts
// src/lib/api.ts
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || `${window.location.origin}/v1`;
```

## Notas

- Cualquier variable de frontend queda embebida en el bundle y es visible públicamente; no coloques secretos aquí.
- Si el frontend sirve desde el mismo dominio que la API, puedes omitir `VITE_API_BASE_URL` y usar rutas relativas (`/v1`). Para desarrollo, podrías configurar un proxy en `vite.config.ts` si lo prefieres.
