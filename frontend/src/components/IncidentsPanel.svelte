<script lang="ts">
  import { incidents, resolvedIncidents, services, activeTab, searchQuery, selectedIncident } from '../stores';
  import { acknowledgeIncident, resolveIncident, pinIncident } from '../services/api';
  import type { Incident, Service } from '../types';
  
  interface GroupedIncidents {
    service: Service | null;
    incidents: Incident[];
  }
  
  // Filter state
  let statusFilter = {
    triggered: true,
    acknowledged: true,
    resolved: false
  };
  let groupByService = false;
  let isLoading = false;
  
  // Get the correct incidents based on active tab
  $: currentIncidents = $activeTab === 'open' ? $incidents : $resolvedIncidents;
  
  // Computed filtered incidents
  $: filteredIncidents = currentIncidents.filter((incident: Incident) => {
    // Status filter
    if (!statusFilter[incident.status]) return false;
    
    // Search filter
    if ($searchQuery) {
      const query = $searchQuery.toLowerCase();
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
  
  function groupIncidentsByService(incidentsList: Incident[]): GroupedIncidents[] {
    const groups: Record<string, GroupedIncidents> = {};
    
    incidentsList.forEach(incident => {
      const serviceName = incident.service || 'unknown';
      if (!groups[serviceName]) {
        const service = $services.find(s => s.name === serviceName) || null;
        groups[serviceName] = {
          service,
          incidents: []
        };
      }
      groups[serviceName].incidents.push(incident);
    });
    
    return Object.values(groups);
  }
  
  function selectIncident(incident: Incident) {
    selectedIncident.set(incident);
  }
  
  function handleIncidentKeydown(event: KeyboardEvent, incident: Incident) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      selectIncident(incident);
    }
  }
  
  async function handleAction(action: string, incident: Incident, event: MouseEvent) {
    event.stopPropagation();
    
    if (isLoading) return;
    isLoading = true;
    
    try {
      switch (action) {
        case 'acknowledge':
          await acknowledgeIncident(incident.id);
          break;
        case 'resolve':
          await resolveIncident(incident.id);
          break;
        case 'pin':
          await pinIncident(incident.id, true);
          break;
        case 'unpin':
          await pinIncident(incident.id, false);
          break;
      }
    } catch (error) {
      console.error(`Failed to ${action} incident:`, error);
    } finally {
      isLoading = false;
    }
  }
  
  function getTimeAgo(date: Date | string): string {
    const now = new Date();
    const then = new Date(date);
    const diffMs = now.getTime() - then.getTime();
    const diff = Math.floor(diffMs / 1000);
    
    if (diff < 60) return `${diff}s ago`;
    if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
    if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
    return `${Math.floor(diff / 86400)}d ago`;
  }
  
  function getEscalationTime(incident: Incident): string | null {
    // Calculate time until escalation based on escalatesIn field
    if (incident.status === 'triggered' && incident.escalatesIn) {
      return `${incident.escalatesIn}m`;
    }
    return null;
  }
  
  function toggleStatusFilter(status: 'triggered' | 'acknowledged' | 'resolved') {
    statusFilter[status] = !statusFilter[status];
  }
  
  function toggleGroupByService() {
    groupByService = !groupByService;
  }
</script>

<div class="incidents-panel">
  <div class="incidents-header">
    <div class="tabs">
      <button 
        class="tab" 
        class:active={$activeTab === 'open'}
        on:click={() => activeTab.set('open')}
      >
        Open ({$incidents.length})
      </button>
      <button 
        class="tab" 
        class:active={$activeTab === 'resolved'}
        on:click={() => activeTab.set('resolved')}
      >
        Resolved ({$resolvedIncidents.length})
      </button>
    </div>
    
    <div class="filters">
      <input 
        type="text" 
        placeholder="Search incidents..."
        bind:value={$searchQuery}
        class="search-input"
      />
      
      <div class="filter-buttons">
        <button 
          class="filter-btn" 
          class:active={statusFilter.triggered}
          on:click={() => toggleStatusFilter('triggered')}
          title="Show/hide triggered incidents"
        >
          Triggered
        </button>
        <button 
          class="filter-btn" 
          class:active={statusFilter.acknowledged}
          on:click={() => toggleStatusFilter('acknowledged')}
          title="Show/hide acknowledged incidents"
        >
          Acknowledged
        </button>
        {#if $activeTab === 'resolved'}
          <button 
            class="filter-btn" 
            class:active={statusFilter.resolved}
            on:click={() => toggleStatusFilter('resolved')}
            title="Show/hide resolved incidents"
          >
            Resolved
          </button>
        {/if}
        <button 
          class="filter-btn" 
          class:active={groupByService}
          on:click={toggleGroupByService}
          title="Group by service"
        >
          Group by Service
        </button>
      </div>
    </div>
  </div>
  
  <div class="incidents-list">
    {#if filteredIncidents.length === 0}
      <div class="empty-state">
        <p>🎉 No incidents matching filters</p>
      </div>
    {:else}
      {#each groupedIncidents as group}
        {#if group.service && groupByService}
          <div class="service-group">
            <div class="service-group-header">
              {group.service.name}
              <span class="incident-count">{group.incidents.length}</span>
            </div>
          </div>
        {/if}
        
        {#each group.incidents as incident}
          <div 
            class="incident-card"
            class:selected={$selectedIncident?.id === incident.id}
            class:pinned={incident.pinned}
            class:urgency-high={incident.urgency === 'high'}
            class:urgency-low={incident.urgency === 'low'}
            on:click={() => selectIncident(incident)}
            on:keydown={(e) => handleIncidentKeydown(e, incident)}
            role="button"
            tabindex="0"
            aria-label="{incident.title} - {incident.status}"
          >
            {#if incident.pinned}
              <div class="pin-indicator" aria-label="Pinned">📌</div>
            {/if}
            
            <div class="incident-header">
              <div class="incident-service">{incident.service}</div>
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
                🕐 {getTimeAgo(incident.createdAt)}
              </span>
              {#if incident.alertCount > 0}
                <span class="meta-item">
                  ⚠️ {incident.alertCount} alerts
                </span>
              {/if}
            </div>
            
            <div class="incident-actions">
              {#if incident.status === 'triggered'}
                <button 
                  class="btn-action acknowledge"
                  on:click={(e) => handleAction('acknowledge', incident, e)}
                  disabled={isLoading}
                  title="Acknowledge"
                >
                  ✓ Ack
                </button>
              {/if}
              
              {#if incident.status !== 'resolved'}
                <button 
                  class="btn-action resolve"
                  on:click={(e) => handleAction('resolve', incident, e)}
                  disabled={isLoading}
                  title="Resolve"
                >
                  ✅ Resolve
                </button>
              {/if}
              
              <button 
                class="btn-action pin"
                on:click={(e) => handleAction(incident.pinned ? 'unpin' : 'pin', incident, e)}
                disabled={isLoading}
                title={incident.pinned ? 'Unpin' : 'Pin'}
              >
                {incident.pinned ? '📌 Unpin' : '📍 Pin'}
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
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-primary);
  }
  
  .incidents-header {
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
    background: var(--bg-secondary);
  }
  
  .tabs {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
  }
  
  .tab {
    padding: 8px 16px;
    background: transparent;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .tab:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }
  
  .tab.active {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    border-color: var(--status-acknowledged);
  }
  
  .filters {
    display: flex;
    gap: 12px;
    align-items: center;
  }
  
  .search-input {
    flex: 1;
    padding: 8px 12px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    color: var(--text-primary);
    outline: none;
  }
  
  .search-input:focus {
    border-color: var(--status-acknowledged);
  }
  
  .filter-buttons {
    display: flex;
    gap: 8px;
  }
  
  .filter-btn {
    padding: 6px 12px;
    background: transparent;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    color: var(--text-secondary);
    cursor: pointer;
    font-size: 12px;
    transition: all 0.2s;
  }
  
  .filter-btn:hover {
    background: var(--bg-hover);
  }
  
  .filter-btn.active {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    border-color: var(--status-acknowledged);
  }
  
  .incidents-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }
  
  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: var(--text-muted);
  }
  
  .service-group {
    margin-bottom: 16px;
  }
  
  .service-group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: var(--bg-secondary);
    border-radius: 4px;
    margin-bottom: 8px;
    font-weight: 600;
    color: var(--text-primary);
  }
  
  .incident-count {
    background: var(--bg-tertiary);
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 12px;
  }
  
  .incident-card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 12px;
    cursor: pointer;
    transition: all 0.2s;
    position: relative;
  }
  
  .incident-card:hover {
    border-color: var(--status-acknowledged);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .incident-card:focus {
    outline: 2px solid var(--status-acknowledged);
    outline-offset: 2px;
  }
  
  .incident-card.selected {
    border-color: var(--status-acknowledged);
    background: var(--bg-tertiary);
  }
  
  .incident-card.pinned {
    border-left: 3px solid var(--status-acknowledged);
  }
  
  .incident-card.urgency-high {
    border-left: 3px solid var(--urgency-high);
  }
  
  .incident-card.urgency-low {
    border-left: 3px solid var(--urgency-low);
  }
  
  .pin-indicator {
    position: absolute;
    top: 8px;
    right: 8px;
    font-size: 16px;
  }
  
  .incident-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
  }
  
  .incident-service {
    font-size: 12px;
    color: var(--text-secondary);
  }
  
  .badge {
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
  }
  
  .status-triggered {
    background: var(--status-triggered);
    color: white;
  }
  
  .status-acknowledged {
    background: var(--status-acknowledged);
    color: white;
  }
  
  .status-resolved {
    background: var(--status-resolved);
    color: white;
  }
  
  .escalation-warning {
    background: var(--status-escalated);
    color: white;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
  }
  
  .incident-title {
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
  }
  
  .incident-description {
    color: var(--text-secondary);
    font-size: 14px;
    margin-bottom: 8px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .incident-meta {
    display: flex;
    gap: 16px;
    color: var(--text-muted);
    font-size: 12px;
    margin-bottom: 12px;
  }
  
  .meta-item {
    display: flex;
    align-items: center;
    gap: 4px;
  }
  
  .incident-actions {
    display: flex;
    gap: 8px;
  }
  
  .btn-action {
    padding: 4px 12px;
    border: none;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .btn-action:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .btn-action.acknowledge {
    background: var(--status-acknowledged);
    color: white;
  }
  
  .btn-action.resolve {
    background: var(--status-resolved);
    color: white;
  }
  
  .btn-action.pin {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
  }
  
  .btn-action:hover:not(:disabled) {
    opacity: 0.9;
  }
</style>