// Modrinth Browser Extension Backend
// Runs in the secure Goja Sandbox

Aether.ui.registerSidebarPage({
    id: "modrinth",
    label: "Modrinth",
    url: "ui/index.html"
});

// Handle IPC messages from the UI iframe
Aether.ui.onMessage(function(msg) {
    if (msg.type === "get_instances") {
        var instances = Aether.instances.list();
        Aether.ui.postMessage({
            type: "instances_result",
            requestId: msg.requestId,
            instances: instances
        });
    }

    if (msg.type === "install_mod") {
        try {
            Aether.instances.installMod(msg.instanceId, msg.jarName, msg.downloadUrl);
            Aether.ui.postMessage({
                type: "install_result",
                requestId: msg.requestId,
                success: true,
                jarName: msg.jarName
            });
        } catch (e) {
            Aether.ui.postMessage({
                type: "install_result",
                requestId: msg.requestId,
                success: false,
                error: String(e)
            });
        }
    }
});
