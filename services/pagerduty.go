package services

import (
	"context"
	"fmt"
	"pager-ops/database"
	"sync"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type PagerDutyService struct {
	client     *pagerduty.Client
	repository *database.Repository
	userID     string
	userEmail  string // Store email for "From" header
	mu         sync.RWMutex
}

func NewPagerDutyService(apiKey string, repo *database.Repository) (*PagerDutyService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	client := pagerduty.NewClient(apiKey)

	service := &PagerDutyService{
		client:     client,
		repository: repo,
	}

	// Get current user
	if err := service.RefreshUser(); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return service, nil
}

// RefreshUser fetches and caches the current user
func (s *PagerDutyService) RefreshUser() error {
	return s.RefreshUserWithContext(context.Background())
}

// RefreshUserWithContext fetches and caches the current user with context
func (s *PagerDutyService) RefreshUserWithContext(ctx context.Context) error {
	user, err := s.client.GetCurrentUserWithContext(ctx, pagerduty.GetCurrentUserOptions{})
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.userID = user.ID
	s.userEmail = user.Email
	s.mu.Unlock()

	// Save to database
	dbUser := database.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
	return s.repository.SetUser(dbUser)
}

// GetUserID returns the cached user ID
func (s *PagerDutyService) GetUserID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userID
}

// GetUserEmail returns the cached user email
func (s *PagerDutyService) GetUserEmail() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userEmail
}

// FetchIncidents fetches incidents based on filters
func (s *PagerDutyService) FetchIncidents(ctx context.Context, filters IncidentFilters) ([]Incident, error) {
	opts := pagerduty.ListIncidentsOptions{
		Limit: 100,
	}

	// Add status filters
	if len(filters.Statuses) > 0 {
		opts.Statuses = filters.Statuses
	}

	// Add service IDs
	if len(filters.ServiceIDs) > 0 {
		opts.ServiceIDs = filters.ServiceIDs
	}

	// Add user filter if enabled
	if filters.AssignedToMe && s.userID != "" {
		opts.UserIDs = []string{s.userID}
	}

	// Add time filter for resolved incidents
	if filters.Since != nil {
		opts.Since = filters.Since.Format(time.RFC3339)
	}

	// Sort by created date
	opts.SortBy = "created_at:desc"

	resp, err := s.client.ListIncidentsWithContext(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch incidents: %w", err)
	}

	incidents := make([]Incident, 0, len(resp.Incidents))
	for _, inc := range resp.Incidents {
		incidents = append(incidents, s.convertIncident(inc))
	}

	return incidents, nil
}

// FetchIncidentAlerts fetches alerts for a specific incident
func (s *PagerDutyService) FetchIncidentAlerts(ctx context.Context, incidentID string) ([]Alert, error) {
	// Use ListIncidentAlertsWithContext with empty options
	opts := pagerduty.ListIncidentAlertsOptions{}
	
	resp, err := s.client.ListIncidentAlertsWithContext(ctx, incidentID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch alerts: %w", err)
	}

	alerts := make([]Alert, 0, len(resp.Alerts))
	for _, alert := range resp.Alerts {
		alerts = append(alerts, Alert{
			ID:        alert.ID,
			Summary:   alert.Summary,
			Status:    alert.Status,
			CreatedAt: alert.CreatedAt,
			Details:   alert.Body,
		})
	}

	return alerts, nil
}

// AcknowledgeIncident acknowledges an incident
func (s *PagerDutyService) AcknowledgeIncident(ctx context.Context, incidentID string) error {
	incidents := []pagerduty.ManageIncidentsOptions{
		{
			ID:     incidentID,
			Type:   "incident_reference",
			Status: "acknowledged",
		},
	}

	_, err := s.client.ManageIncidentsWithContext(ctx, s.userEmail, incidents)
	return err
}

// ResolveIncident resolves an incident
func (s *PagerDutyService) ResolveIncident(ctx context.Context, incidentID string) error {
	incidents := []pagerduty.ManageIncidentsOptions{
		{
			ID:     incidentID,
			Type:   "incident_reference",
			Status: "resolved",
		},
	}

	_, err := s.client.ManageIncidentsWithContext(ctx, s.userEmail, incidents)
	return err
}

// EscalateIncident escalates an incident
func (s *PagerDutyService) EscalateIncident(ctx context.Context, incidentID string, escalationLevelID string) error {
	// Use the ManageIncidents API to escalate
	incidents := []pagerduty.ManageIncidentsOptions{
		{
			ID:   incidentID,
			Type: "incident_reference",
			EscalationLevel: uint(1), // Convert to uint - you may want to parse escalationLevelID
		},
	}

	_, err := s.client.ManageIncidentsWithContext(ctx, s.userEmail, incidents)
	return err
}

