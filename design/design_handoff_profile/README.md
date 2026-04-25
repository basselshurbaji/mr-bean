# Handoff: Profile

## Overview
The Profile screen for **Mr. Bean**. Accessible via the last tab in the bottom nav (replacing Stats). Lets the user view their name and email, edit their first/last name, change their password, and log out.

---

## About the Design Files
`Profile.html` is a **high-fidelity interactive prototype** built in plain HTML + React. It is a design reference — not production code. Recreate it in your React Native / Expo project using your established patterns.

Open `Profile.html` in any browser to interact with the full flow before implementing.

---

## Fidelity
**High-fidelity.** Implement colours, typography, spacing, and interactions as close to this as your platform allows.

---

## Prototype Behaviour (important — read before implementing)

| Behaviour in prototype | What to do in production |
|------------------------|--------------------------|
| Name saves to `localStorage` and updates the avatar initials + display name instantly | Call your update-profile API, then update local state on success |
| Password form validates length ≥ 8 and match, then shows a toast — no real "current password" check | Send `{ currentPassword, newPassword }` to your auth provider; show API error if current password is wrong |
| Log out clears `localStorage` and shows a static goodbye screen | Clear auth session/tokens, navigate back to the Login screen |
| "Sign back in" button on the logout screen resets the prototype state | In production this should navigate to Login |
| User data seeded as `Ada Lovelace / ada@example.com` if nothing in localStorage | Pull from your auth context / user store |
| Closing (Cancel) a section re-mounts it via React key trick — fields reset to saved values | Standard behaviour — just reset local form state on cancel |

---

## Screen: Profile

### Bottom nav
Profile tab is the 4th and final tab (replaces Stats). Nav order:
1. Home
2. My Gear
3. Beans
4. **Profile** ← active on this screen

Same bottom nav style as all other screens (see My Gear handoff for full spec).

### Layout
- `ScrollView`, top padding 12 pt, bottom padding 32 pt
- Horizontal card padding: 20 pt

---

### Header
- Title: "Profile" — Playfair Display 34 pt / weight 900 / `#1C0F07` / tracking −0.7 pt
- Padding: 12 pt top, 24 pt horizontal

---

### Avatar block (centred)
- Circle: 72 × 72 pt, radius 36 pt, background `#1C0F07`, shadow `0 4px 16px rgba(28,15,7,0.2)`
- Initials: first letter of first name + first letter of last name, uppercased
  - DM Sans 26 pt / weight 700 / `#FDF8F2`
- Full name below: DM Sans 18 pt / weight 600 / `#1C0F07`, margin top 14 pt
- Email below name: DM Sans 13 pt / `#8C5340`, margin top 4 pt

---

### Expandable sections

Both sections follow the same card pattern:

**Card container**
- Margin: 0 20 pt 12 pt
- Background `#FAF3E8`, radius 20 pt, padding 18 pt 20 pt
- Shadow: `0 2px 8px rgba(28,15,7,0.07)`

**Collapsed state (row)**
- Left: section label (uppercase 10 pt / weight 700 / tracking 0.08em / `#8C5340`) + preview value below it
- Right: action button — height 30 pt, padding 0 14 pt, radius 9999 pt, background `#F0E4CC`, DM Sans 12 pt / weight 600 / `#1C0F07`

**Expanded state**
- Section label only (no preview value)
- Action button changes to "Cancel" — transparent background, border 1.5 pt `#E8D8B8`, color `#8C5340`
- Form fields expand below the header row with 18 pt gap from header
- Cancelling closes the section and resets form fields to the last saved values

---

#### Section 1 — Name

**Collapsed preview:** `{firstName} {lastName}` — DM Sans 15 pt / weight 600 / `#1C0F07`

**Expanded fields:**
- First name + Last name — 2-column grid, equal width, gap 12 pt
- Both use shared input style
- CTA: "Save name" — full-width, height 50 pt, radius 32 pt, `#1C0F07` bg
- Disabled until first name is non-empty

**On save:**
- Update user object in state + storage
- Collapse section
- Show toast: "Name updated ✓"
- Avatar initials update instantly

