# UniBazaar

## Project Description

**UniBazaar** is a hyper-localized marketplace designed exclusively for university students. It allows users to buy, sell, and trade second-hand items such as textbooks, furniture, electronics, and more within their campus communities. By fostering a trusted, university-specific platform, StudentExchange aims to make affordable goods more accessible while promoting sustainability and reducing waste.

## 🚀 Features

UniBazaar comes packed with features that make it easy to buy, sell, and communicate on your campus.

The React app has the following functionality:

1. **Sign up and verify account with OTP** - Easily create an account and verify it through a One-Time Password (OTP).
2. **Login** - Secure login to access the marketplace.
3. **View user details** - View your profile information.
4. **View products available in the marketplace** - Browse products listed by other users.
5. **Search for products available in the marketplace** - Find the products you're looking for with search functionality.
6. **Message the seller for a product** - Negotiate price, arrange meetings, and discuss details with the seller.
7. **Post a new product that the user wants to sell** - List products you want to sell in the marketplace.
8. **View products posted by the user** - View and manage the products you’ve listed for sale.
9. **Edit/delete products posted by the user** - Update or remove products from your listing.
10. **Logout** - Securely log out from the application when you're done.

---

### 📱 React Frontend

- Built with Vite + pnpm
- User authentication (Sign up, OTP verification, Login/Logout)
- Product listing and search
- Messaging between buyers and sellers
- Product posting, editing, and deletion

### 💻 Backend Services

- **Users Service (Go)**: Manages user authentication, profiles, and sessions.
- **Products Service (Go)**: Handles CRUD for products and search functionality.
- **Messaging Service (Go)**: Real-time messaging and conversations.

---

## ☁️ Cloud Deployment

UniBazaar is fully deployed in the cloud, so you can use it without any local setup.

### 🌐 Live App

### 📱 Frontend (React on Vercel)

#### 🔗 [https://unibazaar.vercel.app](https://unibazaar.vercel.app)

If you don’t want to set up everything locally, just visit the link above to explore the full UniBazaar experience online.

### 💻 Backend Services (Hosted on Azure)

