# MySQL Setup Instructions

## Step 1: Start MySQL

```bash
# Option 1: Using Homebrew
brew services start mysql

# Option 2: Manual start
mysql.server start
```

## Step 2: Connect as root and set up the user

```bash
# Connect as root (you'll be prompted for root password)
mysql -u root -p
```

If you don't know the root password, you can reset it:

### Reset MySQL root password (if needed):

1. **Stop MySQL:**
   ```bash
   brew services stop mysql
   # or
   mysql.server stop
   ```

2. **Start MySQL in safe mode:**
   ```bash
   mysqld_safe --skip-grant-tables &
   ```

3. **Connect without password:**
   ```bash
   mysql -u root
   ```

4. **Reset root password:**
   ```sql
   USE mysql;
   UPDATE user SET authentication_string=PASSWORD('your_new_root_password') WHERE User='root';
   FLUSH PRIVILEGES;
   EXIT;
   ```

5. **Restart MySQL normally:**
   ```bash
   brew services restart mysql
   ```

## Step 3: Create the 'suma' user and database

**Option A: Using the SQL script (recommended)**
```bash
mysql -u root -p < setup_mysql_user.sql
```

**Option B: Manual SQL commands**
```bash
mysql -u root -p
```

Then run:
```sql
CREATE DATABASE IF NOT EXISTS rent;
CREATE USER IF NOT EXISTS 'suma'@'localhost' IDENTIFIED BY 'tMyc6mApj]wgzHl7';
GRANT ALL PRIVILEGES ON rent.* TO 'suma'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

## Step 4: Verify the setup

```bash
mysql -u suma -p
# Password: tMyc6mApj]wgzHl7

# Once connected, verify database access:
SHOW DATABASES;
USE rent;
SHOW TABLES;
```

## Step 5: Start the Go backend

Once MySQL is running and the user is set up:
```bash
cd /Users/suma/Documents/development/project/rentApp
go run main.go
```

You should see: "Successfully connected to the database!"

## Troubleshooting

### If you can't remember root password:
1. Use the safe mode method above to reset it
2. Or use `sudo` if you have admin access

### If user already exists with wrong password:
```sql
ALTER USER 'suma'@'localhost' IDENTIFIED BY 'tMyc6mApj]wgzHl7';
FLUSH PRIVILEGES;
```

### If you want to use a different password:
1. Update `config/database.go` with your new password
2. Then run: `ALTER USER 'suma'@'localhost' IDENTIFIED BY 'your_new_password';`



