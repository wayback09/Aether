<script lang="ts">
  import { createEventDispatcher, onMount, onDestroy, tick } from 'svelte';
  import Icon from './Icon.svelte';

  export let open = false;

  interface Command {
    id: string;
    label: string;
    category?: string;
    action: () => void;
  }

  // Core built-in commands — extensions can push into this array later
  export let commands: Command[] = [];

  const dispatch = createEventDispatcher<{ close: void; navigate: string }>();

  let query = '';
  let selectedIndex = 0;
  let inputEl: HTMLInputElement;
  let listEl: HTMLDivElement;

  $: filtered = query.trim()
    ? commands.filter(
        (cmd) =>
          cmd.label.toLowerCase().includes(query.toLowerCase()) ||
          cmd.category?.toLowerCase().includes(query.toLowerCase())
      )
    : commands;

  // Reset selection on every filter change
  $: if (filtered) selectedIndex = 0;

  // Auto-focus when opened
  $: if (open) {
    tick().then(() => inputEl?.focus());
  }

  function close() {
    open = false;
    query = '';
    dispatch('close');
  }

  function select(cmd: Command) {
    close();
    cmd.action();
  }

  function onGlobalKeydown(e: KeyboardEvent) {
    if (!open) return;
    if (e.key === 'Escape') { e.preventDefault(); close(); return; }
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      selectedIndex = Math.min(selectedIndex + 1, filtered.length - 1);
      scrollToSelected();
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      selectedIndex = Math.max(selectedIndex - 1, 0);
      scrollToSelected();
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (filtered[selectedIndex]) select(filtered[selectedIndex]);
    }
  }

  function scrollToSelected() {
    tick().then(() => {
      const item = listEl?.querySelector('.palette-item.selected') as HTMLElement;
      item?.scrollIntoView({ block: 'nearest' });
    });
  }

  onMount(() => window.addEventListener('keydown', onGlobalKeydown));
  onDestroy(() => window.removeEventListener('keydown', onGlobalKeydown));
</script>

{#if open}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="palette-backdrop" on:click={close}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="palette" on:click|stopPropagation>
      <!-- Search Row -->
      <div class="palette-search">
        <Icon name="search" size={15} color="var(--text-secondary)" />
        <input
          bind:this={inputEl}
          bind:value={query}
          placeholder="Type a command..."
          class="palette-input"
          id="command-palette-input"
          autocomplete="off"
          spellcheck="false"
        />
        <kbd class="esc-badge">ESC</kbd>
      </div>

      <!-- Results -->
      <div class="palette-results" bind:this={listEl}>
        {#if filtered.length > 0}
          {#each filtered as cmd, i}
            <!-- svelte-ignore a11y-mouse-events-have-key-events -->
            <button
              class="palette-item {i === selectedIndex ? 'selected' : ''}"
              on:click={() => select(cmd)}
              on:mousemove={() => (selectedIndex = i)}
            >
              <span class="cmd-label">{cmd.label}</span>
              {#if cmd.category}
                <span class="cmd-category">{cmd.category}</span>
              {/if}
            </button>
          {/each}
        {:else}
          <div class="palette-empty">No commands match <strong>"{query}"</strong></div>
        {/if}
      </div>

      <!-- Footer hint -->
      <div class="palette-footer">
        <span><kbd>↑↓</kbd> navigate</span>
        <span><kbd>↵</kbd> run</span>
        <span><kbd>ESC</kbd> close</span>
      </div>
    </div>
  </div>
{/if}

<style>
  .palette-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.55);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    z-index: 9999;
    display: flex;
    justify-content: center;
    align-items: flex-start;
    padding-top: 14vh;
    /* Animate in */
    animation: fade-in 100ms ease;
  }

  @keyframes fade-in {
    from { opacity: 0; }
    to   { opacity: 1; }
  }

  .palette {
    width: 540px;
    max-width: 92vw;
    background: rgba(18, 18, 18, 0.92);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    box-shadow:
      0 24px 64px rgba(0, 0, 0, 0.8),
      0 0 0 1px rgba(255, 255, 255, 0.04);
    overflow: hidden;
    /* Slide-in */
    animation: slide-in 120ms cubic-bezier(0.16, 1, 0.3, 1);
  }

  @keyframes slide-in {
    from { transform: translateY(-8px); opacity: 0; }
    to   { transform: translateY(0);    opacity: 1; }
  }

  /* Search Row */
  .palette-search {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 14px 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  }

  .palette-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: var(--text-primary);
    font-family: inherit;
    font-size: 14px;
    caret-color: var(--accent-color);
  }

  .palette-input::placeholder {
    color: var(--text-secondary);
  }

  .esc-badge {
    font-family: inherit;
    font-size: 11px;
    font-weight: 500;
    padding: 2px 7px;
    border-radius: 5px;
    background: rgba(255, 255, 255, 0.07);
    color: var(--text-secondary);
    border: 1px solid rgba(255, 255, 255, 0.1);
    cursor: pointer;
  }

  /* Results List */
  .palette-results {
    max-height: 320px;
    overflow-y: auto;
    padding: 6px;
  }

  .palette-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: 9px 12px;
    border-radius: 7px;
    border: none;
    background: transparent;
    color: var(--text-primary);
    font-family: inherit;
    font-size: 13.5px;
    cursor: pointer;
    text-align: left;
    transition: background-color 80ms ease;
  }

  .palette-item.selected {
    background: rgba(59, 130, 246, 0.18);
    color: #fff;
  }

  .palette-item:not(.selected):hover {
    background: rgba(255, 255, 255, 0.05);
  }

  .cmd-label {
    font-weight: 500;
  }

  .cmd-category {
    font-size: 12px;
    color: var(--text-secondary);
    font-weight: 400;
  }

  .palette-empty {
    padding: 28px 16px;
    text-align: center;
    color: var(--text-secondary);
    font-size: 13px;
  }

  .palette-empty strong {
    color: var(--text-primary);
  }

  /* Footer */
  .palette-footer {
    display: flex;
    gap: var(--spacing-md);
    padding: 8px 16px;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
    font-size: 11px;
    color: var(--text-secondary);
  }

  .palette-footer kbd {
    font-family: inherit;
    font-size: 11px;
    padding: 1px 5px;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.07);
    border: 1px solid rgba(255, 255, 255, 0.1);
    margin-right: 3px;
  }
</style>
