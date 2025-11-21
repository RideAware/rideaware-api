Here's a rewritten README for your Go version:

```markdown
# RideAware API

<i>Train with Focus. Ride with Awareness</i>

RideAware API is the backend service for the RideAware platform, built with Go and PostgreSQL. It provides comprehensive endpoints for user authentication, profile management, and cycling performance tracking.

RideAware is a **comprehensive cycling training platform** designed to help riders stay aware of their performance, progress, and goals.

Whether you're building a structured training plan, analyzing ride data, or completing workouts indoors, RideAware keeps you connected to every detail of your ride.

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL
- **ORM**: GORM
- **Router**: Chi v5
- **Auth**: JWT (Access + Refresh tokens)
- **Email**: Resend
- **Containerization**: Podman/Docker

## Getting Started

### Prerequisites

Ensure you have the following installed on your system:

- Go 1.21 or later
- PostgreSQL 12 or later
- Podman or Docker
- Git

### Setting Up the Project

1. **Clone the Repository**

   ```bash
   git clone https://github.com/VeloInnovate/rideaware-api.git
   cd rideaware-api
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   go mod tidy
   ```

3. **Configure Environment Variables**

   Create a `.env` file in the root directory:

   ```env
   # Database
   PG_USER=postgres
   PG_PASSWORD=your_password
   PG_HOST=localhost
   PG_PORT=5432
   PG_DATABASE=rideaware

   # Server
   PORT=5000

   # Security
   JWT_SECRET_KEY=your-super-secret-key-change-in-production

   # Email Service
   RESEND_API_KEY=re_your_resend_api_key
   SENDER_EMAIL=noreply@rideaware.app
   ```

4. **Set Up the Database**

   Ensure PostgreSQL is running and create the database:

   ```bash
   createdb rideaware
   ```

   GORM will automatically run migrations on startup.

### Running Locally

**Option 1: Direct Execution**

```bash
go run main.go
```

The API will be available at `http://localhost:5000`

**Option 2: Build and Run Binary**

```bash
go build -o server .
./server
```

### Running with Podman/Docker

#### Quick Start (with build script)

```bash
chmod +x build.sh
./build.sh --run
```

This will build the image and start a container.

#### Manual Build

1. **Build the Image**

   ```bash
   podman build -t rideaware:latest .
   ```

2. **Run the Container**

   ```bash
   podman run -d \
     --name rideaware-api \
     -p 5000:5000 \
     --env-file .env \
     rideaware:latest
   ```

3. **View Logs**

   ```bash
   podman logs -f rideaware-api
   ```

The API will be available at `http://localhost:5000`

## API Documentation

### Health Check

```bash
GET /health
```

Response: `OK`

### Authentication

#### Sign Up

```bash
POST /api/signup
Content-Type: application/json

{
  "username": "cyclist",
  "password": "SecurePass123",
  "email": "cyclist@example.com",
  "first_name": "John",
  "last_name": "Cyclist"
}
```

Response:
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_in": 900,
  "user_id": 1,
  "username": "cyclist",
  "email": "cyclist@example.com"
}
```

#### Login

```bash
POST /api/login
Content-Type: application/json

{
  "username": "cyclist",
  "password": "SecurePass123"
}
```

#### Request Password Reset

```bash
POST /api/password-reset/request
Content-Type: application/json

{
  "email": "cyclist@example.com"
}
```

#### Confirm Password Reset

```bash
POST /api/password-reset/confirm
Content-Type: application/json

{
  "token": "reset_token_from_email",
  "new_password": "NewSecurePass123"
}
```

#### Logout

```bash
POST /api/logout
```

### Protected Routes

All protected routes require the `Authorization: Bearer <access_token>` header.

#### Get User Profile

```bash
GET /api/protected/profile
Authorization: Bearer <access_token>
```

## Testing

Run the test suite:

```bash
chmod +x test-api.sh
./test-api.sh
```

This will test:
- User signup
- Login
- Protected routes
- Password reset
- Error handling

## Development

### Running Tests

```bash
./test-api.sh
```

### Building a New Binary

```bash
go build -o server .
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
go vet ./...
```

## Deployment

### Environment Variables for Production

```env
JWT_SECRET_KEY=<generate-secure-random-key>
RESEND_API_KEY=<your-resend-production-key>
PG_PASSWORD=<strong-database-password>
```

### Building for Production

```bash
./build.sh -t prod --no-cache --run
```

Or push to registry:

```bash
./build.sh -t prod -r docker.io/username --push
```

## Troubleshooting

### Database Connection Errors

Ensure PostgreSQL is running and `.env` variables are correct:

```bash
psql -h localhost -U postgres -d rideaware
```

### Port Already in Use

Change the PORT in `.env` or stop the running container:

```bash
podman kill rideaware-api
podman rm rideaware-api
```

### Docker Permission Issues

Run podman with sudo or add your user to the podman group:

```bash
sudo usermod -aG podman $USER
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the AGPL-3.0 License - see the LICENSE file for details.

## Support

For issues, questions, or suggestions, please open an issue on GitHub or contact the development team.