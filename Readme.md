# Apollo API

## Project Overview

Backend Service

## Key Dependencies

This project relies on several key open-source libraries:

*   **[Echo](https://echo.labstack.com/):** A high-performance, extensible, and minimalist Go web framework.
*   **[Google Wire](https://github.com/google/wire):** A code generation tool for dependency injection.
*   **[golang-migrate](https://github.com/golang-migrate/migrate):** For database migrations.
*   **[gomail](gopkg.in/gomail.v2):** A simple and efficient package to send emails.
*   **[go-redis](https://github.com/redis/go-redis):** Redis client for Go.
*   **[jwt-go](https://github.com/dgrijalva/jwt-go):** A library for working with JSON Web Tokens.

## Prerequisites

Before you begin, ensure you have the following installed:

*   **Go (version 1.23.0 or higher):** The core programming language for this project.
    *   **Installation:** Follow the official Go installation guide: [https://golang.org/doc/install](https://golang.org/doc/install)
*   **Make:** A build automation tool used to simplify common development tasks.
    *   **Installation on macOS/Linux:** Make is typically pre-installed. If not, use your system's package manager (e.g., `sudo apt-get install build-essential` on Debian/Ubuntu, `xcode-select --install` on macOS).
    *   **Installation on Windows:** You can use Chocolatey (`choco install make`) or install it as part of a development environment like MinGW or Cygwin.
*   **PostgreSQL:** The primary database for the application.
    *   **Installation:** Follow the official PostgreSQL installation guide: [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
*   **Redis:** An in-memory data structure store, used as a database, cache, and message broker.
    *   **Installation:** Follow the official Redis installation guide: [https://redis.io/docs/getting-started/](https://redis.io/docs/getting-started/)
*   **SMTP Server:** A Simple Mail Transfer Protocol (SMTP) server is required for sending emails. The project uses `gopkg.in/gomail.v2` to handle email sending.

## Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/winartodev/apollo-be.git
    cd apollo-be
    ```

2.  **Install dependencies:**

    Run the following command to download and install the required Go modules:

    ```bash
    go mod tidy
    ```

    Alternatively, you can use the Makefile:

    ```bash
    make tidy
    ```

## Configuration

1.  **Set up Configuration File:**

    The application is configured using a YAML file. To set up your local configuration, copy the template file:

    ```bash
    cp files/apollo.development.yaml.template files/apollo.development.yaml
    ```

    Then, edit `files/apollo.development.yaml` to set the required values for your database, Redis, SMTP server, JWT secrets, and other configuration options.

2.  **Database Migrations:**

    Before running the application for the first time, you need to apply the database migrations. You can do this by running:

    ```bash
    make auto-migrate
    ```

    This will create the necessary tables in your PostgreSQL database.

## Usage

To run the application, you can use one of the following methods:

### Standard Go Run

```bash
go run cmd/http/http.go
```

This will start the server.

### Hot Reloading with `make`

For development, it's convenient to use the `make run` command, which utilizes [Air](https://github.com/cosmtrek/air) for live reloading. The server will automatically restart when you make changes to the code.

```bash
make run
```

The API will be available at `http://localhost:8080`.

#### Configuring Air (`.air.toml`)

Before using `make run`, you'll need to configure the live-reloading behavior by creating an `.air.toml` file (you can copy `.air.toml.template` to get started).

The build commands in `.air.toml` need to be adjusted based on your operating system. Open your `.air.toml` and find the `[build]` section.

**For macOS and Linux:**

```toml
[build]
# Just plain old go build and then run the binary
bin = "./tmp/main"
cmd = "go build -o ./tmp/main ./cmd/http/http.go"
```

**For Windows:**

```toml
[build]
# Just plain old go build and then run the binary
bin = "./tmp/main.exe"
cmd = "go build -o ./tmp/main.exe ./cmd/http/http.go"
```

## Working with Modules

The application is divided into modules, each representing a specific domain of the application (e.g., `user`, `auth`, `country`).

### Module Structure

Each module follows a clean architecture pattern and typically contains the following directories:

*   `delivery`: Handles the presentation layer (e.g., HTTP handlers).
*   `domain`: Contains the core business logic and entities.
*   `repository`: Implements the data access layer.
*   `usecase`: Orchestrates the flow of data between the delivery and repository layers.

Each module also contains the following files for dependency injection with [Google Wire](https://github.com/google/wire):

*   `provider.go`: Defines the providers for the module's components.
*   `wire.go`: Defines the Wire injector for the module.
*   `wire_gen.go`: The generated Wire code (do not edit this file directly).

### Creating a New Module

1.  **Create the module directory:**

    ```bash
    mkdir modules/<module_name>
    ```

2.  **Create the subdirectories:**

    ```bash
    mkdir modules/<module_name>/{delivery,domain,repository,usecase}
    ```

3.  **Implement the module's components:**

    Create the necessary `.go` files within each subdirectory to implement the module's functionality.

4.  **Set up dependency injection:**

    *   Create a `provider.go` file to define the providers for the module's components.
    *   Create a `wire.go` file to define the Wire injector for the module.

5.  **Update the main Wire configuration:**

    If your new module needs to be integrated with other modules, you may need to update the `wire.go` file in the `cmd/http` directory to include the new module's providers.

6.  **Generate the Wire code:**

    Run the following command to generate the `wire_gen.go` files:

    ```bash
    make wire
    ```

### Updating an Existing Module

1.  **Modify the module's code:**

    Make the necessary changes to the files within the module's directory.

2.  **Update the dependency injection:**

    If you've added or removed components, update the `provider.go` and `wire.go` files accordingly.

3.  **Regenerate the Wire code:**

    Run `make wire` to update the `wire_gen.go` files.

## Makefile Commands

The `Makefile` provides several useful commands for development:

*   `run`: Starts the HTTP server with live reloading.
*   `test`: Runs the test suite.
*   `test-coverage`: Runs the tests and generates an HTML coverage report.
*   `migrate-create NAME=<migration_name>`: Creates a new database migration file.
*   `auto-migrate`: Applies all pending database migrations.
*   `migrate-up`: Applies all up migrations.
*   `migrate-down`: Rolls back the last migration.
*   `migrate-status`: Shows the current migration status.
*   `generate-api-doc`: Generates the Swagger API documentation.
*   `build`: Creates a production build of the application.
*   `clean`: Removes build artifacts.
*   `wire`: Generates the dependency injection code.
*   `help`: Displays a list of all available Makefile commands.

## API Documentation

The API documentation is generated using Swagger. Once the application is running, you can access the documentation at:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

To regenerate the documentation, run:

```bash
make generate-api-doc
```

## Directory Structure

```
.
├── cmd                                      # Main application entry point
├── config                                   # Configuration loading logic
├── docs                                     # Swagger documentation
├── files                                    # Contains configuration templates and other files
├── helper                                   # Helper functions
├── infrastructure                           # Infrastructure layer (database, etc.)
├── internal                                 # Internal application logic
├── migrations                               # Database migrations
├── modules                                  # Application modules
├── .air.toml.template                       # Template for Air (live-reloading tool)
├── files/apollo.development.yaml.template   # Template for application configuration
├── go.mod                                   # Go module file
├── go.sum                                   # Go module checksums
├── Makefile                                 # Makefile
└── Readme.md                                # This file
```
