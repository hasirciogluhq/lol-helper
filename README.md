# LoL Helper

League of Legends için Go dilinde geliştirilmiş bir yardımcı uygulama.

## Özellikler

- 🎮 **Champion Seçimi**: Oyundaki tüm popüler championları destekler
- 🔮 **Rün Önerileri**: Seçilen champion ve oyun stiline göre otomatik rün önerileri
  - Agresif/Defansif stil seçimi
  - Role özel rün sayfaları (ADC, Support, Mid, Jungle, Top)
- 🛡️ **İtem Önerileri**: Her champion için önerilen item build'leri
- 👥 **Oyuncu Bilgileri**: Oyun içi oyuncu listesi ve detayları
- 🎨 **Modern UI**: LoL temalı koyu tema ile şık arayüz

## Kurulum

### Gereksinimler

- Go 1.23 veya üzeri
- Fyne GUI kütüphanesi bağımlılıkları

### macOS için bağımlılıklar

```bash
# Xcode command line tools (zaten yüklüyse atla)
xcode-select --install
```

### Projeyi Kurma

```bash
# Projeyi klonlayın
cd lol-helper

# Bağımlılıkları indirin
go mod download

# Uygulamayı çalıştırın
go run main.go
```

## Derleme

### macOS için

```bash
# Executable oluştur
go build -o lol-helper main.go

# veya optimize edilmiş versiyon
go build -ldflags="-s -w" -o lol-helper main.go

# Çalıştır
./lol-helper
```

### Windows için (macOS'ta cross-compile)

```bash
GOOS=windows GOARCH=amd64 go build -o lol-helper.exe main.go
```

### Linux için

```bash
GOOS=linux GOARCH=amd64 go build -o lol-helper main.go
```

## Kullanım

1. Uygulamayı başlatın
2. Sol panelden bir champion seçin
3. Agresif veya defansif rün stili seçin
4. Önerilen rünler ve itemler otomatik olarak görüntülenecektir
5. Sağ panelde oyun içi oyuncu bilgilerini görün (demo verisi)

## Proje Yapısı

```text
lol-helper/
├── main.go                 # Ana uygulama entry point
├── go.mod                  # Go modül dosyası
├── internal/
│   ├── lol/               # LoL oyun mantığı
│   │   ├── models.go      # Veri modelleri (Champion, Rune, Item, vb.)
│   │   ├── data.go        # Statik veri (champions, runes)
│   │   └── service.go     # LoL servisi (API çağrıları, veri yönetimi)
│   └── gui/               # GUI katmanı
│       ├── window.go      # Ana pencere ve UI bileşenleri
│       └── theme.go       # Özel LoL teması
└── README.md
```

## Geliştirme Notları

### Gelecek Özellikler

- [ ] Gerçek Riot API entegrasyonu
- [ ] LCU (League Client API) bağlantısı
- [ ] Canlı oyun verisi takibi
- [ ] Counter pick önerileri
- [ ] Detaylı istatistikler
- [ ] Oyun içi overlay modu

### Teknik Detaylar

- **GUI Framework**: Fyne v2.5.0 - Cross-platform Go GUI toolkit
- **Mimari**: Clean architecture ile katmanlı yapı
- **Concurrency**: Goroutine ile arka plan görevleri
- **Tema**: Özel LoL temalı dark mode

## Katkıda Bulunma

1. Fork yapın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Commit yapın (`git commit -m 'feat: add amazing feature'`)
4. Push edin (`git push origin feature/amazing-feature`)
5. Pull Request açın

## Lisans

Bu proje eğitim amaçlıdır. League of Legends, Riot Games, Inc.'in tescilli markasıdır.

## İletişim

Sorular ve öneriler için issue açabilirsiniz.

---

**Not**: Bu uygulama demo amaçlıdır. Gerçek oyun verisi için Riot Games API key'e ihtiyaç vardır.
