package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// all the available themes

type dark_plus struct{}

var _ fyne.Theme = (*dark_plus)(nil)

func (m dark_plus) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m dark_plus) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m dark_plus) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m dark_plus) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

type white_theme struct{}

func (white_theme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	case theme.ColorNameButton:
		return color.Alpha16{A: 0x0}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0xe5, G: 0xe5, B: 0xe5, A: 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x42}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0x7f}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xff}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xf}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x19}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x19}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0xff}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x99}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x33}
	default:
		return theme.DefaultTheme().Color(c, v)
	}
}

func (white_theme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s)
	}
	if s.Bold {
		if s.Italic {
			return theme.DefaultTheme().Font(s)
		}
		return theme.DefaultTheme().Font(s)
	}
	if s.Italic {
		return theme.DefaultTheme().Font(s)
	}
	return theme.DefaultTheme().Font(s)
}

func (white_theme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (white_theme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(s)
	}
}
