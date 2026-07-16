var SKIN_SIZE = 64;
var skinCanvas = document.createElement('canvas');
skinCanvas.width = SKIN_SIZE;
skinCanvas.height = SKIN_SIZE;
var skinCtx = skinCanvas.getContext('2d');

var scene, renderer, camera, playerGroup, skinTexture;
var bodyParts = [];
var isSlim = false;
var selectedPart = null;
var raycaster = new THREE.Raycaster();
var mouse = new THREE.Vector2();

// Face-to-UV mapping for each body part
// Each face maps to [texX, texY, texW, texH] on the 64x64 skin
var PART_REGIONS = {
    head: {
        right:  [0,  8, 8, 8],
        left:   [16, 8, 8, 8],
        top:    [8,  0, 8, 8],
        bottom: [16, 0, 8, 8],
        front:  [8,  8, 8, 8],
        back:   [24, 8, 8, 8]
    },
    body: {
        right:  [16, 20, 4, 12],
        left:   [28, 20, 4, 12],
        top:    [20, 16, 8, 4],
        bottom: [28, 16, 8, 4],
        front:  [20, 20, 8, 12],
        back:   [32, 20, 8, 12]
    },
    right_arm: {
        right:  [40, 20, 4, 12],
        left:   [48, 20, 4, 12],
        top:    [44, 16, 4, 4],
        bottom: [48, 16, 4, 4],
        front:  [44, 20, 4, 12],
        back:   [52, 20, 4, 12]
    },
    left_arm: {
        right:  [32, 52, 4, 12],
        left:   [40, 52, 4, 12],
        top:    [36, 48, 4, 4],
        bottom: [40, 48, 4, 4],
        front:  [36, 52, 4, 12],
        back:   [44, 52, 4, 12]
    },
    right_leg: {
        right:  [0,  20, 4, 12],
        left:   [8,  20, 4, 12],
        top:    [4,  16, 4, 4],
        bottom: [8,  16, 4, 4],
        front:  [4,  20, 4, 12],
        back:   [12, 20, 4, 12]
    },
    left_leg: {
        right:  [16, 52, 4, 12],
        left:   [24, 52, 4, 12],
        top:    [20, 48, 4, 4],
        bottom: [24, 48, 4, 4],
        front:  [20, 52, 4, 12],
        back:   [28, 52, 4, 12]
    }
};

// Adjacent face normals → which face
var FACE_NAMES = {
    '1,0,0': 'right',
    '-1,0,0': 'left',
    '0,1,0': 'top',
    '0,-1,0': 'bottom',
    '0,0,1': 'front',
    '0,0,-1': 'back'
};

function initCoreScene(container) {
    scene = new THREE.Scene();
    var w = container.clientWidth || 600;
    var h = container.clientHeight || 400;
    camera = new THREE.PerspectiveCamera(35, w / h, 1, 200);
    camera.position.set(22, 18, 28);
    camera.lookAt(0, 14, 0);

    renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
    renderer.setSize(w, h);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.setClearColor(0x000000, 0);
    renderer.domElement.style.position = 'absolute';
    renderer.domElement.style.top = '0';
    renderer.domElement.style.left = '0';
    renderer.domElement.style.width = '100%';
    renderer.domElement.style.height = '100%';
    container.style.position = 'relative';
    container.appendChild(renderer.domElement);

    var ambient = new THREE.AmbientLight(0xffffff, 0.5);
    scene.add(ambient);
    var dl = new THREE.DirectionalLight(0xffffff, 0.9);
    dl.position.set(10, 25, 20);
    scene.add(dl);
    var bl = new THREE.DirectionalLight(0xffffff, 0.3);
    bl.position.set(-15, 5, -20);
    scene.add(bl);

    // Floor shadow
    var floor = new THREE.Mesh(
        new THREE.PlaneGeometry(30, 30),
        new THREE.ShadowMaterial({ opacity: 0.12 })
    );
    floor.rotation.x = -Math.PI / 2;
    floor.position.y = -0.5;
    floor.receiveShadow = true;
    scene.add(floor);
}

