# Setup Guide - Remote SIM Gateway

This guide will walk you through setting up the Remote SIM Gateway from scratch. Follow these steps carefully to ensure a smooth installation.

## 📋 Prerequisites

### System Requirements
- **Operating System**: Linux, macOS, or Windows
- **RAM**: Minimum 2GB (4GB recommended)
- **Storage**: 5GB free space
- **Network**: Stable internet connection

### Software Requirements
- **Go**: Version 1.21 or higher
- **Node.js**: Version 18 or higher
- **PostgreSQL**: Version 12 or higher
- **Git**: For cloning the repository
- **Android Studio**: For building the Android app (optional)

### Hardware Requirements
- **Android Phone**: API level 21+ (Android 5.0+)
- **Active SIM Card**: With SMS and calling capabilities
- **Domain & SSL**: For production deployment (optional for development)

## 🚀 Installation Steps

### Step 1: Clone Repository

```bash
# Clone the repository
git clone https://github.com/yourusername/remote-sim-gateway.git
cd remote-sim-gateway

# Check the project structure
ls -la
```

### Step 2: Database Setup

#### Install PostgreSQL
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# Windows
# Download and install from https://www.postgresql.org/download/windows/
```

#### Create Database
```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Create database and user
CREATE DATABASE simgateway;
CREATE USER simgateway_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE simgateway TO simgateway_user;

# Exit PostgreSQL
\q
```

#### Test Database Connection
```bash
psql -h localhost -U simgateway_user -d simgateway
# Enter password when prompted
# If successful, you'll see the psql prompt
```

### Step 3: Backend Setup

#### Install Dependencies
```bash
cd backend

# Download Go modules
go mod download

# Verify Go installation
go version
```

#### Configure Environment
```bash
# Copy environment template
cp .env.example .env

# Edit configuration
nano .env
```

**Edit `.env` file:**
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=simgateway
DB_USER=simgateway_user
DB_PASSWORD=your_secure_password
DB_SSL_MODE=disable

# JWT Configuration (Generate a strong secret!)
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_EXPIRATION=24

# Server Configuration
PORT=8080
GIN_MODE=release

# CORS Configuration
CORS_ORIGINS=http://localhost:3000,https://yourdomain.com

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600
```

#### Start Backend Server
```bash
# Development mode
go run cmd/server/main.go

# Or build and run
go build -o sim-gateway cmd/server/main.go
./sim-gateway
```

**Expected output:**
```
Database connected successfully
Database migrations completed successfully
Server starting on port 8080
```

#### Test Backend
```bash
# Test health endpoint
curl http://localhost:8080/health

# Expected response: {"status":"ok"}
```

### Step 4: Frontend Setup

#### Install Dependencies
```bash
cd ../frontend

# Install npm packages
npm install

# Verify Node.js installation
node --version
npm --version
```

#### Configure Environment
```bash
# Copy environment template
cp .env.local.example .env.local

# Edit configuration
nano .env.local
```

**Edit `.env.local` file:**
```bash
# Backend API URL
NEXT_PUBLIC_API_URL=http://localhost:8080

# WebSocket URL
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws

# App Configuration
NEXT_PUBLIC_APP_NAME=Remote SIM Gateway
NEXT_PUBLIC_APP_VERSION=1.0.0

# Development settings
NEXT_PUBLIC_DEBUG=true
```

#### Start Frontend Server
```bash
# Development mode
npm run dev

# Production build
npm run build
npm start
```

**Expected output:**
```
ready - started server on 0.0.0.0:3000, url: http://localhost:3000
```

#### Test Frontend
1. Open browser to `http://localhost:3000`
2. You should see the login page
3. Try registering a new account

### Step 5: Android App Setup

#### Prerequisites
- Android Studio installed
- Android device with USB debugging enabled
- Phone connected to computer via USB

#### Build Android App
```bash
cd ../android-app

# Open in Android Studio
studio .

# Or build from command line
./gradlew assembleDebug
```

#### Install on Device
```bash
# Install debug APK
adb install app/build/outputs/apk/debug/app-debug.apk

# Or use Android Studio's run button
```

#### Configure Android App
1. Open the app on your phone
2. Grant required permissions:
   - SMS permissions
   - Phone permissions
   - Network access
3. Enter backend server URL: `ws://your-server:8080/ws`
4. Enter device ID (generate unique ID)
5. Test connection

### Step 6: First Test

#### Register User Account
1. Open web dashboard (`http://localhost:3000`)
2. Click "Sign up for free"
3. Enter email and password
4. Complete registration

#### Connect Android Device
1. Open Android app
2. Enter server URL and device ID
3. Tap "Connect"
4. Verify connection in web dashboard

#### Send First SMS
1. Go to "Send SMS" in web dashboard
2. Select your connected device
3. Enter a phone number (your own for testing)
4. Type a test message
5. Click "Send SMS"
6. Check your phone for the message

## 🔧 Advanced Configuration

### SSL Certificate Setup (Production)

#### Using Let's Encrypt
```bash
# Install certbot
sudo apt install certbot

# Generate certificate
sudo certbot certonly --standalone -d yourdomain.com

# Configure nginx
sudo nano /etc/nginx/sites-available/simgateway
```

#### Nginx Configuration
```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Database Optimization

#### PostgreSQL Tuning
```bash
# Edit PostgreSQL configuration
sudo nano /etc/postgresql/14/main/postgresql.conf
```

```sql
# Performance tuning
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
```

### Monitoring Setup

#### Prometheus Configuration
```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'sim-gateway'
    static_configs:
      - targets: ['localhost:8080']
```

#### Grafana Dashboard
```bash
# Install Grafana
sudo apt install grafana

# Start Grafana
sudo systemctl start grafana-server
sudo systemctl enable grafana-server

# Access at http://localhost:3000
# Default credentials: admin/admin
```

## 🔐 Security Hardening

### Production Security Checklist
- [ ] Change default JWT secret
- [ ] Enable SSL/TLS encryption
- [ ] Configure firewall rules
- [ ] Set up rate limiting
- [ ] Enable audit logging
- [ ] Regular security updates
- [ ] Backup encryption
- [ ] Access control policies

### Firewall Configuration
```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp