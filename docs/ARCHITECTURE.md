# Architecture

## Backend
The launcher backend is written in Go to ensure high performance, memory safety, and native compilation across Windows, macOS, and Linux. The backend operates completely headlessly and serves the frontend via Wails bindings. The frontend is a lightweight web view rendering the UI.

## Go Packages
- `main.go`: The main entry point.
- `pkg/instance`: Logic for managing Minecraft instances, resolving dependencies, and constructing launch arguments.
- `pkg/java`: Discovery, installation, and management of Java Runtimes (JRE/JDK).
- `pkg/auth`: Microsoft/Xbox authentication flow and token management.
- `pkg/extensions`: The extension manager and sandbox environment.

## Extension Manager
The Extension Manager (`pkg/extensions`) is responsible for discovering, validating, and executing extensions. 
Extensions are executed in an isolated JavaScript runtime (Goja). They do not have direct access to the host OS. Any action an extension wants to perform must be requested through the launcher's API, which enforces the capability-based permission model.

## Updater
The updater is built into the core launcher but operates independently.
- **Launcher Updates**: The launcher downloads new binaries in the background and swaps them on the next restart.
- **Extension Updates**: Extensions are checked for updates periodically. Updates are applied automatically unless configured otherwise, minimizing user friction.

## Launcher Pipeline
1. **Resolution**: Determine the Minecraft version, loader (Fabric, Forge, etc.), and required libraries.
2. **Verification**: Check if all assets, libraries, and the Java runtime are present. Download missing files.
3. **Authentication**: Ensure the user has a valid session token. Refresh if necessary.
4. **Execution**: Construct the massive Java command line and spawn the child process.
5. **Monitoring**: Pipe standard output/error to the log manager and monitor the process state.

## Diagrams

### Launcher Startup & Extension Loading

```mermaid
sequenceDiagram
    participant OS as Operating System
    participant Core as Aether Core (Go)
    participant UI as Aether UI (Svelte)
    participant ExtManager as Extension Manager
    participant Sandbox as Goja Sandbox

    OS->>Core: Launch Aether
    Core->>UI: Start Wails Webview
    Core->>ExtManager: Initialize()
    
    loop For each Extension
        ExtManager->>ExtManager: Read manifest.json
        ExtManager->>Sandbox: Create new isolated runtime
        ExtManager->>Sandbox: Inject permitted Aether APIs
        ExtManager->>Sandbox: Execute main.js
        Sandbox->>ExtManager: Aether.ui.registerSidebarPage()
        ExtManager->>UI: Emit 'extension:sidebar:add' event
    end
    
    UI-->>User: Render fully loaded UI
```

### Permission Validation

```mermaid
sequenceDiagram
    participant Sandbox as Goja Sandbox
    participant API as Aether Core API
    
    Sandbox->>API: Aether.instances.patch(id, data)
    API->>API: Check Extension Manifest
    
    alt Has 'instances:patch' permission
        API->>API: Apply Patch
        API-->>Sandbox: Success (true)
    else Missing permission
        API-->>Sandbox: Throw Error ("Permission Denied")
    end
```
