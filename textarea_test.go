package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestTextarea(t *testing.T) {
	ctx := context.Background()

	t.Run("renders the value as text without JavaScript", func(t *testing.T) {
		name := gofakeit.Word()

		html := testRender(t, Textarea(TextareaProps{Name: name, Value: "a < b", Rows: 4}))

		assert.Equal(t, `<textarea name="`+name+`" rows="4">a &lt; b</textarea>`, html)
	})

	t.Run("passes its own attributes through", func(t *testing.T) {
		placeholder := gofakeit.Sentence(3)

		attrs, _ := textareaAttrs(ctx, TextareaProps{Placeholder: placeholder, Autocomplete: "off", Form: "profile"})

		assert.Equal(t, placeholder, attrs["placeholder"])
		assert.Equal(t, "off", attrs["autocomplete"])
		assert.Equal(t, "profile", attrs["form"])
		assert.NotContains(t, attrs, "rows")
	})

	t.Run("carries its own state without a Field", func(t *testing.T) {
		attrs, _ := textareaAttrs(ctx, TextareaProps{Disabled: true, Readonly: true, Required: true, Invalid: true})

		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, true, attrs["readonly"])
		assert.Equal(t, true, attrs["required"])
		assert.Equal(t, "true", attrs["aria-invalid"])
		for _, state := range []string{"data-disabled", "data-readonly", "data-required", "data-invalid"} {
			assert.Equal(t, true, attrs[state], state)
		}
	})

	t.Run("Field is the single source of truth", func(t *testing.T) {
		name := gofakeit.Word()
		bio := gofakeit.Sentence(6)

		html := testRender(t,
			Field(FieldProps{Name: name, Value: bio, Description: gofakeit.Sentence(3)}),
			Textarea(TextareaProps{Name: "own", Value: "own"}),
		)

		assert.Contains(t, html, `id="`+name+`"`)
		assert.Contains(t, html, `name="`+name+`"`)
		assert.Contains(t, html, `aria-describedby="`+name+`-description"`)
		assert.Contains(t, html, `>`+bio+`</textarea>`)
		assert.NotContains(t, html, "own")
	})

	t.Run("Alpine adds the local state when asked for it", func(t *testing.T) {
		empty, _ := textareaAttrs(ctx, TextareaProps{Alpine: true})
		filled, _ := textareaAttrs(ctx, TextareaProps{Alpine: true, Value: gofakeit.Sentence(3)})

		assert.Equal(t, textInputStateEmpty, empty["x-data"])
		assert.Equal(t, true, empty["data-empty"])
		assert.Equal(t, "focused = true", empty["@focus"])
		assert.Equal(t, textInputStateFilled, filled["x-data"])
		assert.NotContains(t, filled, "data-empty")
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		attrs, _ := textareaAttrs(ctx, TextareaProps{
			Alpine: true,
			Base: Base{
				Class:  "field-sizing-content",
				Attrs:  templ.Attributes{"class": "w-full", "hx-post": "/notes", "@input": "dirty = true"},
				Events: Events{OnInput: "count++"},
			},
		})

		assert.Equal(t, "w-full field-sizing-content", attrs["class"])
		assert.Equal(t, "/notes", attrs["hx-post"])
		assert.Equal(t, "dirty = true; empty = $event.target.value === ''; count++", attrs["@input"])
	})
}
