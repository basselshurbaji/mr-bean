# Mr. Bean — Mobile

React Native app (iOS, Android). Talks to the Go backend at `../backend` over HTTP.

---

## Stack

| Layer            | Library                                                              |
| ---------------- | -------------------------------------------------------------------- |
| Framework        | Expo SDK 54 + React Native 0.81.5 + React 19                        |
| Language         | TypeScript 5.9, strict mode, `@/*` path alias                       |
| Navigation       | React Navigation 7 — NativeStack + BottomTabs (manual, not Expo Router) |
| Animation        | react-native-reanimated 4 (new arch, worklets)                      |
| Keyboard         | react-native-keyboard-controller                                     |
| Auth tokens      | expo-secure-store                                                    |
| HTTP             | `src/config/api.ts` — thin `fetch` wrapper                          |
| Icons            | react-native-svg (domain icons) + @expo/vector-icons/Feather (UI)   |
| Fonts            | Playfair Display, DM Sans, JetBrains Mono (Google Fonts)            |
| State            | React Context + Hooks                                                |
| New Architecture | Enabled (`newArchEnabled: true`)                                     |

---

## Project Structure

```
app/                            Screen components (grouped by route)
  (auth)/login.tsx              Login / Register
  (tabs)/
    home/
      index.tsx                 Home dashboard + extraction history
      ExtractionModal.tsx       Bottom sheet — log a new extraction
    beans/
      index.tsx                 Bean library list
      [id].tsx                  Bean detail / edit
      BeanSheet.tsx             Bottom sheet — add / edit a bean
    gear/
      index.tsx                 My Gear + Stations
      [id].tsx                  Gear item detail / edit
      GearSheet.tsx             Bottom sheet — add / edit a gear item
      StationSheet.tsx          Bottom sheet — add / edit a station
    profile.tsx                 Profile + logout
assets/                         Images, splash, icon
src/
  api/
    beans.ts                    Bean types + beansApi (CRUD)
    gear.ts                     GearItem, Station types + gearApi (CRUD)
    extractions.ts              Extraction types + extractionApi, computeZone helper
  components/
    GearIcon.tsx                SVG icon resolver by gear type_id
    RoastBubble.tsx             Roast level badge
  config/
    api.ts                      API_URL constant + apiFetch helper
  context/
    AuthContext.tsx             { isAuthenticated, ready, setIsAuthenticated }
    UserContext.tsx             { user, loading, setUser, logout, refreshUser }
    BeansContext.tsx            { beans, loading, refresh, addBean, updateBean, removeBean }
    GearContext.tsx             { gear, stations, loading, refresh, …CRUD helpers }
    ExtractionsContext.tsx      { extractions, loading, refresh, addExtraction, removeExtraction }
  lib/
    apiClient.ts                authorizedFetch — injects Bearer token, auto-refreshes on 401
    auth.ts                     SecureStore token helpers (save / get / clear)
  navigation/
    index.tsx                   Root App component — all navigators and providers defined here
  theme/
    colors.ts                   Semantic color tokens → RN values
    typography.ts               Font families, sizes, line heights, text styles
    spacing.ts                  Spacing scale, border radii, shadow presets
    index.ts                    Re-exports everything from theme/
```

---

## Navigation Architecture

Navigation is implemented manually in `src/navigation/index.tsx` using React Navigation — **not** Expo Router (the `expo-router` plugin in `app.config.ts` is present but routing is handled programmatically).

```
NavigationContainer
└── Root  (conditional on isAuthenticated)
    ├── AuthNavigator    NativeStack → LoginScreen
    └── MainTabs         BottomTabs
        ├── HomeTab      HomeScreen (+ Beans/Gear/Extractions providers)
        ├── BeansNavigator  NativeStack → BeanList → BeanDetail
        ├── GearNavigator   NativeStack → GearList → GearDetail
        └── ProfileScreen
```

Provider order (outermost → innermost):
```
SafeAreaProvider → KeyboardProvider → AuthProvider → UserProvider
  → NavigationContainer → per-tab providers (Beans / Gear / Extractions)
```

---

## Design System

Design tokens live in `../design/design-system/` — use `src/theme/` which is the RN translation.

- **Colors**: `src/theme/colors.ts`
- **Typography**: `src/theme/typography.ts` — Playfair Display (display), DM Sans (body), JetBrains Mono (data)
- **Spacing / Radii / Shadows**: `src/theme/spacing.ts`
- **Brand & voice**: `../design/design-system/README.md`
- **Component specs**: `../design/design_handoff_*/README.md` — read before implementing a screen

Always `import { colors, textStyles, spacing, radii, shadows } from '@/src/theme'` — never hardcode hex values.

---

## HTTP & Auth

