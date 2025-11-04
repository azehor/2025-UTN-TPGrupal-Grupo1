# Build stage
FROM node:20-alpine AS build
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
ARG VITE_API_BASE_URL=/v1
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
RUN npm run build

# Runtime stage (Caddy sirve estáticos)
FROM caddy:2.8-alpine
# Caddy por defecto sirve /usr/share/caddy en :80
COPY --from=build /app/dist /usr/share/caddy
EXPOSE 80