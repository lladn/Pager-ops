<script lang="ts">
  import { selectedIncident, incidents, resolvedIncidents } from '../stores';
  import { acknowledgeIncident, escalateIncident, resolveIncident, addIncidentNote } from '../services/api';
  import type { Note } from '../types';
  
  export let collapsed = false;
  
  let newNote = '';
  let detailsPanelWidth = 400;
  let isResizing = false;
  let startX = 0;
  let startWidth = 0;
  let isLoading = false;
  
  function togglePanel() {
    collapsed = !collapsed;
  }
  
  async function handleAcknowledge() {
    if (!$selectedIncident || isLoading) return;
    
    isLoading = true;
    try {
      await acknowledgeIncident($selectedIncident.id);
    } catch (error) {
      console.error('Failed to acknowledge incident:', error);
      // You might want to show a user-friendly error message here
    } finally {
      isLoading = false;
    }
  }
  
  async function handleEscalate() {
    if (!$selectedIncident || isLoading) return;
    
    isLoading = true;
    try {
      // Using default escalation level "1"
      // You can modify this to get the escalation level from settings or a dropdown
      await escalateIncident($selectedIncident.id, "1");
    } catch (error) {
      console.error('Failed to escalate incident:', error);
    } finally {
      isLoading = false;
    }
  }
  
  async function handleResolve() {
    if (!$selectedIncident || isLoading) return;
    
    isLoading = true;
    try {
      await resolveIncident($selectedIncident.id);
    } catch (error) {
      console.error('Failed to resolve incident:', error);
    } finally {
      isLoading = false;
    }
  }
  
  async function addNote() {
    if (!$selectedIncident || !newNote.trim() || isLoading) return;
    
    isLoading = true;
    try {
      await addIncidentNote($selectedIncident.id, newNote.trim());
      newNote = '';
    } catch (error) {
      console.error('Failed to add note:', error);
    } finally {
      isLoading = false;
    }
  }
  
  function formatTime(date: Date): string {
    const now = new Date();
    const diff = Math.floor((now.getTime() - date.getTime()) / 60000);
    
    if (diff < 1) return 'just now';
    if (diff < 60) return `${diff} min ago`;
    if (diff < 1440) return `${Math.floor(diff / 60)} hours ago`;
    return `${Math.floor(diff / 1440)} days ago`;
  }
  
  function initResize(event: MouseEvent) {
    isResizing = true;
    startX = event.clientX;
    startWidth = detailsPanelWidth;
    
    document.addEventListener('mousemove', doResize);
    document.addEventListener('mouseup', stopResize);
    event.preventDefault();
  }
  
  function doResize(event: MouseEvent) {
    if (!isResizing) return;
    
    const newWidth = startWidth - (event.clientX - startX);
    if (newWidth >= 300 && newWidth <= 600) {
      detailsPanelWidth = newWidth;
    }
  }
  
  function stopResize() {
    isResizing = false;
    document.removeEventListener('mousemove', doResize);
    document.removeEventListener('mouseup', stopResize);
  }
</script>

