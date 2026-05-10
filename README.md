# Mr. Bean

**Mr. Bean** is a coffee tracking system with an MCP server at its core — giving Claude (and other LLMs) direct access to your espresso data. Log shots, manage beans, and get dial-in advice just by talking.

No UI required. Just Claude and your data.

---

## What you can do

The MCP server exposes your full brew journal to Claude. Once connected, you talk — Claude acts.

### Log shots mid-session

> *"Just pulled 18g in, 38g out in 26 seconds on the Hoffman blend. A little fast — grind 11."*

Claude logs the extraction, notes the grind size, and has the context ready for follow-up.

### Dial in a new bean

> *"I'm starting on a washed Ethiopian from Onyx. Light roast. What should I try first?"*

Claude looks at your extraction history, sees how you've done with similar profiles, and gives you a starting point — not a generic recipe from the internet.

### Audit your history

> *"How have my last 10 shots on the Linea Micra compared to the Decent?"*

Claude pulls your extractions, filters by gear, and breaks down the patterns: ratio trends, time variance, where you keep drifting.

### Manage your setup

> *"Add a Niche Zero to my gear list and put it in the home station."*

Done. No forms, no tapping through menus.

### Keep your bean library current

> *"I just finished the Onyx bag. Archive it and add a new one — natural Yemeni from Yemen Mocha, medium roast."*

Claude updates both in one turn.

---

## Tools

21 tools across five domains:

| Domain      | Tools                                                                                                |
| ----------- | ---------------------------------------------------------------------------------------------------- |
| Beans       | `list_beans`, `create_bean`, `update_bean`, `delete_bean`                                            |
| Gear        | `list_gear`, `create_gear`, `get_gear`, `update_gear`, `delete_gear`                                 |
| Stations    | `list_stations`, `create_station`, `update_station`, `delete_station`                                |
| Extractions | `list_extractions`, `create_extraction`, `get_extraction`, `update_extraction`, `delete_extraction`  |
| User        | `get_me`, `update_me`                                                                                |

All tools operate on behalf of the authenticated user identified by the app token.

---

## Stack

| Layer   | Technology                                                        |
| ------- | ----------------------------------------------------------------- |
| Backend | Go (chi router, JWT auth, PostgreSQL via sqlc)                    |
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

`make up` builds all images, runs database migrations, and starts PostgreSQL and the backend. The MCP server will show as **Exited** until you complete step 3 — that's expected.

### 2. Create an account

```bash
curl -s -X POST http://localhost:7489/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"secret"}'
```

Save the `access_token` from the response.

### 3. Create an app token

```bash
curl -s -X POST http://localhost:7489/app-tokens \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"claude"}'
```

Copy the returned token value, add it to `.env`:

```
TOKEN=<your_app_token>
```

Run `make up` again. The MCP server will start and stay running.

### 4. Connect Claude Desktop

Open `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) and add:

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

Restart Claude Desktop. The mr-bean tools appear in the tool picker.

> Claude Desktop spawns the container on-demand via stdio. The long-running MCP container (`make up`) serves HTTP on port 8934 and is useful for testing with other clients.

### Makefile reference

| Command      | Effect                                     |
| ------------ | ------------------------------------------ |
| `make up`    | Build images, start services, show status  |
| `make down`  | Stop and remove containers                 |
| `make logs`  | Stream logs from all services              |
| `make ps`    | Show current service status                |
| `make build` | Rebuild images without starting            |

---

## Repo layout

```text
mr_bean/
├── backend/    # Go API server
├── mcp/        # MCP server — the current focus
├── mobile/     # React Native app (Expo) — on hold, not the primary interface
├── design/     # Design tokens, mockups, and handoff specs
├── product/    # Product notes and planning
└── review/     # Code review artifacts
```
