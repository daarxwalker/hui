package headlessui

import (
	"context"
	"regexp"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToast(t *testing.T) {
	t.Run("toaster is a polite live region with a fixed swap target", func(t *testing.T) {
		html := testRender(t, Toaster(ToasterProps{Class: "fixed bottom-4 right-4"}))
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Equal(t, `<ol aria-label="Notifications" aria-live="polite" class="fixed bottom-4 right-4" data-hui-toaster id="hui-toaster"></ol>`, html)
	})

	t.Run("neither the toaster nor a toast ships a script", func(t *testing.T) {
		assert.NotContains(t, testRender(t, Toaster(ToasterProps{})), "<script>")
		assert.NotContains(t, testRender(t, Toast(ToastProps{})), "<script>")
	})

	t.Run("toast wires its title and description", func(t *testing.T) {
		html := testRender(t, Toast(ToastProps{ID: "saved"}),
			ToastTitle(ToastTitleProps{}),
			ToastDescription(ToastDescriptionProps{}),
			ToastAction(ToastActionProps{Attrs: templ.Attributes{"hx-post": "/undo"}}),
			ToastClose(ToastCloseProps{}),
		)

		assert.Contains(t, html, `aria-describedby="saved-description" aria-labelledby="saved-title"`)
		assert.Contains(t, html, `<div id="saved-title"></div><div id="saved-description"></div>`)
		assert.Contains(t, html, `<button data-enabled hx-post="/undo" type="button"></button>`)
		assert.Contains(t, html, `<button @click="dismiss()" data-enabled type="button"></button>`)
	})

	t.Run("duration defaults to five seconds and can be turned off", func(t *testing.T) {
		assert.Equal(t, "huiToast({ duration: 5000 })", toastAttrs(ToastProps{Base: Base{ID: "a"}})["x-data"])
		assert.Equal(t, "huiToast({ duration: 1500 })", toastAttrs(ToastProps{Base: Base{ID: "a"}, Duration: 1500})["x-data"])
		assert.Equal(t, "huiToast({ duration: 0 })", toastAttrs(ToastProps{Base: Base{ID: "a"}, Duration: -1})["x-data"])
	})

	t.Run("urgent toasts are alerts", func(t *testing.T) {
		assert.Equal(t, "status", toastAttrs(ToastProps{Base: Base{ID: "a"}})["role"])
		assert.Equal(t, "alert", toastAttrs(ToastProps{Base: Base{ID: "a"}, Urgent: true})["role"])
	})

	t.Run("hover and focus pause the timer", func(t *testing.T) {
		attrs := toastAttrs(ToastProps{Base: Base{ID: "a", Attrs: templ.Attributes{"@pointerenter": "track()"}}, OnClose: "refresh()"})

		assert.Equal(t, "pause(); track()", attrs["@pointerenter"])
		assert.Equal(t, "resume()", attrs["@pointerleave"])
		assert.Equal(t, "pause()", attrs["@focusin"])
		assert.Equal(t, "resume($event)", attrs["@focusout"])
		assert.Equal(t, "keydown($event)", attrs["@keydown"])
		assert.Equal(t, "open", attrs["x-show"])
		assert.Equal(t, "refresh()", attrs["@hui-toast:closed.self"])
	})

	t.Run("generates an ID when none is given", func(t *testing.T) {
		_, props := toastContext(WithIDScope(context.Background(), "flash"), ToastProps{})

		assert.Equal(t, "flash-1", props.ID)
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, toastJS)
		assert.Contains(t, toastJS, "window.Alpine.plugin")
		assert.NotContains(t, toastJS, "export")
	})

	t.Run("exposes the duration and the pause for a progress bar", func(t *testing.T) {
		attrs := toastAttrs(ToastProps{Base: Base{ID: "a"}, Duration: 3000})

		assert.Equal(t, "--hui-toast-duration: 3000ms", attrs["style"])
		assert.Equal(t, "paused", attrs[":data-paused"])
		assert.NotContains(t, toastAttrs(ToastProps{Base: Base{ID: "a"}, Duration: -1}), "style")
	})
}
