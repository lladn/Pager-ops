<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { OpenIncidentInBrowser } from '../../wailsjs/go/main/App.js';
    
    export let incident = null;
    export let alerts = [];
    export let note = null;
    
    const dispatch = createEventDispatcher();
    
    let activeTab = 'details';
    let draftNote = {
      incident_id: '',
      why_triggered: '',
      impact: '',
      actions: '',
      links: '',
      note_text: ''
    };
    
    // Update draft when incident changes
    $: if (incident) {
      draftNote.incident_id = incident.id;
      if (note) {
        draftNote = { ...note };
      } else {
        draftNote = {
          incident_id: incident.id,
          why_triggered: '',
          impact: '',
          actions: '',
          links: '',
          note_text: ''
        };
      }
    }
    
    function saveNote() {
      dispatch('saveNote', draftNote);
    }
    
    function openInBrowser() {
      if (incident?.html_url) {
        OpenIncidentInBrowser(incident.html_url);
      }
    }
    
    function copyToClipboard(text) {
      navigator.clipboard.writeText(text);
    }
    
    function formatDate(date) {
      return new Date(date).toLocaleString();
    }
    
    function formatAlertDetails(details) {
      if (!details) return '';
      return JSON.stringify(details, null, 2);
    }
  </script>
  
  <div class="panel alert-details-panel">
    <div class="panel-header">
      <div class="title">Alert Details</div>
      {#if incident}
        <div class="tabs">
          <button 
            class="tab"
            class:active={activeTab === 'details'}
            on:click={() => activeTab = 'details'}
          >
            Details
          </button>
          <button 
            class="tab"
            class:active={activeTab === 'alerts'}
            on:click={() => activeTab = 'alerts'}
          >
            Alerts ({alerts.length})
          </button>
          <button 
            class="tab"
            class:active={activeTab === 'notes'}
            on:click={() => activeTab = 'notes'}
          >
            Notes
          </button>
        </div>
      {/if}
    </div>
    
    <div class="panel-content">
      {#if !incident}
        <div class="empty-state">
          <p>Select an incident to view details</p>
        </div>
      {:else}
        {#if activeTab === 'details'}
          <div class="details-content">
            <div class="detail-section">
              <h3>Incident Information</h3>
              
              <div class="detail-row">
                <span class="label">ID:</span>
                <span class="value">
                  {incident.id}
                  <button 
                    class="copy-btn"
                    on:click={() => copyToClipboard(incident.id)}
                    title="Copy ID"
                  >
                    📋
                  </button>
                </span>
              </div>
              
              <div class="detail-row">
                <span class="label">Number:</span>
                <span class="value">#{incident.incident_number}</span>
              </div>
              
              <div class="detail-row">
                <span class="label">Status:</span>
                <span class="value">
                  <span class="badge status-{incident.status}">
                    {incident.status}
                  </span>
                </span>
              </div>
              
              <div class="detail-row">
                <span class="label">Urgency:</span>
                <span class="value">
                  <span class="badge urgency-{incident.urgency}">
                    {incident.urgency}
                  </span>
                </span>
              </div>
              
              <div class="detail-row">
                <span class="label">Created:</span>
                <span class="value">{formatDate(incident.created_at)}</span>
              </div>
              
              <div class="detail-row">
                <span class="label">Updated:</span>
                <span class="value">{formatDate(incident.updated_at)}</span>
              </div>
            </div>
            
            <div class="detail-section">
              <h3>Description</h3>
              <div class="description">
                {incident.description || 'No description provided'}
              </div>
            </div>
            
            <div class="detail-section">
              <h3>Actions</h3>
              <div class="action-buttons">
                <button class="btn-primary" on:click={openInBrowser}>
                  🌐 Open in Browser
                </button>
                <button 
                  class="btn-success"
                  on:click={() => dispatch('action', { action: 'acknowledge', incidentId: incident.id })}
                  disabled={incident.status !== 'triggered'}
                >
                  ✓ Acknowledge
                </button>
                <button 
                  class="btn-success"
                  on:click={() => dispatch('action', { action: 'resolve', incidentId: incident.id })}
                  disabled={incident.status === 'resolved'}
                >
                  ✅ Resolve
                </button>
              </div>
            </div>
          </div>
        {/if}
        
        {#if activeTab === 'alerts'}
          <div class="alerts-content">
            {#if alerts.length === 0}
              <p class="no-alerts">No alerts for this incident</p>
            {:else}
              {#each alerts as alert}
                <div class="alert-item">
                  <div class="alert-header">
                    <span class="alert-id">Alert: {alert.id}</span>
                    <span class="badge status-{alert.status}">{alert.status}</span>
                  </div>
                  <div class="alert-summary">{alert.summary}</div>
                  {#if alert.details}
                    <details class="alert-details">
                      <summary>View Details</summary>
                      <pre>{formatAlertDetails(alert.details)}</pre>
                    </details>
                  {/if}
                  <div class="alert-time">Created: {alert.created_at}</div>
                </div>
              {/each}
            {/if}
          </div>
        {/if}
        
        {#if activeTab === 'notes'}
          <div class="notes-content">
            <div class="note-form">
              <div class="form-group">
                <label>Why Triggered?</label>
                <textarea 
                  bind:value={draftNote.why_triggered}
                  placeholder="What caused this incident?"
                  rows="3"
                ></textarea>
              </div>
              
              <div class="form-group">
                <label>Impact</label>
                <textarea 
                  bind:value={draftNote.impact}
                  placeholder="What's the impact of this incident?"
                  rows="3"
                ></textarea>
              </div>
              
              <div class="form-group">
                <label>Actions Taken</label>
                <textarea 
                  bind:value={draftNote.actions}
                  placeholder="What actions were taken to resolve?"
                  rows="3"
                ></textarea>
              </div>
              
              <div class="form-group">
                <label>Supporting Links</label>
                <textarea 
                  bind:value={draftNote.links}
                  placeholder="Dashboards, logs, documentation..."
                  rows="2"
                ></textarea>
              </div>
              
              <div class="form-group">
                <label>Additional Notes</label>
                <textarea 
                  bind:value={draftNote.note_text}
                  placeholder="Any other relevant information..."
                  rows="4"
                ></textarea>
              </div>
              
              <button class="btn-primary save-note" on:click={saveNote}>
                💾 Save Draft
              </button>
            </div>
          </div>
        {/if}
      {/if}
    </div>
  </div>
  
  <style>
    .alert-details-panel {
      width: 400px;
      border-right: none;
    }
    
    .tabs {
      display: flex;
      gap: 0.5rem;
      margin-top: 0.5rem;
    }
    
    .tab {
      padding: 0.25rem 0.75rem;
      background: var(--bg-primary);
      color: var(--text-secondary);
      border-radius: 4px 4px 0 0;
      font-size: 0.875rem;
    }
    
    .tab.active {
      background: var(--accent-blue);
      color: white;
    }
    
    .empty-state {
      text-align: center;
      padding: 4rem 2rem;
      color: var(--text-tertiary);
    }
    
    .details-content {
      padding: 1rem;
    }
    
    .detail-section {
      margin-bottom: 1.5rem;
    }
    
    .detail-section h3 {
      margin: 0 0 1rem 0;
      font-size: 1rem;
      font-weight: 600;
      color: var(--text-secondary);
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    
    .detail-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.5rem 0;
      border-bottom: 1px solid var(--border-color);
    }
    
    .label {
      font-weight: 500;
      color: var(--text-secondary);
      font-size: 0.875rem;
    }
    
    .value {
      color: var(--text-primary);
      font-size: 0.875rem;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }
    
    .copy-btn {
      padding: 0.125rem 0.25rem;
      background: transparent;
      font-size: 0.75rem;
      opacity: 0.7;
    }
    
    .copy-btn:hover {
      opacity: 1;
      transform: none;
    }
    
    .description {
      padding: 1rem;
      background: var(--bg-primary);
      border-radius: 4px;
      font-size: 0.875rem;
      line-height: 1.5;
      color: var(--text-primary);
    }
    
    .action-buttons {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
    }
    
    .action-buttons button {
      flex: 1;
      min-width: 120px;
    }
    
    .action-buttons button:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    .alerts-content {
      padding: 1rem;
    }
    
    .no-alerts {
      text-align: center;
      color: var(--text-tertiary);
      padding: 2rem;
    }
    
    .alert-item {
      padding: 1rem;
      background: var(--bg-primary);
      border-radius: 4px;
      margin-bottom: 0.75rem;
    }
    
    .alert-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
    }
    
    .alert-id {
      font-weight: 600;
      color: var(--text-secondary);
      font-size: 0.875rem;
    }
    
    .alert-summary {
      color: var(--text-primary);
      margin-bottom: 0.5rem;
    }
    
    .alert-details {
      margin-top: 0.5rem;
    }
    
    .alert-details summary {
      cursor: pointer;
      color: var(--accent-blue);
      font-size: 0.875rem;
    }
    
    .alert-details pre {
      margin-top: 0.5rem;
      padding: 0.5rem;
      background: var(--bg-secondary);
      border-radius: 4px;
      font-size: 0.75rem;
      overflow-x: auto;
    }
    
    .alert-time {
      font-size: 0.75rem;
      color: var(--text-tertiary);
      margin-top: 0.5rem;
    }
    
    .notes-content {
      padding: 1rem;
    }
    
    .note-form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }
    
    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }
    
    .form-group label {
      font-weight: 500;
      color: var(--text-secondary);
      font-size: 0.875rem;
    }
    
    .form-group textarea {
      padding: 0.5rem;
      background: var(--bg-primary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-primary);
      font-family: inherit;
      font-size: 0.875rem;
      resize: vertical;
    }
    
    .form-group textarea:focus {
      outline: none;
      border-color: var(--accent-blue);
    }
    
    .save-note {
      margin-top: 0.5rem;
    }
    
    .badge.urgency-high {
      background: rgba(252, 129, 129, 0.2);
      color: var(--accent-red);
    }
    
    .badge.urgency-low {
      background: rgba(72, 187, 120, 0.2);
      color: var(--accent-green);
    }
  </style>