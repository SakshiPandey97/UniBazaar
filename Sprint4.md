# Sprint 4

## Overview
**Sprint 4** focused on finalizing key user-facing features and enhancing overall platform stability and usability. This sprint included the implementation of some core authentication flows, significant messaging system improvements, and multiple frontend enhancements for better responsiveness and UX. Backend integration for chat and user services was also completed, alongside updates for session management and search capabilities. The team also addressed several important bugs and optimized various UI flows to ensure a seamless experience. The project was successfully completed and deployed to the cloud, and is now live at: https://uni-bazaar.vercel.app/.

---

## Sprint 4: Completed Issues

| **Issue** | **Status** | **Type** |
|---|---|---|
| ForgotPassword Functionality | Closed | frontend |
| Add new unit tests and fix existing ones for the modified functions | Closed | bug, enhancement |
| Host messaging server on Azure | Closed | backend, enhancement, frontend |
| Fix bug causing ws connection to close while chatting | Closed | backend, bug, frontend |
| Update FrontPage | Closed | frontend |
| Redesign About Us Page | Closed | frontend |
| Use enter to send message in chat app | Closed | enhancement, frontend |
| Redirect to Home on Successful OTP Verification | Closed | frontend |
| Create View Profile | Closed | frontend |
| Fix: Remove current user(self), shown in contact list | Closed | enhancement, frontend |
| chat app layout not consistent for different screen sizes | Closed | enhancement, frontend |
| Add resend OTP functionality | Closed | backend |
| Fix displayusers to use ID | Closed | backend |
| Add functionality for the users to initiate chat from products card | Closed | frontend |
| Adding new users to chat db when login | Closed | backend, frontend |
| Logout on Protected Routes | Closed | frontend |
| Integrate Http cookie for maintaining login state | Closed | frontend |
| Add functionality to edit and delete user's products | Closed | enhancement, frontend |
| Create My products page and integrate with GET by ID API | Closed | enhancement, frontend |
| Link Search from home page to products page | Closed | frontend |
| Add Search Product functionality on Front End | Closed | enhancement, frontend |
| Host User Service on Azure | Closed | backend |

---

## Summary of Work Done in Sprint 4
Throughout Sprint 4, the team brought UniBazaar across the finish line by implementing and polishing all critical user‑facing features, from authentication flows and OTP management to a fully integrated chat system. Frontend enhancements delivered a more responsive, intuitive UI—updating pages like About Us and FrontPage, and streamlining search and messaging interactions. On the backend, user and messaging services were containerized and deployed to Azure, and unit tests were expanded to ensure long‑term stability. Finally, we squashed numerous cross‑platform bugs, optimized our services, and successfully deployed the completed project to Vercel, making the platform live at https://uni-bazaar.vercel.app/.

## UniBazaar User API Documentation
This section of the document describes the UniBazaar backend API for user management. It covers the **design choices**, **request/response formats**, **error handling**, and **sample JSON** requests for each endpoint. 

---

# Users: Backend API Documentation

