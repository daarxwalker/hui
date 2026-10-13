package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScripts(t *testing.T) {
	t.Run("ships every plugin once", func(t *testing.T) {
		ctx := templ.InitializeContext(context.Background())

		html := testRenderIn(t, ctx, templ.Join(Scripts(), Scripts()))

		assert.Equal(t, 1, testCount(html, "<script>"))
		for _, name := range []string{"huiCheckbox", "huiCollapsible", "huiCombobox", "huiDropdown", "huiFileUpload", "huiListbox", "huiModal", "huiRadio", "huiRange", "huiSwitch", "huiTabs", "huiToast", "huiTooltip"} {
			assert.Contains(t, html, name)
		}
		assert.Contains(t, html, "headlessui.pending")
	})

	t.Run("the bundle registers itself", func(t *testing.T) {
		require.NotEmpty(t, scriptsJS)
		assert.Contains(t, scriptsJS, "alpine:init")
		assert.NotContains(t, scriptsJS, "export ")
	})
}
