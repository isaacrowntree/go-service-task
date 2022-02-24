package slicer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesPerChunk(t *testing.T) {
	t.Run("division by zero fix", func(t *testing.T) {
		files := filesPerChunk([]string{}, 1)

		assert.Equal(t, 1, files)
	})

	t.Run("1 worker processes both files", func(t *testing.T) {
		files := filesPerChunk([]string{"file1", "file2"}, 1)

		assert.Equal(t, 2, files)
	})

	t.Run("2 files processed by 4 workers", func(t *testing.T) {
		files := filesPerChunk([]string{"file1", "file2"}, 4)

		assert.Equal(t, 1, files)
	})

	t.Run("12 files processed by 4 workers", func(t *testing.T) {
		files := filesPerChunk([]string{"file", "file", "file", "file", "file", "file", "file", "file", "file", "file", "file", "file"}, 4)

		assert.Equal(t, 3, files)
	})
}

func TestChunkSlice(t *testing.T) {
	t.Run("no files equals nil array", func(t *testing.T) {
		files := ChunkSlice([]string{}, 4)

		assert.Empty(t, files)
	})

	t.Run("with files equals a non-empty array", func(t *testing.T) {
		files := ChunkSlice([]string{"file1", "file2"}, 1)

		assert.Equal(t, ([][]string{{"file1", "file2"}}), files)
	})

	t.Run("sends 2 files to two workers", func(t *testing.T) {
		files := ChunkSlice([]string{"file1", "file2"}, 2)

		assert.Equal(t, ([][]string{{"file1"}, {"file2"}}), files)
	})
}
