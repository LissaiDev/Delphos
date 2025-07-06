# Delphos Server API

A robust system monitoring application, inspired by the Oracle of Delphi, providing real-time insights into your system's health. This API exposes endpoints for retrieving system metrics and integrates with a real-time streaming service for up-to-the-second data updates.

## Key Features

* **Real-time Monitoring:**  Delivers system metrics via Server-Sent Events (SSE), enabling instantaneous updates.
* **Comprehensive Data:** Exposes detailed information on CPU usage, memory, disk space, and network activity.
* **REST API:**  Allows retrieval of all data through a standard REST API for non-real-time access.
* **Scalable Architecture:** Designed to handle various system sizes and load conditions effectively.
* **Robust Logging:**  Provides detailed logging for debugging and troubleshooting, including timestamps, request/response details, error messages, and more.
* **Rate Limiting:** Implemented to prevent abuse and maintain system stability.
* **CORS:** The API supports cross-origin requests for integration with other applications.
* **Security Headers:** Includes essential security headers like `X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection` and `Referrer-Policy`.

## Getting Started

### Prerequisites

* **Go Development Environment:**  Ensure Go (version 1.24.4 or later) is installed and accessible.
* **GOPATH or Go Modules:** Utilize Go's module management system for dependency tracking.
* **Bun:** Ensure Bun is installed for the client-side development.
* **Next.js:** A modern JavaScript framework for building user interfaces (required by the client)

### Running the Server

1.  **Clone the Repository:**
    ```bash
    git clone <repository-url>
    ```

2.  **Install Dependencies:**
    ```bash
    cd Delphos/srv
    go mod download
    ```

3.  **Run the Server:**
    ```bash
    go run cmd/main.go
    ```

This starts the Delphos Server API on the port specified in `.env` (defaulting to :8080).

### Running the Client

1.  **Clone the Repository:**
    ```bash
    git clone <repository-url>
    ```
2.  **Install Client Dependencies:**
    ```bash
    cd Delphos/client
    bun install
    ```
3.  **Run the Client:**
    ```bash
    bun run dev
    ```

This starts the Next.js development server, which will connect to the Delphos Server API on the correct port for real-time monitoring.

## API Endpoints

*   `/api/stats`: Returns JSON with comprehensive system monitoring data.
*   `/api/stats/sse`:  Provides real-time updates via Server-Sent Events (SSE).

## Data Structure

The server returns data in a structured JSON format, which includes fields for:

*   `host`: System hostname and operating system information
*   `memory`: Total, used, and free memory (RAM), along with swap usage.
*   `cpu`: CPU usage (percentage) and information of each core
*   `disk`: Disk usage for each partition (mount point, type, size, used space).
*   `network`: Network interface statistics (total sent/received bytes).

## Contributing

Pull requests are welcome. Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file for details on how to contribute.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
