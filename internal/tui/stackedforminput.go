package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StackedFormInput is a form input field with a stacked label and error message.
type StackedFormInput struct {
	*tview.Flex
	label    *tview.TextView
	input    *tview.InputField
	errorMsg *tview.TextView
}

// NewStackedFormInput creates a new stacked form input with label and placeholder.
func NewStackedFormInput(label, placeholder string) *StackedFormInput {
	s := &StackedFormInput{
		Flex:     tview.NewFlex().SetDirection(tview.FlexRow),
		label:    tview.NewTextView(),
		input:    tview.NewInputField(),
		errorMsg: tview.NewTextView(),
	}

	s.label.SetDynamicColors(true).SetWrap(false)
	s.errorMsg.SetDynamicColors(true).SetWrap(false)

	s.Flex.AddItem(s.label, 1, 0, false)
	s.Flex.AddItem(s.input, 3, 0, false)
	s.Flex.AddItem(s.errorMsg, 1, 0, false)

	if label != "" {
		s.SetLabel(label)
	}

	s.input.SetPlaceholder(placeholder)
	s.input.SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault))

	return s
}

// SetLabel sets the label text.
func (s *StackedFormInput) SetLabel(text string) *StackedFormInput {
	s.label.SetText(text)
	return s
}

// SetLabelAlign sets the label text alignment.
func (s *StackedFormInput) SetLabelAlign(align int) *StackedFormInput {
	s.label.SetTextAlign(align)
	return s
}

// SetLabelColor sets the label text color.
func (s *StackedFormInput) SetLabelColor(color tcell.Color) *StackedFormInput {
	s.label.SetTextColor(color)
	return s
}

// SetError sets the error message text.
func (s *StackedFormInput) SetError(text string) *StackedFormInput {
	s.errorMsg.SetText(text)
	return s
}

// SetErrorAlign sets the error message text alignment.
func (s *StackedFormInput) SetErrorAlign(align int) *StackedFormInput {
	s.errorMsg.SetTextAlign(align)
	return s
}

// SetErrorColor sets the error message text color.
func (s *StackedFormInput) SetErrorColor(color tcell.Color) *StackedFormInput {
	s.errorMsg.SetTextColor(color)
	return s
}

// SetText sets the input field text.
func (s *StackedFormInput) SetText(text string) *StackedFormInput {
	s.input.SetText(text)
	return s
}

// GetText returns the input field text.
func (s *StackedFormInput) GetText() string {
	return s.input.GetText()
}

// SetPlaceholder sets the input field placeholder.
func (s *StackedFormInput) SetPlaceholder(text string) *StackedFormInput {
	s.input.SetPlaceholder(text)
	return s
}

// SetFieldWidth sets the input field width.
func (s *StackedFormInput) SetFieldWidth(width int) *StackedFormInput {
	s.input.SetFieldWidth(width)
	return s
}

// SetDoneFunc sets the done handler function.
func (s *StackedFormInput) SetDoneFunc(fn func(tcell.Key)) *StackedFormInput {
	s.input.SetDoneFunc(fn)
	return s
}

// SetInputCapture sets the input capture function.
func (s *StackedFormInput) SetInputCapture(fn func(*tcell.EventKey) *tcell.EventKey) *StackedFormInput {
	s.input.SetInputCapture(fn)
	return s
}

// SetMaskCharacter sets the mask character for password input.
func (s *StackedFormInput) SetMaskCharacter(r rune) *StackedFormInput {
	s.input.SetMaskCharacter(r)
	return s
}

// SetBorder sets whether to show the border.
func (s *StackedFormInput) SetBorder(show bool) *StackedFormInput {
	s.input.SetBorder(show)
	return s
}

// SetTitle sets the input field title.
func (s *StackedFormInput) SetTitle(title string) *StackedFormInput {
	s.input.SetTitle(title)
	return s
}

// Focus focuses the input field.
func (s *StackedFormInput) Focus(delegate func(p tview.Primitive)) {
	delegate(s.input)
}

// GetInputField returns the underlying input field.
func (s *StackedFormInput) GetInputField() *tview.InputField {
	return s.input
}
