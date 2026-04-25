# Handoff: Your Beans

## Overview

"Your Beans" is the personal bean catalogue section of **Mr. Bean**. Shown to authenticated users, it lets them build a library of coffee products they use or have used — each entry representing a specific roast from a specific roaster, not a purchase or a bag. From the list screen users can browse their catalogue, tap into a bean's detail, and add, edit, or remove entries. The Beans tab lives in the main bottom navigation alongside My Gear and Stats.

---

## About the Design Files

`My Beans.html` is a **high-fidelity interactive prototype built in plain HTML + React (Babel)**. It is a design reference — not production code. Your task is to **recreate these designs in your target codebase** using its established patterns, component library, and framework (React Native, SwiftUI, Expo, etc.). Do not ship the HTML directly.

The prototype is self-contained and interactive — open it in a browser to explore the full flow before implementing.

---

## Fidelity

**High-fidelity.** Pixel-accurate colours, typography, spacing, icons, and interactions. Implement as close to this as your platform allows. All measurements below are from the 390 × 844 pt iPhone canvas (1× logical points).

---

## Prototype Behaviour

| Behaviour in prototype                                                           | What to do in production                                                        |
|----------------------------------------------------------------------------------|---------------------------------------------------------------------------------|
| Bean list seeds from `INIT_BEANS` on first load, then persists in `localStorage` | Fetch from `GET /beans` (ordered by `created_at ASC`) on screen mount           |
| Add bean saves locally and appends to list immediately                           | `POST /beans`, then append response to local state                              |
| Edit bean saves locally and updates in-place                                     | `PUT /beans/:id` (full body, no partial updates), update local state on success |
| Delete bean removes locally, navigates back to list                              | `DELETE /beans/:id`, remove from local state on 204                             |
| Toast message appears for 2.6 s then disappears                                  | ✓ Keep this behaviour                                                           |
| Tapping backdrop dismisses the sheet                                             | ✓ Keep this behaviour                                                           |
| Process and roast level chip pickers toggle off on second tap (deselect)         | ✓ Keep — both fields are optional; `null` means field is unset                  |
| Other nav tabs show a placeholder screen                                         | Wire to real tab screens                                                        |

---

## Screen: Your Beans (list)

### Layout

- Full-screen scroll view with bottom nav
- Top padding: 12 pt; horizontal padding: 24 pt for header, 20 pt for card list
- Fixed bottom navigation bar: 88 pt height, frosted glass (`rgba(253,248,242,0.94)`, `backdrop-filter: blur(20px)`)

### Header

- Title: "Your beans" — Playfair Display 34 pt / weight 900 / tracking −0.7 pt / `#1C0F07`
- Subtitle: "{n} beans in your library" — DM Sans 13 pt / `#8C5340`
- No segment control; no filter chips — this is a single flat list

### Bean cards (list)

- Container: `flex-direction: column`, gap 10 pt, padding `0 20pt`
- Card: background `#FAF3E8`, border-radius 20 pt, padding `14pt 16pt`, shadow `0 2px 8px rgba(28,15,7,0.07)`, no border
- Press state: `scale(0.98)`, 150 ms spring

Card anatomy (horizontal row, left to right):

1. **Roast level bubble** — 50 × 50 pt, border-radius `50 × 0.38` ≈ 19 pt. Background colour is the roast level colour (see Roast Level Colour System). Contains a centred coffee bean SVG mark (23 pt wide) in `#FDF8F2`. When no roast level is set: background `#D8C4A0`, icon colour `#8C5340`.
2. **Text block** (flex 1, min-width 0)
   - Name: DM Sans 15 pt / weight 600 / `#1C0F07`
   - Roaster · Origin: DM Sans 12 pt / `#8C5340` — separated by ` · `; omitted if both are empty
   - Tasting notes preview: DM Sans 12 pt / `#B07060` / single line, truncated with ellipsis
3. **Right meta** (flex-shrink 0, align right)
   - Process badge: background `#F0E4CC`, color `#4A2518`, 10 pt / 600, pill shape
   - Roast level label: DM Sans 11 pt / `#C4A882`

### Empty state

- Centred, 52 pt top padding
- 🫘 emoji (40 pt)
- Title: "No beans on deck." — DM Sans 16 pt / weight 600 / `#1C0F07`
- Body: "Time to fix that, don't you think?" — DM Sans 13 pt / `#B07060`

