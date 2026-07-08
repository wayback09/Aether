<script lang="ts">
  import { onMount } from 'svelte';
  import { GetActiveInstance, LaunchInstance, InstallInstance } from '../../wailsjs/go/main/App.js';
  import { EventsOn } from '../../wailsjs/runtime/runtime.js';

  let currentInstance: any = null;
  let launchState = "Idle";
  let logs: string[] = [];
  let installProgress = 0;

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

    EventsOn("instance:progress", (data: any) => {
      if (currentInstance && data.id === currentInstance.id) {
        installProgress = data.progress;
        launchState = data.status;
        if (data.progress >= 100) {
           currentInstance.installed = true;
           launchState = "Idle";
        }
      }
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

  async function handleInstall() {
    if (!currentInstance) return;
    launchState = "Installing...";
    installProgress = 0;
    try {
      await InstallInstance(currentInstance.id);
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
      
      <div style="display: flex; gap: var(--spacing-md); align-items: center; width: 100%;">
        {#if currentInstance.installed}
          <button class="btn btn-primary play-btn" on:click={handlePlay} disabled={launchState === 'Running'}>
            {launchState === 'Running' ? 'Running' : 'Play'}
          </button>
        {:else}
          <div style="display: flex; flex-direction: column; gap: 8px; width: 100%; max-width: 400px;">
            <button class="btn btn-primary play-btn" on:click={handleInstall} disabled={launchState.includes('Download') || launchState === 'Installing...'}>
              {launchState === 'Idle' || launchState === 'Error' ? 'Install' : 'Installing...'}
            </button>
            {#if installProgress > 0 && installProgress < 100}
              <div style="width: 100%; height: 4px; background: rgba(255,255,255,0.1); border-radius: 2px; overflow: hidden;">
                <div style="height: 100%; background: var(--accent); width: {installProgress}%; transition: width 0.2s;"></div>
              </div>
            {/if}
          </div>
        {/if}
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
      <div style="text-align: center; color: var(--text-secondary); padding: 40px;">
        <h3>No Active Instance</h3>
        <p>Go to the Instances tab to create or select an instance to play.</p>
      </div>
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
