<script lang="ts">
  import { onMount } from 'svelte';
  import { GetActiveAccount } from '../../wailsjs/go/main/App';
  import { auth } from '../../wailsjs/go/models';
  import LoginModal from './LoginModal.svelte';

  let activeAccount: auth.Account | null = null;
  let showLoginModal = false;

  async function loadAccount() {
    try {
      activeAccount = await GetActiveAccount();
    } catch (e) {
      console.error("Failed to load active account", e);
    }
  }

  onMount(loadAccount);

  function handleLogin(event: CustomEvent<auth.Account>) {
    activeAccount = event.detail;
  }
</script>

<div class="account-manager">
  <div class="account-info">
    <div class="avatar">
      {#if activeAccount}
        <!-- Simple initial for offline account -->
        {activeAccount.username.charAt(0).toUpperCase()}
      {:else}
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </svg>
      {/if}
    </div>
    <div class="details">
      {#if activeAccount}
        <span class="username">{activeAccount.username}</span>
        <span class="status">Offline Account</span>
      {:else}
        <span class="username">Guest</span>
        <span class="status">Offline</span>
      {/if}
    </div>
  </div>
  <button class="switch-btn" on:click={() => showLoginModal = true}>
    {activeAccount ? 'Switch' : 'Login'}
  </button>
</div>

<LoginModal bind:showModal={showLoginModal} on:login={handleLogin} />

<style>
  .account-manager {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: var(--border-radius);
    margin-top: 16px;
  }

  .account-info {
    display: flex;
    align-items: center;
    gap: 12px;
    overflow: hidden; /* Added to support text truncation */
    flex: 1; /* Allow to take available space */
  }

  .avatar {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    background: var(--primary-color);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: 600;
    font-size: 16px;
  }

  .avatar svg {
    width: 20px;
    height: 20px;
  }

  .details {
    display: flex;
    flex-direction: column;
    gap: 2px;
    overflow: hidden; /* Support text truncation */
  }

  .username {
    font-size: 13px; /* Slightly smaller */
    font-weight: 600;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .status {
    font-size: 11px;
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .switch-btn {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.05);
    color: var(--text-primary);
    padding: 4px 8px; /* Smaller padding */
    border-radius: 4px;
    font-size: 11px; /* Smaller font */
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: all var(--transition-fast);
    flex-shrink: 0; /* Prevent button from shrinking */
  }

  .switch-btn:hover {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.1);
  }
</style>
