package ai

// AnalysisRequest AI analiz isteği
type AnalysisRequest struct {
	GamePhase   string   `json:"game_phase"`
	Champion    string   `json:"champion"`
	Items       []string `json:"items"`
	Gold        int      `json:"gold"`
	EnemyChamps []string `json:"enemy_champions"`
	GameTime    int      `json:"game_time"` // Saniye cinsinden
}

// AnalysisResponse AI analiz cevabı
type AnalysisResponse struct {
	Suggestion string   `json:"suggestion"`
	NextItems  []string `json:"next_items"`
	Strategy   string   `json:"strategy"`
}
