# Using XAMPP Database Instead of Docker MySQL

## Problem
The backend is currently using Docker MySQL which is a fresh/empty database. Your previous database with all tables (property, floor, notification, etc.) is in XAMPP MySQL.

## Solution Options

### Option 1: Point Backend to XAMPP MySQL (Recommended)

Update `docker-compose.yml` to use XAMPP MySQL instead of Docker MySQL:

```yaml
services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: rent-backend
    ports:
      - "8081:8080"
    environment:
      - DB_HOST=host.docker.internal  # This allows Docker to access host machine
      - DB_PORT=3306                   # XAMPP MySQL port
      - DB_USER=suma
      - DB_PASSWORD=tMyc6mApj]wgzHl7
      - DB_NAME=rent
    networks:
      - rent-network
    restart: unless-stopped
```

Then restart:
```bash
docker-compose restart backend
```

### Option 2: Export XAMPP Database and Import to Docker

1. **Export from XAMPP:**
   ```bash
   mysqldump -usuma -ptMyc6mApj]wgzHl7 -h127.0.0.1 -P3306 rent > xampp_database.sql
   ```

2. **Import to Docker MySQL:**
   ```bash
   docker exec -i rent-mysql mysql -usuma -ptMyc6mApj]wgzHl7 rent < xampp_database.sql
   ```

### Option 3: Create All Missing Tables in Docker

If you want to keep using Docker MySQL, we need to create all the missing tables:
- property
- floor
- notification
- payment
- takes_care_of
- advance
- etc.

## Quick Fix: Use XAMPP Database

The easiest solution is to update docker-compose.yml to use XAMPP MySQL:

1. **Edit docker-compose.yml** - Change DB_HOST to `host.docker.internal`
2. **Restart backend**: `docker-compose restart backend`
3. **Verify**: Check backend logs to confirm connection


