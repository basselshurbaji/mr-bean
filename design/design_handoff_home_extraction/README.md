# Handoff: Home Screen + Extraction Flow

## Overview

This package covers the **Home screen** and **Extraction modal** of the Mr. Bean app — the core feature from which all logging, troubleshooting, and AI predictions flow. The home screen acts as a daily entry point: a warm invitation card that opens a full-height extraction modal where the user times or manually logs a shot, fills in brew parameters, and saves the extraction.

---

## About the Design Files

The files in this bundle (`Home Screen.html`, `Home Screen v1.html`, `Home Screen v2.html`) are **HTML design references** — interactive prototypes built to show intended look, layout, feel, and behavior. They are not production code to copy directly.

Your task is to **recreate these designs in the existing React Native / Expo codebase** (`mobile/`) using its established patterns (Expo Router, StyleSheet, Animated, Context, etc.). Use the HTML as a visual and behavioral spec, not as source.

The existing codebase already has:
- `GearContext` + `gearApi` — gear items and stations
- `BeansContext` + `beansApi` — bean library
- `UserContext` — auth
- Bottom tab navigation in `app/(tabs)/_layout.tsx`
- The home tab stub at `app/(tabs)/index.tsx` (currently empty — this is what you're building)

---

## Fidelity

**High-fidelity.** Colors, typography, spacing, corner radii, shadows, and interactions should be reproduced precisely. The HTML prototypes use the Mr. Bean Design System tokens (`colors_and_type.css`) — map them directly to the values in `src/theme/`.

---

## Screens / Views

### 1. Home Screen (`app/(tabs)/index.tsx`)

The home screen has three stacked sections inside a `ScrollView`:

#### 1a. Header

```
paddingTop: safeAreaInsets.top + 12
paddingHorizontal: 24
paddingBottom: 24
```

- **Greeting line** — `"Good morning"` (or time-appropriate variant)
  - Font: DM Sans 500, 13px, color `espresso400` (`#8C5340`)
  - marginBottom: 4
- **Headline** — `"Ready to pull"` + italic `"a shot?"` on second line
  - Font: Playfair Display 900, 32px, letterSpacing: -0.8, lineHeight: 1.1
  - `"Ready to pull"` color: `espresso800` (`#1C0F07`)
  - `"a shot?"` color: `caramel500` (`#C4782A`), italic
- **Bean mark logo** — top-right, 34×44px SVG (see `assets/logo-icon.svg`). Sits `alignItems: flex-start` on the row.

#### 1b. Extraction Invitation Card

A dark card that opens the extraction modal on tap. Not a form — just an invitation.

```
marginHorizontal: 18
background: #1C0F07 (espresso800)
borderRadius: 28
padding: 22
shadow: 0 8px 40px rgba(28,15,7,0.20)
overflow: hidden
```

**Decorative rings** (non-interactive, purely visual):
- Two concentric partial circles in the top-right corner using `position: absolute`, `border: 16–18px solid rgba(255,255,255,0.03)`, borderRadius full. One at `right: -40, top: -40, size: 180`. One at `right: 16, top: 16, size: 100`.

**Bean pill** (top-left):
```
background: espresso700 (#2D1810)
borderRadius: full
paddingVertical: 5, paddingLeft: 8, paddingRight: 12
marginBottom: 18
flexDirection: row, alignItems: center, gap: 6
```
- Small bean shape: `width:12, height:15, borderRadius: "50% 50% 50% 50% / 60% 60% 40% 40%"` (teardrop), background: `caramel400`
- Text: bean name from last extraction or selected bean. Font: DM Sans 700, 12px, color `cream500` (`#D8C4A0`)
- This pill is purely display — tapping the card opens the modal where the user selects beans.

**Headline:**
```
fontFamily: Playfair Display 700
fontSize: 28
color: cream100 (#FDF8F2)
lineHeight: 1.2
letterSpacing: -0.4
marginBottom: 6
```
Copy: `"Pull a shot."`

**Sub-line** (last shot summary, if exists):
```
fontSize: 13, color: espresso400 (#8C5340)
marginBottom: 20
```
Copy example: `"Last: 27.4s · 1:2.0"`

**Start button** (inside the card, not floating):
```
height: 50
borderRadius: full
background: matcha500 (#4A7C59)
flexDirection: row, alignItems: center, justifyContent: center, gap: 8
shadow: 0 4px 16px rgba(74,124,89,0.40)
```
- Play icon (triangle, filled white, 18×18)
- Label: `"Start extraction"` — DM Sans 700, 15px, white

The entire card is a `Pressable` with `scale(0.98)` on press. Tapping opens the Extraction Modal.

#### 1c. Recent Extractions

```
paddingHorizontal: 20
paddingBottom: 32
```

**Section header row:**
- Left: `"Recent extractions"` — Playfair Display 700, 20px, espresso800
- Right: `"All"` link with chevron → navigates to All Extractions screen (not yet built). Font: DM Sans 700, 13px, color `matcha500` (`#4A7C59`)
- `paddingTop: 32, paddingBottom: 12`

**Extraction cards** (see RecentCard component below). Gap: 10px between cards. Show most recent 3.

---

### 2. Extraction Modal (Bottom Sheet)

Opens on tap of the home card. A full-height bottom sheet (`height: 92%` of screen), slides up with spring animation.

```
background: cream100 (#FDF8F2)
borderRadius: 32 32 0 0
shadow: 0 -8px 48px rgba(28,15,7,0.18)
```

**Backdrop:** `rgba(28,15,7,0.50)` — tapping it closes the modal (only when not mid-extraction).

**Handle:** `width:38, height:4, background: cream500, borderRadius: full, margin: 14 auto 0`

#### Modal Header

```
paddingHorizontal: 22
paddingTop: 16
flexDirection: row, justifyContent: space-between, alignItems: center
```
- Title: `"New extraction"` — Playfair Display 700 italic, 22px, espresso800
- Close button (× icon): `background: cream300, borderRadius: full, width:32, height:32`. Hidden while extraction is running.

#### Bean Selector

Full-width pressable pill below the header:
```
marginTop: 16
marginHorizontal: 20
padding: 12 16
borderRadius: 16
background: espresso800 (selected) OR cream200 (none selected)
border: 1.5px solid espresso800 (selected) OR cream400 (none)
flexDirection: row, alignItems: center, gap: 10
```
- Left: small bean shape icon (28×28 rounded square container, teardrop fill in roast color)
- Name + roaster + roast level (when selected). `"Select beans"` when empty.
- Right: chevron icon

Tapping opens a **Bean Sub-Sheet** (nested bottom sheet within the modal, see below).

#### Timer Ring

Centered, `paddingTop: 16, paddingBottom: 10`.

```
Ring outer diameter: 220px (SVG)
Track ring radius: 88px
Track stroke: cream400 (#E8D8B8), strokeWidth: 10
Ring rotated -90deg so progress starts at 12 o'clock
```

**Background zone arcs** (visible in idle state):
- Under zone arc: `0` to `(target-4)/target` fraction of circumference. Color: `caramel300` (`#E4A855`), opacity 0.30. Indicates under-extracted range.
- Perfect zone arc: `(target-4)/target` to `1.0` of circumference. Color: `matcha400` (`#6A9B74`), opacity 0.35. Rotated to start where under arc ends.

**Progress arc** (hidden in idle, animates during running/done):
- strokeLinecap: round
- Color transitions based on elapsed vs target:
  - `elapsed < target - 4`: caramel400 `#D4872A` (under)
  - `target - 4 ≤ elapsed ≤ target + 4`: matcha500 `#4A7C59` (perfect)
  - `elapsed > target + 4`: error `#C0392B` (over)
- Color transition: 0.6s ease

**Center text (idle):**
- `"ready"` — JetBrains Mono 500, 14px, letterSpacing: 1, cream500
- `"--:--"` — JetBrains Mono 600, 44px, letterSpacing: -2, cream400

**Center text (running/done):**
- `"MM:SS"` — JetBrains Mono 600, 52px, letterSpacing: -2, espresso800
- Zone label below: `"Under"` / `"On target"` / `"Over"` — DM Sans 700, 12px, color matches zone, uppercase, letterSpacing: 0.05em. Transition: color 0.6s ease.

**Pre-infusion phase:** When pre-infusion is active, show small `"Pre-infusion"` label above the time. DM Sans 700, 11px, caramel400, uppercase.

**Zone badge** (done state only, below ring):
```
background: zone color
borderRadius: full
paddingVertical: 6, paddingHorizontal: 20
```
Text: `"On target · 27s"` — DM Sans 700, 13px, white, letterSpacing: 0.04em.

#### Aim For + Pre-infusion Row

Visible in idle state only. `paddingHorizontal: 24, paddingBottom: 24`. `flexDirection: row, justifyContent: space-between, alignItems: center`.

**Aim for stepper:**
- Label: `"Aim for"` — DM Sans 500, 13px, espresso500
- Stepper: `background: cream300, borderRadius: full`. Minus button (34×32) | value in JetBrains Mono 600 14px espresso800 (minWidth: 32) | Plus button (34×32). Buttons have no background, font size 18, espresso500.
- Default: 27s. Range: 10–90s.

**Pre-infusion toggle:**
```
background: caramel100 (on) OR transparent (off)
border: 1.5px solid caramel400 (on) OR cream400 (off)
borderRadius: full, padding: 6 12
flexDirection: row, alignItems: center, gap: 6
```
- Indicator dot: `width:8, height:8, borderRadius:full`. caramel400 (on) or cream500 (off).
- Label: `"Pre-infusion"` — DM Sans 700, 12px. caramel500 (on) or espresso400 (off).
- No duration shown — duration is an internal constant for now.

**Divider:** `height:1, background: cream300, marginHorizontal: 22`

#### Stats Row

`paddingHorizontal: 24, paddingVertical: 24`. Three editable stats separated by 1px dividers.

Each stat:
- Label: DM Sans 700, 11px, uppercase, letterSpacing: 0.07em, espresso400
- Value: JetBrains Mono 600, 28px, letterSpacing: -1, espresso800. Dash (`—`) in cream500 when empty.
- Unit suffix: JetBrains Mono 400, 14px, espresso400

**Tap to edit:** Tapping a stat opens an inline text input (same font/size, no border, just a bottom border in caramel400 while focused). Commit on blur or Return key. Fields: Dose In (g), Yield Out (g), Grind Size (no unit).

**Live ratio:** If both dose and yield are filled, show `"1:X.X ratio"` centered below stats. JetBrains Mono 600, 13px, caramel500.

**Divider** (same as above)

#### Gear Row

Full-width pressable row:
```
paddingHorizontal: 24, paddingVertical: 18
flexDirection: row, justifyContent: space-between, alignItems: center
```
- Label: `"Gear used"` — DM Sans 700, 12px, uppercase, letterSpacing: 0.07em, espresso400, marginBottom: 4
- Value: gear names joined by ` · ` — DM Sans 700, 14px, espresso800. `"Add gear"` when empty.
- Chevron icon (right side). Hidden while running.

Tapping opens **Gear Sub-Sheet** (see below).

**Divider**

#### Tasting Note (done state only)

Appears after extraction is stopped:
```
paddingHorizontal: 22, paddingVertical: 18
```
- Label: `"Tasting note"` — DM Sans 700, 12px, uppercase, letterSpacing: 0.07em, espresso400
- Textarea: `minHeight: 72, padding: 12 14, background: cream200, border: 1.5px solid cream400, borderRadius: 14`. DM Sans 400, 14px, lineHeight: 1.5. Placeholder: `"How did it taste? Sour, sweet, balanced…"`

#### CTA Buttons

`paddingHorizontal: 22, paddingTop: 12, paddingBottom: 32`. `flexDirection: column, gap: 10`.

**Idle, live mode:**
```
height: 56, borderRadius: full
background: matcha500 (#4A7C59)
shadow: 0 6px 24px rgba(74,124,89,0.38)
```
Play icon + `"Start extraction"` — DM Sans 700, 16px, white.

**Running:**
```
background: error500 (#C0392B)
shadow: 0 6px 24px rgba(192,57,43,0.35)
```
Stop icon (filled square) + `"Stop"` — DM Sans 700, 16px, white.

**Done:**
- Reset button (left): `width:56, height:56, borderRadius:full, background:cream300`. Redo icon, espresso500.
- Save button (right, flex:1): `background: caramel400 (#D4872A), shadow: 0 6px 20px rgba(212,135,42,0.38)`. `"Log extraction"` — DM Sans 700, 16px, white.

**Manual entry toggle (idle only):**
- Small text link below the main CTA. `"Add time manually"` → `"← Use live timer instead"` to toggle. DM Sans 500, 13px, espresso400 (off) / caramel500 (on).
- When manual mode active: shows a number input (height:46, JetBrains Mono 600, 16px) + `"Log"` button (espresso800 bg). Input accepts seconds.

---

### Bean Sub-Sheet

Nested bottom sheet within the extraction modal. `maxHeight: 88%`.

Lists all beans from `BeansContext`. Each row:
```
background: espresso800 (selected) OR cream200 (unselected)
border: 1.5px solid espresso800 OR cream400
borderRadius: 16, padding: 13 15
gap: 12
```
- Bean icon: 36×36 rounded square, teardrop shape in roast color
- Name (DM Sans 700, 14px) + roaster + roast level (DM Sans 400, 12px, secondary)
- Check icon (caramel400) when selected

Footer: `"Add new beans"` dashed border button.

---

### Gear Sub-Sheet

**Station presets** at top: two side-by-side buttons loading the station's gear items into selection. Tapping a station adds (not replaces) its gear to the current selection.

**Individual gear list** below, toggle each item. Selected items shown with espresso800 background and caramel400 check.

**Save CTA:** `"Done · N items selected"` — espresso800 bg, full-width pill.

---

### Recent Extraction Card (Home Screen)

```
background: cream200 (#FAF3E8)
borderRadius: 20
padding: 18 18 16
shadow: 0 2px 10px rgba(28,15,7,0.06)
```

**Top row:**
- Left: bean name (DM Sans 700, 14px, espresso800) + roaster · time ago (DM Sans 400, 12px, espresso400)
- Right: zone badge
  - On target: `background: matcha100, color: matcha700, text: "✓ On target"`
  - Over: `background: error100, color: error500, text: "Over"`
  - Under: `background: caramel100, color: #7A4E0A, text: "Under"`
  - All: DM Sans 700, 11px, borderRadius: full, padding: 3 10

**Stats row:**
Four columns (Dose, Yield, Time, Ratio) separated by 1px cream400 dividers. Margin bottom: 12.
- Value: JetBrains Mono 600, 18px, letterSpacing: -0.5. Ratio uses caramel500.
- Unit: JetBrains Mono 400, 10px, espresso400 (superscript-style suffix)
- Label: DM Sans 700, 9px, uppercase, letterSpacing: 0.07em, espresso400, marginTop: 4

**Tasting note:** DM Sans 400, 12px, italic, espresso400, lineHeight: 1.45. Quoted with `"…"`.

---

## Timer State Machine

```
idle  ──tap Start──▶  preinfusion (if toggle on)  ──auto──▶  running  ──tap Stop──▶  done
 ▲                                                                                      │
 └───────────────────────────────tap Reset ─────────────────────────────────────────────┘
```

- **idle**: ring shows zone arcs, time shows `--:--`, fields editable
- **preinfusion**: timer runs for internal duration (no UI countdown shown), ring shows caramel progress, fields locked
- **running**: live timer, ring animates, color transitions by zone, fields locked
- **done**: timer frozen, save/discard CTA appears, tasting note field appears, fields editable again
- **Saving** calls the extractions API (not yet wired) and returns to idle, closing the modal

---

## Extraction Data Model

Based on the existing API patterns, an extraction record should include:

```ts
interface Extraction {
  id: string;
  bean_id: string | null;
  dose_in: number;            // grams
  yield_out: number | null;   // grams
  time: number;               // seconds (float)
  grind_size: number | null;
  gear_ids: string[];         // individual gear item IDs (not station IDs)
  pre_infusion: boolean;
  type: 'espresso' | 'ristretto' | 'lungo' | null;
  tasting_note: string | null;
  created_at: string;
  updated_at: string;
}
```

Note: **extractions are linked to individual gear items, not stations**. Stations are a convenience for pre-selection only.

---

## Animations & Transitions

| Element | Animation | Duration | Easing |
|---|---|---|---|
| Extraction modal open | `translateY: screenHeight → 0` | 360ms | `cubic-bezier(0.32,0,0.15,1.1)` (slight spring overshoot) |
| Modal backdrop | opacity 0→0.5 | 280ms | ease |
| Sub-sheets (bean/gear) | `translateY: 600 → 0` | 300ms | `cubic-bezier(0.33,1,0.68,1)` |
| Ring progress arc | `strokeDashoffset` update | — | rAF (60fps), no CSS transition while running |
| Ring color zone change | `stroke` color | 600ms | ease |
| Press on extraction card | `scale 1 → 0.98` | 120ms | ease |
| Press on nav items | `scale 1 → 0.88` | 120ms | ease |
| Toast | `translateY 16 → 0, opacity 0→1` | 280ms | `cubic-bezier(0.33,1,0.68,1)` |

---

## Design Tokens

All values come from the Mr. Bean Design System (`colors_and_type.css` / `src/theme/`):

### Colors used in this feature

| Token | Hex | Usage |
|---|---|---|
| `espresso800` | `#1C0F07` | Primary text, dark backgrounds, CTA |
| `espresso700` | `#2D1810` | Card sub-backgrounds, icon containers |
| `espresso600` | `#4A2518` | Decorative rings |
| `espresso500` | `#6B3A2A` | Secondary text |
| `espresso400` | `#8C5340` | Tertiary text, icons, labels |
| `caramel500`  | `#C4782A` | Accent (headline italic, ratio) |
| `caramel400`  | `#D4872A` | Under zone, check marks, border focus |
| `caramel300`  | `#E4A855` | Zone arc highlight |
| `caramel100`  | `#F8E4B8` | Pre-infusion toggle bg (on) |
| `cream600`    | `#C4A882` | Muted text on dark |
| `cream500`    | `#D8C4A0` | Handle, inactive nav |
| `cream400`    | `#E8D8B8` | Borders, dividers |
| `cream300`    | `#F0E4CC` | Subtle backgrounds, steppers |
| `cream200`    | `#FAF3E8` | Card backgrounds, inputs |
| `cream100`    | `#FDF8F2` | App background, modal bg |
| `matcha500`   | `#4A7C59` | Perfect zone, start CTA |
| `matcha400`   | `#6A9B74` | Perfect zone arc bg |
| `matcha700`   | `#2D5235` | On target badge text |
| `matcha100`   | `#E8F2EA` | On target badge bg |
| `error500`    | `#C0392B` | Over zone, stop CTA |
| `error100`    | `#FDECEA` | Over badge bg |

### Typography

| Usage | Family | Size | Weight | Notes |
|---|---|---|---|---|
| Display headlines | Playfair Display | 32px | 900 | Home header |
| Modal title | Playfair Display | 22px | 700 | Italic |
| Section headers | Playfair Display | 20px | 700 | Recent section |
| Card headline | Playfair Display | 28px | 700 | Home card |
| Body / labels | DM Sans | 13–15px | 400–700 | General |
| Timer display | JetBrains Mono | 52px | 600 | letterSpacing: -2 |
| Timer idle | JetBrains Mono | 44px | 600 | letterSpacing: -2 |
| Stats values | JetBrains Mono | 28px | 600 | letterSpacing: -1 |
| Recent stats | JetBrains Mono | 18px | 600 | letterSpacing: -0.5 |
| Stepper value | JetBrains Mono | 14px | 600 | |
| Ratio display | JetBrains Mono | 13px | 600 | |

### Spacing
Base unit: 4px. Key values used: 8, 10, 12, 14, 16, 18, 20, 22, 24, 28, 32.

### Radii
- `full`: 9999px (pills, toggles, CTAs)
- `28–32`: modal sheet top corners
- `20`: recent extraction cards
- `16`: bean/gear rows in sheets
- `14`: inputs, textarea
- `12`: nav items, small containers

### Shadows
- Home card: `0 8px 40px rgba(28,15,7,0.20)`
- Start CTA (green): `0 6px 24px rgba(74,124,89,0.38)`
- Stop CTA (red): `0 6px 24px rgba(192,57,43,0.35)`
- Save CTA (caramel): `0 6px 20px rgba(212,135,42,0.38)`
- Modal: `0 -8px 48px rgba(28,15,7,0.18)`
- Recent cards: `0 2px 10px rgba(28,15,7,0.06)`

---

## Assets

- `assets/logo-icon.svg` — bean mark used in header (34px)
- `assets/logo.svg` — full wordmark (not used in this screen)
- Icons: **Lucide Icons** (`@expo/vector-icons` or `lucide-react-native`) — stroke-based, 2px weight, rounded. Key icons: `play`, `square` (stop), `refresh-cw` (reset), `chevron-right`, `x`, `check`, `plus`.

---

## Files in this Package

| File | Purpose |
|---|---|
| `README.md` | This document |
| `Home Screen.html` | Current working design — open in any browser to interact |

---

## Out of Scope (not in this handoff)

- All Extractions history screen
- Prediction engine / AI cards
- Beans tab
- Stats tab
- Extraction API endpoint (to be defined separately)
