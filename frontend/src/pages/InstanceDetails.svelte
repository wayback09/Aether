<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import { GetInstances, UpdateInstance, DeleteInstance, LaunchInstance } from '../../wailsjs/go/main/App.js';

  export let instanceId = '';

  const dispatch = createEventDispatcher();
  let instance: any = null;

  // Form fields
  let editName = '';
  let editMemory = '2048'; // Using dropdown for memory as per plan (2G, 4G, 8G)

  const memoryOptions = [
    { label: '2 GB', value: '2048' },
    { label: '4 GB', value: '4096' },
    { label: '6 GB', value: '6144' },
    { label: '8 GB', value: '8192' },
    { label: '12 GB', value: '12288' },
    { label: '16 GB', value: '16384' },
  ];

  onMount(async () => {
    await loadInstance();
  });

  async function loadInstance() {
    const all = await GetInstances();
    instance = all.find((i: any) => i.id === instanceId);
    if (instance) {
      editName = instance.name;
      editMemory = instance.memory || '2048';
    }
  }

  async function saveChanges() {
    if (!instance) return;
    instance.name = editName;
    instance.memory = editMemory;
    
    try {
      await UpdateInstance(instance);
      // Go back to instances page
      dispatch('navigate', 'instances');
    } catch (e) {
      console.error("Failed to save instance:", e);
    }
  }

  async function deleteInstance() {
    if (!instance) return;
    if (confirm(`Are you sure you want to delete ${instance.name}? This cannot be undone.`)) {
      try {
        await DeleteInstance(instance.id);
        dispatch('navigate', 'instances');
      } catch (e) {
        console.error("Failed to delete instance:", e);
      }
    }
  }

  function launch() {
    if (!instance) return;
    LaunchInstance(instance.id);
    dispatch('navigate', 'home'); // switch to home to see status
  }

  // Consistent gradient generator based on ID
  function generateGradient(id: string) {
    let hash = 0;
    for (let i = 0; i < id.length; i++) {
      hash = id.charCodeAt(i) + ((hash << 5) - hash);
    }
    const hue1 = hash % 360;
    const hue2 = (hash * 2) % 360;
    return `linear-gradient(135deg, hsl(${hue1}, 70%, 60%), hsl(${hue2}, 70%, 40%))`;
  }
</script>

<div class="page page-enter">
  {#if instance}
    <div class="header">
      <button class="btn btn-secondary back-btn" on:click={() => dispatch('navigate', 'instances')}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5M12 19l-7-7 7-7"/>
        </svg>
        Back
      </button>

      <div class="header-content">
        <div class="art-square" style="background: {generateGradient(instance.id)}"></div>
        <div class="info">
          <h1>{instance.name}</h1>
          <p class="meta">
            <span class="badge badge-version">{instance.version}</span>
            {#if instance.loader === 'Fabric'}
              <span class="badge badge-community">Fabric</span>
            {:else}
              <span class="badge badge-official">Vanilla</span>
            {/if}
            <span class="last-played">{instance.lastPlayed ? `Last played ${instance.lastPlayed.split('T')[0]}` : 'Never played'}</span>
          </p>
        </div>
      </div>
    </div>

    <div class="settings-card card">
      <h2>Instance Settings</h2>
      
      <div class="form-group">
        <label for="name">Name</label>
        <input id="name" type="text" bind:value={editName} class="input" />
      </div>

      <div class="form-group">
        <label for="memory">Memory Allocation</label>
        <select id="memory" bind:value={editMemory} class="input select">
          {#each memoryOptions as opt}
            <option value={opt.value}>{opt.label}</option>
          {/each}
        </select>
      </div>

      <div class="actions">
        <button class="btn btn-danger" on:click={deleteInstance}>Delete Instance</button>
        <div class="right-actions">
          <button class="btn btn-secondary" on:click={launch}>Play</button>
          <button class="btn btn-primary" on:click={saveChanges}>Save Changes</button>
        </div>
      </div>
    </div>
  {:else}
    <div class="loading">Loading instance...</div>
  {/if}
</div>

<style>
  .page {
    padding: var(--spacing-xl);
    flex-grow: 1;
    overflow-y: auto;
  }

  .back-btn {
    margin-bottom: var(--spacing-lg);
    gap: var(--spacing-sm);
  }

  .header-content {
    display: flex;
    align-items: center;
    gap: var(--spacing-lg);
    margin-bottom: var(--spacing-xl);
  }

  .art-square {
    width: 96px;
    height: 96px;
    border-radius: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  }

  .info h1 {
    font-size: 32px;
    margin: 0 0 8px 0;
    color: var(--text-primary);
  }

  .meta {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--text-meta);
    font-size: 13px;
  }

  .last-played {
    color: var(--text-secondary);
  }

  .settings-card {
    max-width: 600px;
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .settings-card h2 {
    font-size: 18px;
    font-weight: 600;
    margin: 0;
    color: var(--text-primary);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    padding-bottom: var(--spacing-md);
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  label {
    font-size: 13px;
    font-weight: 500;
    color: var(--text-meta);
  }

  .input {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: var(--text-primary);
    padding: 10px 12px;
    border-radius: var(--border-radius);
    font-size: 14px;
    font-family: inherit;
    outline: none;
    transition: border-color var(--transition-fast);
  }

  .input:focus {
    border-color: var(--accent-color);
  }

  .select {
    appearance: none;
    cursor: pointer;
  }

  .select option {
    background: var(--panel-bg);
    color: var(--text-primary);
  }

  .actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: var(--spacing-md);
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    padding-top: var(--spacing-lg);
  }

  .right-actions {
    display: flex;
    gap: var(--spacing-sm);
  }
</style>
