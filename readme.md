# 🛍️ ECOMM-GO — E-Commerce Backend API

A production-ready e-commerce REST API built with **Go**, **Gin Framework**, and **MongoDB**. Features complete user authentication, product management, shopping cart, and order processing capabilities.

---

## 🎯 Features

- 🔐 **JWT Authentication** — Secure user signup/login with token-based auth
- 👤 **User Management** — Profile management and address handling
- 📦 **Product Catalog** — Browse and search products
- 🛒 **Shopping Cart** — Add, update, and remove cart items
- 💳 **Order System** — Checkout and order history
- 🗄️ **MongoDB** — Persistent data storage with efficient querying
- 🐳 **Docker Ready** — Containerized deployment setup

---

## 🏗️ Project Structure

```
ECOMM-GO/
│
├── controllers/          # Business logic handlers
│   ├── address.go       # Address management
│   ├── cart.go          # Shopping cart operations
│   └── controllers.go   # User & product controllers
│
├── database/            # Database configuration
│   ├── cart.go         # Cart database operations
│   └── databasesetup.go # MongoDB connection setup
│
├── middlewares/         # HTTP middlewares
│   └── middleware.go   # JWT auth middleware
│
├── models/             # Data models & schemas
│   └── models.go       # User, Product, Order models
│
├── routes/             # API route definitions
│   └── routes.go       # All route mappings
│
├── tokens/             # JWT token management
│   └── tokengen.go     # Token generation & validation
│
├── main.go             # Application entry point
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── .env                # Environment variables
├── .gitignore          # Git ignore rules
├── docker-compose.yaml # Docker configuration
└── readme.md           # This file
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22 or higher
- MongoDB (local or Atlas)
- Docker (optional)

### Installation Steps

**1. Clone the repository**

```bash
git clone <your-repo-url>
cd ECOMM-GO
```

**2. Install dependencies**

```bash
go mod download
go mod tidy
```

**3. Configure environment**

Create/update your `.env` file:

```env
PORT=8000
MONGODB_URL=mongodb://localhost:27017/ecommerce
SECRET_KEY=your-jwt-secret-key-here
```

**4. Run the application**

**Option A: Direct Go Run**

```bash
go run main.go
```

**Option B: Using Docker Compose**

```bash
docker-compose up -d
```

Server starts at: `http://localhost:8000`

---

## 🔌 API Endpoints

### Authentication

| Method | Endpoint        | Description           | Auth Required |
| ------ | --------------- | --------------------- | ------------- |
| POST   | `/users/signup` | Register new user     | No            |
| POST   | `/users/login`  | Login & get JWT token | No            |

**Signup Example:**

```json
POST /users/signup
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "securepass123",
  "phone": "+1234567890"
}
```

**Login Example:**

```json
POST /users/login
{
  "email": "john@example.com",
  "password": "securepass123"
}

Response:
{
  "token": "eyJhbGciOiJIUz...",
  "refresh_token": "eyJhbGc...",
  "user_id": "..."
}
```

---

### Products

| Method | Endpoint                    | Description              | Auth Required |
| ------ | --------------------------- | ------------------------ | ------------- |
| GET    | `/users/productview`        | Get all products         | Yes           |
| GET    | `/users/search?name=laptop` | Search products by query | Yes           |

---

### Shopping Cart

| Method | Endpoint                              | Description            | Auth Required |
| ------ | ------------------------------------- | ---------------------- | ------------- |
| POST   | `/addtocart?id=<user>&pid=<product>`  | Add item to cart       | Yes           |
| GET    | `/removeitem?id=<user>&pid=<product>` | Remove item from cart  | Yes           |
| GET    | `/listcart?id=<user>`                 | View user's cart items | Yes           |

---

### Orders & Checkout

| Method | Endpoint                                  | Description                  | Auth Required |
| ------ | ----------------------------------------- | ---------------------------- | ------------- |
| POST   | `/addaddress?id=<user>`                   | Add delivery address         | Yes           |
| PUT    | `/edithomeaddress?id=<user>`              | Edit home address            | Yes           |
| PUT    | `/editworkaddress?id=<user>`              | Edit work address            | Yes           |
| DELETE | `/deleteaddresses?id=<user>`              | Delete user addresses        | Yes           |
| GET    | `/cartcheckout?id=<user>`                 | Checkout all cart items      | Yes           |
| GET    | `/instantbuy?userid=<user>&pid=<product>` | Buy single product instantly | Yes           |

---

## 🗂️ Database Models

### User Schema

