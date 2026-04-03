# Migration Guide

Step-by-step guide for transitioning from the example User Service to your custom service.

## Phase 1: Understanding the Example

Before making changes, familiarize yourself with the template:

1. Run the example service to understand the flow
2. Examine the API endpoints using the generated HTTP handlers
3. Study the database migrations in `internal/user/migrations/`
4. Review the configuration files in `config/user/`

## Phase 2: Gradual Replacement

### Step 1: Copy the Service Structure

Keep the established structure and replace the domain logic:

```bash
# Copy the user service as a template for your domain
cp -r internal/user internal/your-domain
cp -r cmd/user cmd/your-domain
cp -r config/user config/your-domain
```

### Step 2: Update Proto Definitions

Define your service's proto files in the central proto repository:

1. Update `build/proto_modules` with your service name:
   ```bash
   echo "your_service_name" > build/proto_modules
   ```

2. Refresh proto files:
   ```bash
   make proto-refresh
   ```

### Step 3: Implement Domain Models

Replace the user-specific entities with your domain:

1. Update `internal/your-domain/model/model.go` with your entities
2. Create database migrations for your schema in `internal/your-domain/migrations/`
3. Modify repository methods in `internal/your-domain/repo/` for your data access patterns

### Step 4: Update Business Logic

1. Replace service layer methods in `internal/your-domain/service/` with your business operations
2. Update validation rules and error handling
3. Modify `internal/your-domain/server.go` to handle your specific endpoints

## Phase 3: Configuration and Deployment

### Update Configuration

1. Modify `config/your-domain/*.toml` files with service-specific settings
2. Update database connection strings and service settings
3. Configure routing and middleware options in `route_config.toml`

### Update Build Configuration

1. Update `BINS` in the root `Makefile`:
   ```makefile
   BINS ?= your-domain your-domain_migration
   ```

2. Update `IMAGE_PREFIX` if needed:
   ```makefile
   IMAGE_PREFIX ?= your-prefix-
   ```

3. Update GitHub Actions workflows:
   - Modify `.github/workflows/build_example.yaml` with your service name
   - Update image names and binary references

## What to Keep vs. Replace

### ✅ Keep These Patterns
- Overall project structure and directory layout
- Dependency injection and interface design
- Error handling and validation patterns
- Configuration management approach
- Build system and CI/CD workflows
- Logging and metrics integration

### 🔄 Replace These Components
- Domain-specific models and business logic
- Database schema and migrations
- Proto definitions and generated code
- Service-specific configuration values
- API endpoint implementations
- Domain-specific validation rules

## Exploring Without Full Migration

If you want to explore the example service first without setting up your own protos:

1. Keep the existing `build/proto_modules` content
2. Run `make proto-fetch && make proto-generate` to ensure generated code is up-to-date
3. The example user service will work with the existing proto definitions

