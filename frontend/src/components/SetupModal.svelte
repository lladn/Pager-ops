<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { SaveSettings, FetchAvailableServices, SaveService } from '../../wailsjs/go/main/App.js';
    
    const dispatch = createEventDispatcher();
    
    let step = 1;
    let apiKey = '';
    let theme = 'dark';
    let assignedOnly = true;
    let redirectEnabled = false;
    let soundEnabled = true;
    let availableServices = [];
    let selectedServices = [];
    let loading = false;
    let error = '';
    
    async function validateApiKey() {
      if (!apiKey) {
        error = 'API key is required';
        return;
      }
      
      loading = true;
      error = '';
      
      try {
        // Save API key first
        await SaveSettings({
          api_key: apiKey,
          theme: theme,
          assigned_only: assignedOnly,
          redirect_enabled: redirectEnabled,
          sound_enabled: soundEnabled,
          sound_path: ''
        });
        
        // Fetch available services
        availableServices = await FetchAvailableServices();
        
        // Move to service selection step
        step = 2;
      } catch (err) {
        error = 'Invalid API key or failed to connect to PagerDuty';
      } finally {
        loading = false;
      }
    }
    
    function toggleService(service) {
      const index = selectedServices.findIndex(s => s.id === service.id);
      if (index >= 0) {
        selectedServices = selectedServices.filter(s => s.id !== service.id);
      } else {
        selectedServices = [...selectedServices, {
          id: service.id,
          alias: service.name,
          enabled: true
        }];
      }
    }
    
    function isServiceSelected(serviceId) {
      return selectedServices.some(s => s.id === serviceId);
    }
    
    async function completeSetup() {
      loading = true;
      error = '';
      
      try {
        // Save selected services
        for (const service of selectedServices) {
          await SaveService(service);
        }
        
        // Complete setup
        dispatch('complete');
      } catch (err) {
        error = 'Failed to save services';
      } finally {
        loading = false;
      }
    }
    
    function skipServices() {
      dispatch('complete');
    }
  </script>
  
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h2>Welcome to PagerOps Command Center</h2>
        <div class="steps">
          <div class="step" class:active={step === 1}>
            <span class="step-number">1</span>
            <span class="step-label">API Setup</span>
          </div>
          <div class="step" class:active={step === 2}>
            <span class="step-number">2</span>
            <span class="step-label">Select Services</span>
          </div>
        </div>
      </div>
      
      <div class="modal-content">
        {#if step === 1}
          <div class="setup-step">
            <div class="intro">
              <p>Let's get you connected to PagerDuty. You'll need your API key to continue.</p>
              <p class="hint">You can create an API key in PagerDuty under User Settings → API Access</p>
            </div>
            
            <div class="form-section">
              <div class="form-group">
                <label for="api-key">PagerDuty API Key</label>
                <input 
                  id="api-key"
                  type="password"
                  bind:value={apiKey}
                  placeholder="Enter your PagerDuty API key"
                  disabled={loading}
                />
              </div>
              
              <div class="form-group">
                <label for="theme">Theme</label>
                <select id="theme" bind:value={theme}>
                  <option value="dark">Dark</option>
                  <option value="light">Light</option>
                </select>
              </div>
              
              <div class="preferences">
                <h3>Initial Preferences</h3>
                
                <label class="checkbox-label">
                  <input 
                    type="checkbox" 
                    bind:checked={assignedOnly}
                  />
                  <span>Show only incidents assigned to me</span>
                </label>
                
                <label class="checkbox-label">
                  <input 
                    type="checkbox" 
                    bind:checked={redirectEnabled}
                  />
                  <span>Open new incidents in browser automatically</span>
                </label>
                
                <label class="checkbox-label">
                  <input 
                    type="checkbox" 
                    bind:checked={soundEnabled}
                  />
                  <span>Play sound for new incidents</span>
                </label>
              </div>
              
              {#if error}
                <div class="error-message">
                  ⚠️ {error}
                </div>
              {/if}
            </div>
            
            <div class="modal-actions">
              <button 
                class="btn-primary"
                on:click={validateApiKey}
                disabled={loading || !apiKey}
              >
                {loading ? 'Connecting...' : 'Connect to PagerDuty'}
              </button>
            </div>
          </div>
        {/if}
        
        {#if step === 2}
          <div class="setup-step">
            <div class="intro">
              <p>Select the services you want to monitor for incidents.</p>
              <p class="hint">You can change this later in Settings</p>
            </div>
            
            <div class="services-selection">
              {#if availableServices.length === 0}
                <p class="no-services">No services found in your PagerDuty account</p>
              {:else}
                <div class="service-list">
                  {#each availableServices as service}
                    <div 
                      class="service-option"
                      class:selected={isServiceSelected(service.id)}
                      on:click={() => toggleService(service)}
                    >
                      <div class="service-checkbox">
                        {#if isServiceSelected(service.id)}
                          ✓
                        {/if}
                      </div>
                      <div class="service-info">
                        <div class="service-name">{service.name}</div>
                        {#if service.description}
                          <div class="service-description">{service.description}</div>
                        {/if}
                      </div>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
            
            {#if error}
              <div class="error-message">
                ⚠️ {error}
              </div>
            {/if}
            
            <div class="modal-actions">
              <button 
                class="btn-secondary"
                on:click={skipServices}
              >
                Skip for Now
              </button>
              <button 
                class="btn-primary"
                on:click={completeSetup}
                disabled={loading || selectedServices.length === 0}
              >
                {loading ? 'Saving...' : `Monitor ${selectedServices.length} Service${selectedServices.length !== 1 ? 's' : ''}`}
              </button>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
  
  <style>
    .modal-overlay {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.7);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      animation: fadeIn 0.3s;
    }
    
    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    
    .modal {
      background: var(--bg-secondary);
      border-radius: 12px;
      width: 90%;
      max-width: 600px;
      max-height: 80vh;
      display: flex;
      flex-direction: column;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
      animation: slideIn 0.3s;
    }
    
    @keyframes slideIn {
      from { transform: translateY(-20px); opacity: 0; }
      to { transform: translateY(0); opacity: 1; }
    }
    
    .modal-header {
      padding: 2rem;
      background: var(--bg-tertiary);
      border-radius: 12px 12px 0 0;
      border-bottom: 2px solid var(--border-color);
    }
    
    .modal-header h2 {
      margin: 0 0 1.5rem 0;
      font-size: 1.5rem;
      color: var(--text-primary);
    }
    
    .steps {
      display: flex;
      gap: 2rem;
    }
    
    .step {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      opacity: 0.5;
      transition: opacity 0.3s;
    }
    
    .step.active {
      opacity: 1;
    }
    
    .step-number {
      width: 28px;
      height: 28px;
      border-radius: 50%;
      background: var(--bg-primary);
      color: var(--text-secondary);
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      font-size: 0.875rem;
    }
    
    .step.active .step-number {
      background: var(--accent-blue);
      color: white;
    }
    
    .step-label {
      font-size: 0.875rem;
      color: var(--text-secondary);
    }
    
    .step.active .step-label {
      color: var(--text-primary);
      font-weight: 600;
    }
    
    .modal-content {
      flex: 1;
      overflow-y: auto;
      padding: 2rem;
    }
    
    .setup-step {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
    
    .intro {
      text-align: center;
    }
    
    .intro p {
      margin: 0.5rem 0;
      color: var(--text-primary);
    }
    
    .hint {
      font-size: 0.875rem;
      color: var(--text-secondary);
    }
    
    .form-section {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
    
    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .form-group label {
      font-weight: 600;
      color: var(--text-secondary);
      font-size: 0.875rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    
    .form-group input,
    .form-group select {
      padding: 0.75rem;
      background: var(--bg-primary);
      border: 2px solid var(--border-color);
      border-radius: 8px;
      color: var(--text-primary);
      font-size: 1rem;
      font-family: inherit;
    }
    
    .form-group input:focus,
    .form-group select:focus {
      outline: none;
      border-color: var(--accent-blue);
    }
    
    .form-group input:disabled {
      opacity: 0.5;
    }
    
    .preferences {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }
    
    .preferences h3 {
      margin: 0;
      font-size: 1rem;
      color: var(--text-secondary);
    }
    
    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      cursor: pointer;
      padding: 0.5rem;
      border-radius: 4px;
      transition: background 0.2s;
    }
    
    .checkbox-label:hover {
      background: var(--bg-primary);
    }
    
    .checkbox-label input {
      width: 20px;
      height: 20px;
      cursor: pointer;
    }
    
    .checkbox-label span {
      color: var(--text-primary);
      font-size: 0.875rem;
    }
    
    .services-selection {
      max-height: 300px;
      overflow-y: auto;
    }
    
    .no-services {
      text-align: center;
      color: var(--text-tertiary);
      padding: 2rem;
    }
    
    .service-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .service-option {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem;
      background: var(--bg-primary);
      border: 2px solid var(--border-color);
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.2s;
    }
    
    .service-option:hover {
      background: var(--bg-tertiary);
      transform: translateX(4px);
    }
    
    .service-option.selected {
      border-color: var(--accent-blue);
      background: rgba(66, 153, 225, 0.1);
    }
    
    .service-checkbox {
      width: 24px;
      height: 24px;
      border: 2px solid var(--border-color);
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: bold;
      color: white;
      transition: all 0.2s;
    }
    
    .service-option.selected .service-checkbox {
      background: var(--accent-blue);
      border-color: var(--accent-blue);
    }
    
    .service-info {
      flex: 1;
    }
    
    .service-name {
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .service-description {
      font-size: 0.875rem;
      color: var(--text-secondary);
      margin-top: 0.25rem;
    }
    
    .error-message {
      padding: 0.75rem;
      background: rgba(252, 129, 129, 0.1);
      border: 1px solid var(--accent-red);
      border-radius: 4px;
      color: var(--accent-red);
      font-size: 0.875rem;
    }
    
    .modal-actions {
      display: flex;
      gap: 1rem;
      justify-content: flex-end;
      margin-top: 1rem;
    }
    
    .modal-actions button {
      min-width: 120px;
    }
    
    .modal-actions button:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  </style>