# 🍄 Laritmo - Forest Academy

Educational portal with whimsical frog and mushroom theme. Built with Go and Vue 3.

## ✨ Features

- 📚 Course management (lectures, labs, grade sheets, exam questions)
- 🔐 JWT authentication with role-based access (Admin/Student)
- 🎨 Beautiful UI with mushroom and frog theme with animations
- 📱 Responsive design
- 🔬 GitHub integration for lab assignments
- 📝 Markdown support with syntax highlighting
- 🐸 Frog animation on page load

## 🛠️ Tech Stack

**Backend:** Go 1.25 + Gin + MariaDB  
**Frontend:** Vue 3 + TypeScript + PrimeVue + Tailwind CSS

---

## 🚀 Local Setup

### Prerequisites

Install these first:
- **Go 1.23+**: https://go.dev/dl/
- **Node.js 20+**: https://nodejs.org/
- **MariaDB 11**: https://mariadb.org/download/ (or use Docker)
- **mkcert**: for local HTTPS certificates
```bash
# Install mkcert
brew install mkcert              # macOS
choco install mkcert             # Windows
sudo apt install mkcert          # Linux

# Create trusted certificates
mkcert -install
```

---

### Step 1: Clone & Setup Database
```bash
# Clone repository
git clone https://github.com/your-username/laritmo.git
cd laritmo

# Option A: Using Docker (easiest)
docker compose -f docker-compose.local.yml up -d mariadb

# Option B: Using local MariaDB
mysql -u root -p
```
```sql
CREATE DATABASE edu_portal CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'eduuser'@'localhost' IDENTIFIED BY 'edupass';
GRANT ALL PRIVILEGES ON edu_portal.* TO 'eduuser'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

---

### Step 2: Generate SSL Certificates
```bash
# Create certs directory in backend
mkdir -p src/back/certs
cd src/back/certs

# Generate local HTTPS certificates
mkcert -key-file localhost-key.pem -cert-file localhost.pem localhost 127.0.0.1 ::1

cd ../../..
```

---

### Step 3: Configure Backend
```bash
cd src/back

# Config file already exists: configs/config.local.yaml
# Edit if needed (default values work out of the box)
nano configs/config.local.yaml
```

**configs/config.local.yaml** (default):
```yaml
server:
  host: localhost
  mode: debug
  port: 8443
  use_tls: true
  tls_cert_file: "./certs/localhost.pem"
  tls_key_file: "./certs/localhost-key.pem"

database:
  host: localhost
  port: 3306
  user: eduuser
  password: edupass
  name: edu_portal

auth:
  jwt_secret: "super-secret-key-change-in-production-abc123"
  jwt_expiration_hours: 168
```

---

### Step 4: Run Database Migrations
```bash
# Install goose (if not installed)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
cd src/back
goose -dir migrations mysql "eduuser:edupass@tcp(localhost:3306)/edu_portal" up

# Verify
goose -dir migrations mysql "eduuser:edupass@tcp(localhost:3306)/edu_portal" status
```

### Step 5: Create Admin User
```bash
cd src/back

# Run admin creation tool
go run cmd/createadmin/main.go

# Enter credentials when prompted:
# Username: admin
# Email: admin@example.com
# Password: admin123
```

---

### Step 6: Build Frontend
```bash
cd src/front

# Install dependencies
npm install

# Build for production (goes to src/back/web/)
npm run build
```

---

### Step 7: Start Backend
```bash
cd src/back

# Install Go dependencies
go mod download

# Run server
go run cmd/server/main.go
```

You should see:
```
{"level":"INFO","msg":"Config loaded successfully"}
{"level":"INFO","msg":"HTTPS server started","url":"https://localhost:8443"}
{"level":"INFO","msg":"Swagger UI available","url":"https://localhost:8443/swagger/index.html"}
```

---

### Step 8: Access Application

Open in browser:
- **App**: https://localhost:8443
- **Login Page**: https://localhost:8443/auth
- **API Docs**: https://localhost:8443/swagger/index.html
- **Health Check**: https://localhost:8443/health

**Login:**
- Navigate to `/auth` to access the login page
- Username: `admin` (or whatever you created with `createadmin`)
- Password: `admin123` (or whatever you set)

> ⚠️ Browser will show security warning for self-signed certificate - click "Advanced" → "Proceed to localhost"  
> ⚠️ **Important:** Create admin user first using `go run cmd/createadmin/main.go` before logging in!

---

## 🔧 Development Mode

For frontend development with hot reload:
```bash
# Terminal 1: Backend
cd src/back
go run cmd/server/main.go

