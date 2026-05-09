.PHONY: typecheck build ios android github-build clean local-ip

typecheck:
	@npx tsc --noEmit

# Generate ios/ and android/ from Expo config — run this after a clean or first checkout,
# then run `make ios` or `make android` once to compile the native binary.
build:
	@npm install
	@npx expo prebuild --clean

# Compile native code and launch in simulator (needed after `make build` or native changes)
ios:
	@npx expo run:ios

android:
	@npx expo run:android

# Verify the JS bundle compiles (used in CI)
github-build:
	@npx react-native bundle --platform android --dev false \
		--entry-file index.js \
		--bundle-output /tmp/mrbean-bundle.js \
		--reset-cache

# Remove all generated artifacts — follow with `make build` to regenerate
clean:
	@rm -rf ios android node_modules .expo dist web-build

local-ip:
	@IP=$$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null) && \
	[ -n "$$IP" ] || (echo "error: could not detect local IP" && exit 1) && \
	sed -i '' "s|\"server_url\": \"http://[^\"]*\"|\"server_url\": \"http://$$IP:8080\"|" appConfig.json && \
	echo "server_url → http://$$IP:8080"
