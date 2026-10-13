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

func TestCheckbox(t *testing.T) {
	ctx := context.Background()
	scripts := regexp.MustCompile(`(?s)<script>.*?</script>`)

	t.Run("renders a button with the checkbox role and a hidden input", func(t *testing.T) {
		name, value := gofakeit.Word(), gofakeit.Word()

		html := scripts.ReplaceAllString(testRender(t, Checkbox(CheckboxProps{Name: name, Value: value, Checked: true})), "")

		assert.Contains(t, html, `<button :aria-checked="ariaChecked()" :data-checked="checked" :data-indeterminate="indeterminate" @click="toggle()" @keydown="keydown($event)" aria-checked="true" data-checked name="`+name+`" role="checkbox" type="button" value="`+value+`"`)
		assert.Contains(t, html, `><input data-hui-checked name="`+name+`" type="hidden" value="`+value+`"></button>`)
	})

	t.Run("an unchecked checkbox submits nothing", func(t *testing.T) {
		c := checkboxAttrs(ctx, CheckboxProps{Name: gofakeit.Word()})

		assert.Nil(t, c.input)
		assert.Equal(t, "false", c.attrs["aria-checked"])
		assert.NotContains(t, c.attrs, "data-checked")
		assert.NotContains(t, c.attrs, "name")
		assert.NotContains(t, c.attrs, "value")
	})

	t.Run("the value defaults to on like a native checkbox", func(t *testing.T) {
		c := checkboxAttrs(ctx, CheckboxProps{Name: "terms", Checked: true})

		assert.Equal(t, "on", c.input["value"])
		assert.Equal(t, `huiCheckbox({ checked: true, name: "terms", value: "on" })`, c.attrs["x-data"])
		assert.Equal(t, "checked", c.attrs["x-modelable"])
	})

	t.Run("indeterminate is mixed and wins over checked", func(t *testing.T) {
		c := checkboxAttrs(ctx, CheckboxProps{Checked: true, Indeterminate: true})

		assert.Equal(t, "mixed", c.attrs["aria-checked"])
		assert.Equal(t, true, c.attrs["data-indeterminate"])
		assert.Contains(t, c.attrs["x-data"], "indeterminate: true")
	})

	t.Run("indicators show the matching state before Alpine starts", func(t *testing.T) {
		checked := checkableContext(ctx, checkableState{checked: true})
		mixed := checkableContext(ctx, checkableState{indeterminate: true})

		assert.NotContains(t, checkboxIndicatorAttrs(checked, CheckboxIndicatorProps{}), "style")
		assert.Equal(t, "display: none", checkboxIndicatorAttrs(mixed, CheckboxIndicatorProps{})["style"])
		assert.NotContains(t, checkboxIndicatorAttrs(mixed, CheckboxIndicatorProps{Indeterminate: true}), "style")
		assert.Equal(t, "checked && !indeterminate", checkboxIndicatorAttrs(checked, CheckboxIndicatorProps{})["x-show"])
		assert.Equal(t, "true", checkboxIndicatorAttrs(checked, CheckboxIndicatorProps{})["aria-hidden"])
	})

	t.Run("indicators transition unless told otherwise", func(t *testing.T) {
		own := checkboxIndicatorAttrs(ctx, CheckboxIndicatorProps{})
		off := checkboxIndicatorAttrs(ctx, CheckboxIndicatorProps{Base: Base{Attrs: templ.Attributes{"x-transition": false}}})
		custom := checkboxIndicatorAttrs(ctx, CheckboxIndicatorProps{Base: Base{Attrs: templ.Attributes{"x-transition.opacity": true}}})

		assert.Equal(t, true, own["x-transition"])
		assert.Equal(t, false, off["x-transition"])
		assert.NotContains(t, custom, "x-transition")
		assert.Equal(t, true, custom["x-transition.opacity"])
	})

	t.Run("the indicator renders a custom icon", func(t *testing.T) {
		html := scripts.ReplaceAllString(testRender(t, Checkbox(CheckboxProps{}),
			testWithChildren(CheckboxIndicator(CheckboxIndicatorProps{}), templ.Raw("<svg></svg>")),
		), "")

		assert.Contains(t, html, `<span aria-hidden="true" style="display: none" x-show="checked &amp;&amp; !indeterminate" x-transition><svg></svg></span></button>`)
	})

	t.Run("carries its own state without a Field", func(t *testing.T) {
		c := checkboxAttrs(ctx, CheckboxProps{
			ID:       "terms",
			Attrs:    templ.Attributes{"aria-describedby": "terms-hint"},
			Required: true,
			Invalid:  true,
			Disabled: true,
		})

		assert.Equal(t, "terms", c.attrs["id"])
		assert.Equal(t, "true", c.attrs["aria-required"])
		assert.Equal(t, true, c.attrs["disabled"])
		assert.Equal(t, "true", c.attrs["aria-invalid"])
		assert.Equal(t, true, c.attrs["data-required"])
		assert.Equal(t, true, c.attrs["data-disabled"])
		assert.Equal(t, true, c.attrs["data-invalid"])
		assert.Equal(t, "terms-hint", c.attrs["aria-describedby"])
	})

	t.Run("required renders a validation input next to the button", func(t *testing.T) {
		empty := checkboxAttrs(ctx, CheckboxProps{ID: "terms", Required: true, Form: "signup"})
		filled := checkboxAttrs(ctx, CheckboxProps{ID: "terms", Required: true, Checked: true})

		require.NotNil(t, empty.required)
		assert.Equal(t, true, empty.required["required"])
		assert.Equal(t, "terms", empty.required["data-hui-require"])
		assert.Equal(t, "signup", empty.required["form"])
		assert.Equal(t, "", empty.required["value"])
		assert.Equal(t, "1", filled.required["value"])
		assert.NotContains(t, empty.required, "name")
		assert.Nil(t, checkboxAttrs(ctx, CheckboxProps{}).required)
	})

	t.Run("Field is the single source of truth", func(t *testing.T) {
		name := gofakeit.Word()
		value := gofakeit.Word()

		html := testRender(t,
			Field(FieldProps{Name: name, Value: value, Description: gofakeit.Sentence(3), Error: gofakeit.Sentence(3), Required: true}),
			Checkbox(CheckboxProps{ID: "own", Name: "own", Value: "own", Checked: true, Disabled: true}),
			Label(LabelProps{}),
			Description(DescriptionProps{}),
		)

		assert.Contains(t, html, `id="`+name+`"`)
		assert.Contains(t, html, `name="`+name+`"`)
		assert.Contains(t, html, `value="`+value+`"`)
		assert.Contains(t, html, `for="`+name+`"`)
		assert.Contains(t, html, `aria-checked="true"`)
		assert.Contains(t, html, `aria-required="true"`)
		assert.Contains(t, html, `aria-invalid="true"`)
		assert.Contains(t, html, `aria-describedby="`+name+`-description `+name+`-error"`)
		assert.NotContains(t, html, `"own"`)
		assert.NotContains(t, html, ` disabled`)
	})

	t.Run("readonly from a Field never reaches the checkbox", func(t *testing.T) {
		html := testRender(t,
			Field(FieldProps{Name: gofakeit.Word(), Readonly: true}),
			Checkbox(CheckboxProps{}),
		)

		assert.NotContains(t, html, ` readonly`)
		assert.Equal(t, 1, testCount(html, "data-readonly"))
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		c := checkboxAttrs(ctx, CheckboxProps{Base: Base{
			Class: "peer size-4",
			Attrs: templ.Attributes{
				"class":   "rounded",
				"x-model": "agreed",
				"role":    "switch",
				"hx-post": "/settings",
			},
			Events: Events{OnChange: "$dispatch('agreed', checked)", OnClick: "track()"},
		}})

		assert.Equal(t, "rounded peer size-4", c.attrs["class"])
		assert.Equal(t, "agreed", c.attrs["x-model"])
		assert.Equal(t, "/settings", c.attrs["hx-post"])
		assert.Equal(t, "toggle(); track()", c.attrs["@click"])
		assert.Equal(t, "$dispatch('agreed', checked)", c.attrs["@change"])
		assert.Equal(t, "checkbox", c.attrs["role"])
		assert.Equal(t, "button", c.attrs["type"])
	})

	t.Run("belongs to a form outside of it", func(t *testing.T) {
		own := checkboxAttrs(ctx, CheckboxProps{Name: "ids", Form: "bulk", Checked: true})
		inField := testRender(t, Field(FieldProps{Name: gofakeit.Word(), Form: "bulk"}), Checkbox(CheckboxProps{Form: "other"}))

		assert.Equal(t, "bulk", own.attrs["form"])
		assert.Equal(t, "bulk", own.input["form"])
		assert.Contains(t, inField, `form="bulk"`)
		assert.NotContains(t, inField, "other")
	})

	t.Run("renders the same in a server partial", func(t *testing.T) {
		html := testRender(t, Checkbox(CheckboxProps{Name: "row", Value: gofakeit.UUID()}))

		assert.Contains(t, html, `x-data="huiCheckbox(`)
		assert.NotContains(t, html, `<script>`)
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, checkboxJS)
		assert.Contains(t, checkboxJS, "window.Alpine.plugin")
		assert.NotContains(t, checkboxJS, "export")
	})
}
