package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Palette struct {
	Black   tcell.Color
	Red     tcell.Color
	Green   tcell.Color
	Yellow  tcell.Color
	Blue    tcell.Color
	Magenta tcell.Color
	Cyan    tcell.Color
	White   tcell.Color

	BrightBlack   tcell.Color
	BrightRed     tcell.Color
	BrightGreen   tcell.Color
	BrightYellow  tcell.Color
	BrightBlue    tcell.Color
	BrightMagenta tcell.Color
	BrightCyan    tcell.Color
	BrightWhite   tcell.Color

	Error      tcell.Color
	Warning    tcell.Color
	Success    tcell.Color
	Info       tcell.Color
	Border     tcell.Color
	Title      tcell.Color
	Accent     tcell.Color
	Dim        tcell.Color
	Background tcell.Color
	Primary    tcell.Color
}

var DefaultPalette = Palette{
	Black:   tcell.PaletteColor(0),
	Red:     tcell.PaletteColor(1),
	Green:   tcell.PaletteColor(2),
	Yellow:  tcell.PaletteColor(3),
	Blue:    tcell.PaletteColor(4),
	Magenta: tcell.PaletteColor(5),
	Cyan:    tcell.PaletteColor(6),
	White:   tcell.PaletteColor(7),

	BrightBlack:   tcell.PaletteColor(8),
	BrightRed:     tcell.PaletteColor(9),
	BrightGreen:   tcell.PaletteColor(10),
	BrightYellow:  tcell.PaletteColor(11),
	BrightBlue:    tcell.PaletteColor(12),
	BrightMagenta: tcell.PaletteColor(13),
	BrightCyan:    tcell.PaletteColor(14),
	BrightWhite:   tcell.PaletteColor(15),

	Error:      tcell.PaletteColor(1),
	Warning:    tcell.PaletteColor(3),
	Success:    tcell.PaletteColor(2),
	Info:       tcell.PaletteColor(6),
	Border:     tcell.PaletteColor(8),
	Title:      tcell.PaletteColor(7),
	Accent:     tcell.PaletteColor(5),
	Dim:        tcell.PaletteColor(8),
	Background: tcell.PaletteColor(0),
	Primary:    tcell.PaletteColor(4),
}

func Apply() {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorDefault,
		ContrastBackgroundColor:     tcell.ColorDefault,
		MoreContrastBackgroundColor: tcell.ColorDefault,
		BorderColor:                 tcell.ColorDefault,
		TitleColor:                  tcell.ColorDefault,
		GraphicsColor:               tcell.ColorDefault,
		PrimaryTextColor:            tcell.ColorDefault,
		SecondaryTextColor:          tcell.ColorDefault,
		TertiaryTextColor:           tcell.ColorDefault,
		InverseTextColor:            tcell.ColorDefault,
		ContrastSecondaryTextColor:  tcell.ColorDefault,
	}
}
