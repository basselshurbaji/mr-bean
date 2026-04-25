# Handoff: My Gear — Equipment & Stations

## Overview
"My Gear" is the equipment management section of **Mr. Bean**, a coffee enthusiast app centred on espresso extraction logging. Users register their hardware (machines, grinders, accessories), then group them into **Stations** — named presets that pre-select gear when logging a shot. No extraction is ever linked to a station directly; the station is just a quick-select helper.

---

## About the Design Files
The files in this bundle (`My Gear.html`) are **high-fidelity prototypes built in plain HTML + React (Babel)**. They are design references — not production code. Your task is to **recreate these designs in your target codebase** using its established patterns, component library, and framework (React Native, SwiftUI, Expo, etc.). Do not ship the HTML directly.

The prototype is self-contained and interactive — open it in a browser to explore the full flow before implementing.

---

## Fidelity
**High-fidelity.** Pixel-accurate colours, typography, spacing, icons, and interactions. Implement as close to this as your platform allows. All measurements below are from the 390 × 844 pt iPhone canvas (1× logical points).

---

## Design System — Mr. Bean

### Fonts
| Role | Family | Notes |
|------|--------|-------|
| Display / headings | Playfair Display | Weights 700, 900 |
| Body / UI | DM Sans | Weights 400, 500, 600 |
| Numbers / mono | JetBrains Mono | Weights 400, 600 |

### Colour Tokens
| Token | Hex | Usage |
|-------|-----|-------|
| `espresso-800` | `#1C0F07` | Primary text, dark backgrounds, active states |
| `espresso-700` | `#2D1810` | Hover on dark surfaces |
| `espresso-500` | `#6B3A2A` | Secondary text, labels |
| `espresso-300` | `#B07060` | Tertiary / muted text |
| `caramel-500` | `#C4782A` | Accent links, back buttons |
| `caramel-400` | `#D4872A` | Focus rings, selected indicators |
| `cream-100` | `#FDF8F2` | App background |
| `cream-200` | `#FAF3E8` | Card background |
| `cream-300` | `#F0E4CC` | Subtle background, chip deselected |
| `cream-400` | `#E8D8B8` | Borders (default) |
| `cream-500` | `#D8C4A0` | Borders (stronger), drag handles |
| `error-500` | `#C0392B` | Destructive actions |
| `error-100` | `#FDECEA` | Destructive button background |

### Spacing
4 pt base unit. Common values: 4, 8, 12, 16, 20, 24, 32, 40, 48 pt.

### Border Radius
| Name | Value |
|------|-------|
| sm | 6 pt |
| md | 12 pt |
| lg | 20 pt |
| xl | 32 pt |
| full | 9999 pt (pill) |

### Shadows
`0 2px 8px rgba(28,15,7, 0.07)` — card resting  
`0 4px 16px rgba(28,15,7, 0.10)` — elevated sheet  
`0 -4px 32px rgba(28,15,7, 0.16)` — bottom sheet

### Animation
- Ease out: `cubic-bezier(0.33, 1, 0.68, 1)` — default for entries
- Spring: `cubic-bezier(0.34, 1.56, 0.64, 1)` — FAB, buttons
- Fast: 150 ms | Base: 220 ms | Slow: 300–350 ms

---

## Equipment Types (hardcoded list)
Each type has an id, display label, and a custom stroke SVG icon (see prototype source for exact paths):

| id | Label |
|----|-------|
| `machine` | Espresso machine |
| `grinder` | Grinder |
| `scale` | Scale |
| `portafilter` | Portafilter |
| `tamper` | Tamper |
| `distributor` | Distribution tool |
| `wdt` | WDT tool |
| `basket` | Basket |
| `puckscreen` | Puck screen |
| `other` | Other |

Icons are 24 × 24 pt stroke SVGs, `strokeWidth: 1.8`, `strokeLinecap: round`, `strokeLinejoin: round`. All paths are in the prototype source under `const typeIcons = { ... }`.

---

## Screens / Views

### 1. My Gear — Main Screen

**Header (top of scroll)**
- Title: "My Gear" — Playfair Display 34 pt / weight 900 / tracking −0.7 pt / color `espresso-800`
- Subtitle: "{n} pieces · {n} stations" — DM Sans 13 pt / color `espresso-300`
- Padding: 12 pt top, 24 pt horizontal

