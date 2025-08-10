<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    
    export let incidents = [];
    export let services = [];
    export let selectedId = null;
    
    const dispatch = createEventDispatcher();
    
    // Filter state
    let statusFilter = {
      triggered: true,
      acknowledged: true,
      resolved: false
    };
    let searchQuery = '';
    let groupByService = false;
    
    // Computed filtered incidents
    $: filteredIncidents = incidents.filter(incident => {
      // Status filter
      if (!statusFilter[incident.status]) return false;
      
      // Search filter
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        return incident.title.toLowerCase().includes(query) ||
               incident.description?.toLowerCase().includes(query) ||
               incident.id.toLowerCase().includes(query);
      }
      
      return true;
    });
    
    // Group incidents by service if enabled
    $: groupedIncidents = groupByService 
      ? groupIncidentsByService(filteredIncidents)
      : [{ service: null, incidents: filteredIncidents }];
    
    function groupIncidentsByService(incidents) {
      const groups = {};
      
      incidents.forEach(incident => {
        const serviceId = incident.service_id || 'unknown';
        if (!groups[serviceId]) {
          const service = services.find(s => s.id === serviceId);
          groups[serviceId] = {
            service: service || { id: serviceId, alias: 'Unknown Service' },
            incidents: []
          };
        }
        groups[serviceId].incidents.push(incident);
      });
      
      return Object.values(groups);
    }
    
    function selectIncident(incident) {
      dispatch('select', incident);
    }
    
    function handleAction(action, incident, event) {
      event.stopPropagation();
      dispatch('action', {
        action,
        incidentId: incident.id
      });
    }
    
    function getTimeAgo(date) {
      const now = new Date();
      const then = new Date(date);
      const diff = Math.floor((now - then) / 1000);
      
      if (diff < 60) return `${diff}s ago`;
      if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
      if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
      return `${Math.floor(diff / 86400)}d ago`;
    }
    
    function getEscalationTime(incident) {
      // Calculate time until escalation (mock implementation)
      // In real implementation, this would come from PagerDuty escalation policy
      if (incident.status === 'triggered') {
        return '5m';
      }
      return null;
    }
  </script>
  
  <div class="panel incidents-panel">
    <div class="panel-header">
      <div class="title">Incidents</div>
      <div class="filters">
        <div class="status-filters">
          <label class="filter-chip" class:active={statusFilter.triggered}>
            <input 
              type="checkbox" 
              bind:checked={statusFilter.triggered}
            />
            <span>Triggered</span>
          </label>
          <label class="filter-chip" class:active={statusFilter.acknowledged}>
            <input 
              type="checkbox" 
              bind:checked={statusFilter.acknowledged}
            />
            <span>Acknowledged</span>
          </label>
          <label class="filter-chip" class:active={statusFilter.resolved}>
            <input 
              type="checkbox" 
              bind:checked={statusFilter.resolved}
            />
            <span>Resolved</span>
          </label>
        </div>
        
        <div class="search-bar">
          <input 
            type="text" 
            placeholder="Search incidents..."
            bind:value={searchQuery}
            class="search-input"
          />
        </div>
        
        <label class="group-toggle">
          <input 
            type="checkbox" 
            bind:checked={groupByService}
          />
          <span>Group by Service</span>
        </label>
      </div>
    </div>
    
    <div class="panel-content">
      {#if filteredIncidents.length === 0}
        <div class="empty-state">
          <p>🎉 No incidents matching filters</p>
        </div>
      {:else}
        {#each groupedIncidents as group}
          {#if group.service && groupByService}
            <div class="service-group">
              <div class="service-group-header">
                {group.service.alias || group.service.id}
                <span class="incident-count">{group.incidents.length}</span>
              </div>
            </div>
          {/if}
          
          {#each group.incidents as incident}
            <div 
              class="incident-card"
              class:selected={selectedId === incident.id}
              class:pinned={incident.pinned}
              class:urgency-high={incident.urgency === 'high'}
              class:urgency-low={incident.urgency === 'low'}
              on:click={() => selectIncident(incident)}
            >
              {#if incident.pinned}
                <div class="pin-indicator">📌</div>
              {/if}
              
              <div class="incident-header">
                <div class="incident-number">#{incident.incident_number}</div>
                <div class="incident-status">
                  <span class="badge status-{incident.status}">
                    {incident.status}
                  </span>
                </div>
                {#if getEscalationTime(incident)}
                  <div class="escalation-warning" title="Time until escalation">
                    ⏰ {getEscalationTime(incident)}
                  </div>
                {/if}
              </div>
              
              <div class="incident-title">{incident.title}</div>
              
              {#if incident.description}
                <div class="incident-description">
                  {incident.description}
                </div>
              {/if}
              
              <div class="incident-meta">
                <span class="meta-item">
                  🕐 {getTimeAgo(incident.created_at)}
                </span>
                {#if incident.assigned_user_ids?.length > 0}
                  <span class="meta-item">
                    👤 {incident.assigned_user_ids.length} assigned
                  </span>
                {/if}
              </div>
              
              <div class="incident-actions">
                {#if incident.status === 'triggered'}
                  <button 
                    class="btn-action acknowledge"
                    on:click={(e) => handleAction('acknowledge', incident, e)}
                    title="Acknowledge"
                  >
                    ✓ Ack
                  </button>
                {/if}
                
                {#if incident.status !== 'resolved'}
                  <button 
                    class="btn-action resolve"
                    on:click={(e) => handleAction('resolve', incident, e)}
                    title="Resolve"
                  >
                    ✅ Resolve
                  </button>
                {/if}
                
                <button 
                  class="btn-action pin"
                  on:click={(e) => handleAction(incident.pinned ? 'unpin' : 'pin', incident, e)}
                  title={incident.pinned ? 'Unpin' : 'Pin'}
                >
                  {incident.pinned ? '📌' : '📍'}
                </button>
              </div>
            </div>
          {/each}
        {/each}
      {/if}
    </div>
  </div>
  
  <style>
    .incidents-panel {
      flex: 1;
      min-width: 400px;
    }
    
    .filters {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
      margin-top: 0.75rem;
    }
    
    .status-filters {
      display: flex;
      gap: 0.5rem;
    }
    
    .filter-chip {
      display: flex;
      align-items: center;
      gap: 0.25rem;
      padding: 0.25rem 0.75rem;
      background: var(--bg-primary);
      border-radius: 9999px;
      cursor: pointer;
      transition: all 0.2s;
      font-size: 0.875rem;
    }
    
    .filter-chip input {
      display: none;
    }
    
    .filter-chip.active {
      background: var(--accent-blue);
      color: white;
    }
    
    .search-bar {
      display: flex;
    }
    
    .search-input {
      flex: 1;
      padding: 0.5rem;
      background: var(--bg-primary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-primary);
      font-size: 0.875rem;
    }
    
    .search-input:focus {
      outline: none;
      border-color: var(--accent-blue);
    }
    
    .group-toggle {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      color: var(--text-secondary);
    }
    
    .empty-state {
      text-align: center;
      padding: 4rem 2rem;
      color: var(--text-tertiary);
      font-size: 1.125rem;
    }
    
    .service-group {
      margin-bottom: 1rem;
    }
    
    .service-group-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.5rem;
      background: var(--bg-tertiary);
      border-radius: 4px;
      font-weight: 600;
      color: var(--text-secondary);
      font-size: 0.875rem;
    }
    
    .incident-count {
      padding: 0.125rem 0.5rem;
      background: var(--bg-primary);
      border-radius: 9999px;
      font-size: 0.75rem;
    }
    
    .incident-card {
      position: relative;
      padding: 1rem;
      margin-bottom: 0.75rem;
      background: var(--bg-tertiary);
      border-radius: 8px;
      border: 2px solid transparent;
      cursor: pointer;
      transition: all 0.2s;
    }
    
    .incident-card:hover {
      background: var(--bg-primary);
      transform: translateX(4px);
    }
    
    .incident-card.selected {
      border-color: var(--accent-blue);
      background: rgba(66, 153, 225, 0.1);
    }
    
    .incident-card.pinned {
      border-color: var(--accent-yellow);
    }
    
    .incident-card.urgency-high {
      border-left: 4px solid var(--accent-red);
    }
    
    .incident-card.urgency-low {
      border-left: 4px solid var(--accent-green);
    }
    
    .pin-indicator {
      position: absolute;
      top: 0.5rem;
      right: 0.5rem;
      font-size: 1rem;
    }
    
    .incident-header {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-bottom: 0.5rem;
    }
    
    .incident-number {
      font-weight: 600;
      color: var(--text-secondary);
      font-size: 0.875rem;
    }
    
    .escalation-warning {
      margin-left: auto;
      padding: 0.125rem 0.5rem;
      background: rgba(252, 129, 129, 0.2);
      color: var(--accent-red);
      border-radius: 4px;
      font-size: 0.75rem;
      font-weight: 600;
      animation: pulse 2s infinite;
    }
    
    .incident-title {
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 0.5rem;
      word-break: break-word;
    }
    
    .incident-description {
      font-size: 0.875rem;
      color: var(--text-secondary);
      margin-bottom: 0.5rem;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .incident-meta {
      display: flex;
      gap: 1rem;
      margin-bottom: 0.75rem;
      font-size: 0.75rem;
      color: var(--text-tertiary);
    }
    
    .meta-item {
      display: flex;
      align-items: center;
      gap: 0.25rem;
    }
    
    .incident-actions {
      display: flex;
      gap: 0.5rem;
    }
    
    .btn-action {
      padding: 0.25rem 0.75rem;
      font-size: 0.75rem;
      border-radius: 4px;
      background: var(--bg-secondary);
      color: var(--text-secondary);
      transition: all 0.2s;
    }
    
    .btn-action:hover {
      transform: translateY(-1px);
    }
    
    .btn-action.acknowledge:hover {
      background: var(--accent-yellow);
      color: white;
    }
    
    .btn-action.resolve:hover {
      background: var(--accent-green);
      color: white;
    }
    
    .btn-action.pin:hover {
      background: var(--accent-blue);
      color: white;
    }
  </style><script lang="ts">
  import { incidents, resolvedIncidents, selectedIncident, searchQuery, activeTab, filterType } from '../stores';
  import IncidentCard from './IncidentCard.svelte';
  import type { Incident } from '../types';
  
  let filterMenuOpen = false;
  let detailsPanelCollapsed = false;
  
  export { detailsPanelCollapsed };
  
  function toggleFilterMenu() {
    filterMenuOpen = !filterMenuOpen;
  }
  
  function applyFilter(type: string) {
    filterType.set(type);
    filterMenuOpen = false;
    
    switch(type) {
      case 'groupByService':
        incidents.update(items => {
          return [...items].sort((a, b) => a.service.localeCompare(b.service));
        });
        break;
      case 'sortByTime':
        incidents.update(items => {
          return [...items].sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime());
        });
        break;
      case 'sortByAlerts':
        incidents.update(items => {
          return [...items].sort((a, b) => b.alertCount - a.alertCount);
        });
        break;
    }
  }
  
  function selectIncident(incident: Incident) {
    selectedIncident.set(incident);
    if (detailsPanelCollapsed) {
      detailsPanelCollapsed = false;
    }
  }
  
  function toggleDetailsPanel() {
    detailsPanelCollapsed = !detailsPanelCollapsed;
  }
  
  $: displayedIncidents = $activeTab === 'open' 
    ? $incidents.filter(i => i.title.toLowerCase().includes($searchQuery.toLowerCase()))
    : $resolvedIncidents.filter(i => i.title.toLowerCase().includes($searchQuery.toLowerCase()));
</script>

<div class="incidents-panel">
  <div class="incidents-header">
    <div class="incidents-tabs">
      <div class="tabs-container">
        <button 
          class="tab-btn" 
          class:active={$activeTab === 'open'}
          on:click={() => activeTab.set('open')}
        >
          Open Incidents
        </button>
        <button 
          class="tab-btn" 
          class:active={$activeTab === 'resolved'}
          on:click={() => activeTab.set('resolved')}
        >
          Resolved
        </button>
      </div>
      <div class="filter-dropdown">
        <button class="filter-btn" on:click={toggleFilterMenu}>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
          </svg>
        </button>
        {#if filterMenuOpen}
          <div class="filter-menu">
            <div class="filter-option" on:click={() => applyFilter('groupByService')}>Group by Service</div>
            <div class="filter-option" on:click={() => applyFilter('sortByTime')}>Sort by Time</div>
            <div class="filter-option" on:click={() => applyFilter('sortByAlerts')}>Sort by Alert Count</div>
          </div>
        {/if}
      </div>
    </div>
    <div class="incidents-search">
      <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.35-4.35"></path>
      </svg>
      <input 
        type="text" 
        class="search-input" 
        placeholder="Search incidents..."
        bind:value={$searchQuery}
      >
    </div>
  </div>
  <div class="incidents-list">
    {#each displayedIncidents as incident}
      <IncidentCard 
        {incident} 
        selected={$selectedIncident?.id === incident.id}
        on:select={() => selectIncident(incident)}
      />
    {/each}
  </div>
  {#if detailsPanelCollapsed}
    <button class="expand-details-btn" on:click={toggleDetailsPanel}>
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="9 18 15 12 9 6"></polyline>
      </svg>
    </button>
  {/if}
</div>

<style>
  .incidents-panel {
    flex: 1;
    min-width: 320px;
    background: #0f1114;
    display: flex;
    flex-direction: column;
    position: relative;
  }

  .incidents-header {
    padding: 16px;
    border-bottom: 1px solid #2a2d33;
    background: #141619;
  }

  .incidents-tabs {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 12px;
  }

  .tabs-container {
    display: flex;
    gap: 4px;
    background: #0f1114;
    padding: 2px;
    border-radius: 8px;
  }

  .tab-btn {
    padding: 6px 16px;
    background: transparent;
    border: none;
    color: #8b92a9;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.2s;
  }

  .tab-btn.active {
    background: #1e2025;
    color: #e0e6ed;
  }

  .filter-dropdown {
    position: relative;
  }

  .filter-btn {
    width: 32px;
    height: 32px;
    background: transparent;
    border: 1px solid #2a2d33;
    color: #8b92a9;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    transition: all 0.2s;
  }

  .filter-btn:hover {
    background: #1e2025;
    color: #e0e6ed;
    border-color: #3a3d43;
  }

  .filter-menu {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 4px;
    background: #1e2025;
    border: 1px solid #2a2d33;
    border-radius: 8px;
    padding: 4px;
    min-width: 180px;
    box-shadow: 0 10px 40px rgba(0,0,0,0.5);
    z-index: 50;
  }

  .filter-option {
    padding: 8px 12px;
    color: #e0e6ed;
    font-size: 13px;
    cursor: pointer;
    border-radius: 4px;
    transition: background 0.2s;
  }

  .filter-option:hover {
    background: #2a2d33;
  }

  .incidents-search {
    position: relative;
  }

  .search-input {
    width: 100%;
    padding: 8px 12px 8px 36px;
    background: #0f1114;
    border: 1px solid #2a2d33;
    border-radius: 8px;
    color: #e0e6ed;
    font-size: 13px;
    outline: none;
    transition: all 0.2s;
  }

  .search-input:focus {
    border-color: #3b82f6;
    background: #1a1d21;
  }

  .search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: #8b92a9;
    pointer-events: none;
  }

  .incidents-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }

  .expand-details-btn {
    position: absolute;
    left: -24px;
    top: 50%;
    transform: translateY(-50%);
    width: 24px;
    height: 48px;
    background: #141619;
    border: 1px solid #2a2d33;
    border-right: none;
    border-radius: 6px 0 0 6px;
    color: #8b92a9;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    z-index: 10;
  }

  .expand-details-btn:hover {
    background: #1e2025;
    color: #e0e6ed;
  }
</style>
