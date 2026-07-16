# Aether Launcher

<p align="center">
  <img src="frontend/public/logo.png" alt="Aether Logo" width="120" />
</p>

<p align="center">
  <em>A minimal, extensible, lightning-fast Minecraft launcher.</em>
</p>

## Overview

Aether is designed with a single core principle: **This launcher exists to launch Minecraft, and nothing more.**

Every feature that is not essential to launching the game—like downloading mods, managing servers, or viewing logs—is implemented as an extension. By keeping the core launcher intentionally minimal, Aether remains fast, predictable, and free from bloat. Think of it as the VS Code of Minecraft launchers. You install only the features you actually want.

## Features

- **Blazing Fast**: Cold starts in under 2 seconds with an idle memory footprint below 100MB.
- **Minimalist UI**: Clean, native-feeling, elegant design that gets out of your way, complete with smooth custom toast notifications.
- **Snapshot Support**: Effortlessly toggle between stable releases and latest snapshots when creating instances.
- **Extensible Architecture**: Everything from Modrinth integration to server browsers is an extension.
- **Secure Sandbox**: Extensions run in a strict capability-based JavaScript isolate to protect your session tokens and your operating system.

## Documentation

All project documentation is located in the `docs/` directory. If you are looking to contribute, build an extension, or just understand how Aether works, start here:

- **[Project Philosophy & Design](docs/DESIGN.md)** - The core principles guiding Aether's development.
- **[Architecture](docs/ARCHITECTURE.md)** - Overview of the Go backend, Extension Manager, and launcher pipeline.
- **[API & Interoperability](docs/API.md)** - Details on the JSON-RPC endpoints and permission model.
- **[Extensions Guide](docs/EXTENSIONS.md)** - How to build, package, and publish extensions for Aether.
  - *Looking for the official extension registry? Visit the [Aether-Extensions](https://github.com/wayback09/Aether-Extensions) repository.*
- **[Security & Sandboxing](docs/SECURITY.md)** - Threat models, capability isolation, and review guidelines.
- **[UI Specifications](docs/UI.md)** - Layout rules, components, and empty states.
- **[Styleguide](docs/STYLEGUIDE.md)** - Visual language, typography, colors, and animations.

## Getting Started

*(Coming soon - Build instructions and release binaries will be available here.)*