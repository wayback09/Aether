<script lang="ts">
  export let url: string = '';
</script>

<div class="extension-view">
  {#if url}
    <iframe src={url} title="Extension View" sandbox="allow-scripts allow-same-origin"></iframe>
  {:else}
    <div class="placeholder">
      <h2>Extension Loading...</h2>
    </div>
  {/if}
</div>

<style>
  .extension-view {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background: var(--bg-color);
  }

  iframe {
    flex-grow: 1;
    width: 100%;
    height: 100%;
    border: none;
    background: transparent;
  }

  .placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--text-secondary);
  }
</style>
