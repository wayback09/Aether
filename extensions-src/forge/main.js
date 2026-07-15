Aether.launcher.registerModLoader({
    id: "forge",
    name: "Forge",
    description: "The original Minecraft mod loader",
    onLaunch: function(ctx) {
        var mcVersion = ctx.mcVersion;

        // 1. Fetch promotions to get the recommended Forge version
        var promosStr = Aether.http.get("https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json");
        var promos = JSON.parse(promosStr);
        var forgeVer = promos.promos[mcVersion + "-recommended"] || promos.promos[mcVersion + "-latest"];
        if (!forgeVer) {
            throw new Error("Forge is not available for Minecraft " + mcVersion);
        }

        var fullVer = mcVersion + "-" + forgeVer;
        var mavenUrl = "https://maven.minecraftforge.net/";
        var basePath = "net/minecraftforge/forge/" + fullVer + "/forge-" + fullVer;

        // 2. Determine which jar to download (client for modern, universal for legacy)
        var jarType = "client";
        var mainClass = "cpw.mods.modlauncher.Launcher";

        // 3. Try to fetch the Forge version JSON for library listing
        var jsonUrl = mavenUrl + basePath + ".json";
        try {
            var jsonStr = Aether.http.get(jsonUrl);
            var versionInfo = JSON.parse(jsonStr);

            // Download all libraries from the version JSON
            var libs = versionInfo.libraries || [];
            for (var i = 0; i < libs.length; i++) {
                var lib = libs[i];
                if (lib.name && lib.name.indexOf("net.minecraftforge:forge:") === -1) {
                    if (lib.downloads && lib.downloads.artifact && lib.downloads.artifact.url) {
                        try {
                            var localPath = Aether.fs.download(lib.downloads.artifact.url, lib.downloads.artifact.path);
                            ctx.classpath.push(localPath);
                        } catch(e) {}
                    }
                }
            }

            // Use mainClass from version JSON
            var mc = versionInfo.mainClass;
            if (typeof mc === "string") {
                mainClass = mc;
            } else if (mc && mc.client) {
                mainClass = mc.client;
            }
        } catch(e) {
            // If no version JSON, we still proceed with the client jar
        }

        // 4. Download the Forge client/universal jar
        var jarUrl = mavenUrl + basePath + "-" + jarType + ".jar";
        try {
            var jarPath = Aether.fs.download(jarUrl, basePath + "-" + jarType + ".jar");
            ctx.classpath.push(jarPath);
        } catch(e) {
            // Fallback to universal jar for older Forge versions
            jarUrl = mavenUrl + basePath + "-universal.jar";
            var jarPath = Aether.fs.download(jarUrl, basePath + "-universal.jar");
            ctx.classpath.push(jarPath);
        }

        ctx.mainClass = mainClass;
        return ctx;
    }
});
