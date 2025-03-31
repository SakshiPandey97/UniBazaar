# Sprint 3: User Stories and Implementation Plan

## Overview
During **Sprint 3**, our primary focus was on enhancing user management, strengthening account security, improving backend performance, implementing messaging system in backend and refining product-related functionalities. Below is an overview of the issues we tackled, the tasks completed, and the corresponding backend API documentation for user operations.

---

## Sprint 3: Completed Issues

| Issue                                                                                       | Status   | Type                 |
| ------------------------------------------------------------------------------------------- | -------- | -------------------- |
| Pagination support for GET endpoints in products services                                   | ✅ Closed | Backend, Sprint v3  |
| Add a Search API to support full text search over products                                  | ✅ Closed | Backend, Sprint v3  |
| Custom Error Handling in produtcs API to return appropriate status and messages efficiently | ✅ Closed | Backend, Sprint v3  |
| Create products page and load paginated results into the UI                                 | ✅ Closed | Frontend, Sprint v3 |
| Implement frontend of messaging system                                                      | ✅ Closed | Frontend, Sprint v3 |
| Refactor FE messaging to improve modularity                                                 | ✅ Closed | Frontend, Sprint v3 |
| Chat application scrolling issue                                                            | ✅ Closed | Frontend, Sprint v3 |
| Migrated messaging system to AWS                                                            | ✅ Closed | Backend, Sprint v3  | 
| Add unit tests for frontend messaging system                                                | ✅ Closed | Frontend, Sprint v3 |
| Add unit tests for backend messaging system                                                 | ✅ Closed | Backend, Sprint v3  |
| Migrated from npm to pnpm                                                                   | ✅ Closed | Frontend, Sprint v3  |
| Implemented UI for Forgot Password                                                          | ✅ Closed | Frontend, Sprint v3  |
| Implemented UI for OTP Validation                                                           | ✅ Closed | Frontend, Sprint v3  |
| Implemented AboutUs Page                                                                    | ✅ Closed | Frontend, Sprint v3  |
| ReDesigned UI for the Website                                                           | ✅ Closed | Frontend, Sprint v3  |
| Added Dynamic views for all the webpages                                                           | ✅ Closed | Frontend, Sprint v3  |


---

## Summary of Work Done in Sprint 3

### 1. **Products Pagination Support**
   - **Implemented Pagination for Product Listings:** 
     - Added pagination capabilities to the product retrieval logic, which helps in handling large data sets efficiently. Introduced a fast pipeline and aggregation method that fetches products using two key parameters sorted by `ProductId`:
       - `lastID`: Allows fetching products after a specific product ID, facilitating continuous product retrieval.
       - `limit`: Restricts the number of products fetched at a time, optimizing performance and reducing memory usage.
     - This pagination mechanism ensures that users can seamlessly browse through products without overwhelming the system or frontend, even when the number of products increases significantly.
   - **Updated Methods:** 
     - Modified the `GET products/` and `GET products/{userId}/` endpoints to include the `lastID` and `limit` query parameters.

### 2. **Product Search Endpoint**
   - **Implemented Product Search:** 
     - Introduced a search functionality for products, allowing users to query products based on specific terms in the product title and description. This is achieved using a fast pipeline and aggregation method with MongoDB's full-text search capabilities. The search supports two key parameters:
       - `query`: The search query used to match products based on their title and description.
       - `limit`: Restricts the number of products fetched at a time, optimizing performance and reducing memory usage.
     - The search logic utilizes MongoDB’s aggregation pipeline to perform a full-text search across the `ProductTitle` and `ProductDescription` fields, and it returns matching products in a fast and efficient manner.
   - **Updated Endpoints:**
     - Added the `GET /search/products` endpoint to handle product searches, allowing users to search products by query and limit the results. 

### 3. **Custom Error Handling**
   - **Introduced Custom Error Handling:**
     - Implemented a custom error structure to streamline error handling across the application. The new error handling system includes several error types to represent specific error cases:
       - **NotFoundError**: Represents cases where a resource is not found (e.g., no products matching a search query).
       - **DatabaseError**: Represents database-related errors.
       - **S3Error**: Represents errors encountered during interactions with Amazon S3 or other storage solutions.
       - **BadRequestError**: Represents client-side errors due to invalid input or bad requests.
     - Each error type provides a consistent structure with the ability to store an error message, status code, and underlying cause, making it easier to track and respond to errors.

   - **Updated Response Handling:**
     - Refined error response handling to return appropriate status codes and error messages based on the custom error types. The `HandleError` function now processes errors with custom logic to determine the status code, message, and error details.
     - The success response handler (`HandleSuccessResponse`) was also introduced to return successful responses in a consistent format with proper status codes and data serialization.
     - This ensures that all API responses, both error and success, are standardized and informative, improving the API's robustness and user experience.