---

#### Section 2 — Password

**Collapsed preview:** `••••••••` — DM Sans 15 pt / weight 600 / `#1C0F07` / letter-spacing 2 pt

**Expanded fields (in order):**
1. Current password (type: password, placeholder: `Your current password`)
2. New password (type: password, placeholder: `Min. 8 characters`)
3. Confirm new password (type: password, placeholder: `Same again`)

**Validation:**
- New password must be ≥ 8 characters → error: "New password must be at least 8 characters."
- New and confirm must match → error: "Passwords don't match."
- Error shown in `#FDECEA` container (same style as Login error)

**CTA:** "Update password" — disabled until all three fields are non-empty
**On save:** collapse section, reset all three fields, show toast: "Password updated ✓"

---

### Log out button
- Margin: 8 pt top, 20 pt horizontal, 32 pt bottom
- Height 50 pt, radius 32 pt
- Background `#FDECEA`, border 1.5 pt `#FDECEA`, color `#C0392B`
- DM Sans 15 pt / weight 600
- Hover: background `#f5c6c2`
- **On tap:** clear auth session + navigate to Login screen

---

### Toast
- Position: absolute, bottom 104 pt, inset 20 pt horizontal
- Background `#1C0F07`, color `#FDF8F2`, radius 16 pt, padding 14 pt 18 pt
- DM Sans 14 pt / weight 500 / centred
- Slides up on appear: `translateY(100%→0) + opacity(0→1)`, 260 ms ease-out
- Auto-dismiss after 2500 ms

---

## Design Tokens

### Colours
(same as Login — see Login README for full table)

Key values:
| Token | Hex |
|-------|-----|
| App background | `#FDF8F2` |
| Card background | `#FAF3E8` |
| Subtle | `#F0E4CC` |
| Primary text | `#1C0F07` |
| Secondary text | `#8C5340` |
| Label | `#6B3A2A` |
| Border default | `#E8D8B8` |
| Border focus | `#D4872A` |
| Error bg | `#FDECEA` |
| Error text | `#C0392B` |

### Input style
```
height: 50pt
padding: 0 16pt
background: #FAF3E8
border: 1.5pt solid #E8D8B8
border-radius: 14pt
font: DM Sans 15pt / color #1C0F07
focus: border #D4872A + shadow 0 0 0 3pt rgba(212,135,42,0.12)
```

---

## React Native / Expo Notes

### Navigation
- Profile is tab index 3 (0-based) in the bottom tab navigator.
- Tab icon: person/user outline (`M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2 M12 11a4 4 0 100-8 4 4 0 000 8z`) or use SF Symbol `person` on iOS.
- Log out should call `router.replace('/login')` (Expo Router) or `navigation.reset` to clear the stack.

### Expandable sections
- Manage `open` state per section with `useState`.
- Animate height with `LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut)` before setting state — simple and native-feeling.
- Alternatively use `react-native-reanimated` with `useAnimatedStyle` + `withTiming` for more control.
- On cancel: reset the local form state to the current saved user values.

### Keyboard
- Same as Login: `KeyboardAvoidingView` + `ScrollView`.
- Scroll the expanded section into view when it opens — use `scrollTo` on the `ScrollView` ref.

### Password fields
- `secureTextEntry={true}` on all three password inputs.
- Consider adding a show/hide toggle (eye icon) — not in the prototype but expected UX.

### API calls
- **Update name:** `PATCH /users/me` `{ firstName, lastName }` — update local state only on success.
- **Change password:** `POST /auth/change-password` `{ currentPassword, newPassword }` — handle 401 (wrong current password) by showing the error in the error container.
- **Log out:** Invalidate refresh token server-side, clear `SecureStore`, reset navigation.

### Avatar
- Initials are generated client-side from the user object.
- If you add photo upload later, this component becomes a `<Pressable>` that opens an image picker.

### Fonts
- `@expo-google-fonts/playfair-display` + `@expo-google-fonts/dm-sans`

---

## Files
| File | Description |
|------|-------------|
| `Profile.html` | Full interactive prototype — open in any browser |
