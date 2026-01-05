package lol

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"lol-helper/internal/ai"
	"lol-helper/internal/lcu"
)

// Service LoL Helper ana servisi
type Service struct {
	lcuClient     *lcu.Client
	liveClient    *lcu.LiveClient
	aiService     *ai.Service
	state         *HelperState
	stopChan      chan struct{}
	onUpdate      func(*HelperState)
	lastStateHash string // State değişiklik kontrolü için
}

// NewService yeni bir servis oluşturur
func NewService(onUpdate func(*HelperState)) (*Service, error) {
	// LCU Client başlat (bağlantı hatası olsa bile devam et, polling ile deneyecek)
	lcuClient, _ := lcu.NewClient()
	liveClient := lcu.NewLiveClient()

	// AI Service başlat
	aiService, err := ai.NewService()
	if err != nil {
		return nil, fmt.Errorf("AI servisi başlatılamadı: %w", err)
	}

	return &Service{
		lcuClient:  lcuClient,
		liveClient: liveClient,
		aiService:  aiService,
		state:      NewHelperState(),
		stopChan:   make(chan struct{}),
		onUpdate:   onUpdate,
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
	ticker := time.NewTicker(3 * time.Second) // Her 3 saniyede bir güncelle (blinking önlemek için)
	defer ticker.Stop()

	aiTicker := time.NewTicker(20 * time.Second) // AI analizi her 20 saniyede bir
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
	// 1. Önce Live Client (Oyun İçi API) kontrol et
	// Bu API sadece oyun içindeyken çalışır ve en doğru veriyi verir.
	liveData, liveErr := s.liveClient.GetAllGameData()

	if liveErr == nil && liveData != nil {
		// Oyun içindeyiz ve veri alabiliyoruz
		s.state.Game.IsConnected = true
		s.state.Game.Phase = "InProgress"
		s.state.Game.AllPlayers = liveData.AllPlayers
		s.state.Error = nil

		// Aktif oyuncu verilerini güncelle
		for _, p := range liveData.AllPlayers {
			if p.SummonerName == liveData.ActivePlayer.SummonerName {
				s.state.Game.Gold = int(liveData.ActivePlayer.CurrentGold)
				s.state.Game.Champion = p.ChampionName // Şampiyon ismini buradan al

				// İtemleri güncelle
				var items []string
				for _, item := range p.Items {
					items = append(items, item.DisplayName)
				}
				s.state.Game.Items = items
				break
			}
		}

		// LCU bağlantısını arka planda dene ama başarısız olsa bile akışı bozma
		if s.lcuClient == nil || !s.lcuClient.IsConnected() {
			if s.lcuClient == nil {
				client, err := lcu.NewClient()
				if err == nil {
					s.lcuClient = client
				}
			} else {
				s.lcuClient.TryConnect()
			}
		}

		s.notifyUpdate()
		return
	}

	// 2. Eğer Live Client yanıt vermiyorsa, LCU (Client API) kontrol et
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
			// İkisi de yoksa bağlantı yok demektir
			s.state.Game.Phase = "Disconnected"
			s.state.Game.IsConnected = false
			// Disconnected durumunda player listesini temizlemiyoruz
			// Böylece anlık kopmalarda liste kaybolmaz
			s.notifyUpdate()
			return
		}
	}

	// LCU Bağlı, verileri çek
	gameData, err := s.lcuClient.GetActiveGame()
	if err != nil {
		// Oyun yok (Lobby veya başka bir durum)
		s.state.Game.Phase = "Lobby/None"
		s.state.Game.IsConnected = true // Client bağlı ama oyun yok
		s.state.Error = nil             // Bu bir hata değil, durum

		// Lobby'deysek player listesini temizle
		s.state.Game.AllPlayers = nil

		s.notifyUpdate()
		return
	}

	// Oyun var (ChampSelect veya InProgress ama LiveClient henüz hazır değil)
	summoner, err := s.lcuClient.GetCurrentSummoner()
	if err != nil {
		log.Printf("Summoner bilgisi alınamadı: %v", err)
	}

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

// notifyUpdate UI'ı günceller - sadece state değiştiyse
func (s *Service) notifyUpdate() {
	// State'in hash'ini hesapla
	currentHash := s.calculateStateHash()

	// Eğer hash aynıysa, hiçbir şey değişmemiş demektir, UI'ı güncelleme
	if currentHash == s.lastStateHash {
		return
	}

	s.lastStateHash = currentHash
	s.state.LastUpdate = time.Now().UnixMilli()
	if s.onUpdate != nil {
		s.onUpdate(s.state)
	}
}

// calculateStateHash state'in hash'ini hesaplar
func (s *Service) calculateStateHash() string {
	// Player isimlerini topla
	var playerNames string
	for _, p := range s.state.Game.AllPlayers {
		playerNames += p.SummonerName + ","
	}

	// Sadece UI'ı etkileyen alanları hash'le
	data := struct {
		Phase       string
		IsConnected bool
		PlayerNames string
		Gold        int
		Champion    string
		ItemCount   int
	}{
		Phase:       s.state.Game.Phase,
		IsConnected: s.state.Game.IsConnected,
		PlayerNames: playerNames,
		Gold:        s.state.Game.Gold,
		Champion:    s.state.Game.Champion,
		ItemCount:   len(s.state.Game.Items),
	}

	jsonData, _ := json.Marshal(data)
	hash := md5.Sum(jsonData)
	return fmt.Sprintf("%x", hash)
}
