// Register the sidebar UI
Aether.ui.registerSidebarPage({
    id: "mod-manager-ui",
    label: "Mod Manager",
    url: "ui/index.html"
});

// Listen for IPC messages from the frontend UI
Aether.ui.onMessage((payload) => {
    try {
        if (payload.action === 'get_instances') {
            const instances = Aether.instances.list();
            return { action: 'get_instances', success: true, data: instances };
        }

        if (payload.action === 'get_mods') {
            const mods = Aether.instances.listMods(payload.instanceId);
            return { action: 'get_mods', success: true, data: mods };
        }

        if (payload.action === 'toggle_mod') {
            Aether.instances.toggleMod(payload.instanceId, payload.jarName, payload.enable);
            return { action: 'toggle_mod', success: true };
        }

        if (payload.action === 'delete_mod') {
            Aether.instances.deleteMod(payload.instanceId, payload.jarName);
            return { action: 'delete_mod', success: true };
        }
    } catch (e) {
        return { action: payload.action, success: false, error: e.toString() };
    }

    return { success: false, error: "Unknown action" };
});
