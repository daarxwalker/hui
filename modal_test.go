package headlessui

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func modalHTML(t *testing.T, props ModalProps, children ...templ.Component) string {
	t.Helper()
	return regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(testRender(t, Modal(props), children...), "")
}

func modalIn(id string, open bool) context.Context {
	ctx, _ := modalContext(context.Background(), ModalProps{ID: id, Open: open})
	return ctx
}

func TestModal(t *testing.T) {
	t.Run("renders a div and no script", func(t *testing.T) {
		ctx := templ.InitializeContext(context.Background())

		html := testRenderIn(t, ctx, templ.Join(Modal(ModalProps{ID: "a"}), Modal(ModalProps{ID: "b"})))

		assert.Equal(t, 2, strings.Count(html, `x-data="huiModal(`))
		assert.NotContains(t, html, "<script>")
		assert.NotContains(t, html, "<dialog")
	})

	t.Run("parts are wired to the modal ID", func(t *testing.T) {
		html := modalHTML(t, ModalProps{ID: "confirm"},
			ModalTrigger(ModalTriggerProps{}),
			ModalOverlay(ModalOverlayProps{}),
			testWithChildren(ModalContent(ModalContentProps{}),
				ModalTitle(ModalTitleProps{}),
				ModalDescription(ModalDescriptionProps{}),
				ModalClose(ModalCloseProps{}),
			),
		)

		assert.Contains(t, html, `aria-controls="confirm-content"`)
		assert.Contains(t, html, `id="confirm-overlay"`)
		assert.Contains(t, html, `id="confirm-content"`)
		assert.Contains(t, html, `id="confirm-title"`)
		assert.Contains(t, html, `id="confirm-description"`)
		assert.Contains(t, html, `@click="modalHide()"`)
	})

	t.Run("portals container has a fixed ID", func(t *testing.T) {
		attrs := portalsAttrs(PortalsProps{Base: Base{
			ID:    gofakeit.Word(),
			Class: "isolate",
			Attrs: templ.Attributes{"data-layer": "modals"},
		}})

		assert.Equal(t, "hui-portals", attrs["id"])
		assert.Equal(t, "isolate", attrs["class"])
		assert.Equal(t, "modals", attrs["data-layer"])
	})

	t.Run("overlay and content are teleported into the portals", func(t *testing.T) {
		ctx := modalIn("confirm", false)

		_, overlay := modalOverlayAttrs(ctx, ModalOverlayProps{})
		_, content := modalContentAttrs(ctx, ModalContentProps{})
		html := modalHTML(t, ModalProps{ID: "confirm"},
			ModalTrigger(ModalTriggerProps{}),
			ModalOverlay(ModalOverlayProps{}),
			ModalContent(ModalContentProps{}),
		)

		assert.Equal(t, "#hui-portals", overlay)
		assert.Equal(t, "#hui-portals", content)
		assert.Equal(t, 2, testCount(html, `<template x-teleport="#hui-portals">`))
		assert.Equal(t, 2, testCount(html, `</div></template>`))
		assert.Equal(t, 0, testCount(html, `<template x-teleport="#hui-portals"><button`))
	})

	t.Run("content inside an overlay travels with it", func(t *testing.T) {
		_, teleport := modalContentAttrs(modalOverlayContext(modalIn("confirm", false)), ModalContentProps{})
		html := modalHTML(t, ModalProps{ID: "confirm"},
			testWithChildren(ModalOverlay(ModalOverlayProps{}), ModalContent(ModalContentProps{})),
		)

		assert.Empty(t, teleport)
		assert.Equal(t, 1, testCount(html, "<template "))
		assert.Contains(t, html, `id="confirm-content"`)
	})

	t.Run("is closed by default", func(t *testing.T) {
		attrs := modalAttrs(ModalProps{Base: Base{ID: "confirm"}})

		assert.Equal(t, "huiModal({ open: false })", attrs["x-data"])
		assert.NotContains(t, attrs, "data-open")
	})

	t.Run("state is controllable through x-model", func(t *testing.T) {
		attrs := modalAttrs(ModalProps{Base: Base{ID: "confirm", Attrs: templ.Attributes{"x-model": "confirming"}}})

		assert.Equal(t, "open", attrs["x-modelable"])
		assert.Equal(t, "confirming", attrs["x-model"])
		assert.Equal(t, "open", attrs[":data-open"])
	})

	t.Run("an open modal is visible before Alpine starts", func(t *testing.T) {
		ctx := modalIn("confirm", true)

		root := modalAttrs(ModalProps{Base: Base{ID: "confirm"}, Open: true})
		content, _ := modalContentAttrs(ctx, ModalContentProps{})
		overlay, _ := modalOverlayAttrs(ctx, ModalOverlayProps{})
		trigger := modalTriggerAttrs(ctx, ModalTriggerProps{})

		assert.Equal(t, "huiModal({ open: true })", root["x-data"])
		assert.Equal(t, true, root["data-open"])
		for _, attrs := range []templ.Attributes{content, overlay} {
			assert.Equal(t, true, attrs["data-open"])
			assert.NotContains(t, attrs, "style")
		}
		assert.Equal(t, "true", trigger["aria-expanded"])
		assert.Equal(t, true, trigger["data-open"])
	})

	t.Run("a closed modal hides its overlay and content without CSS", func(t *testing.T) {
		ctx := modalIn("confirm", false)

		content, _ := modalContentAttrs(ctx, ModalContentProps{})
		overlay, _ := modalOverlayAttrs(ctx, ModalOverlayProps{})

		for _, attrs := range []templ.Attributes{content, overlay} {
			assert.Equal(t, "display: none", attrs["style"])
			assert.Equal(t, "open", attrs["x-show"])
			assert.Equal(t, "open", attrs[":data-open"])
			assert.NotContains(t, attrs, "data-open")
		}
	})

	t.Run("content is an accessible modal dialog", func(t *testing.T) {
		attrs, _ := modalContentAttrs(modalIn("confirm", false), ModalContentProps{})

		assert.Equal(t, "dialog", attrs["role"])
		assert.Equal(t, "true", attrs["aria-modal"])
		assert.Equal(t, "-1", attrs["tabindex"])
		assert.Equal(t, "confirm-title", attrs["aria-labelledby"])
		assert.Equal(t, "confirm-description", attrs["aria-describedby"])
		assert.NotContains(t, attrs, "@click.outside")
	})

	t.Run("an alert modal is an alertdialog that ignores clicks outside", func(t *testing.T) {
		ctx, props := modalContext(context.Background(), ModalProps{Base: Base{ID: "delete"}, Alert: true})

		attrs, _ := modalContentAttrs(ctx, ModalContentProps{})

		assert.Equal(t, "huiModal({ open: false, alert: true })", modalAttrs(props)["x-data"])
		assert.Equal(t, "alertdialog", attrs["role"])
		assert.Equal(t, "true", attrs["aria-modal"])
	})

	t.Run("the a11y wiring cannot be overwritten", func(t *testing.T) {
		attrs, _ := modalContentAttrs(modalIn("confirm", false), ModalContentProps{Base: Base{
			ID:    gofakeit.Word(),
			Attrs: templ.Attributes{"role": "region", "aria-modal": "false", "@click.outside": "track()"},
		}})

		assert.Equal(t, "confirm-content", attrs["id"])
		assert.Equal(t, "dialog", attrs["role"])
		assert.Equal(t, "true", attrs["aria-modal"])
		assert.Equal(t, "track()", attrs["@click.outside"])
	})

	t.Run("a style from Attrs is appended to the hidden state", func(t *testing.T) {
		attrs, _ := modalOverlayAttrs(modalIn("confirm", false), ModalOverlayProps{Base: Base{
			Attrs: templ.Attributes{"style": "--blur: 4px"},
		}})

		assert.Equal(t, "display: none; --blur: 4px", attrs["style"])
		assert.Equal(t, "true", attrs["aria-hidden"])
	})

	t.Run("trigger opens the modal and describes what it opens", func(t *testing.T) {
		attrs := modalTriggerAttrs(modalIn("confirm", false), ModalTriggerProps{Base: Base{
			Events: Events{OnClick: "count++"},
		}})

		assert.Equal(t, "button", attrs["type"])
		assert.Equal(t, "dialog", attrs["aria-haspopup"])
		assert.Equal(t, "false", attrs["aria-expanded"])
		assert.Equal(t, "open", attrs[":aria-expanded"])
		assert.Equal(t, "modalShow(); count++", attrs["@click"])
	})

	t.Run("a disabled trigger uses the native attribute", func(t *testing.T) {
		attrs := modalTriggerAttrs(modalIn("confirm", false), ModalTriggerProps{Disabled: true})

		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, true, attrs["data-disabled"])
		assert.NotContains(t, attrs, "data-enabled")
	})

	t.Run("buttons carry data-enabled while they can be used", func(t *testing.T) {
		ctx := modalIn("confirm", false)

		assert.Equal(t, true, modalTriggerAttrs(ctx, ModalTriggerProps{})["data-enabled"])
		assert.Equal(t, true, modalCloseAttrs(ctx, ModalCloseProps{})["data-enabled"])
	})

	t.Run("close runs hide before a named click event", func(t *testing.T) {
		attrs := modalCloseAttrs(modalIn("confirm", false), ModalCloseProps{Base: Base{
			Events: Events{OnClick: "track()"},
		}})

		assert.Equal(t, "button", attrs["type"])
		assert.Equal(t, "modalHide(); track()", attrs["@click"])
	})

	t.Run("parts work on their own outside a modal", func(t *testing.T) {
		ctx := context.Background()
		id := gofakeit.Word()

		trigger := modalTriggerAttrs(ctx, ModalTriggerProps{})
		content, teleport := modalContentAttrs(ctx, ModalContentProps{Base: Base{ID: id}})
		_, overlayTeleport := modalOverlayAttrs(ctx, ModalOverlayProps{})
		title := modalTitleAttrs(ctx, ModalTitleProps{Base: Base{ID: id}})
		closer := modalCloseAttrs(ctx, ModalCloseProps{})

		assert.NotContains(t, trigger, "@click")
		assert.NotContains(t, trigger, "aria-controls")
		assert.Equal(t, id, content["id"])
		assert.NotContains(t, content, "x-show")
		assert.NotContains(t, content, "aria-labelledby")
		assert.Empty(t, teleport)
		assert.Empty(t, overlayTeleport)
		assert.Equal(t, id, title["id"])
		assert.NotContains(t, closer, "@click")
	})

	t.Run("generates an ID when none is given", func(t *testing.T) {
		_, props := modalContext(WithIDScope(context.Background(), "page"), ModalProps{})

		assert.Equal(t, "page-1", modalAttrs(props)["id"])
	})

	t.Run("persist is keyed by the modal ID", func(t *testing.T) {
		id := gofakeit.Word()

		attrs := modalAttrs(ModalProps{Base: Base{ID: id}, Persist: true})

		assert.Equal(t, "huiModal({ open: false, persist: '"+id+"' })", attrs["x-data"])
	})

	t.Run("static and allow scroll reach the frontend", func(t *testing.T) {
		attrs := modalAttrs(ModalProps{Base: Base{ID: "confirm"}, Static: true, AllowScroll: true})

		assert.Equal(t, "huiModal({ open: false, static: true, allowScroll: true })", attrs["x-data"])
	})

	t.Run("open and close hooks listen for the dispatched events", func(t *testing.T) {
		attrs := modalAttrs(ModalProps{Base: Base{ID: "confirm"}, OnOpen: "log('open')", OnClose: "refresh()"})

		assert.Equal(t, "log('open')", attrs["@hui-modal:opened.self"])
		assert.Equal(t, "refresh()", attrs["@hui-modal:closed.self"])
	})

	t.Run("hooks are chained with handlers from Attrs", func(t *testing.T) {
		attrs := modalAttrs(ModalProps{
			Base:    Base{ID: "confirm", Attrs: templ.Attributes{"@hui-modal:closed.self": "track()"}},
			OnClose: "refresh()",
		})

		assert.Equal(t, "track(); refresh()", attrs["@hui-modal:closed.self"])
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		attrs, _ := modalContentAttrs(modalIn("confirm", false), ModalContentProps{Base: Base{
			Class: "rounded-xl",
			Attrs: templ.Attributes{
				"class":        "fixed inset-0",
				"x-transition": true,
				"hx-get":       "/confirm",
			},
		}})

		assert.Equal(t, "fixed inset-0 rounded-xl", attrs["class"])
		assert.Equal(t, true, attrs["x-transition"])
		assert.Equal(t, "/confirm", attrs["hx-get"])
	})

	t.Run("overlay fades and content transitions by default", func(t *testing.T) {
		ctx := modalOverlayContext(modalIn("confirm", false))

		overlay, _ := modalOverlayAttrs(ctx, ModalOverlayProps{})
		content, _ := modalContentAttrs(ctx, ModalContentProps{})

		assert.Equal(t, true, overlay["x-transition.opacity"])
		assert.NotContains(t, overlay, "x-transition")
		assert.Equal(t, true, content["x-transition"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, modalJS)
		assert.Contains(t, modalJS, "window.Alpine.plugin")
		assert.Contains(t, modalJS, "alpine:init")
		assert.NotContains(t, modalJS, "export")
		assert.NotContains(t, modalJS, "showModal")
	})
}
