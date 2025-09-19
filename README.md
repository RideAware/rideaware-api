# RideAware API

<i>Train with Focus. Ride with Awareness</i>

RideAware API is the backend service for the RideAware platform, providing endpoints for user authentication and structured workout management.

RideAware is a **comprehensive cycling training platform** designed to help riders stay aware of their performance, progress, and goals.  

Whether you're building a structured training plan, analyzing ride data, or completing workouts indoors, RideAware keeps you connected to every detail of your ride.

## Getting Started

### Prerequisites

Ensure you have the following installed on your system:

- Docker
- Go 1.21 or later
- PostgreSQL (for local development, optional)

### Setting Up the Project

1. **Clone the Repository**  

   ```bash
   git clone https://github.com/rideaware/rideaware-api.git
   cd rideaware-api
   ```

2. **Install Go Dependencies**  
   
	```bash
   go mod tidy
   ```

3. **Build the Application**
   ```bash
   go build -o rideaware-api
   ```

### Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory and define the following variables:

```env
# Database Configuration
PG_HOST=your_postgres_host
PG_PORT=5432
PG_DATABASE=rideaware
PG_USER=your_postgres_user
PG_PASSWORD=your_postgres_password

# Application Configuration
SECRET_KEY=your_secret_key_for_sessions
PORT=8080

# Email Configuration (Optional)
SMTP_SERVER=your_smtp_server
SMTP_PORT=465
SMTP_USER=your_email@domain.com
SMTP_PASSWORD=your_email_password
```

### Running the Application

#### Development Mode

```bash
go run main.go
```

The application will be available at http://localhost:8080.

#### Production Mode

```bash
./rideaware-api
```

### API Endpoints

- `GET /health` - Health check endpoint
- `POST /auth/signup` - User registration
- `POST /auth/login` - User authentication
- `POST /auth/logout` - User logout

### Running with Docker

To run the application in a containerized environment:

1. **Build the Docker Image**:  

```bash
docker build -t rideaware-api .
```

2. **Run the Container**

```bash
docker run -d -p 8080:8080 --env-file .env rideaware-api
```

The application will be available at http://localhost:8080.

### Example Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o rideaware-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/rideaware-api .
CMD ["./rideaware-api"]
```

### Database Migration

The application automatically runs database migrations on startup using GORM's AutoMigrate feature. This will create the necessary tables:

- `users` - User accounts
- `user_profiles` - User profile information

### Running Tests

To run tests:

```bash
go test ./...
```

To run tests with coverage:

```bash
go test -cover ./...
```

### Development

To add new features:

1. Create models in the `models/` directory
2. Add business logic in the `services/` directory  
3. Define API routes in the `routes/` directory
4. Register routes in `main.go`

## Contributing

Contributions are welcome! Please create a pull request or open an issue for any improvements or bug fixes.

## License

This project is licensed under the AGPL-3.0 License.