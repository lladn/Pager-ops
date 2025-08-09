<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime.js';
  import { 
    GetSettings, 
    GetServices, 
    GetIncidents,
    GetCurrentUser,
    ToggleService,
    AcknowledgeIncident,
    ResolveIncident,
    PinIncident,
    GetIncidentAlerts,
    GetDraftNote,
    SaveDraftNote
  } from '../wailsjs/go/main/App.js';
  
  // Import components from the components folder
  import Header from './components/Header.svelte';
  import ServicesPanel from './components/ServicesPanel.svelte';
  import IncidentsPanel from './components/IncidentsPanel.svelte';
  import AlertDetailsPanel from './components/AlertDetailsPanel.svelte';
  import SetupModal from './components/SetupModal.svelte';
  import SettingsModal from './components/SettingsModal.svelte';
  
  // State
  let setupRequired = false;
  let showSettings = false;
  let currentUser = null;
  let settings = null;
  let services = [];
  let incidents = [];
  let selectedIncident = null;
  let incidentAlerts = [];
  let draftNote = null;
  let theme = 'dark';
  
  // Statistics
  let stats = {
    total: 0,
    triggered: 0,
    acknowledged: 0,
    resolved: 0,
    newCount: 0
  };

  // Initialize app
  onMount(async () => {
    // Load initial data
    await loadSettings();
    await loadServices();
    await loadIncidents();
    await loadUser();
    
    // Set up event listeners
    EventsOn('setup:required', () => {
      setupRequired = true;
    });
    
    EventsOn('setup:complete', async () => {
      setupRequired = false;
      await loadSettings();
      await loadServices();
      await loadIncidents();
    });
    
    EventsOn('incidents:updated', async (data) => {
      await loadIncidents();
      stats = {
        total: data.totalCount || 0,
        triggered: data.triggeredCount || 0,
        acknowledged: data.acknowledgedCount || 0,
        resolved: data.resolvedCount || 0,
        newCount: data.newIncidents?.length || 0
      };
    });
    
    EventsOn('services:updated', async () => {
      await loadServices();
    });
    
    EventsOn('settings:updated', async (newSettings) => {
      settings = newSettings;
      theme = settings?.theme || 'dark';
      document.documentElement.setAttribute('data-theme', theme);
    });
  });
  
  onDestroy(() => {
    EventsOff('setup:required');
    EventsOff('setup:complete');
    EventsOff('incidents:updated');
    EventsOff('services:updated');
    EventsOff('settings:updated');
  });
  
  async function loadSettings() {
    try {
      settings = await GetSettings();
      theme = settings?.theme || 'dark';
      document.documentElement.setAttribute('data-theme', theme);
    } catch (err) {
      console.error('Failed to load settings:', err);
    }
  }
  
  async function loadServices() {
    try {
      services = await GetServices();
    } catch (err) {
      console.error('Failed to load services:', err);
    }
  }
  
  async function loadIncidents() {
    try {
      incidents = await GetIncidents(['triggered', 'acknowledged', 'resolved']);
    } catch (err) {
      console.error('Failed to load incidents:', err);
    }
  }
  
  async function loadUser() {
    try {
      currentUser = await GetCurrentUser();
    } catch (err) {
      console.error('Failed to load user:', err);
    }
  }
  
  async function handleServiceToggle(event) {
    const { id, enabled } = event.detail;
    await ToggleService(id, enabled);
  }
  
  async function handleIncidentSelect(event) {
    selectedIncident = event.detail;
    if (selectedIncident) {
      // Load alerts for selected incident
      try {
        incidentAlerts = await GetIncidentAlerts(selectedIncident.id);
        draftNote = await GetDraftNote(selectedIncident.id);
      } catch (err) {
        console.error('Failed to load incident details:', err);
        incidentAlerts = [];
        draftNote = null;
      }
    }
  }
  
  async function handleIncidentAction(event) {
    const { action, incidentId } = event.detail;
    
    try {
      switch (action) {
        case 'acknowledge':
          await AcknowledgeIncident(incidentId);
          break;
        case 'resolve':
          await ResolveIncident(incidentId);
          break;
        case 'pin':
          await PinIncident(incidentId, true);
          break;
        case 'unpin':
          await PinIncident(incidentId, false);
          break;
      }
      await loadIncidents();
    } catch (err) {
      console.error(`Failed to ${action} incident:`, err);
    }
  }
  
  async function handleNoteSave(event) {
    const note = event.detail;
    try {
      await SaveDraftNote(note);
    } catch (err) {
      console.error('Failed to save draft note:', err);
    }
  }
  
  function handleSettingsOpen() {
    showSettings = true;
  }
  
  function handleSettingsClose() {
    showSettings = false;
  }
