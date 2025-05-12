package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type IFPAService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type IFPAPlayer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Country   string `json:"country"`
	City      string `json:"city"`
	State     string `json:"state"`
}

type IFPAPlayerResponse struct {
	Player IFPAPlayer `json:"player"`
}

func NewIFPAService() *IFPAService {
	apiKey := os.Getenv("IFPA_API_KEY")
	if apiKey == "" {
		panic("IFPA_API_KEY environment variable is required")
	}

	return &IFPAService{
		apiKey:  apiKey,
		baseURL: "https://api.ifpapinball.com/v1",
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *IFPAService) GetPlayerByIFPANumber(ifpaNumber int) (*IFPAPlayer, error) {
	url := fmt.Sprintf("%s/player/%d", s.baseURL, ifpaNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IFPA API error: %s - %s", resp.Status, string(body))
	}

	var response IFPAPlayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Player, nil
}
