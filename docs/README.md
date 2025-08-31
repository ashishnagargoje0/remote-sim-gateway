# Remote SIM Gateway Documentation

Welcome to the Remote SIM Gateway documentation! This comprehensive guide will help you understand, install, configure, and use the Remote SIM Gateway system.

## 📚 Documentation Index

### Getting Started
- [Setup Guide](SETUP.md) - Complete installation and configuration instructions
- [Quick Start](#quick-start) - Get up and running in 10 minutes
- [Architecture Overview](#architecture) - Understand how the system works

### User Guides
- [Web Dashboard Guide](#web-dashboard) - How to use the web interface
- [Android App Guide](#android-app) - Setting up the mobile bridge
- [SMS Management](#sms-management) - Sending and managing SMS messages
- [Call Management](#call-management) - Making and tracking calls

### Developer Resources
- [API Documentation](API.md) - Complete REST API reference
- [WebSocket Protocol](#websocket-protocol) - Real-time communication specs
- [Contributing Guide](CONTRIBUTING.md) - How to contribute to the project
- [Security Guide](SECURITY.md) - Security best practices

### Operations
- [Deployment Guide](DEPLOYMENT.md) - Production deployment instructions
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions
- [Monitoring](#monitoring) - Health checks and monitoring setup

## 🚀 Quick Start

### Prerequisites
- Android phone (API 21+) with active SIM card
- Go 1.21+ installed
- Node.js 18+ installed
- PostgreSQL database
- Domain with SSL certificate (for production)

### 1. Clone and Setup
```bash
git clone https://github.com/yourusername/remote-sim-gateway.git
cd remote-sim-gateway

# Copy environment files
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local

# Edit configuration files with your settings
```

### 2. Start Backend
```bash
cd backend
go mod download
go run cmd/server/main.go
```

### 3. Start Frontend
```bash
cd frontend
npm install
npm run dev
```

### 4. Install Android App
```bash
cd android-app
# Open in Android Studio
# Build and install APK on your phone
```

## 🏗️ Architecture

The Remote SIM Gateway consists of four main components:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌─────────────┐
│                 │    │                  │    │                 │    │             │
│   Web App       │◄──►│  Backend Server  │◄──►│   Android App   │◄──►│ SIM Card    │
│   (Dashboard)   │    │   (Go/Golang)   │    │   (Kotlin)      │    │ (SMS/Calls) │
│                 │    │                  │    │                 │    │             │
└─────────────────┘    └──────────────────┘    └─────────────────┘    └─────────────┘
```

### Component Details

1. **Web Dashboard (React/Next.js)**
   - User interface for controlling the system
   - Real-time status monitoring
   - SMS and call management
   - User authentication and authorization

2. **Backend Server (Go)**
   - RESTful API endpoints
   - WebSocket server for real-time communication
   - Database management
   - Authentication and security

3. **Android Bridge App (Kotlin)**
   - Connects to backend via WebSocket
   - Manages SMS and call operations
   - Handles device permissions
   - Background service for 24/7 operation

4. **Database (PostgreSQL)**
   - User management
   - Message and call history
   - Device registration
   - Audit logs

## 🌐 Web Dashboard

### Features
- **Dashboard**: Overview of system status and recent activity
- **SMS Management**: Send individual or bulk SMS messages
- **Call Control**: Initiate calls remotely
- **Device Management**: Monitor and manage connected devices
- **History & Analytics**: Track all communications
- **User Management**: Multi-user support with role-based access

### Navigation
- **Dashboard** (`/`) - Main overview page
- **SMS** (`/sms/*`) - SMS sending and history
- **Calls** (`/calls/*`) - Call management and history
- **Devices** (`/devices`) - Device status and management
- **Analytics** (`/analytics`) - Usage statistics and charts
- **Admin** (`/admin/*`) - Administration panel (admin only)

## 📱 Android App

### Key Features
- **Background Service**: Maintains connection even when app is closed
- **Auto-reconnect**: Handles network interruptions gracefully
- **Battery Optimized**: Minimal power consumption
- **Permission Management**: Requests only necessary permissions
- **Security**: Device authentication and encryption

### Required Permissions
- `SEND_SMS` - Send text messages
- `READ_PHONE_STATE` - Access phone status
- `CALL_PHONE` - Make phone calls
- `INTERNET` - Network connectivity
- `WAKE_LOCK` - Keep service running

### Setup Process
1. Install APK on Android device
2. Grant required permissions
3. Configure backend server URL
4. Enter device authentication token
5. Start background service

## 📊 SMS Management

### Single SMS
- Enter phone number with country code
- Compose message (up to 1600 characters)
- Select target device
- Send and track delivery status

### Bulk SMS
- Upload list of phone numbers (max 100)
- Compose single message for all recipients
- Monitor sending progress in real-time
- View detailed delivery reports

### Message Features
- **Long SMS Support**: Messages over 160 chars split automatically
- **Delivery Tracking**: Real-time status updates
- **History Logs**: Complete message history with search/filter
- **Export Options**: CSV export for reporting

## 📞 Call Management

### Making Calls
- Enter destination phone number
- Select device to make call from
- Monitor call status in real-time
- Track call duration and outcomes

### Call Features
- **Remote Dialing**: Initiate calls from web interface
- **Call History**: Track all outbound calls
- **Duration Tracking**: Automatic call timing
- **Status Updates**: Real-time call progress

## 🔐 Security Features

### Authentication
- JWT token-based authentication
- Password hashing with bcrypt
- Session management and expiration
- Role-based access control

### Communication Security
- TLS encryption for all data in transit
- Secure WebSocket connections (WSS)
- Device authentication and whitelisting
- Input validation and sanitization

### Rate Limiting
- Global API rate limiting (100 requests/hour)
- SMS-specific limits (10 SMS/hour)
- IP-based and user-based limiting
- Configurable limits per user role

## 🔧 Configuration

### Environment Variables

#### Backend (.env)
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=simgateway
DB_USER=postgres
DB_PASSWORD=your_password

# JWT
JWT_SECRET=your-super-secret-key
JWT_EXPIRATION=24

# Server
PORT=8080
CORS_ORIGINS=http://localhost:3000

# Security
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600
```

#### Frontend (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
NEXT_PUBLIC_APP_NAME=Remote SIM Gateway
```

### Database Setup
```sql
-- Create database
CREATE DATABASE simgateway;

-- Create user (optional)
CREATE USER simgateway_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE simgateway TO simgateway_user;
```

## 🚀 Deployment

### Production Checklist
- [ ] Configure SSL certificates
- [ ] Set up domain name
- [ ] Configure environment variables
- [ ] Set up database backups
- [ ] Configure monitoring
- [ ] Set up log rotation
- [ ] Test all functionality

### Docker Deployment
```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d

# Development deployment
docker-compose -f docker-compose.dev.yml up -d
```

### Manual Deployment
See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed instructions.

## 🐛 Troubleshooting

### Common Issues
- **Android app won't connect**: Check network, SSL, and device whitelist
- **SMS not sending**: Verify permissions, SIM status, and rate limits
- **WebSocket disconnections**: Check network stability and SSL configuration
- **Authentication errors**: Verify JWT secret and token expiration

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for detailed solutions.

## 📞 Support

### Getting Help
- **Documentation**: Check this docs folder first
- **GitHub Issues**: Report bugs and request features
- **Discussions**: Ask questions in GitHub Discussions
- **Email**: Contact maintainers directly

### Contributing
We welcome contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### License
This project is licensed under the MIT License. See [LICENSE](../LICENSE) for details.

---

## 📋 Next Steps

1. **Read the [Setup Guide](SETUP.md)** for detailed installation instructions
2. **Check the [API Documentation](API.md)** for development reference
3. **Review [Security Guidelines](SECURITY.md)** for production deployment
4. **Explore [Troubleshooting](TROUBLESHOOTING.md)** for common issues

Happy building! 🚀