<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { ToggleService } from '../../wailsjs/go/main/App.js';
  
    export let services = [];
    
    const dispatch = createEventDispatcher();
    let searchQuery = '';
    let isUpdating = false;
  
    $: filteredServices = services.filter(service => 
      service.displayName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      service.id.toLowerCase().includes(searchQuery.toLowerCase())
    );
  
    $: enabledCount = services.filter(s => s.enabled).length;
    $: totalCount = services.length;
  
    async function handleToggleService(service) {
      if (isUpdating) return;
      
      try {
        isUpdating = true;
        await ToggleService(service.id, !service.enabled);
        // Update local state
        service.enabled = !service.enabled;
        services = services; // Trigger reactivity
      } catch (err) {
        console.error('Failed to toggle service:', err);
      } finally {
        isUpdating = false;
      }
    }
  
    function openSettings() {
      dispatch('openSettings');
    }
  
    function toggleAll(enable: boolean) {
      services.forEach(service => {
        if (service.enabled !== enable) {
          handleToggleService(service);
        }
      });
    }
  </script>
  
  <div class="service-panel">
    <div class="panel-header">
      <div class="header-top">
        <h2>Services</h2>
        <button class="settings-btn" on:click={openSettings} title="Settings">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M12 1v6m0 6v6m4.22-13.22l4.24 4.24M1.54 1.54l4.24 4.24M1 12h6m6 0h6m-13.22 4.22l4.24-4.24m8.44 8.44l-4.24-4.24"></path>
          </svg>
        </button>
      </div>
      <div class="service-stats">
        <span class="stat">{enabledCount} active</span>
        <span class="stat-separator">•</span>
        <span class="stat">{totalCount} total</span>
      </div>
      <div class="search-container">
        <input 
          type="text" 
          placeholder="Search services..." 
          bind:value={searchQuery}
          class="search-input"
        />
      </div>
      {#if totalCount > 0}
        <div class="bulk-actions">
          <button class="bulk-btn" on:click={() => toggleAll(true)}>Enable All</button>
          <button class="bulk-btn" on:click={() => toggleAll(false)}>Disable All</button>
        </div>
      {/if}
    </div>
    
    <div class="service-list">
      {#if filteredServices.length === 0}
        <div class="empty-state">
          {#if services.length === 0}
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="9" y1="9" x2="15" y2="9"></line>
              <line x1="9" y1="12" x2="15" y2="12"></line>
              <line x1="9" y1="15" x2="12" y2="15"></line>
            </svg>
            <p>No services configured</p>
            <p class="hint">Add services in Settings</p>
          {:else}
            <p>No services match your search</p>
          {/if}
        </div>
      {:else}
        {#each filteredServices as service (service.id)}
          <div class="service-item" class:enabled={service.enabled}>
            <button 
              class="service-toggle"
              class:active={service.enabled}
              on:click={() => handleToggleService(service)}
              disabled={isUpdating}
            >
              <div class="toggle-slider"></div>
            </button>
            <div class="service-info">
              <div class="service-name">{service.displayName}</div>
              <div class="service-id">{service.id}</div>
            </div>
            <div class="service-status" class:active={service.enabled}>
              {service.enabled ? 'Active' : 'Inactive'}
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>
  
  <style>
    .service-panel {
      height: 100%;
      display: flex;
      flex-direction: column;
    }
  
    .panel-header {
      padding: 1rem;
      border-bottom: 1px solid var(--border-color);
      background-color: var(--bg-tertiary);
    }
  
    .header-top {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }
  
    .panel-header h2 {
      margin: 0;
      font-size: 1.1rem;
      font-weight: 600;
    }
  
    .settings-btn {
      background: none;
      border: none;
      color: var(--text-secondary);
      cursor: pointer;
      padding: 0.25rem;
      border-radius: 4px;
      transition: all 0.2s;
    }
  
    .settings-btn:hover {
      background-color: var(--bg-secondary);
      color: var(--text-primary);
    }
  
    .service-stats {
      display: flex;
      gap: 0.5rem;
      font-size: 0.85rem;
      color: var(--text-secondary);
      margin-bottom: 0.75rem;
    }
  
    .stat-separator {
      color: var(--text-muted);
    }
  
    .search-container {
      margin-bottom: 0.75rem;
    }
  
    .search-input {
      width: 100%;
      padding: 0.5rem;
      background-color: var(--bg-primary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-primary);
      font-size: 0.9rem;
    }
  
    .search-input::placeholder {
      color: var(--text-muted);
    }
  
    .search-input:focus {
      outline: none;
      border-color: var(--accent-primary);
    }
  
    .bulk-actions {
      display: flex;
      gap: 0.5rem;
    }
  
    .bulk-btn {
      flex: 1;
      padding: 0.375rem 0.75rem;
      background-color: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-secondary);
      font-size: 0.8rem;
      cursor: pointer;
      transition: all 0.2s;
    }
  
    .bulk-btn:hover {
      background-color: var(--bg-primary);
      color: var(--text-primary);
    }
  
    .service-list {
      flex: 1;
      overflow-y: auto;
      padding: 0.5rem;
    }
  
    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: var(--text-muted);
      text-align: center;
      padding: 2rem;
    }
  
    .empty-state svg {
      margin-bottom: 1rem;
      opacity: 0.5;
    }
  
    .empty-state p {
      margin: 0.25rem 0;
    }
  
    .empty-state .hint {
      font-size: 0.85rem;
      color: var(--text-muted);
    }
  
    .service-item {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.75rem;
      background-color: var(--bg-primary);
      border: 1px solid var(--border-color);
      border-radius: 6px;
      margin-bottom: 0.5rem;
      transition: all 0.2s;
    }
  
    .service-item.enabled {
      border-color: var(--accent-primary);
      background-color: rgba(74, 158, 255, 0.05);
    }
  
    .service-toggle {
      position: relative;
      width: 44px;
      height: 24px;
      background-color: var(--bg-tertiary);
      border: 2px solid var(--border-color);
      border-radius: 12px;
      cursor: pointer;
      transition: all 0.3s;
      flex-shrink: 0;
    }
  
    .service-toggle:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  
    .service-toggle.active {
      background-color: var(--accent-primary);
      border-color: var(--accent-primary);
    }
  
    .toggle-slider {
      position: absolute;
      top: 2px;
      left: 2px;
      width: 16px;
      height: 16px;
      background-color: white;
      border-radius: 50%;
      transition: transform 0.3s;
    }
  
    .service-toggle.active .toggle-slider {
      transform: translateX(20px);
    }
  
    .service-info {
      flex: 1;
      min-width: 0;
    }
  
    .service-name {
      font-size: 0.9rem;
      font-weight: 500;
      color: var(--text-primary);
      margin-bottom: 0.125rem;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  
    .service-id {
      font-size: 0.75rem;
      color: var(--text-muted);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  
    .service-status {
      font-size: 0.75rem;
      color: var(--text-muted);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
  
    .service-status.active {
      color: var(--accent-success);
      font-weight: 600;
    }
  </style>