</script>

<main class="app" data-theme={theme}>
  <Header 
    user={currentUser}
    stats={stats}
    on:openSettings={handleSettingsOpen}
  />
  
  <div class="content">
    <div class="panel-container">
      <ServicesPanel 
        {services}
        on:toggle={handleServiceToggle}
      />
      
      <IncidentsPanel 
        {incidents}
        {services}
        selectedId={selectedIncident?.id}
        on:select={handleIncidentSelect}
        on:action={handleIncidentAction}
      />
      
      <AlertDetailsPanel 
        incident={selectedIncident}
        alerts={incidentAlerts}
        note={draftNote}
        on:saveNote={handleNoteSave}
        on:action={handleIncidentAction}
      />
    </div>
  </div>
  
  {#if setupRequired}
    <SetupModal on:complete={() => setupRequired = false} />
  {/if}
  
  {#if showSettings}
    <SettingsModal 
      {settings}
      on:close={handleSettingsClose}
    />
  {/if}
</main>

<style>
  :global(:root) {
    --bg-primary: #1b2636;
    --bg-secondary: #232b3a;
    --bg-tertiary: #2a3444;
    --text-primary: #ffffff;
    --text-secondary: #a0aec0;
    --text-tertiary: #718096;
    --border-color: #2d3748;
    --accent-blue: #4299e1;
    --accent-green: #48bb78;
    --accent-yellow: #f6ad55;
    --accent-red: #fc8181;
    --shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  
  :global([data-theme="light"]) {
    --bg-primary: #ffffff;
    --bg-secondary: #f7fafc;
    --bg-tertiary: #edf2f7;
    --text-primary: #1a202c;
    --text-secondary: #4a5568;
    --text-tertiary: #718096;
    --border-color: #e2e8f0;
  }
  
  :global(html), :global(body) {
    margin: 0;
    padding: 0;
    height: 100vh;
    overflow: hidden;
  }
  
  :global(body) {
    font-family: "Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
      "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
      sans-serif;
    color: var(--text-primary);
    background: var(--bg-primary);
  }
  
  .app {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--bg-primary);
    color: var(--text-primary);
  }
  
  .content {
    flex: 1;
    overflow: hidden;
  }
  
  .panel-container {
    display: flex;
    height: 100%;
  }
  
  :global(.panel) {
    display: flex;
    flex-direction: column;
    background: var(--bg-secondary);
    border-right: 1px solid var(--border-color);
    overflow: hidden;
  }
  
  :global(.panel-header) {
    padding: 1rem;
    background: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-color);
    font-weight: 600;
    font-size: 0.875rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-secondary);
  }
  
  :global(.panel-content) {
    flex: 1;
    overflow-y: auto;
    padding: 1rem;
  }
  
  :global(button) {
    font-family: inherit;
    cursor: pointer;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    font-weight: 500;
    transition: all 0.2s;
  }
  
  :global(button:hover) {
    transform: translateY(-1px);
    box-shadow: var(--shadow);
  }
  
  :global(button:active) {
    transform: translateY(0);
  }
  
  :global(.btn-primary) {
    background: var(--accent-blue);
    color: white;
  }
  
  :global(.btn-success) {
    background: var(--accent-green);
    color: white;
  }
  
  :global(.btn-warning) {
    background: var(--accent-yellow);
    color: white;
  }
  
  :global(.btn-danger) {
    background: var(--accent-red);
    color: white;
  }
  
  :global(.btn-secondary) {
    background: var(--bg-tertiary);
    color: var(--text-primary);
  }
  
  :global(.badge) {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
  }
  
  :global(.status-triggered) {
    background: rgba(252, 129, 129, 0.2);
    color: var(--accent-red);
  }
  
  :global(.status-acknowledged) {
    background: rgba(246, 173, 85, 0.2);
    color: var(--accent-yellow);
  }
  
  :global(.status-resolved) {
    background: rgba(72, 187, 120, 0.2);
    color: var(--accent-green);
  }
  
  :global(.urgency-high) {
    border-left: 4px solid var(--accent-red);
  }
  
  :global(.urgency-low) {
    border-left: 4px solid var(--accent-green);
  }
</style>