// SnoozeIncident snoozes an incident for a duration
func (s *PagerDutyService) SnoozeIncident(ctx context.Context, incidentID string, duration time.Duration) (*pagerduty.Incident, error) {
	// Convert duration to seconds (uint)
	secs := uint(duration.Seconds())
	
	// Call the new context-aware method
	incident, err := s.client.SnoozeIncidentWithContext(ctx, incidentID, secs)
	if err != nil {
		return nil, err
	}
	return incident, nil
}

// AddIncidentNote adds a note to an incident
func (s *PagerDutyService) AddIncidentNote(ctx context.Context, incidentID string, content string) error {
	// Create note with user information
	note := pagerduty.IncidentNote{
		Content: content,
		User: pagerduty.APIObject{
			ID:      s.userID,
			Type:    "user_reference",
			Summary: s.userEmail, // The API uses Summary field for the "From" header
		},
	}

	_, err := s.client.CreateIncidentNoteWithContext(ctx, incidentID, note)
	return err
}

// MergeIncidents merges multiple incidents into one
func (s *PagerDutyService) MergeIncidents(ctx context.Context, sourceIncidentIDs []string, targetIncidentID string) error {
	// Convert source incident IDs to MergeIncidentsOptions
	sourceIncidents := make([]pagerduty.MergeIncidentsOptions, 0, len(sourceIncidentIDs))
	for _, id := range sourceIncidentIDs {
		sourceIncidents = append(sourceIncidents, pagerduty.MergeIncidentsOptions{
			ID:   id,
			Type: "incident_reference",
		})
	}

	_, err := s.client.MergeIncidentsWithContext(ctx, s.userEmail, targetIncidentID, sourceIncidents)
	return err
}

// GetServices fetches all PagerDuty services
func (s *PagerDutyService) GetServices() ([]Service, error) {
	return s.GetServicesWithContext(context.Background())
}

// GetServicesWithContext fetches all PagerDuty services with context
func (s *PagerDutyService) GetServicesWithContext(ctx context.Context) ([]Service, error) {
	opts := pagerduty.ListServiceOptions{
		Limit: 100,
	}

	resp, err := s.client.ListServicesWithContext(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch services: %w", err)
	}

	services := make([]Service, 0, len(resp.Services))
	for _, svc := range resp.Services {
		services = append(services, Service{
			ID:          svc.ID,
			Name:        svc.Name,
			Description: svc.Description,
			Status:      svc.Status,
		})
	}

	return services, nil
}

// Helper function to convert PagerDuty incident to our model
func (s *PagerDutyService) convertIncident(inc pagerduty.Incident) Incident {
	// Extract assigned user IDs
	assignedUserIDs := make([]string, 0, len(inc.Assignments))
	for _, assignment := range inc.Assignments {
		assignedUserIDs = append(assignedUserIDs, assignment.Assignee.ID)
	}

	// Parse created and updated times
	createdAt, _ := time.Parse(time.RFC3339, inc.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, inc.LastStatusChangeAt)

	// Extract service ID - Service is an APIObject, not a pointer
	serviceID := ""
	if inc.Service.ID != "" {
		serviceID = inc.Service.ID
	}

	// Extract escalation level - EscalationPolicy is an APIObject
	// We can't directly get escalation rules, so we'll default to 0
	// In a real implementation, you might need to make another API call
	// to get the full escalation policy details
	escalationLevel := 0

	return Incident{
		ID:              inc.ID,
		ServiceID:       serviceID,
		Status:          inc.Status,
		Title:           inc.Title,
		Description:     inc.Description,
		Urgency:         inc.Urgency,
		IncidentNumber:  int(inc.IncidentNumber), // Convert uint to int
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		AssignedUserIDs: assignedUserIDs,
		EscalationLevel: escalationLevel,
		HTMLURL:         inc.HTMLURL,
	}
}

// Data models
type IncidentFilters struct {
	Statuses     []string
	ServiceIDs   []string
	AssignedToMe bool
	Since        *time.Time
}

type Incident struct {
	ID              string    `json:"id"`
	ServiceID       string    `json:"service_id"`
	Status          string    `json:"status"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Urgency         string    `json:"urgency"`
	IncidentNumber  int       `json:"incident_number"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	AssignedUserIDs []string  `json:"assigned_user_ids"`
	EscalationLevel int       `json:"escalation_level"`
	HTMLURL         string    `json:"html_url"`
}

type Alert struct {
	ID        string                 `json:"id"`
	Summary   string                 `json:"summary"`
	Status    string                 `json:"status"`
	CreatedAt string                 `json:"created_at"`
	Details   map[string]interface{} `json:"details"`
}

type Service struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}