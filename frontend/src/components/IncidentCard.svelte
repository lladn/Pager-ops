<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Incident } from '../types';
  
  export let incident: Incident;
  export let selected = false;
  
  const dispatch = createEventDispatcher();
  
  function formatTime(date: Date): string {
    const now = new Date();
    const diff = Math.floor((now.getTime() - date.getTime()) / 60000);
    
    if (diff < 1) return 'just now';
    if (diff < 60) return `${diff} min ago`;
    if (diff < 1440) return `${Math.floor(diff / 60)} hours ago`;
    return `${Math.floor(diff / 1440)} days ago`;
  }
</script>

<div 
  class="incident-card" 
  class:selected
  on:click={() => dispatch('select')}
>
  {#if incident.escalatesIn}
    <div class="escalation-warning">Escalates in {incident.escalatesIn}m</div>
  {/if}
  <div class="incident-header">
    <span class="incident-status {incident.status}">
      <span class="status-dot {incident.status}"></span>
      {incident.status.charAt(0).toUpperCase() + incident.status.slice(1)}
    </span>
    <span class="incident-urgency {incident.urgency}">{incident.urgency.toUpperCase()}</span>
  </div>
  <div class="incident-title">{incident.title}</div>
  <div class="incident-meta">
    <div class="incident-time">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <polyline points="12 6 12 12 16 14"></polyline>
      </svg>
      <span>{formatTime(incident.createdAt)}</span>
    </div>
    <div class="incident-alerts">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
      </svg>
      <span>{incident.alertCount} {incident.alertCount === 1 ? 'alert' : 'alerts'}</span>
    </div>
  </div>
</div>

<style>
  .incident-card {
    background: #141619;
    border: 1px solid #2a2d33;
    border-radius: 12px;
    padding: 12px;
    margin-bottom: 8px;
    cursor: pointer;
    transition: all 0.2s;
    position: relative;
  }

  .incident-card:hover {
    background: #1a1d21;
    border-color: #3a3d43;
    transform: translateY(-1px);
  }

  .incident-card.selected {
    background: #1e2025;
    border-color: #3b82f6;
  }

  .incident-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 8px;
  }

  .incident-status {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
  }

  .status-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
  }

  .incident-status.triggered {
    background: rgba(239, 68, 68, 0.2);
    color: #ef4444;
  }

  .status-dot.triggered { background: #ef4444; }

  .incident-status.acknowledged {
    background: rgba(234, 179, 8, 0.2);
    color: #eab308;
  }

  .status-dot.acknowledged { background: #eab308; }

  .incident-status.resolved {
    background: rgba(34, 197, 94, 0.2);
    color: #22c55e;
  }

  .status-dot.resolved { background: #22c55e; }

  .incident-urgency {
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
  }

  .incident-urgency.high {
    background: rgba(239, 68, 68, 0.2);
    color: #ef4444;
  }

  .incident-urgency.low {
    background: rgba(107, 114, 128, 0.2);
    color: #9ca3af;
  }

  .incident-title {
    font-size: 14px;
    font-weight: 500;
    color: #e0e6ed;
    margin-bottom: 4px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .incident-meta {
    display: flex;
    gap: 12px;
    font-size: 12px;
    color: #8b92a9;
  }

  .incident-time {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .incident-alerts {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .escalation-warning {
    position: absolute;
    top: 12px;
    right: 12px;
    background: rgba(239, 68, 68, 0.2);
    color: #ef4444;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 11px;
    font-weight: 600;
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.6; }
  }
</style>