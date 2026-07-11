# Extensions Guide

## Architecture Overview
Aether extensions operate in two distinct, isolated layers:
1. **The Backend Sandbox (`main.js`)**: Runs in Aether's secure, headless Goja engine. It has no DOM access, but can interact with Aether's native Go APIs (e.g., to patch instances or read files, based on requested permissions).
2. **The Frontend UI (`ui/index.html`)**: Runs in the Aether Svelte app as a secure `<iframe>`. Aether spins up a lightweight local HTTP server to serve these files.

## How to Build an Extension

1. Create a new directory for your extension.
2. Create a `manifest.json`.
3. Write your `main.js` backend entry point.
4. Create a `ui` folder containing your `index.html` and any CSS/JS.
5. Use the `Aether.ui` API in your backend script to register your interface.

```javascript
// main.js
Aether.ui.registerSidebarPage({
    id: "my-custom-page",
    label: "My Page",
    url: "ui/index.html" // Points to your HTML file relative to your extension folder
});
```

## Packaging and Installation

We recommend using a rich package structure to give your extension a professional presentation:

```text
my-extension.zip
├── manifest.json
├── main.js
├── README.md
├── LICENSE
├── CHANGELOG.md
├── icon.png
├── ui/
│   ├── index.html
│   ├── style.css
│   └── script.js
└── assets/
```

**Installation**: Users can install your extension by simply dropping the `.zip` file into Aether's `Browse Extensions` page, or by placing the extracted folder directly into `%APPDATA%/Aether/extensions/`.

## Manifest
Every extension requires a `manifest.json` at its root.
```json
{
  "id": "com.example.myextension",
  "name": "My Extension",
  "version": "1.0.0",
  "author": "Your Name",
  "description": "Adds a cool new feature to Aether.",
  "main": "main.js",
  "api": "1.0",
  "minApi": "1.0",
  "maxApi": "2.0",
  "homepage": "https://example.com/myextension",
  "repository": "https://github.com/example/myextension",
  "license": "MIT",
  "keywords": ["mods", "fabric"],
  "permissions": [
    "ui:sidebar",
    "instances:patch"
  ]
}
```

## Permissions
Extensions operate under a principle of least privilege. You must explicitly request access to APIs.
- `ui:sidebar`: Register sidebar pages that render your `ui/index.html` in an iframe.
- `ui:dialogs`: Open modal dialogs.
- `instances:patch`: Modify instance JSON files (e.g., to install loaders like Fabric).

## Extension UI Rules
Because extension UIs run inside an `<iframe>`, you have complete control over your DOM. You can use React, Vue, Svelte, Solid, Lit, or plain HTML/CSS. Aether is completely **framework-agnostic**.

However, to maintain a consistent user experience, we recommend matching Aether's dark, frosted-glass aesthetic.

## Examples
Check the `extensions-src/` directory in the Aether repository for complete sample extensions, including the Modrinth Browser and the Fabric mod loader.

## Trust Tiers & Review Guidelines
When an extension is published to the Aether Registry, it is assigned a trust tier. The launcher displays these badges so users know what they are installing:

1. 🔵 **Official**: Developed and maintained directly by the Aether Team.
2. 🟢 **Verified**: Personally reviewed by an Aether maintainer. The code has been thoroughly audited for security, performance, and stability.
3. 🟣 **Community**: Passed automated checks and was merged into the registry via Pull Request, but has not received a manual code audit. Use with caution.
4. 🟡 **Local**: Installed manually from a `.zip` file. Aether treats these as untrusted.

To get an extension **Verified**, it must pass a manual review:
1. **No Malicious Code**: Extensions must not steal tokens, install malware, or attempt to break out of the Goja sandbox.
2. **Performance**: Extensions must not leak memory or block the main thread.
3. **Clear Purpose**: The extension must do exactly what its description claims.

## Developer Experience (Planned SDK)
To simplify extension development in the future, we plan to release a CLI tool (e.g., `aether-cli`) offering commands like:
- `npm create aether-extension`
- `aether init`
- `aether dev` (for hot-reloading)
- `aether pack` (to bundle into `.zip`)
- `aether validate` (to check manifest and security constraints)
