package catapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Breed represents a cat breed from TheCatAPI.
type Breed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client is a client for TheCatAPI.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	breeds     []Breed
	mu         sync.RWMutex
	lastFetch  time.Time
}

// NewClient creates a new CatAPI client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBreeds fetches all cat breeds from TheCatAPI, with caching.
func (c *Client) GetBreeds() ([]Breed, error) {
	c.mu.RLock()
	// Cache for 1 hour
	if time.Since(c.lastFetch) < time.Hour {
		defer c.mu.RUnlock()
		return c.breeds, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double check in case another goroutine just fetched the data
	if time.Since(c.lastFetch) < time.Hour {
		return c.breeds, nil
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/breeds", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch breeds: received status code %d", resp.StatusCode)
	}

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, err
	}

	c.breeds = breeds
	c.lastFetch = time.Now()

	return c.breeds, nil
}

// IsValidBreed checks if a given breed name is valid.
func (c *Client) IsValidBreed(breedName string) (bool, error) {
	breeds, err := c.GetBreeds()
	if err != nil {
		return false, err
	}

	for _, b := range breeds {
		if b.Name == breedName {
			return true, nil
		}
	}

	return false, nil
}
