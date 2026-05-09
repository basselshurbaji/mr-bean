const { getDefaultConfig } = require('expo/metro-config');
const path = require('path');

const projectRoot = path.resolve(__dirname);
const config = getDefaultConfig(__dirname);

// Expo's config may set its own resolveRequest (handles .expo/.virtual-metro-entry etc.)
// Chain our @/ alias on top so both work.
const expoResolveRequest = config.resolver.resolveRequest;

config.resolver.resolveRequest = (context, moduleName, platform) => {
  const resolve = expoResolveRequest ?? context.resolveRequest;
  if (moduleName.startsWith('@/')) {
    return resolve(context, path.resolve(projectRoot, moduleName.slice(2)), platform);
  }
  return resolve(context, moduleName, platform);
};

module.exports = config;
