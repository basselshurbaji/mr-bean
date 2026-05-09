# mr-bean MCP

HTTP MCP server that exposes the mr-bean backend REST API as tools for LLMs.

## Overview

- **SDK**: `github.com/modelcontextprotocol/go-sdk` v1.6.0 (official Go MCP SDK)
- **Transport**: Streamable HTTP (`mcp.NewStreamableHTTPHandler`)
- **Auth**: Bearer token passed via `Authorization` header on every backend request
- **Tools**: 20 tools across 4 domains (beans, gear/stations, extractions, user)

## Running

```sh
TOKEN=your_app_token go run .
# with overrides
MR_BEAN_SERVER_URL=http://localhost:8080 TOKEN=<token> PORT=8081 go run .
```

## Environment Variables

| Variable             | Default                 | Required |
|----------------------|-------------------------|----------|
| `TOKEN`              | —                       | Yes      |
| `MR_BEAN_SERVER_URL` | `http://localhost:8080` | No       |
| `PORT`               | `8081`                  | No       |

`TOKEN` must be a valid app token created via the backend's `/app-tokens` endpoint.

## File Structure

```
main.go          — server init, env vars, tool registration
client.go        — generic HTTP client (auth header, JSON encode/decode, error mapping)
beans.go         — list_beans, create_bean, update_bean, delete_bean
gear.go          — list_gear, create_gear, get_gear, update_gear, delete_gear
                   list_stations, create_station, update_station, delete_station
extractions.go   — list_extractions, create_extraction, get_extraction,
                   update_extraction, delete_extraction
user.go          — get_me, update_me
```

## Adding a New Tool

1. Define a params struct with `json` and `jsonschema` tags. Pointer fields + `omitempty` = optional in schema; non-pointer fields = required.
2. Call `mcp.AddTool(server, &mcp.Tool{Name: "...", Description: "..."}, handlerFunc)` in the appropriate `register*Tools` function.
3. Handler signature: `func(ctx context.Context, req *mcp.CallToolRequest, params MyParams) (*mcp.CallToolResult, any, error)`
4. Return results as `&mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: jsonString}}}`.
5. Register the function in `main.go` if you've added a new file.

## Adding a New Backend Endpoint

If the backend grows new `TagAppAuthenticated` endpoints, add them to the appropriate domain file. The `Client` in `client.go` has `get`, `post`, `put`, `patch`, and `delete` helpers — use those rather than calling `client.do` directly.

## Key Design Decisions

### PUT semantics for updates
Update handlers (`update_bean`, `update_gear`, `update_station`, `update_extraction`) use PUT, meaning they are full replacements. Required fields in the params struct match what the backend requires. For `update_station`, `gear_ids` is always sent (never omitted) — a nil slice is normalized to `[]string{}`, which clears all gear associations. This mirrors the backend's own nil-to-empty normalization.

### Schema inference
Tool input schemas are inferred automatically from the params struct by `google/jsonschema-go`. Use `jsonschema:"description text"` for field descriptions. Do not add `additionalProperties: false` manually — the SDK handles this.

### Error propagation
Backend HTTP errors (4xx/5xx) are returned as Go errors and surface to the LLM as tool errors. The error includes the HTTP status code and the backend's `{"error": "..."}` message. No retry logic — let the LLM decide.

### `time` field in extractions
`CreateExtractionParams.BrewTime` and `UpdateExtractionParams.BrewTime` use `json:"time"` to match the backend's field name. The body maps in `create_extraction` and `update_extraction` explicitly use the string key `"time"` to be unambiguous.

## Connecting to Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mr-bean": {
      "command": "/path/to/mcp/binary",
      "env": {
        "TOKEN": "your_app_token",
        "MR_BEAN_SERVER_URL": "http://localhost:8080"
      }
    }
  }
}
```

The binary auto-detects its transport: when Claude Desktop spawns it (stdin is a pipe) it uses stdio; when run directly in a terminal it falls back to HTTP on port 8081.

The `url` format (`"url": "http://localhost:8081"`) is not supported by Claude Desktop — always use the `command` form.
