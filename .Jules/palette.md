## 2024-05-23 - [Missing Landing Page]
**Learning:** This backend service had no root (`/`) handler, returning a 404. Even for API-only services, a root handler confirming status is a huge UX win for developers.
**Action:** Always verify the root endpoint of a service.
