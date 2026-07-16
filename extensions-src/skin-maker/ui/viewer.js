var viewerContainer = document.getElementById('viewer-3d');
var autoRotate = true;

initCoreScene(viewerContainer);
buildDefaultSkin();
buildModel(scene, false);
startRenderLoop(function () {
    if (autoRotate && playerGroup) {
        playerGroup.rotation.y += 0.008;
    }
});

function loadSkinToViewer(imageUrl, modelType) {
    var img = new Image();
    img.crossOrigin = 'Anonymous';
    img.onload = function () {
        skinCtx.drawImage(img, 0, 0, SKIN_SIZE, SKIN_SIZE);
        updateTexture();
        buildModel(scene, modelType === 'slim');
    };
    img.onerror = function () {
        document.getElementById('viewer-status').textContent = 'Failed to load skin';
    };
    img.src = imageUrl;
}

document.getElementById('search-btn').addEventListener('click', function () {
    var name = document.getElementById('username-input').value.trim();
    if (!name) return;
    document.getElementById('viewer-status').textContent = 'Searching...';
    window.parent.postMessage({ type: "fetch_profile", username: name, reqId: 'viewer' }, '*');
});

document.getElementById('username-input').addEventListener('keydown', function (e) {
    if (e.key === 'Enter') document.getElementById('search-btn').click();
});

window.addEventListener('message', function (e) {
    if (e.data && e.data.type === 'profile_result' && e.data.reqId === 'viewer') {
        if (e.data.error) {
            document.getElementById('viewer-status').textContent = 'Error: ' + e.data.error;
            return;
        }
        if (e.data.skinUrl) {
            document.getElementById('viewer-status').textContent = e.data.username + ' — ' + e.data.uuid.slice(0, 8);
            loadSkinToViewer(e.data.skinUrl, e.data.modelType);
        }
    }
});

document.getElementById('rotate-btn').addEventListener('click', function () {
    autoRotate = !autoRotate;
    this.classList.toggle('active');
});

document.getElementById('slim-btn').addEventListener('click', function () {
    buildModel(scene, !isSlim);
    this.classList.toggle('active');
});

window.addEventListener('resize', function () {
    if (renderer && renderer.domElement && renderer.domElement.parentNode) {
        resizeRenderer(renderer.domElement.parentNode);
    }
});
