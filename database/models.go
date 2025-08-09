package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Settings model
type Setting struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

// Service model
type Service struct {
	ID      string `db:"id" json:"id"`
	Alias   string `db:"alias" json:"alias"`
	Enabled bool   `db:"enabled" json:"enabled"`
}

// User model
type User struct {
	ID    string `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
	Name  string `db:"name" json:"name"`
}

// Incident model for caching
type CachedIncident struct {
	ID              string    `db:"id" json:"id"`
	ServiceID       string    `db:"service_id" json:"service_id"`
	Status          string    `db:"status" json:"status"`
	Title           string    `db:"title" json:"title"`
	Description     string    `db:"description" json:"description"`
	Urgency         string    `db:"urgency" json:"urgency"`
	IncidentNumber  int       `db:"incident_number" json:"incident_number"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	Pinned          bool      `db:"pinned" json:"pinned"`
	AssignedUserIDs JSONArray `db:"assigned_user_ids" json:"assigned_user_ids"`
	EscalationLevel int       `db:"escalation_level" json:"escalation_level"`
	HTMLURL         string    `db:"html_url" json:"html_url"`
}

// DraftNote model
type DraftNote struct {
	IncidentID  string    `db:"incident_id" json:"incident_id"`
	NoteText    string    `db:"note_text" json:"note_text"`
	WhyTriggered string   `db:"why_triggered" json:"why_triggered"`
	Impact      string    `db:"impact" json:"impact"`
	Actions     string    `db:"actions" json:"actions"`
	Links       string    `db:"links" json:"links"`
	LastUpdated time.Time `db:"last_updated" json:"last_updated"`
}

// Template model
type Template struct {
	ID       int    `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	BodyText string `db:"body_text" json:"body_text"`
}

// JSONArray handles JSON array storage in SQLite
type JSONArray []string

func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return fmt.Errorf("cannot scan type %T into JSONArray", value)
	}
}

// Settings keys
const (
	SettingAPIKey          = "api_key"
	SettingTheme           = "theme"
	SettingAssignedOnly    = "assigned_only"
	SettingRedirectEnabled = "redirect_enabled"
	SettingSoundPath       = "sound_path"
	SettingSoundEnabled    = "sound_enabled"
)