package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	t.Run("renders a native select with a custom trigger", func(t *testing.T) {
		html := testRender(t, Select(SelectProps{ID: "fruit", Name: "fruit", Value: "pear"}),
			testWithChildren(SelectTrigger(SelectTriggerProps{}), SelectValue(SelectValueProps{})),
			testWithChildren(SelectItem(SelectItemProps{Value: "apple"}), templ.Raw("Apple")),
			testWithChildren(SelectItem(SelectItemProps{Value: "pear"}), templ.Raw("Pear")),
		)

		assert.Equal(t, `<select id="fruit" name="fruit"><button data-enabled><selectedcontent></selectedcontent></button><option value="apple">Apple</option><option selected value="pear">Pear</option></select>`, html)
	})

	t.Run("attributes come from its own props outside a field", func(t *testing.T) {
		_, attrs := selectContext(context.Background(), SelectProps{Base: Base{ID: "fruit"}, Name: "fruit", Form: "order", Autocomplete: "off", Disabled: true, Required: true, Invalid: true})

		assert.Equal(t, "fruit", attrs["id"])
		assert.Equal(t, "fruit", attrs["name"])
		assert.Equal(t, "order", attrs["form"])
		assert.Equal(t, "off", attrs["autocomplete"])
		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, true, attrs["required"])
		assert.Equal(t, "true", attrs["aria-invalid"])
		assert.Equal(t, true, attrs["data-disabled"])
		assert.Equal(t, true, attrs["data-invalid"])
		assert.Equal(t, true, attrs["data-required"])
	})

	t.Run("Field is the single source of truth", func(t *testing.T) {
		name := gofakeit.Word()
		html := testRender(t, Field(FieldProps{Name: name, Value: "pear", Error: gofakeit.Sentence(3), Required: true}),
			testWithChildren(Select(SelectProps{Name: "ignored", Value: "apple"}),
				testWithChildren(SelectItem(SelectItemProps{Value: "apple"}), templ.Raw("Apple")),
				testWithChildren(SelectItem(SelectItemProps{Value: "pear"}), templ.Raw("Pear")),
			),
		)

		assert.Contains(t, html, `<select aria-describedby="`+name+`-error" aria-invalid="true" data-invalid data-required id="`+name+`" name="`+name+`" required>`)
		assert.Contains(t, html, `<option selected value="pear">`)
		assert.NotContains(t, html, "ignored")
	})

	t.Run("the placeholder follows the trigger and is chosen without a value", func(t *testing.T) {
		empty := testRender(t, Select(SelectProps{Placeholder: "Pick a fruit", Required: true}), SelectTrigger(SelectTriggerProps{}))
		filled := testRender(t, Select(SelectProps{Value: "apple", Placeholder: "Pick a fruit"}), SelectTrigger(SelectTriggerProps{}))

		assert.Contains(t, empty, `</button><option data-placeholder selected value="">Pick a fruit</option></select>`)
		assert.Contains(t, filled, `</button><option data-placeholder value="">Pick a fruit</option></select>`)
	})

	t.Run("the trigger reflects a disabled select", func(t *testing.T) {
		ctx, _ := selectContext(context.Background(), SelectProps{Disabled: true})

		attrs, _, _ := selectTriggerAttrs(ctx, SelectTriggerProps{})

		assert.Equal(t, true, attrs["data-disabled"])
		assert.NotContains(t, attrs, "data-enabled")
	})

	t.Run("items rendered in a partial carry no selection", func(t *testing.T) {
		attrs := selectItemAttrs(context.Background(), SelectItemProps{Value: "apple", Disabled: true})

		assert.Equal(t, "apple", attrs["value"])
		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, true, attrs["data-disabled"])
		assert.NotContains(t, attrs, "selected")
	})

	t.Run("group and separator", func(t *testing.T) {
		html := testRender(t, SelectGroup(SelectGroupProps{Label: "Fruit", Disabled: true}), SelectSeparator(SelectSeparatorProps{}))

		assert.Equal(t, `<optgroup data-disabled disabled label="Fruit"><hr></optgroup>`, html)
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		_, attrs := selectContext(context.Background(), SelectProps{Name: "fruit", Base: Base{
			Class:  "w-56",
			Attrs:  templ.Attributes{"name": "other", "class": "border", "hx-post": "/save"},
			Events: Events{OnChange: "save()"},
		}})

		assert.Equal(t, "fruit", attrs["name"])
		assert.Equal(t, "border w-56", attrs["class"])
		assert.Equal(t, "/save", attrs["hx-post"])
		assert.Equal(t, "save()", attrs["@change"])
	})

	t.Run("slots stay out of the label", func(t *testing.T) {
		left := gofakeit.Word()
		label := gofakeit.Word()

		item := testRender(t, SelectItem(SelectItemProps{Value: label, LeftSlot: templ.Raw(left)}), templ.Raw(label))

		assert.Equal(t, `<option value="`+label+`"><span aria-hidden="true" data-slot="left">`+left+`</span>`+label+`</option>`, item)
	})
}