```go
type User struct {
    ID              primitive.ObjectID `bson:"_id"`
    First_Name      *string           `json:"first_name" validate:"required,min=2,max=30"`
    Last_Name       *string           `json:"last_name"`
    Email           *string           `json:"email" validate:"email,required"`
    Password        *string           `json:"password" validate:"required,min=6"`
    Phone           *string           `json:"phone" validate:"required"`
    Token           *string           `json:"token"`
    Refresh_Token   *string           `json:"refresh_token"`
    Created_At      time.Time         `json:"created_at"`
    Updated_At      time.Time         `json:"updated_at"`
    User_ID         string            `json:"user_id"`
    UserCart        []ProductUser     `json:"usercart" bson:"usercart"`
    Address_Details []Address         `json:"address" bson:"address"`
    Order_Status    []Order           `json:"orders" bson:"orders"`
}
```

### Product Schema

```go
type Product struct {
    Product_ID   primitive.ObjectID `bson:"_id"`
    Product_Name *string            `json:"product_name"`
    Price        *uint64            `json:"price"`
    Rating       *uint              `json:"rating"`
    Image        *string            `json:"image"`
}
```

---

## 🔒 Authentication Flow

1. User signs up → Password is hashed & stored
2. User logs in → JWT token generated
3. Token sent in request headers: `Authorization: Bearer <token>`
4. Middleware validates token for protected routes
5. Token expires → Use refresh token to get new access token

---

## 🐳 Docker Deployment

Your `docker-compose.yaml`:

```yaml
version: "3.8"

services:
  mongodb:
    image: mongo:latest
    container_name: ecomm-mongo
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    environment:
      MONGO_INITDB_DATABASE: ecommerce

  app:
    build: .
    container_name: ecomm-app
    ports:
      - "8000:8000"
    depends_on:
      - mongodb
    environment:
      MONGODB_URL: mongodb://mongodb:27017/ecommerce
      PORT: 8000
      SECRET_KEY: your-secret-key
    volumes:
      - .:/app

volumes:
  mongo_data:
```

**Dockerfile:**

```dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]
```

**Build and run:**

```bash
docker-compose up --build
```

---

## 🧪 Testing

**Test a signup:**

```bash
curl -X POST http://localhost:8000/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "password123",
    "phone": "+1234567890"
  }'
```

**Test a login:**

```bash
curl -X POST http://localhost:8000/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

**View products (with token):**

```bash
curl -X GET http://localhost:8000/users/productview \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 📦 Dependencies

```go
require (
    github.com/gin-gonic/gin
    github.com/golang-jwt/jwt/v4
    go.mongodb.org/mongo-driver
    github.com/joho/godotenv
    golang.org/x/crypto
)
```

Install all:

```bash
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v4
go get go.mongodb.org/mongo-driver/mongo
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
```

---

## ⚙️ Environment Variables

| Variable      | Description               | Example                               |
| ------------- | ------------------------- | ------------------------------------- |
| `PORT`        | Server port               | `8000`                                |
| `MONGODB_URL` | MongoDB connection string | `mongodb://localhost:27017/ecommerce` |
| `SECRET_KEY`  | JWT signing secret        | `my-super-secret-key`                 |

---

## 🛡️ Security Features

- ✅ Password hashing using bcrypt
- ✅ JWT token-based authentication
- ✅ Token expiration and refresh mechanism
- ✅ Protected routes with middleware
- ✅ Input validation
- ✅ CORS support
- ✅ Environment variable protection

---

## 🐛 Troubleshooting

**MongoDB connection failed:**

```
Error: connection refused
Solution: Ensure MongoDB is running on port 27017
```

**JWT token invalid:**

```
Error: token verification failed
Solution: Check SECRET_KEY in .env matches token generation
```

**Port already in use:**

```
Error: bind: address already in use
Solution: Change PORT in .env or kill process using: lsof -ti:8000 | xargs kill
```

---

## 📝 TODO / Future Enhancements

- [ ] Add product categories and filters
- [ ] Implement payment gateway integration
- [ ] Add order tracking system
- [ ] Email notifications
- [ ] Admin dashboard
- [ ] Product reviews and ratings
- [ ] Wishlist functionality
- [ ] Inventory management
- [ ] Analytics and reporting

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License.

---

## 👨‍💻 Author

**Your Name**  
GitHub: [@yourusername](https://github.com/yourusername)  
Email: your.email@example.com

Built with ❤️ using Go, Gin, and MongoDB

---

## 📚 Resources

- [Gin Documentation](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [JWT Best Practices](https://jwt.io/introduction)
- [Go by Example](https://gobyexample.com/)

---

**⭐ If you find this project helpful, please give it a star!**
