## group3-g65-blog-starter-project

# Project Structure Overview

```
group3-g65-blog-starter-project/
├── cmd/
│   └── api/
│       └── main.go                     # Application entry point
├── config/
│   ├── config.go                       # Loads configuration settings
│   └── constants.go                    # Defines application-wide constants
├── delivery/
│   ├── controller/
│   │   ├── auth_controller.go          # Handles authentication endpoints
│   │   ├── blog_controller.go          # Manages blog post endpoints
│   │   ├── user_controller.go          # Manages user-related endpoints
│   │   └── ai_controller.go            # Handles AI-powered features
│   └── route/
│       └── routes.go                   # Centralized route definitions
├── domain/
│   ├── user.go                         # User domain model
│   ├── blog.go                         # Blog post domain model
│   ├── repository.go                   # Repository interface definitions
│   └── usecase.go                      # Use case interface definitions
├── infrastructure/
│   ├── middleware/
│   │   ├── auth_middleware.go          # JWT authentication middleware
│   │   ├── role_middleware.go          # Role-based access control middleware
│   │   ├── logger_middleware.go        # Request logging middleware
│   │   └── recovery_middleware.go      # Error recovery middleware
│   ├── database/
│   │   ├── database.go                 # Database connection and setup
│   │   ├── indexes.go                  # New (index management)
│   │   └── seed/                       # Optional: Seed data scripts
│   ├── auth/
│   │   ├── jwt.go                      # JWT utilities and token struct
│   │   └── password.go                 # Password hashing utilities
│   ├── email/
│   │   └── email_service.go            # Email service for password resets
│   ├── cache/
│   │   └── redis.go                    # Redis client for caching
│   ├── logger.go                       # Logger setup (zap)
│   └── ai/
│       └── ai_client.go                # AI service client integration
├── repository/
│   ├── user_repository.go              # User repository implementation
│   ├── blog_repository.go              # Blog repository implementation
│   └── interaction_repository.go       # Like/dislike tracking repository
│   ├── token_repository.go             # New (for refresh_tokens collection)
│   └── password_reset_repository.go    # New
├── usecase/
│   ├── auth_usecase.go                 # Authentication business logic
│   ├── user_usecase.go                 # User management business logic
│   ├── blog_usecase.go                 # Blog CRUD business logic
│   ├── interaction_usecase.go          # Like/dislike business logic
│   └── ai_usecase.go                   # AI integration business logic
├── utils/
│   ├── pagination.go                   # Pagination helper functions
│   ├── validation.go                   # Input validation helpers
│   ├── context.go                      # Context management utilities
│   ├── Dockerfile                      # Docker container setup
│   ├── .github/
│   │   └── workflows/
│   │       └── ci.yml                  # GitHub Actions CI workflow
│   ├── openapi.yaml                    # OpenAPI specification
│   └── Makefile                        # Build automation script
├── tests/
│   ├── integration/
│   │   ├── auth_test.go                # Authentication integration tests
│   │   └── blog_test.go                # Blog integration tests
│   └── unit/
│       ├── usecase/
│       │   ├── auth_usecase_test.go    # Auth use case unit tests
│       │   └── blog_usecase_test.go    # Blog use case unit tests
│       └── repository/
│           ├── user_repository_test.go # User repository unit tests
│           └── blog_repository_test.go # Blog repository unit tests
├── .gitignore                          # Git ignore rules
├── go.mod                              # Go module definition
├── go.sum                              # Go module checksums
├── render.yaml                         # Render deployment configuration
└── README.md                           # Project documentation
```

# Requirements Mapping

### 1. User Management
- **User Registration & Login**
    - `delivery/controller/auth_controller.go` (API endpoints)
    - `usecase/auth_usecase.go` (business logic)
    - `infrastructure/auth/password.go` (password hashing)
    - `infrastructure/auth/jwt.go` (token generation)
    - `repository/user_repository.go` (database operations)
- **Forgot Password**
    - `infrastructure/email/email_service.go` (email sending)
    - `usecase/auth_usecase.go` (reset logic)
- **Authentication (JWT + Refresh Tokens)**
    - `infrastructure/auth/jwt.go` (token management)
    - `infrastructure/middleware/auth_middleware.go` (JWT validation)
    - `repository/user_repository.go` (refresh token storage)
- **User Promotion/Demotion (Admin)**
    - `delivery/controller/user_controller.go` (API endpoints)
    - `usecase/user_usecase.go` (role change logic)

### 2. Blog Management
- **Blog CRUD Operations**
    - `delivery/controller/blog_controller.go` (API endpoints)
    - `usecase/blog_usecase.go` (business logic)
    - `repository/blog_repository.go` (database operations)
- **Blog Popularity Tracking (Views, Likes, Dislikes)**
    - `repository/interaction_repository.go` (tracking)
    - `usecase/interaction_usecase.go` (logic)
- **Blog Search & Filtration**
    - `repository/blog_repository.go` (search queries)
    - `utils/pagination.go` (pagination)

### 3. AI Integration
- **AI-Generated Content Suggestions**
    - `infrastructure/ai/ai_client.go` (OpenAI/LLM integration)
    - `usecase/ai_usecase.go` (generation logic)
    - `delivery/controller/ai_controller.go` (API endpoints)

### 4. Non-Functional Requirements
- **Security**
    - JWT Authentication: `infrastructure/auth/jwt.go`
    - Password Hashing: `infrastructure/auth/password.go`
    - Role-Based Access Control: `infrastructure/middleware/role_middleware.go`
- **Performance**
    - Caching: `infrastructure/cache/redis.go`
    - Pagination: `utils/pagination.go`
- **Scalability**
    - Database Layer: `infrastructure/database/database.go`
    - Concurrency: Go goroutines (used in use cases)