package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	t.Run("takes text and id from the Field", func(t *testing.T) {
		name := gofakeit.Word()
		text := gofakeit.Sentence(5)
		ignored := gofakeit.Sentence(5)

		html := testRender(t,
			Field(FieldProps{Name: name, Description: text, Required: true}),
			Description(DescriptionProps{ID: "ignored", Class: "hint"}),
		)

		assert.Contains(t, html, `<p class="hint" data-required id="`+name+`-description">`+text+`</p>`)
		assert.NotContains(t, html, ignored)
		assert.NotContains(t, html, `id="ignored"`)
	})

	t.Run("is not rendered when the Field has no description", func(t *testing.T) {
		html := testRender(t, Field(FieldProps{Name: gofakeit.Word()}), Description(DescriptionProps{}))

		assert.NotContains(t, html, "<p")
	})

	t.Run("renders children without a Field", func(t *testing.T) {
		text := gofakeit.Sentence(4)

		html := testRender(t, Description(DescriptionProps{ID: "hint", Class: "hint"}), templ.Raw(text))

		assert.Equal(t, `<p class="hint" id="hint">`+text+`</p>`, html)
	})
}
