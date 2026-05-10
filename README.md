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

| Layer   | Technology                                     |
| ------- | ---------------------------------------------- |
| Backend | Go (chi router, JWT auth, PostgreSQL via sqlc) |
| MCP     | Go - official MCP go-sdk, stdio transport      |

---

## Setup

### 1. Start the backend

```bash
cd backend
docker compose up
```

See `backend/README.md` for full instructions.

### 2. Get an app token

With the backend running, create a long-lived app token:

```bash
curl -X POST http://localhost:8080/app-tokens \
  -H "Authorization: Bearer <your-jwt>" \
  -H "Content-Type: application/json" \
  -d '{"name": "claude"}'
```

Copy the returned token.

### 3. Build the MCP binary

```bash
cd mcp
go build -o mr-bean-mcp .
```

### 4. Add to Claude Desktop

Open `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) and add:

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

Restart Claude Desktop. The mr-bean tools will appear in the tool picker.

See `mcp/README.md` for Docker setup and additional configuration options.

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
