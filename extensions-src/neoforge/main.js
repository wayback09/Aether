Aether.launcher.registerModLoader({
    id: "neoforge",
    name: "NeoForge",
    description: "Modern, community-maintained fork of Forge",
    onLaunch: function(ctx) {
        var mcVersion = ctx.mcVersion;

        // 1. Fetch NeoForge API to get the recommended version
        var apiUrl = "https://api.neoforged.net/api/v1/versions/minecraft";
        var apiStr = Aether.http.get(apiUrl);
        var mcVersions = JSON.parse(apiStr);

        // Find the entry matching our MC version
        var neoforgeVer = null;
        for (var i = 0; i < mcVersions.length; i++) {
            if (mcVersions[i].version === mcVersion) {
                neoforgeVer = mcVersions[i].recommended || mcVersions[i].latest;
                break;
            }
        }
        if (!neoforgeVer) {
            throw new Error("NeoForge is not available for Minecraft " + mcVersion);
        }

        var mavenUrl = "https://maven.neoforged.net/releases/";
        var basePath = "net/neoforged/neoforge/" + neoforgeVer + "/neoforge-" + neoforgeVer;
        var mainClass = "cpw.mods.modlauncher.Launcher";

        // 2. Try to fetch the NeoForge version JSON for library listing
        var jsonUrl = mavenUrl + basePath + ".json";
        try {
            var jsonStr = Aether.http.get(jsonUrl);
            var versionInfo = JSON.parse(jsonStr);

            var libs = versionInfo.libraries || [];
            for (var i = 0; i < libs.length; i++) {
                var lib = libs[i];
                if (lib.name && lib.name.indexOf("net.neoforged:neoforge:") === -1) {
                    if (lib.downloads && lib.downloads.artifact && lib.downloads.artifact.url) {
                        try {
                            var localPath = Aether.fs.download(lib.downloads.artifact.url, lib.downloads.artifact.path);
                            ctx.classpath.push(localPath);
                        } catch(e) {}
                    }
                }
            }

            var mc = versionInfo.mainClass;
            if (typeof mc === "string") {
                mainClass = mc;
            } else if (mc && mc.client) {
                mainClass = mc.client;
            }
        } catch(e) {}

        // 3. Download the NeoForge client jar
        var jarUrl = mavenUrl + basePath + "-client.jar";
        try {
            var jarPath = Aether.fs.download(jarUrl, basePath + "-client.jar");
            ctx.classpath.push(jarPath);
        } catch(e) {
            // Fallback to the main jar (for versions that bundle client)
            jarUrl = mavenUrl + basePath + ".jar";
            var jarPath = Aether.fs.download(jarUrl, basePath + ".jar");
            ctx.classpath.push(jarPath);
        }

        ctx.mainClass = mainClass;
        return ctx;
    }
});
