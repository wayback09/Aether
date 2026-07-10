<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import Icon from './Icon.svelte';

  /** Icon name from Icon.svelte — shown large and faint above the title */
  export let icon: string = 'package';
  export let title: string = 'Nothing here yet';
  export let description: string = '';
  /** Label for the CTA button. Leave empty to hide the button. */
  export let actionLabel: string = '';

  const dispatch = createEventDispatcher<{ action: void }>();
</script>

<div class="empty-state">
  <div class="empty-icon">
    <Icon name={icon} size={72} color="rgba(255,255,255,0.07)" />
  </div>
  <h3 class="empty-title">{title}</h3>
  {#if description}
    <p class="empty-description">{description}</p>
  {/if}
  {#if actionLabel}
    <button class="btn btn-primary empty-action" on:click={() => dispatch('action')}>
      {actionLabel}
    </button>
  {/if}
</div>

<style>
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 50vh;
    gap: var(--spacing-sm);
    text-align: center;
    padding: var(--spacing-xl);
    box-sizing: border-box;
  }

  .empty-icon {
    margin-bottom: var(--spacing-md);
    /* Icon already uses a very low opacity colour, this adds a subtle scale */
    transform: scale(1);
  }

  .empty-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .empty-description {
    font-size: 14px;
    color: var(--text-secondary);
    margin: 0;
    max-width: 280px;
    line-height: 1.6;
  }

  .empty-action {
    margin-top: var(--spacing-md);
  }
</style>
