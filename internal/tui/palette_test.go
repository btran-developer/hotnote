package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestPalette_StandardColors(t *testing.T) {
	p := DefaultPalette

	tests := []struct {
		name  string
		got   tcell.Color
		index int
	}{
		{"Black", p.Black, 0},
		{"Red", p.Red, 1},
		{"Green", p.Green, 2},
		{"Yellow", p.Yellow, 3},
		{"Blue", p.Blue, 4},
		{"Magenta", p.Magenta, 5},
		{"Cyan", p.Cyan, 6},
		{"White", p.White, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tcell.PaletteColor(tt.index)
			if tt.got != expected {
				t.Errorf("got %v, want %v", tt.got, expected)
			}
		})
	}
}

func TestPalette_BrightColors(t *testing.T) {
	p := DefaultPalette

	tests := []struct {
		name  string
		got   tcell.Color
		index int
	}{
		{"BrightBlack", p.BrightBlack, 8},
		{"BrightRed", p.BrightRed, 9},
		{"BrightGreen", p.BrightGreen, 10},
		{"BrightYellow", p.BrightYellow, 11},
		{"BrightBlue", p.BrightBlue, 12},
		{"BrightMagenta", p.BrightMagenta, 13},
		{"BrightCyan", p.BrightCyan, 14},
		{"BrightWhite", p.BrightWhite, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tcell.PaletteColor(tt.index)
			if tt.got != expected {
				t.Errorf("got %v, want %v", tt.got, expected)
			}
		})
	}
}

func TestPalette_SemanticAliases(t *testing.T) {
	p := DefaultPalette

	tests := []struct {
		name  string
		got   tcell.Color
		index int
	}{
		{"Error", p.Error, 1},
		{"Warning", p.Warning, 3},
		{"Success", p.Success, 2},
		{"Info", p.Info, 6},
		{"Border", p.Border, 8},
		{"Title", p.Title, 7},
		{"Accent", p.Accent, 5},
		{"Dim", p.Dim, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tcell.PaletteColor(tt.index)
			if tt.got != expected {
				t.Errorf("got %v, want %v", tt.got, expected)
			}
		})
	}
}

func TestApply(t *testing.T) {
	Apply()

	styles := tview.Styles
	defaultColor := tcell.ColorDefault

	check := func(name string, got, want tcell.Color) {
		if got != want {
			t.Errorf("%s: got %v, want %v", name, got, want)
		}
	}

	check("PrimitiveBackgroundColor", styles.PrimitiveBackgroundColor, defaultColor)
	check("ContrastBackgroundColor", styles.ContrastBackgroundColor, defaultColor)
	check("MoreContrastBackgroundColor", styles.MoreContrastBackgroundColor, defaultColor)
	check("BorderColor", styles.BorderColor, defaultColor)
	check("TitleColor", styles.TitleColor, defaultColor)
	check("GraphicsColor", styles.GraphicsColor, defaultColor)
	check("PrimaryTextColor", styles.PrimaryTextColor, defaultColor)
	check("SecondaryTextColor", styles.SecondaryTextColor, defaultColor)
	check("TertiaryTextColor", styles.TertiaryTextColor, defaultColor)
	check("InverseTextColor", styles.InverseTextColor, defaultColor)
	check("ContrastSecondaryTextColor", styles.ContrastSecondaryTextColor, defaultColor)
}
