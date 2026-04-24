# Feature Design: My Gear

## Overview

My Gear lets authenticated users register their espresso hardware and organise it into **Stations**. A station is a named, ordered collection of gear items that pre-selects tools when logging an extraction shot — it is a convenience helper only and is never directly linked to a shot record.

---

## Entities

### GearItem

Represents a single piece of equipment owned by a user.

| Field        | Type      | Constraints                          |
|--------------|-----------|--------------------------------------|
| `id`         | UUID      | PK, generated on create              |
| `user_id`    | UUID      | FK → users.id, NOT NULL              |
| `type_id`    | string    | NOT NULL, must be a valid type (see Equipment Types) |
| `name`       | string    | NOT NULL                             |
| `brand`      | string    | nullable                             |
| `model`      | string    | nullable                             |
| `year`       | string(4) | nullable, 4-digit year               |
| `notes`      | string    | nullable                             |
| `created_at` | timestamp | set on insert                        |
| `updated_at` | timestamp | set on insert, updated on every write |

### Station

A named, ordered set of gear item references belonging to a user.

| Field        | Type      | Constraints             |
|--------------|-----------|-------------------------|
| `id`         | UUID      | PK, generated on create |
| `user_id`    | UUID      | FK → users.id, NOT NULL |
| `name`       | string    | NOT NULL                |
| `created_at` | timestamp | set on insert           |
| `updated_at` | timestamp | updated on every write  |

### StationGear (join)

Tracks which gear items belong to a station and their display order.

| Field        | Type    | Constraints                              |
|--------------|---------|------------------------------------------|
| `station_id` | UUID    | FK → stations.id, NOT NULL               |
| `gear_id`    | UUID    | FK → gear.id, NOT NULL                   |
| `position`   | integer | NOT NULL, 0-indexed, unique per station  |

PK is `(station_id, gear_id)`.

---

## Equipment Types

Types are a **closed, hardcoded list** validated server-side. They are never stored as a separate table.

| type_id       | Display label      |
|---------------|--------------------|
| `machine`     | Espresso machine   |
| `grinder`     | Grinder            |
| `scale`       | Scale              |
| `portafilter` | Portafilter        |
| `tamper`      | Tamper             |
| `distributor` | Distribution tool  |
| `wdt`         | WDT tool           |
| `basket`      | Basket             |
| `puckscreen`  | Puck screen        |
| `other`       | Other              |

---

## Actions & Endpoints

All endpoints are protected (require authenticated user). Every query is scoped to the caller's `user_id`.

### Gear

| Method   | Path          | Description                                      |
|----------|---------------|--------------------------------------------------|
| `GET`    | `/gear`       | List all gear items for the authenticated user   |
| `POST`   | `/gear`       | Create a new gear item                           |
| `GET`    | `/gear/:id`   | Get a single gear item by ID                     |
| `PUT`    | `/gear/:id`   | Update a gear item                               |
| `DELETE` | `/gear/:id`   | Delete a gear item and remove it from all stations |

#### GET /gear

Response: array of gear items, ordered by `created_at ASC`.

#### POST /gear

Request body:

```json
{
  "type_id": "grinder",
  "name":    "Niche Zero",
  "brand":   "Niche",
  "model":   "Zero",
  "year":    "2022",
  "notes":   "Single dose, 63mm conical burrs."
}
```

- `type_id` and `name` are required; all other fields are optional.
- `type_id` must be one of the 10 valid values — reject with 422 otherwise.
- Returns the created gear item.

#### PUT /gear/:id

Same body shape as POST. Partial updates are not supported — send all fields (omit optional fields to clear them). Returns the updated item.

- Caller must own the item — reject with 404 if not found or belongs to another user.

#### DELETE /gear/:id

- Caller must own the item.
- Removes the gear item from all `station_gear` rows referencing it (cascade or explicit delete).
- Re-sequence remaining positions within each affected station so positions remain contiguous.
- Returns 204 on success.

---

### Stations

| Method   | Path              | Description                                 |
|----------|-------------------|---------------------------------------------|
| `GET`    | `/stations`       | List all stations for the authenticated user |
| `POST`   | `/stations`       | Create a new station                        |
| `PUT`    | `/stations/:id`   | Update a station (name + gear list)         |
| `DELETE` | `/stations/:id`   | Delete a station                            |

#### GET /stations

Response: array of stations. Each station includes its ordered gear item list (full gear objects, not just IDs), so the client can render the icon strip without a second request.

```json
[
  {
    "id":   "...",
    "name": "Morning routine",
    "gear": [
      { "id": "...", "type_id": "machine", "name": "Lelit Bianca", ... },
      { "id": "...", "type_id": "grinder", "name": "Niche Zero",   ... }
    ]
  }
]
```

#### POST /stations

Request body:

```json
{
  "name":     "Morning routine",
  "gear_ids": ["uuid-1", "uuid-2"]
}
```

- `name` is required; `gear_ids` may be empty (`[]`).
- `gear_ids` must reference gear items owned by the caller — reject with 422 if any ID is unknown or unowned.
- The order of `gear_ids` in the request defines the `position` sequence.
- Returns the created station with its full gear list.

#### PUT /stations/:id

Same body shape as POST. Replaces the station's name and full gear list atomically (delete existing `station_gear` rows, insert new ones from the provided `gear_ids` order).

- Caller must own the station.
- Returns the updated station with its full gear list.

#### DELETE /stations/:id

- Caller must own the station.
- Deletes the station and all its `station_gear` rows.
- Does **not** delete any gear items.
- Returns 204 on success.

---

## Business Rules

1. **User scoping.** Every read and write is filtered by `user_id` derived from the auth token. A user can never read or modify another user's gear or stations.
2. **Type validation.** `type_id` is validated against the closed list on every create and update. Invalid values return 422.
3. **Gear deletion cascade.** Deleting a gear item must remove it from every station that includes it. Remaining gear positions in those stations must be re-compacted (no gaps).
4. **Station gear ordering.** The order of `gear_ids` sent by the client is the canonical order. The server stores and returns items in that order.
5. **Station is a helper, not a record.** Stations are never referenced by extraction records (shots). Deleting a station has no impact on shot history.
6. **No shared gear.** Gear items are not shared across users. Each user manages their own independent inventory.

---

## Summary Stats

The UI header displays `{n} pieces · {n} stations`. These counts are derivable client-side from the list responses — no dedicated summary endpoint is needed.

---

## Error Responses

| Scenario                              | Status |
|---------------------------------------|--------|
| Missing required field                | 422    |
| Invalid `type_id`                     | 422    |
| `gear_ids` contains unknown/unowned ID | 422   |
| Item not found or belongs to another user | 404 |
| Unauthenticated request               | 401    |
