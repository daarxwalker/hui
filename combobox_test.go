package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func comboboxIn(props ComboboxProps) context.Context {
	ctx, _ := comboboxContext(context.Background(), props)
	return ctx
}

func TestCombobox(t *testing.T) {
	t.Run("the input lives in the field and the trigger only toggles", func(t *testing.T) {
		ctx, s := comboboxContext(context.Background(), ComboboxProps{Base: Base{ID: "country"}, Name: "country", Value: "cz", Required: true})

		root := comboboxAttrs(s, ComboboxProps{})
		input := comboboxInputAttrs(ctx, ComboboxInputProps{QueryName: "q", Value: "Czechia", Placeholder: "Search"})
		trigger := comboboxTriggerAttrs(ctx, ComboboxTriggerProps{Base: Base{Events: Events{OnClick: "track()"}}})

		assert.Equal(t, `huiCombobox({ open: false, value: "cz", name: "country", searchable: true })`, root["x-data"])
		assert.Equal(t, "anchor-name: --country", root["style"])
		assert.Equal(t, "country-input", input["id"])
		assert.Equal(t, "text", input["type"])
		assert.Equal(t, "combobox", input["role"])
		assert.Equal(t, "list", input["aria-autocomplete"])
		assert.Equal(t, "off", input["autocomplete"])
		assert.Equal(t, "country-list", input["aria-controls"])
		assert.Equal(t, "false", input["aria-expanded"])
		assert.Equal(t, "active", input[":aria-activedescendant"])
		assert.Equal(t, "true", input["aria-required"])
		assert.Equal(t, "q", input["data-hui-query"])
		assert.Equal(t, "Czechia", input["value"])
		assert.Equal(t, "Search", input["placeholder"])
		assert.Equal(t, "input()", input["@input"])
		assert.Equal(t, "keydown($event)", input["@keydown"])
		assert.Equal(t, "inputBlur($event)", input["@blur"])
		assert.Equal(t, "inputClick()", input["@click"])
		assert.NotContains(t, input, "style")
		assert.Equal(t, "country-trigger", trigger["id"])
		assert.Equal(t, "button", trigger["type"])
		assert.Equal(t, "-1", trigger["tabindex"])
		assert.Equal(t, "listbox", trigger["aria-haspopup"])
		assert.Equal(t, "country-list", trigger["aria-controls"])
		assert.Equal(t, true, trigger["@mousedown.prevent"])
		assert.Equal(t, true, trigger["data-enabled"])
		assert.Equal(t, "comboboxToggle(); track()", trigger["@click"])
	})

	t.Run("the query is never submitted with the form", func(t *testing.T) {
		input := comboboxInputAttrs(comboboxIn(ComboboxProps{Name: "country", Form: "f"}), ComboboxInputProps{QueryName: "q"})

		assert.NotContains(t, input, "name")
		assert.NotContains(t, input, "form")
	})

	t.Run("htmx attributes on the input are kept", func(t *testing.T) {
		input := comboboxInputAttrs(comboboxIn(ComboboxProps{Base: Base{ID: "country"}}), ComboboxInputProps{QueryName: "q", Base: Base{
			Attrs:  templ.Attributes{"hx-get": "/countries", "hx-trigger": "input changed delay:200ms", "hx-target": "#country-list", "aria-label": "Search countries"},
			Events: Events{OnInput: "track()"},
		}})

		assert.Equal(t, "/countries", input["hx-get"])
		assert.Equal(t, "#country-list", input["hx-target"])
		assert.Equal(t, "input(); track()", input["@input"])
		assert.Equal(t, "Search countries", input["aria-label"])
		assert.NotContains(t, input, "aria-labelledby")
	})

	t.Run("remove drops one value", func(t *testing.T) {
		attrs := comboboxRemoveAttrs(comboboxIn(ComboboxProps{Multiple: true}), ComboboxRemoveProps{Value: "cz"})

		assert.Equal(t, `remove("cz")`, attrs["@click"])
		assert.Equal(t, "button", attrs["type"])
	})

	t.Run("items rendered by a partial need no combobox around them", func(t *testing.T) {
		html := testRender(t, ComboboxItem(ComboboxItemProps{Value: "sk", Label: "Slovakia"}), ComboboxItemIndicator(ComboboxItemIndicatorProps{}))

		assert.Equal(t, `<div aria-selected="false" data-label="Slovakia" data-value="sk" role="option"><span aria-hidden="true" data-hui-indicator></span></div>`, html)
	})

	t.Run("empty state", func(t *testing.T) {
		html := testRender(t, ComboboxEmpty(ComboboxEmptyProps{Class: "p-2"}), templ.Raw("Nothing"))

		assert.Equal(t, `<div class="p-2" role="presentation">Nothing</div>`, html)
	})

	t.Run("Field gives the input its id and description", func(t *testing.T) {
		html := testRender(t, Field(FieldProps{Name: "country", Value: "cz", Error: "Pick one", Disabled: true}),
			Label(LabelProps{}),
			testWithChildren(Combobox(ComboboxProps{Base: Base{ID: "residence"}}),
				ComboboxInput(ComboboxInputProps{}),
				ComboboxTrigger(ComboboxTriggerProps{}),
				testWithChildren(ComboboxContent(ComboboxContentProps{}), ComboboxList(ComboboxListProps{})),
			),
		)

		assert.Contains(t, html, `<label data-disabled data-invalid for="country" id="country-label">`)
		assert.Contains(t, html, `aria-describedby="country-error" aria-expanded="false" aria-invalid="true" autocomplete="off" data-disabled data-invalid disabled id="country"`)
		assert.Contains(t, html, `aria-labelledby="country-label" data-disabled data-invalid disabled id="residence-trigger"`)
		assert.Contains(t, html, `name="country" type="hidden" value="cz"`)
	})

	t.Run("content is a teleported popup and the list the listbox the server redraws", func(t *testing.T) {
		ctx := comboboxIn(ComboboxProps{Base: Base{ID: "country"}, Multiple: true})

		content, teleport := comboboxContentAttrs(ctx, ComboboxContentProps{})
		list := comboboxListAttrs(ctx, ComboboxListProps{})

		assert.Equal(t, "#hui-portals", teleport)
		assert.Equal(t, "country-content", content["id"])
		assert.NotContains(t, content, "role")
		assert.Equal(t, true, content["x-transition"])
		assert.Contains(t, content["style"], "min-width: anchor-size(width)")
		assert.Equal(t, "country-list", list["id"])
		assert.Equal(t, "listbox", list["role"])
		assert.Equal(t, "country-input", list["aria-labelledby"])
		assert.Equal(t, "true", list["aria-multiselectable"])
		assert.Equal(t, "pick($event)", list["@click"])
		assert.Equal(t, true, list["@mousedown.prevent"])
	})

	t.Run("a loader is the sentinel for the next page", func(t *testing.T) {
		html := testRender(t, ComboboxLoader(ComboboxLoaderProps{Attrs: templ.Attributes{"hx-get": "/countries?page=2", "hx-trigger": "intersect once", "hx-swap": "outerHTML"}}))

		assert.Equal(t, `<div aria-hidden="true" data-hui-loader hx-get="/countries?page=2" hx-swap="outerHTML" hx-trigger="intersect once"></div>`, html)
	})

	t.Run("items rendered in a partial match the full render", func(t *testing.T) {
		props := ComboboxProps{Base: Base{ID: "tags"}, Multiple: true, Values: []string{"go", "htmx"}}

		selected := testRenderIn(t, WithCombobox(context.Background(), props), ComboboxItem(ComboboxItemProps{Value: "htmx"}))
		other := testRenderIn(t, WithCombobox(context.Background(), props), ComboboxItem(ComboboxItemProps{Value: "rust"}))

		assert.Equal(t, `<div aria-selected="true" data-selected data-value="htmx" id="tags-option-htmx" role="option"></div>`, selected)
		assert.Equal(t, `<div aria-selected="false" data-value="rust" id="tags-option-rust" role="option"></div>`, other)
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, comboboxJS)
		assert.Contains(t, comboboxJS, "window.Alpine.plugin")
		assert.NotContains(t, comboboxJS, "export")
	})

	t.Run("slots stay out of the label", func(t *testing.T) {
		left := gofakeit.Word()
		right := gofakeit.Word()
		label := gofakeit.Word()

		item := testRender(t, ComboboxItem(ComboboxItemProps{Value: label, LeftSlot: templ.Raw(left), RightSlot: templ.Raw(right)}), templ.Raw(label))

		assert.Contains(t, item, `<span aria-hidden="true" data-slot="left">`+left+`</span>`+label+`<span aria-hidden="true" data-slot="right">`+right+`</span></div>`)
	})
}
