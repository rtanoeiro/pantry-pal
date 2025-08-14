# PantryPal - Self-Hosted Virtual Pantry & Grocery Planner

PantryPal is a comprehensive, user-friendly application designed to help you manage your pantry efficiently and reduce food waste. It offers a wide range of functionalities, including tracking items, managing expiry dates, integrating with shopping carts, generating recipes, and providing AI-powered meal suggestions.

## Features Overview

### Core Features
- **Pantry Management**: Add, remove, or update items in your virtual pantry.
- **Expiry Tracking**: Monitor items as they approach their expiration date with configurable alerts.
- **Shopping Cart Integration**: Sync your cart with the pantry to track what's on hand and plan your shopping list.

### Next Steps for Development
- **Recipe Suggestions**: Search and filter recipes based on available ingredients. Generate meal plans using your pantry items.
- **AI-Powered Meal Planning**: Use your local AI model (LLM) or access an external AI provider to get personalized recipe recommendations. This can be done once the Go SDK MCP package is complete: https://github.com/modelcontextprotocol/go-sdk

## Future Work
- **Dark Mode**: Choose between light or dark user interfaces for a personalized experience.

## Getting Started
### Installation

#### Option 1: Docker (Recommended)
You can run PantrPal using Docker, use the docker-compose-template.yml file to set up your container. Make sure to define a port, and to create a JWT Secret before running the container. You can rename your file to pantry-pal.yaml and run the following command:

```bash
docker-compose -d up
```

Head over to localhost:8080/login and you should have access to the application. The default user/password is Admin/admin, new users by default are not assigned as admin, but if you log into the Admin account, you can give admin privileges to any account.

### Option 2: Binary
Use this option if you want to run this application without Docker, useful for development purposes, or if you are not familiar with Docker.

1. Pull the repository with:

```bash
git clone https://github.com/rtanoeiro/pantry-pal.git
```

2. You'll need some environment variables. Enter the folder of the repository you just downloaded and run:
``` bash
export PORT=8080
export DATABASE_URL=data/pantry_pal.db
# You can create a Secret with the following command ``echo $(openssl rand -base64 32)```
# Copy the output from the commmand above and paste it in the next line
export JWT_SECRET=THING_YOU_JUST_COPIED
```

3. After that, you can download all dependencies with:

```bash
CGO_ENABLED=1 go build -o pantry_pal
```

4. Now, you should be ready to Go! You can run the application with:

```bash
./pantry_pal
```

### Contributing

### Necessary Tools

- go 1.24.2+
- sqlc [https://sqlc.dev/](https://sqlc.dev/) - Install with go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
- goose [https://github.com/pressly/goose](https://github.com/pressly/goose) - Package for handling migrations in the database. Install with go install github.com/pressly/goose/v3/cmd/goose@latest
- air [https://github.com/air-verse/air](https://github.com/air-verse/air) - Install with go install github.com/air-verse/air@latest
    - Not really necessary, but it helps as it reloads the server on file changes.
- Docker - For running the application in a containerized environment.

Whenever you're contributing to new features, make sure all tests are passing, then open a PR with a clear description of your changes.

For issues, please create one detailing the error with screenshots, if possible.
