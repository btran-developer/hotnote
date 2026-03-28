package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Palette defines the color scheme for the TUI.
type Palette struct {
	// Standard colors (0-7)
	Black   tcell.Color // Black is the standard black color.
	Red     tcell.Color // Red is the standard red color.
	Green   tcell.Color // Green is the standard green color.
	Yellow  tcell.Color // Yellow is the standard yellow color.
	Blue    tcell.Color // Blue is the standard blue color.
	Magenta tcell.Color // Magenta is the standard magenta color.
	Cyan    tcell.Color // Cyan is the standard cyan color.
	White   tcell.Color // White is the standard white color.

	// Bright colors (8-15)
	BrightBlack   tcell.Color // BrightBlack is the bright black color.
	BrightRed     tcell.Color // BrightRed is the bright red color.
	BrightGreen   tcell.Color // BrightGreen is the bright green color.
	BrightYellow  tcell.Color // BrightYellow is the bright yellow color.
	BrightBlue    tcell.Color // BrightBlue is the bright blue color.
	BrightMagenta tcell.Color // BrightMagenta is the bright magenta color.
	BrightCyan    tcell.Color // BrightCyan is the bright cyan color.
	BrightWhite   tcell.Color // BrightWhite is the bright white color.

	// Semantic colors
	Error      tcell.Color // Error is the color for error messages.
	Warning    tcell.Color // Warning is the color for warnings.
	Success    tcell.Color // Success is the color for success messages.
	Info       tcell.Color // Info is the color for informational messages.
	Border     tcell.Color // Border is the color for borders.
	Title      tcell.Color // Title is the color for titles.
	Accent     tcell.Color // Accent is the color for accents.
	Dim        tcell.Color // Dim is the color for dimmed text.
	Background tcell.Color // Background is the background color.
	Primary    tcell.Color // Primary is the primary accent color.
}

// DefaultPalette is the default color palette.
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

// Apply applies the default palette to tview styles.
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
