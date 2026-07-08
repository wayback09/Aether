# Security

## Sandbox
All extensions are executed inside a heavily restricted JavaScript isolate. They do not have access to the Node.js standard library, the DOM of the main launcher window, or the underlying operating system. The only way an extension can affect the outside world is via the exposed `Aether` JSON-RPC API, which validates every request.

## Whitelisted Networking
By default, extensions cannot make outbound network requests. 
If an extension needs to communicate with an external API (e.g., Modrinth or CurseForge), it must declare the specific domains in its `manifest.json` under `hostPermissions`.
The user is informed of these domains during installation.

```json
"hostPermissions": [
  "https://api.modrinth.com/*"
]
```

## Capability Model
The security architecture uses a capability-based model. Instead of giving an extension raw file system access to read logs, the launcher provides a specific `Aether.instances.getLogs(id)` function. This ensures the extension can only read logs and nothing else. If an extension attempts to call an API it lacks permissions for, the call is immediately rejected, and a security warning is logged.

## Review Process
1. **Automated Analysis**: Static analysis tools scan the submitted extension bundle for obfuscated code, forbidden API usage, and known vulnerabilities.
2. **Manual Review**: A human reviewer inspects the manifest, permissions, and core logic.
3. **Community Reporting**: Users can flag suspicious extensions, which triggers an immediate quarantine and manual re-evaluation.

## Threat Model
**Expected Threats:**
- Malicious extensions attempting to steal Minecraft session tokens.
- Extensions attempting to download and execute arbitrary binaries (malware).
- Extensions attempting to read arbitrary files on the user's system (e.g., SSH keys, browser cookies).
- Extensions attempting to execute arbitrary shell commands.
- Extensions attempting to escape the Goja Sandbox.

**Mitigation:**
Every privileged operation is strictly mediated by the launcher. Extensions explicitly **CANNOT**:
- Execute arbitrary shell commands (e.g., via `os/exec`).
- Access the filesystem directly (they can only use scoped `Aether.fs` APIs).
- Read launcher memory.
- Access Microsoft authentication tokens.
- Access another extension's data.
- Escape the Goja sandbox (which has zero bindings by default).
- Access Go APIs directly.

- Token access is strictly forbidden for extensions. Authentication is handled entirely by the core Go backend.
- File system access is abstracted. Extensions cannot read or write outside their designated isolated storage directory unless explicitly granted highly scrutinized permissions.
- Network access is strictly allow-listed.
