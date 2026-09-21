# go-tasks-api

API REST de gestión de tareas con autenticación JWT, hecha en Go y PostgreSQL.

## Estructura del proyecto

```
go-tasks-api/
├── cmd/
│   └── api/
│       └── main.go          # Ejecutable del proyecto
├── internal/
│   ├── config/               # carga de variables de entorno
│   ├── database/             # conexion a la base de datos
│   ├── models/                # Tablas de la base de datos
│   ├── handlers/              # logica de cada endpoint
│   ├── middleware/            # validacion de JWT
│   └── routes/                 # definicion de rutas
├── docker-compose.yml          # levanta Base de datos
├── .env.example
└── go.mod
```

## Como iniciar

1. **Copia el archivo de entorno:**
   ```bash
   cp .env.example .env
   ```

2. **Levanta la base de datos con Docker:**
   ```bash
   docker compose up -d
   ```

3. **Instala las dependencias de Go:**
   ```bash
   go mod tidy
   ```

4. **Corre el servidor:**
   ```bash
   go run cmd/api/main.go
   ```
