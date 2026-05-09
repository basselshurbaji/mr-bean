# mr-bean MCP

An [MCP](https://modelcontextprotocol.io) server that gives Claude (and other LLMs) access to your mr-bean espresso tracking data. Built with the [official Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk).

## What it does

Exposes 20 tools across four domains:

| Domain      | Tools                                                                                               |
|-------------|-----------------------------------------------------------------------------------------------------|
| Beans       | `list_beans`, `create_bean`, `update_bean`, `delete_bean`                                           |
| Gear        | `list_gear`, `create_gear`, `get_gear`, `update_gear`, `delete_gear`                                |
| Stations    | `list_stations`, `create_station`, `update_station`, `delete_station`                               |
| Extractions | `list_extractions`, `create_extraction`, `get_extraction`, `update_extraction`, `delete_extraction` |
| User        | `get_me`, `update_me`                                                                               |

All tools operate on behalf of the authenticated user identified by the app token.

---

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (recommended) **or** Go 1.26+
- A running mr-bean backend (see `../backend`)
- An app token — generate one via the backend's `POST /app-tokens` endpoint

---

## Setup

### 1
. Get an app token

With the backend running, create an app token. You can use the mr-bean mobile app or curl:

```sh
curl -X POST http://localhost:8080/app-tokens \
  -H "Authorization: Bearer <your-jwt-session-token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "claude-mcp"}'
```

Copy the returned token — you'll need it below.

### 2. Run with Docker (recommended)

```sh
# from the mcp/ directory
TOKEN=your_app_token docker compose up --build
```

The MCP server will be available at `http://localhost:8081`.

**If your backend is also running locally** (not in Docker), `host.docker.internal` resolves correctly on macOS and Windows. On Linux, add `--add-host=host.docker.internal:host-gateway` to your Docker run command or set `MR_BEAN_SERVER_URL=http://172.17.0.1:8080`.

To point at a remote backend:

```sh
TOKEN=your_app_token MR_BEAN_SERVER_URL=https://api.yourserver.com docker compose up --build
```

### 3. Build the binary and configure Claude Desktop

```sh
# from the mcp/ directory
go build -o mr-bean-mcp .
```

Open (or create) your Claude Desktop config file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Add the `mr-bean` entry under `mcpServers`:

```json
{
  "mcpServers": {
    "mr-bean": {
      "command": "/absolute/path/to/mcp/mr-bean-mcp",
      "env": {
        "TOKEN": "your_app_token",
        "MR_BEAN_SERVER_URL": "http://localhost:8080"
      }
    }
  }
}
```

Restart Claude Desktop. You should see the mr-bean tools available in the tool picker.

> **Note:** The binary auto-detects its transport. When Claude Desktop spawns it as a subprocess it uses stdio; when run directly in a terminal it falls back to HTTP on port 8081.

---

## Running locally (without Docker)

```sh
# from the mcp/ directory
TOKEN=your_app_token go run .

# with a non-default backend URL or port
MR_BEAN_SERVER_URL=http://localhost:8080 TOKEN=your_app_token PORT=8081 go run .
```

---

## Environment variables

| Variable             | Default                 | Required | Description                         |
| -------------------- | ----------------------- | -------- | ----------------------------------- |
| `TOKEN`              | —                       | Yes      | App token for authenticating to the backend |
| `MR_BEAN_SERVER_URL` | `http://localhost:8080` | No       | Base URL of the mr-bean backend     |
| `PORT`               | `8081`                  | No       | Port the MCP server listens on      |

---

## Example usage in Claude

Once connected, you can ask Claude things like:

- *"List all my coffee beans"*
- *"Log an extraction: 18g in, 36g out, 28 seconds, grind size 12, on the Hoffman blend"*
- *"What was my last extraction?"*
- *"Update my grinder's notes to say it was recalibrated today"*

---

## Building the Docker image manually

```sh
docker build -t mr-bean-mcp .

docker run -p 8081:8081 \
  -e TOKEN=your_app_token \
  -e MR_BEAN_SERVER_URL=http://host.docker.internal:8080 \
  mr-bean-mcp
```
