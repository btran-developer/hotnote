package fsutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtomicWrite(t *testing.T) {
	dir, err := os.MkdirTemp("", "fsutil-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "test.txt")
	content := []byte("hello world")

	err = AtomicWrite(path, content, 0644)
	require.NoError(t, err)

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, got)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
}

func TestAtomicWriteExclusive(t *testing.T) {
	dir, err := os.MkdirTemp("", "fsutil-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "test.txt")
	content := []byte("hello world")

	err = AtomicWriteExclusive(path, content, 0644)
	require.NoError(t, err)

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, got)

	err = AtomicWriteExclusive(path, content, 0644)
	assert.Error(t, err)
}

func TestAtomicWrite_CreatesDirectory(t *testing.T) {
	dir, err := os.MkdirTemp("", "fsutil-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "nested", "dir", "test.txt")
	content := []byte("nested")

	err = AtomicWrite(path, content, 0644)
	require.NoError(t, err)

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, got)
}

func TestAtomicWriteExclusive_CleanupOnError(t *testing.T) {
	dir, err := os.MkdirTemp("", "fsutil-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "test.txt")
	err = AtomicWriteExclusive(path, []byte("first"), 0644)
	require.NoError(t, err)

	err = AtomicWriteExclusive(path, []byte("second"), 0644)
	assert.Error(t, err)

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, []byte("first"), got)
}
