package hui

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"
)

func testRender(t *testing.T, c templ.Component, children ...templ.Component) string {
	t.Helper()
	return testRenderIn(t, context.Background(), c, children...)
}

func testRenderIn(t *testing.T, ctx context.Context, c templ.Component, children ...templ.Component) string {
	t.Helper()
	if len(children) > 0 {
		c = testWithChildren(c, children...)
	}
	var sb strings.Builder
	require.NoError(t, c.Render(ctx, &sb))
	return sb.String()
}

func testWithChildren(c templ.Component, children ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(templ.WithChildren(ctx, templ.Join(children...)), w)
	})
}

func testCount(html, needle string) int {
	return strings.Count(html, needle)
}

func testReadTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := make(map[string]string)
	err := filepath.WalkDir(dir, func(name string, entry fs.DirEntry, walkErr error) error {
		require.NoError(t, walkErr)
		if entry.IsDir() {
			return nil
		}
		content, readErr := os.ReadFile(name)
		require.NoError(t, readErr)
		relative, relErr := filepath.Rel(dir, name)
		require.NoError(t, relErr)
		files[filepath.ToSlash(relative)] = string(content)
		return nil
	})
	require.NoError(t, err)
	return files
}

func testWriteTree(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		target := filepath.Join(dir, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(target), 0o755))
		require.NoError(t, os.WriteFile(target, []byte(content), 0o644))
	}
}