function buildDefaultSkin() {
    skinCtx.fillStyle = '#8b8b8b';
    skinCtx.fillRect(0, 0, SKIN_SIZE, SKIN_SIZE);
    // Head
    skinCtx.fillStyle = '#d4a574';
    fillSkinRegion(8, 8, 8, 8);
    fillSkinRegion(8, 0, 8, 8);
    fillSkinRegion(0, 8, 8, 8);
    fillSkinRegion(16, 8, 8, 8);
    fillSkinRegion(24, 8, 8, 8);
    fillSkinRegion(16, 0, 8, 8);
    // Eyes
    skinCtx.fillStyle = '#000';
    skinCtx.fillRect(13, 12, 1, 2);
    skinCtx.fillRect(16, 12, 1, 2);
    // Hair
    skinCtx.fillStyle = '#5c3a1e';
    skinCtx.fillRect(8, 8, 8, 3);
    // Body
    skinCtx.fillStyle = '#5c5ce0';
    fillSkinRegion(20, 20, 8, 12);
    fillSkinRegion(16, 20, 4, 12);
    fillSkinRegion(28, 20, 4, 12);
    fillSkinRegion(32, 20, 8, 12);
    // Arms
    skinCtx.fillStyle = '#4a4ac4';
    fillSkinRegion(44, 20, 4, 12);
    fillSkinRegion(40, 20, 4, 12);
    fillSkinRegion(52, 20, 4, 12);
    skinCtx.fillStyle = '#3a3a9e';
    fillSkinRegion(36, 52, 4, 12);
    fillSkinRegion(32, 52, 4, 12);
    fillSkinRegion(40, 52, 4, 12);
    fillSkinRegion(44, 52, 4, 12);
    // Legs
    skinCtx.fillStyle = '#3a3a70';
    fillSkinRegion(4, 20, 4, 12);
    fillSkinRegion(0, 20, 4, 12);
    fillSkinRegion(8, 20, 4, 12);
    fillSkinRegion(12, 20, 4, 12);
    fillSkinRegion(20, 52, 4, 12);
    fillSkinRegion(16, 52, 4, 12);
    fillSkinRegion(24, 52, 4, 12);
    fillSkinRegion(28, 52, 4, 12);
    updateTexture();
}

function fillSkinRegion(x, y, w, h) {
    skinCtx.fillRect(x, y, w, h);
}

function updateTexture() {
    if (skinTexture) skinTexture.dispose();
    skinTexture = new THREE.CanvasTexture(skinCanvas);
    skinTexture.colorSpace = THREE.SRGBColorSpace;
    skinTexture.needsUpdate = true;
    if (playerGroup) {
        playerGroup.children.forEach(function (child) {
            if (child.material) child.material.map = skinTexture;
        });
    }
}

