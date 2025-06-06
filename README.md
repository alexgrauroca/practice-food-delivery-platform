# practice-food-delivery-platform

This is a practical exercise to improve my own skills with Kubernetes, Docker, Go, No-SQL and GitHub Actions.

This is the initial use case designed by GitHub Copilot. This is an starting point and will be improved over time by myself.

## Methodology

I will follow a kind of Lean Agile methodology to develop this project. The main idea is to do small and quick iterations, aiming to reduce the risk of failure and to improve the knowledge.

## Food Delivery Platform - Use Case

### Actors

- **Customer**: Browses restaurants, places orders, tracks deliveries.
- **Restaurant**: Manages menu, receives and prepares orders.
- **Delivery Person**: Picks up orders from restaurants and delivers to customers.

---

### Services (Go microservices, each with its own REST API)

1. **Customer Service**
    - User registration/login
    - Profile management
    - Order history

2. **Restaurant Service**
    - Restaurant registration/login
    - Menu management (CRUD for dishes)
    - Order management (receive, update status)

3. **Order Service**
    - Order placement (by customer)
    - Order status tracking (pending, preparing, out for delivery, delivered)
    - Assign delivery person

4. **Delivery Service**
    - Delivery person registration/login
    - View assigned deliveries
    - Update delivery status

---

### Data Storage

- **MongoDB**: Stores users, restaurants, menus, orders, and delivery info.
- **Redis**:
    - Caches popular menus and restaurant data for fast access.
    - Implements distributed locks for order processing (to avoid double assignment).
    - Stores short-lived order status updates for real-time tracking.

---

### Communication

- All services expose REST APIs.
- Internal service-to-service communication also via REST (e.g., Order Service calls Restaurant Service to check menu availability).

---

### Example Workflow

1. **Customer** browses restaurants and menus (Customer Service → Restaurant Service).
2. **Customer** places an order (Customer Service → Order Service).
3. **Order Service** checks menu and availability (calls Restaurant Service), creates order in MongoDB, and uses Redis to lock order processing.
4. **Order Service** assigns a delivery person (calls Delivery Service).
5. **Restaurant** updates order status (preparing, ready).
6. **Delivery Person** picks up and delivers order, updating status via Delivery Service.
7. **Customer** tracks order status in real time (Order Service uses Redis for fast status updates).

---

### Concurrency & Performance

- Redis locks prevent race conditions in order assignment.
- Caching menus and restaurant data reduces MongoDB load.
- Services are stateless and horizontally scalable (Kubernetes).

---

### CI/CD Pipeline (GitHub Actions)

- Runs unit, integration, and e2e tests on non-main branches.
- No deployment on main.
