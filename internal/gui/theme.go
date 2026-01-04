package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type LoLTheme struct{}

var _ fyne.Theme = (*LoLTheme)(nil)

func (t *LoLTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 1, G: 10, B: 19, A: 255} // Çok koyu mavi/siyah
	case theme.ColorNameForeground:
		return color.RGBA{R: 240, G: 230, B: 210, A: 255} // Altın/Bej
	case theme.ColorNamePrimary:
		return color.RGBA{R: 200, G: 170, B: 110, A: 255} // LoL Altın rengi
	case theme.ColorNameButton:
		return color.RGBA{R: 30, G: 35, B: 40, A: 255} // Koyu gri buton
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *LoLTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *LoLTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *LoLTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
