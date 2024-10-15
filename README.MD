# Go UDP Logger

A simple UDP logger application written in Go that listens for incoming UDP messages and provides a web interface to view and search logs.

## Features

- Listens on UDP port 514
- Logs received messages with timestamps
- Web interface to view logs
- Search functionality for log messages (case-insensitive)

## Prerequisites

- Go (version 1.15 or higher)
- Docker (optional, for containerization)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/shakson1/go-udp-logger.git
   cd go-udp-logger
   ```

## Build the Application

### Option 1: Build and Run Locally

1. Build the application:

   ```bash
   go build -o logger main.go
   ```

2. Run the application: 

   ```bash
   ./logger
   ```
   The logger will listen on UDP port 514 and start an HTTP server on port 8080.

### Option 2: Build and Run in Docker

1.	Build the Docker image:

    In the project directory, run:
    ```bash
    docker build -t go-udp-logger-app .
    ```

2.	Run the Docker container:
    
    ```bash
    docker run --name udp-logger-container --rm -p 514:514/udp -p 8080:8080 go-udp-logger-app
    ```

## Usage

1.	**Send a test log message** to the logger:

    ```bash
    echo "Test message" | nc -u -w1 127.0.0.1 514
    ```

2.	**Access the web interface** by opening a web browser and navigating to:

    ```
    http://localhost:8080/logs
    ```

3.	**Search logs** using the search input on the page.

## License

This project is free to use. No licensing restrictions apply.

## Owner

- **Email**: shakokakhadze@gmail.com
- **LinkedIn**: [Shalva Kakhadze](https://www.linkedin.com/in/shalva-kakhadze-966234193/)
- **GitHub**: [shakson1](https://github.com/shakson1)

   ```
   This `README.md` now reflects that the project is free to use, and there are no specific license restrictions. Let me know if you need any other updates!
   ```