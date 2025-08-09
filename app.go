package main

import (
	"context"
	"fmt"
	"pager-ops/database"
	"pager-ops/services"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	repository *database.Repository
	pd         *services.PagerDutyService
	poller     *services.IncidentPoller
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	repo, err := database.NewRepository("pager-ops.db")
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Failed to initialize database: %v", err))
		return
	}
	a.repository = repo

	// Try to initialize PagerDuty service
	a.initializePagerDuty()
}

func (a *App) initializePagerDuty() {
	apiKey, err := a.repository.GetAPIKey()
	if err != nil || apiKey == "" {
		runtime.LogInfo(a.ctx, "No API key configured")
		runtime.EventsEmit(a.ctx, "setup:required", true)
		return
	}

	pd, err := services.NewPagerDutyService(apiKey, a.repository)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to initialize PagerDuty: %v", err))
		runtime.EventsEmit(a.ctx, "setup:required", true)
		return
	}

	a.pd = pd

	// Start incident poller
	a.poller = services.NewIncidentPoller(a.ctx, pd, a.repository)
	a.poller.Start()

	runtime.EventsEmit(a.ctx, "setup:complete", true)
}

// Settings Management

type AppSettings struct {
	APIKey          string `json:"api_key"`
	Theme           string `json:"theme"`
	AssignedOnly    bool   `json:"assigned_only"`
	RedirectEnabled bool   `json:"redirect_enabled"`
	SoundPath       string `json:"sound_path"`
	SoundEnabled    bool   `json:"sound_enabled"`
}

func (a *App) GetSettings() (AppSettings, error) {
	apiKey, _ := a.repository.GetAPIKey()
	theme, _ := a.repository.GetSetting(database.SettingTheme)
	assignedOnly, _ := a.repository.GetSetting(database.SettingAssignedOnly)
	redirectEnabled, _ := a.repository.GetSetting(database.SettingRedirectEnabled)
	soundPath, _ := a.repository.GetSetting(database.SettingSoundPath)
	soundEnabled, _ := a.repository.GetSetting(database.SettingSoundEnabled)

	// Mask API key for security
	if apiKey != "" {
		apiKey = "****" + apiKey[len(apiKey)-4:]
	}

	return AppSettings{
		APIKey:          apiKey,
		Theme:           theme,
		AssignedOnly:    assignedOnly == "true",
		RedirectEnabled: redirectEnabled == "true",
		SoundPath:       soundPath,
		SoundEnabled:    soundEnabled == "true",
	}, nil
}

func (a *App) SaveSettings(settings AppSettings) error {
	// Only update API key if it's not masked
	if settings.APIKey != "" && !startsWith(settings.APIKey, "****") {
		a.repository.SetSetting(database.SettingAPIKey, settings.APIKey)

		// Reinitialize PagerDuty with new API key
		if a.poller != nil {
			a.poller.Stop()
		}
		a.initializePagerDuty()
	}

	a.repository.SetSetting(database.SettingTheme, settings.Theme)
	a.repository.SetSetting(database.SettingAssignedOnly, boolToString(settings.AssignedOnly))
	a.repository.SetSetting(database.SettingRedirectEnabled, boolToString(settings.RedirectEnabled))
	a.repository.SetSetting(database.SettingSoundPath, settings.SoundPath)
	a.repository.SetSetting(database.SettingSoundEnabled, boolToString(settings.SoundEnabled))

	runtime.EventsEmit(a.ctx, "settings:updated", settings)
	return nil
}

// Service Management

func (a *App) GetServices() ([]database.Service, error) {
	return a.repository.GetServices()
}

func (a *App) SaveService(service database.Service) error {
	return a.repository.UpsertService(service)
}

func (a *App) DeleteService(id string) error {
	return a.repository.DeleteService(id)
}

