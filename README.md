# 📈 Stock Ratings API

A simple backend built in **Go** to manage and serve stock rating data.

## 🚀 Features

- Built with Go (Golang)
- RESTful endpoints to fetch and sync ratings
- Connects to a PostgreSQL database
- Provides the best stock recommendation
- Includes simple pagination and sorting

## 📦 Requirements

- [Go](https://golang.org/dl/) 1.20 or higher installed
- [PostgreSQL](https://www.postgresql.org/) database
- `uuid-ossp` PostgreSQL extension enabled
- Internet connection to sync data (if used)

## ⚙️ Environment Configuration

Create a `.env` file in the project root directory with the following variables:

```env
DB_USER=your_postgres_username
DB_PASSWORD=your_postgres_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_database_name
DB_SSL=disable
```
