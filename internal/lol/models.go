package lol

import "lol-helper/internal/lcu"

// GameState oyun durumu
type GameState struct {
	Phase       string
	Champion    string
	Items       []string
	Gold        int
	EnemyChamps []string
	GameTime    int
	IsConnected bool
	AllPlayers  []lcu.LivePlayer
}

// Recommendation AI önerisi
type Recommendation struct {
	Suggestion string
	NextItems  []string
	Strategy   string
}

// HelperState uygulamanın genel durumu
type HelperState struct {
	Game           *GameState
	Recommendation *Recommendation
	LastUpdate     int64
	Error          error
}

// NewHelperState yeni bir durum oluşturur
func NewHelperState() *HelperState {
	return &HelperState{
		Game: &GameState{
			Phase:       "Disconnected",
			Items:       []string{},
			EnemyChamps: []string{},
		},
		Recommendation: &Recommendation{
			NextItems: []string{},
		},
	}
}

// UpdateFromLCU LCU verisiyle durumu günceller
func (s *HelperState) UpdateFromLCU(gameData *lcu.GameData, summoner *lcu.Summoner) {
	s.Game.Phase = gameData.Phase
	s.Game.IsConnected = true

	if gameData.Phase == "ChampSelect" && gameData.ChampSelect != nil {
		// Champ select mantığı buraya eklenebilir
		// Şu anlık basit tutuyoruz
	} else if gameData.Phase == "InProgress" && gameData.InGame != nil {
		s.Game.GameTime = int(gameData.InGame.GameTime)
		// Oyuncu ve item bilgileri buraya eklenecek
		// LCU API'den detaylı oyuncu verisi çekilmesi gerekebilir
	}
}
