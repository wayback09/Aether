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
