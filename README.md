# HTMX-Go-Test

## Overview

This project demonstrates the integration of HTMX with a Go backend using the Echo framework. It includes API endpoints and web routes for managing posts.

## Getting Started

### Prerequisites

- Go installed on your machine
- Air for live reloading (install via `go install github.com/cosmtrek/air@latest`)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/Sadtawad6217/htmx-go-test.git
    cd htmx-go-test
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

### Running the Project

1. Start the project with Air for live reloading:

    ```bash
    air
    ```

2. Open your web browser and navigate to:

    ```
    http://localhost:8080
    ```

    You should see the application running.


### Notes

- The project uses two separate servers: Echo runs on port `8080`, and Fiber runs on port `8081`. Ensure that both ports are available.
- Update `views/*.html` to match your template files.