**Segment Control (Gear / Stations)**
- Full-width pill container: background `cream-300`, border-radius 14 pt, padding 4 pt
- Each segment: height 36 pt, border-radius 10 pt
- Active: background `cream-100`, color `espresso-800`, shadow `0 1px 4px rgba(28,15,7,0.1)`
- Inactive: transparent background, color `espresso-400`
- Font: DM Sans 13 pt / weight 600

---

### 1a. Gear Tab

**Filter chips** (horizontal scroll, no scrollbar)
- Padding: 0 20 pt, gap 8 pt, bottom padding 16 pt
- Chip height: 32 pt, padding 0 14 pt, border-radius 9999 pt
- Active: background `espresso-800`, border `espresso-800`, color `cream-100`
- Inactive: transparent background, border `cream-400`, color `espresso-500`
- Font: DM Sans 12 pt / weight 600
- "All" chip always first, then one chip per type present in the user's gear

**Gear cards (list)**
- Container: `flex-direction: column`, gap 10 pt, padding 0 20 pt
- Card: background `cream-200`, border-radius 20 pt, padding 14 pt 16 pt
- Shadow: `0 2px 8px rgba(28,15,7,0.07)`
- Press state: `scale(0.98)`, 150 ms spring

Card anatomy (horizontal row):
1. **Icon bubble** — 50 × 50 pt, border-radius 50 × 0.38 ≈ 19 pt, background `cream-300`, icon 25 pt / color `espresso-800`
2. **Text block** (flex 1)
   - Name: DM Sans 15 pt / weight 600 / color `espresso-800`
   - Sub: DM Sans 12 pt / color `espresso-500` — "{brand} · {model}"
3. **Right meta**
   - Type badge: background `cream-300`, color `espresso-700`, 10 pt font, pill shape
   - Year: DM Sans 11 pt / color `cream-600`

**Empty state** (no gear or no matching filter)
- Centered, padding 52 pt top
- 40 pt emoji, 16 pt bold title, 13 pt body in `espresso-300`

**FAB (floating action button)**
- Position: absolute, bottom 100 pt, right 22 pt
- Size: 56 × 56 pt, border-radius 28 pt
- Background: `espresso-800`, icon: `+` 24 pt / `cream-100`
- Shadow: `0 6px 20px rgba(28,15,7,0.38)`
- Press: `scale(0.93)`, 220 ms spring

---

### 1b. Stations Tab

**Station cards**
- Same card style as gear cards (background `cream-200`, radius 20, shadow)
- Padding: 18 pt 18 pt 16 pt
- Row 1: station name (15 pt / 600 / `espresso-800`) + item count (12 pt / `espresso-500`) + chevron
- Row 2: icon strip — up to 7 gear icons, each in a 36 × 36 pt / radius 10 pt / `cream-300` tile. "+N" tile if overflow.
- Press: `scale(0.98)`, 150 ms

**Add Station button**
- Height 54 pt, border-radius 20 pt
- Border: 1.5 pt dashed `cream-500`
- Background: transparent → `cream-300` on hover
- Font: DM Sans 14 pt / weight 600 / `espresso-500`
- Label: "+ New station"

---

### 2. Gear Detail Screen

Navigated to by tapping a gear card. Pushes in with a slide-from-right animation (opacity 0→1, translateX 20→0 pt, 260 ms ease-out).

**Nav row**
- Back button (left): 38 × 38 pt circle, background `cream-300`, chevron-left icon 18 pt
- Edit button (right): 34 pt height pill, background `cream-300`, "Edit" label DM Sans 13 pt / 600 / `espresso-800`

**Hero block** — padding 20 pt 24 pt
- Icon bubble: 72 × 72 pt, border-radius 27 pt, background `espresso-800`, icon 36 pt / `cream-100`
- Type badge: `cream-300` background, `espresso-700` text, pill
- Name: Playfair Display 26 pt / weight 700 / `espresso-800`
- Brand · model: DM Sans 14 pt / `espresso-500`
- Year: DM Sans 13 pt / `espresso-300` — "Acquired {year}"

**Notes card** (shown only if notes exist)
- Background `cream-200`, radius 20 pt, padding 18 pt 20 pt
- Label: 10 pt / weight 700 / uppercase / tracking 0.08em / `espresso-400`
- Body: 14 pt / `espresso-800` / line-height 1.65

**Extractions placeholder card**
- Same card style
- Centred content: label (uppercase 10 pt), emoji, 14 pt bold title, 12 pt subtitle

