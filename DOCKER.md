# 🐳 TutupLapak API Docker Setup

This document provides comprehensive instructions for running the TutupLapak API using Docker and Docker Compose.

## 📋 Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd tutuplapak
```

### 2. Environment Configuration

```bash
# Copy environment template
cp env.example .env

# Edit environment variables (optional)
nano .env
```

### 3. Start Services

```bash
# Using Makefile (recommended)
make docker-up-dev

# Or using Docker Compose directly
docker-compose -f docker-compose.dev.yml up -d
```

### 4. Verify Services

```bash
# Check service status
make docker-status

# View logs
make docker-logs
```

## 🛠️ Available Services

### Core Services
- **tutuplapak-api**: Main API application (Port 8080)
- **postgres**: PostgreSQL database (Port 5433)
- **minio**: MinIO object storage (Ports 9000, 9001)
- **redis**: Redis cache (Port 6379)

### Development Tools
- **pgadmin**: Database administration (Port 5050)
- **redis-commander**: Redis administration (Port 8081)

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `PORT` | API server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | Auto-generated |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |
| `MINIO_ENDPOINT` | MinIO server endpoint | `minio:9000` |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin123` |
| `MINIO_BUCKET_NAME` | MinIO bucket name | `tutuplapak-files` |

## 🎯 Makefile Commands

### Development Commands
```bash
make docker-up-dev          # Start development environment
make docker-down-dev        # Stop development environment
make docker-rebuild-dev     # Rebuild and restart dev services
make docker-logs            # View all logs
make docker-logs-app        # View app logs only
```

### Production Commands
```bash
make docker-up-prod         # Start production environment
make docker-down-prod       # Stop production environment
make docker-logs-prod       # View production logs
```

### Database Commands
```bash
make docker-db-shell        # Open database shell
make db-reset              # Reset database (WARNING: deletes data)
```

### Utility Commands
```bash
make docker-status          # Show service status
make docker-shell          # Open shell in app container
make docker-clean          # Clean up Docker resources
```

## 🔧 Docker Compose Files

### Development (`docker-compose.dev.yml`)
- Includes development tools (pgAdmin, Redis Commander)
- Volume mounts for live code reloading
- Debug logging enabled
- Separate data volumes for development

### Production (`docker-compose.prod.yml`)
- Optimized for production use
- Resource limits configured
- No development tools
- Environment variables from external sources

### Default (`docker-compose.yml`)
- Standard configuration
- Suitable for testing and staging

## 🌐 Service URLs

After starting the services, you can access:

- **API**: http://localhost:8080
- **API Health**: http://localhost:8080/v1/health
- **MinIO Console**: http://localhost:9001 (admin/minioadmin123)
- **pgAdmin**: http://localhost:5050 (admin@tutuplapak.com/admin123)
- **Redis Commander**: http://localhost:8081

## 🗄️ Database Access

### Using pgAdmin
1. Open http://localhost:5050
2. Login with `admin@tutuplapak.com` / `admin123`
3. Add server with:
   - Host: `postgres`
   - Port: `5432`
   - Username: `tutuplapak`
   - Password: `postgres`

### Using Command Line
```bash
# Open database shell
make docker-db-shell

# Or directly with docker-compose
docker-compose exec postgres psql -U tutuplapak -d tutuplapak_db
```

## 📁 File Storage

MinIO is configured with:
- Bucket: `tutuplapak-files`
- Access: Public read
- Console: http://localhost:9001

## 🔍 Troubleshooting

### Common Issues

1. **Port conflicts**
   ```bash
   # Check what's using the ports
   lsof -i :8080
   lsof -i :5433
   ```

2. **Database connection issues**
   ```bash
   # Check database logs
   docker-compose logs postgres
   
   # Restart database
   docker-compose restart postgres
   ```

3. **MinIO bucket not created**
   ```bash
   # Check MinIO setup logs
   docker-compose logs minio-setup
   
   # Manually create bucket
   docker-compose exec minio-setup /usr/bin/mc mb myminio/tutuplapak-files
   ```

4. **Permission issues**
   ```bash
   # Fix volume permissions
   sudo chown -R $USER:$USER .
   ```

### Reset Everything
```bash
# Stop and remove all containers, networks, and volumes
make docker-clean

# Start fresh
make docker-up-dev
```

## 🚀 Production Deployment

### 1. Environment Setup
```bash
# Create production environment file
cp env.example .env.prod

# Update with production values
nano .env.prod
```

### 2. Deploy
```bash
# Start production services
make docker-up-prod

# Or with custom environment
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

### 3. Monitoring
```bash
# View production logs
make docker-logs-prod

# Check service health
docker-compose -f docker-compose.prod.yml ps
```

## 📊 Health Checks

All services include health checks:

- **API**: HTTP GET to `/v1/health`
- **PostgreSQL**: `pg_isready` command
- **MinIO**: HTTP GET to `/minio/health/live`
- **Redis**: `redis-cli ping`

## 🔒 Security Notes

- Change default passwords in production
- Use strong JWT secrets
- Configure proper CORS origins
- Enable SSL/TLS for MinIO in production
- Use secrets management for sensitive data

## 📈 Performance Tuning

### Resource Limits
Production compose file includes resource limits:
- API: 512MB RAM, 0.5 CPU
- PostgreSQL: 1GB RAM, 1 CPU

### Database Optimization
- Consider connection pooling
- Tune PostgreSQL settings
- Monitor query performance

### Caching
- Redis is available for caching
- Configure cache strategies in your application

## 🆘 Support

If you encounter issues:

1. Check the logs: `make docker-logs`
2. Verify service status: `make docker-status`
3. Check resource usage: `docker stats`
4. Review this documentation
5. Check Docker and Docker Compose versions

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [MinIO Docker Image](https://hub.docker.com/r/minio/minio)
