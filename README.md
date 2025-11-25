# WebSocket Chat Application

A real-time chat application built with Go, WebSockets, PostgreSQL, and HTMX.

## Features

- User authentication with JWT and bcrypt
- Real-time messaging via WebSockets
- Multiple chat rooms
- Persistent message history
- Responsive UI with dark/light mode

## Tech Stack

- **Backend:** Go, Gin, Gorilla WebSocket, PostgreSQL, JWT, bcrypt
- **Frontend:** HTMX, CSS

## Prerequisites

- Go 1.24.5 or higher
- PostgreSQL database
- Git

## Installation

1. **Clone the repository**
```bash
git clone <your-repo-url>
cd websockets
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up PostgreSQL database**

Create a PostgreSQL database and run the following schema:

```sql
CREATE TABLE users (
    username VARCHAR(255) UNIQUE NOT NULL,
    id VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE rooms (
    id VARCHAR(255) PRIMARY KEY
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    roomid VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (roomid) REFERENCES rooms(id)
);

CREATE INDEX idx_messages_roomid ON messages(roomid);
CREATE INDEX idx_messages_created_at ON messages(created_at);
```

4. **Configure environment variables**

Create a `.env` file in the project root:

```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_SSLMODE=disable
JWT_SECRET=your_secret_key_here
```

## Running the Application

1. **Start the server**
```bash
go run main.go
```

2. **Access the application**
- Open your browser and navigate to: `http://localhost:8080`
- Register a new account at: `http://localhost:8080/register`
- Login at: `http://localhost:8080/login`
