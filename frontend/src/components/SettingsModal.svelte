<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { 
      SaveSettings, 
      GetServices,
      SaveService,
      DeleteService,
      FetchAvailableServices,
      SelectSoundFolder,
      GetTemplates,
      SaveTemplate,
      DeleteTemplate
    } from '../../wailsjs/go/main/App.js';
    
    export let settings = null;
    
    const dispatch = createEventDispatcher();
    
    let activeTab = 'general';
    let loading = false;
    let error = '';
    let successMessage = '';
    
    // General settings
    let apiKey = settings?.api_key || '';
    let theme = settings?.theme || 'dark';
    let assignedOnly = settings?.assigned_only || false;
    let redirectEnabled = settings?.redirect_enabled || false;
    let soundEnabled = settings?.sound_enabled || false;
    let soundPath = settings?.sound_path || '';
    
    // Services
    let services = [];
    let availableServices = [];
    let showServicePicker = false;
    let editingService = null;
    
    // Templates
    let templates = [];
    let editingTemplate = null;
    let newTemplate = { title: '', body_text: '' };
    
    loadServices();
    loadTemplates();
    
    async function loadServices() {
      try {
        services = await GetServices();
      } catch (err) {
        console.error('Failed to load services:', err);
      }
    }
    
    async function loadTemplates() {
      try {
        templates = await GetTemplates();
      } catch (err) {
        console.error('Failed to load templates:', err);
      }
    }
    
    async function saveGeneralSettings() {
      loading = true;
      error = '';
      successMessage = '';
      
      try {
        await SaveSettings({
          api_key: apiKey,
          theme: theme,
          assigned_only: assignedOnly,
          redirect_enabled: redirectEnabled,
          sound_enabled: soundEnabled,
          sound_path: soundPath
        });
        
        successMessage = 'Settings saved successfully';
        setTimeout(() => successMessage = '', 3000);
      } catch (err) {
        error = 'Failed to save settings';
      } finally {
        loading = false;
      }
    }
    
    async function selectSoundFolder() {
      const folder = await SelectSoundFolder();
      if (folder) {
        soundPath = folder;
      }
    }
    
    async function openServicePicker() {
      try {
        availableServices = await FetchAvailableServices();
        showServicePicker = true;
      } catch (err) {
        error = 'Failed to fetch available services';
      }
    }
    
    async function addService(pdService) {
      const service = {
        id: pdService.id,
        alias: pdService.name,
        enabled: true
      };
      
      try {
        await SaveService(service);
        await loadServices();
        showServicePicker = false;
      } catch (err) {
        error = 'Failed to add service';
      }
    }
    
    function editService(service) {
      editingService = { ...service };
    }
    
    async function saveServiceEdit() {
      if (!editingService) return;
      
      try {
        await SaveService(editingService);
        await loadServices();
        editingService = null;
      } catch (err) {
        error = 'Failed to save service';
      }
    }
    
    function cancelServiceEdit() {
      editingService = null;
    }
    
    async function removeService(id) {
      if (confirm('Are you sure you want to remove this service?')) {
        try {
          await DeleteService(id);
          await loadServices();
        } catch (err) {
          error = 'Failed to remove service';
        }
      }
    }
    
    function editTemplate(template) {
      editingTemplate = { ...template };
    }
    
    async function saveTemplateEdit() {
      if (!editingTemplate) return;
      
      try {
        await SaveTemplate(editingTemplate);
        await loadTemplates();
        editingTemplate = null;
      } catch (err) {
        error = 'Failed to save template';
      }
    }
    
    function cancelTemplateEdit() {
      editingTemplate = null;
    }
    
    async function saveNewTemplate() {
      if (!newTemplate.title || !newTemplate.body_text) return;
      
      try {
        await SaveTemplate(newTemplate);
        await loadTemplates();
        newTemplate = { title: '', body_text: '' };
      } catch (err) {
        error = 'Failed to save template';
      }
    }
    
    async function removeTemplate(id) {
      if (confirm('Are you sure you want to delete this template?')) {
        try {
          await DeleteTemplate(id);
          await loadTemplates();
        } catch (err) {
          error = 'Failed to delete template';
        }
      }
    }
    
    function close() {
      dispatch('close');
    }
  </script>
  
  <div class="modal-overlay" on:click={close}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Settings</h2>
        <button class="close-btn" on:click={close}>✕</button>
      </div>
      
      <div class="modal-tabs">
        <button 
          class="tab"
          class:active={activeTab === 'general'}
          on:click={() => activeTab = 'general'}
        >
          General
        </button>
        <button 
          class="tab"
          class:active={activeTab === 'services'}
          on:click={() => activeTab = 'services'}
        >
          Services
        </button>
        <button 
          class="tab"
          class:active={activeTab === 'templates'}
          on:click={() => activeTab = 'templates'}
        >
          Templates
        </button>
      </div>
      
      <div class="modal-content">
        {#if activeTab === 'general'}
          <div class="settings-section">
            <h3>API Configuration</h3>
            <div class="form-group">
              <label for="api-key">PagerDuty API Key</label>
              <input 
                id="api-key"
                type="password"
                bind:value={apiKey}
                placeholder="Enter your PagerDuty API key"
              />
              <p class="help-text">Leave blank to keep current key</p>
            </div>
            
            <h3>Appearance</h3>
            <div class="form-group">
              <label for="theme">Theme</label>
              <select id="theme" bind:value={theme}>
                <option value="dark">Dark</option>
                <option value="light">Light</option>
              </select>
            </div>
            
            <h3>Incident Preferences</h3>
            <div class="checkbox-group">
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
            </div>
            
            <h3>Sound Alerts</h3>
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input 
                  type="checkbox" 
                  bind:checked={soundEnabled}
                />
                <span>Play sound for new incidents</span>
              </label>
            </div>
            
            {#if soundEnabled}
              <div class="form-group">
                <label>Alert Sound Folder</label>
                <div class="folder-picker">
                  <input 
                    type="text"
                    bind:value={soundPath}
                    placeholder="Select folder containing sound files"
                    readonly
                  />
                  <button class="btn-secondary" on:click={selectSoundFolder}>
                    Browse
                  </button>
                </div>
                <p class="help-text">Select a folder containing MP3 or WAV files</p>
              </div>
            {/if}
            
            {#if error}
              <div class="error-message">⚠️ {error}</div>
            {/if}
            
            {#if successMessage}
              <div class="success-message">✅ {successMessage}</div>
            {/if}
            
            <div class="section-actions">
              <button 
                class="btn-primary"
                on:click={saveGeneralSettings}
                disabled={loading}
              >
                {loading ? 'Saving...' : 'Save Settings'}
              </button>
            </div>
          </div>
        {/if}
        
        {#if activeTab === 'services'}
          <div class="settings-section">
            <div class="section-header">
              <h3>Monitored Services</h3>
              <button class="btn-primary" on:click={openServicePicker}>
                + Add Service
              </button>
            </div>
            
            {#if services.length === 0}
              <div class="empty-state">
                <p>No services configured</p>
                <p class="hint">Click "Add Service" to get started</p>
              </div>
            {:else}
              <div class="service-list">
                {#each services as service}
                  <div class="service-item">
                    {#if editingService?.id === service.id}
                      <div class="service-edit">
                        <input 
                          type="text"
                          bind:value={editingService.alias}
                          placeholder="Service alias"
                          class="alias-input"
                        />
                        <div class="edit-actions">
                          <button class="btn-success" on:click={saveServiceEdit}>
                            Save
                          </button>
                          <button class="btn-secondary" on:click={cancelServiceEdit}>
                            Cancel
                          </button>
                        </div>
                      </div>
                    {:else}
                      <div class="service-info">
                        <div class="service-name">{service.alias}</div>
                        <div class="service-id">{service.id}</div>
                      </div>
                      <div class="service-actions">
                        <button 
                          class="btn-icon"
                          on:click={() => editService(service)}
                          title="Edit"
                        >
                          ✏️
                        </button>
                        <button 
                          class="btn-icon"
                          on:click={() => removeService(service.id)}
                          title="Remove"
                        >
                          🗑️
                        </button>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
            
            {#if showServicePicker}
              <div class="service-picker">
                <h4>Available Services</h4>
                <div class="available-services">
                  {#each availableServices as pdService}
                    {#if !services.find(s => s.id === pdService.id)}
                      <div 
                        class="available-service"
                        on:click={() => addService(pdService)}
                      >
                        <div class="service-info">
                          <div class="service-name">{pdService.name}</div>
                          {#if pdService.description}
                            <div class="service-description">{pdService.description}</div>
                          {/if}
                        </div>
                        <button class="btn-add">+</button>
                      </div>
                    {/if}
                  {/each}
                </div>
                <button 
                  class="btn-secondary"
                  on:click={() => showServicePicker = false}
                >
                  Cancel
                </button>
              </div>
            {/if}
          </div>
        {/if}
        
        {#if activeTab === 'templates'}
          <div class="settings-section">
            <h3>Note Templates</h3>
            
            <div class="template-form">
              <h4>Create New Template</h4>
              <input 
                type="text"
                bind:value={newTemplate.title}
                placeholder="Template name"
                class="template-input"
              />
              <textarea 
                bind:value={newTemplate.body_text}
                placeholder="Template content..."
                rows="4"
                class="template-textarea"
              ></textarea>
              <button 
                class="btn-primary"
                on:click={saveNewTemplate}
                disabled={!newTemplate.title || !newTemplate.body_text}
              >
                Save Template
              </button>
            </div>
            
            {#if templates.length > 0}
              <div class="template-list">
                <h4>Saved Templates</h4>
                {#each templates as template}
                  <div class="template-item">
                    {#if editingTemplate?.id === template.id}
                      <div class="template-edit">
                        <input 
                          type="text"
                          bind:value={editingTemplate.title}
                          class="template-input"
                        />
                        <textarea 
                          bind:value={editingTemplate.body_text}
                          rows="3"
                          class="template-textarea"
                        ></textarea>
                        <div class="edit-actions">
                          <button class="btn-success" on:click={saveTemplateEdit}>
                            Save
                          </button>
                          <button class="btn-secondary" on:click={cancelTemplateEdit}>
                            Cancel
                          </button>
                        </div>
                      </div>
                    {:else}
                      <div class="template-info">
                        <div class="template-title">{template.title}</div>
                        <div class="template-preview">{template.body_text}</div>
                      </div>
                      <div class="template-actions">
                        <button 
                          class="btn-icon"
                          on:click={() => editTemplate(template)}
                          title="Edit"
                        >
                          ✏️
                        </button>
                        <button 
                          class="btn-icon"
                          on:click={() => removeTemplate(template.id)}
                          title="Delete"
                        >
                          🗑️
                        </button>
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
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
    }
    
    .modal {
      background: var(--bg-secondary);
      border-radius: 12px;
      width: 90%;
      max-width: 700px;
      max-height: 80vh;
      display: flex;
      flex-direction: column;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    }
    
    .modal-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1.5rem;
      background: var(--bg-tertiary);
      border-radius: 12px 12px 0 0;
      border-bottom: 2px solid var(--border-color);
    }
    
    .modal-header h2 {
      margin: 0;
      font-size: 1.25rem;
      color: var(--text-primary);
    }
    
    .close-btn {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      background: transparent;
      color: var(--text-secondary);
      font-size: 1.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    
    .close-btn:hover {
      background: var(--bg-primary);
      transform: none;
    }
    
    .modal-tabs {
      display: flex;
      gap: 1rem;
      padding: 0 1.5rem;
      background: var(--bg-tertiary);
      border-bottom: 1px solid var(--border-color);
    }
    
    .tab {
      padding: 0.75rem 1rem;
      background: transparent;
      color: var(--text-secondary);
      border-bottom: 2px solid transparent;
      transition: all 0.2s;
    }
    
    .tab:hover {
      color: var(--text-primary);
      transform: none;
    }
    
    .tab.active {
      color: var(--accent-blue);
      border-bottom-color: var(--accent-blue);
    }
    
    .modal-content {
      flex: 1;
      overflow-y: auto;
      padding: 1.5rem;
    }
    
    .settings-section {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
    
    .settings-section h3 {
      margin: 0;
      font-size: 1rem;
      color: var(--text-secondary);
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    
    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    
    .form-group label {
      font-weight: 500;
      color: var(--text-secondary);
      font-size: 0.875rem;
    }
    
    .form-group input,
    .form-group select,
    .form-group textarea {
      padding: 0.5rem;
      background: var(--bg-primary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-primary);
      font-family: inherit;
      font-size: 0.875rem;
    }
    
    .form-group input:focus,
    .form-group select:focus,
    .form-group textarea:focus {
      outline: none;
      border-color: var(--accent-blue);
    }
    
    .help-text {
      font-size: 0.75rem;
      color: var(--text-tertiary);
      margin: 0;
    }
    
    .checkbox-group {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }
    
    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      cursor: pointer;
    }
    
    .checkbox-label input {
      width: 18px;
      height: 18px;
      cursor: pointer;
    }
    
    .checkbox-label span {
      color: var(--text-primary);
      font-size: 0.875rem;
    }
    
    .folder-picker {
      display: flex;
      gap: 0.5rem;
    }
    
    .folder-picker input {
      flex: 1;
    }
    
    .error-message,
    .success-message {
      padding: 0.75rem;
      border-radius: 4px;
      font-size: 0.875rem;
    }
    
    .error-message {
      background: rgba(252, 129, 129, 0.1);
      border: 1px solid var(--accent-red);
      color: var(--accent-red);
    }
    
    .success-message {
      background: rgba(72, 187, 120, 0.1);
      border: 1px solid var(--accent-green);
      color: var(--accent-green);
    }
    
    .section-actions {
      display: flex;
      justify-content: flex-end;
      gap: 1rem;
    }
    
    .empty-state {
      text-align: center;
      padding: 2rem;
      color: var(--text-tertiary);
    }
    
    .empty-state .hint {
      font-size: 0.875rem;
      margin-top: 0.5rem;
    }
    
    .service-list,
    .template-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }
    
    .service-item,
    .template-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem;
      background: var(--bg-primary);
      border-radius: 4px;
    }
    
    .service-info,
    .template-info {
      flex: 1;
    }
    
    .service-name,
    .template-title {
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .service-id,
    .service-description {
      font-size: 0.75rem;
      color: var(--text-tertiary);
      margin-top: 0.25rem;
    }
    
    .template-preview {
      font-size: 0.875rem;
      color: var(--text-secondary);
      margin-top: 0.25rem;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .service-actions,
    .template-actions {
      display: flex;
      gap: 0.5rem;
    }
    
    .btn-icon {
      width: 32px;
      height: 32px;
      padding: 0;
      background: transparent;
      font-size: 1rem;
    }
    
    .btn-icon:hover {
      background: var(--bg-tertiary);
      transform: none;
    }
    
    .service-edit,
    .template-edit {
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }
    
    .alias-input,
    .template-input {
      padding: 0.5rem;
      background: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-primary);
      font-family: inherit;
    }
    
    .template-textarea {
      padding: 0.5rem;
      background: var(--bg-secondary);
      border: 1px solid var(--border-color);
      border-radius: 4px;
      color: var(--text-primary);
      font-family: inherit;
      resize: vertical;
    }
    
    .edit-actions {
      display: flex;
      gap: 0.5rem;
    }
    
    .service-picker {
      margin-top: 1rem;
      padding: 1rem;
      background: var(--bg-primary);
      border-radius: 4px;
    }
    
    .service-picker h4,
    .template-form h4,
    .template-list h4 {
      margin: 0 0 1rem 0;
      font-size: 0.875rem;
      color: var(--text-secondary);
    }
    
    .available-services {
      max-height: 200px;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      margin-bottom: 1rem;
    }
    
    .available-service {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.75rem;
      background: var(--bg-secondary);
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.2s;
    }
    
    .available-service:hover {
      background: var(--bg-tertiary);
      transform: translateX(4px);
    }
    
    .btn-add {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      background: var(--accent-blue);
      color: white;
      font-size: 1.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    
    .template-form {
      padding: 1rem;
      background: var(--bg-primary);
      border-radius: 4px;
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }
  </style>