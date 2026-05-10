# mr-bean MCP

An [MCP](https://modelcontextprotocol.io) server that gives Claude (and other LLMs) access to your mr-bean espresso tracking data. Built with the [official Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk).

## What it does

Exposes 21 tools across five domains:

| Domain      | Tools                                                                                               |
| ----------- | --------------------------------------------------------------------------------------------------- |
| Beans       | `list_beans`, `create_bean`, `update_bean`, `delete_bean`                                           |
| Gear        | `list_gear`, `create_gear`, `get_gear`, `update_gear`, `delete_gear`                                |
| Stations    | `list_stations`, `create_station`, `update_station`, `delete_station`                               |
| Extractions | `list_extractions`, `create_extraction`, `get_extraction`, `update_extraction`, `delete_extraction` |
| User        | `get_me`, `update_me`                                                                               |

All tools operate on behalf of the authenticated user identified by the app token.

---

## Setup

See the root `README.md` for the full setup guide. The short version:

```sh
# from the project root
make up
```

This builds the MCP image alongside the backend. Once you have an app token, set `TOKEN=<value>` in `.env` and run `make up` again.

---

## Connect to Claude Desktop

```json
{
  "mcpServers": {
    "mr-bean": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-e", "TOKEN=<your_app_token>",
        "-e", "MR_BEAN_SERVER_URL=http://host.docker.internal:7489",
        "mr-bean-mcp:latest"
      ]
    }
  }
}
```

Claude Desktop spawns the container on-demand via stdio. No separate server process is needed.

---

## Environment variables

| Variable             | Default                 | Required | Description                              |
| -------------------- | ----------------------- | -------- | ---------------------------------------- |
| `TOKEN`              | —                       | Yes      | App token created via `POST /app-tokens` |
| `MR_BEAN_SERVER_URL` | `http://localhost:7489` | No       | Base URL of the mr-bean backend          |
| `PORT`               | `8934`                  | No       | HTTP port when not running in stdio mode |
| `MCP_TRANSPORT`      | auto-detect             | No       | Force `stdio` or `http`                  |

Transport is auto-detected: stdio when stdin is a pipe (Claude Desktop), HTTP otherwise (docker-compose service).

---

## Running locally without Docker

```sh
# from the mcp/ directory
TOKEN=your_app_token go run .
```

---

## Example prompts

- *"List all my coffee beans"*
- *"Log an extraction: 18g in, 36g out, 28 seconds, grind size 12, on the Hoffman blend"*
- *"What was my last extraction?"*
- *"Update my grinder's notes to say it was recalibrated today"*
