package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Service AI servisi
type Service struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewService yeni bir AI servisi oluşturur
func NewService() (*Service, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	// Model ayarları
	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.ResponseMIMEType = "application/json"

	return &Service{
		client: client,
		model:  model,
	}, nil
}

// AnalyzeGame oyun durumunu analiz eder ve önerilerde bulunur
func (s *Service) AnalyzeGame(req AnalysisRequest) (*AnalysisResponse, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf(`
		You are a League of Legends expert coach. Analyze the current game state and provide advice.
		
		Current State:
		- Phase: %s
		- My Champion: %s
		- Current Items: %s
		- Gold Available: %d
		- Enemy Champions: %s
		- Game Time: %d seconds

		Provide a JSON response with the following structure:
		{
			"suggestion": "General gameplay advice based on the current situation",
			"next_items": ["Item 1", "Item 2"],
			"strategy": "Specific strategic advice (e.g., play safe, roam, freeze lane)"
		}
		
		Focus on the next best item to buy with the available gold and the best strategy against the enemy team composition.
	`, req.GamePhase, req.Champion, strings.Join(req.Items, ", "), req.Gold, strings.Join(req.EnemyChamps, ", "), req.GameTime)

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// JSON yanıtını parse et
	var analysisResp AnalysisResponse

	// Part'ı string'e çevirip JSON unmarshal yap
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), &analysisResp); err != nil {
				// JSON bloğunu temizlemeyi dene (markdown ```json ... ``` varsa)
				cleanJSON := strings.TrimPrefix(string(txt), "```json")
				cleanJSON = strings.TrimPrefix(cleanJSON, "```")
				cleanJSON = strings.TrimSuffix(cleanJSON, "```")

				if err2 := json.Unmarshal([]byte(cleanJSON), &analysisResp); err2 != nil {
					return nil, fmt.Errorf("failed to parse JSON response: %w", err)
				}
			}
			break
		}
	}

	return &analysisResp, nil
}

// Close servisi kapatır
func (s *Service) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
