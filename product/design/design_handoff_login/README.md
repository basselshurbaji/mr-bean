# Handoff: Login & Register

## Overview
The authentication gate for **Mr. Bean**. Shown to users who are not signed in. Two modes — Sign in and Create account — toggled via a segment control on the same screen. On success the user proceeds to the main app.

---

## About the Design Files
`Login.html` is a **high-fidelity interactive prototype** built in plain HTML + React. It is a design reference — not production code. Recreate it in your React Native / Expo project using your established patterns.

Open `Login.html` in any browser to interact with the full flow before implementing.

---

## Fidelity
**High-fidelity.** Implement colours, typography, spacing, and interactions as close to this as your platform allows.

---

## Prototype Behaviour (important — read before implementing)

| Behaviour in prototype | What to do in production |
|------------------------|--------------------------|
| Submitting any valid form shows a success screen inside the prototype | Replace with real API call (sign in / register endpoint) |
| "Forgot password?" button does nothing | Wire to password reset flow |
| User object written to `localStorage` | Use your auth system (JWT, session, Supabase, etc.) |
| No real password validation beyond length ≥ 8 and match | Add server-side validation |
| Confirm password field only shown on register | ✓ Keep this behaviour |
| CTA disabled until all required fields are filled | ✓ Keep this behaviour |

---

## Screen: Login / Register

### Layout
- Full-screen scroll view (no bottom nav — pre-auth)
- Top padding: 40 pt
- Horizontal padding: 24 pt throughout

### Brand hero (top)
- `BeanMark` SVG logo: 48 pt wide
- Title: "Mr. Bean" — Playfair Display 36 pt / weight 900 / `#1C0F07` / tracking −0.8 pt
- Subtitle: "Your espresso, perfected." — DM Sans 14 pt / `#8C5340`
- Centred, text-align center

### Segment control (Sign in / Create account)
- Container: background `#F0E4CC`, border-radius 14 pt, padding 4 pt, margin 0 24 pt 28 pt
- Active segment: background `#FDF8F2`, color `#1C0F07`, shadow `0 1px 4px rgba(28,15,7,0.1)`, radius 10 pt
- Inactive: transparent, color `#8C5340`
- Font: DM Sans 13 pt / weight 600
- Switching mode clears the error state

### Form fields
All fields use the shared input style (see Design Tokens below).

**Sign in mode:**
1. Email (type: email, placeholder: `you@example.com`)
2. Password (type: password, placeholder: `Your password`)
3. "Forgot password?" link below CTA — DM Sans 13 pt / weight 600 / color `#C4782A`, no border

**Create account mode** (additional fields above email):
1. First name + Last name — 2-column grid, equal width, gap 12 pt
2. Email
3. Password (placeholder: `Min. 8 characters`)
4. Confirm password (placeholder: `Same again`)

Required field marker: `*` in `#C4782A` next to label.

### Error state
- Appears between last field and CTA
- Container: background `#FDECEA`, border `1px solid #f5c6c2`, radius 12 pt, padding 10 pt 14 pt
- Text: DM Sans 13 pt / weight 500 / `#C0392B`
- Triggered by: password mismatch, password < 8 chars

### CTA button
- "Sign in" / "Create account"
- Height 54 pt, radius 32 pt, background `#1C0F07`, color `#FDF8F2`
- DM Sans 16 pt / weight 600
- Disabled (opacity 0.38) until:
  - Sign in: email + password filled
  - Register: all fields filled, password ≥ 8 chars, passwords match

---

## Design Tokens

### Colours
| Token | Hex |
|-------|-----|
| App background | `#FDF8F2` |
| Card background | `#FAF3E8` |
| Subtle background | `#F0E4CC` |
| Primary text | `#1C0F07` |
| Secondary text | `#8C5340` |
| Label text | `#6B3A2A` |
| Accent / link | `#C4782A` |
| Focus ring | `rgba(212,135,42,0.12)` |
| Border default | `#E8D8B8` |
| Border focus | `#D4872A` |
| Error background | `#FDECEA` |
| Error border | `#f5c6c2` |
| Error text | `#C0392B` |

### Input style
```
height: 50pt
padding: 0 16pt
background: #FAF3E8
border: 1.5pt solid #E8D8B8
border-radius: 14pt
font: DM Sans 15pt / color #1C0F07
placeholder color: #C4A882
focus: border #D4872A + box-shadow 0 0 0 3pt rgba(212,135,42,0.12)
```

### Typography
| Use | Family | Size | Weight |
|-----|--------|------|--------|
| App title | Playfair Display | 36 pt | 900 |
| Subtitle | DM Sans | 14 pt | 400 |
| Field label | DM Sans | 13 pt | 500 |
| Input value | DM Sans | 15 pt | 400 |
| CTA | DM Sans | 16 pt | 600 |
| Link | DM Sans | 13 pt | 600 |

---

## React Native / Expo Notes

### Navigation
- This screen lives outside the tab navigator — render it conditionally when `user === null` at the root level (e.g. in `app/_layout.tsx` with Expo Router, check auth state and redirect).
- On success: set user in auth context / store, navigate to `/(tabs)/home`.

### Keyboard handling
- Wrap the form in `KeyboardAvoidingView` (`behavior="padding"` on iOS, `behavior="height"` on Android).
- Use `ScrollView` so the form scrolls when keyboard is up.

### Form
- Use `TextInput` with `keyboardType="email-address"` and `autoCapitalize="none"` for email.
- Use `secureTextEntry` for password fields.
- Use `returnKeyType="next"` and `ref` chaining to move focus field-to-field on Return.
- Disable the submit button with `disabled` prop; set `opacity: 0.38` via style.

### Segment control
- Use two `Pressable` buttons in a `View` row — no third-party library needed.
- Animate the active indicator with `Animated.spring` if desired, or just swap styles.

### Fonts
- `@expo-google-fonts/playfair-display` + `@expo-google-fonts/dm-sans`
- Load with `useFonts()` hook; show `SplashScreen` until ready.

### Auth
- Recommended: **Supabase Auth** (`@supabase/supabase-js`) or Firebase Auth.
- Store session in SecureStore (`expo-secure-store`), not AsyncStorage, for tokens.

---

## Files
| File | Description |
|------|-------------|
| `Login.html` | Full interactive prototype — open in any browser |
