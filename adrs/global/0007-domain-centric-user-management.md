# 7. Domain-Centric User Management

## Status

Accepted

Date: 2025-08-07

## Context

Our platform needs to manage different types of users (Customers, Staff, Couriers) who share some common data (like 
authentication credentials) but also have domain-specific data. Following the Principle of Least Surprise, users expect 
to manage their profile data in a single, cohesive interface, regardless of our internal service boundaries.

We need to establish:
- Clear ownership of user management flows
- Consistent API patterns for user operations
- Clean service boundaries
- Simple and intuitive API for clients
- Proper encapsulation of internal service communication

## Decision

We will implement a domain-centric approach to user management where:

1. **Domain Services Own User Journeys**
   - Each domain service (Customer/Staff/Courier) owns its complete user journey
   - Domain services expose cohesive APIs that handle all user operations
   - Internal service communication is hidden from API consumers

2. **Service Coordination Pattern**
   - Domain services coordinate with Authentication Service internally
   - Authentication Service focuses purely on auth mechanisms
   - Domain services handle their specific user data

3. **API Design**
   - Single registration endpoint per user type in its domain service
   - Complete user profile management within domain services
   - Unified response including both domain and auth data

Example Flow:
1. Client -> Domain Service: Register user (complete data)
2. Domain Service -> Auth Service: Handle auth details
3. Domain Service: Store domain-specific profile
4. Domain Service -> Client: Return complete profile

## Consequences

### Positive

- Clean, intuitive APIs that align with user expectations
- Clear ownership of user journeys
- Proper encapsulation of internal service details
- Easy to extend domain-specific user features
- Consistent pattern across all user types
- Better user experience with a unified registration process

### Negative

- Additional coordination logic in domain services
- Need for careful error handling in cross-service communication
- More complex integration testing scenarios
- Potential for data consistency issues between services

### Neutral

- Domain services need to maintain auth service client code
- Need for clear documentation of internal service interactions
- Regular review of domain boundaries may be needed

## Implementation Notes

### API Design

1. **Registration Endpoints**
   ```
   POST /customers/register
   POST /staff/register
   POST /couriers/register
   ```

2. **Request Structure**
   - Keep authentication and domain data at the root level
   - Use consistent field naming across all user types
   - Include all required data in a single request

3. **Response Structure**
   - Return a complete user profile
   - Include authentication tokens
   - Provide clear success/error indicators

### Error Handling

1. **Failure Scenarios**
   - Authentication service unavailable
   - Duplicate user registration
   - Invalid data validation
   - Partial registration failure

2. **Recovery Mechanisms**
   - Implement compensation logic for failed registrations
   - Clear cleanup of partial registrations
   - Maintain an audit trail of registration attempts

### Service Communication

1. **Authentication Flow**
   - Synchronous communication for registration
   - Validate credentials before profile creation
   - Handle auth service errors gracefully

2. **Data Consistency**
   - Use transactional approaches where possible
   - Implement retry mechanisms for temporary failures
   - Regular data reconciliation between services

## Related Documents

- [REST API Request Formatting](../tech-stack/api/0005-rest-api-request-formatting.md)
- [Interface Design Principles](../tech-stack/go/0002-interface-design-principles.md)

## Contributors

- Àlex Grau Roca