## Table of Contents
1. [Design & Implementation Highlights](#design--implementation-highlights)  
2. [Endpoints Overview](#endpoints-overview)  
3. [Detailed Endpoint Documentation](#detailed-endpoint-documentation)  
   - [Health Check (GET `/`)](#1-health-check-get-)  
   - [Sign Up (POST `/signup`)](#2-sign-up-post-signup)  
   - [Verify Email (POST `/verifyEmail`)](#3-verify-email-post-verifyemail)  
   - [Resend OTP (POST `/resendOtp`)](#4-resend-otp-post-resendotp)  
   - [Forgot Password (GET `/forgotPassword`)](#5-forgot-password-post-forgotpassword)  
   - [Update Password (POST `/updatePassword`)](#6-update-password-post-updatepassword)  
   - [Delete User (POST `/deleteUser`)](#7-delete-user-post-deleteuser)  
   - [Display User (GET `/displayUser/:id`)](#8-display-user-get-displayuserid)  
   - [Login (POST `/login`)](#9-login-post-login)  
   - [Logout (POST `/logout`)](#10-logout-post-logout)  
   - [Get JWT (POST `/getjwt`)](#11-get-jwt-post-getjwt)  
   - [Verify JWT (GET `/verifyjwt`)](#12-verify-jwt-get-verifyjwt)  
   - [Update Name (POST `/updateName`)](#13-update-name-post-updatename)  
   - [Update Phone (POST `/updatePhone`)](#14-update-phone-post-updatephone)  
   - [Error Cases & Responses](#error-cases--responses)  
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

| **Endpoint**               | **Method** | **Description**                                                       |
|----------------------------|------------|-----------------------------------------------------------------------|
| `/`                        | GET        | Health check – returns “OK”.                                          |
| `/signup`                  | POST       | Create a new user and send an email‑verification OTP.                 |
| `/verifyEmail`             | POST       | Verify user’s email with the OTP.                                    |
| `/resendOtp`               | POST       | Resend the email‑verification OTP.                                    |
| `/forgotPassword`          | GET        | Generate and email a password‑reset OTP.                              |
| `/updatePassword`          | POST       | Verify reset OTP and update the password.                             |
| `/deleteUser`              | POST       | Delete a user by email.                                               |
| `/displayUser/:id`         | GET        | Fetch user details by **user ID**.                                    |
| `/login`                   | POST       | Authenticate user (email+password) and return a JWT.                 |
| `/logout`                  | POST       | Revoke the current JWT (by `jti`).                                    |
| `/getjwt`                  | POST       | Generate a test JWT from given user data.                             |
| `/verifyjwt`               | GET        | Validate a provided JWT and its revocation status.                   |
| `/updateName`              | POST       | Update the user’s name (requires credentials).                        |
| `/updatePhone`             | POST       | Update the user’s phone number (requires credentials).                |

---

# Detailed Endpoint Documentation

### 1. Health Check (GET `/`)
**Description:**  
Returns **200 OK** and body `"OK"` to indicate the service is up.

---

### 2. Sign Up (POST `/signup`)
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

**Success (200):**
```
Sign‑up successful. Check your email for OTP.
```

### Error Cases
- **400 Bad Request**: Invalid email, weak password, or incorrect phone format.
- **409 Conflict**: User already exists.
- **500 Internal Server Error**: Database or email-sending issues.

---

### 3. Verify Email (POST `/verifyEmail`)
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

**Success (200):**
```json
{ "verified": true }
```

### Error Cases
- **400 Bad Request**: Invalid or expired OTP.
- **404 Not Found**: User does not exist.
- **500 Internal Server Error**: Database or email-sending issues.

---
### 4. Resend OTP (POST `/resendOtp`)
**Description:**  
Generates and emails a new verification OTP.

**Request Body:**
```json
{ "email": "sakshi@ufl.edu" }
```

**Success (200):**
```
OTP resent successfully. Check your e‑mail.
```

**Errors:**
- 400 Bad Request – invalid JSON.
- 404 Not Found – user not found.
- 500 Internal Server Error – email‑send failure.

---
### 5. Forgot Password (GET `/forgotPassword`)
**Description:**  
Initiates a password reset by sending a reset OTP. Requires the user’s email as a query parameter.

**Request:**  
```
GET /forgotPassword?email=user@ufl.edu
```

**Success (200):**  
```
Reset code sent. Check your email.
```

**Errors:**  
- 400 Bad Request – missing or empty `email` query parameter.  
- 404 Not Found – user not found.  
- 500 Internal Server Error – database or email failure.

---
### 6. Verify Reset Code (POST `/updatePassword`)
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

### 7. Delete User (POST `/deleteUser`)
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
### 8. Display User (GET `/displayUser/:id`)

### Description  
Retrieves user details based on their unique user ID.

This endpoint takes a `GET` request with the user ID passed as a URL parameter. The handler validates the ID and fetches the corresponding user from the database.

### Example Request  
```
GET http://localhost:4000/displayUser/101
```

### Success Response  
Returns a JSON object with user details:

```json
{
  "userid": 101,
  "name": "Jinx Silco",
  "email": "getjinxed@ufl.edu",
  "phone": "+15551234567"
}
```

### Error Responses  
- **400 Bad Request**: Invalid or missing user ID.  
- **404 Not Found**: No user found with the given ID.  
- **500 Internal Server Error**: Issue accessing the database or encoding the response.

---
### 9. Login (POST `/login`)
### Description
Authenticates the user using email and password. Returns a JWT valid for 48 hours and the user's ID if credentials are valid and the account is verified.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu",
  "password": "EkkoRew1nd!"
}
```

**Success (200):**
```json
{
  "userId": 101,
  "token": "<JWT>",
  "name": "First Last",
  "email": "user@ufl.edu"
}
```

### Error Cases
- **401 Unauthorized**: Invalid credentials or unverified account.
- **500 Internal Server Error**: Error during hashing or validation.

---
### 10. Logout (POST `/logout`)
**Description:**  
Revokes the current JWT (by storing its `jti`).

**Header:**
```
Authorization: Bearer <JWT>
```

**Success (200):**
```
Logout successful, token revoked.
```

**Errors:**
- 400 Bad Request – malformed header.
- 401 Unauthorized – invalid or expired token.

---
### 11. Get JWT (POST `/getjwt`)
**Description:**  
Test endpoint to generate a JWT from provided user info (no DB lookup).

**Request Body:**
```json
{
  "name":  "First Last",
  "email": "user@ufl.edu",
  "phone": "+15551234567"
}
```

**Success (200):**  
Returns `Authorization: Bearer <JWT>` header and body:
```
JWT generated successfully!
```

### Error Cases
- **500 Internal Server Error**: Token generation failure.

---

### 12. Verify JWT (GET `/verifyjwt`)
### Description
Validates a provided JWT to ensure it is not expired or revoked.

### Header
```
Authorization: Bearer <token>
```

**Success (200):**  
Body:  
```
Token valid. User: {Name: "First Last", Email: "user@ufl.edu", Phone: "+15551234567"}
```

**Errors:**
- 400 Bad Request – malformed header.
- 401 Unauthorized – invalid, expired, or revoked token.


---

### 13. Update Name (POST `/updateName`)
### Description
Updates the user's name after verifying the user's identity through email and password.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu",
  "password": "EkkoRew1nd!",
  "newName": "Jinx Silco Updated"
}
```

### Behavior
- Authenticates user by verifying email and password.
- Updates user's name in the database.

### Success Response (JSON)
```json
{
  "message": "Name updated successfully."
}
```

### Error Cases
- **400 Bad Request**: Invalid JSON input or database update failure.
- **401 Unauthorized**: User not found or invalid credentials.
- **500 Internal Server Error**: Password hashing/verification issues.

---

### 14. Update Phone (POST `/updatePhone`)
### Description
Updates the user's phone number after verifying the user's identity through email and password.

### Request Body (JSON)
```json
{
  "email": "getjinxed@ufl.edu",
  "password": "EkkoRew1nd!",
  "newPhone": "+15559876543"
}
```

### Behavior
- Authenticates user by verifying email and password.
- Updates user's phone number in the database.

### Success Response (JSON)
```json
{
  "message": "Phone updated successfully."
}
```

### Error Cases
- **400 Bad Request**: Invalid JSON input or database update failure.
- **401 Unauthorized**: User not found or invalid credentials.
- **500 Internal Server Error**: Password hashing/verification issues.



### Common Error Codes
| Status Code | Meaning |
|-------------|---------|
| 400 Bad Request | Malformed/missing input data |
| 401 Unauthorized | Invalid credentials or token |
| 404 Not Found | User record not found |
| 500 Internal Server Error | DB/hash/OTP/email issues |


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

## 2. Why Use Argon2id for Password Hashing?
**Argon2id** is the recommended password hashing algorithm by OWASP and is used in UniBazaar due to:

- **Memory-Hardness:** Argon2id is resistant to GPU and ASIC-based brute force attacks due to its high memory requirements.
- **Customizable Parameters:** It allows tuning of memory, iterations, and parallelism for optimal security.
- **Resistance to Timing Attacks:** Unlike bcrypt, Argon2id provides protection against cache-timing attacks.

---

## 3. Password Entropy Enforcement (go-password-validator)

To enforce strong password security, UniBazaar uses `go-password-validator` with a **minimum entropy requirement of 60 bits**. This ensures that passwords are not easily guessable by:

- **Forcing Complexity:** Users must have a unique, hard to guess password.
- **Mitigating Dictionary Attacks:** Prevents users from choosing common or easily cracked passwords.
- **Providing Real-Time Feedback:** If a password is weak, the API returns a message guiding the user to create a stronger one.

### Example of Password Entropy Validation
```go
const minEntropyBits = 60
err := passwordvalidator.Validate(password, minEntropyBits)
if err != nil {
    return fmt.Errorf("password is too weak: %v", err)
}
```

By enforcing **entropy-based validation**, UniBazaar prevents users from choosing weak passwords while maintaining usability.

---

## 4. OTP Generation & Security Measures
To enhance authentication security, UniBazaar uses a **random 6-digit OTP** for email verification and password reset.

- **Generated Using Cryptographic Randomness:** The OTP is created using Go's `crypto/rand` package to ensure unpredictability.
- **Short Expiry Time:** OTP codes expire within a limited time window (e.g., 5 minutes) to reduce brute-force attempts.
- **Failure Tracking & Lockout Mechanism:** If a user enters an incorrect OTP multiple times (3 or more failures), the system sends a **security alert email**.


This ensures **OTP codes are unique, secure, and difficult to guess**, enhancing authentication security.

---
## 5. JWT-Based Authentication
Sprint 3 introduced JWT endpoints to improve session management and stateless authentication:

- **JWT Generation:** The GetJWTHandler endpoint creates a signed token embedding user data. Tokens are generated using a secret stored in an environment variable.

- **JWT Verification:** The VerifyJWTHandler validates tokens and extracts user information. It checks that the token hasn’t been revoked via a global in-memory revocation map.

- **Token Revocation (Logout):** The LogoutHandler revokes the current token by marking its unique identifier (jti) as invalid, ensuring that a logged-out token cannot be reused.

- **Security Considerations:** JWT tokens are signed with HS256, and their expiry and issuance times (iat/exp) are enforced to prevent replay and timing attacks.
---

## Conclusion
These security measures were carefully chosen to **protect user data, prevent unauthorized access, and ensure system integrity**. 
- **Argon2id** provides robust password protection against brute force attacks.
- **SendGrid** ensures reliable OTP delivery with built-in monitoring and security features.
- **Password entropy validation** enforces strong credentials to mitigate password-based attacks.
- **OTP generation and validation** mechanisms enhance security while ensuring user convenience.
- **Error Handling** robust error handling and login/signup flow.
- **JWT functionality** enhances session management by enabling secure token generation, validation, and revocation. This stateless approach reduces server load while maintaining strong security practices.


By implementing these industry-standard best practices, UniBazaar ensures a **secure, scalable, and resilient** authentication system.

## Users: Unit Tests
The users unit tests are located in `unit_test.go` and cover the following functionalities:

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

## Unit Tests in utils_test.go
The tests in `utils_test.go` focus on JWT functionality and the associated helper functions:

### 1. TestGenerateJWT
- **Description:** Validates that a JWT can be generated from a user struct without error.
- **Expected Behavior:** A non-empty token string is returned.

### 2. TestParseJWTValidToken
- **Description:** Ensures that a generated JWT is parsed successfully and contains the correct user claims.
- **Expected Behavior:** The token is valid, and user claims (UserID, Name, Email) match the input.

### 3. TestParseJWTInvalidToken
- **Description:** Verifies that an invalid token (malformed string) is rejected.
- **Expected Behavior:** Parsing returns an error and a nil token.

### 4. TestExpiredToken
- **Description:** Checks that tokens with past expiration dates are rejected.
- **Expected Behavior:** The token is not parsed, and an appropriate error is returned.

## Unit Tests in handler_test.go
These tests simulate HTTP requests to verify that the REST API endpoints behave as expected:

1. **SignUpHandler** (POST `/signup`)
   - **Description:** Simulates a sign-up request with valid user details.
   - **Expected Behavior:** Returns HTTP 200 and creates a new user with `Verified` set to false.

2. **VerifyEmailHandler** (POST `/verifyEmail`)
   - **Description:** Processes email verification using an OTP code.
   - **Expected Behavior:** Returns HTTP 200 and sets the user’s `Verified` status to true.

3. **ForgotPasswordHandler** (GET `/forgotPassword?email={email}`)
   - **Description:** Generates a new OTP for a verified user.
   - **Expected Behavior:** Returns HTTP 200 and updates the user record with a new `OTPCode`.

4. **UpdatePasswordHandler** (POST `/updatePassword`)
   - **Description:** Updates the user’s password using a valid OTP.
   - **Expected Behavior:** Returns HTTP 200, updates and re-hashes the password, and clears the `OTPCode`.

5. **DisplayUserHandler** (GET `/displayUser/{id}`)
   - **Description:** Retrieves user details by user ID.
   - **Expected Behavior:** Returns HTTP 200 and the correct user data in JSON.

6. **UpdateNameHandler** (POST `/updateName`)
   - **Description:** Updates the user’s name after validating the current password. Request body should include `email`, `password`, and `newName`.
   - **Expected Behavior:** Returns HTTP 200 and updates the user’s `Name` field.

7. **UpdatePhoneHandler** (POST `/updatePhone`)
   - **Description:** Updates the user’s phone number after validating the current password. Request body should include `email`, `password`, and `newPhone`.
   - **Expected Behavior:** Returns HTTP 200 and updates the user’s `Phone` field.

8. **DeleteUserHandler** (POST `/deleteUser`)
   - **Description:** Deletes a user record by email. Request body should include `email`.
   - **Expected Behavior:** Returns HTTP 200 and removes the user from the database.

## Running Tests
To execute the unit tests, use the following command:
```sh
 go test -v ./...
```
This will run all test cases and display a detailed output.

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

