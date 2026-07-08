<script lang="ts">
  import { onMount } from 'svelte';
  import { GetInstances, GetAvailableVersions, CreateInstance } from '../../wailsjs/go/main/App.js';
  import Dropdown from '../components/Dropdown.svelte';

  let instances: any[] = [];
  let showModal = false;
  let newInstance = { name: "", version: "", loader: "Vanilla" };
  let availableVersions: string[] = [];
  let isCreating = false;

  async function loadInstances() {
    const res = await GetInstances();
    instances = res || [];
  }

  onMount(async () => {
    await loadInstances();
    availableVersions = await GetAvailableVersions();
    if (availableVersions.length > 0) {
      newInstance.version = availableVersions[0];
    }
  });

  async function handleCreate() {
    if (!newInstance.name || !newInstance.version) return;
    isCreating = true;
    try {
      await CreateInstance(newInstance.name, newInstance.version, newInstance.loader);
      showModal = false;
      newInstance.name = "";
      await loadInstances();
    } catch (err) {
      console.error(err);
      alert("Failed to create instance: " + err);
    } finally {
      isCreating = false;
    }
  }
</script>

<div class="page">
  <header class="page-header">
    <h1>Instances</h1>
    <div class="actions">
      <button class="btn btn-secondary">Import</button>
      <button class="btn btn-primary" on:click={() => showModal = true}>Create New</button>
    </div>
  </header>

  {#if instances.length === 0}
    <div class="empty-state">
      <h3>No Instances Found</h3>
      <p>Create or import a Minecraft instance to get started.</p>
      <button class="btn btn-primary" style="margin-top: var(--spacing-md)" on:click={() => showModal = true}>Create New Instance</button>
    </div>
  {:else}
    <div class="grid">
      {#each instances as instance}
        <div class="card instance-card">
          <div class="card-content">
            <div class="instance-title">{instance.name}</div>
            <div class="instance-meta">
              {instance.version} • {instance.loader}
            </div>
            <div class="instance-last-played">
              Last played: {instance.lastPlayed}
            </div>
          </div>
          <div class="card-actions">
            <button class="btn btn-primary">Play</button>
            <button class="btn btn-secondary">Details</button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if showModal}
    <div class="modal-backdrop" on:click={() => showModal = false}>
      <div class="modal" on:click|stopPropagation>
        <h2>Create Instance</h2>
        
        <div class="form-group">
          <label>Instance Name</label>
          <input type="text" bind:value={newInstance.name} placeholder="e.g. My Survival World" />
        </div>

        <div class="form-group">
          <label>Version</label>
          <Dropdown 
            options={availableVersions.map(v => ({ label: v, value: v }))} 
            bind:value={newInstance.version} 
            direction="up"
          />
        </div>

        <div class="form-group">
          <label>Mod Loader</label>
          <Dropdown 
            options={[
              { label: 'Vanilla', value: 'Vanilla' },
              { label: 'Fabric (Coming Soon)', value: 'Fabric', disabled: true }
            ]} 
            bind:value={newInstance.loader} 
            direction="up"
          />
        </div>

        <div class="modal-actions">
          <button class="btn btn-secondary" on:click={() => showModal = false}>Cancel</button>
          <button class="btn btn-primary" on:click={handleCreate} disabled={isCreating || !newInstance.name}>
            {isCreating ? 'Creating...' : 'Create'}
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .page {
    padding: var(--spacing-xl);
    height: 100%;
    box-sizing: border-box;
    overflow-y: auto;
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--spacing-xl);
  }
  
  .actions {
    display: flex;
    gap: var(--spacing-md);
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 50vh;
    color: var(--text-secondary);
  }

  .empty-state h3 {
    color: var(--text-primary);
    margin-bottom: var(--spacing-sm);
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: var(--spacing-lg);
  }

  .instance-card {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .card-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .instance-title {
    font-size: 18px;
    font-weight: 600;
  }

  .instance-meta {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .instance-last-played {
    font-size: 12px;
    color: var(--text-secondary);
    margin-top: var(--spacing-sm);
  }

  .card-actions {
    display: flex;
    gap: var(--spacing-sm);
    margin-top: auto;
  }
  
  .card-actions .btn {
    flex: 1;
  }

  /* Modal Styles */
  .modal-backdrop {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0,0,0,0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
  }

  .modal {
    background: rgba(25, 25, 25, 0.7);
    backdrop-filter: blur(24px);
    -webkit-backdrop-filter: blur(24px);
    padding: var(--spacing-xl);
    border-radius: var(--border-radius);
    width: 400px;
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
    box-shadow: 0 10px 40px rgba(0,0,0,0.6);
    border: 1px solid rgba(255,255,255,0.1);
  }

  .modal h2 {
    margin: 0;
    font-size: 20px;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .form-group label {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .form-group input {
    background: rgba(0,0,0,0.2);
    border: 1px solid rgba(255,255,255,0.1);
    color: white;
    padding: 10px 12px;
    border-radius: var(--border-radius);
    font-family: inherit;
    font-size: 14px;
    outline: none;
    transition: border-color 0.2s;
  }

  .form-group input:focus {
    border-color: var(--accent);
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-md);
    margin-top: var(--spacing-sm);
  }
</style>
