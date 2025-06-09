# practice-food-delivery-platform

This is a practical exercise to improve my own skills with Kubernetes, Docker, Go, No-SQL, and GitHub Actions.

This project simulates a simplified food delivery platform and is used as a training ground for infrastructure design, clean architecture, service boundaries, and automation.

## Methodology

This project follows a Lean Agile development style with quick, incremental iterations to minimize risk and continuously improve the solution.

A **Domain Service Architecture** is used, where each microservice encapsulates a clear and distinct domain. This improves separation of concerns and system maintainability.

Testing follows **TDD (Test-Driven Development)** and **BDD (Behavior-Driven Development)** principles, ensuring behavior and correctness are covered through automation.

Coverage reports will not be generated at this stage. Ideally the pipelines should generate coverage reports with tools like courtney and gocovmerge, to then upload them to an analysis tool like Sonar. Setting up 100% coverage is considered a bad practice, as testing should focus on quality not quantity, but coverage reports can help to detect missed cases scenarios. During the development of this project the Goland coverage tool was used to detect those cases. TDD should help to reduce those cases as much as possible, because we just create the needed code based on the expected behaviour.

---

## Actors

- **Customer**: Browses restaurants, places orders, tracks deliveries.
- **Restaurant User**: Manages restaurant profile, menus, and orders.
- **Delivery Person**: Handles the delivery of orders to customers.

---

## Services (Go microservices, each exposing a REST API)

### 1. Authentication Service
- Handles identity and session management for:
    - **Customers**
    - **Restaurant Users**
- Issues JWTs (access & refresh tokens)
- Manages independent auth flows for each user type
- Example Endpoints:
    - `POST /customers/register`
    - `POST /restaurants/login`
    - `POST /customers/refresh-token`

### 2. Customer Service
- Manages domain data and behavior for customers:
    - Profile (address, preferences, etc.)
    - Order history (linked by user ID from JWT)

### 3. Restaurant Service
- Manages restaurants and their menus:
    - Restaurant profile
    - Menu management (CRUD)
    - Order status updates (preparing, ready)

### 4. Order Service
- Handles full order lifecycle:
    - Order placement
    - Order status tracking
    - Delivery assignment
    - Menu validation (via Restaurant Service)

### 5. Delivery Service
- Manages delivery personnel and tasks:
    - Registration and login (via Auth Service)
    - Assigned deliveries
    - Delivery status updates

### 6. API Gateway
- The only publicly exposed entrypoint
- Routes requests to the appropriate internal service
- Optionally validates JWTs
- Manages CORS, rate limiting, logging, etc.

---

## Data Storage

- **MongoDB**: Primary store for persistent data (users, restaurants, orders, menus, deliveries)
- **Redis**:
    - Caches popular restaurants and menus
    - Stores short-lived order status for real-time tracking
    - Implements distributed locks to prevent race conditions (e.g., double delivery assignment)

---

## Communication

- Services communicate via REST over HTTP
- API Gateway maps external routes (e.g., `/customers/login`) to the appropriate service
- Services remain stateless and communicate using secure JWT tokens

---

## Example Workflow

1. A **customer** registers via `POST /customers/register` (Auth Service)
2. They login and receive JWT tokens
3. They browse restaurants and menus (`Customer Service → Restaurant Service`)
4. They place an order (`Customer Service → Order Service`)
5. Order is validated, stored, and assigned a delivery person
6. Restaurant updates order status as it prepares
7. Delivery person picks up and updates delivery status
8. Customer tracks delivery in real-time (via Redis-backed status from Order Service)

---

## JWT & Security

- JWTs include claims such as `sub` (user ID), `role` (customer, restaurant), and `user_type`
- Tokens are verified by services or the API Gateway
- Refresh tokens allow session extension
- Future support for MFA is planned (initially for restaurant users)

---

## Concurrency & Performance

- Redis locks prevent double assignment of deliveries
- Redis caching reduces MongoDB load on frequent reads
- Services are stateless and horizontally scalable via Kubernetes

---

## CI/CD Pipeline (GitHub Actions)

- Runs on non-main branches:
    - Unit tests
    - Integration tests
    - End-to-end tests
- No deployment occurs on `main` to avoid accidental production pushes
- Future plans: add linting, security checks, and deployment to a test cluster

---

## Planned Improvements

- Add MFA support for restaurant users
- Add delivery person role to the Authentication Service
- Add observability tools (e.g., Prometheus, Grafana)
- Implement circuit breakers and retries between services
- Extend order tracking with websockets or server-sent events (SSE)

