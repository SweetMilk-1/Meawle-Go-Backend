---
alwaysApply: true
---

Maintain detailed understanding of Meawle's clean architecture: cmd/api/ contains main entry point, config loading, dependency injection, routing, and server management; internal/ implements business logic with clear separation: config/ for environment variables, database/ for DB connections, handlers/ for HTTP endpoints, middleware/ for auth and other middleware, models/ for data structures, repositories/ for data access, services/ for business logic

Follow Meawle's clean architecture pattern with these specific directories and responsibilities:

**cmd/api/** - Application entry point and orchestration:
- main.go - Application entry, dependency initialization, server startup
- config/ - Configuration loading from environment variables
- di/ - Dependency injection container, service initialization
- routes/ - HTTP routing setup, middleware configuration
- server/ - HTTP server lifecycle management, graceful shutdown

**internal/** - Business logic with clear separation:
- config/ - Environment variable handling, application configuration
- database/ - Database connection management, migrations
- handlers/ - HTTP request handlers, response formatting
- middleware/ - Authentication, authorization, and other middleware
- models/ - Data structures, domain entities
- repositories/ - Data access layer, database operations
- services/ - Business logic, use case implementation

**Key architectural principles:**
- Clear separation between layers (handlers -> services -> repositories)
- Dependency injection for testability
- Single responsibility for each component
- Clean boundaries between HTTP concerns and business logic

**Database schema updates:**
- Table `cats` has foreign key `cat_breed_id` referencing `cat_breeds(id)`
- Field `cat_breed_id` is optional (nullable) in both models and API
- All cat-related endpoints support `cat_breed_id` field in requests and responses

**New cat houses functionality:**
- Table `cat_houses` stores houses with name, capacity, and current occupancy
- Table `cat_house_occupancy` links cats to houses with unique constraint on cat_id
- CRUD operations for cat houses (create, read, update, delete)
- Ability to add/remove cats from houses with capacity validation
- One cat can only be in one house at a time
- Houses cannot be deleted if capacity is reduced below current occupancy