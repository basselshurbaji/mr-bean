.PHONY: up down build logs ps

up:
	@docker compose up --build -d
	@echo ""
	@docker compose ps
	@echo ""
	@echo "================================================================"
	@echo " mr-bean"
	@echo "================================================================"
	@echo ""
	@echo " Backend API:  http://localhost:7489"
	@echo " MCP server:   http://localhost:8934  (starts once TOKEN is set)"
	@echo ""
	@echo " ---- First time? Set up your account: -------------------------"
	@echo ""
	@echo " 1. Register:"
	@echo "    curl -s -X POST http://localhost:7489/auth/register \\"
	@echo "      -H 'Content-Type: application/json' \\"
	@echo "      -d '{\"email\":\"you@example.com\",\"password\":\"secret\"}'"
	@echo ""
	@echo " 2. Create an app token (replace <jwt> with the access_token):"
	@echo "    curl -s -X POST http://localhost:7489/app-tokens \\"
	@echo "      -H 'Authorization: Bearer <jwt>' \\"
	@echo "      -H 'Content-Type: application/json' \\"
	@echo "      -d '{\"name\":\"claude\"}'"
	@echo ""
	@echo " 3. Copy the returned token, add it to .env:"
	@echo "    TOKEN=<your_app_token>"
	@echo "    Then run 'make up' again to start the MCP server."
	@echo ""
	@echo " ---- Connect Claude Desktop: -----------------------------------"
	@echo ""
	@echo " Add to ~/Library/Application Support/Claude/claude_desktop_config.json:"
	@echo ""
	@echo "   {"
	@echo "     \"mcpServers\": {"
	@echo "       \"mr-bean\": {"
	@echo "         \"command\": \"docker\","
	@echo "         \"args\": ["
	@echo "           \"run\", \"--rm\", \"-i\","
	@echo "           \"-e\", \"TOKEN=<your_app_token>\","
	@echo "           \"-e\", \"MR_BEAN_SERVER_URL=http://host.docker.internal:8080\","
	@echo "           \"mr-bean-mcp:latest\""
	@echo "         ]"
	@echo "       }"
	@echo "     }"
	@echo "   }"
	@echo ""
	@echo " Restart Claude Desktop. The mr-bean tools will appear in the tool picker."
	@echo "================================================================"

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

ps:
	docker compose ps