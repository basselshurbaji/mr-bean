# Product — Feature Requirements

This directory is the source of truth for **what we are building and why**. It is the contract between product and backend.

---

## What belongs here

- Feature scope and intent (what problem does this solve, what does it enable)
- Data entities: fields, types, constraints, relationships
- Core actions: what a user can do, expressed as API endpoints
- Business rules: validation logic, ownership, cascades, edge cases
- Error handling: which scenarios return which HTTP status codes

## What does NOT belong here

- Code, queries, or implementation details → those live in `backend/` or `mcp/`

---

## Who reads this

| Audience | Uses it for                                                   |
| -------- | ------------------------------------------------------------- |
| Backend  | Defining models, migrations, API routes, and validation logic |

---

## Directory structure

Each feature gets its own subdirectory:

```
product/
  <feature-name>/
    feature-design.md    Core requirements document
```

Reference: `my-gear/feature-design.md` is the canonical example of a well-formed document.

---

## feature-design.md format

Every feature document should cover these sections, in order:

### 1. Overview
One paragraph. What is this feature? What problem does it solve? What is the user able to do that they could not do before?

### 2. Entities
For each data model: a table of fields with type and constraints. Include relationships (FKs, nullable, uniqueness). If values are constrained to a closed list, define that list here.

### 3. Actions & Endpoints
A table of HTTP method + path + description, then a section per endpoint with:
- Request body (JSON shape, required vs optional fields)
- Validation rules
- Response shape
- Ownership/auth rules (can user X act on resource owned by user Y?)

### 4. Business Rules
Numbered list. Each rule is a constraint that goes beyond field-level validation — cascade behaviour, ordering, limits, invariants. These are the rules a backend engineer needs to encode in service logic.

### 5. Error Responses
A table mapping scenarios to HTTP status codes.

---

## Tone and precision

- Be specific enough that a backend engineer can write the migration and handlers without asking follow-up questions.
- Do not over-specify implementation. Say *what* must happen, not *how* to implement it.
- Use plain language. Avoid vague words like "handle", "manage", or "support" — say exactly what should happen.