**Remove button**
- Height 48 pt, radius 32 pt
- Background `error-100`, border `error-100` 1.5 pt
- Color `error-500`, DM Sans 15 pt / 600
- Label: "Remove from my gear"

---

### 3. Add Gear Sheet (Bottom Sheet)

**Step 1 — Type picker**
- 3-column grid of type buttons
- Each button: background `cream-200`, border 1.5 pt `cream-400`, radius 18 pt, padding 16 pt 8 pt 14 pt
- Icon: 30 pt / `espresso-800`
- Label: DM Sans 10 pt / weight 600 / centered / `espresso-800`
- Hover: background `cream-300`, border `cream-500`
- Press: `scale(0.95)`

**Step 2 — Form** (skipped when editing, type locked)
- Back link (editing mode hidden): caramel-400, 13 pt / 600
- Title: Playfair Display 24 pt / 700
- Subtitle: current type label, 13 pt / `espresso-400`

Fields:
- **Name** (required — marked with `*` in caramel): full-width input
- **Brand + Model**: 2-column grid, equal width
- **Year purchased**: full-width, maxLength 4
- **Notes**: textarea, 3 rows

Input style: height 50 pt, padding 0 16 pt, background `cream-200`, border 1.5 pt `cream-400`, radius 14 pt, DM Sans 15 pt / `espresso-800`. Focus: border `caramel-400`, box-shadow `0 0 0 3px rgba(212,135,42,0.12)`.

CTA button: "Add to my gear" / "Save changes" — full-width, height 54 pt, radius 32 pt, background `espresso-800`, DM Sans 16 pt / 600 / `cream-100`. Disabled: opacity 0.38.

---

### 4. Add / Edit Station Sheet (Bottom Sheet)

- Title: "New station" / "Edit station" — Playfair Display 24 pt / 700
- Station name field (same input style as above, required)
- Explanation text: 12 pt / `espresso-300` / italic — "Stations pre-select tools when logging — you can always tweak before saving a shot."

**Gear selection list**
- On open, sort: selected items first (in their original order), then unselected. Order stays fixed for the session — items do NOT re-sort as you tap.
- Each row: same horizontal structure as gear card but interactive toggle
  - Selected: background `espresso-800`, border `espresso-800`; name `cream-100`; sub `cream-600`; icon bubble dark (`espresso-800` bg, `cream-100` icon)
  - Unselected: background `cream-200`, border `cream-400`
  - Toggle indicator: 20 × 20 pt circle, right side. Selected: `caramel-400` fill + white checkmark. Unselected: `cream-500` border, transparent fill.
  - Transition: all 150 ms

**CTAs**
- Primary: "Create station · {n} items" / "Save changes" — same full-width style
- Destructive (edit mode only): "Delete station" — height 46 pt, background `error-100`, border `error-100`, color `error-500`

---

## Interactions & Behaviour

| Trigger | Action |
|---------|--------|
| Tap gear card | Navigate to Gear Detail (push animation) |
| Tap Edit on Gear Detail | Open Add/Edit Gear sheet, pre-filled, step 2 |
| Tap back on Gear Detail | Return to list |
| Tap + FAB | Open Add Gear sheet (step 1) |
| Tap type in step 1 | Advance to step 2 |
| Tap back in step 2 | Return to step 1 |
| Save gear (add) | Append to list, show toast "X added ✓" |
| Save gear (edit) | Update in list, show toast "X updated ✓" |
| Delete gear | Remove from list + remove from all stations, show toast |
| Tap station card | Open Edit Station sheet pre-filled |
| Save station | Upsert in list, show toast |
| Delete station | Remove from list, show toast |
| Tap backdrop | Dismiss sheet |

**Toast**
- Position: absolute, bottom 104 pt, horizontal inset 20 pt
- Background: `espresso-800`, color `cream-100`, radius 16 pt, padding 14 pt 18 pt
- DM Sans 14 pt / 500 / centred
- Auto-dismiss after 2600 ms
- Slides up on appearance (same sheet animation)

**Bottom sheet animation**
- `translateY(100%) → translateY(0)`, 300 ms `cubic-bezier(0.33, 1, 0.68, 1)`
- Backdrop: `rgba(28,15,7, 0.42)`
- Drag handle: 36 × 4 pt, `cream-500`, radius 9999 pt, centred, margin 14 pt top

---

## State Management

