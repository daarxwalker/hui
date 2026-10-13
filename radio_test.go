package headlessui

import (
	"context"
	"regexp"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRadio(t *testing.T) {
	ctx := context.Background()

	t.Run("renders a button with the radio role", func(t *testing.T) {
		name, value := gofakeit.Word(), gofakeit.Word()

		html := testRender(t, Radio(RadioProps{Name: name, Value: value, Checked: true}), RadioIndicator(RadioIndicatorProps{}))
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `@click="select()" @keydown="keydown($event)" aria-checked="true" data-checked data-hui-radio="`+name+`" name="`+name+`" role="radio" type="button" value="`+value+`"`)
		assert.Contains(t, html, `<input data-hui-checked name="`+name+`" type="hidden" value="`+value+`">`)
		assert.Contains(t, html, `<span aria-hidden="true" x-show="checked" x-transition></span></button>`)
	})

	t.Run("carries its own state without a Field", func(t *testing.T) {
		c := radioAttrs(ctx, RadioProps{Disabled: true, Required: true, Invalid: true, Form: "checkout"})

		assert.Equal(t, true, c.attrs["disabled"])
		assert.Equal(t, "true", c.attrs["aria-invalid"])
		assert.Equal(t, "checkout", c.attrs["form"])
		assert.NotContains(t, c.attrs, "aria-required")
		assert.NotNil(t, c.required)
		for _, state := range []string{"data-disabled", "data-required", "data-invalid"} {
			assert.Equal(t, true, c.attrs[state], state)
		}
		assert.NotContains(t, c.attrs, "data-checked")
	})

	t.Run("is wired into a Field", func(t *testing.T) {
		name, value := gofakeit.Word(), gofakeit.Word()

		html := testRender(t,
			Field(FieldProps{Name: name, Value: value, Readonly: true, Error: gofakeit.Sentence(3)}),
			Radio(RadioProps{Name: "own", Value: "own", Checked: true}),
		)

		assert.Contains(t, html, `id="`+name+`"`)
		assert.Contains(t, html, `name="`+name+`"`)
		assert.Contains(t, html, `value="`+value+`"`)
		assert.Contains(t, html, `aria-describedby="`+name+`-error"`)
		assert.Contains(t, html, `aria-checked="true"`)
		assert.NotContains(t, html, " readonly")
	})

	t.Run("only the checked option of a group is tabbable", func(t *testing.T) {
		chosen := fieldsetIn(FieldsetProps{Name: "plan", Value: "pro"})
		empty := fieldsetIn(FieldsetProps{Name: "plan"})

		assert.Equal(t, "-1", radioAttrs(chosen, RadioProps{Value: "free"}).attrs["tabindex"])
		assert.NotContains(t, radioAttrs(chosen, RadioProps{Value: "pro"}).attrs, "tabindex")
		assert.NotContains(t, radioAttrs(empty, RadioProps{Value: "free"}).attrs, "tabindex")
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		c := radioAttrs(ctx, RadioProps{Base: Base{
			Class:  "size-4",
			Attrs:  templ.Attributes{"role": "checkbox", "hx-post": "/plan", "hx-trigger": "change"},
			Events: Events{OnChange: "save()", OnKeydown: "track($event)"},
		}})

		assert.Equal(t, "radio", c.attrs["role"])
		assert.Equal(t, "size-4", c.attrs["class"])
		assert.Equal(t, "/plan", c.attrs["hx-post"])
		assert.Equal(t, "save()", c.attrs["@change"])
		assert.Equal(t, "keydown($event); track($event)", c.attrs["@keydown"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, radioJS)
		assert.Contains(t, radioJS, "window.Alpine.plugin")
	})
}
