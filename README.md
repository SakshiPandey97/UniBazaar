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
VITE_USER_BASE_URL=\
VITE_PRODUCT_BASE_URL=\
VITE_CHAT_USERS_BASE_URL=

### ⚙️ Backend/products/.env
MONGO_URI=\
AWS_REGION=\
AWS_ACCESS_KEY_ID=\
AWS_SECRET_ACCESS_KEY=\
AWS_S3_BUCKET=\
AWS_CONSOLE=\
AWS_USER=\
AWS_PWD=

### ⚙️ Backend/messaging/.env
TODO

### ⚙️ Backend/users/.env
SENDGRID_API_KEY=

## 🛠️ Running Locally

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

### 4. Run React frontend
```bash
cd FE-UniBazaar
pnpm install
pnpm run dev
```
The app will run at http://localhost:3000

## 🛠️ Local Database Setup

### Mongo DB (Products Service)
#### 1. Install MongoDB Compass
   Download from [here](https://www.mongodb.com/try/download/compass).

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

## ☁️ AWS S3 Setup (for Image Uploads)

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
AWS_CONSOLE=https://console.aws.amazon.com/
AWS_USER=<your-iam-username>
AWS_PWD=<your-iam-user-password>
```


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

| Method | Endpoint           | Description      |
| ------ | ------------------ | ---------------- |
| POST   | `/signup`          | Register a user  |
| POST   | `/verifyEmail`     | OTP verification |
| POST   | `/resendOtp`       | Resend OTP       |
| POST   | `/login`           | Login            |
| POST   | `/logout`          | Logout           |
| POST   | `/getjwt`          | Get JWT token    |
| GET    | `/verifyjwt`       | Verify JWT token |
| GET    | `/displayUser/:id` | Get user details |

--- 

### Products Service

| Method | Endpoint                         | Description          |
| ------ | -------------------------------- | -------------------- |
| POST   | `/products`                      | Create new product   |
| GET    | `/products`                      | Get all products     |
| GET    | `/products/{UserId}`             | Get products by user |
| PUT    | `/products/{UserId}/{ProductId}` | Update product       |
| DELETE | `/products/{UserId}/{ProductId}` | Delete product       |
| GET    | `/search/products`               | Search products      |

--- 

### Messaging Service

| Method | Endpoint                                | Description                    |
| ------ | --------------------------------------- | ------------------------------ |
| GET    | `/api/conversation/{user1ID}/{user2ID}` | Get conversation               |
| POST   | `/messages`                             | Send a message                 |
| GET    | `/users`                                | Get users                      |
| POST   | `/api/users/sync`                       | Sync user data                 |
| GET    | `/api/unread-senders`                   | Get users with unread messages |
| WS     | `/ws`                                   | WebSocket for real-time chat   |

---

## ☁️ Cloud Deployment

UniBazaar is fully deployed in the cloud, so you can use it without any local setup.

### 🌐 Live App
### 📱 Frontend (React on Vercel)
#### 🔗 [https://unibazaar.vercel.app](https://unibazaar.vercel.app)

### 💻 Backend Services (Hosted on Azure)
| Service               | URL                                                                                          |
| --------------------- | -------------------------------------------------------------------------------------------- |
| **Users Service**     | [https://unibazaar-users.azurewebsites.net](https://unibazaar-users.azurewebsites.net)       |
| **Products Service**  | [https://unibazaar-products.azurewebsites.net](https://unibazaar-products.azurewebsites.net) |
| **Messaging Service** | `ws://messaging.eastus.cloudapp.azure.com:8000` (WebSocket URL)                              |

Make sure your frontend `.env` file points to these URLs to connect to the production services.

### 🗄️ Databases
| Component        | Provider                                                        |
| ---------------- | --------------------------------------------------------------- |
| **Users DB**     | PostgreSQL on [Neon Tech](https://neon.tech)                    |
| **Products DB**  | MongoDB on [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) |
| **Messaging DB** | PostgreSQL on [Neon Tech](https://neon.tech)                    |

---

This cloud setup enables seamless access across all platforms without requiring local infrastructure. Great for demos, testing, or getting straight to buying and selling on UniBazaar!

---

## 💡 Why UniBazaar?
### 🌱 Sustainability First
*Encouraging reuse and reducing campus waste.*

### 🔗 Community Driven  
*Built to connect students through trust and purpose.*

---

### 🎯 Our Vision  
*To make the circular economy a campus culture.*

---

### 📬 Contributing  
*We welcome contributions! Please fork the repo and open a PR with clear commits and comments.*

---

## Members

- **Tanmay Saxena**
- **Shubham Singh**
- **Avaneesh Khandekar**
- **Sakshi Pandey**
--- 
