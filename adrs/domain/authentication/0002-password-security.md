# 2. Password Security

## Status

Accepted

Date: 2025-06-18

## Context

The food delivery platform needs a secure and robust password management system that:

- Protects user credentials across all user types (customers, staff, couriers)
- Follows industry best practices for password storage
- Balances security with performance
- Ensures compliance with security standards
- Provides secure password verification

Current challenges:

- Need for secure password storage
- Protection against common password attacks
- Maintaining performance with secure hashing
- Future-proofing against evolving security threats

## Decision

We will implement password security using the following approach:

1. **Password Hashing**
    - Use Argon2id as the primary hashing algorithm
    - Configuration parameters:
      ```
      Memory: 64MB
      Iterations: 3
      Parallelism: 2
      SaltLength: 16 bytes
      KeyLength: 32 bytes
      ```

2. **Password Requirements**
    - Minimum length: 8 characters
    - Character requirements:
        - At least one uppercase letter
        - At least one lowercase letter
        - At least one number
        - At least one special character
    - Maximum length: 128 characters
    - No common password patterns are allowed

3. **Password Storage**
    - Store only hashed passwords
    - Never log or transmit plain-text passwords
    - Use dedicated password fields with `writeOnly: true` in API specs

4. **Password Verification**
    - Constant-time comparison for hash verification
    - Rate limiting on verification attempts
    - Secure error messages (no information disclosure)

## Consequences

### Positive

- Strong protection against password attacks
- Industry-standard security compliance
- Future-proof hashing algorithm
- Clear password requirements for users
- Protection against rainbow table attacks

### Negative

- Higher computational cost for password operations
- Increased registration/login latency
- More complex password validation logic
- Need for a password upgrade mechanism if parameters change

### Neutral

- Regular security audits required
- Need for password strength metrics
- User education on password requirements

## Implementation Notes

1. **Password Hashing Flow**
   ```
   1. Validate password requirements
   2. Generate random salt
   3. Apply Argon2id with configured parameters
   4. Format: $argon2id$v=19$m=65536,t=3,p=2$[salt]$[hash]
   ```

2. **Password Verification Process**
   ```
   1. Extract parameters from stored hash
   2. Hash input password with same parameters
   3. Constant-time comparison of hashes
   ```

3. **Security Considerations**
    - Use secure random number generator for salt
    - Clear password buffers after use
    - Implement rate limiting per IP/user
    - Log failed attempts (without passwords)

4. **API Response Codes**
    - 400: Invalid password format
    - 401: Invalid credentials
    - Never indicate which part of credentials failed

## Related Documents

- [ADR-1 Authentication Strategy](./0001-authentication-strategy.md)
- [ADR-4 User Types and Roles Structure](./../../global/0004-user-types-and-roles-structure.md)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- Password Security Standards Document (to be created)

## Contributors

- Àlex Grau Roca