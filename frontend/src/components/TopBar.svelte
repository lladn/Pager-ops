<script lang="ts">
  import { statusCounts } from '../stores';
  import SettingsModal from './SettingsModal.svelte';
  
  let showSettings = false;
  
  function openSettings() {
    showSettings = true;
  }
</script>

<div class="top-bar">
  <div class="app-title">PagerOps</div>
  <div class="top-bar-right">
    <div class="status-indicators">
      <div class="status-item">
        <span class="status-dot triggered"></span>
        <span>Tri: {$statusCounts.triggered}</span>
      </div>
      <div class="status-item">
        <span class="status-dot acknowledged"></span>
        <span>Ack: {$statusCounts.acknowledged}</span>
      </div>
      <div class="status-item">
        <span class="status-dot resolved"></span>
        <span>Res: {$statusCounts.resolved}</span>
      </div>
    </div>
    <button class="settings-btn" on:click={openSettings}>
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="3"></circle>
        <path d="M12 1v6m0 6v6m4.22-13.22l4.24 4.24M1.54 1.54l4.24 4.24M21 12h-6m-6 0H3m13.22 4.22l4.24 4.24M1.54 21.54l4.24-4.24"></path>
      </svg>
    </button>
  </div>
</div>

{#if showSettings}
  <SettingsModal bind:show={showSettings} />
{/if}

<style>
  .top-bar {
    height: 48px;
    background: #141619;
    border-bottom: 1px solid #2a2d33;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    position: relative;
    z-index: 100;
  }

  .app-title {
    font-size: 16px;
    font-weight: 600;
    color: #e0e6ed;
    letter-spacing: -0.5px;
  }

  .top-bar-right {
    display: flex;
    align-items: center;
    gap: 24px;
  }

  .status-indicators {
    display: flex;
    gap: 16px;
    font-size: 13px;
    font-weight: 500;
  }

  .status-item {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }

  .status-dot.triggered { background: #ef4444; }
  .status-dot.acknowledged { background: #eab308; }
  .status-dot.resolved { background: #22c55e; }

  .settings-btn {
    width: 32px;
    height: 32px;
    background: transparent;
    border: none;
    color: #8b92a9;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    transition: all 0.2s;
  }

  .settings-btn:hover {
    background: #1e2025;
    color: #e0e6ed;
  }
</style>