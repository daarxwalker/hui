package headlessui

import (
	"context"
	"io"
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
