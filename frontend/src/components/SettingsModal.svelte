<script lang="ts">
  import { settings } from '../stores';
  import { fade, scale } from 'svelte/transition';
  
  export let show = false;
  
  let localSettings = { ...$settings };
  
  function saveSettings() {
    settings.set(localSettings);
    show = false;
  }
  
  function cancel() {
    localSettings = { ...$settings };
    show = false;
  }
  
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      cancel();
    }
  }
  
  function handleBackdropKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      cancel();
    }
  }
</script>

{#if show}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div 
    class="modal-backdrop" 
    on:click={handleBackdropClick} 
    on:keydown={handleBackdropKeydown}
    transition:fade={{ duration: 200 }}
    role="dialog"
    aria-modal="true"
    aria-labelledby="settings-title"
  >
    <div class="modal" transition:scale={{ duration: 200, start: 0.95 }}>
      <div class="modal-header">
        <h2 id="settings-title">Settings</h2>
        <button class="close-btn" on:click={cancel} aria-label="Close settings">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
      
      <div class="modal-content">
        <div class="settings-section">
          <h3>General</h3>
          
          <div class="setting-item">
            <label for="refresh-interval">Refresh Interval (seconds)</label>
            <input 
              id="refresh-interval"
              type="number" 
              min="1" 
              max="60" 
              bind:value={localSettings.refreshInterval}
              class="setting-input"
            />
          </div>
          
          <div class="setting-item">
            <label>
              <input 
                type="checkbox" 
                bind:checked={localSettings.soundEnabled}
              />
              Enable Sound Alerts
            </label>
          </div>
          
          <div class="setting-item">
            <label>
              <input 
                type="checkbox" 
                bind:checked={localSettings.desktopNotifications}
              />
              Enable Desktop Notifications
            </label>
          </div>
        </div>
        
        <div class="settings-section">
          <h3>Appearance</h3>
          
          <div class="setting-item">
            <label for="theme-select">Theme</label>
            <select 
              id="theme-select"
              bind:value={localSettings.theme}
              class="setting-select"
            >
              <option value="dark">Dark</option>
              <option value="light">Light</option>
            </select>
          </div>
          
          <div class="setting-item">
            <label>
              <input 
                type="checkbox" 
                bind:checked={localSettings.compactView}
              />
              Compact View
            </label>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="btn btn-secondary" on:click={cancel}>Cancel</button>
        <button class="btn btn-primary" on:click={saveSettings}>Save Changes</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  
  .modal {
    background: #1e2025;
    border: 1px solid #2a2d33;
    border-radius: 12px;
    width: 90%;
    max-width: 500px;
    max-height: 80vh;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  
  .modal-header {
    padding: 20px;
    border-bottom: 1px solid #2a2d33;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .modal-header h2 {
    font-size: 18px;
    font-weight: 600;
    color: #e0e6ed;
    margin: 0;
  }
  
  .close-btn {
    background: transparent;
    border: none;
    color: #8b92a9;
    cursor: pointer;
    padding: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    transition: all 0.2s;
  }
  
  .close-btn:hover {
    background: #2a2d33;
    color: #e0e6ed;
  }
  
  .modal-content {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
  }
  
  .settings-section {
    margin-bottom: 24px;
  }
  
  .settings-section h3 {
    font-size: 14px;
    font-weight: 600;
    color: #8b92a9;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 16px;
  }
  
  .setting-item {
    margin-bottom: 16px;
  }
  
  .setting-item label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #e0e6ed;
    cursor: pointer;
  }
  
  .setting-item input[type="checkbox"] {
    width: 18px;
    height: 18px;
    accent-color: #3b82f6;
    cursor: pointer;
  }
  
  .setting-input {
    width: 100%;
    padding: 8px 12px;
    background: #141619;
    border: 1px solid #2a2d33;
    border-radius: 6px;
    color: #e0e6ed;
    font-size: 14px;
    outline: none;
    transition: all 0.2s;
    margin-top: 8px;
  }
  
  .setting-input:focus {
    border-color: #3b82f6;
    background: #1a1d21;
  }
  
  .setting-select {
    width: 100%;
    padding: 8px 12px;
    background: #141619;
    border: 1px solid #2a2d33;
    border-radius: 6px;
    color: #e0e6ed;
    font-size: 14px;
    outline: none;
    transition: all 0.2s;
    margin-top: 8px;
    cursor: pointer;
  }
  
  .setting-select:focus {
    border-color: #3b82f6;
    background: #1a1d21;
  }
  
  .modal-footer {
    padding: 20px;
    border-top: 1px solid #2a2d33;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
  
  .btn {
    padding: 8px 20px;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .btn-primary {
    background: #3b82f6;
    color: white;
  }
  
  .btn-primary:hover {
    background: #2563eb;
  }
  
  .btn-secondary {
    background: transparent;
    color: #e0e6ed;
    border: 1px solid #2a2d33;
  }
  
  .btn-secondary:hover {
    background: #2a2d33;
  }
</style>