### FAB (floating action button)

- Position: absolute, bottom 100 pt, right 22 pt
- Size: 56 × 56 pt, border-radius 28 pt
- Background: `#1C0F07`, icon: `+` 24 pt stroke / `#FDF8F2`
- Shadow: `0 6px 20px rgba(28,15,7,0.38)`
- Press: `scale(0.93)`, 220 ms spring (`cubic-bezier(0.34,1.56,0.64,1)`)

---

## Screen: Bean detail

Navigated to by tapping any bean card. Pushes in with a slide-from-right animation (opacity 0→1, translateX 20→0 pt, 260 ms ease-out).

### Nav row

- Back button (left): 38 × 38 pt circle, background `#F0E4CC`, left-arrow icon 18 pt
- Edit button (right): 34 pt height pill, background `#F0E4CC`, "Edit" — DM Sans 13 pt / 600 / `#1C0F07`

### Hero block — padding `20pt 24pt 24pt`

- **Roast bubble**: 72 × 72 pt, border-radius `72 × 0.38` ≈ 27 pt. Same colour system as list card.
- **Badge row** (shown before name): process badge + roast level badge side by side, gap 6 pt
  - Process badge: background `#F0E4CC`, color `#4A2518`
  - Roast level badge: background = roast level colour, color `#FDF8F2`
- **Name**: Playfair Display 26 pt / weight 700 / `#1C0F07` / line-height 1.15
- **Roaster**: DM Sans 14 pt / `#6B3A2A`, margin-top 5 pt
- **Origin**: DM Sans 13 pt / `#B07060`, margin-top 2 pt

### Tasting notes card (shown only if `tasting_notes` is set)

- Background `#FAF3E8`, border-radius 20 pt, padding `18pt 20pt`, shadow `0 2px 8px rgba(28,15,7,0.07)`
- Section label: "TASTING NOTES" — DM Sans 10 pt / weight 700 / uppercase / tracking 0.08em / `#8C5340`
- Content: parse `tasting_notes` by splitting on `, ` — render each descriptor as a small pill chip: background `#F0E4CC`, color `#4A2518`, 12 pt / 600, padding `4pt 10pt`, border-radius `9999pt`

### Personal notes card (shown only if `notes` is set)

- Same card style as tasting notes card
- Section label: "NOTES"
- Content: plain text, DM Sans 14 pt / `#1C0F07` / line-height 1.65

### Remove button

- Full width, height 48 pt, border-radius 32 pt
- Background `#FDECEA`, border `1.5pt solid #FDECEA`, color `#C0392B`
- DM Sans 15 pt / 600
- Label: "Remove this bean"
- Hover: background `#f5c6c2`

---

## Add / Edit Bean Sheet (bottom sheet)

Triggered by the FAB (add mode) or the Edit button on the detail screen (edit mode).

### Sheet container

- `translateY(100%) → translateY(0)`, 300 ms `cubic-bezier(0.33,1,0.68,1)`
- Background `#FDF8F2`, border-radius `28pt 28pt 0 0`
- Drag handle: 36 × 4 pt, `#D8C4A0`, centred, margin 14 pt top
- Backdrop: `rgba(28,15,7,0.42)` — tap to dismiss

### Header

- Title: "Add a bean" / "Edit bean" — Playfair Display 24 pt / weight 700 / `#1C0F07`
- Close button (top right): 34 × 34 pt circle, background `#F0E4CC`, × icon 14 pt / `#6B3A2A`

### Fields

All fields use the shared input style (see Design Tokens below).

| # | Field         | Control     | Required | Notes                                                                                                |
|---|---------------|-------------|----------|------------------------------------------------------------------------------------------------------|
| 1 | Name          | Text        | Yes (\*) | Placeholder: "e.g. Ethiopia Yirgacheffe"                                                             |
| 2 | Roaster       | Text        | No       | Left column of a 2-column grid (gap 12 pt)                                                           |
| 3 | Origin        | Text        | No       | Right column of same grid; placeholder "e.g. Ethiopia"                                               |
| 4 | Process       | Chip picker | No       | 5 pill chips (Washed, Natural, Honey, Anaerobic, Other). Single-select; tap again to deselect.       |
| 5 | Roast level   | Chip picker | No       | 5 pill chips with coloured dot. See Roast Level Colour System. Single-select; tap again to deselect. |
| 6 | Tasting notes | Textarea    | No       | 2 rows; placeholder "e.g. Jasmine, bergamot, peach."                                                 |
| 7 | Notes         | Textarea    | No       | 2 rows; placeholder "Your impressions, ratios, anything worth remembering."                          |

