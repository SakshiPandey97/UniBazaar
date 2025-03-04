# Sprint 2: User Stories and Implementation Plan

## Overview
During **Sprint 2**, our primary focus was on enhancing user management, strengthening account security, improving backend performance, implementing messaging system in backend and refining product-related functionalities. Below is an overview of the issues we tackled, the tasks completed, and the corresponding backend API documentation for user operations.

---

## Sprint 2: Completed Issues

| Issue                                                                          | Status     | Type                    |
|--------------------------------------------------------------------------------|------------|-------------------------|
| Provide API Support for user to add phone number                               | ✅ Closed  | Backend, Sprint v2      |
| Provide API Support for user to update phone number                            | ✅ Closed  | Backend, Sprint v2      |
| Validate US phone numbers in user profiles                                     | ✅ Closed  | Backend, Sprint v2      |
| Detect and notify users of suspicious login attempts                           | ✅ Closed  | Backend, Sprint v2      |
| Email verification for new users                                              | ✅ Closed  | Backend, Sprint v2      |
| Modify Product struct to have user id as a foreign key                         | ✅ Closed  | Backend, Sprint v2      |
| Create API spec and add Postman Collection                                    | ✅ Closed  | Backend, Documentation   |
| Performance improvement for API calls                                         | ✅ Closed  | Backend, Sprint v2      |
| Fix strict validations reload boundaries                                       | ✅ Closed  | Backend, Sprint v2      |
| Update Product API (Delete image from S3)                                      | ✅ Closed  | Backend, Sprint v2      |
| Linking Frontend Register Page with Backend                                    | ✅ Closed  | Frontend, Sprint v2     |
| Backend Implementation for user updating product using product Id             | ✅ Closed  | Backend, Sprint v2      |
| Backend Implementation for user deleting product using product Id             | ✅ Closed  | Backend, Sprint v2      |
| User Registration and DB creation                                             | ✅ Closed  | Backend, Sprint v2      |
| Backend Implementation for Forgot Password Functionality                       | ✅ Closed  | Backend, Sprint v2      |
| S3 support from backend to upload and download images                          | ✅ Closed  | Backend, Sprint v2      |
| Implementation of messaging system for chat using websocket connection         | ✅ Closed  | Backend, Sprint v2      |
| Use Go routines to speed up requests by running mongo and s3 in parallel       | ✅ Closed  | Backend, Sprint v2      |

---

## Summary of Work Done in Sprint 2

1.  **Forgot Password Flow**
    - **OTP-Based Reset:** Implemented an OTP-based reset mechanism, allowing users to securely reset their passwords if they forget them.  

2. **Email Verification & Suspicious Login Detection**  
   - **Email Verification:** Users now receive a unique OTP via SendGrid when registering. This OTP is required to activate their accounts.  
   - **Suspicious Login Attempts:** Added logic to detect repeated failed login attempts and send security alerts to the user’s email.

3. **Performance & Validation Improvements**  
   - Optimized API calls and database queries, particularly for user lookups and product retrieval.

4. **Product Updates**  
   - **User ID as Foreign Key:** Modified the product schema to associate each product with the user who created it.  
   - **Delete Image from S3:** On product deletion or update, images are now removed from AWS S3 to maintain storage cleanliness.

5. **Phone Number Support & Validation**  
   - **Add & Update Phone Number:** Implemented new endpoints to add or update phone numbers for users.  
   - **Validation:** Ensured only valid US numbers are accepted.
   
6. **Frontend Unit Testing**
    - Developed and executed unit tests for frontend components to ensure reliability and maintainability.

7. **Real-Time Messaging with WebSockets:** 
    - Established WebSocket-based real-time communication between users, enabling instant message delivery.  

8. **Chat History Persistence:** 
    - Implemented database storage for messages, ensuring that conversations remain accessible even if users disconnect.  

9. **Message Retrieval API:** 
    - Developed REST endpoints to fetch conversation history between users, supporting both individual and group chats.  

10. **Read Receipts & Unread Messages:** 
    - Added tracking for message read status, allowing users to see if their messages have been viewed. 

11. **Connection Handling & Reconnection Support:** 
    - Ensured stable WebSocket connections with mechanisms for automatic reconnection if a user disconnects.  

---

# UniBazaar User API Documentation

This section of the document describes the UniBazaar backend API for user management. It covers the **design choices**, **request/response formats**, **error handling**, and **sample JSON** requests for each endpoint. 


---
# Users: Backend API Documentation

