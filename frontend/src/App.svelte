<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  
  let appTitle: string = "PagerDuty Incident Command Center"
</script>

<main>
  <div class="header">
    <img alt="PagerDuty Command Center" id="logo" src="{logo}">
    <h1>{appTitle}</h1>
  </div>
  
  <div class="container">
    <!-- Three-column layout for PagerDuty interface -->
    <div class="columns">
      <!-- Column 1: Services Panel -->
      <div class="column services-panel">
        <h2>Services</h2>
        <div class="panel-content">
          <!-- Service list will go here -->
          <p>Service health overview coming soon...</p>
        </div>
      </div>
      
      <!-- Column 2: Incidents Panel -->
      <div class="column incidents-panel">
        <h2>Active Incidents</h2>
        <div class="panel-content">
          <!-- Incident list will go here -->
          <p>Loading incidents...</p>
        </div>
      </div>
      
      <!-- Column 3: Alert Details Panel -->
      <div class="column details-panel">
        <h2>Alert Details</h2>
        <div class="panel-content">
          <!-- Alert details will go here -->
          <p>Select an incident to view details</p>
        </div>
      </div>
    </div>
  </div>
</main>

<style>
  .header {
    padding: 20px;
    text-align: center;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  #logo {
    display: block;
    width: 60px;
    height: 60px;
    margin: 0 auto 10px;
  }

  h1 {
    margin: 0;
    font-size: 1.5rem;
    color: white;
  }

  .container {
    height: calc(100vh - 120px);
    padding: 20px;
  }

  .columns {
    display: flex;
    gap: 20px;
    height: 100%;
  }

  .column {
    flex: 1;
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    padding: 15px;
    display: flex;
    flex-direction: column;
  }

  .column h2 {
    margin: 0 0 15px 0;
    font-size: 1.2rem;
    color: white;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    padding-bottom: 10px;
  }

  .panel-content {
    flex: 1;
    overflow-y: auto;
    color: rgba(255, 255, 255, 0.8);
  }

  .services-panel {
    flex: 0.8;
  }

  .incidents-panel {
    flex: 1.2;
  }

  .details-panel {
    flex: 1;
  }

  /* Scrollbar styling */
  .panel-content::-webkit-scrollbar {
    width: 8px;
  }

  .panel-content::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
  }

  .panel-content::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 4px;
  }

  .panel-content::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.3);
  }
</style>