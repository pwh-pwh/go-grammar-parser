package ui

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type MyTheme struct {
}

//go:embed ShangShouJianSongXianXiTi-2.ttf
var ShangShouJianSongXianXiTi []byte

var resourceShangShouJianSongXianXiTi2Ttf = &fyne.StaticResource{
	StaticName:    "ShangShouJianSongXianXiTi-2.ttf",
	StaticContent: ShangShouJianSongXianXiTi,
}

func (m *MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return resourceShangShouJianSongXianXiTi2Ttf
}

func (m *MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var _ fyne.Theme = (*MyTheme)(nil)
