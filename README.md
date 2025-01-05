# RideAware API

<i>Train with Focus. Ride with Awareness</i>

RideAware API is the backend service for the RideAware platform, providing endpoints for user authentication and structured workout management.

RideAware is a **comprehensive cycling training platform** designed to help riders stay aware of their performance, progress, and goals.  

Whether you're building a structured training plan, analyzing ride data, or completing workouts indoors, RideAware keeps you connected to every detail of your ride.

## Getting Started

### Prerequisites

Ensure you have the following installed on your system:

- Docker
- Python 3.10 or later
- pip

### Setting Up the Project

1. **Clone the Repository**  

   ```bash
   git clone https://github.com/VeloInnovate/rideaware-api.git
   cd rideaware-api
   ```

2. **Create a Virtual Environment**  
   It is recommended to use a Python virtual environment to isolate dependencies.
   
	```bash
   python3 -m venv .venv
   ```

3. **Activate the Virtual Environment**
   - On Linux/Mac:
     ```bash
     source .venv/bin/activate
     ```
   - On Windows:
     ```cmd
     .venv\Scripts\activate
     ```

4. **Install Requirements**
   Install the required Python packages using pip:
   ```bash
   pip install -r requirements.txt
   ```

### Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory and define the following variables:

```
DATABASE=<your_database_connection_string>
```
- Replace `<your_database_connection_string>` with the URI of your database (e.g., SQLite, PostgreSQL).

### Running with Docker

To run the application in a containerized environment, you can use the provided Dockerfile.

1. **Build the Docker Image**:  

```bash
docker build -t rideaware-api .
```

2. **Run the Container**

```bash
docker run -d -p 5000:5000 --env-file .env rideaware-api
```

The application will be available at http://127.0.0.1:5000.

### Running Tests

To be added.

## Contributing

Contributions are welcome! Please create a pull request or open an issue for any improvements or bug fixes.

## License

This project is licensed under the AGPL-3.0 License.

