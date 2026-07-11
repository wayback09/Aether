# Aether Future Architecture Plan

This document outlines architectural improvements to be implemented in a future refactoring phase. Our current sandbox architecture works well, but as we add more extensions, we will need a stricter Service-Based model.

## 1. Extension Capabilities (Provides/Consumes)
Transition the `manifest.json` from a generic `permissions` array to a strict `provides` and `consumes` schema.

**Example:**
```json
{
  "provides": [
    "service:loader:fabric",
    "ui:sidebar:settings"
  ],
  "consumes": [
    "core:instances:patch",
    "core:downloads"
  ]
}
```
**Goal**: Make it instantly obvious what an extension does and what rights it needs, allowing the launcher to gracefully disable extensions if dependencies are missing.

## 2. Universal Service Discovery API
Instead of hardcoding APIs in the Goja runtime, Aether should expose a central Service Registry. 

**Implementation Idea:**
- Expose `Aether.services.list()`
- Extensions can query: `Aether.services.get("loader:fabric")`
- Aether acts as the middleman, routing requests between extensions while maintaining complete sandbox isolation. Extensions never share memory or talk directly; they only pass JSON messages through Aether's core bus.

## 3. Standardized Provider Interfaces
When extensions register as a provider (e.g., a Mod Loader Installer), they should conform to a strict interface.

**Implementation Idea:**
- Instead of raw callbacks, Aether enforces that a Loader Provider must return an object with specific methods: `download()`, `getClasspath()`, `getMainClass()`.
- If an extension fails to provide these, Aether rejects the registration.

---
*Note: These changes are non-urgent. They should be implemented when the extension ecosystem grows large enough to require strict interplay between community extensions.*

## 4. Appearance Packs (UI Extensions)
A UI Extension customizes how Aether looks and feels. Unlike feature extensions, UI extensions cannot add launcher functionality. Their only purpose is presentation.

**Examples:** Themes, Icon Packs, Layout Variants, Typography Packs, Density Presets, Animation Packs.
- They never install mods.
- They never manage instances.
- They never access launcher data unless explicitly granted.

### Architecture
```
Core Launcher
       ↓
UI Manager
       ↓
Installed UI Extension
       ↓
Theme Tokens / Component Overrides / Icons / Animations
```

### Rules
**UI Extensions MAY:** Change colors, icons, fonts, spacing, border radius, component density, animation timing, register component styles, register custom illustrations, provide custom empty states.

**UI Extensions MAY NOT:** Replace the sidebar, replace navigation, hide security warnings, inject arbitrary HTML into launcher chrome, override dialogs, move permission prompts, modify another extension, execute privileged API calls.

### Component System & Overrides
The launcher owns every component (Button, Card, Input, Search, Sidebar, Dialog, Toast, ProgressBar, Checkbox, Select, Badge, Tabs). UI extensions only change how these components look. They never replace them.

**Theme Tokens Example:**
```json
{
  "colors": {
    "background": "#0B0B0B",
    "surface": "#141414",
    "accent": "#3B82F6",
    "danger": "#DC2626"
  },
  "radius": 8,
  "spacing": "comfortable",
  "font": "Inter",
  "density": "normal"
}
```
The launcher maps these tokens to every component. Extensions may optionally override the appearance of individual components (e.g., Button, Card, Sidebar Item), but the launcher still controls behaviour.

### Icon Packs & Assets
A UI extension may register an icon pack (e.g., Default, Lucide, Minecraft Style, Windows Fluent, Material Symbols). The launcher automatically swaps icons.
- **Allowed Assets:** SVG, PNG, Fonts, CSS Variables
- **Not Allowed:** Executable code, Dynamic scripts, External CDN assets

### Settings & Layouts
Every UI Extension automatically receives its own settings page (Appearance → Theme, Icons, Typography, Density, Advanced). Supported layout presets include Compact, Comfortable, Touch, Wide. The launcher adapts spacing automatically; extensions cannot invent their own layouts.

### Animations
UI extensions may change: Duration, Easing, Fade, Slide Distance, Scale.
They may not: Create particle effects, Add unnecessary motion, Reduce accessibility.

### Compatibility & Composability
UI extensions declare compatibility:
```json
{
    "type": "ui",
    "supports": {
        "api": "1.0",
        "components": "1.0",
        "theme": "1.0"
    }
}
```
UI extensions should be composable. Instead of installing one massive theme, users can mix and match UI packs however they like (e.g., Catppuccin theme + Fluent icons + Minimal animations).

### Design Philosophy
Feature Extensions change **WHAT** Aether can do. UI Extensions change **HOW** Aether feels. Neither should interfere with the other. The launcher always owns the layout, navigation, security, and behaviour.

## 5. CI/CD & Automation (GitHub Actions)
To maintain stability and security across the ecosystem, both the core launcher and the extension registry will rely on GitHub Actions for automated workflows.

### Aether (Core Launcher)
- **Multi-OS Builds (CD)**: Automated pipelines will run `wails build` across Windows, macOS, and Linux runners whenever a new version tag (e.g., `v1.2.0`) is pushed.
- **Automated Releases**: Compiled binaries (`.exe`, `.app`, `.AppImage`) will be automatically packaged and attached to GitHub Releases for immediate user download.
- **Testing (CI)**: Go unit tests and Svelte component tests will run on every Pull Request to prevent regressions in the core API and UI.

### Aether-Extensions (Registry)
- **Automated Security Audits (CI)**: When a developer submits a PR to add their extension to `index.json`, an Action will download their packaged `.zip`, extract it, and validate the `manifest.json`. It will verify that no hidden or dangerous permissions (e.g., `fs:write` to system directories) are requested.
- **Community Tier Assignment**: If the automated security and structure checks pass, the workflow will approve the PR and automatically assign the `community` trust badge.
- **Registry Linter**: A JSON validation step will guarantee `index.json` is perfectly formatted, preventing a missing comma from breaking the in-app gallery for all users.
