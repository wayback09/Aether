<script lang="ts">
  import { onMount } from 'svelte';
  import { GetActiveInstance, LaunchInstance } from '../../wailsjs/go/main/App.js';
  import { EventsOn } from '../../wailsjs/runtime/runtime.js';

  let currentInstance: any = null;
  let launchState = "Idle";
  let logs: string[] = [];

  onMount(async () => {
    currentInstance = await GetActiveInstance();

    EventsOn("instance:state", (data: any) => {
      if (currentInstance && data.id === currentInstance.id) {
        launchState = data.state;
      }
    });

    EventsOn("instance:log", (line: string) => {
      logs = [...logs, line].slice(-10); // keep last 10
    });
  });

  async function handlePlay() {
    if (!currentInstance) return;
    launchState = "Starting...";
    try {
      await LaunchInstance(currentInstance.id);
    } catch (err) {
      launchState = "Error";
      console.error(err);
    }
  }
</script>

<div class="page home-page">
  <div class="spacer"></div>
  
  <div class="launch-container">
    {#if currentInstance}
      <div class="instance-info">
        <div class="instance-name">{currentInstance.name}</div>
        <div class="instance-meta">
          <span>{currentInstance.version}</span> • 
          <span>{currentInstance.loader}</span> • 
          <span>{currentInstance.memory}</span>
        </div>
      </div>
      
      <div style="display: flex; gap: var(--spacing-md); align-items: center;">
        <button class="btn btn-primary play-btn" on:click={handlePlay} disabled={launchState === 'Running'}>
          {launchState === 'Running' ? 'Running' : 'Play'}
        </button>
        <span style="color: var(--text-secondary); font-size: 14px;">{launchState !== 'Idle' ? launchState : ''}</span>
      </div>

      {#if logs.length > 0}
        <div style="font-family: monospace; font-size: 12px; color: var(--text-secondary); margin-top: 20px; background: rgba(0,0,0,0.2); padding: 10px; border-radius: 4px; width: 100%;">
          {#each logs as log}
            <div>{log}</div>
          {/each}
        </div>
      {/if}
    {:else}
      <div class="instance-info">Loading...</div>
    {/if}
  </div>
  
  <!-- "Large whitespace below" as requested in UI.md -->
</div>

<style>
  .page {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: var(--spacing-xl);
    box-sizing: border-box;
  }
  
  .spacer {
    flex-grow: 1; /* Pushes launch container down slightly or centers it */
  }

  .launch-container {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-lg);
    margin-bottom: 20vh; /* Keeps it grounded but with large whitespace below */
  }

  .instance-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .instance-name {
    font-size: 48px;
    font-weight: 700;
    letter-spacing: -1px;
  }

  .instance-meta {
    font-size: 14px;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .play-btn {
    font-size: 18px;
    padding: 16px 48px;
    border-radius: var(--border-radius);
    text-transform: uppercase;
    letter-spacing: 1px;
    font-weight: 700;
    box-shadow: 0 4px 14px rgba(59, 130, 246, 0.4); /* Slight glow for the primary action */
  }
</style>
