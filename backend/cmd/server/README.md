# cmd/server

Binary entrypoint. This file contains wiring only — no business logic.

Responsibilities:
- Load config
- Open DB connection
- Instantiate repos, services, and token service
- Register routes (public and protected)
- Start the HTTP server

Any adapter code needed to bridge incompatible interfaces between packages also lives here (e.g. `userStoreAdapter`).

Nothing in this directory should be imported by other packages.
