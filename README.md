# Test GridWhiz Microservice with Clean Architecture

A microservice implementation using Clean Architecture with Authentication and User Management services.

## Tech Stack
- **Protocol**: gRPC
- **Database**: MongoDB
- **Programming Language**: Go
- **Authentication**: JWT
- **Architecture**: Clean Architecture

## Features

### Authentication Service (Port 50051)
1. **Register** - Create new user accounts with email validation and password strength checks
2. **Login** - Authenticate users with rate limiting (5 attempts per minute)
3. **Logout** - Revoke JWT tokens

### User Management Service (Port 50052)
1. **List Users** - Get paginated user list with filtering by name and email
2. **Get Profile** - Retrieve user profile by ID
3. **Update Profile** - Update user's own profile with validation
4. **Delete Profile** - Delete user's own profile

## Project Structure
```
microservice/
├── cmd/                  # Application entrypoints
├── internal/             # Private application code
│   ├── auth/             # Authentication domain
│   ├── user/             # User management domain
│   └── pkg/              # Shared packages
├── proto/                # Protocol buffer definitions
├── pb/                   # Generated protobuf code
└── docker-compose.yml    # Docker configuration
```

## Setup Instructions

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Protocol Buffer compiler (protoc)
- Go plugins for protoc:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

### Installation & Running Locally

1. Clone the repository and navigate to the project directory

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Generate protobuf files:
   ```bash
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto

    แล้วย้าย *.pb.go ทั้งหมดไปไว้ folder pb
   ```

4. Start MongoDB:
   ```bash
   docker componse up -d
   ```

5. Run Auth Service:
   ```bash
   go run cmd/auth/main.go
   ```

6. Run User Service (in another terminal):
   ```bash
   go run cmd/user/main.go
   ```

## API Usage Examples

### Using grpcurl

1. Register a new user:
   ```bash
   grpcurl -plaintext -d '{
     "email": "user@example.com",
     "password": "Password123",
     "name": "John Doe"
   }' localhost:50051 auth.AuthService/Register
   ```

2. Login:
   ```bash
   grpcurl -plaintext -d '{
     "email": "user@example.com",
     "password": "Password123"
   }' localhost:50051 auth.AuthService/Login
   ```

3. List users (requires authentication):
   ```bash
   grpcurl -plaintext \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{
       "page": 1,
       "limit": 10
     }' localhost:50052 user.UserService/ListUsers
   ```

4. Get profile:
   ```bash
   grpcurl -plaintext \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{
       "user_id": "USER_ID"
     }' localhost:50052 user.UserService/GetProfile
   ```

## Security Features

1. **Password Security**:
   - Bcrypt hashing
   - Password strength validation (min 8 chars, uppercase, lowercase, number)

2. **JWT Token Management**:
   - Token expiration (24h default)
   - Token revocation on logout
   - Token validation middleware

3. **Rate Limiting**:
   - 5 login attempts per minute per email
   - Prevents brute force attacks

4. **Input Validation**:
   - Email format validation
   - Input sanitization
   - Authorization checks for profile updates/deletes

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| MONGODB_URI | MongoDB connection string | mongodb://root:example@localhost:27017/mydb?authSource=admin |
| JWT_SECRET | Secret key for JWT signing | your-secret-key-here-change-this-in-production |
| JWT_EXPIRY | JWT token expiration time | 24h |
| AUTH_SERVICE_PORT | Auth service gRPC port | 50051 |
| USER_SERVICE_PORT | User service gRPC port | 50052 |
| RATE_LIMIT_ATTEMPTS | Max login attempts | 5 |
| RATE_LIMIT_WINDOW | Rate limit time window | 60s |

## License

This project is licensed under the MIT License.
