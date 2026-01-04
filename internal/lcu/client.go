package lcu

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// Client LCU (League Client Update) API istemcisi
type Client struct {
	host      string
	port      string
	token     string
	client    *http.Client
	baseURL   string
	connected bool
}

// NewClient yeni bir LCU client oluşturur
func NewClient() (*Client, error) {
	lcu := &Client{
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	if err := lcu.connect(); err != nil {
		return nil, err
	}

	return lcu, nil
}

// connect League Client'a bağlanır
func (c *Client) connect() error {
	// League Client process'inden bilgileri al
	port, token, err := c.getClientInfo()
	if err != nil {
		return fmt.Errorf("League Client bulunamadı: %w", err)
	}

	c.port = port
	c.token = token
	c.host = "127.0.0.1"
	c.baseURL = fmt.Sprintf("https://%s:%s", c.host, c.port)
	c.connected = true

	return nil
}

// getClientInfo League Client'ın port ve token bilgilerini alır
func (c *Client) getClientInfo() (string, string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("wmic", "PROCESS", "WHERE", "name='LeagueClientUx.exe'", "GET", "commandline")
	case "darwin": // macOS
		cmd = exec.Command("bash", "-c", "ps aux | grep 'LeagueClientUx'")
	case "linux":
		cmd = exec.Command("bash", "-c", "ps aux | grep 'LeagueClientUx'")
	default:
		return "", "", fmt.Errorf("desteklenmeyen işletim sistemi: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", "", err
	}

	cmdLine := string(output)

	// Port'u bul
	portRegex := regexp.MustCompile(`--app-port=(\d+)`)
	portMatch := portRegex.FindStringSubmatch(cmdLine)
	if len(portMatch) < 2 {
		return "", "", fmt.Errorf("port bulunamadı")
	}

	// Token'ı bul
	tokenRegex := regexp.MustCompile(`--remoting-auth-token=([a-zA-Z0-9_-]+)`)
	tokenMatch := tokenRegex.FindStringSubmatch(cmdLine)
	if len(tokenMatch) < 2 {
		return "", "", fmt.Errorf("token bulunamadı")
	}

	return portMatch[1], tokenMatch[1], nil
}

// IsConnected client'ın bağlı olup olmadığını kontrol eder
func (c *Client) IsConnected() bool {
	return c.connected
}

// makeRequest LCU API'sine istek yapar
func (c *Client) makeRequest(method, endpoint string) ([]byte, error) {
	if !c.connected {
		return nil, fmt.Errorf("LCU'ya bağlı değil")
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Basic Auth ekle
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("riot:%s", c.token)))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, endpoint)
	}

	return io.ReadAll(resp.Body)
}

// GetCurrentSummoner aktif summoner bilgisini alır
func (c *Client) GetCurrentSummoner() (*Summoner, error) {
	data, err := c.makeRequest("GET", "/lol-summoner/v1/current-summoner")
	if err != nil {
		return nil, err
	}

	var summoner Summoner
	if err := json.Unmarshal(data, &summoner); err != nil {
		return nil, err
	}

	return &summoner, nil
}

// GetActiveGame aktif oyun bilgisini alır
func (c *Client) GetActiveGame() (*GameData, error) {
	data, err := c.makeRequest("GET", "/lol-gameflow/v1/session")
	if err != nil {
		return nil, err
	}

	var session GameFlowSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	// Oyun durumunu kontrol et
	if session.Phase != "InProgress" && session.Phase != "ChampSelect" {
		return nil, fmt.Errorf("aktif oyun yok, durum: %s", session.Phase)
	}

	gameData := &GameData{
		Phase:    session.Phase,
		GameTime: 0,
	}

	// Champion select'teyse
	if session.Phase == "ChampSelect" {
		champSelect, err := c.GetChampSelectSession()
		if err == nil {
			gameData.ChampSelect = champSelect
		}
	}

	// Oyun içindeyse
	if session.Phase == "InProgress" {
		gameInfo, err := c.GetInGameInfo()
		if err == nil {
			gameData.InGame = gameInfo
		}
	}

	return gameData, nil
}

// GetChampSelectSession champion select bilgisini alır
func (c *Client) GetChampSelectSession() (*ChampSelectSession, error) {
	data, err := c.makeRequest("GET", "/lol-champ-select/v1/session")
	if err != nil {
		return nil, err
	}

	var session ChampSelectSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// GetInGameInfo oyun içi bilgileri alır
func (c *Client) GetInGameInfo() (*InGameInfo, error) {
	// Aktif oyuncu bilgisi
	playerData, err := c.makeRequest("GET", "/lol-gameflow/v1/session")
	if err != nil {
		return nil, err
	}

	var session GameFlowSession
	if err := json.Unmarshal(playerData, &session); err != nil {
		return nil, err
	}

	info := &InGameInfo{
		GameTime: 0, // Gerçek game time için ayrı endpoint gerekli
		Players:  []Player{},
	}

	return info, nil
}

// TryConnect bağlantı denemesi yapar (hata döndürmez)
func (c *Client) TryConnect() bool {
	if err := c.connect(); err != nil {
		c.connected = false
		return false
	}
	return true
}

// Reconnect yeniden bağlanmayı dener
func (c *Client) Reconnect() error {
	c.connected = false
	return c.connect()
}

// GetLockfile lockfile'dan bilgileri okur (alternatif yöntem)
func (c *Client) GetLockfile() (string, string, error) {
	var lockfilePath string

	switch runtime.GOOS {
	case "windows":
		lockfilePath = os.Getenv("LOCALAPPDATA") + "\\Riot Games\\League of Legends\\lockfile"
	case "darwin":
		lockfilePath = os.Getenv("HOME") + "/Library/Application Support/Riot Games/League of Legends/lockfile"
	default:
		return "", "", fmt.Errorf("desteklenmeyen platform")
	}

	data, err := os.ReadFile(lockfilePath)
	if err != nil {
		return "", "", err
	}

	// Format: LeagueClient:PID:PORT:TOKEN:PROTOCOL
	parts := strings.Split(string(data), ":")
	if len(parts) < 5 {
		return "", "", fmt.Errorf("geçersiz lockfile formatı")
	}

	return parts[2], parts[3], nil
}
