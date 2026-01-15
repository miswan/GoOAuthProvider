## 2024-05-21 - [Form Binding for Accessibility]
**Learning:** Go structs often lack `form` tags, defaulting to JSON. This prevents standard HTML `<form>` submissions (application/x-www-form-urlencoded), forcing users/devs to use API tools.
**Action:** Always add `form:"name"` alongside `json:"name"` in request models to support both API and basic HTML interfaces.
