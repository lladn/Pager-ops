<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Incident } from '../types';
  
  export let incident: Incident;
  export let selected = false;
  
  const dispatch = createEventDispatcher<{
    select: Incident;
    action: { type: string; incident: Incident };
  }>();
  
  function formatTime(date: Date | string): string {
    const now = new Date();
    const then = new Date(date);
    const diffMs = now.getTime() - then.getTime();
    const diff = Math.floor(diffMs / 60000); // Convert to minutes
    
    if (diff < 1) return 'just now';
    if (diff < 60) return `${diff} min ago`;
    if (diff < 1440) return `${Math.floor(diff / 60)} hours ago`;
    return `${Math.floor(diff / 1440)} days ago`;
  }
  
  function handleClick() {
    dispatch('select', incident);
  }
  
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleClick();
    }
  }
  
  function handleAction(event: MouseEvent, type: string) {
    event.stopPropagation();
    dispatch('action', { type, incident });
  }
  
  function getStatusIcon(status: string): string {
    switch (status) {
      case 'triggered':
        return '🔴';
      case 'acknowledged':
        return '🟡';
      case 'resolved':
        return '🟢';
      default:
        return '⚪';
    }
  }
</script>

<div 
  class="incident-card" 
  class:selected
  class:pinned={incident.pinned}
  class:urgency-high={incident.urgency === 'high'}
  class:urgency-low={incident.urgency === 'low'}
  on:click={handleClick}
  on:keydown={handleKeydown}
  role="button"
  tabindex="0"
  aria-label="{incident.title} - {incident.status} - {incident.urgency} urgency"
  aria-pressed={selected}
>
  {#if incident.escalatesIn}
    <div class="escalation-warning" aria-label="Escalates in {incident.escalatesIn} minutes">
      ⏰ Escalates in {incident.escalatesIn}m
    </div>
  {/if}
  
  {#if incident.pinned}
    <div class="pin-indicator" aria-label="Pinned incident">📌</div>
  {/if}
  
  <div class="incident-header">
    <span class="incident-status {incident.status}">
      <span class="status-icon" aria-hidden="true">{getStatusIcon(incident.status)}</span>
      <span class="status-text">{incident.status.charAt(0).toUpperCase() + incident.status.slice(1)}</span>
    </span>
    <span class="incident-urgency {incident.urgency}" aria-label="{incident.urgency} urgency">
      {incident.urgency.toUpperCase()}
    </span>
  </div>
  
  <div class="incident-service">{incident.service}</div>
  
  <div class="incident-title">{incident.title}</div>
  
  {#if incident.description}
    <div class="incident-description">{incident.description}</div>
  {/if}
  
  <div class="incident-meta">
    <div class="incident-time">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
        <circle cx="12" cy="12" r="10"></circle>
        <polyline points="12 6 12 12 16 14"></polyline>
      </svg>
      <span>{formatTime(incident.createdAt)}</span>
    </div>
    
    {#if incident.alertCount > 0}
      <div class="incident-alerts">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
          <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
        </svg>
        <span>{incident.alertCount} {incident.alertCount === 1 ? 'alert' : 'alerts'}</span>
      </div>
    {/if}
  </div>
  
  <div class="incident-actions">
    {#if incident.status === 'triggered'}
      <button 
        class="action-btn acknowledge"
        on:click={(e) => handleAction(e, 'acknowledge')}
        aria-label="Acknowledge incident"
      >
        Acknowledge
      </button>
    {/if}
    
    {#if incident.status !== 'resolved'}
      <button 
        class="action-btn resolve"
        on:click={(e) => handleAction(e, 'resolve')}
        aria-label="Resolve incident"
      >
        Resolve
      </button>
    {/if}
    
    <button 
      class="action-btn more"
      on:click={(e) => handleAction(e, 'more')}
      aria-label="More actions"
    >
      ⋯
    </button>
  </div>
</div>

<style>
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
  
  .escalation-warning {
    position: absolute;
    top: 8px;
    right: 8px;
    background: var(--status-escalated);
    color: white;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
  }
  
  .pin-indicator {
    position: absolute;
    top: 8px;
    left: 8px;
    font-size: 16px;
  }
  
  .incident-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
  }
  
  .incident-status {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
  }
  
  .status-icon {
    font-size: 10px;
  }
  
  .incident-status.triggered {
    color: var(--status-triggered);
  }
  
  .incident-status.acknowledged {
    color: var(--status-acknowledged);
  }
  
  .incident-status.resolved {
    color: var(--status-resolved);
  }
  
  .incident-urgency {
    font-size: 10px;
    font-weight: 600;
    padding: 2px 6px;
    border-radius: 4px;
  }
  
  .incident-urgency.high {
    background: rgba(239, 68, 68, 0.1);
    color: var(--urgency-high);
  }
  
  .incident-urgency.low {
    background: rgba(153, 153, 153, 0.1);
    color: var(--urgency-low);
  }
  
  .incident-service {
    font-size: 12px;
    color: var(--text-secondary);
    margin-bottom: 4px;
  }
  
  .incident-title {
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
    line-height: 1.4;
  }
  
  .incident-description {
    color: var(--text-secondary);
    font-size: 14px;
    margin-bottom: 12px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .incident-meta {
    display: flex;
    gap: 16px;
    margin-bottom: 12px;
    font-size: 12px;
    color: var(--text-muted);
  }
  
  .incident-time,
  .incident-alerts {
    display: flex;
    align-items: center;
    gap: 4px;
  }
  
  .incident-actions {
    display: flex;
    gap: 8px;
  }
  
  .action-btn {
    padding: 6px 12px;
    border: none;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .action-btn:hover {
    opacity: 0.9;
  }
  
  .action-btn.acknowledge {
    background: var(--status-acknowledged);
    color: white;
  }
  
  .action-btn.resolve {
    background: var(--status-resolved);
    color: white;
  }
  
  .action-btn.more {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
    padding: 6px 8px;
  }
  
  .action-btn.more:hover {
    background: var(--bg-hover);
  }
</style>