{#if !collapsed}
  <div class="alert-details-panel" style="width: {detailsPanelWidth}px">
    <div class="resize-handle" on:mousedown={initResize}></div>
    <div class="alert-details-header">
      <span class="panel-title">Alert Details</span>
      <button class="collapse-btn" on:click={togglePanel}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9 18 15 12 9 6"></polyline>
        </svg>
      </button>
    </div>
    {#if $selectedIncident}
      <div class="alert-details-content">
        <div class="detail-section">
          <div class="section-title">Incident Information</div>
          <div class="detail-field">
            <div class="field-label">Service</div>
            <div class="field-value">{$selectedIncident.service}</div>
          </div>
          <div class="detail-field">
            <div class="field-label">Title</div>
            <div class="field-value">{$selectedIncident.title}</div>
          </div>
          {#if $selectedIncident.description}
            <div class="detail-field">
              <div class="field-label">Description</div>
              <div class="field-value">{$selectedIncident.description}</div>
            </div>
          {/if}
        </div>

        <div class="detail-section">
          <div class="section-title">Timeline</div>
          <div class="detail-field">
            <div class="field-label">Created</div>
            <div class="field-value">
              {$selectedIncident.createdAt instanceof Date 
                ? $selectedIncident.createdAt.toLocaleString()
                : new Date($selectedIncident.createdAt).toLocaleString()}
            </div>
          </div>
        </div>

        <div class="detail-section">
          <div class="section-title">Notes</div>
          <div class="notes-container">
            {#if $selectedIncident.notes && $selectedIncident.notes.length > 0}
              {#each $selectedIncident.notes as note}
                <div class="note-item">
                  <div class="note-author">
                    {note.author} - {formatTime(note.timestamp instanceof Date ? note.timestamp : new Date(note.timestamp))}
                  </div>
                  <div class="note-content">{note.content}</div>
                </div>
              {/each}
            {:else}
              <div class="no-notes">No notes yet</div>
            {/if}
          </div>
          <textarea 
            class="add-note-input" 
            placeholder="Add a note..."
            bind:value={newNote}
            disabled={isLoading}
            on:keydown={(e) => e.key === 'Enter' && !e.shiftKey && addNote()}
          ></textarea>
        </div>
      </div>
      <div class="action-buttons">
        {#if $selectedIncident.status === 'triggered'}
          <button 
            class="action-btn primary" 
            on:click={handleAcknowledge}
            disabled={isLoading}
          >
            {isLoading ? 'Loading...' : 'Acknowledge'}
          </button>
        {/if}
        {#if $selectedIncident.status !== 'resolved'}
          <button 
            class="action-btn secondary" 
            on:click={handleResolve}
            disabled={isLoading}
          >
            {isLoading ? 'Loading...' : 'Resolve'}
          </button>
        {/if}
        <button 
          class="action-btn tertiary" 
          on:click={handleEscalate}
          disabled={isLoading}
        >
          {isLoading ? 'Loading...' : 'Escalate'}
        </button>
      </div>
    {:else}
      <div class="no-selection">
        <p>Select an incident to view details</p>
      </div>
    {/if}
  </div>
{/if}

<style>
  .alert-details-panel {
    display: flex;
    flex-direction: column;
    background: var(--bg-secondary);
    height: 100%;
    position: relative;
  }

  .resize-handle {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 4px;
    cursor: col-resize;
    background: transparent;
  }

  .resize-handle:hover {
    background: var(--border-color);
  }

  .alert-details-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .panel-title {
    font-weight: 600;
    font-size: 16px;
  }

  .collapse-btn {
    background: none;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    padding: 4px;
  }

  .collapse-btn:hover {
    color: var(--text-primary);
  }

  .alert-details-content {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }

  .detail-section {
    margin-bottom: 24px;
  }

  .section-title {
    font-weight: 600;
    margin-bottom: 12px;
    color: var(--text-primary);
  }

  .detail-field {
    margin-bottom: 12px;
  }

  .field-label {
    font-size: 12px;
    color: var(--text-secondary);
    margin-bottom: 4px;
  }

  .field-value {
    color: var(--text-primary);
  }

  .notes-container {
    max-height: 200px;
    overflow-y: auto;
    margin-bottom: 12px;
  }

  .note-item {
    padding: 8px;
    background: var(--bg-tertiary);
    border-radius: 4px;
    margin-bottom: 8px;
  }

  .note-author {
    font-size: 12px;
    color: var(--text-secondary);
    margin-bottom: 4px;
  }

  .note-content {
    color: var(--text-primary);
    white-space: pre-wrap;
  }

  .no-notes {
    color: var(--text-muted);
    font-style: italic;
    padding: 8px;
  }

  .add-note-input {
    width: 100%;
    padding: 8px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    color: var(--text-primary);
    resize: vertical;
    min-height: 60px;
  }

  .add-note-input:focus {
    outline: none;
    border-color: var(--status-acknowledged);
  }

  .add-note-input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .action-buttons {
    padding: 16px;
    border-top: 1px solid var(--border-color);
    display: flex;
    gap: 8px;
  }

  .action-btn {
    flex: 1;
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    transition: background-color 0.2s;
  }

  .action-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .action-btn.primary {
    background: var(--status-acknowledged);
    color: white;
  }

  .action-btn.primary:hover:not(:disabled) {
    opacity: 0.9;
  }

  .action-btn.secondary {
    background: var(--status-resolved);
    color: white;
  }

  .action-btn.secondary:hover:not(:disabled) {
    opacity: 0.9;
  }

  .action-btn.tertiary {
    background: var(--status-escalated);
    color: white;
  }

  .action-btn.tertiary:hover:not(:disabled) {
    opacity: 0.9;
  }

  .no-selection {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--text-muted);
  }
</style>