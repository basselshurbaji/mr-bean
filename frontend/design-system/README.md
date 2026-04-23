# Mr. Bean Design System

## Overview

**Mr. Bean** is a mobile app for coffee enthusiasts — from the "just got an espresso machine" beginner to the obsessive who weighs their beans to the tenth of a gram. The app lets users log their coffee extractions (espresso shots, pour-overs, etc.), troubleshoot and dial in recipes, and get AI-powered predictions for new beans based on their history.

Think of it as a "brew journal + shot whisperer" — equal parts notebook, calculator, and trusted coffee nerd friend.

**Source materials:** No Figma or codebase was provided. This design system was built from scratch based on the brand brief.

---

## Products

| Product | Description |
|---|---|
| **Mr. Bean iOS/Android App** | Core mobile product. Extraction logging, shot dialing, bean library, prediction engine. |

---

## CONTENT FUNDAMENTALS

### Voice & Tone

Mr. Bean speaks like your most knowledgeable coffee friend — the one who genuinely cares about your shot quality but also knows when to make a joke about it. Never condescending, always helpful, occasionally hilarious.

**Core traits:**
- **Warm** — celebrates your wins, commiserates your sours
- **Humorous** — coffee puns are not just allowed, they're encouraged
- **Direct** — gets to the point; no padding
- **Encouraging** — every bad shot is data, not failure

### Pronouns & POV
- Refers to the user as **"you"** — personal and direct
- App speaks in **first person** only for personality lines ("I've seen worse extractions. Actually, no I haven't.")
- UI labels and actions are **second person**: "Your Beans", "Add Shot", "Your History"

### Casing
- **Sentence case** for all UI labels, buttons, and headings
- **Title Case** only for product names (Mr. Bean, Extraction Log)
- No ALL CAPS except for acronyms (TDS, EY)
- Numbers formatted with units: `18.5 g`, `28 s`, `1:2.1`

### Emoji Usage
- Used **sparingly** and only in onboarding, empty states, and celebratory moments
- Never in navigation or functional UI
- Preferred coffee set: ☕ 🫘 ⚗️ 🌡️

### Humor Guidelines
- Puns on "bean", "shot", "grind", "extraction", "roast" — all fair game
- Self-aware about the absurdity of coffee obsession ("Yes, you do need a refractometer.")
- Never punches down at beginners — humor is always inclusive
- Toast messages and empty states are the main playground for humor

### Copy Examples
| Context | Copy |
|---|---|
| Empty bean library | "No beans on deck. Time to fix that, don't you think?" |
| First shot logged | "Your first log! The journey from 'meh' to 'magnificent' starts here." |
| Prediction ready | "Based on your history, here's my best guess. I've been wrong before, but not often." |
| Bad extraction | "Oof. That one's data, not failure. Adjust and pull again." |
| Loading | "Grinding the numbers…" |
| Error | "Something went sideways. Even my grinder jams sometimes." |

---

## VISUAL FOUNDATIONS

### Color Philosophy
Warm, earthy, and inviting — like the inside of a specialty coffee shop at 8am. Rooted in the actual colors of coffee: from the near-black of a ristretto to the pale gold of cream in a latte. Accented with a deep matcha green that nods to both nature and the specialty coffee world.

### Typography
- **Display**: Playfair Display — serif, characterful, warm. Used for big headlines, feature numbers.
- **Body**: DM Sans — clean, modern, highly readable at small sizes.
- **Mono**: JetBrains Mono — for all numbers, measurements, ratios. Gives precision and readability to data.
- Google Fonts substitutions in use (see `colors_and_type.css` for import URLs).

### Backgrounds
- Default: `--color-foam` (`#FDF8F2`) — very warm off-white
- Cards: `--color-cream` (`#FAF3E8`) with subtle shadow
- Dark/rich sections: `--color-espresso` (`#1C0F07`)
- No gradient backgrounds in main UI; gradients only for hero/onboarding moments

### Spacing & Layout
- Base unit: `4px`. Scale: 4, 8, 12, 16, 24, 32, 48, 64
- Mobile-first; max content width ~390px
- Fixed bottom navigation bar (48px safe area + 64px nav)
- 16px side padding standard; 24px for cards

