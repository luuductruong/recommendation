# 1. System Design

## Architecture Diagram

```
+------------+         +-----------------+          +-----------------+          +------------------------+
| Client / UI| <-----> |  API Gateway /  | <------> | Product Service | <------> |    PostgreSQL DB       |
| - Sends    |         |  Backend Layer  |          | - Get product   |          | - product              |
|   requests |         | - Routes API    |          | - Recommend     |          | - user_view_history    |
| - Receives |         | - Auth, Logging |          | - Record views  |          | - category_view_history|
|   responses|         |                 |          |                 |          |                        |
+------------+         +-----------------+          +-----------------+          +------------------------+
```

## Components Explanation

1. **Client / UI:**
    - Sends requests to get product details and user recommendations.
    - Receives responses to display or process (though UI is out of scope here).

2. **API Gateway / Backend Layer:**
    - Handles incoming API requests.
    - Responsible for routing requests to appropriate services.
    - Manages authentication, logging, and basic request validations.

3. **Product Service:**
    - Handles business logic for products.
    - Provides APIs for product detail retrieval and recommendation.
    - Records user view history asynchronously.
    - Implements recommendation algorithms based on user viewing history and popular products.

4. **PostgreSQL Database:**
    - Stores product information (`product` table).
    - Stores user viewing history (`user_view_history` table).
    - Stores category-level view counts for popularity tracking (`category_view_history` table).


## Internal service
```
       +----------------------------------  Product Service  ---------------------------------+ 
       |                                                                                      |
       | +---------------------+       +--------------------+       +--------------------+    |
       | |   Application Layer |       |    Domain Layer    |       |   External Layer   |    |
       | |  (gRPC Handler, Req | <---> | (Business Logic,   | <---> | (DB Access, Cache, |    |
       | |   Validation, etc.) |       |  Domain Entities)  |       |  Messaging, etc.)  |    |
       | +----------|----------+       +----------|---------+       +----------|---------+    |
       |            |                             |                            |              |
       |            |                             |                            |              |
       |            v                             v                            v              |
       |   +----------------+          +----------------------+     +----------------------+  |
       |   | gRPC Handler   |          | Recommendation Logic  |    | Database (Postgres)  |  |
       |   +----------------+          +----------------------+     +----------------------+  |
       |            |                             |                            |              |
       |   +-------------------+       +----------------------+     +----------------------+  |
       |   | Request Validation|       | Domain Entities &    |     | Cache (Redis, etc.)  |  |
       |   +-------------------+       | Services             |     +----------------------+  |
       |                               +----------------------+                               |
       +--------------------------------------------------------------------------------------+
```
## Design Flow

1. **User views product:**
    - Client sends a request with `userID` and `productID` to the Product Service via API Gateway.
2. **Product detail retrieval:**
    - Product Service queries the `product` table in PostgreSQL.
    - Returns product details to the client.
3. **Record user view history:**
    - Product Service asynchronously records the view in `user_view_history` (userID, productID, timestamp).
    - Updates `category_view_history` by incrementing view count and updating last view timestamp.
4. **User requests recommendations:**
    - Client sends a request with `userID` and optional parameters.
    - Product Service fetches recent user views and popular products.
    - Combines results to generate a recommendation list.
5. **Return recommendations:**
    - Product Service returns the recommended product list to the client.


# 2. Data Model

### Product Information Schema
- `product` table stores product details:
    - `product_id` (integer, primary key)
    - `name` (text)
    - `price` (numeric)
    - `category_id` (text)

### User Viewing History Schema
- `user_view_history` table records each user's product views with timestamps:
    - `id` (text, primary key)
    - `user_id` (text)
    - `product_id` (integer)
    - `view_at` (timestamp with time zone)
- Index on `(user_id, product_id, view_at DESC)` to optimize querying recent views.

