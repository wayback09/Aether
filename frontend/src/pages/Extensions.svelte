<script lang="ts">
  import { onMount } from 'svelte';
  import { GetExtensions, SelectAndInstallExtension } from '../../wailsjs/go/main/App.js';
  import EmptyState from '../components/EmptyState.svelte';

  let extensions: any[] = [];
  let isInstalling = false;

  async function loadExtensions() {
    try {
      const exts = await GetExtensions();
      extensions = exts || [];
    } catch (e) {
      console.error(e);
    }
  }

  async function handleLocalInstall() {
    if (isInstalling) return;
    isInstalling = true;
    try {
      const installed = await SelectAndInstallExtension();
      if (installed) await loadExtensions();
    } catch (e) {
      console.error('Installation failed:', e);
    } finally {
      isInstalling = false;
    }
  }

  onMount(loadExtensions);

  function trustBadge(trust: string | undefined): { cls: string; label: string } {
    switch (trust) {
      case 'official':  return { cls: 'badge-official',  label: 'Official' };
      case 'verified':  return { cls: 'badge-verified',  label: 'Verified' };
      case 'local':     return { cls: 'badge-experimental', label: 'Local' };
      case 'community':
      default:          return { cls: 'badge-community', label: 'Community' };
    }
  }

  function extGradient(name: string): string {
    const g = [
      'linear-gradient(135deg, #3b82f6, #1d4ed8)',
      'linear-gradient(135deg, #8b5cf6, #6d28d9)',
      'linear-gradient(135deg, #06b6d4, #0284c7)',
      'linear-gradient(135deg, #10b981, #047857)',
      'linear-gradient(135deg, #f59e0b, #b45309)',
      'linear-gradient(135deg, #ec4899, #be185d)',
    ];
    let h = 0;
    for (let i = 0; i < name.length; i++) h = name.charCodeAt(i) + ((h << 5) - h);
    return g[Math.abs(h) % g.length];
  }

  function statusColor(status: string): string {
    if (!status) return 'rgba(255,255,255,0.2)';
    const s = status.toLowerCase();
    if (s === 'running') return '#22c55e';
    if (s === 'error')   return '#ef4444';
    return 'rgba(255,255,255,0.2)';
  }
</script>

<div class="page page-enter">
  <header class="page-header">
    <h1>Extensions</h1>
    <button class="btn btn-secondary" on:click={handleLocalInstall} disabled={isInstalling}>
      {isInstalling ? 'Installing...' : 'Install from .zip'}
    </button>
  </header>

  {#if extensions.length === 0}
    <EmptyState
      icon="puzzle"
      title="No extensions installed"
      description="Install a .zip extension to add new capabilities to Aether. The Extension Gallery is coming soon."
      actionLabel="Install from .zip"
      on:action={handleLocalInstall}
    />
  {:else}
    <div class="grid">
      {#each extensions as ext}
        {@const badge = trustBadge(ext.trust)}
        {@const grad  = extGradient(ext.name)}
        {@const dot   = statusColor(ext.status)}

        <div class="card ext-card">
          <div class="card-accent" style="background: {grad};"></div>
          <div class="card-body">
            <div class="ext-header">
              <div class="ext-icon" style={!ext.iconUrl ? `background: ${grad};` : ''}>
                {#if ext.iconUrl}
                  <img src={ext.iconUrl} alt={ext.name} class="ext-icon-img" />
                {:else}
                  <span class="ext-icon-letter">{ext.name.charAt(0).toUpperCase()}</span>
                {/if}
              </div>
              <div class="ext-info">
                <div class="ext-name-row">
                  <span class="ext-name">{ext.name}</span>
                </div>
                <div class="ext-badges">
                  <span class="badge {badge.cls}">{badge.label}</span>
                  {#if ext.version}
                    <span class="badge badge-version">v{ext.version}</span>
                  {/if}
                </div>
                <div class="ext-author">by {ext.author}</div>
              </div>
            </div>

            <div class="ext-stats">
              <div class="stat">
                <span class="status-dot" style="background: {dot};"></span>
                <span class="stat-value">{ext.status}</span>
              </div>
              <div class="stat">
                <span class="stat-label">Memory</span>
                <span class="stat-value">{ext.memory || '0 MB'}</span>
              </div>
              <div class="stat">
                <span class="stat-label">CPU</span>
                <span class="stat-value">{ext.cpu || '0%'}</span>
              </div>
            </div>

            <div class="card-actions">
              <button class="btn btn-secondary">Settings</button>
              <button class="btn btn-secondary">Restart</button>
              <button class="btn btn-danger">Disable</button>
            </div>
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

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
    gap: var(--spacing-lg);
  }

  .ext-card {
    position: relative;
    padding: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .card-accent {
    height: 3px;
    width: 100%;
    flex-shrink: 0;
  }

  .card-body {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
    padding: var(--spacing-md);
    flex: 1;
  }

  .ext-header {
    display: flex;
    gap: var(--spacing-md);
    align-items: flex-start;
  }

  .ext-icon {
    width: 44px;
    height: 44px;
    flex-shrink: 0;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  .ext-icon-letter {
    font-size: 20px;
    font-weight: 800;
    color: rgba(255, 255, 255, 0.92);
    line-height: 1;
  }

  .ext-icon-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 10px;
  }

  .ext-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 0;
  }

  .ext-name-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .ext-name {
    font-size: 15px;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .ext-badges {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
  }

  .ext-author {
    font-size: 12px;
    color: var(--text-meta);
  }

  .ext-stats {
    display: flex;
    gap: var(--spacing-md);
    background: rgba(0, 0, 0, 0.18);
    border: 1px solid rgba(255, 255, 255, 0.04);
    padding: 10px 12px;
    border-radius: var(--border-radius);
    font-size: 12px;
  }

  .stat {
    display: flex;
    align-items: center;
    gap: 5px;
  }

  .stat-label {
    color: var(--text-secondary);
  }

  .stat-value {
    color: var(--text-primary);
    font-weight: 500;
  }

  .status-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .card-actions {
    display: flex;
    gap: var(--spacing-sm);
    margin-top: auto;
  }

  .card-actions .btn {
    flex: 1;
  }
</style>
