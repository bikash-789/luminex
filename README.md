# Luminex

Backend service for the GitHub Insights Dashboard project.

## Configuration

The application can be configured using the following approaches:

### Main Configuration

The main configuration is stored in `configs/config.yaml`. This file contains general settings for the application such as server ports, database connection details, and other non-sensitive configuration.

### Secrets Management

The application uses a modular secrets management system where each type of secret is stored in its own file under the `configs/secrets/` directory. This provides better separation and security for different types of credentials.

Secrets are loaded in the following order:

1. From individual JSON files in the `configs/secrets/` directory (primary)
2. From environment variables (fallback)

#### Secret Files Structure

Each type of secret has its own JSON file:

- `configs/secrets/github.json` - GitHub credentials
- `configs/secrets/database.json` - Database credentials
- `configs/secrets/api_keys.json` - Various API keys

Example GitHub secrets file (`github.json`):

```json
{
  "token": "your_github_token_here"
}
```

Example Database secrets file (`database.json`):

```json
{
  "username": "database_user",
  "password": "database_password",
  "connection_string": "user:pass@localhost:3306/mydb"
}
```

Example API keys file (`api_keys.json`):

```json
{
  "stripe": "sk_test_sample_key",
  "aws": "AKIA_sample_key",
  "third_party": {
    "mailchimp": "your_mailchimp_key",
    "sendgrid": "your_sendgrid_key"
  }
}
```

#### Accessing Secrets in Code

Secrets can be accessed using namespace (filename without .json) and key:

```go
// Simple key access
githubToken := secretsManager.GetStringWithEnvFallback("entity", "token", "GITHUB_TOKEN")

// Nested key access (for nested JSON objects)
mailchimpKey := secretsManager.GetNestedStringWithEnvFallback("api_keys", "third_party.mailchimp", "MAILCHIMP_API_KEY")
```

#### Security Best Practices

- Add the entire `configs/secrets/` directory to your `.gitignore` to prevent committing secrets to version control
- For production deployments, use environment variables or a secure secrets management service
- Limit file permissions on the secrets files to only the necessary users/processes

## Running the Application

```bash
# With secrets files
go run cmd/luminex-service/main.go

# Using environment variables (fallback)
export GITHUB_TOKEN=your_token_here
go run cmd/luminex-service/main.go
```

## Docker Support

The application can be containerized using Docker:

```bash
# Build the Docker image
docker build -t luminex-service:latest .

# Run the container with environment variables
docker run -p 8000:8000 -p 9000:9000 -e GITHUB_TOKEN=your_token_here luminex-service:latest

# Or use a local config directory
docker run -p 8000:8000 -p 9000:9000 -v $(pwd)/configs:/app/data/conf luminex-service:latest
```

### Docker Compose (Optional)

You can also use Docker Compose for easier deployment:

```yaml
version: '3'
services:
  luminex-service:
    build:
      context: .
    ports:
      - "8000:8000"
      - "9000:9000"
    environment:
      - GITHUB_TOKEN=${GITHUB_TOKEN:-}
    volumes:
      - ./configs:/app/data/conf
      - ./logs:/var/log/luminex
    restart: unless-stopped
```

Save this as `docker-compose.yml` and run with:

```bash
# Set the GitHub token if needed
export GITHUB_TOKEN=your_github_api_token

# Run with Docker Compose
docker-compose up
```

## API Endpoints

The application exposes both HTTP and gRPC endpoints:

- HTTP: `http://localhost:8000/api/`
- gRPC: `localhost:9000`

### Available Endpoints

- `/api/health` - Health check
- `/api/metrics` - PR metrics for a repository
- `/api/monthly-stats` - Monthly statistics for a repository
- `/api/repo-stats` - Repository statistics
- `/api/contributor-stats` - Contributor statistics
- `/api/issue-stats` - Issue statistics
- `/api/detailed-pr-stats` - Detailed PR statistics

## Project Structure

This project follows the [go-kratos](https://go-kratos.dev/) microservice framework structure:

```
├── api/                    # API definitions (Protocol Buffers)
│   └── github/v1/          # GitHub API v1
├── cmd/                    # Application entry points
│   └── server/             # Main server command
├── configs/                # Configuration files
│   ├── config.yaml         # Main configuration
│   └── secrets/            # Secret files (not in version control)
│       ├── github.json     # GitHub credentials
│       ├── database.json   # Database credentials
│       └── api_keys.json   # Various API keys
├── internal/               # Private application code
│   ├── biz/                # Business logic
│   ├── conf/               # Configuration processing
│   ├── data/               # Data processing
│   ├── pkg/                # Internal packages
│   │   └── secrets/        # Secrets manager
│   ├── server/             # Server implementation
│   └── service/            # Service implementation
└── third_party/            # Third-party code and protobuf imports
```

## Features

- Pull Request statistics (merge time, open PRs, etc.)
- Monthly repository activity
- Contributor statistics
- Repository overview metrics
- Issue statistics and analysis
- Detailed pull request analysis

## Development

For local development:

1. Install the required tools:

```bash
make init
```

2. Generate code from protobuf definitions:

```bash
make api
```

3. Generate dependency injection code:

```bash
make wire
```

