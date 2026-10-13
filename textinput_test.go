package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestTextInput(t *testing.T) {
	t.Run("renders no JavaScript by default", func(t *testing.T) {
		name := gofakeit.Word()

		html := testRender(t, TextInput(TextInputProps{Name: name}))

		assert.Equal(t, `<input name="`+name+`" type="text">`, html)
	})

	t.Run("Alpine adds the local state when asked for it", func(t *testing.T) {
		html := testRender(t, TextInput(TextInputProps{Alpine: true}))

		assert.Contains(t, html, `x-data="{ focused: false, empty: true }"`)
		assert.Contains(t, html, `@focus="focused = true"`)
		assert.Contains(t, html, `@blur="focused = false"`)
		assert.Contains(t, html, `@input="empty = $event.target.value === &#39;&#39;"`)
		assert.Contains(t, html, `:data-focused="focused"`)
		assert.Contains(t, html, `:data-empty="empty"`)
		assert.Contains(t, html, ` data-empty `)
	})

	t.Run("a value switches the Alpine empty state off", func(t *testing.T) {
		email := gofakeit.Email()

		html := testRender(t, TextInput(TextInputProps{Type: "email", Value: email, Alpine: true}))

		assert.Contains(t, html, `value="`+email+`"`)
		assert.Contains(t, html, `x-data="{ focused: false, empty: false }"`)
		assert.NotContains(t, html, ` data-empty`)
	})

	t.Run("named events work without Alpine", func(t *testing.T) {
		html := testRender(t, TextInput(TextInputProps{
			Events: Events{OnBlur: "$store.form.validate('email')"},
		}))

		assert.Contains(t, html, `@blur="$store.form.validate(&#39;email&#39;)"`)
		assert.NotContains(t, html, "x-data")
	})

	t.Run("passes its own attributes through", func(t *testing.T) {
		name := gofakeit.Word()
		placeholder := gofakeit.Sentence(3)

		html := testRender(t, TextInput(TextInputProps{
			Type:         "tel",
			Name:         name,
			Placeholder:  placeholder,
			Autocomplete: "tel",
			Inputmode:    "tel",
		}))

		assert.Contains(t, html, `type="tel"`)
		assert.Contains(t, html, `name="`+name+`"`)
		assert.Contains(t, html, `placeholder="`+placeholder+`"`)
		assert.Contains(t, html, `autocomplete="tel"`)
		assert.Contains(t, html, `inputmode="tel"`)
	})

	t.Run("carries its own state without a Field", func(t *testing.T) {
		html := testRender(t, TextInput(TextInputProps{
			ID:       "note",
			Required: true,
			Readonly: true,
			Invalid:  true,
			Attrs:    templ.Attributes{"aria-describedby": "note-hint"},
		}))

		assert.Contains(t, html, `id="note"`)
		assert.Contains(t, html, ` required`)
		assert.Contains(t, html, ` readonly`)
		assert.Contains(t, html, `aria-invalid="true"`)
		assert.Contains(t, html, `data-invalid`)
		assert.Contains(t, html, `data-required`)
		assert.Contains(t, html, `data-readonly`)
		assert.Contains(t, html, `aria-describedby="note-hint"`)
	})

	t.Run("a disabled input reports it in both ways", func(t *testing.T) {
		html := testRender(t, TextInput(TextInputProps{Disabled: true}))

		assert.Contains(t, html, ` disabled`)
		assert.Contains(t, html, `data-disabled`)
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		html := testRender(t, TextInput(TextInputProps{
			Alpine: true,
			Class:  "input input--lg",
			Attrs: templ.Attributes{
				"x-model":      "email",
				"@keyup.enter": "submit()",
				"type":         "number",
			},
			Events: Events{OnBlur: "$store.form.validate('email')", OnFocus: "log()"},
		}))

		assert.Contains(t, html, `class="input input--lg"`)
		assert.Contains(t, html, `x-model="email"`)
		assert.Contains(t, html, `@keyup.enter="submit()"`)
		assert.Contains(t, html, `@blur="focused = false; $store.form.validate(&#39;email&#39;)"`)
		assert.Contains(t, html, `@focus="focused = true; log()"`)
		assert.Contains(t, html, `type="text"`)
		assert.NotContains(t, html, `type="number"`)
	})

	t.Run("belongs to a form outside of it", func(t *testing.T) {
		attrs := textInputAttrs(context.Background(), TextInputProps{Form: "search"})

		assert.Equal(t, "search", attrs["form"])
	})
}
