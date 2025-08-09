package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"pager-ops/database"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type IncidentPoller struct {
	ctx           context.Context
	pd            *PagerDutyService
	repository    *database.Repository
	interval      time.Duration
	stopCh        chan struct{}
	lastIncidents map[string]Incident
	mu            sync.RWMutex
	soundPlayer   *SoundPlayer
}

func NewIncidentPoller(ctx context.Context, pd *PagerDutyService, repo *database.Repository) *IncidentPoller {
	return &IncidentPoller{
		ctx:           ctx,
		pd:            pd,
		repository:    repo,
		interval:      3 * time.Second,
		stopCh:        make(chan struct{}),
		lastIncidents: make(map[string]Incident),
		soundPlayer:   NewSoundPlayer(repo),
	}
}

func (p *IncidentPoller) Start() {
	go p.pollLoop()
}

func (p *IncidentPoller) Stop() {
	close(p.stopCh)
}

func (p *IncidentPoller) pollLoop() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// Initial poll
	p.poll()

	for {
		select {
		case <-ticker.C:
			p.poll()
		case <-p.stopCh:
			return
		}
	}
}

func (p *IncidentPoller) poll() {
	// Get settings
	assignedOnly, _ := p.repository.GetSetting(database.SettingAssignedOnly)
	redirectEnabled, _ := p.repository.GetSetting(database.SettingRedirectEnabled)

	// Get enabled services
	services, err := p.repository.GetServices()
	if err != nil {
		wailsruntime.LogError(p.ctx, fmt.Sprintf("Failed to get services: %v", err))
		return
	}

	enabledServiceIDs := []string{}
	for _, svc := range services {
		if svc.Enabled {
			enabledServiceIDs = append(enabledServiceIDs, svc.ID)
		}
	}

	if len(enabledServiceIDs) == 0 {
		return
	}

	// Prepare filters
	filters := IncidentFilters{
		ServiceIDs:   enabledServiceIDs,
		AssignedToMe: assignedOnly == "true",
	}

	// Fetch incidents in parallel
	var wg sync.WaitGroup
	wg.Add(3)

	var triggeredIncidents, acknowledgedIncidents, resolvedIncidents []Incident
	var errTriggered, errAcknowledged, errResolved error

	// Fetch triggered incidents
	go func() {
		defer wg.Done()
		f := filters
		f.Statuses = []string{"triggered"}
		triggeredIncidents, errTriggered = p.pd.FetchIncidents(p.ctx, f)
	}()

	// Fetch acknowledged incidents
	go func() {
		defer wg.Done()
		f := filters
		f.Statuses = []string{"acknowledged"}
		acknowledgedIncidents, errAcknowledged = p.pd.FetchIncidents(p.ctx, f)
	}()

	// Fetch resolved incidents (last 24 hours)
	go func() {
		defer wg.Done()
		f := filters
		f.Statuses = []string{"resolved"}
		since := time.Now().Add(-24 * time.Hour)
		f.Since = &since
		resolvedIncidents, errResolved = p.pd.FetchIncidents(p.ctx, f)
	}()

	wg.Wait()

	// Handle errors
	if errTriggered != nil {
		wailsruntime.LogError(p.ctx, fmt.Sprintf("Failed to fetch triggered incidents: %v", errTriggered))
	}
	if errAcknowledged != nil {
		wailsruntime.LogError(p.ctx, fmt.Sprintf("Failed to fetch acknowledged incidents: %v", errAcknowledged))
	}
	if errResolved != nil {
		wailsruntime.LogError(p.ctx, fmt.Sprintf("Failed to fetch resolved incidents: %v", errResolved))
	}

	// Combine all incidents
	allIncidents := append(triggeredIncidents, acknowledgedIncidents...)
	allIncidents = append(allIncidents, resolvedIncidents...)

	// Detect new incidents
	newIncidents := p.detectNewIncidents(allIncidents)

	// Handle new incidents
	if len(newIncidents) > 0 {
		// Play sound for new triggered incidents
		for _, inc := range newIncidents {
			if inc.Status == "triggered" {
				p.soundPlayer.PlayAlert()

				// Open in browser if enabled
				if redirectEnabled == "true" && inc.HTMLURL != "" {
					wailsruntime.BrowserOpenURL(p.ctx, inc.HTMLURL)
				}
			}
		}
	}

	// Cache incidents in database
	for _, inc := range allIncidents {
		cached := database.CachedIncident{
			ID:              inc.ID,
			ServiceID:       inc.ServiceID,
			Status:          inc.Status,
			Title:           inc.Title,
			Description:     inc.Description,
			Urgency:         inc.Urgency,
			IncidentNumber:  inc.IncidentNumber,
			CreatedAt:       inc.CreatedAt,
			UpdatedAt:       inc.UpdatedAt,
			AssignedUserIDs: database.JSONArray(inc.AssignedUserIDs),
			EscalationLevel: inc.EscalationLevel,
			HTMLURL:         inc.HTMLURL,
		}
		p.repository.UpsertIncident(cached)
	}

	// Clean old resolved incidents
	p.repository.CleanOldIncidents(time.Now().Add(-48 * time.Hour))

	// Emit update event to frontend
	wailsruntime.EventsEmit(p.ctx, "incidents:updated", map[string]interface{}{
		"incidents":         allIncidents,
		"newIncidents":      newIncidents,
		"totalCount":        len(allIncidents),
		"triggeredCount":    len(triggeredIncidents),
		"acknowledgedCount": len(acknowledgedIncidents),
		"resolvedCount":     len(resolvedIncidents),
	})
}

