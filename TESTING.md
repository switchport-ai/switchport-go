# Switchport Go SDK - Testing Guide

This directory contains comprehensive test suites for the Switchport Go SDK.

## Test Coverage

The test suite achieves **91.4% code coverage** and includes:

- **Client Tests** ([switchport/client_test.go](switchport/client_test.go)) - 7 tests
  - Client initialization with API key
  - Environment variable configuration
  - Custom base URL support
  - API key validation (format and presence)
  - Sub-client initialization

- **Prompts Tests** ([switchport/prompts_test.go](switchport/prompts_test.go)) - 11 tests
  - Successful prompt execution with mocked responses
  - Context handling (map and string formats)
  - Variable substitution
  - Error handling (404, 401, 500 responses)
  - Invalid JSON response handling
  - Complex nested variables and context

- **Metrics Tests** ([metrics_test.go](metrics_test.go)) - 14 tests
  - Metric recording with different value types (float, int, bool, string)
  - Context handling (map and string formats)
  - Timestamp formatting (ISO 8601 / RFC3339)
  - Error handling (404, 401, 400 responses)
  - Complex nested context
  - Optional timestamp parameter

- **Types Tests** ([types_test.go](types_test.go)) - 15 tests
  - Context normalization for different input types
  - Timestamp formatting in various scenarios
  - Edge cases (nil, empty, invalid types)
  - JSON serialization verification

- **Errors Tests** ([errors_test.go](errors_test.go)) - 13 tests
  - All error types (AuthenticationError, PromptNotFoundError, MetricNotFoundError, APIError)
  - Error message formatting
  - Type assertions
  - Complex response data handling
  - Special characters in error messages

## Running Tests

### Run all tests
```bash
cd switchport-go
go test ./switchport
```

### Run with verbose output
```bash
go test ./switchport -v
```

### Run with coverage report
```bash
go test ./switchport -cover
```

### Generate detailed coverage report
```bash
# Generate coverage profile
go test ./switchport -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

### Run specific test
```bash
# Run a specific test function
go test ./switchport -run TestNewClient_WithAPIKey

# Run tests matching a pattern
go test ./switchport -run TestPrompts
```

## Test Strategy

The test suite follows these principles:

1. **No External Dependencies**: Tests use mocked HTTP servers (`httptest.NewServer`) instead of making real API calls
2. **Isolation**: Each test is independent and can run in any order
3. **Comprehensive Coverage**: Tests cover success cases, error cases, and edge cases
4. **Mock Responses**: All API responses are mocked, so tests don't rely on actual server data
5. **Type Safety**: Tests verify correct types and type assertions
6. **Real-World Scenarios**: Tests include complex nested data structures similar to production usage

## Test Organization

```
switchport/
├── client.go         → client_test.go
├── prompts.go        → prompts_test.go
├── metrics.go        → metrics_test.go
├── types.go          → types_test.go
├── errors.go         → errors_test.go
└── TESTING.md        (this file)
```

Each `*_test.go` file contains tests for its corresponding source file.

## Writing New Tests

When adding new functionality:

1. Create tests in the corresponding `*_test.go` file
2. Use `httptest.NewServer` for mocking HTTP responses
3. Test both success and error cases
4. Include edge cases (nil, empty, invalid inputs)
5. Verify error types and messages
6. Run `go test ./switchport -cover` to ensure coverage remains high

## Continuous Integration

To integrate with CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run tests
  run: go test ./switchport -v -cover

- name: Check coverage
  run: |
    go test ./switchport -coverprofile=coverage.out
    go tool cover -func=coverage.out
```

## Notes

- Tests do not require `SWITCHPORT_API_KEY` environment variable (mocked)
- Tests do not make real network requests
- Tests run quickly (< 10 seconds total)
- All tests are compatible with Go 1.21+
