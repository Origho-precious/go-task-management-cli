### Project Overview:

The **Go Task Management CLI** is a small personal project built with Go (Golang) I built to learn Go (Golang). It utilizes the `huh` package for handling terminal forms and offers basic CRUD operations for managing tasks. The CLI allows users to add tasks with descriptions and due dates, mark tasks as completed, delete tasks, and exit the program seamlessly.

### Features:

1. **Add Task**: Users can add a task by providing a description and due date.
2. **Mark as Completed**: Tasks can be marked as completed.
3. **Delete Task**: Users can delete a task.
4. **Exit CLI**: Provides a clean exit option for the CLI without needing to use `ctrl + C`.

### Getting Started and Installation:

To get started with the project, clone the repository and run:

```bash
go mod download
```

or

```bash
go mod tidy
```

 You need to create a file in your root directory called `.env` and add the following environment variables:

```bash
DB_HOST=MY_HOST
DB_PORT=PORT (default is 5432)
DB_USER=DB_ADMIN
DB_PASSWORD=MYSUPERPASSWORD
DB_NAME=MY_DATABASE_NAME
SSL_MODE=disable # leave as is for development
ADMINER_PORT=ADMINER_PORT
```

Make sure to replace the placeholders with your actual database credentials.

The project utilizes Docker for simplicity. By running `docker compose up`, both PostgreSQL and Adminer (for database visualization) will be set up automatically.

### Architecture:

The project comprises two main packages: `main` and `configs`. The `main` package contains the core functionality of the CLI, while the `configs` package handles database connection and environment setup.

### Technologies Used:

- Go (Golang)
- `huh` package for terminal forms handling
- `database/sql` package for database operations
- `pgx` package for connecting to PostgreSQL database
- `docker` for setting up PostgreSQL and Adminer containers

### Usage:

- Customize port settings and database credentials in the `docker-compose.yml` file.
- Run the CLI using `go run main.go`.

### Known Issues/Limitations:

- Lack of authentication functionality.
- No unit tests currently implemented.

### Future Enhancements:

- Implement unit tests for better code reliability.
- Consider adding authentication for broader usage scenarios.

### Credits:

- The project utilizes the `huh` package for terminal forms handling.
- PostgreSQL database connection is simplified with the `pgx` package.
- Docker containers for PostgreSQL and Adminer facilitated easy database setup and management.

### References:

- [GitHub Repository](https://github.com/Origho-precious/go-task-management-cli)
- [huh Package](https://github.com/charmbracelet/huh)
- [pgx Package](https://github.com/jackc/pgx)
- [Adminer](https://github.com/dockage/adminer)
