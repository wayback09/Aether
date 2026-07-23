// Skin colors palette (common Minecraft colors)
var COLOR_PALETTE = [
    // Skin tones
    '#d4a574', '#c49a6c', '#b08050', '#8b6348', '#5c3a1e', '#2d1b0e',
    // Hair
    '#1c1c1c', '#3a1e0c', '#543819', '#8b6914', '#c4a434', '#e5d4a0',
    // Eyes
    '#000000', '#333333', '#5555cc', '#33aa33',
    // Shirts
    '#5555ff', '#4a4acc', '#cc3333', '#33cc33', '#ffaa00', '#ffffff',
    // Pants
    '#333366', '#444488', '#666666', '#333333',
    // Accents
    '#ff5555', '#ffaa00', '#c4a434', '#55ff55', '#55ffff', '#ff55ff'
];

var partNames = {
    head: 'Head', body: 'Body',
    right_arm: 'Right Arm', left_arm: 'Left Arm',
    right_leg: 'Right Leg', left_leg: 'Left Leg'
};

var selectedPartId = null;
var selectedFace = null;
var colorHistory = ['#d4a574', '#5c3a1e', '#000000', '#5555ff', '#333366', '#ffffff'];

function showNotice(message) {
    var notice = document.getElementById('notice');
    notice.textContent = message;
    notice.classList.remove('hidden');
    clearTimeout(showNotice.timer);
    showNotice.timer = setTimeout(function () {
        notice.classList.add('hidden');
    }, 3500);
}

// Build the color grid
function buildColorGrid() {
    var grid = document.getElementById('color-grid');
    COLOR_PALETTE.forEach(function (c) {
        var btn = document.createElement('button');
        btn.className = 'color-swatch';
        btn.style.background = c;
        btn.dataset.color = c;
        btn.addEventListener('click', function () {
            var color = this.dataset.color;
            skinCtx.fillStyle = color;
            document.getElementById('custom-color').value = color;
            document.querySelectorAll('.color-swatch').forEach(function (b) {
                b.classList.remove('selected');
            });
            this.classList.add('selected');
            if (selectedPartId) {
                fillPartRegion(selectedPartId);
                updateTexture();
                updateSelectedInfo();
            }
        });
        grid.appendChild(btn);
    });
}

// Part buttons
function buildPartButtons() {
    var container = document.getElementById('part-buttons');
    Object.keys(partNames).forEach(function (id) {
        var btn = document.createElement('button');
        btn.className = 'part-btn';
        btn.dataset.part = id;
        btn.textContent = partNames[id];
        btn.addEventListener('click', function () {
            selectPart(this.dataset.part);
        });
        container.appendChild(btn);
    });
}

function selectPart(partId) {
    selectedPartId = partId;
    selectedFace = null;
    document.querySelectorAll('.part-btn').forEach(function (b) {
        b.classList.toggle('active', b.dataset.part === partId);
    });
    updateSelectedInfo();
}

function updateSelectedInfo() {
    var info = document.getElementById('selected-info');
    if (selectedPartId) {
        info.textContent = partNames[selectedPartId] + (selectedFace ? ' — ' + selectedFace : '');
    } else {
        info.textContent = 'Click on the model or select a part above';
    }
}

// Click on 3D model to select part
var activeContainer = null;
var makerActive = false;

function enableModelClicking(container) {
    activeContainer = container;
    makerActive = true;
    container.addEventListener('click', onModelClick);
}

function onModelClick(e) {
    if (!makerActive) return;
    var result = getClickedPart(e, activeContainer);
    if (result) {
        selectPart(result.partId);
        selectedFace = result.face;
        updateSelectedInfo();
    }
}

// Fill buttons
document.getElementById('fill-part-btn').addEventListener('click', function () {
    if (!selectedPartId) return;
    fillPartRegion(selectedPartId);
    updateTexture();
    updateSelectedInfo();
});

document.getElementById('fill-face-btn').addEventListener('click', function () {
    if (!selectedPartId || !selectedFace) return;
    fillFaceRegion(selectedPartId, selectedFace);
    updateTexture();
    updateSelectedInfo();
});

// Custom color
document.getElementById('custom-color').addEventListener('input', function () {
    skinCtx.fillStyle = this.value;
    document.querySelectorAll('.color-swatch').forEach(function (b) {
        b.classList.toggle('selected', b.dataset.color === this.value);
    }.bind(this));
});

document.getElementById('apply-custom').addEventListener('click', function () {
    if (selectedPartId) {
        fillPartRegion(selectedPartId);
        updateTexture();
        updateSelectedInfo();
    }
});

// Import skin into editor
document.getElementById('import-editor-btn').addEventListener('click', function () {
    var name = document.getElementById('username-input').value.trim();
    if (!name) { showNotice('Enter a username in the search bar first'); return; }
    window.parent.postMessage({ type: "fetch_profile", username: name, reqId: 'editor' }, '*');
});

window.addEventListener('message', function (e) {
    if (e.data && e.data.type === 'profile_result' && e.data.reqId === 'editor') {
        if (e.data.error) { showNotice(e.data.error); return; }
        if (!e.data.skinUrl) { showNotice('No skin found'); return; }
        var img = new Image();
        img.crossOrigin = 'Anonymous';
        img.onload = function () {
            skinCtx.drawImage(img, 0, 0, SKIN_SIZE, SKIN_SIZE);
            updateTexture();
            buildModel(scene, e.data.modelType === 'slim');
        };
        img.onerror = function () { showNotice('Failed to load skin'); };
        img.src = e.data.skinUrl;
    }
});

// Export
document.getElementById('export-btn').addEventListener('click', function () {
    var exportCanvas = document.createElement('canvas');
    exportCanvas.width = SKIN_SIZE;
    exportCanvas.height = SKIN_SIZE;
    var ctx = exportCanvas.getContext('2d');
    ctx.drawImage(skinCanvas, 0, 0);
    var link = document.createElement('a');
    link.download = 'skin.png';
    link.href = exportCanvas.toDataURL('image/png');
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
});

// Clear / New Skin
document.getElementById('clear-skin-btn').addEventListener('click', function () {
    skinCtx.fillStyle = '#8b8b8b';
    skinCtx.fillRect(0, 0, SKIN_SIZE, SKIN_SIZE);
    skinCtx.fillStyle = '#d4a574';
    fillSkinRegion(8, 8, 8, 8);
    fillSkinRegion(8, 0, 8, 8);
    skinCtx.fillRect(13, 12, 1, 2);
    skinCtx.fillRect(16, 12, 1, 2);
    updateTexture();
    buildModel(scene, false);
});

// Tab switching — move the 3D view with the tab
document.querySelectorAll('.tab-btn').forEach(function (btn) {
    btn.addEventListener('click', function () {
        document.querySelectorAll('.tab-btn, .tab-content').forEach(function (el) {
            el.classList.remove('active');
        });
        this.classList.add('active');
        var tab = document.getElementById('tab-' + this.dataset.tab);
        tab.classList.add('active');
        var container = tab.querySelector('.model-container');
        if (container) {
            attachRendererToContainer(container);
            makerActive = this.dataset.tab === 'maker';
            if (makerActive) {
                enableModelClicking(container);
            }
        }
    });
});

// Init
buildColorGrid();
buildPartButtons();
updateSelectedInfo();
