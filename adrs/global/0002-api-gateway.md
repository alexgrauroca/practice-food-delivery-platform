# 2. API Gateway

## Status

Accepted

Date: 2025-06-18

## Context

The current microservices architecture exposes multiple endpoints directly to external clients, which creates several
challenges:

- Security concerns with direct access to internal services
- Difficulty in monitoring and controlling access
- Complex client-side service discovery
- No unified way to handle cross-cutting concerns

## Decision

We will implement an API Gateway as the single entry point for all external requests. The API Gateway will:

- Route requests to appropriate internal services
- Implement rate limiting and throttling
- Enable monitoring and analytics
- Handle SSL termination

## Consequences

### Positive

- Single entry point for all external traffic
- Improved monitoring and analytics capabilities
- Simplified client integration
- Better control over API versioning
- Reduced client-side complexity
- Centralized SSL/TLS management
- Services maintain full ownership of their API contracts

### Negative

- Additional network hop for all requests
- Potential single point of failure
- Increased operational complexity
- Need for careful capacity planning
- Additional latency for requests

### Neutral

- Need for additional infrastructure management
- Changes in deployment and testing procedures
- Different skill set requirements for the team
- Modified monitoring and alerting setup

## Implementation Notes

Implementation details will be available at the dedicated tech-stack folder for the API Gateway.

## Related Documents


## Contributors

- Àlex Grau Roca