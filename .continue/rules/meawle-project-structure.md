---
description: Detailed Meawle project structure for efficient navigation and
  consistent development
alwaysApply: true
---

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