package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StackedFormInput struct {
	*tview.Flex
	label    *tview.TextView
	input    *tview.InputField
	errorMsg *tview.TextView
}

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

func (s *StackedFormInput) SetLabel(text string) *StackedFormInput {
	s.label.SetText(text)
	return s
}

func (s *StackedFormInput) SetLabelAlign(align int) *StackedFormInput {
	s.label.SetTextAlign(align)
	return s
}

func (s *StackedFormInput) SetLabelColor(color tcell.Color) *StackedFormInput {
	s.label.SetTextColor(color)
	return s
}

func (s *StackedFormInput) SetError(text string) *StackedFormInput {
	s.errorMsg.SetText(text)
	return s
}

func (s *StackedFormInput) SetErrorAlign(align int) *StackedFormInput {
	s.errorMsg.SetTextAlign(align)
	return s
}

func (s *StackedFormInput) SetErrorColor(color tcell.Color) *StackedFormInput {
	s.errorMsg.SetTextColor(color)
	return s
}

func (s *StackedFormInput) SetText(text string) *StackedFormInput {
	s.input.SetText(text)
	return s
}

func (s *StackedFormInput) GetText() string {
	return s.input.GetText()
}

func (s *StackedFormInput) SetPlaceholder(text string) *StackedFormInput {
	s.input.SetPlaceholder(text)
	return s
}

func (s *StackedFormInput) SetFieldWidth(width int) *StackedFormInput {
	s.input.SetFieldWidth(width)
	return s
}

func (s *StackedFormInput) SetDoneFunc(fn func(tcell.Key)) *StackedFormInput {
	s.input.SetDoneFunc(fn)
	return s
}

func (s *StackedFormInput) SetInputCapture(fn func(*tcell.EventKey) *tcell.EventKey) *StackedFormInput {
	s.input.SetInputCapture(fn)
	return s
}

func (s *StackedFormInput) SetMaskCharacter(r rune) *StackedFormInput {
	s.input.SetMaskCharacter(r)
	return s
}

func (s *StackedFormInput) SetBorder(show bool) *StackedFormInput {
	s.input.SetBorder(show)
	return s
}

func (s *StackedFormInput) SetTitle(title string) *StackedFormInput {
	s.input.SetTitle(title)
	return s
}

func (s *StackedFormInput) Focus(delegate func(p tview.Primitive)) {
	delegate(s.input)
}

func (s *StackedFormInput) GetInputField() *tview.InputField {
	return s.input
}
