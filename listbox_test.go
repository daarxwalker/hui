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

func listboxIn(props ListboxProps) context.Context {
	ctx, _ := listboxContext(context.Background(), props)
	return ctx
}

func TestListbox(t *testing.T) {
	t.Run("wires trigger, value, listbox and options", func(t *testing.T) {
		html := testRender(t, Listbox(ListboxProps{ID: "fruit", Name: "fruit", Value: "apple"}),
			testWithChildren(ListboxTrigger(ListboxTriggerProps{}), testWithChildren(ListboxValue(ListboxValueProps{}), templ.Raw("Apple"))),
			testWithChildren(ListboxContent(ListboxContentProps{}),
				testWithChildren(ListboxList(ListboxListProps{}),
					testWithChildren(ListboxItem(ListboxItemProps{Value: "apple"}), ListboxItemIndicator(ListboxItemIndicatorProps{})),
					testWithChildren(ListboxItem(ListboxItemProps{Value: "pear"}), ListboxItemIndicator(ListboxItemIndicatorProps{})),
				),
			),
		)
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `<input data-hui-input="fruit" name="fruit" type="hidden" value="apple">`)
		assert.Contains(t, html, `id="fruit-trigger" role="combobox"`)
		assert.Contains(t, html, `aria-controls="fruit-list"`)
		assert.Contains(t, html, `<span :data-placeholder="!selected().length" data-hui-value="fruit">Apple</span>`)
		assert.Contains(t, html, `<template x-teleport="#hui-portals">`)
		assert.Contains(t, html, `id="fruit-content" style=`)
		assert.Contains(t, html, `aria-labelledby="fruit-trigger" id="fruit-list" role="listbox"`)
		assert.Contains(t, html, `aria-selected="true" data-selected data-value="apple" id="fruit-option-apple" role="option"`)
		assert.Contains(t, html, `aria-selected="false" data-value="pear" id="fruit-option-pear" role="option"`)
		assert.Equal(t, 1, testCount(html, `data-hui-indicator style="display: none"`))
	})

	t.Run("root holds the value for x-model", func(t *testing.T) {
		_, s := listboxContext(context.Background(), ListboxProps{Base: Base{ID: "fruit"}, Name: "fruit", Value: "apple", Form: "order", Placeholder: "Pick"})

		attrs := listboxAttrs(s, ListboxProps{OnChange: "save()"})

		assert.Equal(t, `huiListbox({ open: false, value: "apple", name: "fruit", form: "order", placeholder: "Pick" })`, attrs["x-data"])
		assert.Equal(t, "value", attrs["x-modelable"])
		assert.Equal(t, "fruit", attrs["id"])
		assert.Equal(t, "save()", attrs["@hui-listbox:change.self"])
	})

	t.Run("multiple keeps every value in its own hidden input", func(t *testing.T) {
		_, s := listboxContext(context.Background(), ListboxProps{Base: Base{ID: "tags"}, Name: "tags", Multiple: true, Values: []string{"a", "", "b"}, Disabled: true})

		inputs := listInputAttrs(s)

		assert.Contains(t, s.data(), `multiple: true, value: ["a", "b"]`)
		require.Len(t, inputs, 2)
		assert.Equal(t, "b", inputs[1]["value"])
		assert.Equal(t, true, inputs[1]["disabled"])
	})

	t.Run("required uses a hidden validator that hands focus to the trigger", func(t *testing.T) {
		_, empty := listboxContext(context.Background(), ListboxProps{Name: "fruit", Required: true})
		_, filled := listboxContext(context.Background(), ListboxProps{Name: "fruit", Value: "apple", Required: true})

		assert.Equal(t, "", listRequiredAttrs(empty)["value"])
		assert.Equal(t, "1", listRequiredAttrs(filled)["value"])
		assert.Equal(t, true, listRequiredAttrs(empty)["required"])
		assert.Equal(t, "focusControl()", listRequiredAttrs(empty)["@focus"])
		assert.NotContains(t, listRequiredAttrs(empty), "name")
	})

	t.Run("trigger is a combobox button anchored for the listbox", func(t *testing.T) {
		attrs := listboxTriggerAttrs(listboxIn(ListboxProps{Base: Base{ID: "fruit"}, Invalid: true}), ListboxTriggerProps{Base: Base{Events: Events{OnClick: "track()"}}})

		assert.Equal(t, "button", attrs["type"])
		assert.Equal(t, "listbox", attrs["aria-haspopup"])
		assert.Equal(t, "false", attrs["aria-expanded"])
		assert.Equal(t, "active", attrs[":aria-activedescendant"])
		assert.Equal(t, "true", attrs["aria-invalid"])
		assert.Equal(t, true, attrs["data-placeholder"])
		assert.Equal(t, "anchor-name: --fruit", attrs["style"])
		assert.Equal(t, "listboxToggle(); track()", attrs["@click"])
		assert.Equal(t, "keydown($event)", attrs["@keydown"])
	})

	t.Run("a searchable trigger is a plain button and the search drives the list", func(t *testing.T) {
		ctx := listboxIn(ListboxProps{Base: Base{ID: "country"}, Searchable: true, Required: true})

		trigger := listboxTriggerAttrs(ctx, ListboxTriggerProps{})
		search := listboxSearchAttrs(ctx, ListboxSearchProps{QueryName: "q", Placeholder: "Search"})

		assert.NotContains(t, trigger, "role")
		assert.NotContains(t, trigger, ":aria-activedescendant")
		assert.NotContains(t, trigger, "aria-required")
		assert.Equal(t, "country-list", trigger["aria-controls"])
		assert.Equal(t, "listbox", trigger["aria-haspopup"])
		assert.Equal(t, "country-search", search["id"])
		assert.Equal(t, "combobox", search["role"])
		assert.Equal(t, "list", search["aria-autocomplete"])
		assert.Equal(t, "country-list", search["aria-controls"])
		assert.Equal(t, "active", search[":aria-activedescendant"])
		assert.Equal(t, "country-trigger", search["aria-labelledby"])
		assert.Equal(t, "q", search["name"])
		assert.Equal(t, "Search", search["placeholder"])
		assert.Equal(t, "input()", search["@input"])
		assert.Equal(t, "searchKeydown($event)", search["@keydown"])
		assert.NotContains(t, search, "form")
		_, s := listboxContext(context.Background(), ListboxProps{Searchable: true})
		assert.Contains(t, s.data(), "searchable: true")
	})

	t.Run("the search keeps an aria-label of its own", func(t *testing.T) {
		search := listboxSearchAttrs(listboxIn(ListboxProps{Base: Base{ID: "country"}, Searchable: true}), ListboxSearchProps{Base: Base{Attrs: templ.Attributes{"aria-label": "Search countries", "hx-get": "/countries"}}})

		assert.Equal(t, "Search countries", search["aria-label"])
		assert.NotContains(t, search, "aria-labelledby")
		assert.Equal(t, "/countries", search["hx-get"])
	})

	t.Run("value shows the placeholder until something is picked", func(t *testing.T) {
		html := testRenderIn(t, listboxIn(ListboxProps{Base: Base{ID: "fruit"}, Placeholder: "Pick a fruit"}), ListboxValue(ListboxValueProps{}), templ.Raw("ignored"))

		assert.Equal(t, `<span :data-placeholder="!selected().length" data-hui-value="fruit" data-placeholder>Pick a fruit</span>`, html)
	})

	t.Run("Field is the single source of truth", func(t *testing.T) {
		description := gofakeit.Sentence(3)

		html := testRender(t, Field(FieldProps{Name: "country", Value: "cz", Description: description, Required: true}),
			Label(LabelProps{}),
			testWithChildren(Listbox(ListboxProps{Name: "own", Value: "own"}),
				ListboxTrigger(ListboxTriggerProps{}),
				testWithChildren(ListboxContent(ListboxContentProps{}), ListboxList(ListboxListProps{})),
			),
		)

		assert.Contains(t, html, `<label data-required for="country" id="country-label">`)
		assert.Contains(t, html, `id="country" role="combobox"`)
		assert.Regexp(t, `aria-controls="listbox-[a-z0-9]+-list"`, html)
		assert.Contains(t, html, `aria-describedby="country-description"`)
		assert.Contains(t, html, `aria-required="true"`)
		assert.Contains(t, html, `name="country" type="hidden" value="cz"`)
		assert.Contains(t, html, `aria-labelledby="country-label"`)
		assert.NotContains(t, html, `"own"`)
	})

	t.Run("content is a teleported popup no narrower than the trigger", func(t *testing.T) {
		attrs, teleport := listboxContentAttrs(listboxIn(ListboxProps{Base: Base{ID: "tags"}, Multiple: true}), ListboxContentProps{})

		assert.Equal(t, "#hui-portals", teleport)
		assert.Equal(t, "tags-content", attrs["id"])
		assert.NotContains(t, attrs, "role")
		assert.NotContains(t, attrs, "@mousedown.prevent")
		assert.Equal(t, "open", attrs["x-show"])
		assert.Equal(t, true, attrs["x-transition"])
		assert.Contains(t, attrs["style"], "min-width: anchor-size(width)")
		assert.Contains(t, attrs["style"], "display: none")
	})

	t.Run("list is the listbox the server redraws", func(t *testing.T) {
		attrs := listboxListAttrs(listboxIn(ListboxProps{Base: Base{ID: "tags"}, Multiple: true}), ListboxListProps{Base: Base{Attrs: templ.Attributes{"hx-get": "/tags"}}})

		assert.Equal(t, "tags-list", attrs["id"])
		assert.Equal(t, "listbox", attrs["role"])
		assert.Equal(t, "tags-trigger", attrs["aria-labelledby"])
		assert.Equal(t, "true", attrs["aria-multiselectable"])
		assert.Equal(t, "pick($event)", attrs["@click"])
		assert.Equal(t, "hover($event)", attrs["@pointermove"])
		assert.Equal(t, true, attrs["@mousedown.prevent"])
		assert.Equal(t, "/tags", attrs["hx-get"])
	})

	t.Run("the default transition can be replaced or turned off", func(t *testing.T) {
		ctx := listboxIn(ListboxProps{Base: Base{ID: "tags"}})

		custom, _ := listboxContentAttrs(ctx, ListboxContentProps{Base: Base{Attrs: templ.Attributes{"x-transition:enter": "fade"}}})
		off, _ := listboxContentAttrs(ctx, ListboxContentProps{Base: Base{Attrs: templ.Attributes{"x-transition": false}}})

		assert.NotContains(t, custom, "x-transition")
		assert.Equal(t, "fade", custom["x-transition:enter"])
		assert.Equal(t, false, off["x-transition"])
		assert.NotContains(t, testRenderIn(t, ctx, ListboxContent(ListboxContentProps{Base: Base{Attrs: templ.Attributes{"x-transition": false}}})), "x-transition")
	})

	t.Run("loader and empty state stay out of the accessibility tree of options", func(t *testing.T) {
		html := testRender(t, ListboxLoader(ListboxLoaderProps{Attrs: templ.Attributes{"hx-get": "/tags?page=2", "hx-trigger": "intersect once", "hx-swap": "outerHTML"}}), templ.Raw("Loading"))
		empty := testRender(t, ListboxEmpty(ListboxEmptyProps{}), templ.Raw("Nothing"))

		assert.Equal(t, `<div aria-hidden="true" data-hui-loader hx-get="/tags?page=2" hx-swap="outerHTML" hx-trigger="intersect once">Loading</div>`, html)
		assert.Equal(t, `<div role="presentation">Nothing</div>`, empty)
	})

	t.Run("items rendered in a partial match the full render", func(t *testing.T) {
		props := ListboxProps{Base: Base{ID: "fruit"}, Name: "fruit", Value: "pear"}
		item := testWithChildren(ListboxItem(ListboxItemProps{Value: "pear"}), ListboxItemIndicator(ListboxItemIndicatorProps{}))

		full := testRender(t, Listbox(props), testWithChildren(ListboxContent(ListboxContentProps{}), testWithChildren(ListboxList(ListboxListProps{}), item)))
		partial := testRenderIn(t, WithListbox(context.Background(), props), item)
		bare := testRender(t, item)

		assert.Contains(t, full, partial)
		assert.Equal(t, `<div aria-selected="true" data-selected data-value="pear" id="fruit-option-pear" role="option"><span aria-hidden="true" data-hui-indicator></span></div>`, partial)
		assert.Equal(t, `<div aria-selected="false" data-value="pear" role="option"><span aria-hidden="true" data-hui-indicator></span></div>`, bare)
	})

	t.Run("an item carries its label and disabled state", func(t *testing.T) {
		attrs := listboxItemAttrs(listboxIn(ListboxProps{Base: Base{ID: "fruit"}}), ListboxItemProps{Value: "a b", Label: "A B", Disabled: true})

		assert.Equal(t, "fruit-option-a-b", attrs["id"])
		assert.Equal(t, "A B", attrs["data-label"])
		assert.Equal(t, "true", attrs["aria-disabled"])
		assert.Equal(t, true, attrs["data-disabled"])
	})

	t.Run("group is labelled by its label", func(t *testing.T) {
		html := testRenderIn(t, WithIDScope(context.Background(), "s"), ListboxGroup(ListboxGroupProps{}), ListboxGroupLabel(ListboxGroupLabelProps{}), ListboxSeparator(ListboxSeparatorProps{}))

		assert.Equal(t, `<div aria-labelledby="s-1-label" id="s-1" role="group"><div id="s-1-label"></div><div aria-hidden="true"></div></div>`, html)
	})

	t.Run("parts work outside a listbox", func(t *testing.T) {
		ctx := context.Background()

		trigger := listboxTriggerAttrs(ctx, ListboxTriggerProps{})
		content, teleport := listboxContentAttrs(ctx, ListboxContentProps{})
		item := listboxItemAttrs(ctx, ListboxItemProps{Value: "a"})

		assert.NotContains(t, trigger, "@click")
		assert.Empty(t, teleport)
		assert.NotContains(t, content, "x-show")
		assert.Equal(t, "false", item["aria-selected"])
		assert.NotContains(t, item, "id")
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		ctx := listboxIn(ListboxProps{Base: Base{ID: "fruit"}})

		trigger := listboxTriggerAttrs(ctx, ListboxTriggerProps{Base: Base{
			ID:    "own",
			Class: "w-56",
			Attrs: templ.Attributes{"role": "button", "style": "color: red", "hx-get": "/x"},
		}})
		item := listboxItemAttrs(ctx, ListboxItemProps{Value: "a", Base: Base{Class: "px-2", Attrs: templ.Attributes{"role": "menuitem", "title": "A"}}})

		assert.Equal(t, "fruit-trigger", trigger["id"])
		assert.Equal(t, "combobox", trigger["role"])
		assert.Equal(t, "anchor-name: --fruit; color: red", trigger["style"])
		assert.Equal(t, "/x", trigger["hx-get"])
		assert.Equal(t, "option", item["role"])
		assert.Equal(t, "A", item["title"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, listboxJS)
		assert.Contains(t, listboxJS, "window.Alpine.plugin")
		assert.NotContains(t, listboxJS, "export")
	})

	t.Run("slots stay out of the label", func(t *testing.T) {
		left := gofakeit.Word()
		right := gofakeit.Word()
		label := gofakeit.Word()

		item := testRender(t, ListboxItem(ListboxItemProps{Value: label, LeftSlot: templ.Raw(left), RightSlot: templ.Raw(right)}), templ.Raw(label))
		trigger := testRender(t, ListboxTrigger(ListboxTriggerProps{RightSlot: templ.Raw(right)}), templ.Raw(label))

		assert.Contains(t, item, `<span aria-hidden="true" data-slot="left">`+left+`</span>`+label+`<span aria-hidden="true" data-slot="right">`+right+`</span></div>`)
		assert.Contains(t, trigger, label+`<span aria-hidden="true" data-slot="right">`+right+`</span></button>`)
	})
}
