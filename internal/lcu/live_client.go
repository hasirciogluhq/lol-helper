package lcu

import (
"crypto/tls"
"encoding/json"
"fmt"
"io"
"net/http"
"time"
)

// LiveClient Live Client Data API istemcisi
type LiveClient struct {
	client  *http.Client
	baseURL string
}

// NewLiveClient yeni bir LiveClient oluşturur
func NewLiveClient() *LiveClient {
	return &LiveClient{
		client: &http.Client{
			Timeout: 2 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Self-signed sertifika için gerekli
				},
			},
		},
		baseURL: "https://127.0.0.1:2999/liveclientdata",
	}
}

// GetAllGameData tüm oyun verisini çeker
func (c *LiveClient) GetAllGameData() (*LiveGameData, error) {
	resp, err := c.client.Get(c.baseURL + "/allgamedata")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Live Client API hatası: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data LiveGameData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
