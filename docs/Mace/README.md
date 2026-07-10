# Mace

Mace is a lightweight, cross-platform desktop application for managing Minecraft server instances. It is built using the Wails framework with a Go backend and a React/TypeScript frontend.

<p align="center">
  <img src="appicon.png" width="128" height="128" alt="Mace Icon" />
</p>

## Screenshots

<details>
  <summary>📸 Click to view application screenshots</summary>
  <br/>

  | Dashboard Overview | Resource Usage & Players |
  | :---: | :---: |
  | ![Dashboard Overview](Screenshots/Screenshot%202026-01-01%20064218.png) | ![Resource Usage](Screenshots/Screenshot%202026-01-01%20065045.png) |

  | Mods & Plugins Manager | Intelligent Crash Analyzer |
  | :---: | :---: |
  | ![Mods & Plugins Manager](Screenshots/Screenshot%202026-01-01%20064642.png) | ![Intelligent Crash Analyzer](Screenshots/Screenshot%202026-01-01%2006444g4.png) |

  | Create New Instance | Server Import |
  | :---: | :---: |
  | ![Create New Instance](Screenshots/Screenshot%202026-01-01%20064408.png) | ![Server Import](Screenshots/Screenshot%202026-01-01%20064417.png) |
</details>

## Features

- **Isolated Server Instances**: Create and manage multiple isolated server instances in separate directories.
- **Import Existing Servers**: Hook into an existing Minecraft server directory to manage it without moving your files. Auto-detects loader type.
- **Multiple Loaders Supported**: Automatic download and setup of **Vanilla**, **Spigot**, **Paper**, **Fabric**, **Quilt**, **Forge**, and **NeoForge** servers.
- **Automatic Java Detection**: Scans the system for installed Java versions to ensure compatibility.
- **Real-Time Console & Player List**: Interactive terminal console with support for sending commands directly to the server input stream, combined with live active player list tracking.
- **Resource Monitoring**: Live tracking of CPU usage, RAM utilization, and server uptime.
- **Config Editor**: Edit `server.properties` visually and directly from the interface.
- **Crash Watchdog & Intelligent Analyzer**:
  - **Crash Watchdog**: Automatically monitors running server processes and auto-restarts them upon crash detection.
  - **Intelligent Crash Analyzer**: Inspects console logs and crash reports on unexpected stops to identify common issues (Java version mismatch, port conflicts, out-of-memory errors, missing mod dependencies, or unaccepted EULAs) and suggest resolutions.
- **Content Management (Mods & Plugins)**:
  - Native integrations with **Modrinth**, **CurseForge**, **Hangar** (PaperMC), and **Spiget** (SpigotMC) APIs to search, browse, and install mods or plugins.
  - Toggle mods and plugins on or off with a simple switch (manages renaming to/from `.jar.disabled`).
  - Validate JAR signatures during import or download to ensure compatibility with your chosen loader (checks for `plugin.yml`, `fabric.mod.json`, `mods.toml`, etc.).
- **Modpack Support**:
  - Support for applying Modrinth `.mrpack` files (with automated dependency resolution and asset downloading) and generic `.zip` overrides (CurseForge modpacks).
- **Server Backups**:
  - Compress world directories and configurations (including `server.properties`, `whitelist.json`, `ops.json`, etc.) into timestamped ZIP archives.
  - Safe restoration: automatically stops the server if running, creates a temporary pre-restore backup, and extracts the zip safely with zip-slip prevention.
  - Customize backup directories per server instance.

## Prerequisites

To run or build the application from source, you need:

- **Go**: 1.22.0 or higher
- **Node.js & npm**: For frontend assets compilation
- **Wails CLI**: Install via `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Java**: Java Development Kit (JDK) installed and configured on your path (relevant version depending on the Minecraft version you plan to run).

## Project Structure

- `backend/pkg/`: Core Go logic
  - `downloader/`: Fetches version manifests, installer URLs, handles async jar downloads, and integrates with remote APIs (Modrinth, CurseForge, Hangar, Spiget).
  - `launcher/`: Manages Java execution processes, stdin/stdout piping, process logs, player list monitoring, and the crash watchdog/analyzer.
  - `servermanager/`: High-level CRUD operations for managing server configs, backups, and mod/plugin directories.
- `frontend/src/`: Frontend React application
  - `components/`: UI components (Console, ConfigEditor, ResourceMonitor, ServerCard, ContentManager, etc.).
  - `pages/`: Application screens (Instances list, Create Server, Settings).

## Getting Started

### Development

To run the application in live-development mode with hot-reloading:

```bash
wails dev
```

### Build

To package the application into a production executable:

```bash
wails build
```

The compiled binary will be located in the `build/bin/` folder.

## Links

- **Website**: [MACE website](https://mace-7yf.pages.dev)
- **Discord**: [Join the community](https://discord.com/invite/zrrHQC4QKF)

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
