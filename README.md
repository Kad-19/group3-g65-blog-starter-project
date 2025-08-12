# Group3 G65 Blog Starter Project

Starter project for a blog platform using Go. Includes modular architecture, RESTful APIs, authentication, caching, and integrations with AI and cloud services.

## Features
- User authentication (JWT, OAuth)
- Blog CRUD operations
- AI integration
- Caching (in-memory)
- Email notifications
- Image upload (Cloudinary)
- Middleware for auth, cache, and roles
- RESTful API structure
- Modular domain, repository, usecase, and delivery layers

## Project Structure

```
cmd/
    api/
        main.go
config/
    config.go
delivery/
    controller/
        ai_controller.go
        auth_controller.go
        blog_controller.go
        interaction_controller.go
        oauth_controller.go
        user_controller.go
    route/
        routes.go
domain/
    blog.go
    interface_infrastructure.go
    pagination.go
    repository.go
    token.go
    usecase.go
    user.go
infrastructure/
    logger.go
    ai/
        ai_integration.go
    auth/
        jwt.go
        password.go
    cache/
        cache.go
        inmemory.go
    database/
        database.go
    email/
        email.go
    image/
        cloudinary.go
    middleware/
        auth_middleware.go
        cache_middleware.go
        role_middleware.go
repository/
    blog_repository_cache.go
    blog_repository.go
    interaction_repository.go
    passowrd_reset_repository.go
    token_repository.go
    unactive_user_repository.go
    user_repository.go
tmp/
    build-errors.log
usecase/
    ai_usecase.go
    auth_usecase.go
    blog_usecase.go
    interaction_usecase.go
    oauth_usecase.go
    user_usecase.go
utils/
    activation_success.html
    token_impl.go
app.log
Dockerfile
go.mod
go.sum
README.md
```

## Getting Started
1. **Clone the repository:**
    ```bash
    git clone https://github.com/Kad-19/group3-g65-blog-starter-project.git
    cd group3-g65-blog-starter-project
    ```
2. **Configure environment variables:**
    - Edit the `.env` file with your settings (DB, JWT secret, Cloudinary, etc).

## Environment Variables (.env)

The `.env` file contains sensitive configuration for the application. You must set these values before running the project. Example variables:

```
MONGODB_DB           # MongoDB database name
MONGODB_URI          # MongoDB connection URI

ACCESS_TOKEN_SECRET  # JWT access token secret
REFRESH_TOKEN_SECRET # JWT refresh token secret
ACCESS_TOKEN_EXPIRY  # Access token expiry (e.g., 15m)
REFRESH_TOKEN_EXPIRY # Refresh token expiry (e.g., 168h)

GOOGLE_OAUTH_CLIENT_ID     # Google OAuth client ID
GOOGLE_OAUTH_CLIENT_SECRET # Google OAuth client secret
OAUTH_STATE_STRING         # OAuth state string

GEMINI_AI_API_KEY    # Gemini AI API key

SMTP_HOST            # SMTP server host
SMTP_PORT            # SMTP server port
SMTP_USER            # SMTP username
SMTP_PASS            # SMTP password
FROM_EMAIL           # Sender email address

CLOUD_NAME           # Cloudinary cloud name
API_KEY              # Cloudinary API key
API_SECERT           # Cloudinary API secret
```

**Note:** Never commit your `.env` file to public repositories. Keep your secrets safe.
3. **Install dependencies:**
    ```bash
    go mod tidy
    ```
4. **Run the application:**
    ```bash
    go run cmd/api/main.go
    ```

## Live Reloading with Air
This project supports live reloading using [Air](https://github.com/cosmtrek/air).

### Install Air
```bash
go install github.com/cosmtrek/air@latest
```

### Run with Air
```bash
air
```

Air will automatically reload the server when you change Go files.

## Docker Support
To run with Docker:
```bash
docker build -t blog-app .
docker run -p 8080:8080 --env-file .env blog-app
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the MIT License.