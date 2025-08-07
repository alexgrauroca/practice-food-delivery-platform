# 6. Test-Driven Development and Behavior-Driven Development Adoption

## Status

Accepted

Date: 2025-07-01

## Context

As our project continues to grow with multiple services, we need to establish a consistent development 
methodology that ensures high code quality, maintainability, and reliability.

We need a standardized approach to software development that promotes quality from the beginning of the development 
lifecycle and helps document behavior through tests.

## Decision

We will adopt a hybrid approach that combines Test-Driven Development (TDD) with Behavior-Driven Development (BDD) 
principles. This integrated approach will be our primary development methodology across all services.

### Integrated TDD/BDD Approach

We will follow the classic TDD cycle, but with tests written in a behavior-oriented style:

1. **Red**: Write a failing test that defines the expected behavior descriptively
2. **Green**: Write the minimal implementation code to make the test pass
3. **Refactor**: Improve the implementation while keeping tests passing

### Key Practices

- Write tests before implementation code
- Express all tests in terms of expected behavior, not implementation details
- Use descriptive test names that document the expected behavior (e.g., "when a user logs in with valid credentials,
  then it should return a token")
- Structure test cases to clearly represent different scenarios and edge cases
- Maintain the TDD cycle discipline across all feature development
- Use behavior-oriented assertions that make the expected outcomes clear
- Organize tests in logical groupings that reflect user scenarios

### Test Naming and Structure

Tests should follow a consistent format that clearly expresses behavior:

- Test names should follow the pattern: `when [condition], then [expected result]`
- Test organization should group related behaviors together
- Test cases should be organized to cover all relevant scenarios, including edge cases
- Tests should serve as documentation of the system's expected behavior

## Consequences

### Positive

- Higher test coverage across the codebase
- Better design through a test-first approach
- Clearer understanding of requirements through tests
- Tests that serve as living documentation
- Faster detection of bugs and regressions
- Improved confidence during refactoring
- More modular and testable code architecture
- Easier collaboration between technical and non-technical team members

### Negative

- Initial development may seem slower
- Additional time needed for test writing and maintenance
- Learning curve for writing behavior-oriented tests
- Risk of over-testing or testing implementation details

### Neutral

- Need for continuous training and reinforcement
- Regular retrospectives to refine the approach
- Need to balance between detailed specifications and maintainable tests

## Implementation Notes

### TDD with Behavior-Oriented Tests Example

```go
// 1. Write a failing behavior-oriented test first
func TestService_LoginCustomer(t *testing.T) {
    tests := []struct {
        name      string
        input     LoginCustomerInput
        mockSetup func(repo *mocks.MockRepository, tokenService *mocks.MockService)
        want      LoginCustomerOutput
        wantErr   error
    }{
        {
            name: "when there is not an active customer with the same email, " +
                "then it should return an invalid credentials error",
            input: LoginCustomerInput{
                Email:    "test@example.com",
                Password: "ValidPassword123",
            },
            mockSetup: func(repo *mocks.MockRepository, _ *mocks.MockService) {
                repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
                    Return(Customer{}, ErrCustomerNotFound)
            },
            want:    LoginCustomerOutput{},
            wantErr: ErrInvalidCredentials,
        },
        // More test cases covering different behaviors...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            repo := mocks.NewMockRepository(ctrl)
            tokenService := mocks.NewMockService(ctrl)
            if tt.mockSetup != nil {
                tt.mockSetup(repo, tokenService)
            }

            service := NewService(logger, repo, tokenService)

            // Act
            got, err := service.LoginCustomer(context.Background(), tt.input)

            // Assert
            assert.ErrorIs(t, err, tt.wantErr)
            assert.Equal(t, tt.want, got)
        })
    }
}

// 2. Implement the minimal code to make it pass
func (s *service) LoginCustomer(ctx context.Context, input LoginCustomerInput) (LoginCustomerOutput, error) {
    // Initial implementation to make the test pass
    _, err := s.repo.FindByEmail(ctx, input.Email)
    if err != nil {
        return LoginCustomerOutput{}, ErrInvalidCredentials
    }
    return LoginCustomerOutput{}, nil
}

// 3. Refactor while adding more test cases and improving implementation
```

## Related Documents

- [What is Test-Driven Development?](https://agilealliance.org/glossary/tdd/)
- [What is Behavior-Driven Development?](https://agilealliance.org/glossary/bdd/)

## Contributors

- Àlex Grau Roca
