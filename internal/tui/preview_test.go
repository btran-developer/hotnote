package tui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreviewPane_NewPreviewPane(t *testing.T) {
	p := NewPreviewPane()
	assert.NotNil(t, p)
	assert.Equal(t, "", p.currentFile)
}

func TestPreviewPane_LoadNote(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("# Test Note\nContent here")
	assert.NoError(t, err)
	tmpFile.Close()

	p := NewPreviewPane()
	err = p.LoadNote(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, tmpFile.Name(), p.currentFile)
}

func TestPreviewPane_LoadNote_FileNotFound(t *testing.T) {
	p := NewPreviewPane()
	err := p.LoadNote("/nonexistent/file.md")
	assert.Error(t, err)
}

func TestPreviewPane_Clear(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("content")
	assert.NoError(t, err)
	tmpFile.Close()

	p := NewPreviewPane()
	p.LoadNote(tmpFile.Name())
	assert.Equal(t, tmpFile.Name(), p.currentFile)

	p.Clear()
	assert.Equal(t, "", p.currentFile)
}

func TestPreviewPane_GetCurrentFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.md")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("content")
	assert.NoError(t, err)
	tmpFile.Close()

	p := NewPreviewPane()
	assert.Equal(t, "", p.GetCurrentFile())

	p.LoadNote(tmpFile.Name())
	assert.Equal(t, tmpFile.Name(), p.GetCurrentFile())
}
