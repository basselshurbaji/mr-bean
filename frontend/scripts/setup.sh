#!/usr/bin/env bash
# Mr. Bean — Frontend dev environment setup
# Run once after cloning: bash scripts/setup.sh

set -euo pipefail

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BOLD='\033[1m'; NC='\033[0m'

info()    { echo -e "${BOLD}[setup]${NC} $*"; }
success() { echo -e "${GREEN}[setup]${NC} $*"; }
warn()    { echo -e "${YELLOW}[setup]${NC} $*"; }
die()     { echo -e "${RED}[setup] ERROR:${NC} $*"; exit 1; }

# ── 1. Homebrew ──────────────────────────────────────────────────────────────
if ! command -v brew &>/dev/null; then
  info "Installing Homebrew…"
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  # Add brew to PATH for Apple Silicon
  if [[ -f /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
  fi
else
  success "Homebrew already installed ($(brew --version | head -1))"
fi

# ── 2. Node.js ───────────────────────────────────────────────────────────────
if ! command -v node &>/dev/null; then
  info "Installing Node.js via Homebrew…"
  brew install node
else
  success "Node $(node --version) already installed"
fi

# Ensure npm/npx are reachable
export PATH="/opt/homebrew/bin:$PATH"

# ── 3. Expo CLI ──────────────────────────────────────────────────────────────
if ! command -v expo &>/dev/null && ! npx expo --version &>/dev/null 2>&1; then
  info "Installing Expo CLI globally…"
  npm install -g expo-cli
else
  success "Expo CLI available"
fi

# ── 4. npm dependencies ──────────────────────────────────────────────────────
info "Installing npm dependencies…"
npm install

# ── 5. .env ──────────────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

if [[ ! -f "$PROJECT_ROOT/.env" ]]; then
  cp "$PROJECT_ROOT/.env.example" "$PROJECT_ROOT/.env"
  warn ".env created from .env.example — edit EXPO_PUBLIC_API_URL if needed"
else
  success ".env already exists"
fi

# ── 6. Optional: iOS simulator tooling ───────────────────────────────────────
if [[ "$(uname)" == "Darwin" ]]; then
  if ! command -v xcrun &>/dev/null; then
    warn "Xcode command-line tools not found. Install with: xcode-select --install"
  else
    success "Xcode CLT available"
  fi
fi

# ── Done ─────────────────────────────────────────────────────────────────────
echo ""
echo -e "${BOLD}Setup complete.${NC} Start the dev server with:"
echo ""
echo -e "  ${GREEN}npm run ios${NC}       — iOS Simulator"
echo -e "  ${GREEN}npm run android${NC}   — Android emulator"
echo -e "  ${GREEN}npm run web${NC}       — Browser"
echo -e "  ${GREEN}npm start${NC}         — Expo Go / interactive menu"
echo ""
echo -e "Backend URL: ${YELLOW}$(grep EXPO_PUBLIC_API_URL "$PROJECT_ROOT/.env" | tail -1)${NC}"
echo -e "To change it, edit ${BOLD}.env${NC} and restart the dev server."
