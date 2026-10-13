package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestAlert(t *testing.T) {
	t.Run("is an alert", func(t *testing.T) {
		text := gofakeit.Sentence()

		html := testRender(t, Alert(AlertProps{Class: "rounded border p-4"}), templ.Raw(text))

		assert.Equal(t, `<div class="rounded border p-4" role="alert">`+text+`</div>`, html)
	})

	t.Run("polite alert is a status", func(t *testing.T) {
		assert.Equal(t, "status", alertAttrs(AlertProps{Polite: true})["role"])
	})

	t.Run("role cannot be overridden, the rest passes through", func(t *testing.T) {
		attrs := alertAttrs(AlertProps{Base: Base{ID: "form-error", Attrs: templ.Attributes{"role": "note", "hx-swap-oob": "true"}}})

		assert.Equal(t, "alert", attrs["role"])
		assert.Equal(t, "form-error", attrs["id"])
		assert.Equal(t, "true", attrs["hx-swap-oob"])
	})
}
