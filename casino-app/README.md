# Casino App

## Pages

- Loading Page
- User Create / Login
- User Management
- Game Area
    - Lobby (list of Games)
    - Game Popups
- Admin
    - User List & Moderation

## Environment Variables

The application uses environment variables to configure the WebSocket connection URL. This allows for seamless switching between development and production environments.

### Available Variables

- `PUBLIC_WS_URL`: The WebSocket server URL
  - Production: `wss://casino-host.stmbl.dev/ws`
  - Development: `ws://localhost:9000/ws`

### Setup

1. Create a `.env` file in the project root (or copy from `.env.example`)
2. Set the appropriate variables:

```
# For local development
PUBLIC_WS_URL=ws://localhost:9000/ws

# For production (default if not specified)
# PUBLIC_WS_URL=wss://casino-host.stmbl.dev/ws
```

The `.env` file is ignored by git to prevent sensitive information from being committed.