function buildModel(targetScene, slim) {
    isSlim = slim;
    if (playerGroup) {
        targetScene.remove(playerGroup);
        playerGroup = null;
    }

    var mat = new THREE.MeshStandardMaterial({
        map: skinTexture,
        side: THREE.DoubleSide,
        roughness: 0.6,
        metalness: 0.0
    });

    var group = new THREE.Group();
    bodyParts = [];

    var armW = slim ? 3 : 4;

    function addPart(id, cx, cy, cz, w, h, d) {
        var faceOrder = ['right', 'left', 'top', 'bottom', 'front', 'back'];
        var posMap = {
            right: [w / 2, 0, 0], left: [-w / 2, 0, 0],
            top: [0, h / 2, 0], bottom: [0, -h / 2, 0],
            front: [0, 0, d / 2], back: [0, 0, -d / 2]
        };
        var rotMap = {
            right: [0, Math.PI / 2, 0], left: [0, -Math.PI / 2, 0],
            top: [-Math.PI / 2, 0, 0], bottom: [Math.PI / 2, 0, 0],
            front: [0, 0, 0], back: [0, Math.PI, 0]
        };

        var regions = PART_REGIONS[id];

        faceOrder.forEach(function (side) {
            var pw = (side === 'right' || side === 'left') ? d : w;
            var ph = (side === 'top' || side === 'bottom') ? d : h;
            var geo = new THREE.PlaneGeometry(
                (side === 'top' || side === 'bottom') ? w : pw,
                (side === 'top' || side === 'bottom') ? d : ph
            );
            var uv = geo.attributes.uv.array;
            var reg = regions[side];
            var u1 = reg[0] / 64, v1 = (64 - reg[1] - reg[3]) / 64;
            var u2 = (reg[0] + reg[2]) / 64, v2 = (64 - reg[1]) / 64;
            uv[0] = u1; uv[1] = v1;
            uv[2] = u2; uv[3] = v1;
            uv[4] = u1; uv[5] = v2;
            uv[6] = u2; uv[7] = v2;
            geo.attributes.uv.needsUpdate = true;

            var mesh = new THREE.Mesh(geo, mat);
            var p = posMap[side];
            var r = rotMap[side];
            mesh.position.set(cx + p[0], cy + p[1], cz + p[2]);
            mesh.rotation.set(r[0], r[1], r[2]);
            mesh.userData = { partId: id, face: side };
            group.add(mesh);
            bodyParts.push(mesh);
        });
    }

    // Head
    addPart('head', 0, 24, 0, 8, 8, 8);
    // Body
    addPart('body', 0, 12, 0, 8, 12, 4);
    // Right Arm
    addPart('right_arm', armW === 3 ? -5.5 : -6, 12, 0, armW, 12, armW);
    // Left Arm
    addPart('left_arm', armW === 3 ? 5.5 : 6, 12, 0, armW, 12, armW);
    // Right Leg
    addPart('right_leg', -2, 0, 0, 4, 12, 4);
    // Left Leg
    addPart('left_leg', 2, 0, 0, 4, 12, 4);

    playerGroup = group;
    targetScene.add(group);
}

function getClickedPart(event, container) {
    var rect = container.getBoundingClientRect();
    mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
    mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
    raycaster.setFromCamera(mouse, camera);
    var meshes = playerGroup ? playerGroup.children.filter(function (c) { return c.isMesh; }) : [];
    var hits = raycaster.intersectObjects(meshes);
    if (hits.length > 0) {
        var hit = hits[0];
        var faceNormal = hit.face.normal;
        var faceName = FACE_NAMES[Math.round(faceNormal.x) + ',' + Math.round(faceNormal.y) + ',' + Math.round(faceNormal.z)];
        if (faceName && hit.object && hit.object.userData && hit.object.userData.partId) {
            return {
                partId: hit.object.userData.partId,
                face: faceName
            };
        }
    }
    return null;
}

function fillPartRegion(partId) {
    var regions = PART_REGIONS[partId];
    if (!regions) return;
    var color = skinCtx.fillStyle;
    var keys = Object.keys(regions);
    for (var i = 0; i < keys.length; i++) {
        var r = regions[keys[i]];
        skinCtx.fillRect(r[0], r[1], r[2], r[3]);
    }
    updateTexture();
}

function fillFaceRegion(partId, face) {
    var regions = PART_REGIONS[partId];
    if (!regions || !regions[face]) return;
    var r = regions[face];
    skinCtx.fillRect(r[0], r[1], r[2], r[3]);
    updateTexture();
}

function resizeRenderer(container) {
    if (renderer && container) {
        var w = container.clientWidth;
        var h = container.clientHeight;
        if (w === 0 || h === 0) return;
        renderer.setSize(w, h);
        camera.aspect = w / h;
        camera.updateProjectionMatrix();
    }
}

function startRenderLoop(callback) {
    function loop() {
        requestAnimationFrame(loop);
        if (callback) callback();
        renderer.render(scene, camera);
    }
    loop();
}

function attachRendererToContainer(container) {
    if (renderer && renderer.domElement && container) {
        if (renderer.domElement.parentNode) {
            renderer.domElement.parentNode.removeChild(renderer.domElement);
        }
        container.style.position = 'relative';
        container.appendChild(renderer.domElement);
        resizeRenderer(container);
    }
}
