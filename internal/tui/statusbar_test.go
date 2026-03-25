package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusBar_NewStatusBar(t *testing.T) {
	sb := NewStatusBar()
	assert.NotNil(t, sb)
}

func TestStatusBar_SetStatus(t *testing.T) {
	sb := NewStatusBar()
	sb.SetStatus("test status")
	assert.Equal(t, "test status", sb.TextView.GetText(false))
}

func TestStatusBar_SetStatus_UpdatesContent(t *testing.T) {
	sb := NewStatusBar()
	sb.SetStatus("first")
	sb.SetStatus("second")
	assert.Equal(t, "second", sb.TextView.GetText(false))
}
