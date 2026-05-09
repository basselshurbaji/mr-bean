.PHONY: typecheck build clean local-ip

typecheck:
	@npx tsc --noEmit

build:
	@npx react-native bundle --platform android --dev false --entry-file index.js --bundle-output /tmp/mrbean-bundle.js --reset-cache

local-ip:
	@IP=$$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null) && \
	[ -n "$$IP" ] || (echo "error: could not detect local IP" && exit 1) && \
	sed -i '' "s|\"server_url\": \"http://[^\"]*\"|\"server_url\": \"http://$$IP:8080\"|" appConfig.json && \
	echo "server_url → http://$$IP:8080"