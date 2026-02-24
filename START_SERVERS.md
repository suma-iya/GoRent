# How to Start Database and Backend Server

## Option 1: If MySQL is running in Docker

1. **Start MySQL Docker container:**
   ```bash
   docker run -d \
     --name mysql-rent \
     -e MYSQL_ROOT_PASSWORD=your_root_password \
     -e MYSQL_DATABASE=rent \
     -e MYSQL_USER=suma \
     -e MYSQL_PASSWORD=tMyc6mApj]wgzHl7 \
     -p 3306:3306 \
     mysql:8.0
   ```

2. **Or if you already have a MySQL container:**
   ```bash
   docker start mysql-rent
   ```

## Option 2: If MySQL is installed locally

1. **Start MySQL service:**
   ```bash
   # macOS (using Homebrew)
   brew services start mysql
   
   # Or manually
   mysql.server start
   ```

2. **Verify MySQL is running:**
   ```bash
   mysql -u suma -p -h localhost -P 3306
   # Password: tMyc6mApj]wgzHl7
   ```

## Start the Go Backend Server

1. **Navigate to project directory:**
   ```bash
   cd /Users/suma/Documents/development/project/rentApp
   ```

2. **Start the backend:**
   ```bash
   go run main.go
   ```

   Or if you have a compiled binary:
   ```bash
   ./main
   ```

3. **Verify it's running:**
   - You should see: "Successfully connected to the database!"
   - Server should start on: `http://192.168.0.230:8080`
   - Test endpoint: `http://localhost:8080/test/fcm-connection-public`

## Troubleshooting

### Database Connection Error
If you see "Failed to initialize database", check:
- MySQL is running: `docker ps` or `brew services list`
- Database credentials in `config/database.go` are correct
- Database `rent` exists: `mysql -u suma -p -e "SHOW DATABASES;"`
- Port 3306 is accessible

### Port Already in Use
If port 8080 is already in use:
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Update Database Host
If MySQL is running locally (not in Docker), update `config/database.go`:
```go
DBHost = "localhost"  // Change from "host.docker.internal"
```



