package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	tokenURL = "https://id.twitch.tv/oauth2/token"
	apiURL   = "https://api.igdb.com/v4/games"
)

// TokenResponse represents the OAuth2 token response from Twitch
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GameSearchResult represents a game search result from IGDB
type GameSearchResult struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Client handles IGDB API requests with automatic token refresh
type Client struct {
	clientID     string
	clientSecret string
	accessToken  string
	expiresAt    time.Time
	mu           sync.RWMutex
	httpClient   *http.Client
}

// NewClient creates a new IGDB client
func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// getAccessToken requests a new access token from Twitch OAuth
func (c *Client) getAccessToken() error {
	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&grant_type=client_credentials",
		tokenURL, c.clientID, c.clientSecret)

	resp, err := c.httpClient.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.mu.Lock()
	c.accessToken = tokenResp.AccessToken
	// Set expiration with 5 minute buffer to refresh before actual expiration
	c.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn)*time.Second - 5*time.Minute)
	c.mu.Unlock()

	return nil
}

// ensureValidToken checks if the token is valid and refreshes if needed
func (c *Client) ensureValidToken() error {
	c.mu.RLock()
	needsRefresh := c.accessToken == "" || time.Now().After(c.expiresAt)
	c.mu.RUnlock()

	if needsRefresh {
		return c.getAccessToken()
	}

	return nil
}

// SearchGames searches for games by name using the IGDB API
func (c *Client) SearchGames(query string) ([]GameSearchResult, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, fmt.Errorf("failed to ensure valid token: %w", err)
	}

	// Build IGDB query - search by name, return only main games (game_type = 0)
	body := fmt.Sprintf(`search "%s"; fields name,url; where game_type = 0;`, query)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.mu.RLock()
	accessToken := c.accessToken
	c.mu.RUnlock()

	req.Header.Set("Client-ID", c.clientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search games: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var results []GameSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return results, nil
}
