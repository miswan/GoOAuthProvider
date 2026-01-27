## 2026-01-27 - Frontend Verification without Database
**Learning:** When the main application has hard database dependencies that prevent startup in the test environment, create a temporary, standalone verification server (e.g., `verification/server.go`) that registers only the handler being tested.
**Action:** Use this pattern for future frontend verifications of isolated components to bypass complex environment setup.
