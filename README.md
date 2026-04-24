# Mr. Bean

**Mr. Bean** is a mobile app for coffee enthusiasts — from the "just got an espresso machine" beginner to the obsessive who weighs their beans to the tenth of a gram. Log extractions, dial in recipes, manage your gear, and get AI-powered predictions for new beans based on your history.

Think of it as a "brew journal + shot whisperer" — equal parts notebook, calculator, and trusted coffee nerd friend.

---

## What it does

| Feature            | Description                                                                                                                               |
|--------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| **Extraction log** | Record espresso shots and pour-overs: dose, yield, time, TDS, notes, and rating                                                           |
| **Shot dialing**   | Troubleshoot and iterate on recipes; track changes across pulls                                                                           |
| **Bean library**   | Catalog your beans with origin, roast, and tasting notes                                                                                  |
| **My Gear**        | Register hardware (machines, grinders, scales, baskets, etc.) and group them into named Stations that pre-select gear when logging a shot |
| **AI predictions** | Get dial-in suggestions for new beans based on your extraction history                                                                    |
| **Profile**        | Edit name and password; avatar auto-generated from initials                                                                               |

---

## Screens designed so far

| Screen               | Status          | Notes                                                        |
|----------------------|-----------------|--------------------------------------------------------------|
| Login / Register     | High-fidelity   | Two-mode (sign in / create account) on one screen; JWT auth  |
| My Gear              | High-fidelity   | Gear list + detail + add/edit sheet; Stations tab            |
| Profile              | High-fidelity   | Edit name, change password, log out                          |
| Home                 | Planned         | Extraction feed                                              |
| Beans                | Planned         | Bean library                                                 |

---

## Stack

| Layer     | Technology                                        |
|-----------|---------------------------------------------------|
| Mobile    | React Native / Expo (Expo Router, file-based nav) |
| Backend   | Go (chi router, JWT auth, PostgreSQL via sqlc)    |
| AI        | Planned                                           |

---

## Design system

Warm, earthy palette rooted in actual coffee colours — near-black espresso (`#1C0F07`) through caramel (`#C4782A`) to pale cream (`#FDF8F2`). Typography: **Playfair Display** for headings, **DM Sans** for body/UI, **JetBrains Mono** for measurements and data. Design files live in `frontend/design/`.

---

## Repo layout

```
mr_bean/
├── backend/    # Go API server — see backend/README.md
└── frontend/   # Expo app + design handoffs
    └── design/
        ├── design-system/          # Colour, type, spacing tokens and previews
        ├── design_handoff_login/   # Login & register screen
        ├── design_handoff_profile/ # Profile screen
        └── design_handoff_my_gear/ # My Gear + Stations screens
```
