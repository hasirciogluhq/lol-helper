package main

import (
	"github.com/joho/godotenv"

	"lol-helper/internal/gui"
)

func main() {
	// .env ve .env.local dosyalarını yükle
	// .env.local varsa öncelikli olur
	godotenv.Load(".env.local", ".env")

	// GUI'yi başlat
	mainWindow := gui.NewMainWindow()
	mainWindow.Start()
}
