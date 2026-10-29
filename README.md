# Laritmo

Educational portal: courses with lectures, labs, grade sheets and exam questions. Backend in Go, frontend in Vue 3.

## Features

- Courses with lectures (Markdown, syntax highlighting), labs with GitHub links, grade sheets and exam question banks
- Random exam ticket selection and ticket document generation
- JWT authentication with three roles: owner, admin, student
- User management for the owner
- Full database dump to JSON for admins
- Light and dark themes, responsive layout
- Session cache for read-only data on the client

## Stack

- Backend: Go 1.25, Gin, Squirrel, MariaDB 11, goose migrations
- Frontend: Vue 3, TypeScript, Pinia, PrimeVue, Tailwind CSS, Vite

## Local setup

Requirements: Go 1.25+, Node.js 20+, Docker (or a local MariaDB 11), [mkcert](https://github.com/FiloSottile/mkcert), [goose](https://github.com/pressly/goose).

1. Start the database:

   ```bash
   docker compose -f docker-compose.local.yml up -d mariadb
   ```

   This creates the `edu_portal` database with user `eduuser` / `edupass`, matching `src/back/configs/config.local.yaml`. With a local MariaDB, create the database and user manually with the same names.

2. Generate local HTTPS certificates:

   ```bash
   mkcert -install
   mkdir -p src/back/certs
   cd src/back/certs
   mkcert -key-file localhost-key.pem -cert-file localhost.pem localhost 127.0.0.1 ::1
   ```

3. Apply migrations:

   ```bash
   cd src/back
   goose -dir migrations mysql "eduuser:edupass@tcp(localhost:3306)/edu_portal" up
   ```

4. Create the first user and promote it to owner. `createadmin` creates a user with the `admin` role:

   ```bash
   go run ./cmd/createadmin
   ```

   ```sql
   UPDATE users SET role = 'owner' WHERE username = '<your username>';
   ```

5. Build the frontend. The output goes to `src/back/web/`, which the backend serves:

   ```bash
   cd src/front
   npm install
   npm run build
   ```

6. Run the server:

   ```bash
   cd src/back
   go run ./cmd/server
   ```

The app is available at https://localhost:8443, the login page at `/auth`, Swagger at `/swagger/index.html` (debug mode only) and the health check at `/health`.

### Frontend development

Run the backend as above, then in a second terminal:

```bash
cd src/front
npm run dev
```

The dev server runs on http://localhost:5173 and proxies API requests to the backend.

## Project layout

```
src/back/
  cmd/server/       entry point and routing
  cmd/createadmin/  CLI for creating a user
  internal/         handlers, repository, middleware, models, auth, services
  configs/          local and production configs
  migrations/       goose migrations
  seed/             development seed data
  web/              built frontend (generated, not in git)
src/front/
  src/views/        pages
  src/components/   shared components
  src/stores/       Pinia stores
  src/api/          API client and session cache
```

## Roles

| Role    | Access                                                              |
|---------|---------------------------------------------------------------------|
| owner   | Everything an admin can do, plus user management                    |
| admin   | Create, edit and delete content; generate tickets; download DB dump |
| student | Read content, random tickets, lecture export                        |

The owner account cannot be edited or deleted through the UI.

## API

Interactive documentation is in Swagger. Main groups:

- Public: `POST /api/auth/login`, `GET /api/courses`, `/api/lectures`, `/api/labs`, `/api/grade-sheets`, `/api/exam-questions`, `GET /api/courses/:id/tickets/random`
- Authenticated: `GET /api/lectures/export`, `PUT /api/auth/me/password`
- Admin (`/api/admin/*`): create, update and delete for courses, lectures, labs, grade sheets and exam questions; `POST /api/admin/courses/:id/tickets/generate`; `GET /api/admin/dump`
- Owner (`/api/admin/users/*`): user CRUD, activate, deactivate, password reset

## Client cache

Courses, lectures, labs, grade sheets and exam questions rarely change, so GET responses for them are cached in `sessionStorage` for up to 30 minutes. Any successful write to `/api/admin/*`, logout or a 401 clears the cache. See `src/front/src/api/cache.ts`.

## Database dump

Admins and the owner get a "Дамп БД" button in the header. It downloads `GET /api/admin/dump` as JSON with all tables (courses, lectures, labs, grade sheets, exam questions, users). Password hashes are left out. The dump is not a restorable backup; for backups use `mysqldump` against the database directly.

## Deployment

The `Dockerfile` builds the frontend and backend into one image. `docker-compose.prod.yml` runs it with MariaDB; secrets (`LARITMO_DATABASE_PASSWORD`, `LARITMO_AUTH_JWT_SECRET`, `MYSQL_ROOT_PASSWORD`) are passed through environment variables. Production config is `src/back/configs/config.prod.yaml`.

## Troubleshooting

- Connection refused: check `curl -k https://localhost:8443/health` and that `src/back/certs/` contains both PEM files.
- Database connection failed: check that the container is up (`docker compose -f docker-compose.local.yml ps`) and the credentials match `configs/config.local.yaml`.
- Swagger returns 500: regenerate the docs with `cd src/back && swag init -g cmd/server/main.go`.
- Certificate warning in the browser: run `mkcert -install` and regenerate the certificates.
