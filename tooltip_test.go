package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tooltipIn(id string, open bool) context.Context {
	ctx, _ := tooltipContext(context.Background(), TooltipProps{Base: Base{ID: id}, Open: open})
	return ctx
}

func TestTooltip(t *testing.T) {
	t.Run("any element controls it and there is no built-in trigger", func(t *testing.T) {
		html := testRender(t, Tooltip(TooltipProps{ID: "save"}),
			templ.Raw(`<button type="button" aria-describedby="save-content" @click="`+TooltipToggle+`">Save</button>`),
			testWithChildren(TooltipContent(TooltipContentProps{}), TooltipClose(TooltipCloseProps{})),
		)

		assert.Contains(t, html, `@click="tooltipToggle($el, $event)"`)
		assert.NotContains(t, html, `save-trigger`)
		assert.NotContains(t, tooltipJS, `-trigger`)
		assert.Contains(t, html, `<template x-teleport="#hui-portals">`)
		assert.Contains(t, html, `id="save-content" role="tooltip"`)
		assert.Contains(t, html, `@click="tooltipDismiss()"`)
	})

	t.Run("delays reach the frontend", func(t *testing.T) {
		assert.Equal(t, "huiTooltip({ open: false })", tooltipAttrs(TooltipProps{})["x-data"])
		assert.Equal(t, "huiTooltip({ open: true, delay: 100, closeDelay: 50 })",
			tooltipAttrs(TooltipProps{Open: true, Delay: 100, CloseDelay: 50})["x-data"])
	})

	t.Run("controls are unique per component so nested scopes cannot collide", func(t *testing.T) {
		controls := []string{
			TooltipShow, TooltipHide, TooltipToggle, TooltipShowDelayed, TooltipHideDelayed, TooltipFocus, TooltipBlur, TooltipDismiss,
			DropdownShow, DropdownHide, DropdownToggle, DropdownKeydown,
			ModalShow, ModalHide, ModalToggle, CollapsibleShow, CollapsibleHide, CollapsibleToggle,
			ListboxShow, ListboxHide, ListboxToggle, ComboboxShow, ComboboxHide, ComboboxToggle,
		}
		seen := map[string]bool{}
		for _, control := range controls {
			assert.False(t, seen[control], control)
			seen[control] = true
			assert.NotRegexp(t, `^(show|hide|toggle)\(`, control)
		}
		for _, method := range []string{"tooltipShow", "tooltipHide", "tooltipToggle", "tooltipEnter", "tooltipLeave", "tooltipFocus", "tooltipBlur", "tooltipDismiss"} {
			assert.Contains(t, tooltipJS, method)
		}
	})

	t.Run("content sits above the anchor by default", func(t *testing.T) {
		content, _ := tooltipContentAttrs(tooltipIn("save", false), TooltipContentProps{})
		below, _ := tooltipContentAttrs(tooltipIn("save", false), TooltipContentProps{Placement: Placement{Side: "bottom", Align: "start"}})

		assert.Contains(t, content["style"], "position-area: top center")
		assert.Contains(t, below["style"], "position-area: bottom span-right")
		assert.Equal(t, "display: none", attrString(content["style"])[len(attrString(content["style"]))-13:])
	})

	t.Run("hovering the content keeps it open", func(t *testing.T) {
		content, _ := tooltipContentAttrs(tooltipIn("save", true), TooltipContentProps{})

		assert.Equal(t, TooltipShowDelayed, content["@pointerenter"])
		assert.Equal(t, TooltipHideDelayed, content["@pointerleave"])
		assert.Equal(t, true, content["data-open"])
		assert.Equal(t, true, content["x-transition"])
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		content, _ := tooltipContentAttrs(tooltipIn("save", false), TooltipContentProps{Base: Base{
			ID:    gofakeit.Word(),
			Class: "bubble",
			Attrs: templ.Attributes{"@pointerenter": "track()", "data-tone": "dark", "hx-get": "/hint"},
		}})

		assert.Equal(t, "save-content", content["id"])
		assert.Equal(t, "tooltipEnter($event, $el); track()", content["@pointerenter"])
		assert.Equal(t, "bubble", content["class"])
		assert.Equal(t, "dark", content["data-tone"])
		assert.Equal(t, "/hint", content["hx-get"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, tooltipJS)
		assert.Contains(t, tooltipJS, "window.Alpine.plugin")
		assert.NotContains(t, tooltipJS, "export")
	})
}
