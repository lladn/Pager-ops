<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    
    export let user = null;
    export let stats = {
      total: 0,
      triggered: 0,
      acknowledged: 0,
      resolved: 0,
      newCount: 0
    };
    
    const dispatch = createEventDispatcher();
    
    function openSettings() {
      dispatch('openSettings');
    }
  </script>
  
  <header class="header">
    <div class="header-left">
      <h1 class="app-title">
        <span class="logo">🚨</span>
        PagerOps Command Center
      </h1>
      {#if user}
        <div class="user-info">
          <span class="user-name">{user.name}</span>
          <span class="user-email">{user.email}</span>
        </div>
      {/if}
    </div>
    
    <div class="header-center">
      <div class="stats">
        <div class="stat" class:has-incidents={stats.triggered > 0}>
          <span class="stat-value">{stats.triggered}</span>
          <span class="stat-label">Triggered</span>
        </div>
        <div class="stat" class:has-incidents={stats.acknowledged > 0}>
          <span class="stat-value">{stats.acknowledged}</span>
          <span class="stat-label">Acknowledged</span>
        </div>
        <div class="stat">
          <span class="stat-value">{stats.resolved}</span>
          <span class="stat-label">Resolved (24h)</span>
        </div>
      </div>
      {#if stats.newCount > 0}
        <div class="new-badge">
          {stats.newCount} new
        </div>
      {/if}
    </div>
    
    <div class="header-right">
      <button class="settings-btn" on:click={openSettings} title="Settings">
        ⚙️
      </button>
    </div>
  </header>
  
  <style>
    .header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1rem 1.5rem;
      background: var(--bg-tertiary);
      border-bottom: 2px solid var(--border-color);
      box-shadow: var(--shadow);
    }
    
    .header-left {
      display: flex;
      align-items: center;
      gap: 2rem;
    }
    
    .app-title {
      margin: 0;
      font-size: 1.25rem;
      font-weight: 700;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      color: var(--text-primary);
    }
    
    .logo {
      font-size: 1.5rem;
      animation: pulse 2s infinite;
    }
    
    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.7; }
    }
    
    .user-info {
      display: flex;
      flex-direction: column;
      font-size: 0.875rem;
    }
    
    .user-name {
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .user-email {
      color: var(--text-secondary);
      font-size: 0.75rem;
    }
    
    .header-center {
      display: flex;
      align-items: center;
      gap: 1rem;
    }
    
    .stats {
      display: flex;
      gap: 2rem;
    }
    
    .stat {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 0.5rem 1rem;
      border-radius: 8px;
      background: var(--bg-secondary);
      transition: all 0.3s;
    }
    
    .stat.has-incidents {
      background: rgba(252, 129, 129, 0.1);
      animation: glow 2s infinite;
    }
    
    @keyframes glow {
      0%, 100% { box-shadow: 0 0 0 0 rgba(252, 129, 129, 0); }
      50% { box-shadow: 0 0 20px 5px rgba(252, 129, 129, 0.3); }
    }
    
    .stat-value {
      font-size: 1.5rem;
      font-weight: 700;
      color: var(--text-primary);
    }
    
    .stat.has-incidents .stat-value {
      color: var(--accent-red);
    }
    
    .stat-label {
      font-size: 0.75rem;
      color: var(--text-secondary);
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    
    .new-badge {
      padding: 0.5rem 1rem;
      background: var(--accent-red);
      color: white;
      border-radius: 9999px;
      font-size: 0.875rem;
      font-weight: 600;
      animation: bounce 1s infinite;
    }
    
    @keyframes bounce {
      0%, 100% { transform: translateY(0); }
      50% { transform: translateY(-5px); }
    }
    
    .header-right {
      display: flex;
      align-items: center;
      gap: 1rem;
    }
    
    .settings-btn {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: var(--bg-secondary);
      border: 2px solid var(--border-color);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      transition: all 0.3s;
    }
    
    .settings-btn:hover {
      background: var(--bg-primary);
      transform: rotate(90deg);
    }
  </style>