# Palette's Journal

## 2024-05-22 - Embedded HTML Status Pages
**Learning:** For backend-heavy services, providing a lightweight, embedded HTML status page at the root URL significantly improves developer experience and immediate system feedback compared to a 404.
**Action:** Utilize inline SVGs and scoped CSS within `c.HTML()` handlers to create self-contained, accessible (ARIA `role="status"`) landing pages without external dependencies.
