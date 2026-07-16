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

**Installation**: Users can install your extension by simply dropping the `.zip` file into Aether's `Browse Extensions` page, or by placing the extracted folder directly into the Aether extensions directory:
- **Windows**: `%APPDATA%/Aether/extensions/`
- **macOS**: `~/Library/Application Support/Aether/extensions/`
- **Linux**: `~/.config/Aether/extensions/`

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

## Developer Experience (Planned: Aether CLI)

The Aether CLI (`aether`) is the planned official toolkit for creating, developing, testing, packaging, validating, and publishing Aether extensions. The goal is a new developer can go from nothing to a working extension in under five minutes.

### Project Creation

```
aether init
```

Starts an interactive folder generator. You will be asked:

- **Extension Name**
- **Extension ID** (e.g. `com.example.my-extension`)
- **Author**
- **Version**
- **Description**
- **License**
- **Extension Type**: Feature or Appearance Pack
- **Framework**: Vanilla, React, Vue, Svelte, or Solid
- **Homepage** *(optional)*
- **Repository URL** *(optional)*

Scaffolds the following structure:

```
my-extension/
├── manifest.json
├── package.json
├── README.md
├── LICENSE
├── src/
│   └── main.js
├── ui/
│   ├── index.html
│   ├── main.js
│   └── styles.css
├── assets/
└── .gitignore
```

> **Note**: Choosing a framework (React, Vue, Svelte, Solid) scaffolds a Vite build config inside `ui/` and requires a build step before the extension runs. Vanilla skips the build step entirely.

---

### Development

```
aether dev
```

Watches source files and hot-reloads the extension. Requires an active Aether instance to be running — if one is not detected, the command will print an error and exit rather than silently doing nothing. Shows extension logs, API calls, permission usage, and runtime errors with source maps.

---

### Validation

```
aether validate
```

Checks manifest syntax, missing files, invalid permissions, API compatibility, version format, duplicate IDs, invalid assets, missing icons, and missing metadata. Returns a clear pass/fail with specific error messages.

---

### Packaging

```
aether build
```

Produces a `my-extension.aex` file. Automatically minifies, compresses, validates, generates a checksum, and strips development files. The `.aex` format is Aether's planned first-class extension container. *(Currently, standard `.zip` files are used for packaging)*.

---

### Testing

```
aether test
```

Runs optional API mocks, permission tests, UI snapshot tests, manifest validation, and integration tests.

---

### Version Management

```
aether version patch
aether version minor
aether version major
```

Automatically updates `manifest.json`, `package.json`, and `CHANGELOG`.

---

### Utilities

| Command | Purpose |
|---|---|
| `aether lint` | Warns about unused permissions, deprecated APIs, missing metadata |
| `aether fmt` | Formats `manifest.json`, source, and configuration |
| `aether clean` | Removes build files |
| `aether info` | Shows extension ID, version, API version, permissions, build size, author |
| `aether migrate` | Upgrades manifest and API usage for new API versions |
| `aether permissions` | Scans source code and suggests required permissions |
| `aether docs <api>` | Searches and prints API documentation inline |
| `aether examples` | Generates example code for sidebar pages, dialogs, loaders, etc. |

---

### Registry Commands

```
aether search <query>
aether install <extension-id>
aether remove <extension-id>
aether update <extension-id>
```

---

### Developer Tools

```
aether console    # Open extension console
aether inspect    # Inspect a running extension
aether trace      # View live API calls
aether profile    # View performance metrics
```

---

### Future Commands

| Command | Purpose |
|---|---|
| `aether benchmark` | Measures startup time and memory usage |
| `aether doctor` | Checks the development environment |
| `aether sdk update` | Updates SDK templates to the latest version |
| `aether create provider` | Scaffolds a new Loader Provider extension |
| `aether create theme` | Scaffolds a new Appearance Pack |
| `aether create loader` | Scaffolds a new Mod Loader extension |
