## 2024-05-23 - [Adding Landing Page to API]
**Learning:** Even purely backend/API services benefit significantly from a human-readable root endpoint (`/`). It provides immediate "system operational" verification and acts as entry-level documentation, reducing the "black box" feeling for new developers or integrators.
**Action:** For API projects, always ensure `GET /` returns a friendly HTML status page instead of 404 or JSON error.
