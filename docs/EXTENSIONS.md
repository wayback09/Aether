# Extensions Guide

## How to Build an Extension
Extensions are written in JavaScript/TypeScript. They interact with the launcher via the injected `Aether` global object.

1. Create a new directory for your extension.
2. Create a `manifest.json`.
3. Write your `index.js` entry point.
4. Use the `Aether.ui` API to register your interface.

```javascript
// index.js
Aether.ui.registerSidebarPage({
    id: "my-custom-page",
    title: "My Page",
    icon: "box",
    render: (container) => {
        container.innerHTML = "<h1>Hello from my extension!</h1>";
    }
});
```

## Manifest
Every extension requires a `manifest.json` at its root.
```json
{
  "id": "com.example.myextension",
  "name": "My Extension",
  "version": "1.0.0",
  "author": "Your Name",
  "description": "Adds a cool new feature to Aether.",
  "main": "index.js",
  "api": "^1.0.0",
  "permissions": [
    "ui:sidebar"
  ]
}
```

## Permissions
Extensions operate under a principle of least privilege. You must explicitly request access to APIs.
- `ui:sidebar`: Register sidebar pages.
- `ui:dialogs`: Open modal dialogs.
- `instances:read`: Read instance data.
- `network:fetch`: Make network requests (requires a `hostPermissions` array).

## Extension UI Rules
Extensions cannot draw arbitrary windows.
Extensions may register:
- Sidebar Pages
- Toolbar Buttons
- Context Menus
- Settings Pages
- Dialogs
- Notifications

Nothing else.
The launcher owns the layout.
Extensions only provide content.

## Examples
Check the `examples/` directory in the Aether repository for complete sample extensions, including a basic mod downloader and a custom theme engine.

## Review Guidelines
Before an extension can be published to the official Aether registry, it must pass a review:
1. **No Malicious Code**: Extensions must not steal tokens, install malware, or bypass the sandbox.
2. **Minimal Performance Impact**: Extensions must not leak memory or block the main thread.
3. **UI Consistency**: Extensions must use the provided UI components and adhere to the `STYLEGUIDE.md`.
4. **Clear Purpose**: The extension must do what its description claims.