### Corner Radii
- `--radius-sm`: 6px (tags, badges, chips)
- `--radius-md`: 12px (inputs, small cards)
- `--radius-lg`: 20px (bottom sheets, main cards)
- `--radius-xl`: 32px (pill buttons, FAB)

### Shadows
- Cards: `0 2px 8px rgba(28,15,7,0.08)` — warm, subtle
- Elevated sheets: `0 -4px 24px rgba(28,15,7,0.12)`
- No hard drop shadows; warmth through softness

### Animation
- Easing: `cubic-bezier(0.32, 0, 0.67, 0)` for exits; `cubic-bezier(0.33, 1, 0.68, 1)` for entrances
- Duration: 200ms micro, 350ms standard, 500ms complex
- Bottom sheets: spring-like slide up with subtle overshoot
- No gratuitous animation; motion serves meaning

### Hover & Press States
- Hover: slight background tint (`--color-cream` overlay at 40%)
- Press: scale `0.97` + darken background — tactile feel
- Destructive: red tint on hover/press
- No underlines on links; color change only

### Iconography
See ICONOGRAPHY section below.

### Cards
- Rounded corners `--radius-lg`
- Background `--color-cream`
- Shadow: `0 2px 8px rgba(28,15,7,0.08)`
- 16px padding internal
- Border: none (shadow defines the edge)

### Imagery
- Warm, slightly golden color grading
- Shallow depth of field (coffee-shop aesthetic)
- No stock photography; prefer illustrations or blank placeholders
- When using photography: always warm tones, never cold/blue

### Transparency & Blur
- Blur used on bottom nav overlay: `backdrop-filter: blur(20px)`
- Bottom nav: semi-transparent `--color-foam` at 85% opacity
- Header on scroll: same treatment

---

## ICONOGRAPHY

No custom icon font or sprite sheet provided. The app uses **Lucide Icons** (CDN), matched to the stroke weight and rounded style of the brand.

- Icon library: [Lucide Icons](https://lucide.dev) via CDN
- Style: Stroke-based, 2px weight, rounded caps/joins
- Size: 20px standard in nav; 24px in headers; 16px inline
- Color: inherits from text context
- No filled icons except active nav state (filled variant)

Key icons used:
| Purpose | Lucide name |
|---|---|
| Add extraction | `plus-circle` |
| Beans library | `package` |
| Analytics | `bar-chart-2` |
| Settings | `settings` |
| Timer | `timer` |
| Temperature | `thermometer` |
| Weight/dose | `scale` |
| Notes | `file-text` |
| Star/favorite | `star` |
| Prediction | `sparkles` |

---

## File Index

```
README.md                          — This file
SKILL.md                           — Agent skill definition
colors_and_type.css               — CSS variables: colors, typography, spacing, radii
assets/
  logo.svg                        — Mr. Bean wordmark + bean icon
  logo-icon.svg                   — Icon-only mark
  logo-dark.svg                   — Light version for dark backgrounds
preview/
  colors-base.html                — Base color palette
  colors-semantic.html            — Semantic / functional colors
  colors-dark.html                — Dark mode palette
  type-scale.html                 — Type scale specimen
  type-specimens.html             — Display, body, mono specimens
  spacing-tokens.html             — Spacing + radius tokens
  shadows-elevation.html          — Shadow & elevation system
  components-buttons.html         — Button variants
  components-inputs.html          — Form inputs
  components-cards.html           — Card variants
  components-badges.html          — Tags, badges, chips
  brand-logo.html                 — Logo usage
  brand-voice.html                — Voice & tone quick-ref
ui_kits/
  app/
    README.md                     — App UI kit documentation
    index.html                    — Interactive app prototype
    components/
      Navigation.jsx              — Bottom nav bar
      Header.jsx                  — App header variants
      ExtractionCard.jsx          — Shot log card
      BeanCard.jsx                — Bean library card
      ShotForm.jsx                — Log a shot form
      PredictionCard.jsx          — Prediction result card
```