Required field marker: `*` in `#C4782A` next to label.

### Chip picker — process

Chips: flex-wrap row, gap 8 pt.

- Inactive: transparent background, border `1.5pt #E8D8B8`, color `#6B3A2A`, 13 pt / 600
- Active: background `#1C0F07`, border `#1C0F07`, color `#FDF8F2`
- Height 32 pt, padding `0 14pt`, border-radius 9999 pt
- Transition: all 150 ms

### Chip picker — roast level

Same style as process chips, but each chip has a 10 × 10 pt coloured dot (left of label, gap 6 pt):

- Inactive: dot is the roast level colour; chip has transparent background
- Active: dot turns `#FDF8F2`; chip has `#1C0F07` background

### CTA button

- "Add to your beans" / "Save changes"
- Full width, height 54 pt, border-radius 32 pt, background `#1C0F07`, DM Sans 16 pt / 600 / `#FDF8F2`
- Disabled (opacity 0.38) until `name` field is non-empty

---

## Interactions & Behaviour

| Trigger                          | Action                                                                       |
|----------------------------------|------------------------------------------------------------------------------|
| Tap bean card                    | Navigate to Bean detail (push animation)                                     |
| Tap Edit on detail               | Open Add/Edit sheet pre-filled with bean data                                |
| Tap back on detail               | Return to list                                                               |
| Tap FAB                          | Open Add sheet (empty form)                                                  |
| Save bean (add)                  | Append to list, show toast "{name} added ✓"                                  |
| Save bean (edit)                 | Update in list, show toast "{name} updated ✓"                                |
| Remove this bean (detail)        | Delete from list, navigate back, show toast "Removed from your beans."       |
| Tap backdrop                     | Dismiss sheet                                                                |

### Toast

- Position: absolute, bottom 104 pt, horizontal inset 20 pt
- Background `#1C0F07`, color `#FDF8F2`, border-radius 16 pt, padding `14pt 18pt`
- DM Sans 14 pt / 500 / centred
- Slides up on appearance (same animation as sheet)
- Auto-dismiss after 2600 ms

---

## Design Tokens

Only tokens new or specific to this feature are listed. All standard tokens follow `design-system/colors_and_type.css`.

### Roast level colour system

Maps roast levels to existing design system palette tokens.

| Roast level    | Display label | Hex       | Token           |
|----------------|---------------|-----------|-----------------|
| `light`        | Light         | `#D4872A` | `caramel-400`   |
| `medium_light` | Medium light  | `#9B5E1A` | `caramel-600`   |
| `medium`       | Medium        | `#6B3A2A` | `espresso-500`  |
| `medium_dark`  | Medium dark   | `#4A2518` | `espresso-600`  |
| `dark`         | Dark          | `#1C0F07` | `espresso-800`  |
| (unset)        | —             | `#D8C4A0` | `cream-500`     |

Used on: roast bubble background (list + detail), roast level badge background (detail), coloured dot in roast chip picker.

### Colours

| Role                              | Hex                     |
|-----------------------------------|-------------------------|
| App background                    | `#FDF8F2`               |
| Card background                   | `#FAF3E8`               |
| Subtle background / chip inactive | `#F0E4CC`               |
| Primary text                      | `#1C0F07`               |
| Secondary text                    | `#6B3A2A`               |
| Tertiary / muted text             | `#8C5340`               |
| Faint text (tasting note preview) | `#B07060`               |
| Badge text (dark-on-cream)        | `#4A2518`               |
| Accent / required marker          | `#C4782A`               |
| Border default                    | `#E8D8B8`               |
| Border focus                      | `#D4872A`               |
| Focus ring                        | `rgba(212,135,42,0.12)` |
| Error background                  | `#FDECEA`               |
| Error text / destructive          | `#C0392B`               |

### Input style

```
height:        50 pt
padding:       0 16 pt
background:    #FAF3E8
border:        1.5 pt solid #E8D8B8
border-radius: 14 pt
font:          DM Sans 15 pt / color #1C0F07
placeholder:   #C4A882
focus:         border #D4872A + box-shadow 0 0 0 3pt rgba(212,135,42,0.12)

textarea:
  padding:     14pt 16pt
  rows:        2
  resize:      none
  line-height: 1.5
```