**Request chain:**
1. `apiFetch` (`src/config/api.ts`) — base wrapper, sets `Content-Type`, reads `API_URL`
2. `authorizedFetch` (`src/lib/apiClient.ts`) — adds `Authorization: Bearer`, retries once after token refresh on 401
3. API modules (`src/api/`) call `authorizedFetch` — never call `fetch` directly

**Tokens** are stored in `expo-secure-store` (keys: `mr_bean_access`, `mr_bean_refresh`). Never use AsyncStorage for tokens.

---

## Backend API

Base URL from `appConfig.json` → `server_url` (default `http://localhost:8080`).

| Method | Path                  | Auth   | Notes                                  |
| ------ | --------------------- | ------ | -------------------------------------- |
| POST   | /auth/register        | public | First name, last name, email, password |
| POST   | /auth/login           | public | Returns access + refresh tokens        |
| POST   | /auth/refresh         | public | Rotate refresh token                   |
| GET    | /user/me              | bearer | Current user profile                   |
| PATCH  | /user/me              | bearer | Update profile                         |
| POST   | /user/change-password | bearer | Change password                        |
| GET    | /health               | public | Liveness check                         |
| GET    | /gear                 | bearer | List gear items                        |
| POST   | /gear                 | bearer | Create a gear item                     |
| PUT    | /gear/:id             | bearer | Update a gear item                     |
| DELETE | /gear/:id             | bearer | Delete a gear item (204)               |
| GET    | /stations             | bearer | List stations                          |
| POST   | /stations             | bearer | Create a station                       |
| PUT    | /stations/:id         | bearer | Update a station                       |
| DELETE | /stations/:id         | bearer | Delete a station (204)                 |
| GET    | /beans                | bearer | List beans                             |
| POST   | /beans                | bearer | Create a bean                          |
| PUT    | /beans/:id            | bearer | Update a bean                          |
| DELETE | /beans/:id            | bearer | Delete a bean (204)                    |
| GET    | /extractions          | bearer | List extractions (limit: 20)           |
| POST   | /extractions          | bearer | Log an extraction                      |
| DELETE | /extractions/:id      | bearer | Delete an extraction (204)             |

---

## Backend Host Configuration

Edit `appConfig.json` at the project root. The file is committed with `localhost` defaults:

```json
{
  "server_url": "http://192.168.1.42:8080"
}
```

On simulators `localhost` works. On a physical device use your machine's LAN IP (`make local-ip` writes it automatically).

---

## Getting Started

```bash
npm install

npm run ios                   # iOS Simulator
npm run android               # Android emulator

make typecheck                # tsc --noEmit
make local-ip                 # write machine LAN IP into appConfig.json
```

---

## Conventions

- Sentence case for all UI labels, buttons, and headings (Title Case only for "Mr. Bean")
- Numbers with units: `18.5 g`, `28 s`, `1:2.1` — always use JetBrains Mono (`textStyles.mono`)
- Icons: `react-native-svg` for domain/gear icons, `@expo/vector-icons/Feather` for UI icons
- No hardcoded colors or font names — always go through `src/theme`
- `StyleSheet.create()` for all styles; no inline style objects in JSX

### Keyboard Avoidance

`react-native-keyboard-controller` provides frame-synced keyboard avoidance via Reanimated worklets. `KeyboardProvider` wraps the root in `src/navigation/index.tsx`. No `Platform.OS` checks needed.

**Full-screen screens** — replace `ScrollView` with `KeyboardAwareScrollView`:

```tsx
import { KeyboardAwareScrollView } from 'react-native-keyboard-controller';

<KeyboardAwareScrollView keyboardShouldPersistTaps="handled">
  {/* inputs */}
</KeyboardAwareScrollView>
```

**Bottom sheet modals** — lift the sheet at the `Modal` level, scroll inside the sheet:

```tsx
import { KeyboardAvoidingView, KeyboardAwareScrollView } from 'react-native-keyboard-controller';

<Modal transparent>
  <KeyboardAvoidingView style={{ flex: 1 }} behavior="padding">
    <View style={styles.overlay}>
      <Animated.View style={styles.sheet}>
        <KeyboardAwareScrollView keyboardShouldPersistTaps="handled">
          {/* inputs */}
        </KeyboardAwareScrollView>
      </Animated.View>
    </View>
  </KeyboardAvoidingView>
</Modal>
```

### Safe Area

`SafeAreaProvider` is in `src/navigation/index.tsx`. Navigators with `headerShown: false` do not handle insets automatically — each screen must handle its own.

| Pattern                                       | When to use                                             |
| --------------------------------------------- | ------------------------------------------------------- |
| `<SafeAreaView edges={['top']}>`              | Most screens — simple top inset                        |
| `<SafeAreaView edges={['top', 'bottom']}>`    | Auth screens — full wrap including keyboard scroll area|
| `useSafeAreaInsets()` → apply `insets.top`    | Screens with absolute-positioned overlays              |

Never hardcode status-bar offsets. Import from `react-native-safe-area-context`.