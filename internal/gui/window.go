package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"lol-helper/internal/lol"
)

// MainWindow ana pencere yapısı
type MainWindow struct {
	app     fyne.App
	window  fyne.Window
	service *lol.Service

	// UI Components
	statusLabel     *widget.Label
	phaseLabel      *widget.Label
	championLabel   *widget.Label
	goldLabel       *widget.Label
	suggestionLabel *widget.Label
	itemsList       *widget.List
	nextItemsLabel  *widget.Label
	strategyLabel   *widget.Label
}

// NewMainWindow yeni bir ana pencere oluşturur
func NewMainWindow() *MainWindow {
	a := app.New()
	a.Settings().SetTheme(&LoLTheme{})

	w := a.NewWindow("LoL Helper AI")
	w.Resize(fyne.NewSize(400, 600))

	mw := &MainWindow{
		app:    a,
		window: w,
	}

	mw.setupUI()
	return mw
}

// setupUI arayüz bileşenlerini oluşturur
func (mw *MainWindow) setupUI() {
	// Status Section
	mw.statusLabel = widget.NewLabel("Durum: Başlatılıyor...")
	mw.phaseLabel = widget.NewLabel("Oyun Fazı: -")

	// Game Info Section
	mw.championLabel = widget.NewLabel("Şampiyon: -")
	mw.goldLabel = widget.NewLabel("Altın: 0")

	// AI Suggestion Section
	mw.suggestionLabel = widget.NewLabel("Öneri: Bekleniyor...")
	mw.suggestionLabel.Wrapping = fyne.TextWrapWord

	mw.nextItemsLabel = widget.NewLabel("Sonraki İtemler: -")
	mw.strategyLabel = widget.NewLabel("Strateji: -")
	mw.strategyLabel.Wrapping = fyne.TextWrapWord

	// Layout
	content := container.NewVBox(
		widget.NewCard("Bağlantı Durumu", "", container.NewVBox(
			mw.statusLabel,
			mw.phaseLabel,
		)),
		widget.NewCard("Oyun Bilgisi", "", container.NewVBox(
			mw.championLabel,
			mw.goldLabel,
		)),
		widget.NewCard("AI Koç", "", container.NewVBox(
			mw.suggestionLabel,
			mw.nextItemsLabel,
			mw.strategyLabel,
		)),
	)

	mw.window.SetContent(content)
}

// Start uygulamayı başlatır
func (mw *MainWindow) Start() {
	// Servisi başlat
	service, err := lol.NewService(mw.UpdateUI)
	if err != nil {
		mw.statusLabel.SetText(fmt.Sprintf("Hata: %v", err))
	} else {
		mw.service = service
		mw.service.Start()
	}

	mw.window.ShowAndRun()

	// Kapanırken servisi durdur
	if mw.service != nil {
		mw.service.Stop()
	}
}

// UpdateUI arayüzü günceller (thread-safe)
func (mw *MainWindow) UpdateUI(state *lol.HelperState) {
	// Fyne UI güncellemeleri main thread'de yapılmalı
	// Ancak widget.Label.SetText thread-safe'dir

	if state.Error != nil {
		mw.statusLabel.SetText(fmt.Sprintf("Hata: %v", state.Error))
	} else if state.Game.IsConnected {
		mw.statusLabel.SetText("Durum: Bağlı")
	} else {
		mw.statusLabel.SetText("Durum: Bağlantı Bekleniyor...")
	}

	mw.phaseLabel.SetText(fmt.Sprintf("Oyun Fazı: %s", state.Game.Phase))
	mw.championLabel.SetText(fmt.Sprintf("Şampiyon: %s", state.Game.Champion))
	mw.goldLabel.SetText(fmt.Sprintf("Altın: %d", state.Game.Gold))

	if state.Recommendation != nil {
		mw.suggestionLabel.SetText(fmt.Sprintf("Öneri: %s", state.Recommendation.Suggestion))
		mw.nextItemsLabel.SetText(fmt.Sprintf("Sonraki İtemler: %s", strings.Join(state.Recommendation.NextItems, ", ")))
		mw.strategyLabel.SetText(fmt.Sprintf("Strateji: %s", state.Recommendation.Strategy))
	}
}
