package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("takes text and id from the Field and announces itself", func(t *testing.T) {
		name := gofakeit.Word()
		text := gofakeit.Sentence(5)

		html := testRender(t,
			Field(FieldProps{Name: name, Error: text}),
			Error(ErrorProps{Class: "error"}),
		)

		assert.Contains(t, html, `id="`+name+`-error"`)
		assert.Contains(t, html, `role="alert"`)
		assert.Contains(t, html, `data-invalid`)
		assert.Contains(t, html, `>`+text+`</p>`)
	})

	t.Run("is not rendered when the Field has no error", func(t *testing.T) {
		html := testRender(t, Field(FieldProps{Name: gofakeit.Word()}), Error(ErrorProps{}))

		assert.NotContains(t, html, "<p")
	})

	t.Run("renders children without a Field", func(t *testing.T) {
		text := gofakeit.Sentence(4)

		html := testRender(t, Error(ErrorProps{Class: "error"}), templ.Raw(text))

		assert.Equal(t, `<p class="error" role="alert">`+text+`</p>`, html)
	})
}
