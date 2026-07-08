<script lang="ts">
  import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime.js';

  function minimize() {
    WindowMinimise();
  }
  function maximize() {
    WindowToggleMaximise();
  }
  function closeApp() {
    Quit();
  }
</script>

<div class="titlebar" style="--wails-draggable: drag">
  <div class="title">
    <div class="logo">W</div>
    <span>Aether</span>
  </div>
  
  <div class="controls" style="--wails-draggable: no-drag">
    <button class="win-btn" on:click={minimize}>
      <svg width="10" height="10" viewBox="0 0 10 10" stroke="currentColor" stroke-width="1">
        <line x1="0" y1="5" x2="10" y2="5"></line>
      </svg>
    </button>
    <button class="win-btn" on:click={maximize}>
      <svg width="10" height="10" viewBox="0 0 10 10" stroke="currentColor" stroke-width="1" fill="none">
        <rect x="1.5" y="1.5" width="7" height="7"></rect>
      </svg>
    </button>
    <button class="win-btn close-btn" on:click={closeApp}>
      <svg width="10" height="10" viewBox="0 0 10 10" stroke="currentColor" stroke-width="1">
        <line x1="1" y1="1" x2="9" y2="9"></line>
        <line x1="9" y1="1" x2="1" y2="9"></line>
      </svg>
    </button>
  </div>
</div>

<style>
  .titlebar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 32px;
    background: #0d0d0d;
    user-select: none;
    -webkit-user-select: none;
    flex-shrink: 0;
    border-bottom: 1px solid rgba(255,255,255,0.05);
  }

  .title {
    display: flex;
    align-items: center;
    padding-left: 12px;
    font-size: 12px;
    color: #a0a0a0;
    gap: 8px;
    font-family: inherit;
  }
  
  .logo {
    background: #ffffff;
    color: #000000;
    font-size: 10px;
    font-weight: 800;
    padding: 2px 4px;
    border-radius: 2px;
    line-height: 1;
  }

  .controls {
    display: flex;
    height: 100%;
  }

  .win-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 46px;
    height: 100%;
    background: transparent;
    border: none;
    color: #a0a0a0;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
    outline: none;
  }

  .win-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #ffffff;
  }

  .close-btn:hover {
    background: #e81123;
    color: #ffffff;
  }
</style>