### Typography

| Use                         | Family           | Size  | Weight |
|-----------------------------|------------------|-------|--------|
| Screen title ("Your beans") | Playfair Display | 34 pt | 900    |
| Detail bean name            | Playfair Display | 26 pt | 700    |
| Sheet title                 | Playfair Display | 24 pt | 700    |
| Card bean name              | DM Sans          | 15 pt | 600    |
| Card secondary (roaster)    | DM Sans          | 12 pt | 400    |
| Card tasting preview        | DM Sans          | 12 pt | 400    |
| Field label                 | DM Sans          | 13 pt | 500    |
| Input value                 | DM Sans          | 15 pt | 400    |
| Chip label                  | DM Sans          | 13 pt | 600    |
| Badge                       | DM Sans          | 10 pt | 600    |
| Section label (detail card) | DM Sans          | 10 pt | 700    |
| CTA button                  | DM Sans          | 16 pt | 600    |

---

## React Native / Expo Notes

### Navigation

- Beans is a **tab** in the bottom tab bar. Inside the tab, use a native stack:
  - `BeansScreen` (index) — list view
  - `BeanDetailScreen` — push on card tap
- The Add/Edit sheet is a modal overlay, not a stack screen. Use **`@gorhom/bottom-sheet`** with `snapPoints={['92%']}`.

### List

- Use `FlatList` for the bean list — not `ScrollView + .map()` — for performance on larger catalogues.
- `keyExtractor` by `bean.id`.

### Chip pickers

- Both chip pickers (process and roast level) are plain `Pressable` rows in a `View` with `flexWrap: 'wrap'`.
- No library needed. Each chip toggles a local state string (`null` when deselected).

### Fonts

- `@expo-google-fonts/playfair-display` + `@expo-google-fonts/dm-sans`
- Load with `useFonts()` hook and hold the splash screen until ready.

### Bean mark SVG

- Port the inline SVG to `react-native-svg` (`Svg`, `Ellipse`, `Path`).
- Wrap in a `<RoastBubble roastId={bean.roast_level} size={50} />` helper that switches background colour per the colour system table above.

### Tasting note chips (detail)

- Split `tasting_notes` on `/,\s*/` in the render function, strip trailing period.
- Render each chip as a `Text` in a `View` with `flexWrap: 'wrap'` — no library needed.

### Keyboard handling

- Wrap the Add/Edit sheet form in `KeyboardAvoidingView` (`behavior="padding"` on iOS).
- Use `returnKeyType="next"` and `ref` chaining to advance focus field-to-field on Return.

### Data persistence

- Store beans in **`@react-native-async-storage/async-storage`** under key `mrbean_beans` while offline / optimistic.
- On mount: load local cache first, then fetch `GET /beans` and reconcile.
- On write: call the API, update local cache on success.

### Animations

- Sheet slide-up: `@gorhom/bottom-sheet` native gesture handler.
- Card press scale: `react-native-reanimated` with `withSpring`.
- FAB press: `withSpring({ damping: 8, stiffness: 180 })` for the bouncy feel.
- Detail screen push: standard `@react-navigation/native-stack` slide animation.
- Toast: `Animated.timing` with `translateY` + `opacity`, auto-dismiss with `setTimeout`.

### Haptics

- `expo-haptics`: `Haptics.notificationAsync(NotificationFeedbackType.Success)` on save.
- `Haptics.impactAsync(ImpactFeedbackStyle.Light)` on card press.

### Recommended libraries summary

| Purpose      | Library                                         |
|--------------|-------------------------------------------------|
| Navigation   | `expo-router` or `@react-navigation/native`     |
| Bottom sheet | `@gorhom/bottom-sheet`                          |
| SVG icons    | `react-native-svg`                              |
| Fonts        | `@expo-google-fonts/playfair-display` + dm-sans |
| Persistence  | `@react-native-async-storage/async-storage`     |
| Animations   | `react-native-reanimated`                       |
| Blur (tab)   | `expo-blur`                                     |
| Haptics      | `expo-haptics`                                  |

---

## Files

| File            | Description                                      |
|-----------------|--------------------------------------------------|
| `My Beans.html` | Full interactive prototype — open in any browser |
| `README.md`     | This handoff spec                                |
