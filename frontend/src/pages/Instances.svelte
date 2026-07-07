<script lang="ts">
  import { onMount } from 'svelte';
  import { GetInstances } from '../../wailsjs/go/main/App.js';

  let instances: any[] = [];

  onMount(async () => {
    instances = await GetInstances();
  });
</script>

<div class="page">
  <header class="page-header">
    <h1>Instances</h1>
    <div class="actions">
      <button class="btn btn-secondary">Import</button>
      <button class="btn btn-primary">Create New</button>
    </div>
  </header>

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
</style>
