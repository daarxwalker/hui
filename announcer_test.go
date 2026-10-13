package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncer(t *testing.T) {
	t.Run("is a polite live region with a fixed swap target", func(t *testing.T) {
		html := testRender(t, Announcer(AnnouncerProps{Class: "sr-only"}))

		assert.Equal(t, `<div aria-atomic="true" aria-live="polite" class="sr-only" data-hui-announcer id="hui-announcer" role="status"></div>`, html)
	})

	t.Run("urgent announcer is assertive with its own target", func(t *testing.T) {
		attrs := announcerAttrs(AnnouncerProps{Urgent: true})

		assert.Equal(t, "alert", attrs["role"])
		assert.Equal(t, "assertive", attrs["aria-live"])
		assert.Equal(t, "hui-announcer-urgent", attrs["id"])
	})

	t.Run("ID can be changed", func(t *testing.T) {
		id := gofakeit.Word()

		assert.Equal(t, id, announcerAttrs(AnnouncerProps{Base: Base{ID: id}})["id"])
		assert.Equal(t, id, announcerAttrs(AnnouncerProps{Base: Base{ID: id}, Urgent: true})["id"])
	})

	t.Run("renders its initial message", func(t *testing.T) {
		text := gofakeit.Sentence()

		assert.Contains(t, testRender(t, Announcer(AnnouncerProps{}), templ.Raw(text)), `role="status">`+text+`</div>`)
	})

	t.Run("live region wiring cannot be overridden", func(t *testing.T) {
		attrs := announcerAttrs(AnnouncerProps{Base: Base{Attrs: templ.Attributes{"aria-live": "off", "role": "log", "hx-preserve": true}}})

		assert.Equal(t, "polite", attrs["aria-live"])
		assert.Equal(t, "status", attrs["role"])
		assert.Equal(t, true, attrs["hx-preserve"])
	})
}
