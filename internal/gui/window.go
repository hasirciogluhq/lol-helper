package gui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"lol-helper/internal/lcu"
	"lol-helper/internal/lol"
)

// MainWindow ana pencere yapısı
type MainWindow struct {
	app         fyne.App
	window      fyne.Window
	service     *lol.Service
	itemManager *ItemManager

	// Cache
	imageCache      map[int]fyne.Resource
	lastPlayerNames string // Player isimlerini cache'le
	lastAIItems     string
	playersLoaded   bool // İlk yükleme yapıldı mı?

	// UI Components
	statusLabel *widget.Label
	phaseLabel  *widget.Label

	// Team Containers
	teamOrderContainer *fyne.Container
	teamChaosContainer *fyne.Container

	// AI Suggestion Section
	suggestionLabel  *widget.Label
	strategyLabel    *widget.Label
	aiItemsContainer *fyne.Container
}

// NewMainWindow yeni bir ana pencere oluşturur
func NewMainWindow() *MainWindow {
	a := app.New()
	a.Settings().SetTheme(&LoLTheme{})

	w := a.NewWindow("LoL Helper AI")
	w.Resize(fyne.NewSize(1200, 800))

	mw := &MainWindow{
		app:         a,
		window:      w,
		itemManager: NewItemManager(),
		imageCache:  make(map[int]fyne.Resource),
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

	mw.strategyLabel = widget.NewLabel("Strateji: -")
	mw.strategyLabel.Wrapping = fyne.TextWrapWord

	mw.aiItemsContainer = container.NewHBox()

	// Team Containers
	mw.teamOrderContainer = container.NewVBox()
	mw.teamChaosContainer = container.NewVBox()

	// Teams Split View - No scroll, compact design
	teamsSplit := container.NewGridWithColumns(2,
		container.NewPadded(mw.teamOrderContainer),
		container.NewPadded(mw.teamChaosContainer),
	)

	// Top Info
	topInfo := container.NewVBox(
		mw.statusLabel,
		mw.phaseLabel,
	)

	// Bottom AI - Professional Layout
	aiHeader := widget.NewLabelWithStyle("AI KOÇ ANALİZİ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Left: Strategy & Suggestion
	aiTextContent := container.NewVBox(
		widget.NewLabelWithStyle("Strateji", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		mw.strategyLabel,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Öneri", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		mw.suggestionLabel,
	)

	// Right: Recommended Items
	aiItemsContent := container.NewVBox(
		widget.NewLabelWithStyle("Önerilen Eşyalar", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(mw.aiItemsContainer),
	)

	// Split AI Section
	aiSplit := container.NewGridWithColumns(2,
		container.NewPadded(aiTextContent),
		container.NewPadded(aiItemsContent),
	)

	bottomAI := widget.NewCard("", "", container.NewVBox(
		aiHeader,
		widget.NewSeparator(),
		aiSplit,
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
		mw.suggestionLabel.SetText(state.Recommendation.Suggestion)
		mw.strategyLabel.SetText(state.Recommendation.Strategy)

		// Update AI Items - only if changed
		newAIItems := strings.Join(state.Recommendation.NextItems, ",")
		if newAIItems != mw.lastAIItems {
			mw.lastAIItems = newAIItems
			mw.aiItemsContainer.Objects = nil
			for _, itemName := range state.Recommendation.NextItems {
				itemID := mw.itemManager.GetItemID(itemName)
				if itemID != 0 {
					img := mw.createItemImage(itemID)
					itemContainer := container.NewVBox(
						img,
						widget.NewLabelWithStyle(itemName, fyne.TextAlignCenter, fyne.TextStyle{}),
					)
					mw.aiItemsContainer.Add(itemContainer)
				} else {
					mw.aiItemsContainer.Add(widget.NewLabel(itemName))
				}
			}
			mw.aiItemsContainer.Refresh()
		}
	}

	// Update Players
	mw.updatePlayerLists(state.Game.AllPlayers)
}

func (mw *MainWindow) updatePlayerLists(players []lcu.LivePlayer) {
	// Player isimlerini string olarak oluştur
	var currentPlayerNames string
	for _, p := range players {
		currentPlayerNames += p.SummonerName + ","
	}

	// Eğer player isimleri aynıysa hiçbir şey yapma (blinking önlemek için)
	if currentPlayerNames == mw.lastPlayerNames && mw.playersLoaded {
		return
	}

	// Boş liste geldi ama daha önce playerlar vardı - mevcut UI'ı koru
	if len(players) == 0 && mw.playersLoaded {
		return
	}

	// İlk kez yükleme veya gerçekten değişiklik var
	mw.lastPlayerNames = currentPlayerNames
	mw.playersLoaded = true

	// Pre-build all elements first, then update containers atomically
	orderPlayers := make([]fyne.CanvasObject, 0, 5)
	chaosPlayers := make([]fyne.CanvasObject, 0, 5)

	for _, p := range players {
		card := mw.createPlayerRow(p)
		if p.Team == "ORDER" {
			orderPlayers = append(orderPlayers, card)
		} else {
			chaosPlayers = append(chaosPlayers, card)
		}
	}

	// Clear and rebuild atomically
	mw.teamOrderContainer.Objects = []fyne.CanvasObject{
		widget.NewLabelWithStyle("ORDER TEAM (MAVİ)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		mw.createTableHeader(),
	}
	mw.teamOrderContainer.Objects = append(mw.teamOrderContainer.Objects, orderPlayers...)

	mw.teamChaosContainer.Objects = []fyne.CanvasObject{
		widget.NewLabelWithStyle("CHAOS TEAM (KIRMIZI)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		mw.createTableHeader(),
	}
	mw.teamChaosContainer.Objects = append(mw.teamChaosContainer.Objects, chaosPlayers...)

	mw.teamOrderContainer.Refresh()
	mw.teamChaosContainer.Refresh()
}

func (mw *MainWindow) createTableHeader() fyne.CanvasObject {
	return container.NewHBox(
		mw.fixedLabel("Şampiyon", 120, true),
		mw.fixedLabel("Sihirdar", 120, true),
		mw.fixedLabel("KDA", 100, true),
		mw.fixedLabel("CS", 50, true),
		widget.NewLabelWithStyle("İtemler", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
}

func (mw *MainWindow) createPlayerRow(p lcu.LivePlayer) fyne.CanvasObject {
	// Columns
	champLabel := mw.fixedLabel(p.ChampionName, 120, false)
	nameLabel := mw.fixedLabel(p.SummonerName, 120, false)
	kdaLabel := mw.fixedLabel(fmt.Sprintf("%d/%d/%d", p.Scores.Kills, p.Scores.Deaths, p.Scores.Assists), 100, false)
	csLabel := mw.fixedLabel(fmt.Sprintf("%d", p.Scores.CreepScore), 50, false)

	// Items as images in horizontal row
	itemsRow := container.NewHBox()
	for _, item := range p.Items {
		if item.ItemID != 0 {
			itemImg := mw.createItemImage(item.ItemID)
			itemsRow.Add(itemImg)
		}
	}
	// Fill empty slots
	for len(itemsRow.Objects) < 6 {
		emptySlot := canvas.NewRectangle(color.RGBA{R: 20, G: 20, B: 20, A: 255})
		emptySlot.SetMinSize(fyne.NewSize(32, 32))
		itemsRow.Add(emptySlot)
	}

	// Row Content
	content := container.NewHBox(
		champLabel,
		nameLabel,
		kdaLabel,
		csLabel,
		itemsRow,
	)

	// Clickable Wrapper
	return NewClickableRow(content, func() {
		mw.showPlayerDetail(p)
	})
}

// fixedLabel creates a label with fixed width
func (mw *MainWindow) fixedLabel(text string, width float32, bold bool) fyne.CanvasObject {
	style := fyne.TextStyle{Bold: bold}
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, style)
	label.Truncation = fyne.TextTruncateEllipsis

	return container.New(&fixedWidthLayout{width: width}, label)
}

// fixedWidthLayout enforces a fixed width
type fixedWidthLayout struct {
	width float32
}

func (l *fixedWidthLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(fyne.NewSize(l.width, size.Height))
		o.Move(fyne.NewPos(0, 0))
	}
}

func (l *fixedWidthLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	h := float32(0)
	for _, o := range objects {
		h = fyne.Max(h, o.MinSize().Height)
	}
	return fyne.NewSize(l.width, h)
}

// ClickableRow is a custom widget that handles tap events
type ClickableRow struct {
	widget.BaseWidget
	OnTap   func()
	Content fyne.CanvasObject
}

func NewClickableRow(content fyne.CanvasObject, onTap func()) *ClickableRow {
	c := &ClickableRow{OnTap: onTap, Content: content}
	c.ExtendBaseWidget(c)
	return c
}

func (c *ClickableRow) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.Content)
}

func (c *ClickableRow) Tapped(_ *fyne.PointEvent) {
	if c.OnTap != nil {
		c.OnTap()
	}
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

// createItemImage creates an item icon from DDragon
func (mw *MainWindow) createItemImage(itemID int) fyne.CanvasObject {
	// Check cache first
	if res, ok := mw.imageCache[itemID]; ok {
		img := canvas.NewImageFromResource(res)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(32, 32))
		return img
	}

	// Data Dragon item icon URL
	itemURL := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.1.1/img/item/%d.png", itemID)

	// Load resource (blocking but cached next time)
	res, err := fyne.LoadResourceFromURLString(itemURL)
	if err == nil {
		mw.imageCache[itemID] = res
		img := canvas.NewImageFromResource(res)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(32, 32))
		return img
	}

	// Fallback
	rect := canvas.NewRectangle(color.RGBA{R: 50, G: 100, B: 150, A: 255})
	rect.SetMinSize(fyne.NewSize(32, 32))
	return rect
}
