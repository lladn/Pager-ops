<script lang="ts">
  import { services } from '../stores';
  import type { Service } from '../types';
  
  let collapsed = false;
  
  function togglePanel() {
    collapsed = !collapsed;
  }
  
  function toggleService(service: Service) {
    service.active = !service.active;
    services.update(s => s);
  }
</script>

<div class="services-panel" class:collapsed>
  <div class="panel-header">
    <span class="panel-title">Services</span>
    <button class="collapse-btn" on:click={togglePanel}>
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points={collapsed ? "9 18 15 12 9 6" : "15 18 9 12 15 6"}></polyline>
      </svg>
    </button>
  </div>
  <div class="services-list">
    {#each $services as service}
      <div 
        class="service-item" 
        class:active={service.active}
        class:inactive={!service.active}
        on:click={() => toggleService(service)}
      >
        <div class="service-info">
          <div class="service-icon">{service.abbreviation}</div>
          <span class="service-name">{service.name}</span>
        </div>
        {#if service.incidentCount > 0}
          <span class="service-count">{service.incidentCount}</span>
        {/if}
      </div>
    {/each}
  </div>
</div>

<style>
  .services-panel {
    width: 240px;
    background: #141619;
    border-right: 1px solid #2a2d33;
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    position: relative;
    display: flex;
    flex-direction: column;
  }

  .services-panel.collapsed {
    width: 60px;
  }

  .panel-header {
    padding: 16px;
    border-bottom: 1px solid #2a2d33;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .panel-title {
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    color: #8b92a9;
    letter-spacing: 0.5px;
  }

  .collapse-btn {
    width: 24px;
    height: 24px;
    background: transparent;
    border: none;
    color: #8b92a9;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    transition: all 0.2s;
  }

  .collapse-btn:hover {
    background: #1e2025;
    color: #e0e6ed;
  }

  .services-list {
    padding: 8px;
    overflow-y: auto;
    flex: 1;
  }

  .service-item {
    padding: 10px 12px;
    margin-bottom: 4px;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: relative;
  }

  .service-item.active {
    background: #1a1d21;
    border-left: 3px solid #3b82f6;
  }

  .service-item.inactive {
    opacity: 0.5;
    border-left: 3px solid #ef4444;
  }

  .service-item:hover {
    background: #1e2025;
  }

  .service-info {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
  }

  .service-icon {
    width: 20px;
    height: 20px;
    border-radius: 4px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 10px;
    font-weight: bold;
    color: white;
    flex-shrink: 0;
  }

  .service-name {
    font-size: 13px;
    color: #e0e6ed;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .service-count {
    background: #ef4444;
    color: white;
    font-size: 11px;
    font-weight: 600;
    padding: 2px 6px;
    border-radius: 10px;
    min-width: 20px;
    text-align: center;
  }

  .services-panel.collapsed .panel-header {
    padding: 16px 8px;
  }

  .services-panel.collapsed .panel-title {
    display: none;
  }

  .services-panel.collapsed .service-name {
    display: none;
  }

  .services-panel.collapsed .service-count {
    position: absolute;
    top: 2px;
    right: 2px;
    font-size: 9px;
    padding: 1px 4px;
  }
</style>