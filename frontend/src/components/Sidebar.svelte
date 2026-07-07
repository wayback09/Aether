<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
  export let activePage: string = 'home';
  
  const dispatch = createEventDispatcher();
  
  const topNav = [
    { id: 'home', label: 'Home' },
    { id: 'instances', label: 'Instances' },
    { id: 'extensions', label: 'Extensions' }
  ];
  
  const bottomNav = [
    { id: 'settings', label: 'Settings' },
    { id: 'account', label: 'Account' }
  ];
  
  function navigate(pageId: string) {
    dispatch('navigate', pageId);
  }
</script>

<aside class="sidebar">
  <div class="logo">Aether</div>
  
  <nav class="top-nav">
    {#each topNav as item}
      <button 
        class="nav-item {activePage === item.id ? 'active' : ''}" 
        on:click={() => navigate(item.id)}
      >
        {item.label}
      </button>
    {/each}
  </nav>
  
  <nav class="bottom-nav">
    {#each bottomNav as item}
      <button 
        class="nav-item {activePage === item.id ? 'active' : ''}" 
        on:click={() => navigate(item.id)}
      >
        {item.label}
      </button>
    {/each}
  </nav>
</aside>

<style>
  .sidebar {
    width: 240px;
    min-width: 240px;
    background-color: var(--sidebar-bg);
    display: flex;
    flex-direction: column;
    padding: 24px 16px;
    box-sizing: border-box;
    border-right: 1px solid rgba(255, 255, 255, 0.05);
  }

  .logo {
    font-size: 24px;
    font-weight: 700;
    margin-bottom: 40px;
    padding: 0 16px;
    letter-spacing: -0.5px;
  }

  .top-nav {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .bottom-nav {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-top: auto;
  }

  .nav-item {
    background: transparent;
    border: none;
    text-align: left;
    padding: 10px 16px;
    border-radius: var(--border-radius);
    color: var(--text-secondary);
    font-family: inherit;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color var(--transition-fast), color var(--transition-fast);
  }

  .nav-item:hover {
    background-color: rgba(255, 255, 255, 0.05);
    color: var(--text-primary);
  }

  .nav-item.active {
    background-color: rgba(255, 255, 255, 0.1);
    color: var(--text-primary);
    font-weight: 600;
  }
</style>
