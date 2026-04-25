# Mr. Bean — Mobile

React Native / Expo cross-platform app (iOS, Android, Web). Talks to the Go backend at `../backend` over HTTP.

---

## Stack

| Layer       | Library                                                        |
|-------------|----------------------------------------------------------------|
| Framework   | Expo SDK 54 + Expo Router 6 (file-based routing)               |
| Language    | TypeScript 5.9, strict mode                                    |
| Navigation  | Expo Router — `app/` directory                                 |
| State       | TBD (context + hooks for now)                                  |
| Auth tokens | `expo-secure-store`                                            |
| HTTP        | `src/config/api.ts` — thin `fetch` wrapper                     |
| Fonts       | `@expo-google-fonts/{playfair-display,dm-sans,jetbrains-mono}` |
| Icons       | `@expo/vector-icons` (Lucide-compatible)                       |

---

## Project structure

```
app/                        Expo Router screens
  _layout.tsx               Root layout — loads fonts, shows splash
  index.tsx                 Redirect to (auth)/login or (tabs) based on auth
  (auth)/
    _layout.tsx             Stack, no header
    login.tsx               Login / Register screen
  (tabs)/
    _layout.tsx             Bottom tab navigator
    index.tsx               Home
    beans.tsx               Bean library
    gear/
      _layout.tsx           Stack layout for gear screens
      index.tsx             My Gear screen (gear list + stations)
      [id].tsx              Gear item detail / edit screen
      GearSheet.tsx         Bottom sheet — add / edit a gear item
      StationSheet.tsx      Bottom sheet — add / edit a station
assets/                     Images, splash, icon
src/
  api/
    gear.ts                 Gear + station API calls (gearApi)
  components/
    GearIcon.tsx            Icon resolver by gear type_id
  config/
    api.ts                  API_URL constant + apiFetch helper
  context/
    GearContext.tsx         Gear + station state provider
    UserContext.tsx         Auth / user state provider
  lib/
    apiClient.ts            authorizedFetch — token refresh logic
    auth.ts                 Secure-store token helpers
  theme/
    colors.ts               Design system color tokens → RN values
    typography.ts           Font families, sizes, line heights, text styles
    spacing.ts              Spacing scale, border radii, shadow presets
    index.ts                Re-exports everything from theme/
scripts/
  setup.sh                  One-shot dev environment setup (run after clone)
```

---

## Design system

Design tokens live in `../design/design-system/` — use `src/theme/` which is the RN translation.

- **Colors**: `src/theme/colors.ts` — mirrors `../design/design-system/colors_and_type.css`
- **Typography**: `src/theme/typography.ts` — Playfair Display (display), DM Sans (body), JetBrains Mono (data)
- **Spacing / Radii / Shadows**: `src/theme/spacing.ts`
- **Brand & voice**: `../design/design-system/README.md`
- **Component specs**: `../design/design_handoff_login/`, `../design/design_handoff_profile/`, `../design/design_handoff_my_gear/`
  - Each folder has an HTML prototype (`*.html`) and a `README.md` spec. Always read the README before implementing a screen.

Always `import { colors, textStyles, spacing, radii, shadows } from '@/src/theme'` — never hardcode hex values.

---

## Backend API

Base URL is read from `EXPO_PUBLIC_API_URL` env var (falls back to `http://localhost:8080`).
Configured in `app.config.ts` → `extra.apiUrl` → consumed by `src/config/api.ts`.

Known endpoints (backend default port 8080):

| Method | Path                      | Auth   | Notes                                  |
|--------|---------------------------|--------|----------------------------------------|
| POST   | /auth/register            | public | First name, last name, email, password |
| POST   | /auth/login               | public | Returns access + refresh tokens        |
| POST   | /auth/refresh             | public | Rotate refresh token                   |
| GET    | /user/me                  | bearer | Current user profile                   |
| PATCH  | /user/me                  | bearer | Update profile                         |
| POST   | /user/change-password     | bearer | Change password                        |
| GET    | /health                   | public | Liveness check                         |
| GET    | /gear                     | bearer | List all gear items                    |
| POST   | /gear                     | bearer | Create a gear item                     |
| PUT    | /gear/:id                 | bearer | Update a gear item                     |
| DELETE | /gear/:id                 | bearer | Delete a gear item (returns 204)       |
| GET    | /stations                 | bearer | List all stations                      |
| POST   | /stations                 | bearer | Create a station                       |
| PUT    | /stations/:id             | bearer | Update a station                       |
| DELETE | /stations/:id             | bearer | Delete a station (returns 204)         |

Use `apiFetch` from `src/config/api.ts` for all requests. Pass the JWT as `token` option.

---

## Backend host configuration

Edit `.env` (copy from `.env.example` on first clone):

```
EXPO_PUBLIC_API_URL=http://192.168.1.42:8080   # LAN IP for physical device
```

Restart `expo start` after changing `.env`. On simulators `localhost` works; on a physical device use your machine's LAN IP.

---

## Getting started

```bash
bash scripts/setup.sh   # install Node, deps, copy .env — run once after clone
npm run ios             # iOS Simulator
npm run android         # Android emulator
npm run web             # Browser (Metro + React Native Web)
npm start               # Expo Go / QR code
```

---

## Conventions

- Sentence case for all UI labels, buttons, headings (Title Case only for "Mr. Bean")
- Numbers with units: `18.5 g`, `28 s`, `1:2.1` — always use JetBrains Mono (`textStyles.mono`)
- Lucide icons via `@expo/vector-icons/Feather` (closest stroke-weight match)
- No hardcoded colours or font names — always go through `src/theme`
- `StyleSheet.create` for all styles; no inline style objects in JSX
- Auth screens live outside `(tabs)` — gate in `app/index.tsx` by checking auth state

### Safe area

`SafeAreaProvider` lives in `app/_layout.tsx` (root). Every screen must handle its own insets — the navigators do **not** do this automatically when `headerShown: false`.

| Pattern                                       | When to use                                                                             |
|-----------------------------------------------|-----------------------------------------------------------------------------------------|
| `<SafeAreaView edges={['top']} …>`            | Simple screens with a plain `View` root (Home, Beans, etc.)                             |
| `<SafeAreaView edges={['top', 'bottom']} …>`  | Auth screens inside `KeyboardAvoidingView` (wrap KAV, not replace it)                  |
| `useSafeAreaInsets()` → apply `insets.top`    | Screens with an absolute-positioned overlay (e.g. toast) where SAV can't be outermost  |

Never hardcode status-bar offsets. Import from `react-native-safe-area-context` (already in deps).
