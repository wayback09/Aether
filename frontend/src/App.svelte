<script lang="ts">
  import TitleBar from './components/TitleBar.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import Home from './pages/Home.svelte';
  import Instances from './pages/Instances.svelte';
  import Extensions from './pages/Extensions.svelte';
  import ExtensionView from './pages/ExtensionView.svelte';

  let activePage = 'home';
  let extensionRoutes: Record<string, string> = {};

  function handleNavigate(event: CustomEvent<string>) {
    activePage = event.detail;
  }

  function handleRegisterExtensionRoute(event: CustomEvent<{id: string, url: string}>) {
    extensionRoutes[event.detail.id] = event.detail.url;
  }
</script>

<div class="app-container">
  <TitleBar />
  <div class="layout">
    <Sidebar {activePage} on:navigate={handleNavigate} on:registerExtensionRoute={handleRegisterExtensionRoute} />

  <main class="content">
    {#if activePage === 'home'}
      <Home />
    {:else if activePage === 'instances'}
      <Instances />
    {:else if activePage === 'extensions'}
      <Extensions />
    {:else if extensionRoutes[activePage]}
      <ExtensionView url={extensionRoutes[activePage]} />
    {:else}
      <div class="placeholder">
        <h2>{activePage.charAt(0).toUpperCase() + activePage.slice(1)}</h2>
        <p>This page is under construction.</p>
      </div>
    {/if}
  </main>
  </div>
</div>

<style>
  :global(body) {
    margin: 0;
    overflow: hidden;
  }

  .app-container {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    overflow: hidden;
    background: var(--bg-dark, #0d0d0d);
  }

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