### 4. **Products Page with Pagination and Load More**
   - **Paginated Listings:**
     - Developed a `ProductsPage` with a "Load More" button, fetching 12 products per request and using the `lastId` for pagination.
     - Implemented `useAllProducts` hook to track the last fetched product and avoid duplicates.
   
   - **Error Handling and Loading:**
     - Displayed a loading spinner when products are being fetched and showed a relevant error message if the request fails.
     - Added a fallback UI that allows the user to retry fetching products in case of an error.
   
   - **Efficient Data Fetching:**
     - Used `useRef` to track already fetched products and stop requests when all products are loaded.
     - Used conditional logic to stop making requests once all products have been fetched, providing a smooth browsing experience for users.

---
# UniBazaar User API Documentation

This section of the document describes the UniBazaar backend API for user management. It covers the **design choices**, **request/response formats**, **error handling**, and **sample JSON** requests for each endpoint. 

---
---

# Products: Backend API Documentation
For detailed information, you can view the Swagger API specification here:

[Swagger API Specification](https://github.com/SakshiPandey97/UniBazaar/blob/main/Backend/products/docs/swagger.yaml)

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

| **Endpoint**                                      | **Method** | **Description**                                  |
| ------------------------------------------------- | ---------- | ------------------------------------------------ |
| `/products?lastId={lastId}&limit={limit}`         | GET        | Get all products from database.                  |
| `/producs/{userId}?lastId={lastId}&limit={limit}` | GET        | Get all products belonging to a particular user. |
| `/search/products?query={query}&limit={limit}`    | GET        | Search products with given query.                |
| `/products`                                       | POST       | Create a new product in database.                |
| `/products/{userId}/{productId}`                  | PUT        | Update a user in database.                       |
| `/products/{userId}/{productId}`                  | DELETE     | Delete a product from database.                  |

---


## 1. **Get All Products** - `GET /products`

### Description
Fetches all products from the system, regardless of the user ID. If no products are found, an error is returned.

### Query Parameters (Pagination)
| Name   | Type   | Required | Description                                                                                  |
| ------ | ------ | -------- | -------------------------------------------------------------------------------------------- |
| lastId | string | ❌ No     | The lastId of the product to start the pagination from (If empty starts from first product). |
| limit  | int    | ❌ No     | The maximum number of results to return (default is 10).                                     |


### Response
- **200 OK**: Returns a list of all products.
- **404 Not Found**: If no products are found in the system.
- **500 Internal Server Error**: For database issues.

### Example Response (200 OK)
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

## 2. Create a New Product - POST /products

### Description
Creates a new product by parsing form data, uploading images to S3, and saving the details in the database. The product is linked to the user via their User ID.

### Request Parameters (Form Data)
- **UserId** (integer): *Required*. The user ID linking to the product.
- **productTitle** (string): *Required*. Title of the product.
- **productDescription** (string): *Optional*. Description of the product.
- **productPrice** (number): *Required*. Price of the product.
- **productCondition** (integer): *Required*. Condition of the product.
- **productLocation** (string): *Required*. Location of the product.
- **productImage** (file): *Required*. Image of the product.

### Response
- **201 Created**: The product was successfully created.
- **400 Bad Request**: Invalid User ID or form data.
- **500 Internal Server Error**: Server issues.

### Example Response (201 OK)
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

### Example Response (400 Bad Request)
```json
{
    "error": "Error reading image",
    "details": "error retrieving file: http: no such file"
}
```

## 3. Get Products by User ID - GET /products/{UserId}

### Description
Fetches all products listed by a user, identified by their user ID. If no products are found, an error is returned.

### Request Parameters
- **UserId** (integer): *Required*. The unique user ID.

### Query Parameters
| Name   | Type   | Required | Description                                              |
| ------ | ------ | -------- | -------------------------------------------------------- |
| lastId | string | ✅ Yes    | The lastId of the product to start the pagination from.  |
| limit  | int    | ❌ No     | The maximum number of results to return (default is 10). |

### Response
- **200 OK**: Returns a list of products.
- **400 Bad Request**: Invalid user ID.
- **404 Not Found**: No products found for the given user ID.
- **500 Internal Server Error**: Server issues.

### Example Response (200 OK)
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

### Example Response (404 Not Found)
```json
{
    "error": "Error fetching products for user",
    "details": "No products found for user ID: 45: no products found"
}
```

## 4. Update Product by User ID and Product ID - PUT /products/{UserId}/{ProductId}

### Description
Updates a product's details based on the user ID and product ID. The product image is also updated if provided.

### Request Parameters
- **UserId** (integer): *Required*. The unique user ID.
- **ProductId** (string): *Required*. The unique product ID.
- **productTitle** (string): *Required*. The updated title of the product.
- **productDescription** (string): *Optional*. The updated description of the product.
- **productPrice** (number): *Required*. The updated price of the product.
- **productCondition** (integer): *Required*. The updated condition of the product.
- **productLocation** (string): *Required*. The updated location of the product.
- **productImage** (string): *Optional*. Image of the product.

### Response
- **200 OK**: The product was updated successfully.
- **400 Bad Request**: Invalid request data.
- **404 Not Found**: Product not found.
- **500 Internal Server Error**: Server issues.

### Example Response Body
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

### Example Response (404 Not Found)
```json
  {
    "error": "Error reading image",
    "details": "error retrieving file: http: no such file"
}
```

## 5. Delete Product by User ID and Product ID - DELETE /products/{UserId}/{ProductId}

### Description
Deletes a product from the system based on the user ID and product ID. This also removes the associated image from S3 if available.

### Request Parameters
- **UserId** (integer): *Required*. The unique user ID.
- **ProductId** (string): *Required*. The unique product ID.

### Response
- **204 No Content**: The product was successfully deleted.
- **400 Bad Request**: Invalid request data.
- **404 Not Found**: Product not found.
- **500 Internal Server Error**: Server issues.

### Example Response (204 No Content)
```json
{}
```

### Example Response (404 Not Found)
```json
{
    "error": "Error fetching product",
    "details": "Product not found for UserId: 456 and ProductId: 26678cba-459c-45a8-b856-a333ae4e0356: no products found"
}
```
## 6. **Search Products** - `GET /search/products`

### Description
Fetches all products from the system, regardless of the user ID. If no products are found, an error is returned.

### Parameters
| Name  | Type   | Required | Description                                              |
| ----- | ------ | -------- | -------------------------------------------------------- |
| query | string | ✅ Yes    | The search query for filtering products.                 |
| limit | int    | ❌ No     | The maximum number of results to return (default is 10). |

### Response
- **200 OK**: Returns a list of all products.
- **404 Not Found**: If no products are found in the system.
- **500 Internal Server Error**: For database issues.

### Example Response (200 OK)
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
---

## Conclusion

The **Products API** is designed to provide a robust, efficient, and scalable solution for managing product data within the UniBazaar platform. By utilizing well-structured endpoints, advanced database features, and optimized data handling mechanisms, the API ensures:

- **Fast retrieval of product data** with efficient filtering, searching, and data indexing.
- **Secure CRUD operations** with validation mechanisms.
- **Scalability** Uses go routines to handle a growing product catalog and high traffic.

With this API, UniBazaar can effectively manage its product offerings while ensuring smooth, secure, and optimized interactions for users, sellers, and administrators alike. Future enhancements can focus on expanding features like product search and security to further enrich the user experience.

## Unit Tests
## Product Handler: Unit Tests

#### 1. Search Products (NEW) 🚀
- **Test:** `TestSearchProductsHandler`
- **Description:** Tests the `SearchProductsHandler` function to ensure it returns the correct list of products based on the search query and limit.
- **Expected Behavior:** Returns a `200 OK` status with the list of products and their pre-signed image URLs.

#### 2. Handle Error (NEW) 🚀
- **Test:** `TestHandleError`
- **Description:** Tests the `HandleError` function to ensure proper error handling and response formatting.
- **Scenarios:**
  - **Database Error:** Returns `500 Internal Server Error` with the correct error message.
  - **Not Found Error:** Returns `404 Not Found` when the error is a resource-not-found error.
  - **S3 Error:** Returns `500 Internal Server Error` for S3-related issues.
  - **Bad Request Error:** Returns `400 Bad Request` for invalid input errors.
  - **Default Error:** Returns `500 Internal Server Error` for generic errors.
  - **Nil Error:** Returns `500 Internal Server Error` with a message indicating no specific error.

#### 3. Handle Success Response (NEW) 🚀
- **Test:** `TestHandleSuccessResponse`
- **Description:** Tests the `HandleSuccessResponse` function to verify proper response handling for successful operations.
- **Scenarios:**
  - **Success with Data:** Returns `200 OK` with a JSON response.
  - **Success with Empty Data:** Returns `201 Created` with a `null` response body.
  - **Success with Integer Data:** Returns `200 OK` with an integer in the response body.
  - **Success with String Data:** Returns `200 OK` with a string in the response body.

#### 4. Handle Success Response with Encoding Error (NEW) 🚀
- **Test:** `TestHandleSuccessResponse_EncodingError`
- **Description:** Tests how `HandleSuccessResponse` handles encoding failures.
- **Expected Behavior:** If the data cannot be encoded into JSON, it returns `500 Internal Server Error` with an appropriate error message.

#### 5. Create Product
- **Test:** `TestCreateProductHandler`
- **Description:** Verifies that a product can be created successfully with an image upload.
- **Expected Behavior:** The product is created, and the status code is `201 Created`.

#### 6. Get All Products
- **Test:** `TestGetAllProductsHandler`
- **Description:** Tests retrieving all products from the database.
- **Expected Behavior:** A list of products is returned with a status code of `200 OK`.

#### 7. Get Products by User ID
- **Test:** `TestGetAllProductsByUserIDHandler`
- **Description:** Tests retrieving products specific to a user based on the user ID.
- **Expected Behavior:** A list of products for the user is returned with a status code of `200 OK`.

#### 8. Update Product
- **Test:** `TestUpdateProductHandler`
- **Description:** Ensures that an existing product can be updated successfully, including updating the product image.
- **Expected Behavior:** The product is updated, and the status code is `200 OK`.

#### 9. Delete Product
- **Test:** `TestDeleteProductHandler`
- **Description:** Verifies that a product can be deleted, including the associated image.
- **Expected Behavior:** The product is deleted, and the status code is `204 No Content`.
---
## Mongo Repository: Unit Tests

#### 1. Test Search Products Success (NEW) 🚀
- **Test:** `TestSearchProducts_Success`
- **Description:** Verifies that products are returned correctly when a valid search query is provided.
- **Expected Behavior:** The function should return a list of products matching the query.
- **Test Case:**
    - Search query: `"laptop"`, Limit: `2`.
    - Expected products: 
        - `ProductTitle: "Laptop X1"`, `ProductID: "prod123"`, `UserID: 1`.
        - `ProductTitle: "Gaming Laptop"`, `ProductID: "prod456"`, `UserID: 2`.

#### 2. Test Search Products No Results (NEW) 🚀
- **Test:** `TestSearchProducts_NoResults`
- **Description:** Tests the behavior when no products match the search query.
- **Expected Behavior:** The function should return an error with a message indicating no products were found.
- **Test Case:**
    - Search query: `"laptop"`, Limit: `2`.
    - Expected behavior: 
        - Error: `"no products found"`.
        - An empty product list.

#### 3. Test Search Products Database Error (NEW) 🚀
- **Test:** `TestSearchProducts_DatabaseError`
- **Description:** Verifies the behavior when there is a database error during the product search.
- **Expected Behavior:** The function should return an error indicating a database issue.
- **Test Case:**
    - Search query: `"laptop"`, Limit: `2`.
    - Expected behavior:
        - Error: `"database error"`.
        - An empty product list.

#### 4. Test Get All Products Pagination Success (NEW) 🚀
- **Test:** `TestGetAllProducts_Pagination_Success`
- **Description:** Verifies that products are returned correctly when pagination is used with a valid last product ID.
- **Expected Behavior:** The function should return a paginated list of products based on the `lastID` and `limit`.
- **Test Case:**
    - Last product ID: `"prod100"`, Limit: `2`.
    - Expected products:
        - `ProductID: "prod123"`, `UserID: 1`.
        - `ProductID: "prod456"`, `UserID: 2`.

#### 5. Test Get All Products Pagination No Last ID Success (NEW) 🚀
- **Test:** `TestGetAllProducts_Pagination_NoLastID_Success`
- **Description:** Verifies that products are returned correctly when pagination is used without a last product ID (starting from the first product).
- **Expected Behavior:** The function should return a paginated list of products based on the `limit` provided.
- **Test Case:**
    - Last product ID: `""` (empty), Limit: `3`.
    - Expected products:
        - `ProductID: "prod123"`, `UserID: 1`.

        - `ProductID: "prod456"`, `UserID: 2`.
        - `ProductID: "prod789"`, `UserID: 3`.

#### 6. Test Get All Products Pagination Empty Results (NEW) 🚀
- **Test:** `TestGetAllProducts_Pagination_EmptyResults`
- **Description:** Verifies the behavior when no products are found for a given `lastID` in pagination.
- **Expected Behavior:** The function should return an empty list when no products are found.
- **Test Case:**
    - Last product ID: `"prod999"`, Limit: `2`.
    - Expected behavior: 
        - Empty product list.

#### 7. Test CreateProduct Success
- **Test:** `TestCreateProduct`
- **Description:** Tests the successful creation of a product in the repository.
- **Expected Behavior:** The function should return no error when the product is successfully created.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: No error should be returned.

#### 8. Test CreateProduct Error
- **Test:** `TestCreateProduct_Error`
- **Description:** Tests the failure scenario for creating a product in the repository.
- **Expected Behavior:** The function should return an error when the product creation fails.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: Error should be returned due to the insertion failure.

#### 9. Test GetAllProducts Success
- **Test:** `TestGetAllProducts`
- **Description:** Tests the successful retrieval of all products from the repository.
- **Expected Behavior:** The function should return all the products with no error.
- **Test Case:** 
    - Two products with IDs `prod123` and `prod456`.
    - Expected behavior: The result should contain both products with correct IDs.

#### 10. Test GetProductsByUserID Success
- **Test:** `TestGetProductsByUserID`
- **Description:** Tests the successful retrieval of products for a specific user.
- **Expected Behavior:** The function should return all products for the specified user.
- **Test Case:** 
    - Two products for `UserID: 1` with IDs `prod123` and `prod456`.
    - Expected behavior: The result should contain both products associated with the given user ID.

#### 11. Test UpdateProduct Success
- **Test:** `TestUpdateProduct`
- **Description:** Tests the successful update of a product in the repository.
- **Expected Behavior:** The function should return no error when the product is successfully updated.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: No error should be returned.

#### 12. Test DeleteProduct Success
- **Test:** `TestDeleteProduct`
- **Description:** Tests the successful deletion of a product from the repository.
- **Expected Behavior:** The function should return no error when the product is successfully deleted.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: No error should be returned.

#### 13. Test FindProductByUserAndId Success
- **Test:** `TestFindProductByUserAndId`
- **Description:** Tests the successful retrieval of a product by user and product ID.
- **Expected Behavior:** The function should return the product if found.
- **Test Case:** 
    - Product with `UserID: 1` and `ProductID: "prod123"`.
    - Expected behavior: The product should be returned with the correct product ID.

---

## Config: Unit Tests

#### 1. Test ConnectDB Success (NEW) 🚀
- **Test:** `TestConnectDB_Success`
- **Description:** Verifies successful connection to the MongoDB database using a valid URI.
- **Expected Behavior:** The function should establish a connection, and `client.Ping()` should succeed.
- **Test Case:**
  - Set environment variable `MONGO_URI` to `"mongodb://localhost:27017"`.
  - Ensure the database client is not nil and no errors occur during connection, ping, and disconnection.

#### 2. Test ConnectDB Default URI (NEW) 🚀
- **Test:** `TestConnectDB_DefaultURI`
- **Description:** Verifies connection to MongoDB when the `MONGO_URI` environment variable is not set (default URI).
- **Expected Behavior:** The function should establish a connection to the default URI.
- **Test Case:**
  - Unset `MONGO_URI` environment variable.
  - Ensure the database client is not nil and no errors occur during connection, ping, and disconnection.

#### 3. Test GetCollection Success (NEW) 🚀
- **Test:** `TestGetCollection_Success`
- **Description:** Verifies that the correct collection is returned when the `MONGO_URI` is set and a connection is established.
- **Expected Behavior:** The function should return the correct collection, and the collection name should match the input.
- **Test Case:**
  - Set `MONGO_URI` to `"mongodb://localhost:27017"`.
  - Retrieve the collection named `"testCollection"`, ensure no error occurs, and that the collection name and database name match.

#### 4. Test GetCollection DB Not Nil (NEW) 🚀
- **Test:** `TestGetCollection_DBNotNil`
- **Description:** Verifies that when the database client is already connected, the function retrieves the collection without establishing a new connection.
- **Expected Behavior:** The function should return the correct collection and database.
- **Test Case:**
  - Set `MONGO_URI` to `"mongodb://localhost:27017"`.
  - Use an already established DB connection and retrieve the collection `"testCollection2"`.
  - Ensure no error occurs and that the collection name and database name match.

#### 5. Test ConnectDB Connection Failure (NEW) 🚀
- **Test:** `TestConnectDB_ConnectionFailure`
- **Description:** Tests the failure scenario when the MongoDB URI is invalid.
- **Expected Behavior:** The function should return an error indicating the connection failure.

#### 6. Test GetCollection ConnectDB Failure (NEW) 🚀
- **Test:** `TestGetCollection_ConnectDBFailure`
- **Description:** Tests the scenario when the MongoDB URI is invalid, causing `ConnectDB` to fail.
- **Expected Behavior:** The function should return an error and nil collection.

#### 7. Test GetAWSClientInstance Failure (NEW) 🚀
- **Test:** `TestGetAWSClientInstance_Failure`
- **Description:** Verifies failure when loading AWS config results in an error.
- **Expected Behavior:** The function should return an error and nil client.
- **Test Case:**
  - Mock the `LoadDefaultConfig` method to return an error.
  - Ensure the AWS client is nil and an error is returned.

#### 8. Test GetAWSClientInstance Success (NEW) 🚀
- **Test:** `TestGetAWSClientInstance_Success`
- **Description:** Verifies successful loading of AWS config and creation of AWS client.
- **Expected Behavior:** The function should successfully create and return an AWS client.
- **Test Case:**
  - Mock the `LoadDefaultConfig` method to return a valid `aws.Config`.
  - Ensure the AWS client is created with the correct region.

#### 9. Test GetAWSClientInstance Singleton (NEW) 🚀
- **Test:** `TestGetAWSClientInstance_Singleton`
- **Description:** Verifies that `GetAWSClientInstance` returns the same AWS client instance across multiple calls.
- **Expected Behavior:** The function should always return the same AWS client instance.
- **Test Case:**
  - Mock the `LoadDefaultConfig` method to return a valid `aws.Config`.
  - Call `GetAWSClientInstance` twice and verify that both calls return the same client instance.

#### 10. Test LoadAWSConfig Error (NEW) 🚀
- **Test:** `TestLoadAWSConfig_Error`
- **Description:** Tests the failure scenario when loading the AWS config results in an error.
- **Expected Behavior:** The function should return an error and nil client.

#### 11. Test LoadAWSConfig Success (NEW) 🚀
- **Test:** `TestLoadAWSConfig_Success`
- **Description:** Verifies successful loading of AWS config and client creation.
- **Expected Behavior:** The function should successfully load the AWS config and create the AWS client.
  
---
## Custom Errors: Unit Tests

#### 1. Test Custom Error (NEW) 🚀
- **Test:** `TestCustomError`
- **Description:** Verifies the behavior of custom errors with a non-nil cause and a nil cause.
- **Expected Behavior:** The function should return the correct message, status code, cause, and formatted error message.
- **Test Case:** 
    - Error with message: `"test message"`, status code: `http.StatusInternalServerError`, cause: `"test cause"`.
    - Expected behavior: 
        - Message: `"test message"`.
        - Status code: `500 Internal Server Error`.
        - Cause: `"test cause"`.
        - Error message: `"Error: test message, Cause: test cause"`.
    - Error with `nil` cause:
        - Message: `"test message"`, status code: `http.StatusInternalServerError`, cause: `nil`.
        - Expected behavior: Error message: `"Error: test message, Cause: <nil>"`.

#### 2. Test Not Found Error (NEW) 🚀
- **Test:** `TestNotFoundError`
- **Description:** Verifies the behavior of a `NotFoundError` with a non-nil cause.
- **Expected Behavior:** The function should return the correct message, status code, cause, and formatted error message.
- **Test Case:** 
    - Error with message: `"not found"`, cause: `"test cause"`.
    - Expected behavior:
        - Message: `"not found"`.
        - Status code: `404 Not Found`.
        - Cause: `"test cause"`.
        - Error message: `"Error: not found, Cause: test cause"`.

#### 3. Test Database Error (NEW) 🚀
- **Test:** `TestDatabaseError`
- **Description:** Verifies the behavior of a `DatabaseError` with a non-nil cause.
- **Expected Behavior:** The function should return the correct message, status code, cause, and formatted error message.
- **Test Case:** 
    - Error with message: `"database error"`, cause: `"test cause"`.
    - Expected behavior:
        - Message: `"database error"`.
        - Status code: `500 Internal Server Error`.
        - Cause: `"test cause"`.
        - Error message: `"Error: database error, Cause: test cause"`.

#### 4. Test S3 Error (NEW) 🚀
- **Test:** `TestS3Error`
- **Description:** Verifies the behavior of an `S3Error` with a non-nil cause.
- **Expected Behavior:** The function should return the correct message, status code, cause, and formatted error message.
- **Test Case:** 
    - Error with message: `"s3 error"`, cause: `"test cause"`.
    - Expected behavior:
        - Message: `"s3 error"`.
        - Status code: `500 Internal Server Error`.
        - Cause: `"test cause"`.
        - Error message: `"Error: s3 error, Cause: test cause"`.

#### 5. Test Bad Request Error (NEW) 🚀
- **Test:** `TestBadRequestError`
- **Description:** Verifies the behavior of a `BadRequestError` with a non-nil cause.
- **Expected Behavior:** The function should return the correct message, status code, cause, and formatted error message.
- **Test Case:** 
    - Error with message: `"bad request"`, cause: `"test cause"`.
    - Expected behavior:
        - Message: `"bad request"`.
        - Status code: `400 Bad Request`.
        - Cause: `"test cause"`.
        - Error message: `"Error: bad request, Cause: test cause"`.

#### 6. Test Error With Nil Cause (NEW) 🚀
- **Test:** `TestErrorWithNilCause`
- **Description:** Verifies the behavior of an error with a `nil` cause.
- **Expected Behavior:** The function should return the correct formatted error message when the cause is `nil`.
- **Test Case:** 
    - Error with message: `"bad request"`, cause: `nil`.
    - Expected behavior:
        - Error message: `"Error: bad request, Cause: <nil>"`.
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
| Package                  | Time    | Coverage             |
| ------------------------ | ------- | -------------------- |
| `web-service/config`     | 20.160s | 93.9% of statements  |
| `web-service/errors`     | 0.079s  | 100.0% of statements |
| `web-service/handler`    | 0.186s  | 71.3% of statements  |
| `web-service/helper`     | 0.162s  | 83.3% of statements  |
| `web-service/model`      | 0.068s  | 100.0% of statements |
| `web-service/repository` | 0.111s  | 27.0% of statements  |
| `web-service/routes`     | 0.092s  | 100.0% of statements |




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

| **Endpoint**      | **Method** | **Description**                                 |
| ----------------- | ---------- | ----------------------------------------------- |
| `/signup`         | POST       | Create a new user, send verification OTP email. |
| `/verifyEmail`    | POST       | Verify user email with OTP code.                |
| `/forgotPassword` | POST       | Initiate a password reset (send reset OTP).     |
| `/updatePassword` | POST       | Update user password (OTP-based).               |
| `/deleteUser`     | POST       | Remove user from the system.                    |
| `/displayUser`    | POST       | Retrieve user information.                      |
| `/login`          | POST       | Authenticate user with email & password.        |
| `/updateName`     | POST       | Update user's display name.                     |
| `/updatePhone`    | POST       | Update user's phone number.                     |

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

| **Method** | **Endpoint**                                                | **Description**                                 | **Usage**                                 |
| ---------- | ----------------------------------------------------------- | ----------------------------------------------- | ----------------------------------------- |
| `GET`      | `/ws?user_id={user_id}`                                     | WebSocket connection for real-time messaging.   | Establishes WebSocket connection.         |
| `POST`     | `/send`                                                     | Sends a message from one user to another.       | Accepts a JSON payload to send a message. |
| `GET`      | `/messages?sender_id={sender_id}&receiver_id={receiver_id}` | Retrieves messages exchanged between two users. | Fetches messages by sender and receiver.  |
| `GET`      | `/users`                                                    | Retrieves all registered users in the system.   | Fetches a list of all users.              |

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

## Messaging Unit Tests

This module contains unit tests for the **messaging service**, ensuring correct behavior of message-related operations using **Testify** for mocking and assertions.

### **Testing Framework**
- **Go Testing Package (`testing`)** – Standard testing framework in Go.
- **Testify (`github.com/stretchr/testify`)**:
  - `mock` – Used to create a mock message repository.
  - `assert` – Used for validating expected and actual outcomes.

### **Mocking Message Repository**
A **MockMessageRepository** is created to simulate database interactions without an actual database.

### **Tested Functions**

1. **`SaveMessage`**  
   - **Test:** Ensures a message is successfully saved.  
   - **Mocked Call:** `SaveMessage(models.Message)`  
   - **Expected Behavior:** No errors returned.  

2. **`SaveMessageError`**  
   - **Test:** Handles database error while saving a message.  
   - **Mocked Call:** `SaveMessage(models.Message)`  
   - **Expected Behavior:** Returns an expected error.  

3. **`GetLatestMessages`**  
   - **Test:** Retrieves the latest messages with a limit.  
   - **Mocked Call:** `GetLatestMessages(limit int)`  
   - **Expected Behavior:** Returns an expected list of messages.  

4. **`GetLatestMessagesError`**  
   - **Test:** Handles error while fetching latest messages.  
   - **Mocked Call:** `GetLatestMessages(limit int)`  
   - **Expected Behavior:** Returns an expected error.  

5. **`MarkMessageAsRead`**  
   - **Test:** Marks a message as read by its ID.  
   - **Mocked Call:** `MarkMessageAsRead(messageID string)`  
   - **Expected Behavior:** No errors returned.  

6. **`MarkMessageAsReadError`**  
   - **Test:** Handles error while marking a message as read.  
   - **Mocked Call:** `MarkMessageAsRead(messageID string)`  
   - **Expected Behavior:** Returns an expected error.  

7. **`GetUnreadMessages`**  
   - **Test:** Fetches unread messages for a user.  
   - **Mocked Call:** `GetUnreadMessages(userID uint)`  
   - **Expected Behavior:** Returns unread messages.  

8. **`GetUnreadMessagesError`**  
   - **Test:** Handles error while retrieving unread messages.  
   - **Mocked Call:** `GetUnreadMessages(userID uint)`  
   - **Expected Behavior:** Returns an expected error.  

9. **`GetConversation`**  
   - **Test:** Retrieves conversation history for a user.  
   - **Mocked Call:** `GetConversation(user1ID, user2ID uint)`  
   - **Expected Behavior:** Returns a list of exchanged messages.  

10. **`GetConversationError`**  
    - **Test:** Handles error while fetching conversation history.  
    - **Mocked Call:** `GetConversation(user1ID, user2ID uint)`  
    - **Expected Behavior:** Returns an expected error.  

11. **`GetAllUsers`**  
    - **Test:** Retrieves a list of all users.  
    - **Mocked Call:** `GetAllUsers()`  
    - **Expected Behavior:** Returns a list of users.  

12. **`GetAllUsersError`**  
    - **Test:** Handles error while fetching user list.  
    - **Mocked Call:** `GetAllUsers()`  
    - **Expected Behavior:** Returns an expected error.  

13. **`NewUserRepository`**  
    - **Test:** Ensures a new UserRepository is instantiated correctly.  
    - **Mocked Call:** `NewUserRepository(*sql.DB)`  
    - **Expected Behavior:** Returns a valid repository instance.  

14. **`NewMessageRepository`**  
    - **Test:** Ensures a new MessageRepository is instantiated correctly.  
    - **Mocked Call:** `NewMessageRepository(*sql.DB)`  
    - **Expected Behavior:** Returns a valid repository instance.  

### **Running Tests**
To execute the test suite, use:
```bash
go test ./...
```

---
# Frontend Messaging System 
---

## Overview
The Frontend Messaging System allows users to send and receive messages in real-time using WebSockets. The system supports user selection, message history fetching, live message updates, and a typing indicator.

## Features
| Feature                       | Function Name                                                       | Description                                           |
|-------------------------------|---------------------------------------------------------------------|-------------------------------------------------------|
| **User Selection**            | `useFetchUsers(userId)`                                             | Fetches the list of users for selection.              |
| **Real-time Messaging**       | `useWebSocket(userId, handleMessageReceived)`                       | Manages WebSocket connection for live message updates.|
| **Message History**           | `useFetchMessages(userId, selectedUser, setMessages)`               | Fetches previously exchanged messages.                |
| **Typing Indicator**          | `useTypingIndicator(setInput)`                                      | Detects and notifies when a user is typing.           |
| **Send Messages**             | `useSendMessage(userId, selectedUser, users, ws, input, setInput)`  | Sends messages via WebSocket.                         |
| **Receive Messages**          | `handleMessageReceived(message)`                                    | Updates message state when a new message arrives.     |
| **User Authentication**       | `getCurrentUserId()`                                                | Retrieves the logged-in user’s ID.                    |

## Dependencies
- React
- Custom hooks:
  - `useWebSocket`: Manages WebSocket connection.
  - `useFetchMessages`: Fetches previous messages.
  - `useTypingIndicator`: Detects typing activity.
  - `useFetchUsers`: Fetches available users.
  - `useSendMessage`: Handles sending messages.
- Utility function:
  - `getCurrentUserId`: Retrieves the logged-in user’s ID.
- Custom Component:
  - `MessageDisplay`: Displays chat messages.

## Component Breakdown
### **Chat Component**
#### **State Variables:**
- `messages`: Stores chat messages.
- `input`: Tracks user input.
- `selectedUser`: Stores the currently selected chat partner.

#### **Hooks Used:**
1. `getCurrentUserId()`: Retrieves the current user's ID.
2. `useFetchUsers(userId)`: Fetches the list of users.
3. `useWebSocket(userId, handleMessageReceived)`: Establishes a WebSocket connection.
4. `useFetchMessages(userId, selectedUser, setMessages)`: Loads past messages.
5. `useTypingIndicator(setInput)`: Detects when the user is typing.
6. `useSendMessage(userId, selectedUser, users, ws, input, setInput)`: Sends messages over WebSockets.

#### **Event Handlers:**
- `handleMessageReceived(message)`: Updates message state when a new message arrives.
- `handleTyping(event)`: Updates input and triggers typing indication.

## Required Function Implementations
To complete the system, provide implementations for:
1. `useWebSocket(userId, handleMessageReceived)`: Manages WebSocket connection and listens for messages.
2. `useFetchMessages(userId, selectedUser, setMessages)`: Fetches messages from the database.
3. `useTypingIndicator(setInput)`: Detects when the user is typing and notifies the chat partner.
4. `useFetchUsers(userId)`: Retrieves the list of available users.
5. `useSendMessage(userId, selectedUser, users, ws, input, setInput)`: Handles sending messages via WebSocket.

## UI Structure
### **User List (Sidebar)**
- Displays all available users.
- Highlights the selected user.

### **Chat Window**
- Shows chat history.
- Allows users to send messages.
- Displays typing indicator when the chat partner is typing.

## Future Enhancements
- **Read Receipts:** Indicate when a message has been read.
- **Message Timestamps:** Show when messages were sent.
- **Group Chats:** Allow conversations with multiple users.



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

### 17. **User Selection**
- **Test:** Should fetch the list of users for selection.  
- **Expected Behavior:** Calls `useFetchUsers(userId)` and returns a list of users.  

---  

### 18. **Real-time Messaging**  
- **Test:** Should establish a WebSocket connection and receive live updates.  
- **Expected Behavior:** Calls `useWebSocket(userId, handleMessageReceived)` and updates messages in real time.  

---  

### 19. **Message History**  
- **Test:** Should fetch previously exchanged messages when a user is selected.  
- **Expected Behavior:** Calls `useFetchMessages(userId, selectedUser, setMessages)` and retrieves message history.  

---  

### 20. **Typing Indicator**  
- **Test:** Should detect and notify when a user is typing.  
- **Expected Behavior:** Calls `useTypingIndicator(setInput)` and displays a typing indicator when a user types.  

---  

### 21. **Send Messages**  
- **Test:** Should send messages via WebSocket.  
- **Expected Behavior:** Calls `useSendMessage(userId, selectedUser, users, ws, input, setInput)` and updates the chat.  

---  

### 22. **Receive Messages**  
- **Test:** Should update the message state when a new message arrives.  
- **Expected Behavior:** Calls `handleMessageReceived(message)` and updates the UI.  

---  

### 23. **User Authentication**  
- **Test:** Should retrieve the logged-in user’s ID.  
- **Expected Behavior:** Calls `getCurrentUserId()` and returns a valid user ID.  


---

### **End-to-End Testing (Cypress)**  
Cypress was used for **end-to-end testing**, with a successful test for the **login functionality**, as demonstrated in the recorded video.  

#### **Tested Scenario: Login Flow**  
- **Steps Covered:**  
  1. Navigate to the login page.  
  2. Enter valid credentials.  
  3. Click the login button.  
  4. Verify successful authentication and redirection to the dashboard.  

- **Expected Behavior:**  
  - User should be authenticated and redirected to the dashboard upon successful login.  
  - Incorrect credentials should display an error message.  

### **Running Tests**  

#### **Cypress E2E Tests**  
To run Cypress tests, use:  
```bash
npx cypress open
```
### **Video** 

[https://youtu.be/zvNfPwMtuNg](https://youtu.be/zvNfPwMtuNg)

---

