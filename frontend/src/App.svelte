<script lang="ts">
  import Sidebar from './components/Sidebar.svelte';
  import Home from './pages/Home.svelte';
  import Instances from './pages/Instances.svelte';
  import Extensions from './pages/Extensions.svelte';

  let activePage = 'home';

  function handleNavigate(event: CustomEvent<string>) {
    activePage = event.detail;
  }
</script>

<div class="layout">
  <Sidebar {activePage} on:navigate={handleNavigate} />

  <main class="content">
    {#if activePage === 'home'}
      <Home />
    {:else if activePage === 'instances'}
      <Instances />
    {:else if activePage === 'extensions'}
      <Extensions />
    {:else}
      <div class="placeholder">
        <h2>{activePage.charAt(0).toUpperCase() + activePage.slice(1)}</h2>
        <p>This page is under construction.</p>
      </div>
    {/if}
  </main>
</div>

<style>
  .layout {
    display: flex;
    width: 100%;
    height: 100%;
  }

  .content {
    flex-grow: 1;
    background-color: var(--bg-color);
    overflow: hidden; /* Each page handles its own scroll if needed */
  }

  .placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--text-secondary);
  }
  
  .placeholder h2 {
    color: var(--text-primary);
    margin-bottom: var(--spacing-sm);
  }
</style>
