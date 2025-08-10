<script lang="ts">
  import { selectedIncident, incidents, resolvedIncidents } from '../stores';
  import type { Note } from '../types';
  
  export let collapsed = false;
  
  let newNote = '';
  let detailsPanelWidth = 400;
  let isResizing = false;
  let startX = 0;
  let startWidth = 0;
  
  function togglePanel() {
    collapsed = !collapsed;
  }
  
  async function acknowledgeIncident() {
    if (!$selectedIncident) return;
    
    try {
      // await AcknowledgeIncident($selectedIncident.id);
      $selectedIncident.status = 'acknowledged';
      incidents.update(items => items);
    } catch (error) {
      console.error('Failed to acknowledge incident:', error);
    }
  }
  
  async function escalateIncident() {
    if (!$selectedIncident) return;
    
    try {
      // await EscalateIncident($selectedIncident.id);
      console.log('Escalating incident:', $selectedIncident.id);
    } catch (error) {
      console.error('Failed to escalate incident:', error);
    }
  }
  
  async function resolveIncident() {
    if (!$selectedIncident) return;
    
    try {
      // await ResolveIncident($selectedIncident.id);
      const resolved = { ...$selectedIncident, status: 'resolved' as const };
      
      incidents.update(items => items.filter(i => i.id !== $selectedIncident!.id));
      resolvedIncidents.update(items => [resolved, ...items]);
      selectedIncident.set(null);
    } catch (error) {
      console.error('Failed to resolve incident:', error);
    }
  }
  
  async function addNote() {
    if (!$selectedIncident || !newNote.trim()) return;
    
    try {
      // await AddIncidentNote($selectedIncident.id, newNote);
      const note: Note = {
        author: 'Current User',
        content: newNote,
        timestamp: new Date()
      };
      
      if (!$selectedIncident.notes) {
        $selectedIncident.notes = [];
      }
      $selectedIncident.notes = [note, ...$selectedIncident.notes];
      newNote = '';
      incidents.update(items => items);
    } catch (error) {
      console.error('Failed to add note:', error);
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
            <div class="field-value">{$selectedIncident.createdAt.toLocaleString()}</div>
          </div>
        </div>

        <div class="detail-section">
          <div class="section-title">Notes</div>
          <div class="notes-container">
            {#if $selectedIncident.notes}
              {#each $selectedIncident.notes as note}
                <div class="note-item">
                  <div class="note-author">{note.author} - {formatTime(note.timestamp)}</div>
                  <div class="note-content">{note.content}</div>
                </div>
              {/each}
            {/if}
          </div>
          <textarea 
            class="add-note-input" 
            placeholder="Add a note..."
            bind:value={newNote}
            on:keydown={(e) => e.key === 'Enter' && !e.shiftKey && addNote()}
          ></textarea>
        </div>
      </div>
      <div class="action-buttons">
        {#if $selectedIncident.status === 'triggered'}
          <button class="action-btn primary" on:click={acknowledgeIncident}>Acknowledge</button>
        {/if}
        <button class="action-btn secondary" on:click={escalateIncident}>Escalate</button>
        {#if $selectedIncident.status !== 'resolved'}
          <button class="action-btn danger" on:click={resolveIncident}>Resolve</button>
        {/if}
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
    background: #141619;
    border-left: 1px solid #2a2d33;
    display: flex;
    flex-direction: column;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
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
    transition: background 0.2s;
  }

  .resize-handle:hover {
    background: #3b82f6;
  }

  .alert-details-header {
    padding: 16px;
    border-bottom: 1px solid #2a2d33;
    display: flex;
    justify-content: space-between;
    align-items: center;
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

  .alert-details-content {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }

  .no-selection {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #8b92a9;
    font-size: 14px;
  }

  .detail-section {
    margin-bottom: 24px;
  }

  .section-title {
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    color: #8b92a9;
    margin-bottom: 12px;
    letter-spacing: 0.5px;
  }

  .detail-field {
    margin-bottom: 12px;
  }

  .field-label {
    font-size: 11px;
    color: #8b92a9;
    margin-bottom: 4px;
  }

  .field-value {
    font-size: 13px;
    color: #e0e6ed;
    padding: 8px 12px;
    background: #0f1114;
    border-radius: 6px;
    word-break: break-word;
  }

  .notes-container {
    background: #0f1114;
    border-radius: 8px;
    padding: 12px;
    margin-bottom: 12px;
  }

  .note-item {
    padding: 8px;
    border-left: 2px solid #3b82f6;
    margin-bottom: 8px;
  }

  .note-item:last-child {
    margin-bottom: 0;
  }

  .note-author {
    font-size: 11px;
    color: #8b92a9;
    margin-bottom: 4px;
  }

  .note-content {
    font-size: 13px;
    color: #e0e6ed;
  }

  .add-note-input {
    width: 100%;
    padding: 8px 12px;
    background: #0f1114;
    border: 1px solid #2a2d33;
    border-radius: 6px;
    color: #e0e6ed;
    font-size: 13px;
    font-family: inherit;
    outline: none;
    resize: vertical;
    min-height: 60px;
  }

  .add-note-input:focus {
    border-color: #3b82f6;
    background: #1a1d21;
  }

  .action-buttons {
    padding: 16px;
    border-top: 1px solid #2a2d33;
    display: flex;
    gap: 8px;
  }

  .action-btn {
    flex: 1;
    padding: 8px 16px;
    border: none;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn.primary {
    background: #3b82f6;
    color: white;
  }

  .action-btn.primary:hover {
    background: #2563eb;
  }

  .action-btn.secondary {
    background: #1e2025;
    color: #e0e6ed;
    border: 1px solid #2a2d33;
  }

  .action-btn.secondary:hover {
    background: #2a2d33;
  }

  .action-btn.danger {
    background: rgba(239, 68, 68, 0.2);
    color: #ef4444;
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .action-btn.danger:hover {
    background: rgba(239, 68, 68, 0.3);
  }
</style>