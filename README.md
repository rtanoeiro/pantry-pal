# PantryPal - Self-Hosted Virtual Pantry & Grocery Planner

PantryPal is a comprehensive, user-friendly application designed to help you manage your pantry efficiently and reduce food waste. It offers a wide range of functionalities, including tracking items, managing expiry dates, integrating with shopping carts, generating recipes, and providing AI-powered meal suggestions.

## Features Overview

### Core Features
- **Pantry Management**: Add, remove, or update items in your virtual pantry.
- **Expiry Tracking**: Monitor items as they approach their expiration date with configurable alerts.
- **Shopping Cart Integration**: Sync your cart with the pantry to track what's on hand and plan your shopping list.
- **Recipe Suggestions**: Search and filter recipes based on available ingredients. Generate meal plans using your pantry items.
- **AI-Powered Meal Planning**: Use your local AI model (LLM) or access an external AI provider to get personalized recipe recommendations.

### Advanced Features
- **Multi-User Support**: Collaborate with friends and family by sharing pantries and managing recipes collectively.
- **Dark Mode**: Choose between light or dark user interfaces for a personalized experience.
- **Admin Panel**: Manage users, permissions, and system settings from one central location.

## Getting Started

### Prerequisites
2. Have Docker installed if you're setting up the application on a server.

### Installation

#### Option 1: Docker (Recommended)
Run these commands to set up PantryPal:

```bash

docker pull mrramonster/pantrypal:latest
docker run -d -p PORT-CHOICE:PORT-CHOICE --name pantry-pal \
          --env-file .env \
          mrramonster/pantrypal:latest
```

### Option 2: Binary

Pull the repository and run the script at scripts/start_server.sh. The server should start at the designed port on your env file.


### Contributing

Necessary Tools for Dev:

- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) - Used to generate safe sql
- [air](https://github.com/air-verse/air) - Used to automatically build/start your webserver uppon changes on files

Whenever you're contributing to new features, make sure all tests are passing, then open a PR with a clear description of your changes.

For issues, please create one detailing the error with screenshots, if possible.