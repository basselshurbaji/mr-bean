# Feature Design: Beans

## Overview

Beans lets authenticated users build a personal catalogue of coffee beans they use or have used. A bean represents a coffee product — a specific roast from a specific roaster — not a physical bag or purchase. The same bean entry can be referenced across multiple brew sessions without being tied to a single purchase date or bag size.

---

## Entities

### Bean

Represents a coffee product in the user's catalogue.

| Field           | Type      | Constraints                                           |
|-----------------|-----------|-------------------------------------------------------|
| `id`            | UUID      | PK, generated on create                               |
| `user_id`       | UUID      | FK → users.id, NOT NULL                               |
| `name`          | string    | NOT NULL (e.g. "Ethiopia Yirgacheffe")                |
| `roaster`       | string    | nullable                                              |
| `origin`        | string    | nullable — country or region (e.g. "Colombia, Huila") |
| `process`       | string    | nullable — must be a valid process (see Process Types) |
| `roast_level`   | string    | nullable — must be a valid level (see Roast Levels)   |
| `tasting_notes` | string    | nullable — free-text flavor descriptors               |
| `notes`         | string    | nullable — personal notes                             |
| `created_at`    | timestamp | set on insert                                         |
| `updated_at`    | timestamp | set on insert, updated on every write                 |

---

## Process Types

A **closed, hardcoded list** validated server-side. Never stored as a separate table.

| process_id  | Display label |
|-------------|---------------|
| `washed`    | Washed        |
| `natural`   | Natural       |
| `honey`     | Honey         |
| `anaerobic` | Anaerobic     |
| `other`     | Other         |

---

## Roast Levels

A **closed, hardcoded list** validated server-side. Never stored as a separate table.

| roast_level_id | Display label |
|----------------|---------------|
| `light`        | Light         |
| `medium_light` | Medium Light  |
| `medium`       | Medium        |
| `medium_dark`  | Medium Dark   |
| `dark`         | Dark          |

---

## Actions & Endpoints

All endpoints are protected (require authenticated user). Every query is scoped to the caller's `user_id`.

| Method   | Path          | Description                                    |
|----------|---------------|------------------------------------------------|
| `GET`    | `/beans`      | List all beans for the authenticated user      |
| `POST`   | `/beans`      | Add a new bean                                 |
| `PUT`    | `/beans/:id`  | Update a bean                                  |
| `DELETE` | `/beans/:id`  | Delete a bean                                  |

### GET /beans

Response: array of beans, ordered by `created_at ASC`.

### POST /beans

Request body:

```json
{
  "name":          "Ethiopia Yirgacheffe",
  "roaster":       "Onyx Coffee Lab",
  "origin":        "Ethiopia, Yirgacheffe",
  "process":       "washed",
  "roast_level":   "light",
  "tasting_notes": "Jasmine, bergamot, peach.",
  "notes":         "Best pulled at 1:2.5 ratio."
}
```

- `name` is required; all other fields are optional.
- `process` must be one of the 5 valid values — reject with 422 otherwise.
- `roast_level` must be one of the 5 valid values — reject with 422 otherwise.
- Returns the created bean.

### PUT /beans/:id

Same body shape as POST. Partial updates are not supported — send all fields (omit optional fields to clear them). Returns the updated bean.

- Caller must own the bean — reject with 404 if not found or belongs to another user.

### DELETE /beans/:id

- Caller must own the bean — reject with 404 if not found or belongs to another user.
- Returns 204 on success.

---

## Business Rules

1. **User scoping.** Every read and write is filtered by `user_id` derived from the auth token. A user can never read or modify another user's beans.
2. **Process validation.** `process` is validated against the closed list on every create and update. Invalid values return 422.
3. **Roast level validation.** `roast_level` is validated against the closed list on every create and update. Invalid values return 422.
4. **Bean is a product, not a purchase.** A bean entry represents a coffee product (roaster + name + profile) and is intentionally decoupled from purchase events, bag weights, or roast dates. Those concerns belong to a future purchase or inventory feature.
5. **No shared beans.** Beans are not shared across users. Each user manages their own independent catalogue.

---

## Error Responses

| Scenario                               | Status |
|----------------------------------------|--------|
| Missing required field (`name`)        | 422    |
| Invalid `process`                      | 422    |
| Invalid `roast_level`                  | 422    |
| Bean not found or belongs to another user | 404 |
| Unauthenticated request                | 401    |