### Pre-computed Recommendations / Association Rules Storage
- Currently, precomputed recommendations are not stored persistently.
- Recommendations are generated on-the-fly using:
    - Recent views from `user_view_history` (via API `GetRecommendationForUser`).
    - Popular products computed from recent views within a time range (last 2 days).
- `category_view_history` stores aggregated category view counts:
    - `id` (text, primary key)
    - `category_id` (text, unique)
    - `total_view` (integer)
    - `last_view_at` (timestamp)
- This supports recommending related products by category popularity.

---

### Proto Message Schemas (example)

```proto
message GetProductDetailReq {
  int64 product_id = 1;  // product to query
  string user_id = 2;    // optional user id for tracking
}

message GetProductDetailResp {
  Product product = 1;
}

message GetRecommendationForUserReq {
  string user_id = 1;    // user to recommend for
  int32 limit = 2;       // max number of products to recommend
}

message GetRecommendationForUserResp {
  repeated SummaryProductView products = 1;  // recommended products, but uses this for testing
}

message Product {
  int64 product_id = 1;
  string name = 2;
  string category_id = 3;
  double price = 4;
}

message SummaryProductView {
  int64 product_id = 1;
  int64 view_count = 2;
  google.protobuf.Timestamp view_at = 3;
}
```

# 3. API Endpoints

### 3.1 Record User's Product View

- **Purpose:** To record when a user views a product.
- **Endpoint:** `GetProductDetail` API (currently `POST /v1/product/detail`, TODO: change to `GET /v1/product/{product_id}/detail`)
- **Input:**
    - `product_id` (int64) — product to view
    - `user_id` (string) — user making the request (detected from header or input)
- **Process:**
    - Retrieve product details from the database.
    - Asynchronously record the user's view into `user_view_history` and update `category_view_history`.
- **Response:** Product details.

---

### 3.2 Get Product Recommendations for User

- **Purpose:** To get a list of recommended products for a user.
- **Endpoint:** `GetRecommendationForUser` API (currently `POST /v1/product/user/recommendation`)
- **Input:**
    - `user_id` (string) — user to generate recommendations for (detected from header)
    - `limit` (int32) — number of recommended products requested (optional, default 10)
- **Logic:**
    - If no user ID is provided, return the most popular products.
    - Retrieve recent viewed products by the user (limit by requested number).
    - If the number of recent views is less than limit, append popular products until limit reached.
    - If no history found, return a random selection of products.
- **Response:** List of recommended product summaries.

---

### Proto Service Example

```proto
service ProductService {
  rpc GetProductDetail(dto.GetProductDetailReq) returns (dto.GetProductDetailResp) {
    option (google.api.http) = {
      post: "/v1/product/detail"
      body: "*"
    };
  }
  rpc GetRecommendationForUser(dto.GetRecommendationForUserReq) returns (dto.GetRecommendationForUserResp) {
    option (google.api.http) = {
      post: "/v1/product/user/recommendation"
      body: "*"
    };
  }
}
```

### 3.3 Summary of API

- `GetProductDetail` serves both as product info API and user view recording trigger.
- `GetRecommendationForUser` returns a tailored recommendation list based on user viewing history and product popularity.
- Both APIs currently use POST with body but planned to be changed to GET with user ID from headers and product ID from URL for RESTful design.


# 4. Technology Choices

### Programming Language & Framework
- **Golang:**
    - High performance, concurrency-friendly, suitable for backend services.
    - Strong standard library and community support for microservices and APIs.

### Database
- **PostgreSQL:**
    - Reliable relational database to store product data, user view history, and category view stats.
    - Supports indexing and advanced querying needed for recommendation logic.

### Data Processing & Algorithms
- Simple custom algorithms for:
    - Tracking user views and product popularity.
    - Generating recommendations based on recent views and popular items.
- Possible future use of lightweight machine learning or collaborative filtering libraries if needed.

