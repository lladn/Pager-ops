<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    
    export let services = [];
    
    const dispatch = createEventDispatcher();
    let collapsed = false;
    let hovering = false;
    
    function toggleService(service) {
      dispatch('toggle', {
        id: service.id,
        enabled: !service.enabled
      });
    }
    
    function toggleAll(enabled) {
      services.forEach(service => {
        if (service.enabled !== enabled) {
          dispatch('toggle', {
            id: service.id,
            enabled: enabled
          });
        }
      });
    }
    
    function handleMouseEnter() {
      if (collapsed) {
        hovering = true;
      }
    }
    
    function handleMouseLeave() {
      hovering = false;
    }
    
    $: panelClass = collapsed && !hovering ? 'panel services-panel collapsed' : 'panel services-panel';
    $: enabledCount = services.filter(s => s.enabled).length;
  </script>
  
  <div 
    class={panelClass}
    on:mouseenter={handleMouseEnter}
    on:mouseleave={handleMouseLeave}
  >
    <div class="panel-header">
      <div class="header-content">
        <span class="title">Services ({enabledCount}/{services.length})</span>
        <button 
          class="collapse-btn"
          on:click={() => collapsed = !collapsed}
          title={collapsed ? 'Expand' : 'Collapse'}
        >
          {collapsed ? '→' : '←'}
        </button>
      </div>
      {#if !collapsed || hovering}
        <div class="actions">
          <button class="btn-sm" on:click={() => toggleAll(true)}>All On</button>
          <button class="btn-sm" on:click={() => toggleAll(false)}>All Off</button>
        </div>
      {/if}
    </div>
    
    {#if !collapsed || hovering}
      <div class="panel-content">
        {#if services.length === 0}
          <div class="empty-state">
            <p>No services configured</p>
            <p class="hint">Add services in Settings</p>
          </div>
        {:else}
          <div class="services-list">
            {#each services as service}
              <div 
                class="service-item"
                class:enabled={service.enabled}
                on:click={() => toggleService(service)}
              >
                <div class="service-toggle">
                  <span class="toggle-indicator" class:active={service.enabled}>
                    {service.enabled ? '✓' : ''}
                  </span>
                </div>
                <div class="service-info">
                  <div class="service-name">{service.alias || service.id}</div>
                  {#if service.alias && service.alias !== service.id}
                    <div class="service-id">{service.id}</div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    {/if}
  </div>
  
  <style>
    .services-panel {
      width: 250px;
      transition: width 0.3s ease;
    }
    
    .services-panel.collapsed {
      width: 60px;
    }
    
    .header-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .title {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    
    .collapse-btn {
      padding: 0.25rem;
      background: transparent;
      color: var(--text-secondary);
      font-size: 1rem;
      min-width: 24px;
    }
    
    .collapse-btn:hover {
      background: var(--bg-primary);
      transform: none;
    }
    
    .actions {
      display: flex;
      gap: 0.5rem;
      margin-top: 0.5rem;
    }
    
    .btn-sm {
      padding: 0.25rem 0.5rem;
      font-size: 0.75rem;
      background: var(--bg-primary);
      color: var(--text-secondary);
      flex: 1;
    }
    
    .btn-sm:hover {
      background: var(--accent-blue);
      color: white;
    }
    
    .empty-state {
      text-align: center;
      padding: 2rem 1rem;
      color: var(--text-tertiary);
    }
    
    .empty-state p {
      margin: 0.5rem 0;
    }
    
    .hint {
      font-size: 0.875rem;
      opacity: 0.7;
    }
    
    .services-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .service-item {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.75rem;
      background: var(--bg-tertiary);
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.2s;
      border: 2px solid transparent;
    }
    
    .service-item:hover {
      background: var(--bg-primary);
      transform: translateX(4px);
    }
    
    .service-item.enabled {
      border-color: var(--accent-green);
      background: rgba(72, 187, 120, 0.1);
    }
    
    .service-toggle {
      width: 24px;
      height: 24px;
      border: 2px solid var(--border-color);
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.2s;
    }
    
    .service-item.enabled .service-toggle {
      background: var(--accent-green);
      border-color: var(--accent-green);
    }
    
    .toggle-indicator {
      color: white;
      font-weight: bold;
      font-size: 0.875rem;
    }
    
    .service-info {
      flex: 1;
      min-width: 0;
    }
    
    .service-name {
      font-weight: 600;
      color: var(--text-primary);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    
    .service-id {
      font-size: 0.75rem;
      color: var(--text-tertiary);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  </style>