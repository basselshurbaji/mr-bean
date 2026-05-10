# Mr. Bean

[![build](https://github.com/basselshurbaji/mr-bean/actions/workflows/go.yml/badge.svg)](https://github.com/basselshurbaji/mr-bean/actions/workflows/go.yml)

Personal espresso journal with an MCP server. Claude gets direct read/write access to your shots, beans, and gear — log a pull, check trends, or update your setup without leaving the chat.

Backend is a Go REST API (chi, sqlc, PostgreSQL). The MCP server speaks stdio to Claude Desktop.

---

## What you can do

Connect the MCP server and Claude has your full brew history. Some things worth knowing:

### Log a shot mid-session

> "18g in, 36.5g out, 27 seconds. A bit fast — bump the grind to 12."

Claude logs the extraction and notes the grind adjustment. Each shot stores dose, yield, time, your target time, and grind size — zone classification (under / perfect / over) is derived from how far off your actual time was from target.

### Get a starting point on a new bean

> "Starting on a washed Ethiopian from Onyx. Light roast. Where should I begin?"

Claude checks your history for extractions on similar profiles (same process, roast level) and gives you a starting point based on what's actually worked — not a generic internet recipe.

### Review a session

> "How have my last 10 shots on the Linea Micra compared to the Decent?"

Claude filters by gear, pulls the numbers, and breaks down ratio trends and time variance.

### Update your setup

> "Add a Niche Zero to my gear list. Put it in the home station."

Gear items live independently of stations — stations are just ordered pre-selection groups for logging. Deleting a station doesn't touch your shot history.

### Keep your bean catalogue current

> "Just finished the Onyx bag. Add a new one — natural Yemeni from Qima, medium roast."

---

## Tools

20 tools across five domains:

| Domain      | Tools                                                                                               |
|-------------|-----------------------------------------------------------------------------------------------------|
| Beans       | `list_beans`, `create_bean`, `update_bean`, `delete_bean`                                           |
| Gear        | `list_gear`, `create_gear`, `get_gear`, `update_gear`, `delete_gear`                                |
| Stations    | `list_stations`, `create_station`, `update_station`, `delete_station`                               |
| Extractions | `list_extractions`, `create_extraction`, `get_extraction`, `update_extraction`, `delete_extraction` |
| User        | `get_me`, `update_me`                                                                               |

All tools operate on behalf of the authenticated user identified by the app token.

---

## Stack

| Layer   | Technology                                                       |
|---------|------------------------------------------------------------------|
| Backend | Go (chi router, JWT auth, PostgreSQL via sqlc)                   |
| MCP     | Go — official MCP go-sdk, stdio (Claude Desktop) + HTTP (Docker) |

---

## Setup

Everything runs in Docker. The only prerequisite is Docker with Compose.

### 1. Clone and start

```bash
git clone https://github.com/basselshurbaji/mr_bean
cd mr_bean
cp .env.example .env
make up
```

`make up` builds all images (including `mr-bean-mcp:latest`), runs database migrations, and starts PostgreSQL and the backend.

### 2. Connect Claude

```bash
make setup
```

Asks whether you have an account, walks you through login or registration, creates an app token, and prints the config block to paste into Claude. One command, done.

### 3. Add the config to Claude

Paste the block printed by `make setup` into the appropriate file:

- **Claude Desktop** — `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Claude Code** — `~/.claude.json` (under `mcpServers`)

Restart Claude Desktop. The mr-bean tools appear in the tool picker.

> Claude Desktop spawns the container on-demand via stdio. The `mr-bean-mcp:latest` image is built by `make up` but not run as a service — Claude Desktop manages its lifecycle.

### Step by step

If you prefer to run each step manually:

```bash
make account    # register a new account
make login      # or log in to an existing one
make app-token  # create an app token and print the config
```

### Makefile reference

| Command                                                              | Effect                                              |
|----------------------------------------------------------------------|-----------------------------------------------------|
| `make up`                                                            | Build images, start services, show status           |
| `make down`                                                          | Stop and remove containers                          |
| `make logs`                                                          | Stream logs from all services                       |
| `make ps`                                                            | Show current service status                         |
| `make build`                                                         | Rebuild images without starting                     |
| `make health`                                                        | Verify the backend is up                            |
| `make setup`                                                         | Register or log in, create app token, print config  |
| `make account`                                                       | Register a new account                              |
| `make login`                                                         | Log in to an existing account                       |
| `make app-token`                                                     | Create an app token and print the Claude config     |

---

## Repo layout

```text
mr_bean/
├── backend/    # Go API server
├── mcp/        # MCP server
└── product/    # Product notes and planning
```
