# Sprint 1: User Stories and Implementation Plan

## User Stories

### Authentication & User Management

**US1: User Registration**
- As a university student, I want to sign up using my email so that I can access UniBazaar's services.
- Acceptance Criteria:
  - User can enter their name, email, and password.
  - User receives an OTP for verification.
  - Only university email addresses (`.edu`) are accepted.
  - User account is created upon successful OTP verification.

**US2: User Login**
- As a registered user, I want to log in using my email and password so that I can access my account.
- Acceptance Criteria:
  - User can enter email and password to log in.
  - System validates credentials against the PostgreSQL database.
  - On successful login, a session is created.
  - If credentials are invalid, an error message is displayed.

**US3: Session Management**
- As a user, I want my session to be tracked so that I don’t need to log in every time.
- Acceptance Criteria:
  - Implement session tracking using JWT or cookies.
  - Sessions should expire after a set duration.

### Product Management

**US4: Post a Product**
- As a seller, I want to post a product for sale so that I can find buyers.
- Acceptance Criteria:
  - User can upload a product image (stored in AWS S3).
  - User enters product details (title, description, price, category).
  - Product data is stored in NoSQL database.
  - Product appears in the marketplace upon successful posting.

**US5: Browse Products**
- As a buyer, I want to browse available products so that I can find what I need.
- Acceptance Criteria:
  - Products are categorized for easy browsing.
  - Users can filter and search for specific items.
  - Product images and details are displayed.

**US6: Send a Quote to Buy a Product**
- As a buyer, I want to send a quote to the seller so that I can negotiate the price.
- Acceptance Criteria:
  - Users can click "Send Quote" on a product.
  - A message is sent to the seller.
  - The system logs the quote request.

### Messaging & Communication

**US7: Chat System Implementation**
- As a user, I want to chat with buyers/sellers so that I can negotiate and finalize deals.
- Acceptance Criteria:
  - A messaging system is implemented.
  - Users can send and receive messages in real time.
  - Chat history is stored and accessible.

### Sprint 1: Planned Issues and Assignments

**Frontend Tasks:**
1. Implement the Landing Screen ✅ *(Completed by Tanmay)*
2. Implement Login & Registration screens with form validation ✅ *(Completed by Shubham)*
4. Implement Product Listing UI ✅ *(Completed by Tanmay)*

**Backend Tasks:**
1. **Set up PostgreSQL for user authentication** ✅ *(Completed by Sakshi)*
2. **Implement User Signup/Login API** ✅ *(Completed by Sakshi)*
- **Validation:**
     - Ensures only Florida institutes' `.edu` emails are accepted.
     - First and Last name must be entered.
     - Email Format validation.
     - Uses `go-password-validator` for password entropy validation.
     - Enforces strong passwords with at least 60 entropy bits
   - Passwords hashed using **Argon2id**.
   - **Security:** Uses `argon2id` with the following parameters:
     - Memory: `128 * 1024`
     - Iterations: `4`
     - Parallelism: `NumCPU()`
     - SaltLength: `16`
     - KeyLength: `32`
   - **CRUD Functions Implemented:**
     - `Create`: Add a new user with hashed password.
     - `Read`: Fetches user details.
     - `Update`: Updates the user's password securely.
     - `Delete`: Deletes user account.
   

3. **Set up NoSQL database for products** ✅ *(Completed by Avaneesh)*
- Setup connection with Mongo DB to store and update product details.
- Mongo DB hosted on Mongo DB Atlas with Free Tier.

4. **Set up NoSQL database for products** ✅ *(Completed by Shubham and Avaneesh)*
- Setup connection with AWS to store and retrieve product images and generate presigned urls for sharng with UI.
- AWS S3 setup Free Tier.

4. **Implement Product APIs** ✅ *(Completed by Avaneesh)*
- POST products/
  - Create a new product entry in Mongo DB with all the details.
  - Upload images received in POST request body to AWS S3 and save the location in Mongo DB.
- GET products/ 
  - Get all the products in the dabase to populate the landing page.
  - Returns presigned URL in aboev JSON for product images to populate on databases.
- Update products/ queries: userId, productId
  - Updates Product details using `productId` belonging to the user `userId`.
- Delete products/ queries: userId, productId
  - Deleted Product details using `productId` belonging to the user `userId`.


### Sprint 1 Completion Status

| Task | Status |
|------|--------|
| Implement Landing Screen | ✅ Completed |
| Implement Login & Registration screens | ✅ Completed |
| Implement Product Listing UI | ✅ Completed |
| Set up PostgreSQL for user authentication | ✅ Completed |
| Implement User Signup/Login API | ✅ Completed |
| Set up NoSQL database for products | ✅ Completed |
| Implement Product APIs | ✅ Completed |
| Set up AWS S3 to store and retrieve images for products | ✅ Completed |

### Challenges Faced
- **Frontend and Backend Sync Issues:** Some APIs require modifications based on frontend feedback.
- **Login Integration Challenges:** Encountered issues integrating the login and registration UI with the backend.
- **Responsive Design Difficulties:** Faced challenges implementing dynamic page layout adjustments based on screen size.

### Next Steps for Sprint 2
- Complete pending tasks (OTP verification, messaging UI, quote sending API integration).
- Conduct testing for user authentication and product posting.
- Implement user profile management.
- Improve UI/UX based on feedback.
- Work on data consistency and roll back mechanism between S3 and Mongo DB.
- Work on Updating and Deleting images from S3.
- Work on image compression.
- Work on validation checks in API.
- Setup API spec and contracts to enforece validaiton.
---
- Dev-Team Videos:
- **Frontend Demo:** [https://youtu.be/kYsLyq-FL60](https://youtu.be/kYsLyq-FL60)
- **Backend Demo:** [https://youtu.be/WtBir6qzsLs](https://youtu.be/WtBir6qzsLs)


