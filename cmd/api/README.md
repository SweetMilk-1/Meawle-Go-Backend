# API Server Structure

This directory contains the API server implementation with a clean separation of concerns.

## Structure

- **main.go** - Entry point of the application
- **config/** - Configuration loading and setup
- **di/** - Dependency injection and service initialization
- **routes/** - HTTP routing setup
- **server/** - HTTP server configuration and lifecycle management

## Files Description

### main.go
- Application entry point
- Orchestrates the startup process
- Handles graceful shutdown

### config/config.go
- Loads application configuration
- Sets up logging

### di/dependencies.go
- Initializes all application dependencies
- Manages database connections, repositories, services, handlers, and middleware
- Provides a clean dependency injection container

### routes/router.go
- Defines all HTTP routes
- Configures middleware chains
- Separates public and protected routes

### server/server.go
- Manages HTTP server lifecycle
- Handles graceful shutdown
- Provides server configuration

## Benefits of This Structure

1. **Separation of Concerns** - Each component has a single responsibility
2. **Testability** - Dependencies can be easily mocked for testing
3. **Maintainability** - Clear boundaries between different parts of the application
4. **Scalability** - Easy to add new features without affecting existing code
5. **Reusability** - Components can be reused in different contexts

## Development

To add new features:
1. Add new routes in `routes/router.go`
2. Add new dependencies in `di/dependencies.go`
3. The main file remains clean and focused on orchestration