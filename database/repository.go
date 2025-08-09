package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repo := &Repository{db: db}
	if err := repo.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *Repository) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS services (
		id TEXT PRIMARY KEY,
		alias TEXT NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS incidents (
		id TEXT PRIMARY KEY,
		service_id TEXT NOT NULL,
		status TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		urgency TEXT,
		incident_number INTEGER,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		pinned INTEGER NOT NULL DEFAULT 0,
		assigned_user_ids TEXT DEFAULT '[]',
		escalation_level INTEGER DEFAULT 0,
		html_url TEXT,
		FOREIGN KEY (service_id) REFERENCES services(id)
	);

	CREATE TABLE IF NOT EXISTS draft_notes (
		incident_id TEXT PRIMARY KEY,
		note_text TEXT,
		why_triggered TEXT,
		impact TEXT,
		actions TEXT,
		links TEXT,
		last_updated DATETIME NOT NULL,
		FOREIGN KEY (incident_id) REFERENCES incidents(id)
	);

	CREATE TABLE IF NOT EXISTS templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		body_text TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_incidents_status ON incidents(status);
	CREATE INDEX IF NOT EXISTS idx_incidents_service ON incidents(service_id);
	CREATE INDEX IF NOT EXISTS idx_incidents_updated ON incidents(updated_at);
	`

	_, err := r.db.Exec(schema)
	return err
}

// Settings operations
func (r *Repository) GetSetting(key string) (string, error) {
	var value string
	err := r.db.Get(&value, "SELECT value FROM settings WHERE key = ?", key)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (r *Repository) SetSetting(key, value string) error {
	// Encrypt API key before storing
	if key == SettingAPIKey && value != "" {
		encrypted, err := encryptValue(value)
		if err != nil {
			return err
		}
		value = encrypted
	}

	_, err := r.db.Exec(`
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	return err
}

func (r *Repository) GetAPIKey() (string, error) {
	encrypted, err := r.GetSetting(SettingAPIKey)
	if err != nil || encrypted == "" {
		return "", err
	}
	return decryptValue(encrypted)
}

// Service operations
func (r *Repository) GetServices() ([]Service, error) {
	var services []Service
	err := r.db.Select(&services, "SELECT * FROM services ORDER BY alias")
	return services, err
}

func (r *Repository) UpsertService(service Service) error {
	_, err := r.db.Exec(`
		INSERT INTO services (id, alias, enabled) VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET alias = excluded.alias, enabled = excluded.enabled
	`, service.ID, service.Alias, service.Enabled)
	return err
}

func (r *Repository) DeleteService(id string) error {
	_, err := r.db.Exec("DELETE FROM services WHERE id = ?", id)
	return err
}

func (r *Repository) ToggleService(id string, enabled bool) error {
	_, err := r.db.Exec("UPDATE services SET enabled = ? WHERE id = ?", enabled, id)
	return err
}

// User operations
func (r *Repository) GetUser() (*User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM user LIMIT 1")
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) SetUser(user User) error {
	// Clear existing user
	r.db.Exec("DELETE FROM user")
	
	_, err := r.db.Exec(`
		INSERT INTO user (id, email, name) VALUES (?, ?, ?)
	`, user.ID, user.Email, user.Name)
	return err
}

// Incident operations
func (r *Repository) GetIncidents(status []string) ([]CachedIncident, error) {
	query, args, err := sqlx.In(`
		SELECT * FROM incidents 
		WHERE status IN (?)
		ORDER BY pinned DESC, created_at DESC
	`, status)
	if err != nil {
		return nil, err
	}

	var incidents []CachedIncident
	err = r.db.Select(&incidents, r.db.Rebind(query), args...)
	return incidents, err
}

func (r *Repository) UpsertIncident(incident CachedIncident) error {
	_, err := r.db.Exec(`
		INSERT INTO incidents (
			id, service_id, status, title, description, urgency, 
			incident_number, created_at, updated_at, pinned, 
			assigned_user_ids, escalation_level, html_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET 
			status = excluded.status,
			title = excluded.title,
			description = excluded.description,
			urgency = excluded.urgency,
			updated_at = excluded.updated_at,
			assigned_user_ids = excluded.assigned_user_ids,
			escalation_level = excluded.escalation_level
	`, incident.ID, incident.ServiceID, incident.Status, incident.Title,
		incident.Description, incident.Urgency, incident.IncidentNumber,
		incident.CreatedAt, incident.UpdatedAt, incident.Pinned,
		incident.AssignedUserIDs, incident.EscalationLevel, incident.HTMLURL)
	return err
}

func (r *Repository) PinIncident(id string, pinned bool) error {
	_, err := r.db.Exec("UPDATE incidents SET pinned = ? WHERE id = ?", pinned, id)
	return err
}

func (r *Repository) CleanOldIncidents(olderThan time.Time) error {
	_, err := r.db.Exec(`
		DELETE FROM incidents 
		WHERE status = 'resolved' AND updated_at < ?
	`, olderThan)
	return err
}

// Draft notes operations
func (r *Repository) GetDraftNote(incidentID string) (*DraftNote, error) {
	var note DraftNote
	err := r.db.Get(&note, "SELECT * FROM draft_notes WHERE incident_id = ?", incidentID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &note, err
}

func (r *Repository) SaveDraftNote(note DraftNote) error {
	note.LastUpdated = time.Now()
	_, err := r.db.Exec(`
		INSERT INTO draft_notes (
			incident_id, note_text, why_triggered, impact, actions, links, last_updated
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(incident_id) DO UPDATE SET 
			note_text = excluded.note_text,
			why_triggered = excluded.why_triggered,
			impact = excluded.impact,
			actions = excluded.actions,
			links = excluded.links,
			last_updated = excluded.last_updated
	`, note.IncidentID, note.NoteText, note.WhyTriggered, note.Impact, 
		note.Actions, note.Links, note.LastUpdated)
	return err
}

// Template operations
func (r *Repository) GetTemplates() ([]Template, error) {
	var templates []Template
	err := r.db.Select(&templates, "SELECT * FROM templates ORDER BY title")
	return templates, err
}

func (r *Repository) SaveTemplate(template Template) error {
	if template.ID == 0 {
		_, err := r.db.Exec(`
			INSERT INTO templates (title, body_text) VALUES (?, ?)
		`, template.Title, template.BodyText)
		return err
	}
	
	_, err := r.db.Exec(`
		UPDATE templates SET title = ?, body_text = ? WHERE id = ?
	`, template.Title, template.BodyText, template.ID)
	return err
}

func (r *Repository) DeleteTemplate(id int) error {
	_, err := r.db.Exec("DELETE FROM templates WHERE id = ?", id)
	return err
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}

// Helper functions for encryption
func encryptValue(value string) (string, error) {
	// For production, implement proper AES encryption
	// This is a simple example - in production use proper key management
	// For now, we'll just base64 encode with a prefix to identify encrypted values
	encoded := base64.StdEncoding.EncodeToString([]byte(value))
	return "enc:" + encoded, nil
}

func decryptValue(encrypted string) (string, error) {
	// Check if value is encrypted (has our prefix)
	if !strings.HasPrefix(encrypted, "enc:") {
		// Return as-is if not encrypted (for backward compatibility)
		return encrypted, nil
	}
	
	// Remove prefix and decode
	encoded := strings.TrimPrefix(encrypted, "enc:")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt value: %w", err)
	}
	
	return string(decoded), nil
}