# Terminal 2: Frontend dev server
cd src/front
npm run dev
```

Frontend will run on `http://localhost:5173` with API proxy to backend.

---

## 📁 Project Structure
```
laritmo/
├── src/
│   ├── back/                   # Go backend
│   │   ├── cmd/server/         # Main application
│   │   ├── internal/           # Business logic
│   │   │   ├── handlers/       # HTTP handlers
│   │   │   ├── repository/     # Database layer
│   │   │   ├── middleware/     # Auth middleware
│   │   │   ├── models/         # Data models
│   │   │   └── auth/           # JWT manager
│   │   ├── configs/            # Configuration files
│   │   ├── migrations/         # Database migrations
│   │   ├── certs/              # SSL certificates (local)
│   │   └── web/                # Built frontend (generated)
│   └── front/                  # Vue 3 frontend
│       ├── src/
│       │   ├── views/          # Page components
│       │   ├── components/     # Reusable components
│       │   ├── stores/         # Pinia stores
│       │   ├── router/         # Vue Router
│       │   └── api/            # API client
│       └── public/             # Static assets
├── docker-compose.local.yml     # Docker setup for local development
└── README.md
```

---

## 🐳 Docker Setup (Alternative)
```bash
# Start database with Docker
docker compose -f docker-compose.local.yml up -d mariadb

# View logs
docker compose -f docker-compose.local.yml logs -f mariadb

# Stop
docker compose -f docker-compose.local.yml down

# Stop and remove volumes (clean slate)
docker compose -f docker-compose.local.yml down -v
```

---

## 🧪 API Documentation

Swagger UI available at: https://localhost:8443/swagger/index.html (debug mode only)

### Key Endpoints:

**Public:**
- `POST /auth/login` - Login
- `GET /api/courses` - List courses
- `GET /api/lectures/:id` - Get lecture
- `GET /api/labs/:id` - Get lab

**Admin (requires JWT):**
- `POST /api/admin/courses` - Create course
- `PUT /api/admin/lectures/:id` - Update lecture
- `DELETE /api/admin/labs/:id` - Delete lab

---

## 🔐 Authentication

Authentication is handled via a dedicated login page at `/auth`:
1. Navigate to `/auth` to access the login form
2. Enter credentials → Receive JWT token
3. Token stored in localStorage → Automatically attached to requests
4. Redirected to home page after successful login

Admin features (Create/Edit/Delete) only visible when logged in as admin.

---

## 🎨 Theme

Whimsical forest theme with:
- 🍄 Mushroom icons and decorations
- 🐸 Frog mascot with entrance animation
- 🌿 Nature-inspired color palette
- Playful yet educational design
- Smooth animations and transitions

---

## 📝 License

MIT License - feel free to use for your own educational projects!

---

## 🤝 Contributing

This is a learning project, but suggestions and improvements are welcome!

---

## ❓ Troubleshooting

**"Connection refused" when accessing app:**
- Check backend is running: `curl https://localhost:8443/health`
- Check certificates exist: `ls src/back/certs/`

**"Database connection failed":**
- Verify MariaDB is running: `docker compose -f docker-compose.local.yml ps` or `mysql -u root -p`
- Check credentials in `configs/config.local.yaml`
- If using Docker, ensure container is up: `docker compose -f docker-compose.local.yml up -d mariadb`

**"Swagger returns 500":**
- Regenerate docs: `cd src/back && swag init -g cmd/server/main.go`

**Browser security warning:**
- This is normal for self-signed certificates
- Click "Advanced" → "Proceed to localhost (unsafe)"
- Or use mkcert for trusted certificates

---

## 📧 Contact

Questions? Open an issue on GitHub!