### API Design
- gRPC with HTTP/REST gateway using protobuf.
- Easy to evolve and maintain with clear contract between client and service.

### Logging & Monitoring
- Structured logging for debugging and traceability.
- Metrics and alerts for system health and recommendation quality.

---

# 5. Development Roadmap (Overview)

Here are some potential features and improvements for the product recommendation system:

- **Incorporate Purchase History:**  
  Enhance recommendation quality by combining user purchase data with product view history to better understand user preferences.

- **More Advanced Algorithms:**  
  Implement collaborative filtering, matrix factorization, or content-based filtering techniques to provide more personalized recommendations.

- **Real-time Updates:**  
  Move from batch processing to real-time streaming data processing to keep recommendations always up-to-date.

- **A/B Testing Framework:**  
  Apply A/B testing to compare the effectiveness of different recommendation strategies, measuring impact on user engagement and revenue.

- **Advanced Personalization:**  
  Use demographic information, session context, or product attributes to create more relevant recommendations.

- **Scalability Enhancements:**  
  Optimize data storage and retrieval for large volumes of users and products; consider caching layers or distributed data systems.

---

# 6. Emergency Handling

**Scenario:**  
After deploying a new recommendation algorithm, users report that recommendations are irrelevant (e.g., winter coats shown in summer, or unrelated products), negatively impacting user experience and revenue.

**Steps for handling and prevention:**

1. **Quick Root Cause Analysis:**
    - Immediately roll back to the previous stable algorithm version to stop spreading incorrect recommendations.
    - Check logs and monitoring data for errors or anomalies in the new algorithm.
    - Verify input data correctness (e.g., category labels, timestamps).
    - Review recent code or data structure changes that might cause issues.

2. **Mitigation Measures:**
    - Redeploy the proven old algorithm.
    - Notify stakeholders and customer support teams.
    - Temporarily disable personalized recommendations and switch to popular or random product suggestions.

3. **Prevention:**
    - Build a staging environment with simulated data and traffic to test algorithms before production deployment.
    - Set up automated tests to validate recommendation relevance (e.g., category matching, seasonality).
    - Use feature flags to gradually roll out new algorithms, enabling quick toggling.
    - Monitor key metrics (click-through rate, conversion rate, user feedback) and set alerts for anomalies.
    - Maintain detailed version control and clear rollback plans for both models and code.

---


# 7. Bonus: Advanced Considerations for Recommendation System

## 7.1 Scalability & Cold Start Solutions

### Scalability
- Use indexed PostgreSQL queries with proper LIMIT/OFFSET and time range filters.
- Add Redis cache for top products or precomputed results to reduce DB load.
- Horizontally scale the recommendation service with load balancers.

### Cold Start
- **New Users:** Recommend trending or popular products by category.
- **New Products:** Promote in category feeds to get early interactions.
- Consider hybrid models (collaborative + content-based filtering).

---

## 7.2 Real-time vs Batch Update Strategy

- **Batch Updates:** Precompute recommendations periodically (e.g., every 2 hours).
- **Real-Time Updates:** Use message queues (Kafka/NATS) to track new views instantly.
- Balance freshness vs performance by caching precomputed results for popular users/products.

---

## 7.3 Measuring Effectiveness

Key metrics to monitor:
- **CTR (Click-Through Rate):** How many users click recommended items.
- **Conversion Rate:** How many users purchase items from recommendations.
- **Session Duration:** Do recommendations keep users engaged longer?

---

## 7.4 Personalization Beyond Co-occurrence

- Combine product metadata (tags, category, price range) with view history.
- Use user demographics, preferences, and browsing session context.
- Explore user clustering or lightweight ML-based embeddings for deeper personalization.

---

## 7.5 Extended Roadmap

- Implement collaborative filtering using similarity scores.
- Integrate external signals (search keywords, wishlists).
- Run A/B experiments with multiple recommendation strategies.
- Add admin UI to track recommendation effectiveness over time.
