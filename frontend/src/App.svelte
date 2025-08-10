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
  let pollInterval: number;
  
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
      pollInterval = setInterval(pollData, $settings.refreshInterval);
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
    // This is handled by individual components
  }
</script>

<div class="app">
  <TopBar />
  <div class="main-container">
    <ServicesPanel />
    <IncidentsPanel bind:detailsPanelCollapsed />
    <AlertDetailsPanel collapsed={detailsPanelCollapsed} />
  </div>
</div>

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }
  
  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    background: #0a0b0d;
    color: #e0e6ed;
    overflow: hidden;
  }
  
  .app {
    height: 100vh;
    display: flex;
    flex-direction: column;
  }
  
  .main-container {
    display: flex;
    height: calc(100vh - 48px);
    position: relative;
  }
  
  /* Global Scrollbar Styling */
  :global(::-webkit-scrollbar) {
    width: 8px;
    height: 8px;
  }
  
  :global(::-webkit-scrollbar-track) {
    background: #0f1114;
  }
  
  :global(::-webkit-scrollbar-thumb) {
    background: #2a2d33;
    border-radius: 4px;
  }
  
  :global(::-webkit-scrollbar-thumb:hover) {
    background: #3a3d43;
  }
</style>