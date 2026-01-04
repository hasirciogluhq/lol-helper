package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"lol-helper/internal/lcu"
	"lol-helper/internal/lol"
)

// MainWindow ana pencere yapısı
type MainWindow struct {
	app     fyne.App
	window  fyne.Window
	service *lol.Service

	// UI Components
	statusLabel *widget.Label
	phaseLabel  *widget.Label

	// Team Containers
	teamOrderContainer *fyne.Container
	teamChaosContainer *fyne.Container

	// AI Suggestion Section
	suggestionLabel *widget.Label
	nextItemsLabel  *widget.Label
	strategyLabel   *widget.Label
}

// NewMainWindow yeni bir ana pencere oluşturur
func NewMainWindow() *MainWindow {
	a := app.New()
	a.Settings().SetTheme(&LoLTheme{})

	w := a.NewWindow("LoL Helper AI")
	w.Resize(fyne.NewSize(1000, 800))

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

	// AI Suggestion Section
	mw.suggestionLabel = widget.NewLabel("Öneri: Bekleniyor...")
	mw.suggestionLabel.Wrapping = fyne.TextWrapWord

	mw.nextItemsLabel = widget.NewLabel("Sonraki İtemler: -")
	mw.strategyLabel = widget.NewLabel("Strateji: -")
	mw.strategyLabel.Wrapping = fyne.TextWrapWord

	// Team Containers
	mw.teamOrderContainer = container.NewVBox()
	mw.teamChaosContainer = container.NewVBox()

	// Teams Split View
	teamsSplit := container.NewHSplit(
		container.NewVScroll(container.NewPadded(mw.teamOrderContainer)),
		container.NewVScroll(container.NewPadded(mw.teamChaosContainer)),
	)
	teamsSplit.SetOffset(0.5)

	// Top Info
	topInfo := container.NewVBox(
		mw.statusLabel,
		mw.phaseLabel,
	)

	// Bottom AI
	bottomAI := widget.NewCard("AI Koç", "", container.NewVBox(
		mw.suggestionLabel,
		mw.nextItemsLabel,
		mw.strategyLabel,
	))

	// Main Layout
	content := container.NewBorder(
		topInfo,
		bottomAI,
		nil, nil,
		teamsSplit,
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
	if state.Error != nil {
		mw.statusLabel.SetText(fmt.Sprintf("Hata: %v", state.Error))
	} else if state.Game.IsConnected {
		mw.statusLabel.SetText("Durum: Bağlı")
	} else {
		mw.statusLabel.SetText("Durum: Bağlantı Bekleniyor...")
	}

	mw.phaseLabel.SetText(fmt.Sprintf("Oyun Fazı: %s", state.Game.Phase))

	if state.Recommendation != nil {
		mw.suggestionLabel.SetText(fmt.Sprintf("Öneri: %s", state.Recommendation.Suggestion))
		mw.nextItemsLabel.SetText(fmt.Sprintf("Sonraki İtemler: %s", strings.Join(state.Recommendation.NextItems, ", ")))
		mw.strategyLabel.SetText(fmt.Sprintf("Strateji: %s", state.Recommendation.Strategy))
	}

	// Update Players
	mw.updatePlayerLists(state.Game.AllPlayers)
}

func (mw *MainWindow) updatePlayerLists(players []lcu.LivePlayer) {
	mw.teamOrderContainer.Objects = nil
	mw.teamChaosContainer.Objects = nil

	// Header for teams
	mw.teamOrderContainer.Add(widget.NewLabelWithStyle("ORDER TEAM (MAVİ)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	mw.teamChaosContainer.Add(widget.NewLabelWithStyle("CHAOS TEAM (KIRMIZI)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	for _, p := range players {
		card := mw.createPlayerCard(p)
		if p.Team == "ORDER" {
			mw.teamOrderContainer.Add(card)
		} else {
			mw.teamChaosContainer.Add(card)
		}
	}

	mw.teamOrderContainer.Refresh()
	mw.teamChaosContainer.Refresh()
}

func (mw *MainWindow) createPlayerCard(p lcu.LivePlayer) fyne.CanvasObject {
	// Basic Info
	nameLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)", p.SummonerName, p.ChampionName), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	kdaLabel := widget.NewLabel(fmt.Sprintf("KDA: %d/%d/%d | CS: %d", p.Scores.Kills, p.Scores.Deaths, p.Scores.Assists, p.Scores.CreepScore))

	// Items Summary (Just count or first few)
	itemCount := 0
	for _, item := range p.Items {
		if item.ItemID != 0 {
			itemCount++
		}
	}
	itemsLabel := widget.NewLabel(fmt.Sprintf("İtemler: %d", itemCount))

	// Detail Button
	detailBtn := widget.NewButton("Detay", func() {
		mw.showPlayerDetail(p)
	})

	// Layout
	content := container.NewVBox(
		nameLabel,
		kdaLabel,
		itemsLabel,
		detailBtn,
	)

	return widget.NewCard("", "", content)
}

func (mw *MainWindow) showPlayerDetail(p lcu.LivePlayer) {
	// Detailed Stats
	stats := fmt.Sprintf("Seviye: %d\nAltın: %.0f\nKDA: %d/%d/%d\nCS: %d\nWard: %.1f",
		p.Level, 0.0, // Gold might not be available for enemies usually, but struct has it? LiveActivePlayer has it. LivePlayer might not.
		p.Scores.Kills, p.Scores.Deaths, p.Scores.Assists, p.Scores.CreepScore, p.Scores.WardScore)

	// Items
	var itemNames []string
	for _, item := range p.Items {
		if item.DisplayName != "" {
			itemNames = append(itemNames, fmt.Sprintf("- %s (x%d)", item.DisplayName, item.Count))
		}
	}
	itemsStr := strings.Join(itemNames, "\n")

	// Runes
	runesStr := fmt.Sprintf("Keystone: %s\nPrimary: %s\nSecondary: %s",
		p.Runes.Keystone.DisplayName, p.Runes.PrimaryRuneTree.DisplayName, p.Runes.SecondaryRuneTree.DisplayName)

	content := container.NewVBox(
		widget.NewLabelWithStyle(fmt.Sprintf("%s - %s", p.SummonerName, p.ChampionName), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("İstatistikler", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(stats),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("İtemler", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(itemsStr),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Rünler", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(runesStr),
	)

	d := dialog.NewCustom("Oyuncu Detayı", "Kapat", container.NewVScroll(content), mw.window)
	d.Resize(fyne.NewSize(400, 500))
	d.Show()
}
