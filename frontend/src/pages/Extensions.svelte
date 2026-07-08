<script lang="ts">
  import { onMount } from 'svelte';
  import { GetExtensions, SelectAndInstallExtension } from '../../wailsjs/go/main/App.js';

  let extensions: any[] = [];
  let isInstalling = false;

  async function loadExtensions() {
    const res = await GetExtensions();
    extensions = res || [];
  }

  async function handleInstall() {
    if (isInstalling) return;
    isInstalling = true;
    try {
      const installed = await SelectAndInstallExtension();
      if (installed) {
        await loadExtensions(); // Refresh the list
      }
    } catch (e) {
      console.error("Installation failed:", e);
      // TODO: show notification via Wails event or UI store
    } finally {
      isInstalling = false;
    }
  }

  onMount(loadExtensions);
</script>

<div class="page">
  <header class="page-header">
    <h1>Extensions</h1>
    <button class="btn btn-primary" on:click={handleInstall} disabled={isInstalling}>
      {isInstalling ? 'Installing...' : 'Browse Extensions'}
    </button>
  </header>

  {#if extensions.length === 0}
    <div class="empty-state">
      <h3>No Extensions Installed</h3>
      <p>Install your first extension to add new functionality.</p>
      <button class="btn btn-primary" style="margin-top: var(--spacing-md)" on:click={handleInstall} disabled={isInstalling}>
        {isInstalling ? 'Installing...' : 'Browse Extensions'}
      </button>
    </div>
  {:else}
    <div class="grid">
      {#each extensions as ext}
        <div class="card ext-card">
          <div class="ext-header">
            <div class="ext-icon"></div>
            <div class="ext-info">
              <div class="ext-name">{ext.name}</div>
              <div class="ext-meta">v{ext.version} by {ext.author}</div>
            </div>
          </div>
          
          <div class="ext-stats">
            <div class="stat"><span class="label">Status:</span> {ext.status}</div>
            <div class="stat"><span class="label">Memory:</span> {ext.memory}</div>
            <div class="stat"><span class="label">CPU:</span> {ext.cpu}</div>
          </div>
          
          <div class="card-actions">
            <button class="btn btn-secondary">Settings</button>
            <button class="btn btn-secondary">Restart</button>
            <button class="btn btn-danger">Disable</button>
          </div>
        </div>
      {/each}
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
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: var(--spacing-lg);
  }

  .ext-card {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .ext-header {
    display: flex;
    gap: var(--spacing-md);
    align-items: center;
  }

  .ext-icon {
    width: 48px;
    height: 48px;
    border-radius: var(--border-radius);
    background-color: rgba(255,255,255,0.1);
  }

  .ext-info {
    display: flex;
    flex-direction: column;
  }

  .ext-name {
    font-size: 16px;
    font-weight: 600;
  }

  .ext-meta {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .ext-stats {
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 13px;
    background-color: rgba(0,0,0,0.2);
    padding: var(--spacing-md);
    border-radius: var(--border-radius);
  }
  
  .stat .label {
    color: var(--text-secondary);
    display: inline-block;
    width: 60px;
  }

  .card-actions {
    display: flex;
    gap: var(--spacing-sm);
    margin-top: auto;
  }
</style>