| Service               | URL                                                                                            |
| --------------------- | ---------------------------------------------------------------------------------------------- |
| **Users Service**     | [https://unibazaar-users.azurewebsites.net](https://unibazaar-users.azurewebsites.net)         |
| **Products Service**  | [https://unibazaar-products.azurewebsites.net](https://unibazaar-products.azurewebsites.net)   |
| **Messaging Service** | [wss://unibazaar-messaging.azurewebsites.net](WebSocket url)                                   |
| **Messaging Service** | [https://unibazaar-messaging.azurewebsites.net](https://unibazaar-messaging.azurewebsites.net) |

Make sure your frontend `.env` file points to these URLs to connect to the production services.

### 🗄️ Databases

| Component        | Provider                                                        |
| ---------------- | --------------------------------------------------------------- |
| **Users DB**     | PostgreSQL on [Neon Tech](https://neon.tech)                    |
| **Products DB**  | MongoDB on [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) |
| **Messaging DB** | PostgreSQL on [AWS RDS](https://aws.amazon.com/free/database/)  |

---

This cloud setup enables seamless access across all platforms without requiring local infrastructure. Great for demos, testing, or getting straight to buying and selling on UniBazaar!

---

## 🧪 Requirements

### 🔧 Prerequisites

- Node.js (v18+)
- Go (v1.20+)
- pnpm (v8+)
- MogoDB (for local database setup)
- PostgreSQL (for local database setup)

## ⚙️ Environment Setup

Each service and the client app requires a `.env` file.

### ⚙️ FE-UniBazaar/.env

```env
VITE_USER_BASE_URL="https://unibazaar-users.azurewebsites.net"
VITE_PRODUCT_BASE_URL="https://unibazaar-products.azurewebsites.net"
VITE_CHAT_USERS_BASE_URL="https://unibazaar-messaging.azurewebsites.net"
VITE_CHAT_USERS_WS_URL="wss://unibazaar-messaging.azurewebsites.net"
```

### ⚙️ Backend/products/.env

```env
MONGO_URI=<MONGO_DB_CONNECTION_STRING>
AWS_REGION=<AWS_REGION>
AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>
AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY>
AWS_S3_BUCKET=<AWS_S3_BUCKET_NAME>
AWS_CONSOLE=<AWS_CONSOLE_URL>
AWS_USER=<AWS_USER_ID>
AWS_PWD=<AWS_USER_ID_PASSWORD>
```

### ⚙️ Backend/messaging/.env

```env
CHAT_DB_URI=<YOUR_POSTGRESQL_CONNECTION_STRING>
```

### ⚙️ Backend/users/.env

```env
DSN=<POSTGRES_DB_CONNECTION_STRING>
JWT_SECRET=<JWT_SECRET>
SENDGRID_API_KEY=<API_KEY>
```

# 🛠️ Running Locally

## 🗃️ Local Database Setup

### Mongo DB (Products Service)

#### 1. Install MongoDB Compass

- Download from [here](https://www.mongodb.com/try/download/compass).

#### 2. Create DB and Collection

- Open Compass and connect to `mongodb://localhost:27017`.
- Create a database named `unibazaar`.
- Inside it, create a collection named `products`.

#### 3. Set Environment Variable

- Copy the connection string (e.g., `mongodb://localhost:27017/unibazaar`).
- In `Backend/Products/.env`, add (e.g.):
  ```env
  MONGO_URI=mongodb://localhost:27017/unibazaar
  ```

### PostgreSQL (Users Service)

#### 1. Install PostgreSQL

- Download and install from the [official PostgreSQL website](https://www.postgresql.org/download/).
- During setup, remember your **username (`postgres`)** and **password**.

#### 2. Create Database and Table

- Open **pgAdmin** or use the terminal:

  ```bash
  psql -U postgres
  ```

- Create the database:

  ```sql
  CREATE DATABASE unibazaar;
  ```

- Connect to it:

  ```bash
  \c unibazaar
  ```

- Create the `users` table:
  ```sql
  CREATE TABLE users (
      userid BIGSERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      email TEXT UNIQUE NOT NULL,
      password TEXT NOT NULL,
      phone TEXT,
      verified BOOLEAN NOT NULL DEFAULT false,
      otp_code TEXT,
      failed_reset_attempts INT NOT NULL DEFAULT 0
  );
  ```

#### 3. Set Environment Variable

- In `Backend/Users/.env`, add:
  ```env
  DSN=postgres://postgres:admin@localhost/unibazaar?sslmode=disable
  (This is an example change the connection string to match your local setup.)
  ```

### PostgreSQL (Messaging Service)

#### 1. Install PostgreSQL (if not already done)

- Follow the steps in the "PostgreSQL (Users Service)" section if you haven't installed PostgreSQL yet.

#### 2. Create Database and User

- Open `psql` (PostgreSQL command-line tool) or use a GUI tool like pgAdmin.
- Connect as the `postgres` user (or your superuser).
- Create a dedicated user and database for the messaging service:

  ```sql
  CREATE USER unibazaar_msg_user WITH PASSWORD 'your_messaging_db_password';

  CREATE DATABASE unibazaar_messaging OWNER unibazaar_msg_user;
  ```

#### 3. Create Database Tables

- Connect to the newly created `unibazaar_messaging` database.
- You'll need to create two main tables: one to store user information relevant for messaging, and another to store the chat messages themselves.

  **`users` Table:**
  This table holds basic user details needed by the messaging service. It is expected to be synchronized automatically when new users login.

  - `id`: A unique identifier for the user (Primary Key).
  - `name`: The user's display name.
  - `email`: The user's unique email address.

  **`messages` Table:**
  This table stores the individual chat messages exchanged between users.

  - `id`: A unique identifier for the message (Primary Key, like a UUID string).
  - `sender_id`: The ID of the user who sent the message (Foreign Key referencing `users.id`).
  - `receiver_id`: The ID of the user who received the message (Foreign Key referencing `users.id`).
  - `content`: The text content of the message.
  - `timestamp`: When the message was sent (stored as a Unix timestamp).
  - `read`: A boolean flag indicating if the message has been read by the receiver.
  - `sender_name`: The name of the sender.

  **Example `CREATE TABLE` Statements:**

  ```sql
  CREATE TABLE users (
      id    BIGSERIAL PRIMARY KEY,
      name  VARCHAR(255) NOT NULL,
      email VARCHAR(255) UNIQUE NOT NULL
  );

  CREATE TABLE messages (
      id          VARCHAR(255) PRIMARY KEY,
      sender_id   BIGINT NOT NULL REFERENCES users(id),
      receiver_id BIGINT NOT NULL REFERENCES users(id),
      content     TEXT NOT NULL,
      timestamp   BIGINT NOT NULL,
      read        BOOLEAN NOT NULL DEFAULT false,
      sender_name VARCHAR(255) NOT NULL
  );
  ```

- **Note:** Adjust column types (like `BIGINT`, `VARCHAR`, `TEXT`) and constraints based on your specific needs and PostgreSQL version.

#### 4. Set Environment Variable

- In `Backend/messaging/.env`, set the `CHAT_DB_URI` variable with the appropriate connection string:

  ```env
  # --- Choose ONE of the following examples ---

  # Example using the default 'postgres' user:
  # CHAT_DB_URI="host=localhost user=postgres password=your_postgres_password dbname=unibazaar_messaging port=5432 sslmode=disable"

  # Example using the dedicated 'unibazaar_msg_user':
  CHAT_DB_URI="host=localhost user=unibazaar_msg_user password=your_messaging_db_password dbname=unibazaar_messaging port=5432 sslmode=disable"

  # Example for AWS RDS (replace placeholders):
  # CHAT_DB_URI="host=your_rds_endpoint user=your_rds_user password=your_rds_password dbname=unibazaar_messaging port=5432 sslmode=require"
  ```

- Replace placeholders like `your_postgres_password` or `your_messaging_db_password` with the actual passwords you set up.

### ☁️ AWS S3 Setup (for Image Uploads)

To enable image uploads in the Products Service using AWS S3, follow these steps:

#### 1. Create an S3 Bucket

- Go to the [AWS S3 Console](https://s3.console.aws.amazon.com/s3).
- Click **Create bucket**, give it a name (e.g., `unibazaar`), and choose a region.

#### 2. Create an IAM User

- Go to the [IAM Console](https://console.aws.amazon.com/iam/).
- Create a new user with **Programmatic access**.
- Attach the **AmazonS3FullAccess** policy (or a custom policy with limited access to your bucket).

#### 3. Copy Credentials

- Note down the **Access Key ID** and **Secret Access Key** for the IAM user.

#### 4. Set Up Environment Variables

In `Backend/Products/.env`, add:

```env
AWS_REGION=<your-region>
AWS_ACCESS_KEY_ID=<your-access-key-id>
AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
AWS_S3_BUCKET=<your-s3-bucket-name>
AWS_CONSOLE=<aws-console-url>
AWS_USER=<your-iam-username>
AWS_PWD=<your-iam-user-password>
```

## 🔌 Local Code Setup

### 1. Clone the repo

```bash
git clone https://github.com/SakshiPandey97/UniBazaar.git
cd unibazaar
```

### 2. Set up environment files

Create a .env file in each folder as shown above.

### 3. Run Go backend services

```bash
cd Backend/users
go mod tidy
go run main.go
```

The user service will run at http://localhost:4000

```bash
cd Backend/products
go mod tidy
go run main.go
```

The products service will run at http://localhost:8080

```bash
cd Backend/messaging
go mod tidy
go run main.go
```

The messaging service will run at http://localhost:8000

Make sure to update `FE-UniBazaar/.env` as:

```
VITE_USER_BASE_URL="http://localhost:4000"
VITE_PRODUCT_BASE_URL="http://localhost:8080"
VITE_CHAT_USERS_BASE_URL="http://localhost:8000"
VITE_CHAT_USERS_WS_URL="ws://localhost:8000"
```

### 4. Run React frontend

```bash
cd FE-UniBazaar
pnpm install
pnpm run dev
```

The app will run at http://localhost:3000

## 📡 Frontend Routes

| Route           | Description                                               |
| --------------- | --------------------------------------------------------- |
| `/sell`         | Page to post a new product for sale                       |
| `/messaging`    | Messaging center to chat with other users                 |
| `/products`     | Explore all products listed by other users                |
| `/userproducts` | View, edit, or delete products listed by the current user |
| `/about`        | Learn about UniBazaar’s mission, vision, and values       |

---

## 📡 API Routes

### Users Service

| Method | Endpoint            | Description             |
| ------ | ------------------- | ----------------------- |
| POST   | `/signup`           | Register a user         |
| POST   | `/verifyEmail`      | OTP verification        |
| POST   | `/resendOtp`        | Resend OTP              |
| POST   | `/login`            | Login                   |
| POST   | `/logout`           | Logout                  |
| POST   | `/getjwt`           | Get JWT token           |
| GET    | `/verifyjwt`        | Verify JWT token        |
| GET    | `/displayUser/{id}` | Get user details by id  |
| POST   | `/forgotPassword`   | Forgot password handler |
| POST   | `/updatePassword`   | Update password handler |

---

### Products Service

| Method | Endpoint                                          | Description          |
| ------ | ------------------------------------------------- | -------------------- |
| POST   | `/products`                                       | Create new product   |
| GET    | `/products?lastId={lastId}&limit={limit}`         | Get all products     |
| GET    | `/producs/{userId}?lastId={lastId}&limit={limit}` | Get products by user |
| PUT    | `/products/{UserId}/{ProductId}`                  | Update product       |
| DELETE | `/products/{UserId}/{ProductId}`                  | Delete product       |
| GET    | `/search/products?query={query}&limit={limit}`    | Search products      |

---

### Messaging Service

| Method | Endpoint                                | Description                    |
| ------ | --------------------------------------- | ------------------------------ |
| GET    | `/api/conversation/{user1ID}/{user2ID}` | Get conversation               |
| POST   | `/messages`                             | Send a message                 |
| GET    | `/users`                                | Get users                      |
| POST   | `/api/users/sync`                       | Sync latest user data          |
| GET    | `/api/unread-senders`                   | Get users with unread messages |
| WSS    | `/wss`                                  | WebSocket for real-time chat   |

---

## 💡 Why UniBazaar?

### 🌱 Sustainability First

_Encouraging reuse and reducing campus waste._

### 🔗 Community Driven

_Built to connect students through trust and purpose._

---

### 🎯 Our Vision

_To make the circular economy a campus culture._

---

### 📬 Contributing

_We welcome contributions! Please fork the repo and open a PR with clear commits and comments._

---

## Members

- **Tanmay Saxena**
- **Shubham Singh**
- **Avaneesh Khandekar**
- **Sakshi Pandey**

---
