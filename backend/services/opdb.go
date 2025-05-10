package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"backend/models"

	"gorm.io/gorm"
)

type OPDBService struct {
	db *gorm.DB
}

type OPDBMachineResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Manufacturer string `json:"manufacturer"`
	Year      int    `json:"year"`
	IPDBID    int    `json:"ipdb_id"`
	Type      string `json:"type"`
	OPDBID    string `json:"opdb_id"`
	IsPinball bool   `json:"is_pinball"`
	IsGroup   bool   `json:"is_group"`
	IsAlias   bool   `json:"is_alias"`
}

type OPDBRequest struct {
	APIToken string `json:"api_token"`
}

func NewOPDBService(db *gorm.DB) *OPDBService {
	return &OPDBService{db: db}
}

func (s *OPDBService) GetMachine(opdbID string) (*models.Machine, error) {
	// First, try to get from database
	var machine models.Machine
	result := s.db.Where("opdb_id = ?", opdbID).First(&machine)

	// If found and recently updated (within last 24 hours), return cached data
	if result.Error == nil && time.Since(machine.UpdatedAt) < 24*time.Hour {
		return &machine, nil
	}

	// If not found or data is stale, fetch from OPDB API
	apiToken := os.Getenv("OPDB_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("OPDB_API_TOKEN not set")
	}

	// Make POST request with URL parameters
	fmt.Println("Fetching from OPDB API:", opdbID)
	params := url.Values{}
	params.Add("api_token", apiToken)
	url := fmt.Sprintf("https://opdb.org/api/machines/%s?%s", opdbID, params.Encode())
	fmt.Println("URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from OPDB API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OPDB API returned status code: %d", resp.StatusCode)
	}

	var opdbMachine OPDBMachineResponse
	if err := json.NewDecoder(resp.Body).Decode(&opdbMachine); err != nil {
		fmt.Println("OPDB Machine:", opdbMachine)
		fmt.Println("Error decoding OPDB response:", err)
		fmt.Println("Response body:", resp.Body)
		return nil, fmt.Errorf("failed to decode OPDB response: %v", err)
	}
	// Convert OPDB response to our model
	machine = models.Machine{
		OPDBID: opdbMachine.OPDBID,
		Name:   opdbMachine.Name,
		// Manufacturer:  opdbMachine.Manufacturer,
		Year:      opdbMachine.Year,
		IPDBID:    opdbMachine.IPDBID,
		Type:      opdbMachine.Type,
		IsPinball: opdbMachine.IsPinball,
		IsGroup:   opdbMachine.IsGroup,
		IsAlias:   opdbMachine.IsAlias,
		UpdatedAt: time.Now(),
	}

	// Save to database (create or update)
	if result.Error == nil {
		// Update existing record
		if err := s.db.Model(&machine).Updates(machine).Error; err != nil {
			return nil, fmt.Errorf("failed to update machine in database: %v", err)
		}
	} else {
		// Create new record
		if err := s.db.Create(&machine).Error; err != nil {
			return nil, fmt.Errorf("failed to create machine in database: %v", err)
		}
	}

	return &machine, nil
}
