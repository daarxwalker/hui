package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestLabel(t *testing.T) {
	t.Run("takes for, id and state from the Field", func(t *testing.T) {
		name := gofakeit.Word()
		text := gofakeit.Name()

		html := testRender(t,
			Field(FieldProps{Name: name, Required: true}),
			testWithChildren(Label(LabelProps{Class: "label"}), templ.Raw(text)),
		)

		assert.Contains(t, html, `<label class="label" data-required for="`+name+`" id="`+name+`-label">`+text+`</label>`)
	})

	t.Run("falls back to its own For without a Field", func(t *testing.T) {
		text := gofakeit.Name()

		html := testRender(t, Label(LabelProps{For: "note", Class: "label"}), templ.Raw(text))

		assert.Equal(t, `<label class="label" for="note">`+text+`</label>`, html)
	})

	t.Run("has no for attribute when nothing provides it", func(t *testing.T) {
		html := testRender(t, Label(LabelProps{}))

		assert.NotContains(t, html, "for=")
	})

	t.Run("chains a handler from Attrs with a named event", func(t *testing.T) {
		html := testRender(t, Label(LabelProps{
			Attrs:  templ.Attributes{"@click": "track()"},
			Events: Events{OnClick: "open()"},
		}))

		assert.Contains(t, html, `@click="track(); open()"`)
	})
}
