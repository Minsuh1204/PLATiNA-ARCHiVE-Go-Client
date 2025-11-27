package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

func (m MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return theme.DefaultTheme().Font(style)
	}
	if style.Bold {
		if style.Italic {
			return theme.DefaultTheme().Font(style) // Fallback for Bold+Italic if not available
		}
		return notoSansKRBold
	}
	if style.Italic {
		return theme.DefaultTheme().Font(style) // Fallback for Italic if not available
	}
	return notoSansKRRegular
}

func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var notoSansKRRegular, _ = fyne.LoadResourceFromPath("assets/Noto_Sans_KR/static/NotoSansKR-Regular.ttf")
var notoSansKRBold, _ = fyne.LoadResourceFromPath("assets/Noto_Sans_KR/static/NotoSansKR-Bold.ttf")
