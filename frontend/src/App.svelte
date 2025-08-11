<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime.js';
  import TopBar from './components/TopBar.svelte';
  import ServicesPanel from './components/ServicesPanel.svelte';
  import IncidentsPanel from './components/IncidentsPanel.svelte';
  import AlertDetailsPanel from './components/AlertDetailsPanel.svelte';
  import { incidents, resolvedIncidents, selectedIncident, services, settings } from './stores';
  import { loadIncidents, loadServices, pollData } from './services/api';
  import type { Incident, Service } from './types';
  
  let detailsPanelCollapsed = false;
  let pollInterval: ReturnType<typeof setInterval> | undefined;
  
  onMount(async () => {
    // Set up event listeners for backend updates
    EventsOn('incident:update', handleIncidentUpdate);
    EventsOn('service:update', handleServiceUpdate);
    EventsOn('incidents:refresh', handleIncidentsRefresh);
    
    // Load initial data
    await initialize();
    
    // Set up polling based on settings
    const unsubscribe = settings.subscribe(($settings) => {
      if (pollInterval) {
        clearInterval(pollInterval);
      }
      if ($settings.refreshInterval > 0) {
        pollInterval = setInterval(pollData, $settings.refreshInterval);
      }
    });
    
    // Handle click outside for dropdowns
    document.addEventListener('click', handleClickOutside);
    
    return () => {
      unsubscribe();
      if (pollInterval) {
        clearInterval(pollInterval);
      }
      document.removeEventListener('click', handleClickOutside);
    };
  });
  
  onDestroy(() => {
    EventsOff('incident:update');
    EventsOff('service:update');
    EventsOff('incidents:refresh');
  });
  
  async function initialize() {
    try {
      await loadServices();
      await loadIncidents();
    } catch (error) {
      console.error('Failed to initialize:', error);
      // You might want to show a user-friendly error message here
    }
  }
  
  function handleIncidentUpdate(incident: Incident) {
    incidents.update(items => {
      const index = items.findIndex(i => i.id === incident.id);
      if (index >= 0) {
        items[index] = incident;
      } else {
        items.push(incident);
      }
      return [...items];
    });
    
    // Update service counts
    updateServiceCounts();
  }
  
  function handleServiceUpdate(service: Service) {
    services.update(items => {
      const index = items.findIndex(s => s.id === service.id);
      if (index >= 0) {
        items[index] = service;
      }
      return [...items];
    });
  }
  
  async function handleIncidentsRefresh() {
    await loadIncidents();
  }
  
  function handleClickOutside(event: MouseEvent) {
    // Close dropdowns when clicking outside
    const target = event.target as HTMLElement;
    if (!target.closest('.dropdown-container')) {
      // Dispatch custom event to close all dropdowns
      window.dispatchEvent(new CustomEvent('close-dropdowns'));
    }
  }
  
  function updateServiceCounts() {
    services.update(items => {
      const currentIncidents = $incidents;
      
      items.forEach(service => {
        service.incidentCount = currentIncidents.filter(
          i => i.service === service.name && i.status !== 'resolved'
        ).length;
      });
      
      return [...items];
    });
  }
  
  function toggleDetailsPanel() {
    detailsPanelCollapsed = !detailsPanelCollapsed;
  }
  
  // Subscribe to selectedIncident to ensure details panel shows when needed
  $: if ($selectedIncident) {
    detailsPanelCollapsed = false;
  }
</script>

<main class="app">
  <TopBar />
  
  <div class="main-content">
    <ServicesPanel />
    
    <div class="incidents-container">
      <IncidentsPanel />
    </div>
    
    {#if $selectedIncident}
      <div class="details-panel" class:collapsed={detailsPanelCollapsed}>
        <button class="collapse-toggle" on:click={toggleDetailsPanel}>
          {detailsPanelCollapsed ? '◀' : '▶'}
        </button>
        <AlertDetailsPanel />
      </div>
    {/if}
  </div>
</main>

<style>
  .app {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: var(--bg-primary);
    color: var(--text-primary);
    overflow: hidden;
  }
  
  .main-content {
    display: flex;
    flex: 1;
    overflow: hidden;
  }
  
  .incidents-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }
  
  .details-panel {
    width: 400px;
    border-left: 1px solid var(--border-color);
    background: var(--bg-secondary);
    transition: width 0.3s ease;
    position: relative;
  }
  
  .details-panel.collapsed {
    width: 40px;
  }
  
  .collapse-toggle {
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-left: none;
    border-radius: 0 4px 4px 0;
    padding: 8px 4px;
    color: var(--text-secondary);
    cursor: pointer;
    z-index: 10;
  }
  
  .collapse-toggle:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }
  
  :global(:root) {
    --bg-primary: #0f0f0f;
    --bg-secondary: #1a1a1a;
    --bg-tertiary: #252525;
    --bg-hover: #2a2a2a;
    --text-primary: #ffffff;
    --text-secondary: #999999;
    --text-muted: #666666;
    --border-color: #2a2a2a;
    --status-triggered: #ff4444;
    --status-acknowledged: #ffaa00;
    --status-resolved: #00aa00;
    --status-escalated: #ff6600;
    --urgency-high: #ff4444;
    --urgency-low: #999999;
  }
  
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
</style>