func (p *IncidentPoller) detectNewIncidents(incidents []Incident) []Incident {
	p.mu.Lock()
	defer p.mu.Unlock()

	newIncidents := []Incident{}
	currentMap := make(map[string]Incident)

	for _, inc := range incidents {
		currentMap[inc.ID] = inc

		// Check if this is a new incident
		if lastInc, exists := p.lastIncidents[inc.ID]; !exists {
			newIncidents = append(newIncidents, inc)
		} else if lastInc.Status != inc.Status && inc.Status == "triggered" {
			// Status changed to triggered
			newIncidents = append(newIncidents, inc)
		}
	}

	p.lastIncidents = currentMap
	return newIncidents
}

// SoundPlayer handles alert sound playback
type SoundPlayer struct {
	repository *database.Repository
}

func NewSoundPlayer(repo *database.Repository) *SoundPlayer {
	return &SoundPlayer{repository: repo}
}

func (s *SoundPlayer) PlayAlert() {
	soundEnabled, _ := s.repository.GetSetting(database.SettingSoundEnabled)
	if soundEnabled != "true" {
		return
	}

	soundPath, _ := s.repository.GetSetting(database.SettingSoundPath)
	if soundPath == "" {
		return
	}

	// Find first mp3 or wav file in the directory
	files, err := ioutil.ReadDir(soundPath)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".mp3" || ext == ".wav" {
			soundFile := filepath.Join(soundPath, file.Name())
			// Play sound (platform specific implementation needed)
			s.playSound(soundFile)
			break
		}
	}
}

func (s *SoundPlayer) playSound(file string) {
	// Platform-specific sound playback
	// This is a simplified version - you'll need to implement actual audio playback

	switch runtime.GOOS {
	case "darwin": // macOS
		if _, err := os.Stat("/usr/bin/afplay"); err == nil {
			go func() {
				cmd := exec.Command("/usr/bin/afplay", file)
				cmd.Run()
			}()
		}
	case "windows":
		// Windows: Use PowerShell or Windows Media Player
		go func() {
			cmd := exec.Command("powershell", "-c", fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", file))
			cmd.Run()
		}()
	case "linux":
		// Linux: Try multiple audio players
		players := []string{"aplay", "paplay", "ffplay"}
		for _, player := range players {
			if path, err := exec.LookPath(player); err == nil {
				go func() {
					cmd := exec.Command(path, file)
					cmd.Run()
				}()
				break
			}
		}
	}
}