```
gear: GearItem[]           // persisted
stations: Station[]        // persisted
activeSheet: null | 'add-gear' | 'add-station' | GearItem | Station
activeTab: 'gear' | 'stations'
filterTypeId: string       // 'all' | typeId
detailItem: GearItem | null
toast: string | null

GearItem {
  id: string
  typeId: string           // one of the 10 type ids
  name: string             // required
  brand?: string
  model?: string
  year?: string
  notes?: string
}

Station {
  id: string
  name: string             // required
  gearIds: string[]        // ordered list of GearItem ids
}
```

Persistence: store `gear` and `stations` in AsyncStorage (React Native) or equivalent. Restore on app launch.

---

## Assets
- All icons are inline SVG paths defined in the prototype. No external icon library needed.
- Fonts: Playfair Display, DM Sans, JetBrains Mono — load from Google Fonts or bundle.

---

## Files
| File | Description |
|------|-------------|
| `My Gear.html` | Full interactive prototype — open in any browser |

---

## React Native / Expo Implementation Notes

### Navigation
- Use **Expo Router** (file-based) or **React Navigation** with a native stack.
- My Gear lives as a tab in the bottom tab bar. Recommended stack inside the tab:
  - `GearScreen` (index) → `GearDetailScreen` (push)
- Sheets (Add Gear, Add/Edit Station) should use **`@gorhom/bottom-sheet`** — it gives native gesture handling, snap points, and backdrop. Set `snapPoints={['92%']}` to match the prototype's `max-height: 92%`.

### Icons
- Port the inline SVGs to **`react-native-svg`** (`Svg`, `Path`, `Circle`, `Line` components).
- Wrap each icon in a helper: `<GearIcon typeId="machine" size={24} color="#1C0F07" />` that switches on `typeId` and renders the matching SVG.

### Fonts
- Load via **`expo-font`** or `@expo-google-fonts/playfair-display` + `@expo-google-fonts/dm-sans` + `@expo-google-fonts/jetbrains-mono`.
- Use `useFonts()` hook and show a splash until loaded.

### Persistence
- Use **`@react-native-async-storage/async-storage`** for gear and stations.
- Store as JSON strings under keys `mrbean_gear` and `mrbean_stations`.
- Load on app init; write on every mutation.

### Layout
- All measurements are in logical points (same as pt on iOS, dp on Android) — no conversion needed.
- Use `StyleSheet.create` for performance. Avoid inline styles in list items.
- `FlatList` for the gear list and station gear picker — do not use `ScrollView` + `.map()` for long lists.
- Filter chips: `FlatList` with `horizontal={true}` and `showsHorizontalScrollIndicator={false}`.

### Bottom Tab Bar
- Custom tab bar component to match the design (transparent frosted glass, 88 pt height with home indicator space).
- Use `BlurView` from `expo-blur` for the backdrop blur effect.

### Animations
- Sheet slide-up: handled by `@gorhom/bottom-sheet` natively.
- Card press scale (`scale(0.98)`): use `Animated.spring` or **`react-native-reanimated`** with `withSpring`.
- FAB press: `withSpring({ damping: 8, stiffness: 180 })` for the bouncy feel.
- Toast: slide up from bottom using `Animated.timing` with `translateY` + `opacity`.

### Haptics
- Add `expo-haptics` on gear save / delete: `Haptics.notificationAsync(NotificationFeedbackType.Success)`.
- Light impact on card press: `Haptics.impactAsync(ImpactFeedbackStyle.Light)`.

### Platform notes
- The prototype uses `espresso-800` (`#1C0F07`) for the status bar — set `<StatusBar style="dark" />` (Expo) or match with `SystemUI`.
- On Android, ensure the bottom sheet clears the navigation bar — use `react-native-edge-to-edge` or `WindowInsetsCompat`.
- The `__editGear` sheet routing hack in the prototype: implement with a proper navigation param or a dedicated `useGearSheet` context/store.

### Recommended libraries summary
| Purpose | Library |
|---------|---------|
| Navigation | `expo-router` or `@react-navigation/native` |
| Bottom sheets | `@gorhom/bottom-sheet` |
| SVG icons | `react-native-svg` |
| Fonts | `@expo-google-fonts/*` |
| Persistence | `@react-native-async-storage/async-storage` |
| Animations | `react-native-reanimated` |
| Blur (tab bar) | `expo-blur` |
| Haptics | `expo-haptics` |
