# Feature Design: Home Screen & Extraction Logging

## Overview

The Home Screen is the daily entry point of the Mr. Bean app. It greets the user, surfaces their most recent extractions, and provides a single-tap path into the **Extraction Modal** — a full-height sheet where the user times or manually logs a shot, fills in brew parameters, and saves the record. Every logged extraction becomes part of the user's history and feeds future AI predictions. This feature defines the `Extraction` entity and all endpoints needed to create, read, update, and delete extraction records.

---

## Entities

### Extraction

| Field          | Type      | Constraints                                                              |
|----------------|-----------|--------------------------------------------------------------------------|
| `id`           | UUID      | PK, generated on create                                                  |
| `user_id`      | UUID      | FK → users.id, NOT NULL                                                  |
| `bean_id`      | UUID      | FK → beans.id, NOT NULL                                                  |
| `dose_in`      | float     | NOT NULL, > 0 (grams)                                                    |
| `yield_out`    | float     | NOT NULL, > 0 (grams)                                                    |
| `time`         | float     | NOT NULL, > 0 (seconds; extraction duration only, excludes pre-infusion) |
| `target_time`  | float     | NOT NULL, > 0 (seconds; the user's aim-for value at time of logging)     |
| `grind_size`   | float     | NOT NULL, > 0 (unitless numeric scale)                                   |
| `pre_infusion` | boolean   | NOT NULL, default false                                                  |
| `tasting_note` | string    | nullable                                                                 |
| `created_at`   | timestamp | set on insert                                                            |
| `updated_at`   | timestamp | set on insert, updated on every write                                    |

### ExtractionGear (join)

Links an extraction to the individual gear items used during that shot.

| Field           | Type | Constraints                   |
|-----------------|------|-------------------------------|
| `extraction_id` | UUID | FK → extractions.id, NOT NULL |
| `gear_id`       | UUID | FK → gear.id, NOT NULL        |

PK is `(extraction_id, gear_id)`.

---

## Zone Classification

Zone is a **derived value** computed from `time` and `target_time`. It is never stored; callers compute it from the returned fields.

| Condition                                  | Zone      |
|--------------------------------------------|-----------|
| `time < target_time - 4`                   | `under`   |
| `target_time - 4 ≤ time ≤ target_time + 4` | `perfect` |
| `time > target_time + 4`                   | `over`    |

---

## Actions & Endpoints

All endpoints are protected (require authenticated user). Every query is scoped to the caller's `user_id`.

| Method   | Path               | Description                                 |
|----------|--------------------|---------------------------------------------|
| `GET`    | `/extractions`     | List extractions for the authenticated user |
| `POST`   | `/extractions`     | Log a new extraction                        |
| `GET`    | `/extractions/:id` | Get a single extraction by ID               |
| `PUT`    | `/extractions/:id` | Update an extraction                        |
| `DELETE` | `/extractions/:id` | Delete an extraction                        |

---

### GET /extractions

Returns extractions ordered by `created_at DESC`.

Query parameters:

| Param    | Type    | Default | Description                                    |
|----------|---------|---------|------------------------------------------------|
| `limit`  | integer | 20      | Number of records per page                     |
| `page`   | integer | 1       | Page number (1-indexed). First page is page 1. |

Response: array of extraction objects. Each object includes the full bean object and the full list of gear objects, so the client can render cards without additional requests.

```json
[
  {
    "id": "...",
    "user_id": "...",
    "bean": {
      "id": "...",
      "name": "Yirgacheffe Natural",
      "roaster": "Square Mile",
      "roast": "light"
    },
    "dose_in": 18.0,
    "yield_out": 36.5,
    "time": 27.4,
    "target_time": 27.0,
    "grind_size": 14.0,
    "pre_infusion": false,
    "tasting_note": "Juicy and bright.",
    "gear": [
      { "id": "...", "type_id": "machine", "name": "Rocket Appartamento" },
      { "id": "...", "type_id": "grinder", "name": "Niche Zero" }
    ],
    "created_at": "2026-04-28T07:12:00Z",
    "updated_at": "2026-04-28T07:12:00Z"
  }
]
```

---

### POST /extractions

Logs a new extraction.

Request body:

```json
{
  "bean_id":      "uuid",
  "dose_in":      18.0,
  "yield_out":    36.5,
  "time":         27.4,
  "target_time":  27.0,
  "grind_size":   14.0,
  "gear_ids":     ["uuid-1", "uuid-2"],
  "pre_infusion": false,
  "tasting_note": "Juicy and bright."
}
```

All fields are required except `tasting_note` (which may be omitted or sent as `null`). `gear_ids` must be present but may be an empty array.

Validation:
- `dose_in`, `yield_out`, `time`, `target_time`, and `grind_size` must be > 0.
- `bean_id` must reference a bean visible to the caller — reject with 422 otherwise.
- Each ID in `gear_ids` must reference a gear item owned by the caller — reject with 422 if any ID is unknown or unowned.

Returns the created extraction in the same shape as the GET /extractions list item (with embedded bean and gear objects).

---

### GET /extractions/:id

Returns a single extraction in the same shape as the list item.

- Returns 404 if the extraction does not exist or belongs to another user.

---

### PUT /extractions/:id

Updates an extraction. Same body shape as POST. Full replacement — send all fields; omit `tasting_note` or send `null` to clear it.

Validation rules are identical to POST.

- Caller must own the extraction — return 404 otherwise.
- Returns the updated extraction.

---

### DELETE /extractions/:id

- Caller must own the extraction.
- Deletes the extraction record and all associated `extraction_gear` rows.
- Returns 204 on success.

---

## Business Rules

1. **User scoping.** Every read and write is filtered by `user_id` derived from the auth token. A user can never read or modify another user's extractions.
2. **Gear is linked to individual items, not stations.** The `gear_ids` array holds gear item IDs. Stations are a client-side convenience for pre-selection; they are never referenced on an extraction record.
3. **`target_time` must be stored.** The user's aim-for value at the time of logging is stored on the record so that zone classification (`under` / `perfect` / `over`) is reproducible for historical extractions without relying on a current preference.
4. **Zone is derived, not stored.** The backend returns `time` and `target_time`; the client computes zone using the classification table above.
5. **`time` is extraction duration only.** When pre-infusion is used, the pre-infusion phase duration is excluded from `time`. `pre_infusion: true` signals that a pre-infusion phase occurred, but its duration is not stored (it is a fixed internal constant in the client).
6. **Bean visibility.** `bean_id` must reference a bean record accessible to the caller. Beans owned by other users are not valid references.
7. **No hard limit on extractions per user.** There is no cap on the number of extraction records per user.

---

## Error Responses

| Scenario                                                              | Status |
|-----------------------------------------------------------------------|--------|
| Missing required field                                                | 422    |
| `dose_in`, `yield_out`, `time`, `target_time`, or `grind_size` is ≤ 0 | 422    |
| `bean_id` is unknown or inaccessible                                  | 422    |
| `gear_ids` contains an unknown or unowned ID                          | 422    |
| Extraction not found or belongs to another user                       | 404    |
| Unauthenticated request                                               | 401    |
