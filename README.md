# practice-food-delivery-platform

This is a practical exercise to improve my own skills with Kubernetes, Docker, Go, No-SQL, and GitHub Actions.

This project simulates a simplified food delivery platform and is used as a training ground for infrastructure design, clean architecture, service boundaries, and automation.

- [practice-food-delivery-platform](#practice-food-delivery-platform)
  - [Linter](#linter)
  - [Methodology](#methodology)
  - [Actors](#actors)
  - [Services (Go microservices, each exposing a REST API)](#services-go-microservices-each-exposing-a-rest-api)
    - [1. Authentication Service](#1-authentication-service)
    - [2. Customer Service](#2-customer-service)
    - [3. Restaurant Service](#3-restaurant-service)
    - [4. Order Service](#4-order-service)
    - [5. Delivery Service](#5-delivery-service)
    - [6. API Gateway](#6-api-gateway)
  - [Data Storage](#data-storage)
  - [Communication](#communication)
  - [Example Workflow](#example-workflow)
  - [JWT \& Security](#jwt--security)
  - [Concurrency \& Performance](#concurrency--performance)
  - [CI/CD Pipeline (GitHub Actions)](#cicd-pipeline-github-actions)
  - [Planned Improvements](#planned-improvements)

## Linter

It is recommended to use golangci-lint in the project. To install it globaly use the following command:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

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

## Services

### 1. Authentication Service
- Handles identity and access management for all users:
  - Customers
  - Staff (Restaurant Users)
  - Couriers
- Core responsibilities:
  - User registration and authentication
  - JWT issuance (access & refresh tokens)
  - Session management
- Example Endpoints:
  - `POST /auth/customers/register`
  - `POST /auth/restaurants/login`
  - `POST /auth/refresh-token`

### 2. Customer Service
- Manages customer-specific domain data:
  - Customer profiles
  - Delivery addresses
  - Preferences
  - Payment methods
- Example Endpoints:
  - `GET /customers/profile`
  - `PUT /customers/addresses`
  - `GET /customers/payment-methods`

### 3. Restaurant Service
- Manages restaurant domain:
  - Restaurant profiles and details
  - Menu management (CRUD)
  - Operating hours
  - Available locations
- Example Endpoints:
  - `GET /restaurants`
  - `GET /restaurants/{id}/menu`
  - `PUT /restaurants/{id}/menu-items`
  - `GET /restaurants/search`

### 4. Order Service
- Central coordination for order lifecycle:
  - Order creation and validation
  - Order status management
  - Order history
  - Coordinates communication between Restaurant and Delivery services
- Status tracking through order lifecycle
- Example Endpoints:
  - `POST /orders`
  - `GET /orders/{id}`
  - `GET /orders/history`
  - `GET /orders/{id}/status`

### 5. Delivery Service
- Manages delivery operations:
  - Courier management
  - Delivery task assignment
  - Real-time delivery tracking
  - Delivery status updates
- Example Endpoints:
  - `GET /deliveries/active`
  - `PUT /deliveries/{id}/status`
  - `GET /couriers/available`
  - `POST /deliveries/assignments`

### 6. API Gateway
- Single entry point for external clients
- Core responsibilities:
  - Request routing
  - Authentication (JWT validation)
  - Rate limiting
  - CORS management
  - Request/Response logging
  - Basic request validation


---

## Data Storage

- **MongoDB**: Primary store for persistent data (users, restaurants, orders, menus, deliveries)
- **Redis**:
    - Caches popular restaurants and menus
    - Stores short-lived order status for real-time tracking
    - Implements distributed locks to prevent race conditions (e.g., double delivery assignment)

---

## Service Communication

Services communicate through well-defined interfaces following these principles:

1. **Domain Ownership**
  - Each service owns and manages its domain data
  - No direct database access across service boundaries

2. **Event Notification**
  - Services notify relevant changes to interested parties
  - Example: Restaurant Service notifies Order Service about preparation status

3. **Status Management**
  - Order Service acts as the central coordinator for order status
  - Other services report status changes to Order Service

4. **Data Consistency**
  - Each service maintains its data consistency
  - Cross-service consistency through eventual consistency patterns

## Example Workflow

1. Customer authenticates via Authentication Service
2. Customer browses restaurants and menus directly through Restaurant Service
3. Customer places order through Order Service
4. Order Service:
  - Validates menu items with Restaurant Service
  - Creates order record
  - Requests delivery assignment from Delivery Service
5. Restaurant Service updates order preparation status to Order Service
6. Delivery Service updates delivery status to Order Service
7. Customer tracks order status through Order Service


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

