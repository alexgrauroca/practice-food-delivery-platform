# 2. Table-Driven Tests Implementation

## Status

Accepted

Date: 2025-07-01

## Context

Our unit testing approach (defined in [ADR-0001](0001-unit-testing-standards-and-mocking.md)) requires efficient 
ways to test multiple scenarios with minimal duplication. Traditional unit testing patterns can lead to excessive 
code repetition when testing various input/output combinations for the same function or method.

Additionally, our behavior-driven approach to testing (defined in 
[ADR-0006](../../global/0006-test-driven-development-and-behavior-driven-development-adoption.md)) needs a 
structured way to represent different behavior scenarios clearly and consistently.

We need to standardize our approach to organizing multiple test cases that:
- Reduces code duplication
- Improves readability and maintainability
- Facilitates comprehensive scenario coverage
- Aligns with our behavior-driven testing approach
- Supports our existing mocking practices

## Decision

We will adopt table-driven tests as our standard approach for organizing unit tests with multiple scenarios. This 
pattern uses a slice of structs to define test cases, which are then executed in a loop.

### Table-Driven Test Structure

1. **Test Case Definition**
   - Define a struct that contains all necessary test case data
   - For complex tests, use a named struct that can be reused across similar tests
   - For simpler tests or unique structs, use anonymous structs inline

2. **Test Case Fields**
   - `name`: Descriptive name following the "when [condition], then [expected result]" pattern. Mandatory for all tests
   - `mocksSetup`: Functions for configuring mocks (optional)
   - For the expected outcome, use `want` as a common naming. Use it as a prefix to specify the type of want, e.g. 
     `wantErr` for errors or `wantStatus` for statuses
   - Use generics when the types are dynamic, but the name can be consistent. For example, the input type can be 
     different depending on the function that is being tested, but they are still `inputs`
   - Use functions for setup-related params, like inserting documents to the database.

3. **Test Execution Loop**
   - Use `t.Run()` with the test case name to create subtests
   - Follow the Arrange-Act-Assert (AAA) pattern inside each test case
   - Keep assertions focused on behavior verification

### Naming Conventions

- Test case names must follow the "when [condition], then [expected result]" pattern
- Group test cases logically by behavior category
- Order test cases from the simplest to the most complex
- Include edge cases and error scenarios

## Consequences

### Positive

- Reduced code duplication through reuse of test setup and assertion logic
- Improved test readability with clear scenario definitions
- Easier addition of new test cases without modifying existing test code
- Better test organization with logical grouping of related scenarios
- Consistent structure across the codebase
- Self-documenting tests that clearly describe system behavior
- Improved test coverage by encouraging comprehensive test cases

### Negative

- Initial complexity when setting up reusable test structures
- Potential learning curve for developers not familiar with the pattern
- Error messages may be less specific without additional context

### Neutral

- Need for careful test case design to balance reusability with clarity
- Need to balance between shared setup and test case independence

## Implementation Notes

### Basic Table-Driven Test Pattern

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    SomeInput
        want     SomeOutput
        wantErr  error
    }{
        {
            name:    "when valid input is provided, then it returns expected output",
            input:   SomeInput{Field1: "value1", Field2: 42},
            want:    SomeOutput{Result: "expected"},
            wantErr: nil,
        },
        {
            name:    "when invalid input is provided, then it returns an error",
            input:   SomeInput{Field1: "", Field2: -1},
            want:    SomeOutput{},
            wantErr: ErrInvalidInput,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange (setup code if needed)

            // Act
            got, err := FunctionUnderTest(tt.input)

            // Assert
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Reusable Test Case Structure with Generics

```go
// Define a reusable test case structure that can be used across similar tests
type serviceTestCase[I, W any] struct {
    name       string
    input      I
    mocksSetup func(repo *mocks.MockRepository, otherDep *mocks.MockDependency)
    want       W
    wantErr    error
}

func TestServiceMethod(t *testing.T) {
    tests := []serviceTestCase[SomeInput, SomeOutput]{
        {
            name: "when condition X exists, then result Y is returned",
            input: SomeInput{...},
            mocksSetup: func(repo *mocks.MockRepository, otherDep *mocks.MockDependency) {
                repo.EXPECT().SomeMethod(gomock.Any(), gomock.Any()).Return(someResult, nil)
                // Configure other mocks as needed
            },
            want: SomeOutput{...},
            wantErr: nil,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            repo := mocks.NewMockRepository(ctrl)
            otherDep := mocks.NewMockDependency(ctrl)
            if tt.mocksSetup != nil {
                tt.mocksSetup(repo, otherDep)
            }

            service := NewService(repo, otherDep)

            // Act
            got, err := service.Method(context.Background(), tt.input)

            // Assert
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Best Practices

1. **Test Independence**

   Each test case should be completely independent and not rely on the execution of previous tests.

2. **Clear Test Case Names**

   Use descriptive names that clearly indicate the scenario being tested and the expected outcome.

3. **Group Related Cases**

   Organize test cases in logical groups, with related scenarios next to each other.

4. **Complete Coverage**

   Include happy paths, error paths, edge cases, and boundary conditions.

5. **Minimal-Shared Setup**

   Keep shared setup code minimal. If extensive setup is needed, consider helper functions.

6. **Focused Assertions**

   Assert only what is relevant to the current test case. Avoid asserting unrelated properties.

7. **Consistent Structure**

   Maintain a consistent structure across all table-driven tests in the codebase.

## Related Documents

- [ADR-0001: Unit Testing Standards and Mocking](0001-unit-testing-standards-and-mocking.md)
- [ADR-0005: Testing Strategy and Types](../../global/0005-testing-strategy-and-types.md)
- [ADR-0006: Test-Driven Development and Behavior-Driven Development Adoption](../../global/0006-test-driven-development-and-behavior-driven-development-adoption.md)
- [Dave Cheney: Prefer table-driven tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Wiki: TableDrivenTests](https://github.com/golang/go/wiki/TableDrivenTests)

## Contributors

- Àlex Grau Roca
