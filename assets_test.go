package hui

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssets(t *testing.T) {
	t.Run("ships the module and its types under a header", func(t *testing.T) {
		files := Assets()
		require.Len(t, files, 2)
		for name, source := range map[string]string{"hui.js": moduleJS, "hui.d.ts": moduleTypes} {
			asset, ok := files[name]
			require.True(t, ok, name)

			module, version, found := parseAssetHeader(asset.Content)
			require.True(t, found, name)
			assert.Equal(t, assetsModule, module)
			assert.Equal(t, asset.Version, version)
			assert.Equal(t, assetsModule, asset.Module)
			assert.True(t, strings.HasSuffix(string(asset.Content), source), name)
		}
	})

	t.Run("the types declare every export of the module", func(t *testing.T) {
		block := regexp.MustCompile(`(?s)export \{(.*?)\};`).FindStringSubmatch(moduleJS)
		require.Len(t, block, 2)
		var exported []string
		for item := range strings.SplitSeq(block[1], ",") {
			fields := strings.Fields(item)
			exported = append(exported, fields[len(fields)-1])
		}
		var declared []string
		for _, match := range regexp.MustCompile(`export (?:declare (?:const|function) (\w+)|default \w+|\{ \w+ as (\w+) \})`).FindAllStringSubmatch(moduleTypes, -1) {
			switch {
			case match[1] != "":
				declared = append(declared, match[1])
			case match[2] != "":
				declared = append(declared, match[2])
			default:
				declared = append(declared, "default")
			}
		}
		slices.Sort(exported)
		slices.Sort(declared)

		assert.Equal(t, exported, declared)
	})

	t.Run("writes the files and the list into a new directory", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "hui")

		require.NoError(t, WriteAssets(dir, Assets()))

		tree := testReadTree(t, dir)
		assert.Equal(t, string(Assets()["hui.js"].Content), tree["hui.js"])
		assert.Equal(t, string(Assets()["hui.d.ts"].Content), tree["hui.d.ts"])
		version := Assets()["hui.js"].Version
		assert.Equal(t, assetsListHead+"\n"+
			"hui.d.ts "+assetsModule+" "+version+"\n"+
			"hui.js "+assetsModule+" "+version+"\n", tree[assetsList])
	})

	t.Run("two runs give the same bytes and touch nothing", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, WriteAssets(dir, Assets()))
		before := testReadTree(t, dir)
		listInfo, err := os.Stat(filepath.Join(dir, assetsList))
		require.NoError(t, err)

		require.NoError(t, WriteAssets(dir, Assets()))

		assert.Equal(t, before, testReadTree(t, dir))
		listAfter, err := os.Stat(filepath.Join(dir, assetsList))
		require.NoError(t, err)
		assert.Equal(t, listInfo.ModTime(), listAfter.ModTime())
		changes, err := CheckAssets(dir, Assets())
		require.NoError(t, err)
		assert.Empty(t, changes)
	})

	t.Run("a wrong dir with the app sources fails and deletes nothing", func(t *testing.T) {
		dir := t.TempDir()
		sources := map[string]string{
			"generate.go": "package web\n",
			"index.ts":    gofakeit.Sentence(),
			"hui/hui.js":  string(Assets()["hui.js"].Content),
		}
		testWriteTree(t, dir, sources)

		err := WriteAssets(dir, Assets())

		require.ErrorIs(t, err, ErrForeignAssets)
		assert.Contains(t, err.Error(), "generate.go, index.ts")
		assert.Equal(t, sources, testReadTree(t, dir))
	})

	t.Run("a foreign file in a generated directory fails and stays", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, WriteAssets(dir, Assets()))
		own := gofakeit.Sentence()
		testWriteTree(t, dir, map[string]string{"app.js": own})
		before := testReadTree(t, dir)

		err := WriteAssets(dir, map[string]Asset{"hui.js": Assets()["hui.js"]})

		require.ErrorIs(t, err, ErrForeignAssets)
		assert.Contains(t, err.Error(), "app.js")
		assert.Equal(t, before, testReadTree(t, dir))
		changes, checkErr := CheckAssets(dir, Assets())
		require.NoError(t, checkErr)
		assert.Equal(t, []string{"app.js: foreign, not in .hui-assets"}, changes)
	})

	t.Run("hidden files are left alone", func(t *testing.T) {
		dir := t.TempDir()
		hidden := map[string]string{
			".DS_Store":       gofakeit.Sentence(),
			".gitkeep":        "",
			".cache/state":    gofakeit.Sentence(),
			"fonts/.DS_Store": gofakeit.Sentence(),
		}
		testWriteTree(t, dir, hidden)

		require.NoError(t, WriteAssets(dir, Assets()))
		changes, err := CheckAssets(dir, Assets())

		require.NoError(t, err)
		assert.Empty(t, changes)
		tree := testReadTree(t, dir)
		for name, content := range hidden {
			assert.Equal(t, content, tree[name], name)
		}
	})

	t.Run("removes only the listed files the new set does not have", func(t *testing.T) {
		dir := t.TempDir()
		files := Assets()
		files["old.js"] = Asset{Module: assetsModule, Version: "v0.1.2", Content: []byte(gofakeit.Sentence())}
		require.NoError(t, WriteAssets(dir, files))
		require.NoError(t, os.Remove(filepath.Join(dir, "hui.d.ts")))

		require.NoError(t, WriteAssets(dir, Assets()))

		assert.NoFileExists(t, filepath.Join(dir, "old.js"))
		assert.FileExists(t, filepath.Join(dir, "hui.d.ts"))
		assert.NotContains(t, testReadTree(t, dir)[assetsList], "old.js")
	})

	t.Run("a binary file without a header is listed and removed with its directory", func(t *testing.T) {
		dir := t.TempDir()
		font := []byte{0x77, 0x4f, 0x46, 0x32, 0x00, 0x01}
		files := Assets()
		files["fonts/x.woff2"] = CreateAsset("example.com/pdui", font)

		require.NoError(t, WriteAssets(dir, files))

		written, err := os.ReadFile(filepath.Join(dir, "fonts", "x.woff2"))
		require.NoError(t, err)
		assert.Equal(t, font, written)
		assert.Contains(t, testReadTree(t, dir)[assetsList], "fonts/x.woff2 example.com/pdui "+assetsDevel+"\n")

		require.NoError(t, WriteAssets(dir, Assets()))

		assert.NoDirExists(t, filepath.Join(dir, "fonts"))
	})

	t.Run("adopts a directory from v0.1.2 without a list", func(t *testing.T) {
		dir := t.TempDir()
		testWriteTree(t, dir, map[string]string{
			"hui.js": assetsHeadStart + assetsModule + " v0.1.2" + assetsHeadEnd + "\n" + gofakeit.Sentence(),
		})

		changes, err := CheckAssets(dir, Assets())
		require.NoError(t, err)
		assert.Equal(t, []string{
			"hui.d.ts: missing",
			"hui.js: differs, on disk from github.com/daarxwalker/hui v0.1.2",
			".hui-assets: missing",
		}, changes)

		require.NoError(t, WriteAssets(dir, Assets()))

		changes, err = CheckAssets(dir, Assets())
		require.NoError(t, err)
		assert.Empty(t, changes)
	})

	t.Run("a directory without a list and without headers is refused", func(t *testing.T) {
		dir := t.TempDir()
		testWriteTree(t, dir, map[string]string{"hui.js": gofakeit.Sentence()})

		assert.ErrorIs(t, WriteAssets(dir, Assets()), ErrForeignAssets)
	})

	t.Run("check reports every difference with its origin", func(t *testing.T) {
		dir := t.TempDir()
		files := Assets()
		files["hui.js"] = Asset{Module: assetsModule, Version: "v0.1.2", Content: []byte(gofakeit.Sentence())}
		files["old.js"] = Asset{Module: assetsModule, Version: "v0.1.2", Content: []byte(gofakeit.Sentence())}
		delete(files, "hui.d.ts")
		require.NoError(t, WriteAssets(dir, files))

		changes, err := CheckAssets(dir, Assets())

		require.NoError(t, err)
		assert.Equal(t, []string{
			"hui.d.ts: missing",
			"hui.js: differs, on disk from github.com/daarxwalker/hui v0.1.2",
			"old.js: stale, from github.com/daarxwalker/hui v0.1.2",
			".hui-assets: differs",
		}, changes)
	})

	t.Run("check reports a missing directory", func(t *testing.T) {
		changes, err := CheckAssets(filepath.Join(t.TempDir(), "hui"), Assets())

		require.NoError(t, err)
		assert.Equal(t, []string{"hui.d.ts: missing", "hui.js: missing", ".hui-assets: missing"}, changes)
	})

	t.Run("rejects names it could not list or would hide", func(t *testing.T) {
		for _, name := range []string{"../escape.js", "/abs.js", ".hidden.js", "fonts/.x.woff2", "with space.js", assetsList} {
			files := Assets()
			files[name] = CreateAsset(assetsModule, []byte(gofakeit.Sentence()))

			assert.Error(t, WriteAssets(t.TempDir(), files), name)
		}
	})

	t.Run("a list pointing outside the directory is refused", func(t *testing.T) {
		dir := t.TempDir()
		testWriteTree(t, dir, map[string]string{assetsList: "../../etc/passwd " + assetsModule + " v0.1.2\n"})

		assert.Error(t, WriteAssets(dir, Assets()))
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
