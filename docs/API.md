# API & Interoperability

## JSON-RPC Endpoints
The launcher exposes a local JSON-RPC 2.0 API over a named pipe or Unix domain socket for extensions to communicate with the core launcher.

### Core Endpoints
- `launcher.getVersion`: Returns the current launcher version.
- `launcher.getInstances`: Returns a list of all installed Minecraft instances.
- `launcher.launchInstance(id)`: Initiates the launch sequence for a specific instance.
- `launcher.killInstance(id)`: Forcefully terminates a running instance.

### UI Endpoints
- `ui.registerSidebarPage(id, title, icon, url)`: Registers a new page in the sidebar.
- `ui.showNotification(title, message, type)`: Displays a notification (info, warning, error).
- `ui.openDialog(options)`: Opens a modal dialog and returns the user's action.

## Permissions
Extensions run in an isolated environment and must request permissions in their `manifest.json`.

- `instances:read`: View installed instances and their metadata.
- `instances:write`: Create, modify, or delete instances.
- `process:launch`: Start the Minecraft process.
- `network:http`: Make external HTTP requests (requires whitelisted domains).
- `fs:read`: Read from specific launcher directories.
- `fs:write`: Write to extension-specific storage.

## Versioning
The API follows Semantic Versioning (SemVer). 
- **Major version bumps** indicate breaking changes.
- **Minor version bumps** indicate new features.
- Extensions must declare their compatible API version range in `manifest.json`.

## Extension Lifecycle
1. **Installed**: The extension is downloaded and extracted.
2. **Initialized**: The launcher reads the `manifest.json` and registers the extension.
3. **Activated**: The `main` script of the extension is executed in the sandbox. This happens either on startup or when the extension's UI is accessed (lazy loading).
4. **Suspended**: If the extension is inactive and consuming resources, it may be suspended.
5. **Terminated**: The extension is stopped, and its resources are freed (e.g., when uninstalled or the launcher closes).