## Table of Contents
1. [Design & Implementation Highlights](#design--implementation-highlights)
2. [Endpoints Overview](#endpoints-overview)
3. [Detailed Endpoint Documentation](#detailed-endpoint-documentation)
   - [Sign Up (POST /signup)](#1-sign-up-post-signup)
   - [Verify Email (POST /verifyEmail)](#2-verify-email-post-verifyemail)
   - [Forgot Password (POST /forgotpassword)](#3-forgot-password-post-forgotpassword)
   - [Verify Reset Code (POST /verifyresetcode)](#4-verify-reset-code-post-verifyresetcode)
   - [Update Password (POST /updatepassword)](#5-update-password-post-updatepassword)
   - [Delete User (POST /deleteuser)](#6-delete-user-post-deleteuser)
   - [Display User (POST /displayuser)](#7-display-user-post-displayuser)
   - [Login (POST /login)](#8-login-post-login)
   - [Update Name (POST /updatename)](#9-update-name-post-updatename)
   - [Update Phone (POST /updatephone)](#10-update-phone-post-updatephone)
   - [Error Cases & Responses](#error-cases--responses)
4. [Appendix: Security & Design Choices](#security-design-choices)

---

## Design & Implementation Highlights

### 1. Password Hashing (Argon2id)
- We use **Argon2id** (via [alexedwards/argon2id](https://github.com/alexedwards/argon2id)) for secure password storage.
- Parameters:
  - **Memory**: 128 MB
  - **Iterations**: 4
  - **Parallelism**: Number of CPU cores
  - **SaltLength**: 16 bytes
  - **KeyLength**: 32 bytes
- **Why Argon2id?**  
  Argon2id is recommended by OWASP for modern password hashing. It is resistant to GPU-cracking attacks and provides configurable memory hardness.

### 2. Password Complexity
- Minimum **60 bits of entropy** enforced via [go-password-validator](https://github.com/wagslane/go-password-validator).
- If a password is too weak, the server returns an error indicating insufficient complexity.

### 3. Email Validation
- Only **.edu** domains from specific Florida universities are allowed to register (`ufl.edu`, `fsu.edu`, `ucf.edu`, etc.).
- If the domain is not recognized, the request is rejected.

### 4. Phone Validation
- Regex is used to ensure a valid US number with 10 digits, optionally prefixed by `+1` is being used.

### 5. One-Time Password (OTP) Verification for Email upon Registration
- Upon user registration, a 6-digit OTP is generated and sent to the user's registered email via SendGrid.
- The user must enter the correct OTP to complete the email verification process.
- If the OTP is incorrect or expired, verification will fail, and the user must request a new OTP.

### 6. One-Time Password (OTP) Verification for Forgot Password
- When a user initiates a password reset, a 6-digit OTP is generated and sent to their registered email.
- The user must enter the correct OTP to proceed with resetting their password.
- If the user enters an incorrect OTP three or more times, a security alert email is triggered, notifying them of suspicious activity.

### 7. Database (GORM)
- GORM is used to interact with the database.
- **User** model:
  ```go
  type User struct {
      UserID              int    `gorm:"column:userid;primaryKey" json:"userid"`
      Name                string `json:"name"`
      Email               string `json:"email"`
      Password            string `json:"-"`
      OTPCode             string `json:"-"`
      FailedResetAttempts int    `json:"-"`
      Verified            bool   `json:"-"`
      Phone               string `json:"phone"`
  }

## Endpoints Overview

Below is a quick reference to each endpoint:

| **Endpoint**       | **Method** | **Description**                                        |
|--------------------|-----------|--------------------------------------------------------|
| `/signup`         | POST      | Create a new user, send verification OTP email.       |
| `/verifyEmail`    | POST      | Verify user email with OTP code.                      |
| `/forgotPassword` | POST      | Initiate a password reset (send reset OTP).           |
| `/updatePassword` | POST      | Update user password (OTP-based).      |
| `/deleteUser`     | POST      | Remove user from the system.                          |
| `/displayUser`    | POST      | Retrieve user information.                            |
| `/login`          | POST      | Authenticate user with email & password.              |
| `/updateName`     | POST      | Update user's display name.                           |
| `/updatePhone`    | POST      | Update user's phone number.                           |

---

# Detailed Endpoint Documentation

## 1. Sign Up (POST `/signup`)
### Description
Creates a new user, stores a hashed password, and emails an OTP code for verification.

### Request Body (JSON)
```json
{
  "name": "Jinx Silco",
  "email": "getjinxed@ufl.edu",
  "password": "MonkeyBomb#5567",
  "phone": "+15551234567"
}
```

### Behavior
- Validates email domain (must be recognized .edu).
- Validates password strength (≥ 60 bits).
- Validates USA phone number format.
- Hashes password using Argon2id.
- Saves user to the database.
- Generates a 6-digit OTP code, saves it, and sends it via email (through SendGrid's API).

### Success Response (JSON)
```json
{
  "message": "User created successfully. Please check your email for the OTP code."
}
```

### Error Cases
- **400 Bad Request**: Invalid email, weak password, or incorrect phone format.
- **409 Conflict**: User already exists.
- **500 Internal Server Error**: Database or email-sending issues.

---

## 2. Verify Email (POST `/verifyEmail`)
### Description
Verifies the user's email with the OTP code sent during sign-up.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu",
  "code": "223466"
}
```

### Behavior
- Checks if OTP matches.
- If valid, marks user as Verified.
- If invalid, increments failure counter. After 3 failures, sends a security alert email.

### Success Response (JSON)
```json
{
  "message": "Email verified successfully!"
}
```

### Error Cases
- **400 Bad Request**: Invalid or expired OTP.
- **404 Not Found**: User does not exist.
- **500 Internal Server Error**: Database or email-sending issues.

---

## 3. Forgot Password (POST `/forgotPassword`)
### Description
Initiates a password reset by generating and emailing a reset OTP.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu"
}
```

### Behavior
- Ensures user exists.
- Resets `FailedResetAttempts` to 0.
- Generates a 6-digit OTP and emails it.
- Stores OTP in the database.

### Success Response (JSON)
```json
{
  "message": "A reset code has been sent to your email."
}
```

### Error Cases
- **404 Not Found**: User does not exist.
- **500 Internal Server Error**: Database or email issues.

---

## 4. Verify Reset Code (POST `/updatePassword`)
### Description
Verifies the OTP code sent for password reset, and if correct, sets the new password.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu",
  "otp_code": "654321",
  "new_password": "EkkoRew1nd!"
}
```

### Behavior
- Checks if OTP matches.
- Validates new password strength.
- Hashes new password and updates user record.
- Resets `FailedResetAttempts` to 0 and clears OTP.
- After 3 failed attempts, sends a security alert email.

### Success Response (JSON)
```json
{
  "message": "Password has been reset successfully."
}
```

### Error Cases
- **400 Bad Request**: Invalid OTP or weak password.
- **404 Not Found**: User does not exist.
- **500 Internal Server Error**: Database or hashing issues.

---

## 6. Delete User (POST `/deleteUser`)
### Description
Deletes a user from the database.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu"
}
```

### Behavior
- Finds and deletes user record.

### Success Response (JSON)
```json
{
  "message": "User has been deleted."
}
```

### Error Cases
- **404 Not Found**: User does not exist.
- **500 Internal Server Error**: Database operation failure.

---

## 7. Display User (POST `/displayUser`)
### Description
Fetches user details by email.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu"
}
```

### Success Response (JSON)
```json
{
  "userid": 101,
  "name": "Jinx Silco",
  "email": "getjinxed@ufl.edu",
  "phone": "+15551234567"
}
```

### Error Cases
- **404 Not Found**: User does not exist.
- **500 Internal Server Error**: Database issues.

---

## 8. Login (POST `/login`)
### Description
Authenticates the user with email and password.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu",
  "password": "EkkoRew1nd!"
}
```

### Success Response (JSON)
```json
{
  "message": "Login successful."
}
```

### Error Cases
- **401 Unauthorized**: Incorrect password.
- **404 Not Found**: User does not exist.
- **500 Internal Server Error**: Database or hashing issues.

# Appendix: Security & Design Choices

This section explains key security and design decisions behind the UniBazaar user authentication system.

## 1. Why Use SendGrid for Email Notifications?
SendGrid was chosen as the email provider for sending OTP codes and security alerts because:
- **Reliability:** SendGrid offers a scalable and highly available cloud-based email service.
- **Security:** It supports **DKIM, SPF, and TLS encryption**, ensuring secure email transmission.
- **Rate Limits & Throttling:** SendGrid provides rate limiting to prevent abuse (e.g., brute-force OTP requests).
- **Easy API Integration:** The SendGrid API is well-documented and integrates smoothly with Go-based backend services.
- **Monitoring & Analytics:** Logs and analytics help track email delivery rates and failures.

### Alternative Considerations
Other email services such as **AWS SES** and **Mailgun** were considered, but SendGrid was selected due to its **free tier for transactional emails** and superior developer tooling.

---

## 2. OTP Generation & Security Measures
To enhance authentication security, UniBazaar uses a **random 6-digit OTP** for email verification and password reset.

- **Generated Using Cryptographic Randomness:** The OTP is created using Go's `crypto/rand` package to ensure unpredictability.
- **Short Expiry Time:** OTP codes expire within a limited time window (e.g., 5 minutes) to reduce brute-force attempts.
- **Failure Tracking & Lockout Mechanism:** If a user enters an incorrect OTP multiple times (3 or more failures), the system sends a **security alert email**.


This ensures **OTP codes are unique, secure, and difficult to guess**, enhancing authentication security.


---
## 3. No Unit Tests for Email-Sent OTPs
Currently, there are no unit tests for verifying whether OTP emails are actually sent.

- **Relies on External Services:** Email delivery depends on **SendGrid**, making it difficult to test in an isolated backend environment.
- **Better suited for Integration Testing:** Instead of unit tests, **integration tests** with a mock SendGrid API should be used to validate OTP email sending.
- **Logs & Provider Monitoring:** The best way to ensure email delivery works correctly is through **server logs and SendGrid’s dashboard** rather than unit tests.

---

## 4. No OTP for Phone Verification
At the moment, phone verification is **not implemented** due to the following reasons:

- **Cost Considerations:** SendGrid (owned by Twilio) provides a free tier for **email-based OTPs**, but **phone-based OTP verification via Twilio SMS requires payment**.
- **Target Audience:** Since **UniBazaar is a student-focused platform**, we aim to minimize operational costs.
- **Exploring Alternatives:** Future phone verification strategies may include:
  - **Optional Twilio verification for premium users**.
  - **Using an alternative SMS provider with lower costs**.
  - **Validating phone numbers via a third-party API instead of OTP verification**.

---

## 5. Incomplete JWT Implementation
Currently, the JWT implementation is **incomplete**, and further enhancements are planned in the next sprint.

- **What Has Been Implemented So Far:**
  - **JWT Verification Handler:** Extracts, parses, and validates JWT tokens from API requests.
  - **JWT Generation Handler:** Creates JWT tokens upon user authentication.
- **What’s Missing & Planned for Next Sprint:**
  - **Unit tests for token parsing and claim extraction**.
  - **Integration tests for end-to-end authentication flow**.
  - **Token expiration handling and refresh token implementation**.
  - **Stronger JWT validation mechanisms to prevent misuse**.

By implementing these security improvements in the upcoming sprint, **UniBazaar will have a more robust authentication system** with tested JWT handling and improved security.

---

## 6. Why Use Argon2id for Password Hashing?
**Argon2id** is the recommended password hashing algorithm by OWASP and is used in UniBazaar due to:

- **Memory-Hardness:** Argon2id is resistant to GPU and ASIC-based brute force attacks due to its high memory requirements.
- **Customizable Parameters:** It allows tuning of memory, iterations, and parallelism for optimal security.
- **Resistance to Timing Attacks:** Unlike bcrypt, Argon2id provides protection against cache-timing attacks.

---

## 7. Password Entropy Enforcement (go-password-validator)

To enforce strong password security, UniBazaar uses `go-password-validator` with a **minimum entropy requirement of 60 bits**. This ensures that passwords are not easily guessable by:

- **Forcing Complexity:** Users must have a unique, hard to guess password.
- **Mitigating Dictionary Attacks:** Prevents users from choosing common or easily cracked passwords.
- **Providing Real-Time Feedback:** If a password is weak, the API returns a message guiding the user to create a stronger one.

By enforcing **entropy-based validation**, UniBazaar prevents users from choosing weak passwords while maintaining usability.

---

## Conclusion
These security measures were carefully chosen to **protect user data, prevent unauthorized access, and ensure system integrity**. 
- **Argon2id** provides robust password protection against brute force attacks.
- **SendGrid** ensures reliable OTP delivery with built-in monitoring and security features.
- **Password entropy validation** enforces strong credentials to mitigate password-based attacks.
- **OTP generation and validation** mechanisms enhance security while ensuring user convenience.

By implementing these industry-standard best practices, UniBazaar ensures a **secure, scalable, and resilient** authentication system.

## Users: Unit Tests
The unit tests are located in `unit_test.go` and cover the following functionalities:

### 1. User Insertion
- **Test:** `TestUserInsert`
- **Description:** Verifies that a user can be inserted into the database.
- **Expected Behavior:** The insertion function returns no error.

### 2. User Retrieval
- **Test:** `TestUserRead`
- **Description:** Tests reading a user from the database.
- **Expected Behavior:** The correct user object is returned without errors.

### 3. Update User Name
- **Test:** `TestUpdateUserName`
- **Description:** Ensures that a user’s name can be updated successfully.
- **Expected Behavior:** The update function completes without errors.

### 4. Update User Phone
- **Test:** `TestUpdateUserPhone`
- **Description:** Validates the ability to update a user’s phone number.
- **Expected Behavior:** The function executes successfully without errors.

### 5. Delete User
- **Test:** `TestDeleteUser`
- **Description:** Tests the deletion of a user from the database.
- **Expected Behavior:** The delete function completes without errors.

### 6. Initiate Password Reset
- **Test:** `TestPasswordReset`
- **Description:** Ensures that a password reset request can be initiated.
- **Expected Behavior:** The function returns no errors.

### 7. Verify Reset Code and Set New Password
- **Test:** `TestVerifyResetCodeAndSetNewPassword`
- **Description:** Verifies that a reset code can be validated and a new password can be set.
- **Expected Behavior:** The function returns no errors.

### 8. Validate `.edu` Email Addresses
- **Test:** `TestValidateEduEmail`
- **Description:** Checks whether only `.edu` email addresses are accepted.
- **Expected Behavior:** Valid `.edu` emails pass, while non-`.edu` emails return errors.

### 9. Validate Password Strength
- **Test:** `TestValidatePassword`
- **Description:** Ensures that passwords meet security requirements.
- **Expected Behavior:** Weak passwords return errors, strong passwords pass validation.

### 10. Validate Phone Numbers
- **Test:** `TestValidatePhone`
- **Description:** Checks the format of phone numbers.
- **Expected Behavior:** Valid numbers pass, invalid numbers return errors.

## Running Tests
To execute the unit tests, use the following command:
```sh
 go test -v ./...
```
This will run all test cases and display detailed output.

## Dependencies
The tests utilize the following dependencies:
- `github.com/stretchr/testify/assert` for assertions
- `github.com/stretchr/testify/mock` for mocking user model methods

Ensure these dependencies are installed before running tests:
```sh
go mod tidy
go get github.com/stretchr/testify
```

## Conclusion
These unit tests help ensure the reliability of the backend user management functionalities by validating the core operations such as user creation, update, deletion, authentication, and validation processes.


---


# Messaging System API Documentation

## Table of Contents
1. [Overview](#overview)
2. [Endpoints Overview](#endpoints-overview)
3. [Database Connection](#1-database-connection)
4. [WebSocket Messaging](#2-websocket-messaging)
5. [Message Handling](#3-message-handling)
6. [User Handling](#4-user-handling)
7. [Data Models](#5-data-models)
8. [Repository Functions](#6-repository-functions)
9. [Conclusion](#conclusion)

---

## Overview
The Messaging System API provides functionality for real-time messaging between users. It supports WebSocket connections for live message transmission, REST endpoints for fetching messages, and user management operations. The system also includes chat history persistence, ensuring messages remain accessible even if a user gets disconnected.

---

## Endpoints Overview
Below is a summary of the available endpoints in this API:

| **Method** | **Endpoint**                   | **Description**                                                                 | **Usage**                                |
|------------|---------------------------------|---------------------------------------------------------------------------------|------------------------------------------|
| `GET`      | `/ws?user_id={user_id}`         | WebSocket connection for real-time messaging.                                  | Establishes WebSocket connection.       |
| `POST`     | `/send`                         | Sends a message from one user to another.                                       | Accepts a JSON payload to send a message.|
| `GET`      | `/messages?sender_id={sender_id}&receiver_id={receiver_id}` | Retrieves messages exchanged between two users.                               | Fetches messages by sender and receiver.|
| `GET`      | `/users`                        | Retrieves all registered users in the system.                                   | Fetches a list of all users.            |

---

## 1. Database Connection

### `ConnectDB()`
**Description:** Establishes a connection to the PostgreSQL database.  
**Returns:** `*sql.DB` (database connection instance)

---

## 2. WebSocket Messaging

### `HandleWebSocket(w http.ResponseWriter, r *http.Request)`
**Method:** `GET`  
**Endpoint:** `/ws?user_id={user_id}`  
**Description:** Upgrades an HTTP connection to WebSocket and registers the user as a client.  
**Query Parameters:**
- `user_id` (integer) - The ID of the user connecting to the WebSocket.

**Connection Details:**
- When a client connects, the system assigns a persistent session.
- If a user gets disconnected, the system retains their chat history for seamless recovery upon reconnection.
- Heartbeat signals ensure the connection remains active, and reconnections are handled automatically.

---

## 3. Message Handling

### `HandleSendMessage(w http.ResponseWriter, r *http.Request)`
**Method:** `POST`  
**Endpoint:** `/send`  
**Description:** Accepts a JSON payload to send a message.  
**Request Body:**
```json
{
  "sender_id": 1,
  "receiver_id": 2,
  "content": "Hello!"
}
```
**Response:**
```json
{
  "status": "message sent"
}
```

### `HandleGetMessages(w http.ResponseWriter, r *http.Request)`
**Method:** `GET`  
**Endpoint:** `/messages?sender_id={sender_id}&receiver_id={receiver_id}`  
**Description:** Retrieves messages exchanged between two users.  
**Query Parameters:**
- `sender_id` (integer) - Sender user ID.
- `receiver_id` (integer) - Receiver user ID.

**Response:** List of messages:
```json
[
  {
    "id": 1,
    "sender_id": 1,
    "receiver_id": 2,
    "content": "Hello!",
    "timestamp": 1700000000,
    "read": false
  }
]
```

---

## 4. User Handling

### `GetUsersHandler(w http.ResponseWriter, r *http.Request)`
**Method:** `GET`  
**Endpoint:** `/users`  
**Description:** Fetches all registered users.  
**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
]
```

---

## 5. Data Models

### `Message`
```json
{
  "id": 1,
  "sender_id": 1,
  "receiver_id": 2,
  "content": "Hello!",
  "timestamp": 1700000000,
  "read": false
}
```

### `User`
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

---

## 6. Repository Functions

### `SaveMessage(msg models.Message) error`
- Inserts a new message into the database.

### `GetLatestMessages(limit int) ([]models.Message, error)`
- Retrieves the latest messages up to the specified limit.

### `MarkMessageAsRead(messageID int) error`
- Updates a message as read.

### `GetUnreadMessages(userID uint) ([]models.Message, error)`
- Fetches all unread messages for a user.

### `GetConversation(userID uint) ([]models.Message, error)`
- Retrieves all messages where the user is either the sender or receiver.

---

## Conclusion
This API facilitates real-time and stored messaging functionalities through WebSockets and REST endpoints, enabling seamless communication between users. The system ensures chat history persistence, so conversations remain intact even if users experience connectivity issues. 

---
# Products: Backend API Documentation
For detailed information, you can view the Swagger API specification here:

[Swagger API Specification](https://github.com/SakshiPandey97/UniBazaar/blob/product-enhancements-and-fixes/Backend/products/docs/swagger.yaml)

## Table of Contents
1. [Endpoints Overview](#endpoints)
2. [Detailed Endpoint Documentation](#endpoint-documentation)
   - [Get All Products (GET /products)](#1-get-products)
   - [Get Produtcs By User Id (GET /producs/{userId})](#2-get-products-user)
   - [Create Product (POST /products)](#3-create-product)
   - [Update Product (PUT /products/{userId}/{productId})](#4-update-product)
   - [Delete Product (DELETE /products/{userId}/{productId})](#5-delete-product)
3. [Unit Tests](#unit-tests)

---

## Endpoints Overview

Below is a quick reference to each endpoint:

| **Endpoint**       | **Method** | **Description**                                       |
|--------------------|-----------|--------------------------------------------------------|
| `/products`         | GET      | Get all products from database.                      |
| `/producs/{userId}`    | GET      | Get all products belonging to a particular user.  |
| `/products` | POST      | Create a new product in database.                           |
| `/products/{userId}/{productId}` | PUT      | Update a user in database.              |
| `/products/{userId}/{productId}`     | DELETE      | Delete a product from database.  |

---


### 1. **Get All Products** - `GET /products`

#### Description
Fetches all products from the system, regardless of the user ID. If no products are found, an error is returned.

#### Response
- **200 OK**: Returns a list of all products.
- **404 Not Found**: If no products are found in the system.
- **500 Internal Server Error**: For database issues.

#### Example Response (200 OK)
```json
[
  {
    "productId": "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
    "productTitle": "Laptop",
    "productDescription": "A high-performance laptop",
    "productPrice": 999.99,
    "productCondition": 4,
    "productLocation": "University of Florida",
    "productImage": "https://example.com/laptop.jpg",
    "productPostDate": "02-20-2025",
    "userId": 123
  }
]
```

### 2. Create a New Product - POST /products

#### Description
Creates a new product by parsing form data, uploading images to S3, and saving the details in the database. The product is linked to the user via their User ID.

## Request Parameters (Form Data)
- **UserId** (integer): *Required*. The user ID linking to the product.
- **productTitle** (string): *Required*. Title of the product.
- **productDescription** (string): *Optional*. Description of the product.
- **productPrice** (number): *Required*. Price of the product.
- **productCondition** (integer): *Required*. Condition of the product.
- **productLocation** (string): *Required*. Location of the product.
- **productImage** (file): *Required*. Image of the product.

## Response
- **201 Created**: The product was successfully created.
- **400 Bad Request**: Invalid User ID or form data.
- **500 Internal Server Error**: Server issues.

#### Example Response (201 OK)
```json
[
  {
    "productId": "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
    "productTitle": "Laptop",
    "productDescription": "A high-performance laptop",
    "productPrice": 999.99,
    "productCondition": 4,
    "productLocation": "University of Florida",
    "productImage": "https://example.com/laptop.jpg",
    "productPostDate": "02-20-2025",
    "userId": 123
  }
]
```

#### Example Response (400 Bad Request)
```json
{
    "error": "Error reading image",
    "details": "error retrieving file: http: no such file"
}
```

### 3. Get Products by User ID - GET /products/{UserId}

#### Description
Fetches all products listed by a user, identified by their user ID. If no products are found, an error is returned.

## Request Parameters
- **UserId** (integer): *Required*. The unique user ID.

## Response
- **200 OK**: Returns a list of products.
- **400 Bad Request**: Invalid user ID.
- **404 Not Found**: No products found for the given user ID.
- **500 Internal Server Error**: Server issues.

#### Example Response (200 OK)
```json
[
  {
    "productId": "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
    "productTitle": "Laptop",
    "productDescription": "A high-performance laptop",
    "productPrice": 999.99,
    "productCondition": 4,
    "productLocation": "University of Florida",
    "productImage": "https://example.com/laptop.jpg",
    "productPostDate": "02-20-2025",
    "userId": 123
  }
]
```

#### Example Response (404 Not Found)
```json
{
    "error": "Error fetching products for user",
    "details": "No products found for user ID: 45: no products found"
}
```

### 4. Update Product by User ID and Product ID - PUT /products/{UserId}/{ProductId}

#### Description
Updates a product's details based on the user ID and product ID. The product image is also updated if provided.

## Request Parameters
- **UserId** (integer): *Required*. The unique user ID.
- **ProductId** (string): *Required*. The unique product ID.
- **productTitle** (string): *Required*. The updated title of the product.
- **productDescription** (string): *Optional*. The updated description of the product.
- **productPrice** (number): *Required*. The updated price of the product.
- **productCondition** (integer): *Required*. The updated condition of the product.
- **productLocation** (string): *Required*. The updated location of the product.
- **productImage** (string): *Optional*. Image of the product.

## Response
- **200 OK**: The product was updated successfully.
- **400 Bad Request**: Invalid request data.
- **404 Not Found**: Product not found.
- **500 Internal Server Error**: Server issues.

#### Example Response Body
```json
{
    "productId": "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
    "productTitle": "Laptop",
    "productDescription": "A high-performance laptop",
    "productPrice": 999.99,
    "productCondition": 4,
    "productLocation": "University of Florida",
    "productImage": "Image",
    "productPostDate": "02-20-2025",
    "userId": 123
  }
  ```

#### Example Response (404 Not Found)
```json
  {
    "error": "Error reading image",
    "details": "error retrieving file: http: no such file"
}
```

### 5. Delete Product by User ID and Product ID - DELETE /products/{UserId}/{ProductId}

#### Description
Deletes a product from the system based on the user ID and product ID. This also removes the associated image from S3 if available.

## Request Parameters
- **UserId** (integer): *Required*. The unique user ID.
- **ProductId** (string): *Required*. The unique product ID.

## Response
- **204 No Content**: The product was successfully deleted.
- **400 Bad Request**: Invalid request data.
- **404 Not Found**: Product not found.
- **500 Internal Server Error**: Server issues.

#### Example Response (204 No Content)
```json
{}
```

#### Example Response (404 Not Found)
```json
{
    "error": "Error fetching product",
    "details": "Product not found for UserId: 456 and ProductId: 26678cba-459c-45a8-b856-a333ae4e0356: no products found"
}
```

---

## Conclusion

The **Products API** is designed to provide a robust, efficient, and scalable solution for managing product data within the UniBazaar platform. By utilizing well-structured endpoints, advanced database features, and optimized data handling mechanisms, the API ensures:

- **Fast retrieval of product data** with efficient filtering, searching, and data indexing.
- **Secure CRUD operations** with validation mechanisms.
- **Scalability** Uses go routines to handle a growing product catalog and high traffic.

With this API, UniBazaar can effectively manage its product offerings while ensuring smooth, secure, and optimized interactions for users, sellers, and administrators alike. Future enhancements can focus on expanding features like product search and security to further enrich the user experience.

## Unit Tests
## Product Handler: Unit Tests

### 1. Create Product
- **Test:** `TestCreateProductHandler`
- **Description:** Verifies that a product can be created successfully with an image upload.
- **Expected Behavior:** The product is created, and the status code is `201 Created`.

### 2. Get All Products
- **Test:** `TestGetAllProductsHandler`
- **Description:** Tests retrieving all products from the database.
- **Expected Behavior:** A list of products is returned with a status code of `200 OK`.

### 3. Get Products by User ID
- **Test:** `TestGetAllProductsByUserIDHandler`
- **Description:** Tests retrieving products specific to a user based on the user ID.
- **Expected Behavior:** A list of products for the user is returned with a status code of `200 OK`.

### 4. Update Product
- **Test:** `TestUpdateProductHandler`
- **Description:** Ensures that an existing product can be updated successfully, including updating the product image.
- **Expected Behavior:** The product is updated, and the status code is `200 OK`.

### 5. Delete Product
- **Test:** `TestDeleteProductHandler`
- **Description:** Verifies that a product can be deleted, including the associated image.
- **Expected Behavior:** The product is deleted, and the status code is `204 No Content`.
---
### Helper Functions: Unit Tests

#### 1. GetUserID: Valid Input
- **Test:** `TestGetUserID_ValidInput`
- **Description:** Verifies that a valid user ID returns the correct result.
- **Expected Behavior:** The function returns `123` when input `"123"` is passed.

#### 2. GetUserID: Invalid Input
- **Test:** `TestGetUserID_InvalidInput`
- **Description:** Tests the `GetUserID` function with invalid inputs.
- **Expected Behavior:** An error is returned for inputs such as `"abc"`, `""`, `"12.34"`, and `"-"`.

#### 3. ParseFormAndCreateProduct: Valid Data
- **Test:** `TestParseFormAndCreateProduct_ValidData`
- **Description:** Tests the creation of a product when valid form data is provided.
- **Expected Behavior:** A product is created with correct user ID, title, and other fields.

#### 4. ParseFormAndCreateProduct: Missing or Invalid Data
- **Test:** `TestParseFormAndCreateProduct_MissingOrInvalidData`
- **Description:** Tests handling of invalid or missing form data.
- **Expected Behavior:** The function returns an error for missing or invalid data such as missing product title or invalid product condition and price.

#### 5. ParseNumericalFormValues: Valid Data
- **Test:** `TestParseNumericalFormValues_ValidData`
- **Description:** Verifies that numerical form values are correctly parsed.
- **Expected Behavior:** The correct product condition and price are set in the product object.

#### 6. ParseNumericalFormValues: Invalid Data
- **Test:** `TestParseNumericalFormValues_InvalidData`
- **Description:** Tests the function with invalid numerical values.
- **Expected Behavior:** The function returns an error when invalid numerical data is provided for condition or price.
---
## Helper Functions: Unit Tests

#### 1. Create Mock Image
- **Test:** `CreateMockImage`
- **Description:** Creates a mock image (either JPEG or PNG) for testing purposes.
- **Expected Behavior:** A red-colored 100x100 image is created and returned as a byte array.

#### 2. Parse Product Image: Error Retrieving File
- **Test:** `TestParseProductImage_ErrorRetrievingFile`
- **Description:** Tests the `ParseProductImage` function when there is an error retrieving the file from the request.
- **Expected Behavior:** The function returns an error containing the string "error retrieving file".

#### 3. Parse Product Image: Error Encoding JPEG
- **Test:** `TestParseProductImage_ErrorEncodingJPEG`
- **Description:** Tests the `ParseProductImage` function when there is an error encoding the image in JPEG format.
- **Expected Behavior:** The function returns an error containing the string "error encoding compressed image".

#### 4. Parse Product Image: Error Encoding PNG
- **Test:** `TestParseProductImage_ErrorEncodingPNG`
- **Description:** Tests the `ParseProductImage` function when there is an error encoding the image in PNG format.
- **Expected Behavior:** The function returns an error containing the string "error encoding compressed image".

#### 5. Parse Product Image: Unsupported Format
- **Test:** `TestParseProductImage_UnsupportedFormat`
- **Description:** Tests the `ParseProductImage` function when an unsupported image format (GIF) is uploaded.
- **Expected Behavior:** The function returns an error containing the string "error decoding image".
---
## Model Functions: Unit Tests

#### 1. Test Error Response Serialization
- **Test:** `TestErrorResponseSerialization`
- **Description:** Tests the serialization of the `ErrorResponse` struct to JSON.
- **Expected Behavior:** The function should correctly marshal the `ErrorResponse` into the expected JSON string.
- **Test Case:**
    ```json
    {"error":"Error updating product","details":"ProductPrice: cannot be empty or zero, Product not found"}
    ```

#### 2. Test Error Response Deserialization
- **Test:** `TestErrorResponseDeserialization`
- **Description:** Tests the deserialization of JSON into the `ErrorResponse` struct.
- **Expected Behavior:** The function should correctly unmarshal the JSON string into the `ErrorResponse` struct, with the error and details fields populated correctly.
- **Test Case (JSON Input):**
    ```json
    {"error":"Error updating product","details":"ProductPrice: cannot be empty or zero, Product not found"}
    ```

#### 3. Test Error Response Serialization Without Details
- **Test:** `TestErrorResponseSerializationWithoutDetails`
- **Description:** Tests serialization of the `ErrorResponse` struct with no details field.
- **Expected Behavior:** The function should correctly marshal the `ErrorResponse` into the expected JSON string with only the `error` field.
- **Test Case:**
    ```json
    {"error":"Error updating product"}
    ```

#### 4. Test Empty Error Response
- **Test:** `TestEmptyErrorResponse`
- **Description:** Tests serialization of an empty `ErrorResponse` struct.
- **Expected Behavior:** The function should correctly marshal the empty `ErrorResponse` struct into the expected JSON string with only the `error` field as an empty string.
- **Test Case:**
    ```json
    {"error":""}
    ```
---
## Model Functions: Unit Tests

#### 1. Test Product Validation
- **Test:** `TestProductValidation`
- **Description:** Tests the validation of both valid and invalid `Product` instances.
- **Expected Behavior:** The valid product should not return an error, while the invalid product should return validation errors.
- **Test Case:**
    - Valid Product:
    ```json
    {
        "UserID": 123,
        "ProductID": "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
        "ProductTitle": "Laptop",
        "ProductDescription": "A high-performance laptop",
        "ProductPostDate": "02-20-2025",
        "ProductCondition": 4,
        "ProductPrice": 999.99,
        "ProductLocation": "University of Florida",
        "ProductImage": "https://example.com/laptop.jpg"
    }
    ```
    - Invalid Product:
    ```json
    {
        "UserID": 0,
        "ProductID": "",
        "ProductTitle": "",
        "ProductPostDate": "02-20-2025",
        "ProductCondition": 0,
        "ProductPrice": 0,
        "ProductLocation": "",
        "ProductImage": ""
    }
    ```

#### 2. Test Product Post Date Validation
- **Test:** `TestProductPostDateValidation`
- **Description:** Tests validation of the `ProductPostDate` field to ensure it is in MM-DD-YYYY format.
- **Expected Behavior:** The valid date should not return an error, while the invalid date should trigger an error with the message "productPostDate must be in MM-DD-YYYY format".
- **Test Case (Valid Date):**
    ```json
    {
        "ProductPostDate": "02-20-2025"
    }
    ```
- **Test Case (Invalid Date):**
    ```json
    {
        "ProductPostDate": "2025-02-20"
    }
    ```

#### 3. Test Product Validation with Empty Fields
- **Test:** `TestProductValidationWithEmptyFields`
- **Description:** Tests the validation of a `Product` with empty or zero values for required fields.
- **Expected Behavior:** The function should return an error for missing or invalid required fields.
- **Test Case:**
    ```json
    {
        "UserID": 0,
        "ProductID": "",
        "ProductTitle": "",
        "ProductCondition": 0,
        "ProductPrice": 0,
        "ProductLocation": "",
        "ProductImage": ""
    }
    ```

#### 4. Test Format Validation Error
- **Test:** `TestFormatValidationError`
- **Description:** Tests the formatting of validation error messages.
- **Expected Behavior:** The function should format the error message as expected, e.g., `"ProductTitle: cannot be empty or zero"`.
- **Test Case:**
    ```json
    {
        "Error": "ProductTitle: zero value"
    }
    ```

#### 5. Test Empty Product
- **Test:** `TestEmptyProduct`
- **Description:** Tests validation of an empty `Product` struct.
- **Expected Behavior:** The function should return an error indicating that the product is empty and invalid.
- **Test Case:**
    ```json
    {
        "UserID": 0,
        "ProductID": "",
        "ProductTitle": "",
        "ProductDescription": "",
        "ProductCondition": 0,
        "ProductPrice": 0,
        "ProductLocation": "",
        "ProductImage": ""
    }
    ```

#### 6. Test Format Validation Error with Nil Error
- **Test:** `TestFormatValidationErrorWithNilError`
- **Description:** Tests the behavior of `formatValidationError` when a nil error is passed.
- **Expected Behavior:** The function should return nil when the input error is nil.
---
## S3 Repository: Unit Tests

#### 1. Test UploadImage Failure
- **Test:** `TestUploadImage_Failure`
- **Description:** Tests the failure scenario for uploading an image to S3.
- **Expected Behavior:** The function should return an error when the upload fails and the URL should be empty.
- **Test Case:**
    - Mocked error response for S3 upload.
    - Expected behavior: Error should be returned, URL should be empty.

#### 2. Test DeleteImage Failure
- **Test:** `TestDeleteImage_Failure`
- **Description:** Tests the failure scenario for deleting an image from S3.
- **Expected Behavior:** The function should return an error when the delete operation fails.
- **Test Case:**
    - Mocked error response for S3 delete.
    - Expected behavior: Error should be returned.

#### 3. Test GeneratePresignedURL Failure
- **Test:** `TestGeneratePresignedURL_Failure`
- **Description:** Tests the failure scenario for generating a presigned URL for an image in S3.
- **Expected Behavior:** The function should return an error when generating the presigned URL fails and the URL should be empty.
- **Test Case:**
    - Mocked error response for presigned URL generation.
    - Expected behavior: Error should be returned, URL should be empty.

#### 4. Test GetPreSignedURLs Success
- **Test:** `TestGetPreSignedURLs_Success`
- **Description:** Tests the success scenario for generating presigned URLs for multiple products.
- **Expected Behavior:** The function should return the presigned URLs for the provided products.
- **Test Case:**
    - Sample Products.
    - Expected behavior: The result should contain 2 items, each with a presigned URL.

---
## Mongo Repository: Unit Tests

#### 1. Test CreateProduct Success
- **Test:** `TestCreateProduct`
- **Description:** Tests the successful creation of a product in the repository.
- **Expected Behavior:** The function should return no error when the product is successfully created.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: No error should be returned.

#### 2. Test CreateProduct Error
- **Test:** `TestCreateProduct_Error`
- **Description:** Tests the failure scenario for creating a product in the repository.
- **Expected Behavior:** The function should return an error when the product creation fails.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: Error should be returned due to the insertion failure.

#### 3. Test GetAllProducts Success
- **Test:** `TestGetAllProducts`
- **Description:** Tests the successful retrieval of all products from the repository.
- **Expected Behavior:** The function should return all the products with no error.
- **Test Case:** 
    - Two products with IDs `prod123` and `prod456`.
    - Expected behavior: The result should contain both products with correct IDs.

#### 4. Test GetProductsByUserID Success
- **Test:** `TestGetProductsByUserID`
- **Description:** Tests the successful retrieval of products for a specific user.
- **Expected Behavior:** The function should return all products for the specified user.
- **Test Case:** 
    - Two products for `UserID: 1` with IDs `prod123` and `prod456`.
    - Expected behavior: The result should contain both products associated with the given user ID.

#### 5. Test UpdateProduct Success
- **Test:** `TestUpdateProduct`
- **Description:** Tests the successful update of a product in the repository.
- **Expected Behavior:** The function should return no error when the product is successfully updated.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: No error should be returned.

#### 6. Test DeleteProduct Success
- **Test:** `TestDeleteProduct`
- **Description:** Tests the successful deletion of a product from the repository.
- **Expected Behavior:** The function should return no error when the product is successfully deleted.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: No error should be returned.

#### 7. Test FindProductByUserAndId Success
- **Test:** `TestFindProductByUserAndId`
- **Description:** Tests the successful retrieval of a product by user and product ID.
- **Expected Behavior:** The function should return the product if found.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: The product should be returned with the correct product ID.

---
## Routes: Unit Tests

#### 1. Test Register Product Routes
- **Test:** `TestRegisterProductRoutes`
- **Description:** Tests the registration of product routes, including POST and GET requests.
- **Expected Behavior:** 
    - POST request to `/products` should trigger `CreateProduct` on the `MockProductRepository`.
    - GET request to `/products` should trigger `GetAllProducts` on the `MockProductRepository`.
- **Test Case:**
    - **POST /products:** The handler should call `CreateProduct` with the provided product and return a success response.
    - **GET /products:** The handler should call `GetAllProducts` and return an empty list in the response.

#### 2. Test CORS Headers
- **Test:** `TestCORSHeaders`
- **Description:** Tests that the CORS headers are correctly set on the response.
- **Expected Behavior:** 
    - The response should include `Access-Control-Allow-Origin: *`, `Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS`, and `Access-Control-Allow-Headers: Content-Type,Authorization`.
- **Test Case:**
    - **OPTIONS /products:** The handler should set the correct CORS headers in the response. The `Access-Control-Allow-Origin` should be `*`, `Access-Control-Allow-Methods` should include all necessary HTTP methods, and `Access-Control-Allow-Headers` should include `Content-Type` and `Authorization`.


## Running Tests
To execute the unit tests, use the following command:
```sh
 go test -coverprofile=coverage ./... 
```
```sh
 go tool cover -html=coverage
```
This will run all test cases and display detailed coverage in browser.

Current Coverage:
| Package                     | Duration | Coverage           |
|-----------------------------|----------|--------------------|
| `web-service/handler`        | 0.423s   | 62.2% of statements|
| `web-service/helper`         | 0.391s   | 83.0% of statements|
| `web-service/model`          | 0.176s   | 100.0% of statements|
| `web-service/repository`     | 0.130s   | 28.8% of statements|
| `web-service/routes`         | 0.122s   | 100.0% of statements|



## Dependencies
The tests utilize the following dependencies:
- `github.com/stretchr/testify/assert` for assertions
- `github.com/stretchr/testify/mock` for mocking user model methods

Ensure these dependencies are installed before running tests:
```sh
go mod tidy
go get github.com/stretchr/testify
```

## Conclusion

These unit tests help ensure the reliability of the backend product management functionalities by validating core operations such as product creation, update, deletion, retrieval, and validation processes.

---
# Frontend Unit Testing

---

## Testing Summary

This project utilizes **Vitest** for running unit and integration tests to ensure proper functionality, performance, and reliability of the application. **React Testing Library** is used for rendering components, simulating user interactions, and checking component behavior through assertions.

### Key Features:
- **Unit Tests**: Validating individual functions and components to ensure expected behavior.
- **UI Tests**: Verifying the rendering of UI elements and interactions, such as button clicks or form submissions.
- **Mocking**: Mocking external API calls and services to test components in isolation.
- **Assertions**: Using `expect()` to check if components or values meet the expected results.

### How to Run Tests:
1. Install dependencies using `npm install`.
2. Run tests with the command:
   ```bash
   npx vitest
   ```

The tests cover critical areas of the app, including UI rendering, state management, API interactions, and more, helping maintain code quality and application stability.

---

### 1. **User Login API**
- **Test:** Should handle successful login
- **Expected Behavior:** Returns user ID '12345' and calls `localStorage.setItem`.

### 2. **User Registration API**
- **Test:** Should handle successful registration
- **Expected Behavior:** Returns success object and calls `axios.post`.

### 3. **Fetch All Users API**
- **Test:** Should fetch all users
- **Expected Behavior:** Returns an array of users excluding the specified user.

### 4. **Fetch All Products API**
- **Test:** Should fetch all products
- **Expected Behavior:** Returns an array of products.

### 5. **Post Product API**
- **Test:** Should post a new product
- **Expected Behavior:** Returns the newly created product.

---

### 6. **Banner Rendering**
- **Test:** Should render banner text
- **Expected Behavior:** Renders "Uni", "Bazaar", and "Connecting students for buying/selling".

---

### 7. **Input Rendering**
- **Test:** Should render input field with label
- **Expected Behavior:** Renders input with the correct label.

### 8. **Disabled Input**
- **Test:** Should disable input when submitting
- **Expected Behavior:** Input is disabled when `isSubmitting` is true.

---

### 9. **Product Rendering**
- **Test:** Should render title, image, and price
- **Expected Behavior:** Renders product title, image, and price correctly.

---

### 10. **Products Loading**
- **Test:** Should display loading spinner
- **Expected Behavior:** Spinner shown when `loading` is true.

### 11. **Products Error**
- **Test:** Should display error message
- **Expected Behavior:** Shows error message "Error fetching products".

---

### 12. **Spinner Rendering**
- **Test:** Should render spinner
- **Expected Behavior:** Spinner is displayed.

---

### 13. **Valid Login**
- **Test:** Should allow login with valid credentials
- **Expected Behavior:** `handleSubmit` called with valid credentials.

---

### 14. **Registration Submission**
- **Test:** Should call handleSubmit with correct values
- **Expected Behavior:** Calls `handleSubmit` with correct form values.

---

### 15. **Initial State**
- **Test:** Should initialize with `isAnimating` as false
- **Expected Behavior:** `isAnimating` is `false` initially.

---

### 16. **Navbar State**
- **Test:** Should return initial menu and dropdown state
- **Expected Behavior:** `isMenuOpen` and `isDropdownOpen` are both `false`.

---
