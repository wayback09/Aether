const elements = {
    instancesDropdown: document.getElementById('instances'),
    loadingState: document.getElementById('loading'),
    emptyState: document.getElementById('empty-state'),
    modsContainer: document.querySelector('.mods-container'),
    modRowTemplate: document.getElementById('mod-row-template'),
    confirmation: document.getElementById('confirmation')
};

let currentInstanceId = null;
let pendingDelete = null;

// Helper: Send a message to the Goja backend
function sendMessage(action, payload = {}) {
    window.parent.postMessage({ action, ...payload }, '*');
}

// Format jar name to look nicer (e.g., "fabric-api-0.90.0+1.20.1.jar" -> "Fabric API")
function formatModName(filename) {
    let name = filename.replace(/\.jar(\.disabled)?$/, '');
    // Replace dashes and underscores with spaces
    name = name.replace(/[-_]/g, ' ');
    // Capitalize words
    name = name.replace(/\b\w/g, l => l.toUpperCase());
    return name;
}

// Render the list of mods
function renderMods(mods) {
    elements.modsContainer.innerHTML = '';
    
    if (!mods || mods.length === 0) {
        elements.loadingState.classList.add('hidden');
        elements.emptyState.classList.remove('hidden');
        return;
    }

    elements.emptyState.classList.add('hidden');
    elements.loadingState.classList.add('hidden');

    mods.forEach(modFilename => {
        const isEnabled = !modFilename.endsWith('.disabled');
        const clone = elements.modRowTemplate.content.cloneNode(true);
        const row = clone.querySelector('.mod-row');
        
        if (!isEnabled) {
            row.classList.add('disabled');
        }

        row.querySelector('.mod-name').textContent = formatModName(modFilename);
        row.querySelector('.mod-filename').textContent = modFilename;
        
        const toggle = row.querySelector('.mod-toggle');
        toggle.checked = isEnabled;
        
        // Toggle event listener
        toggle.addEventListener('change', (e) => {
            const enable = e.target.checked;
            sendMessage('toggle_mod', {
                instanceId: currentInstanceId,
                jarName: modFilename,
                enable: enable
            });
            // Optimistic update
            if (enable) {
                row.classList.remove('disabled');
            } else {
                row.classList.add('disabled');
            }
        });

        // Delete event listener
        const deleteBtn = row.querySelector('.delete-btn');
        deleteBtn.addEventListener('click', () => {
            pendingDelete = { modFilename, row };
            elements.confirmation.classList.remove('hidden');
        });

        elements.modsContainer.appendChild(row);
    });
}

// Listen for messages from the Goja backend
window.addEventListener('message', (event) => {
    const data = event.data;
    
    if (data.action === 'get_instances' && data.success) {
        const instances = data.data;
        elements.instancesDropdown.innerHTML = '';
        
        if (instances.length === 0) {
            elements.instancesDropdown.innerHTML = '<option>No instances found</option>';
            elements.loadingState.textContent = 'Create an instance first to manage mods.';
            return;
        }

        instances.forEach(inst => {
            const option = document.createElement('option');
            option.value = inst.id;
            option.textContent = `${inst.name} (${inst.loader} ${inst.version})`;
            elements.instancesDropdown.appendChild(option);
        });

        // Select the first instance automatically
        currentInstanceId = instances[0].id;
        sendMessage('get_mods', { instanceId: currentInstanceId });
    }

    if (data.action === 'get_mods' && data.success) {
        renderMods(data.data);
    }
});

// Dropdown change listener
elements.instancesDropdown.addEventListener('change', (e) => {
    currentInstanceId = e.target.value;
    if (currentInstanceId) {
        elements.modsContainer.innerHTML = '';
        elements.emptyState.classList.add('hidden');
        elements.loadingState.classList.remove('hidden');
        elements.loadingState.textContent = 'Loading mods...';
        sendMessage('get_mods', { instanceId: currentInstanceId });
    }
});

// Init
sendMessage('get_instances');

document.getElementById('cancel-delete').addEventListener('click', () => {
    pendingDelete = null;
    elements.confirmation.classList.add('hidden');
});

document.getElementById('confirm-delete').addEventListener('click', () => {
    if (!pendingDelete) return;
    const { modFilename, row } = pendingDelete;
    sendMessage('delete_mod', {
        instanceId: currentInstanceId,
        jarName: modFilename
    });
    row.remove();
    if (elements.modsContainer.children.length === 0) {
        elements.emptyState.classList.remove('hidden');
    }
    pendingDelete = null;
    elements.confirmation.classList.add('hidden');
});

elements.confirmation.addEventListener('click', (event) => {
    if (event.target === elements.confirmation) {
        pendingDelete = null;
        elements.confirmation.classList.add('hidden');
    }
});
