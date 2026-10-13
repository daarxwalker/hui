package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func fieldsetIn(props FieldsetProps) context.Context {
	return context.WithValue(context.Background(), fieldsetKey, props.normalize(context.Background()))
}

func TestFieldset(t *testing.T) {
	t.Run("wires legend, description, options and error together", func(t *testing.T) {
		description := gofakeit.Sentence(4)
		message := gofakeit.Sentence(3)

		html := testRender(t,
			Fieldset(FieldsetProps{Name: "plan", Value: "pro", Required: true, Description: description, Error: message}),
			testWithChildren(Legend(LegendProps{}), templ.Raw("Plan")),
			Description(DescriptionProps{}),
			testWithChildren(Field(FieldProps{Value: "free"}), Radio(RadioProps{}), Label(LabelProps{})),
			testWithChildren(Field(FieldProps{Value: "pro"}), Radio(RadioProps{}), Label(LabelProps{})),
			Error(ErrorProps{}),
		)

		assert.Contains(t, html, `<fieldset aria-describedby="plan-description plan-error" data-field="plan" data-invalid data-required id="plan-field">`)
		assert.Contains(t, html, `<legend data-invalid data-required id="plan-legend">Plan</legend>`)
		assert.Contains(t, html, `id="plan-description">`+description+`</p>`)
		assert.Contains(t, html, `id="plan-error" role="alert">`+message+`</p>`)
		assert.Contains(t, html, `aria-checked="false" aria-invalid="true" data-hui-radio="plan" data-invalid data-required id="plan-free" role="radio" tabindex="-1" type="button"`)
		assert.Contains(t, html, `aria-checked="true" aria-invalid="true" data-checked data-hui-radio="plan" data-invalid data-required id="plan-pro" name="plan" role="radio" type="button" value="pro"`)
		assert.Contains(t, html, `<input data-hui-checked name="plan" type="hidden" value="pro">`)
		assert.Contains(t, html, `for="plan-pro"`)
		assert.Equal(t, 2, testCount(html, `name="plan"`))
		assert.Equal(t, 2, testCount(html, `data-hui-require=`))
	})

	t.Run("disabled is native on the fieldset and on every input", func(t *testing.T) {
		attrs := fieldsetAttrs(FieldsetProps{Base: Base{ID: "plan"}, Disabled: true})
		radio := radioAttrs(fieldsetIn(FieldsetProps{Name: "plan", Disabled: true}), RadioProps{Value: "pro"}).attrs

		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, true, attrs["data-disabled"])
		assert.Equal(t, true, radio["disabled"])
		assert.Equal(t, true, radio["data-disabled"])
	})

	t.Run("generates an ID when neither ID nor Name is set", func(t *testing.T) {
		html := testRenderIn(t, WithIDScope(context.Background(), "form"),
			Fieldset(FieldsetProps{}), Legend(LegendProps{}),
		)

		assert.Contains(t, html, `id="form-1-field"`)
		assert.Contains(t, html, `id="form-1-legend"`)
	})

	t.Run("an option takes name, state and checked from the fieldset", func(t *testing.T) {
		ctx := fieldsetIn(FieldsetProps{Name: "tags", Values: []string{"go", "htmx"}, Invalid: true, Form: "post"})

		option := FieldProps{Value: "go", Disabled: true}.normalize(ctx)
		other := FieldProps{Value: "rust"}.normalize(ctx)

		assert.Equal(t, "tags", option.Name)
		assert.Equal(t, "tags-go", option.ID)
		assert.Equal(t, "post", option.Form)
		assert.True(t, option.Invalid)
		assert.True(t, option.Disabled)
		assert.True(t, option.checked)
		assert.False(t, other.checked)
		assert.False(t, other.Disabled)
	})

	t.Run("option IDs are cleaned up and never empty", func(t *testing.T) {
		ctx := WithIDScope(fieldsetIn(FieldsetProps{Name: "size"}), "page")

		assert.Equal(t, "size-x-large", FieldProps{Value: "x large"}.normalize(ctx).ID)
		assert.Equal(t, "page-1", FieldProps{}.normalize(ctx).ID)
	})

	t.Run("the option wrapper is not a field swap target", func(t *testing.T) {
		html := testRender(t,
			Fieldset(FieldsetProps{Name: "plan", Value: "pro"}),
			testWithChildren(Field(FieldProps{Value: "pro"}), Radio(RadioProps{})),
		)

		assert.Contains(t, html, `<div data-option="pro" id="plan-pro-field">`)
		assert.Equal(t, 1, testCount(html, `data-field=`))
	})

	t.Run("checkboxes in a group do not require every option", func(t *testing.T) {
		ctx := fieldsetIn(FieldsetProps{Name: "tags", Values: []string{"go"}, Required: true})

		checkbox := checkboxAttrs(ctx, CheckboxProps{Value: "go"})
		toggle := switchAttrs(ctx, SwitchProps{Value: "go"})
		radio := radioAttrs(ctx, RadioProps{Value: "go"})

		assert.Nil(t, checkbox.required)
		assert.Nil(t, toggle.required)
		assert.NotContains(t, checkbox.attrs, "aria-required")
		assert.NotNil(t, radio.required)
		assert.Equal(t, "true", checkbox.attrs["aria-checked"])
		assert.Equal(t, true, checkbox.attrs["data-required"])
	})

	t.Run("the fieldset decides what is checked", func(t *testing.T) {
		ctx := fieldsetIn(FieldsetProps{Name: "plan", Value: "free"})

		c := radioAttrs(ctx, RadioProps{Value: "pro", Checked: true})

		assert.Equal(t, "plan", c.attrs["data-hui-radio"])
		assert.Equal(t, "false", c.attrs["aria-checked"])
		assert.Nil(t, c.input)
	})

	t.Run("a field with its own name is independent", func(t *testing.T) {
		html := testRender(t,
			Fieldset(FieldsetProps{ID: "address", Disabled: true, Invalid: true, Required: true}),
			testWithChildren(Field(FieldProps{Name: "street"}), Label(LabelProps{}), TextInput(TextInputProps{})),
		)

		assert.Contains(t, html, `<input data-disabled disabled id="street" name="street" type="text">`)
		assert.Contains(t, html, `data-field="street"`)
	})

	t.Run("descriptions and errors belong to the closest owner", func(t *testing.T) {
		group, option := gofakeit.Sentence(3), gofakeit.Sentence(3)

		html := testRender(t,
			Fieldset(FieldsetProps{Name: "plan", Description: group}),
			Description(DescriptionProps{}),
			testWithChildren(Field(FieldProps{Value: "pro", Description: option}), Radio(RadioProps{}), Description(DescriptionProps{})),
		)

		assert.Contains(t, html, `id="plan-description">`+group+`</p>`)
		assert.Contains(t, html, `id="plan-pro-description">`+option+`</p>`)
		assert.Contains(t, html, `aria-describedby="plan-pro-description"`)
	})

	t.Run("inputs directly in a fieldset join the group by value", func(t *testing.T) {
		ctx := fieldsetIn(FieldsetProps{Name: "tags", Values: []string{"go"}, Invalid: true})

		c := checkboxAttrs(ctx, CheckboxProps{Value: "go"})
		own := checkboxAttrs(ctx, CheckboxProps{Name: "newsletter", Value: "yes"})

		assert.Equal(t, "tags", c.input["name"])
		assert.Equal(t, "true", c.attrs["aria-checked"])
		assert.Equal(t, "true", c.attrs["aria-invalid"])
		assert.NotContains(t, c.attrs, "id")
		assert.Contains(t, own.attrs["x-data"], `name: "newsletter"`)
		assert.NotContains(t, own.attrs, "aria-invalid")
	})

	t.Run("text inputs directly in a fieldset only take disabled and form", func(t *testing.T) {
		ctx := fieldsetIn(FieldsetProps{Name: "plan", Disabled: true, Invalid: true, Form: "checkout"})

		attrs := textInputAttrs(ctx, TextInputProps{Name: "note"})

		assert.Equal(t, "note", attrs["name"])
		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, "checkout", attrs["form"])
		assert.NotContains(t, attrs, "aria-invalid")
	})

	t.Run("form reaches the fieldset and the controls inside it", func(t *testing.T) {
		html := testRender(t,
			Fieldset(FieldsetProps{Name: "plan", Form: "checkout"}),
			testWithChildren(Field(FieldProps{Value: "pro"}), Radio(RadioProps{})),
		)

		assert.Equal(t, 2, testCount(html, `form="checkout"`))
	})

	t.Run("legend outside a fieldset renders its children", func(t *testing.T) {
		html := testRender(t, Legend(LegendProps{Class: "sr-only"}), templ.Raw("Plan"))

		assert.Equal(t, `<legend class="sr-only">Plan</legend>`, html)
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		attrs := fieldsetAttrs(FieldsetProps{
			Base: Base{
				ID:     "plan",
				Class:  "grid gap-2",
				Attrs:  templ.Attributes{"class": "border", "aria-describedby": "other", "hx-post": "/plan", "id": "own"},
				Events: Events{OnChange: "dirty = true"},
			},
			Description: gofakeit.Sentence(3),
		})

		assert.Equal(t, "plan-field", attrs["id"])
		assert.Equal(t, "border grid gap-2", attrs["class"])
		assert.Equal(t, "plan-description", attrs["aria-describedby"])
		assert.Equal(t, "/plan", attrs["hx-post"])
		assert.Equal(t, "dirty = true", attrs["@change"])
	})
}
