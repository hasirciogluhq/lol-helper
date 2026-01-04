package lol

import (
	"fmt"
	"log"
	"time"

	"lol-helper/internal/ai"
	"lol-helper/internal/lcu"
)

// Service LoL Helper ana servisi
type Service struct {
	lcuClient *lcu.Client
	aiService *ai.Service
	state     *HelperState
	stopChan  chan struct{}
	onUpdate  func(*HelperState)
}

// NewService yeni bir servis oluşturur
func NewService(onUpdate func(*HelperState)) (*Service, error) {
	// LCU Client başlat (bağlantı hatası olsa bile devam et, polling ile deneyecek)
	lcuClient, _ := lcu.NewClient()

	// AI Service başlat
	aiService, err := ai.NewService()
	if err != nil {
		return nil, fmt.Errorf("AI servisi başlatılamadı: %w", err)
	}

	return &Service{
		lcuClient: lcuClient,
		aiService: aiService,
		state:     NewHelperState(),
		stopChan:  make(chan struct{}),
		onUpdate:  onUpdate,
	}, nil
}

// Start polling işlemini başlatır
func (s *Service) Start() {
	go s.pollLoop()
}

// Stop polling işlemini durdurur
func (s *Service) Stop() {
	close(s.stopChan)
	s.aiService.Close()
}

// pollLoop sürekli olarak LCU'dan veri çeker
func (s *Service) pollLoop() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	aiTicker := time.NewTicker(10 * time.Second) // AI analizi her 10 saniyede bir
	defer aiTicker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.updateGameState()
		case <-aiTicker.C:
			s.runAIAnalysis()
		}
	}
}

// updateGameState oyun durumunu günceller
func (s *Service) updateGameState() {
	// Bağlantı kontrolü
	if s.lcuClient == nil || !s.lcuClient.IsConnected() {
		if s.lcuClient == nil {
			client, err := lcu.NewClient()
			if err == nil {
				s.lcuClient = client
			}
		} else {
			s.lcuClient.TryConnect()
		}

		if s.lcuClient == nil || !s.lcuClient.IsConnected() {
			s.state.Game.Phase = "Disconnected"
			s.state.Game.IsConnected = false
			s.notifyUpdate()
			return
		}
	}

	// Verileri çek
	gameData, err := s.lcuClient.GetActiveGame()
	if err != nil {
		// Oyun yoksa veya hata varsa
		s.state.Game.Phase = "Lobby/None"
		s.state.Error = err
		s.notifyUpdate()
		return
	}

	summoner, err := s.lcuClient.GetCurrentSummoner()
	if err != nil {
		log.Printf("Summoner bilgisi alınamadı: %v", err)
	}

	// State güncelle
	s.state.UpdateFromLCU(gameData, summoner)
	s.state.Error = nil
	s.notifyUpdate()
}

// runAIAnalysis AI analizi yapar
func (s *Service) runAIAnalysis() {
	// Sadece oyun içindeyse veya şampiyon seçimindeyse analiz yap
	if s.state.Game.Phase != "InProgress" && s.state.Game.Phase != "ChampSelect" {
		return
	}

	req := ai.AnalysisRequest{
		GamePhase:   s.state.Game.Phase,
		Champion:    s.state.Game.Champion,
		Items:       s.state.Game.Items,
		Gold:        s.state.Game.Gold,
		EnemyChamps: s.state.Game.EnemyChamps,
		GameTime:    s.state.Game.GameTime,
	}

	resp, err := s.aiService.AnalyzeGame(req)
	if err != nil {
		log.Printf("AI analizi hatası: %v", err)
		s.state.Error = err
		s.notifyUpdate()
		return
	}

	s.state.Recommendation = &Recommendation{
		Suggestion: resp.Suggestion,
		NextItems:  resp.NextItems,
		Strategy:   resp.Strategy,
	}
	s.notifyUpdate()
}

// notifyUpdate UI'ı günceller
func (s *Service) notifyUpdate() {
	s.state.LastUpdate = time.Now().UnixMilli()
	if s.onUpdate != nil {
		s.onUpdate(s.state)
	}
}
