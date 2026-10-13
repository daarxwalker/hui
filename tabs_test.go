package headlessui

import (
	"context"
	"regexp"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tabsIn(id, value, orientation string) context.Context {
	ctx, _ := tabsContext(context.Background(), TabsProps{Base: Base{ID: id}, Value: value, Orientation: orientation})
	return ctx
}

func TestTabs(t *testing.T) {
	t.Run("wires triggers and panels by value", func(t *testing.T) {
		html := testRender(t, Tabs(TabsProps{ID: "settings", Value: "profile"}),
			testWithChildren(TabsList(TabsListProps{}),
				TabsTrigger(TabsTriggerProps{Value: "profile"}),
				TabsTrigger(TabsTriggerProps{Value: "billing"}),
			),
			TabsContent(TabsContentProps{Value: "profile"}),
			TabsContent(TabsContentProps{Value: "billing"}),
		)
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `role="tablist"`)
		assert.Contains(t, html, `id="settings-trigger-profile"`)
		assert.Contains(t, html, `aria-controls="settings-content-billing"`)
		assert.Contains(t, html, `aria-labelledby="settings-trigger-billing"`)
		assert.Equal(t, 1, testCount(html, `aria-selected="true"`))
		assert.Equal(t, 1, testCount(html, `style="display: none"`))
	})

	t.Run("root holds the value", func(t *testing.T) {
		attrs := tabsAttrs(TabsProps{Value: "it's", Activation: "manual", OnChange: "save($event.detail.value)"})

		assert.Equal(t, `huiTabs({ value: "it's", manual: true })`, attrs["x-data"])
		assert.Equal(t, "value", attrs["x-modelable"])
		assert.Equal(t, "horizontal", attrs["data-orientation"])
		assert.Equal(t, "save($event.detail.value)", attrs["@hui-tabs:change.self"])
	})

	t.Run("the selected trigger is the only tab stop", func(t *testing.T) {
		ctx := tabsIn("settings", "profile", "")

		selected := tabsTriggerAttrs(ctx, TabsTriggerProps{Value: "profile"})
		other := tabsTriggerAttrs(ctx, TabsTriggerProps{Value: "billing"})

		assert.Equal(t, "0", selected["tabindex"])
		assert.Equal(t, true, selected["data-selected"])
		assert.Equal(t, "-1", other["tabindex"])
		assert.Equal(t, "false", other["aria-selected"])
		assert.Equal(t, `value === "billing" ? 0 : -1`, other[":tabindex"])
		assert.Equal(t, `select("billing")`, other["@click"])
		assert.NotContains(t, other, "data-selected")
	})

	t.Run("list handles the arrow keys in its orientation", func(t *testing.T) {
		attrs := tabsListAttrs(tabsIn("settings", "", "vertical"), TabsListProps{Base: Base{Attrs: templ.Attributes{"aria-label": "Settings"}}})

		assert.Equal(t, "vertical", attrs["aria-orientation"])
		assert.Equal(t, "keydown($event)", attrs["@keydown"])
		assert.Equal(t, "Settings", attrs["aria-label"])
	})

	t.Run("a panel is focusable and shown by value", func(t *testing.T) {
		attrs := tabsContentAttrs(tabsIn("settings", "profile", ""), TabsContentProps{Value: "profile"})

		assert.Equal(t, "tabpanel", attrs["role"])
		assert.Equal(t, "0", attrs["tabindex"])
		assert.Equal(t, `value === "profile"`, attrs["x-show"])
		assert.Equal(t, true, attrs["data-selected"])
		assert.NotContains(t, attrs, "style")
	})

	t.Run("a disabled trigger uses the native attribute", func(t *testing.T) {
		attrs := tabsTriggerAttrs(tabsIn("settings", "", ""), TabsTriggerProps{Value: "danger", Disabled: true})

		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, true, attrs["data-disabled"])
		assert.NotContains(t, attrs, "data-enabled")
	})

	t.Run("parts work outside tabs", func(t *testing.T) {
		trigger := tabsTriggerAttrs(context.Background(), TabsTriggerProps{Value: "a"})
		content := tabsContentAttrs(context.Background(), TabsContentProps{Value: "a"})

		assert.NotContains(t, trigger, "@click")
		assert.NotContains(t, content, "x-show")
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		attrs := tabsTriggerAttrs(tabsIn("settings", "profile", ""), TabsTriggerProps{Value: "billing", Base: Base{
			ID:     "own",
			Class:  "px-3",
			Attrs:  templ.Attributes{"role": "button", "hx-get": "/billing", "hx-trigger": "click once"},
			Events: Events{OnClick: "track()"},
		}})

		assert.Equal(t, "settings-trigger-billing", attrs["id"])
		assert.Equal(t, "tab", attrs["role"])
		assert.Equal(t, "/billing", attrs["hx-get"])
		assert.Equal(t, `select("billing"); track()`, attrs["@click"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, tabsJS)
		assert.Contains(t, tabsJS, "window.Alpine.plugin")
		assert.NotContains(t, tabsJS, "export")
	})
}