func (a *App) ToggleService(id string, enabled bool) error {
	err := a.repository.ToggleService(id, enabled)
	if err == nil {
		runtime.EventsEmit(a.ctx, "services:updated", nil)
	}
	return err
}

func (a *App) FetchAvailableServices() ([]services.Service, error) {
	if a.pd == nil {
		return nil, fmt.Errorf("PagerDuty not initialized")
	}
	return a.pd.GetServicesWithContext(a.ctx)
}

// Incident Management

func (a *App) GetIncidents(statuses []string) ([]database.CachedIncident, error) {
	return a.repository.GetIncidents(statuses)
}

func (a *App) AcknowledgeIncident(incidentID string) error {
	if a.pd == nil {
		return fmt.Errorf("PagerDuty not initialized")
	}

	err := a.pd.AcknowledgeIncident(a.ctx, incidentID)
	if err == nil {
		runtime.EventsEmit(a.ctx, "incident:acknowledged", incidentID)
	}
	return err
}

func (a *App) ResolveIncident(incidentID string) error {
	if a.pd == nil {
		return fmt.Errorf("PagerDuty not initialized")
	}

	err := a.pd.ResolveIncident(a.ctx, incidentID)
	if err == nil {
		runtime.EventsEmit(a.ctx, "incident:resolved", incidentID)
	}
	return err
}

func (a *App) EscalateIncident(incidentID string, escalationLevelID string) error {
	if a.pd == nil {
		return fmt.Errorf("PagerDuty not initialized")
	}

	return a.pd.EscalateIncident(a.ctx, incidentID, escalationLevelID)
}

func (a *App) SnoozeIncident(incidentID string, minutes int) error {
	if a.pd == nil {
		return fmt.Errorf("PagerDuty not initialized")
	}

	duration := time.Duration(minutes) * time.Minute
	_, err := a.pd.SnoozeIncident(a.ctx, incidentID, duration)
	if err != nil {
		return fmt.Errorf("failed to snooze incident: %w", err)
	}

	runtime.EventsEmit(a.ctx, "incident:snoozed", incidentID)
	return nil
}

func (a *App) PinIncident(incidentID string, pinned bool) error {
	return a.repository.PinIncident(incidentID, pinned)
}

func (a *App) MergeIncidents(sourceIDs []string, targetID string) error {
	if a.pd == nil {
		return fmt.Errorf("PagerDuty not initialized")
	}

	return a.pd.MergeIncidents(a.ctx, sourceIDs, targetID)
}

// Alert Management

func (a *App) GetIncidentAlerts(incidentID string) ([]services.Alert, error) {
	if a.pd == nil {
		return nil, fmt.Errorf("PagerDuty not initialized")
	}

	return a.pd.FetchIncidentAlerts(a.ctx, incidentID)
}

// Notes Management

func (a *App) GetDraftNote(incidentID string) (*database.DraftNote, error) {
	return a.repository.GetDraftNote(incidentID)
}

func (a *App) SaveDraftNote(note database.DraftNote) error {
	return a.repository.SaveDraftNote(note)
}

func (a *App) AddIncidentNote(incidentID string, content string) error {
	if a.pd == nil {
		return fmt.Errorf("PagerDuty not initialized")
	}

	return a.pd.AddIncidentNote(a.ctx, incidentID, content)
}

// Template Management

func (a *App) GetTemplates() ([]database.Template, error) {
	return a.repository.GetTemplates()
}

func (a *App) SaveTemplate(template database.Template) error {
	return a.repository.SaveTemplate(template)
}

func (a *App) DeleteTemplate(id int) error {
	return a.repository.DeleteTemplate(id)
}

// User Management

func (a *App) GetCurrentUser() (*database.User, error) {
	return a.repository.GetUser()
}

// UI Helper Functions

func (a *App) SelectSoundFolder() string {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Alert Sound Folder",
	})
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to select folder: %v", err))
		return ""
	}
	return folder
}

func (a *App) OpenIncidentInBrowser(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// Helper functions

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
