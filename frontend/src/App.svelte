<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import TitleBar from './components/TitleBar.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import CommandPalette from './components/CommandPalette.svelte';
  import ToastContainer from './components/ToastContainer.svelte';
  import Home from './pages/Home.svelte';
  import Instances from './pages/Instances.svelte';
  import Extensions from './pages/Extensions.svelte';
  import ExtensionView from './pages/ExtensionView.svelte';
  import InstanceDetails from './pages/InstanceDetails.svelte';
  import Settings from './pages/Settings.svelte';

  let activePage = 'home';
  let targetInstanceId = '';
  let activeInstanceId = '';
  let extensionRoutes: Record<string, { url: string; extensionId: string }> = {};
  let paletteOpen = false;

  // ── Navigation ──────────────────────────────────────────────────────────────

  function navigate(page: string) {
    activePage = page;
  }

  function handleNavigate(event: CustomEvent<string>) {
    const payload = event.detail;
    if (payload.startsWith('instance-details:')) {
      targetInstanceId = payload.split(':')[1];
      navigate('instance-details');
    } else if (payload.startsWith('home:instance:')) {
      // Navigate home and pre-select the given instance ID
      activeInstanceId = payload.split(':')[2];
      navigate('home');
    } else {
      navigate(payload);
    }
  }

  function handleRegisterExtensionRoute(event: CustomEvent<{ id: string; url: string; extensionId: string }>) {
    extensionRoutes[event.detail.id] = { url: event.detail.url, extensionId: event.detail.extensionId };
  }

  // ── Command Palette ──────────────────────────────────────────────────────────

  /** Core commands — every extension can push more in later */
  const commands = [
    { id: 'go-home',       label: 'Go to Home',        category: 'Navigation', action: () => navigate('home')       },
    { id: 'go-instances',  label: 'Go to Instances',   category: 'Navigation', action: () => navigate('instances')  },
    { id: 'go-extensions', label: 'Go to Extensions',  category: 'Navigation', action: () => navigate('extensions') },
    { id: 'go-settings',   label: 'Go to Settings',    category: 'Navigation', action: () => navigate('settings')   },
    { id: 'open-gallery',  label: 'Open Gallery',       category: 'Extensions', action: () => navigate('extensions') },
    {
      id: 'create-instance',
      label: 'Create Instance',
      category: 'Instances',
      action: () => {
        navigate('instances');
        // Slight delay so the page mounts before opening the modal
        setTimeout(() => {
          window.dispatchEvent(new CustomEvent('aether:open-create-instance'));
        }, 80);
      },
    },
  ];

  function openPalette() {
    paletteOpen = true;
  }

  function onGlobalKeydown(e: KeyboardEvent) {
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'P') {
      e.preventDefault();
      openPalette();
    }
  }

  onMount(() => window.addEventListener('keydown', onGlobalKeydown));
  onDestroy(() => window.removeEventListener('keydown', onGlobalKeydown));
</script>

<div class="app-container">
  <TitleBar />
  <div class="layout">
    <Sidebar
      {activePage}
      on:navigate={handleNavigate}
      on:registerExtensionRoute={handleRegisterExtensionRoute}
    />

    <main class="content">
      {#if activePage === 'home'}
        <Home on:navigate={handleNavigate} activeInstanceId={activeInstanceId} />
      {:else if activePage === 'instances'}
        <Instances on:navigate={handleNavigate} />
      {:else if activePage === 'instance-details'}
        <InstanceDetails instanceId={targetInstanceId} on:navigate={handleNavigate} />
      {:else if activePage === 'extensions'}
        <Extensions />
      {:else if activePage === 'settings'}
        <Settings />
      {:else if extensionRoutes[activePage]}
        <ExtensionView url={extensionRoutes[activePage].url} extID={extensionRoutes[activePage].extensionId} />
      {:else}
        <div class="placeholder">
          <h2>{activePage.charAt(0).toUpperCase() + activePage.slice(1)}</h2>
          <p>This page is under construction.</p>
        </div>
      {/if}
    </main>
  </div>

  <!-- Command Palette — rendered above everything -->
  <CommandPalette bind:open={paletteOpen} {commands} on:close={() => (paletteOpen = false)} />

  <!-- Global Toasts -->
  <ToastContainer />
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
