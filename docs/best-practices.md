# Best Practices

Coding standards and architectural guidelines to maintain when building services from this template.

## Layered Architecture

- **Maintain separation** between handler, service, and repository layers
- **Use dependency injection** for testability and flexibility
- **Define clear interfaces** between layers to enable mocking and testing
- Keep business logic in the service layer, not in handlers

## Error Handling

- Use structured errors with proper error codes from the `internal/errors` package
- Separate public-facing error messages from internal error details
- Implement proper error logging and metrics for observability
- Always wrap errors with context when propagating them up

Example:
```go
if err != nil {
    return nil, errpkg.WrapError(err)
}
```

## Configuration Management

- Use environment-specific configuration files (`default.toml`, `dev.toml`, etc.)
- Keep sensitive data in environment variables, not in config files
- Document all configuration options in the default config file
- Use the Foundation library's config management utilities

## Database Patterns

- Use migrations for all schema changes (never modify schema directly)
- Implement proper connection pooling and timeout configurations
- Add database health checks to monitor connectivity
- Use transactions appropriately for operations that modify multiple tables
- Keep raw SQL queries in the repository layer only

## Testing Patterns

- Write unit tests for business logic in the service layer
- Use mocks for external dependencies (database, external APIs)
- Implement integration tests for critical paths
- Test error conditions and edge cases
- Maintain test coverage above 70% for business logic

## Code Organization

- Keep domain logic in the `internal/` directory (private to the service)
- Use meaningful package names that reflect domain concepts
- Avoid cyclic dependencies between packages
- Group related functionality together in the same package

## API Design

- Follow RESTful principles for HTTP endpoints
- Use proper HTTP status codes and error responses
- Validate all inputs using proto validation rules
- Version your APIs appropriately (`/v1/`, `/v2/`)
- Document your API endpoints in proto files

## Logging and Observability

- Use structured logging with the telemetry interface
- Include request IDs in all log statements
- Log at appropriate levels (debug, info, warn, error)
- Emit metrics for critical operations (latency, error rates)
- Add tracing context for distributed tracing

## Security

- Validate and sanitize all user inputs
- Use parameterized queries to prevent SQL injection
- Implement proper authentication and authorization
- Keep dependencies up to date for security patches
- Never log sensitive information (passwords, tokens, PII)

