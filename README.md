# 🎫 Ticket Queen

A full-stack ticket purchasing web application built with React, Go, and MySQL. Features real-time inventory management, user authentication, and concurrency control to prevent overselling.

![Ticket Queen](https://img.shields.io/badge/React-18.2.0-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

### User Features

- 🔐 **User Authentication** - Secure login and registration system
- 🎭 **Event Browsing** - View all available events with real-time availability
- 🛒 **Shopping Cart** - Add multiple tickets to cart before checkout
- 📊 **User Dashboard** - View purchase history and ticket status
- ✅ **Ticket Management** - Mark tickets as used/unused
- 🔄 **Real-time Updates** - Automatic inventory updates after purchases

### Technical Features

- 🔒 **Concurrency Control** - Mutex locks prevent overselling
- 💾 **Database Transactions** - ACID compliance with row-level locking
- 🚀 **RESTful API** - Clean API design with proper HTTP methods
- 📱 **Responsive Design** - Works seamlessly on mobile and desktop
- ⚡ **Performance Optimized** - Efficient database queries with indexes
- 🛡️ **Error Handling** - Graceful error handling throughout the stack

## 🛠️ Tech Stack

### Frontend

- **React** 18.2.0 - UI framework

### Backend

- **Go** 1.21+ - Server-side language
- **Gorilla Mux** - HTTP router
- **MySQL Driver** - Database connectivity

### Database

- **MySQL** 8.0+ - Relational database
- **ACID Transactions** - Data integrity
- **Row-level Locking** - Concurrency control

---

## 🏗️ Architecture

```
┌─────────────────┐
│  React Frontend │ (Port 3000)
│   - UI/UX       │
│   - State Mgmt  │
└────────┬────────┘
         │ HTTP/REST
         │
┌────────▼────────┐
│   Go Backend    │ (Port 8080)
│   - API Server  │
│   - Mutex Lock  │
└────────┬────────┘
         │ SQL
         │
┌────────▼────────┐
│  MySQL Database │ (Port 3306)
│   - Events      │
│   - Users       │
│   - Purchases   │
└─────────────────┘
```

---

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js** (v16 or higher)

- **Go** (v1.21 or higher)

- **MySQL** (v8.0 or higher)

- **Git**

## 🔐 Concurrency Control

### Problem: Overselling

When multiple users try to purchase the last few tickets simultaneously, without proper concurrency control, the system might sell more tickets than available.

### Solution: Two-Layer Protection

#### 1. Application-Level Mutex Lock (Go)

```go
s.mutex.Lock()
defer s.mutex.Unlock()
```

Ensures only one purchase request is processed at a time.

#### 2. Database-Level Row Locking (MySQL)

```sql
SELECT available FROM events WHERE id = ? FOR UPDATE
```

Locks the specific event row during the transaction.

### Two Layers?

Since the database supports transactions with row-level locking with SELECT ... FOR UPDATE locks the selected row for the duration of the transaction.
Two concurrent purchases for the same event will not cause overbooking. The Go mutex does not add any extra protection for database rows, it only prevents concurrent Go requests, which is unnecessary for correctness. I decided to still keep it just for a demo of use of mutex.

### Concurrency Control Flow

```
User Request → Mutex Lock → Begin Transaction →
Check Availability → Update Inventory →
Record Purchase → Commit Transaction → Release Lock
```

---

## 🚀 Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/yanranw1/tickets.git
   cd tickets
   ```

2. **Set up the database**

   ```bash
   cd database
   mysql -u root -p < database/schema.sql
   ```

3. **Start the backend**

   ```bash
    cd backend
    # Initialize Go module
    go mod init ticket-queen

    # Install dependencies
    go get github.com/go-sql-driver/mysql
    go get github.com/gorilla/mux
    # run
    go run .
   ```

4. **Start the frontend**
   ```bash
   cd frontend
   npm install
   npm start
   ```

![Screenshot](images/screenshot1.png)
![Screenshot](images/screenshot2.png)
