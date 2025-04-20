# Luminex 🚀

## Overview 🔍

Luminex is a Go-based microservice. Currently, it provides analytics and insights for GitHub repositories. It extracts and processes GitHub data to deliver meaningful metrics and statistics for engineering teams and project managers.

## Configuration ⚙️

### Main Configuration 📝

The main configuration is stored in `configs/config.yaml`. This file contains general settings for the application such as server ports, database connection details, and other non-sensitive configuration.

### Secrets Management 🔐

Secrets are managed through JSON files in the `configs/secrets/` directory or environment variables:

- `configs/secrets/github.json` - GitHub credentials
- `configs/secrets/database.json` - Database credentials
- `configs/secrets/api_keys.json` - Various API keys

For production deployments, use environment variables or a secure secrets management service.

## Running the Application 🏃‍♂️

```bash
# With secrets files
go run cmd/luminex-service/main.go

# Using environment variables (fallback)
export GITHUB_TOKEN=your_token_here
go run cmd/luminex-service/main.go
```

## Docker Support 🐳

```bash
# Build the Docker image
docker build -t luminex-service:latest .

# Run the container
docker run -p 8000:8000 -p 9000:9000 -e GITHUB_TOKEN=your_token_here luminex-service:latest
```

## API Endpoints 🌐

The application exposes both HTTP and gRPC endpoints:

- HTTP: `http://localhost:8000/v1/`
- gRPC: `localhost:9000`

### Key Endpoints 🔑

- `/v1/health` - Health check
- `/v1/metrics` - PR metrics for a repository
- `/v1/monthly-stats` - Monthly statistics
- `/v1/repo-stats` - Repository statistics
- `/v1/contributor-stats` - Contributor statistics
- `/v1/issue-stats` - Issue statistics
- `/v1/detailed-pr-stats` - Detailed PR statistics

## Project Structure 📂

```
├── cmd/                    # Application entry points
│   └── luminex-service/    # Main server command
├── configs/                # Configuration files
├── constants/              # Application constants
├── internal/               # Private application code
│   ├── biz/                # Business logic
│   ├── conf/               # Configuration processing
│   ├── data/               # Data processing
│   ├── helpers/            # Helper functions
│   ├── server/             # Server implementation
│   ├── service/            # Service implementation
│   └── interfaces/         # Interface definitions
├── models/                 # Data models
└── utils/                  # Utility functions
```

## Features ✨

- Pull Request analytics (merge time, open PRs, etc.)
- Monthly repository activity metrics
- Contributor statistics and leaderboards
- Repository overview metrics
- Issue tracking and analysis
- Detailed pull request analysis

## Future Roadmap 🗺️

This service is under active development. Will be adding some cool features. 😉

## Development 👨‍💻

For local development:

1. Clone the repository
2. Set up your GitHub token in configs/secrets or as an environment variable
3. Run `go mod download` to install dependencies
4. Use `go run cmd/luminex-service/main.go` to start the service

