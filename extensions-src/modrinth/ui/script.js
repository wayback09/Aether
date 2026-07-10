// ── IPC Bridge ────────────────────────────────────────────────────────────────
const pending = {};
let reqCounter = 0;

function sendMessage(payload) {
    return new Promise((resolve) => {
        const id = ++reqCounter;
        payload.requestId = id;
        pending[id] = resolve;
        window.parent.postMessage(payload, '*');
    });
}

window.addEventListener('message', (e) => {
    const msg = e.data;
    if (!msg || !msg.requestId) return;
    const resolve = pending[msg.requestId];
    if (resolve) {
        delete pending[msg.requestId];
        resolve(msg);
    }
});

// ── State ─────────────────────────────────────────────────────────────────────
let currentMod = null;
let currentVersions = [];

// ── Elements ──────────────────────────────────────────────────────────────────
const searchInput    = document.getElementById('searchInput');
const resultsDiv     = document.getElementById('resultsContainer');
const modal          = document.getElementById('installModal');
const modalModName   = document.getElementById('modalModName');
const modalModAuthor = document.getElementById('modalModAuthor');
const modalModIcon   = document.getElementById('modalModIcon');
const versionSelect  = document.getElementById('versionSelect');
const instanceSelect = document.getElementById('instanceSelect');
const installBtn     = document.getElementById('installBtn');
const installBtnText = document.getElementById('installBtnText');
const cancelBtn      = document.getElementById('cancelBtn');
const modalClose     = document.getElementById('modalClose');
const installStatus  = document.getElementById('installStatus');

// ── Search ────────────────────────────────────────────────────────────────────
async function search(query) {
    if (!query.trim()) return;
    resultsDiv.innerHTML = '<div class="loading"><div class="spinner"></div><p>Searching Modrinth...</p></div>';

    try {
        const res = await fetch(`https://api.modrinth.com/v2/search?query=${encodeURIComponent(query)}&limit=20&facets=[["project_type:mod"]]`);
        if (!res.ok) throw new Error('API Error ' + res.status);
        const data = await res.json();

        if (!data.hits.length) {
            resultsDiv.innerHTML = '<div class="placeholder-wrap"><div class="placeholder-icon">📦</div><p>No mods found.</p></div>';
            return;
        }

        resultsDiv.innerHTML = '<div class="grid">' + data.hits.map(hit => {
            const icon = hit.icon_url
                ? `<img src="${hit.icon_url}" class="card-icon" alt="" />`
                : `<div class="card-icon card-icon-placeholder">${hit.title.charAt(0)}</div>`;
            const dl = hit.downloads >= 1000000
                ? (hit.downloads / 1000000).toFixed(1) + 'M'
                : hit.downloads >= 1000
                    ? (hit.downloads / 1000).toFixed(0) + 'K'
                    : hit.downloads;
            return `
            <div class="card" data-id="${hit.project_id}" data-title="${encodeURIComponent(hit.title)}" data-author="${encodeURIComponent(hit.author)}" data-icon="${encodeURIComponent(hit.icon_url || '')}">
                <div class="card-top">
                    ${icon}
                    <div class="card-info">
                        <div class="card-title">${hit.title}</div>
                        <div class="card-author">by ${hit.author}</div>
                    </div>
                </div>
                <p class="card-desc">${hit.description}</p>
                <div class="card-footer">
                    <span class="card-downloads">⬇ ${dl}</span>
                    <button class="btn-install-card" data-id="${hit.project_id}">Install</button>
                </div>
            </div>`;
        }).join('') + '</div>';

        // Attach click handlers to Install buttons
        document.querySelectorAll('.btn-install-card').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const card = btn.closest('.card');
                openInstallModal({
                    id: card.dataset.id,
                    title: decodeURIComponent(card.dataset.title),
                    author: decodeURIComponent(card.dataset.author),
                    icon: decodeURIComponent(card.dataset.icon)
                });
            });
        });

    } catch (err) {
        resultsDiv.innerHTML = `<div class="placeholder-wrap"><p>Error: ${err.message}</p></div>`;
    }
}

// Debounce search
let debounceTimer;
searchInput.addEventListener('input', () => {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => search(searchInput.value), 400);
});
searchInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') { clearTimeout(debounceTimer); search(searchInput.value); }
});

// ── Install Modal ─────────────────────────────────────────────────────────────
async function openInstallModal(mod) {
    currentMod = mod;
    currentVersions = [];

    // Populate header
    modalModName.textContent = mod.title;
    modalModAuthor.textContent = 'by ' + mod.author;
    if (mod.icon) {
        modalModIcon.innerHTML = `<img src="${mod.icon}" alt="" />`;
    } else {
        modalModIcon.innerHTML = `<div class="icon-letter">${mod.title.charAt(0)}</div>`;
    }

    // Reset state
    versionSelect.innerHTML = '<option>Loading versions...</option>';
    instanceSelect.innerHTML = '<option>Loading instances...</option>';
    installStatus.classList.add('hidden');
    installStatus.textContent = '';
    installBtnText.textContent = 'Install';
    installBtn.disabled = false;
    modal.classList.remove('hidden');

    // Fetch versions and instances in parallel
    const [versionsRes, instancesMsg] = await Promise.all([
        fetch(`https://api.modrinth.com/v2/project/${mod.id}/version`).then(r => r.json()),
        sendMessage({ type: 'get_instances' })
    ]);

    // Populate versions
    currentVersions = versionsRes;
    versionSelect.innerHTML = versionsRes.length
        ? versionsRes.map((v, i) =>
            `<option value="${i}">${v.version_number} — ${v.name} (${v.game_versions.slice(0,3).join(', ')})</option>`
          ).join('')
        : '<option>No versions found</option>';

    // Populate instances
    const instances = instancesMsg.instances || [];
    instanceSelect.innerHTML = instances.length
        ? instances.map(inst =>
            `<option value="${inst.id}">${inst.name} (${inst.version} • ${inst.loader})</option>`
          ).join('')
        : '<option value="">No instances found</option>';
}

function closeModal() {
    modal.classList.add('hidden');
    currentMod = null;
}

modalClose.addEventListener('click', closeModal);
cancelBtn.addEventListener('click', closeModal);
modal.addEventListener('click', (e) => { if (e.target === modal) closeModal(); });

installBtn.addEventListener('click', async () => {
    const vIdx = parseInt(versionSelect.value, 10);
    const instanceId = instanceSelect.value;

    if (isNaN(vIdx) || !instanceId) {
        showStatus('Please select a version and instance.', 'error');
        return;
    }

    const version = currentVersions[vIdx];
    // Pick the primary jar file
    const file = version.files.find(f => f.primary) || version.files[0];
    if (!file) {
        showStatus('No downloadable file found for this version.', 'error');
        return;
    }

    installBtnText.textContent = 'Installing...';
    installBtn.disabled = true;
    installStatus.classList.add('hidden');

    const result = await sendMessage({
        type: 'install_mod',
        instanceId,
        jarName: file.filename,
        downloadUrl: file.url
    });

    if (result.success) {
        showStatus(`✓ ${file.filename} installed successfully!`, 'success');
        installBtnText.textContent = 'Done!';
    } else {
        showStatus(`✗ ${result.error}`, 'error');
        installBtnText.textContent = 'Install';
        installBtn.disabled = false;
    }
});

function showStatus(msg, type) {
    installStatus.textContent = msg;
    installStatus.className = `install-status ${type}`;
}

// Initial search
search('fabric');
