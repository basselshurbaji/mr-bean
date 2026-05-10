.PHONY: up down build logs ps health setup login account app-token

up:
	@docker compose up --build -d
	@docker build -q -t mr-bean-mcp:latest ./mcp
	@echo ""
	@docker compose ps
	@echo ""
	@echo "================================================================"
	@echo " mr-bean"
	@echo "================================================================"
	@echo ""
	@echo " Backend API:  http://localhost:7489"
	@echo " MCP image:    mr-bean-mcp:latest  (built, ready for Claude Desktop)"
	@echo ""
	@echo " ---- Connect Claude in one step: ------------------------------"
	@echo ""
	@echo "  make setup"
	@echo ""
	@echo " Or step by step:"
	@echo ""
	@echo "  make account  /  make login"
	@echo "  make app-token"
	@echo ""
	@echo "================================================================"

down:
	docker compose down

build:
	docker compose build
	docker build -q -t mr-bean-mcp:latest ./mcp

logs:
	docker compose logs -f

ps:
	docker compose ps

health:
	@bash scripts/health.sh

setup: health
	@bash scripts/setup.sh

login: health
	@bash scripts/login.sh

account: health
	@bash scripts/account.sh

app-token: health
	@bash scripts/app_token.sh
