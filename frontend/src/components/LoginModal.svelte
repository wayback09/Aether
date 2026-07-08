<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { LoginOffline } from '../../wailsjs/go/main/App';

  export let showModal = false;

  const dispatch = createEventDispatcher();

  let username = '';
  let isLoading = false;
  let errorMsg = '';

  async function handleLogin() {
    if (!username.trim()) {
      errorMsg = 'Username cannot be empty';
      return;
    }
    
    isLoading = true;
    errorMsg = '';
    
    try {
      const account = await LoginOffline(username.trim());
      dispatch('login', account);
      showModal = false;
      username = '';
    } catch (err) {
      errorMsg = err.toString();
    } finally {
      isLoading = false;
    }
  }
</script>

{#if showModal}
  <div class="modal-backdrop" on:click={() => showModal = false} on:keydown={(e) => e.key === 'Escape' && (showModal = false)} role="button" tabindex="0">
    <div class="modal" on:click|stopPropagation on:keydown|stopPropagation role="button" tabindex="0">
      <h2>Offline Login</h2>
      <p class="subtitle">Enter a username to play offline.</p>
      
      <div class="form-group">
        <label for="username">Username</label>
        <input 
          id="username"
          type="text" 
          bind:value={username} 
          placeholder="e.g. Notch" 
          on:keydown={(e) => e.key === 'Enter' && handleLogin()}
          autofocus
        />
        {#if errorMsg}
          <span class="error">{errorMsg}</span>
        {/if}
      </div>

      <div class="modal-actions">
        <button class="btn secondary" on:click={() => showModal = false}>Cancel</button>
        <button class="btn primary" on:click={handleLogin} disabled={isLoading}>
          {isLoading ? 'Logging in...' : 'Login'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: rgba(30, 30, 30, 0.6);
    backdrop-filter: blur(24px);
    -webkit-backdrop-filter: blur(24px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    padding: 32px;
    border-radius: 16px;
    width: 400px;
    box-shadow: 0 24px 48px rgba(0, 0, 0, 0.4);
    color: var(--text-primary);
  }

  h2 {
    margin: 0 0 8px 0;
    font-size: 24px;
    font-weight: 600;
  }

  .subtitle {
    margin: 0 0 24px 0;
    color: var(--text-secondary);
    font-size: 14px;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 24px;
  }

  label {
    font-size: 13px;
    font-weight: 500;
    color: var(--text-secondary);
  }

  input {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: var(--text-primary);
    padding: 12px 16px;
    border-radius: var(--border-radius);
    font-size: 14px;
    font-family: inherit;
    transition: border-color var(--transition-fast);
  }

  input:focus {
    outline: none;
    border-color: var(--primary-color);
  }

  .error {
    color: #ff5555;
    font-size: 12px;
    margin-top: 4px;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }

  .btn {
    padding: 10px 20px;
    border-radius: var(--border-radius);
    font-size: 14px;
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    border: none;
    transition: all var(--transition-fast);
  }

  .btn.secondary {
    background: transparent;
    color: var(--text-secondary);
  }

  .btn.secondary:hover {
    background: rgba(255, 255, 255, 0.05);
    color: var(--text-primary);
  }

  .btn.primary {
    background: var(--primary-color);
    color: white;
  }

  .btn.primary:hover {
    filter: brightness(1.1);
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
