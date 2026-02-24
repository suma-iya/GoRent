# Docker Setup Guide

## Prerequisites
- Docker Desktop installed and running
- Docker Compose (usually included with Docker Desktop)

## Quick Start

### 1. Start MySQL and Backend with Docker Compose

```bash
cd /Users/suma/Documents/development/project/rentApp
docker-compose up -d
```

This will:
- Start MySQL database container
- Build and start the Go backend container
- Create the `rent` database automatically
- Set up the `suma` user with the correct password

### 2. Check if services are running

```bash
docker-compose ps
```

You should see both `rent-mysql` and `rent-backend` running.

### 3. View backend logs

```bash
docker-compose logs -f backend
```

You should see:
```
Successfully connected to the database!
Server starting on http://192.168.0.230:8080
```

### 4. Test the backend

```bash
curl http://localhost:8080/test/fcm-connection-public
```

## Useful Commands

### Stop services
```bash
docker-compose down
```

### Stop and remove volumes (clean slate)
```bash
docker-compose down -v
```

### Rebuild and restart
```bash
docker-compose up -d --build
```

### View MySQL logs
```bash
docker-compose logs -f mysql
```

### Access MySQL directly
```bash
docker-compose exec mysql mysql -u suma -p
# Password: tMyc6mApj]wgzHl7
```

### Access backend container shell
```bash
docker-compose exec backend sh
```

## Database Initialization

If you need to run SQL scripts to initialize the database schema:

```bash
# Copy SQL files into MySQL container and execute
docker-compose exec -T mysql mysql -u suma -p'tMyc6mApj]wgzHl7' rent < your_script.sql
```

Or run multiple scripts:
```bash
for file in *.sql; do
  docker-compose exec -T mysql mysql -u suma -p'tMyc6mApj]wgzHl7' rent < "$file"
done
```

## Troubleshooting

### Backend can't connect to database
- Check if MySQL is healthy: `docker-compose ps`
- Check MySQL logs: `docker-compose logs mysql`
- Wait a few seconds for MySQL to fully start

### Port already in use
If port 8080 or 3306 is already in use:
```bash
# Find what's using the port
lsof -i :8080
lsof -i :3306

# Or change ports in docker-compose.yml
```

### Rebuild after code changes
```bash
docker-compose up -d --build backend
```

## Environment Variables

The docker-compose.yml sets these environment variables for the backend:
- `DB_HOST=mysql` (service name in Docker network)
- `DB_PORT=3306`
- `DB_USER=suma`
- `DB_PASSWORD=tMyc6mApj]wgzHl7`
- `DB_NAME=rent`

You can override these in docker-compose.yml or use a `.env` file.


