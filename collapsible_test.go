package headlessui

import (
	"context"
	"regexp"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollapsible(t *testing.T) {
	t.Run("trigger controls the content", func(t *testing.T) {
		html := testRender(t, Collapsible(CollapsibleProps{ID: "more"}),
			CollapsibleTrigger(CollapsibleTriggerProps{}),
			CollapsibleContent(CollapsibleContentProps{}),
		)
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `x-data="huiCollapsible({ open: false })"`)
		assert.Contains(t, html, `aria-controls="more-content" aria-expanded="false"`)
		assert.Contains(t, html, `@click="collapsibleToggle()"`)
		assert.Contains(t, html, `id="more-content" style="display: none" x-collapse x-show="open"`)
		assert.NotContains(t, html, "anchor-name")
	})

	t.Run("open is rendered from the server", func(t *testing.T) {
		ctx, props := collapsibleContext(context.Background(), CollapsibleProps{Base: Base{ID: "more"}, Open: true, OnOpen: "load()"})

		root := collapsibleAttrs(props)
		content := collapsibleContentAttrs(ctx, CollapsibleContentProps{})

		assert.Equal(t, true, root["data-open"])
		assert.Equal(t, "load()", root["@hui-collapsible:opened.self"])
		assert.Equal(t, true, content["data-open"])
		assert.NotContains(t, content, "style")
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		ctx, _ := collapsibleContext(context.Background(), CollapsibleProps{Base: Base{ID: "more"}})

		attrs := collapsibleTriggerAttrs(ctx, CollapsibleTriggerProps{Base: Base{
			Class:  "underline",
			Attrs:  templ.Attributes{"aria-controls": "other", "hx-get": "/more"},
			Events: Events{OnClick: "track()"},
		}})

		assert.Equal(t, "more-content", attrs["aria-controls"])
		assert.Equal(t, "collapsibleToggle(); track()", attrs["@click"])
		assert.Equal(t, "/more", attrs["hx-get"])
	})

	t.Run("a header row holds the title and a control button", func(t *testing.T) {
		html := testRender(t, Collapsible(CollapsibleProps{ID: "more", Open: true}),
			testWithChildren(CollapsibleHeader(CollapsibleHeaderProps{}),
				testWithChildren(CollapsibleTitle(CollapsibleTitleProps{}), templ.Raw("Details")),
				CollapsibleTrigger(CollapsibleTriggerProps{}),
			),
		)
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `<div :data-open="open" data-open><span :data-open="open" data-open id="more-title">Details</span><button`)
		assert.Contains(t, html, `aria-labelledby="more-title more-trigger"`)
		assert.NotContains(t, html, "data-clickable")
		assert.NotContains(t, html, "collapsibleHeaderClick")
	})

	t.Run("a clickable header toggles from anywhere in the row", func(t *testing.T) {
		ctx, _ := collapsibleContext(context.Background(), CollapsibleProps{Base: Base{ID: "more"}})

		attrs := collapsibleHeaderAttrs(ctx, CollapsibleHeaderProps{Clickable: true, Base: Base{Events: Events{OnClick: "track()"}}})
		alone := collapsibleHeaderAttrs(context.Background(), CollapsibleHeaderProps{Clickable: true})

		assert.Equal(t, "collapsibleHeaderClick($event); track()", attrs["@click"])
		assert.Equal(t, true, attrs["data-clickable"])
		assert.NotContains(t, alone, "@click")
	})

	t.Run("a trigger outside a header keeps its own name", func(t *testing.T) {
		ctx, _ := collapsibleContext(context.Background(), CollapsibleProps{Base: Base{ID: "more"}})

		assert.NotContains(t, collapsibleTriggerAttrs(ctx, CollapsibleTriggerProps{}), "aria-labelledby")
	})

	t.Run("an aria-label from Attrs wins over the title", func(t *testing.T) {
		ctx, _ := collapsibleContext(context.Background(), CollapsibleProps{Base: Base{ID: "more"}})

		attrs := collapsibleTriggerAttrs(collapsibleHeaderContext(ctx), CollapsibleTriggerProps{Attrs: templ.Attributes{"aria-label": "Show details"}})

		assert.Equal(t, "Show details", attrs["aria-label"])
		assert.NotContains(t, attrs, "aria-labelledby")
	})

	t.Run("content collapses with x-collapse unless told otherwise", func(t *testing.T) {
		ctx, _ := collapsibleContext(context.Background(), CollapsibleProps{Base: Base{ID: "more"}})

		own := collapsibleContentAttrs(ctx, CollapsibleContentProps{})
		fade := collapsibleContentAttrs(ctx, CollapsibleContentProps{Base: Base{Attrs: templ.Attributes{"x-transition.opacity": true}}})
		off := collapsibleContentAttrs(ctx, CollapsibleContentProps{Base: Base{Attrs: templ.Attributes{"x-collapse": false}}})

		assert.Equal(t, true, own["x-collapse"])
		assert.NotContains(t, own, "x-transition")
		assert.NotContains(t, fade, "x-collapse")
		assert.Equal(t, true, fade["x-transition.opacity"])
		assert.Equal(t, false, off["x-collapse"])
		assert.NotContains(t, off, "x-transition")
	})

	t.Run("controls reach the plugin", func(t *testing.T) {
		for _, method := range []string{"collapsibleShow", "collapsibleHide", "collapsibleToggle", "collapsibleHeaderClick"} {
			assert.Contains(t, collapsibleJS, method)
		}
		assert.Contains(t, collapsibleJS, `"collapse"`)
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, collapsibleJS)
		assert.Contains(t, collapsibleJS, "window.Alpine.plugin")
	})
}
