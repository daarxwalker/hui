package hui

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssets(t *testing.T) {
	t.Run("ships the ES module under a header", func(t *testing.T) {
		js, ok := Assets()["hui.js"]
		require.True(t, ok)

		header, found := assetHeader(js)
		require.True(t, found)
		assert.True(t, strings.HasPrefix(header, "github.com/daarxwalker/hui "))
		assert.True(t, strings.HasSuffix(string(js), moduleJS))
		assert.Contains(t, moduleJS, "export ")
	})

	t.Run("writes into a new directory", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "hui")

		require.NoError(t, WriteAssets(dir, Assets()))

		written, err := os.ReadFile(filepath.Join(dir, "hui.js"))
		require.NoError(t, err)
		assert.Equal(t, Assets()["hui.js"], written)
	})

	t.Run("writing again changes nothing", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, WriteAssets(dir, Assets()))
		before, err := os.Stat(filepath.Join(dir, "hui.js"))
		require.NoError(t, err)

		require.NoError(t, WriteAssets(dir, Assets()))

		after, err := os.Stat(filepath.Join(dir, "hui.js"))
		require.NoError(t, err)
		assert.Equal(t, before.ModTime(), after.ModTime())
		changes, err := CheckAssets(dir, Assets())
		require.NoError(t, err)
		assert.Empty(t, changes)
	})

	t.Run("removes files the new version does not have", func(t *testing.T) {
		dir := t.TempDir()
		previous := Assets()
		previous["old.js"] = []byte(AssetHeader(assetsModule) + gofakeit.Sentence())
		previous["fonts/old.woff2"] = []byte(gofakeit.Sentence())
		require.NoError(t, WriteAssets(dir, previous))

		require.NoError(t, WriteAssets(dir, Assets()))

		assert.NoFileExists(t, filepath.Join(dir, "old.js"))
		assert.NoDirExists(t, filepath.Join(dir, "fonts"))
		assert.FileExists(t, filepath.Join(dir, "hui.js"))
	})

	t.Run("writes nested files", func(t *testing.T) {
		dir := t.TempDir()
		files := Assets()
		files["fonts/plex.woff2"] = []byte(gofakeit.Sentence())

		require.NoError(t, WriteAssets(dir, files))

		assert.FileExists(t, filepath.Join(dir, "fonts", "plex.woff2"))
	})

	t.Run("refuses a directory it did not generate", func(t *testing.T) {
		dir := t.TempDir()
		own := []byte(gofakeit.Sentence())
		require.NoError(t, os.WriteFile(filepath.Join(dir, "app.js"), own, 0o644))

		err := WriteAssets(dir, Assets())

		require.ErrorIs(t, err, ErrForeignAssets)
		kept, readErr := os.ReadFile(filepath.Join(dir, "app.js"))
		require.NoError(t, readErr)
		assert.Equal(t, own, kept)
		assert.NoFileExists(t, filepath.Join(dir, "hui.js"))
	})

	t.Run("rejects names outside the directory", func(t *testing.T) {
		files := Assets()
		files["../escape.js"] = []byte(gofakeit.Sentence())

		assert.Error(t, WriteAssets(t.TempDir(), files))
	})

	t.Run("check reports every difference with its origin", func(t *testing.T) {
		dir := t.TempDir()
		previous := Assets()
		previous["hui.js"] = []byte(assetsHeadStart + assetsModule + " v0.0.1" + assetsHeadEnd + "\n" + gofakeit.Sentence())
		previous["old.js"] = []byte(gofakeit.Sentence())
		require.NoError(t, WriteAssets(dir, previous))

		changes, err := CheckAssets(dir, Assets())

		require.NoError(t, err)
		assert.Equal(t, []string{
			"hui.js: differs, on disk from github.com/daarxwalker/hui v0.0.1",
			"old.js: stale, on disk not generated",
		}, changes)
	})

	t.Run("check reports a missing directory", func(t *testing.T) {
		changes, err := CheckAssets(filepath.Join(t.TempDir(), "hui"), Assets())

		require.NoError(t, err)
		assert.Equal(t, []string{"hui.js: missing"}, changes)
	})

	t.Run("the version comes from a replacement", func(t *testing.T) {
		required := debug.Module{Path: assetsModule, Version: "v1.2.3"}
		local := debug.Module{Path: assetsModule, Version: "v1.2.3", Replace: &debug.Module{Path: "../hui"}}
		fork := debug.Module{Path: assetsModule, Version: "v1.2.3", Replace: &debug.Module{Path: "example.com/hui", Version: "v1.2.4"}}

		assert.Equal(t, "v1.2.3", buildVersion(required))
		assert.Equal(t, assetsDevel, buildVersion(local))
		assert.Equal(t, "v1.2.4", buildVersion(fork))
	})
}
