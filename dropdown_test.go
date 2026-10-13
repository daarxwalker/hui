package headlessui

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dropdownIn(id string, open bool) context.Context {
	ctx, _ := dropdownContext(context.Background(), DropdownProps{Base: Base{ID: id}, Open: open})
	return ctx
}

func TestDropdown(t *testing.T) {
	t.Run("wires trigger and content together", func(t *testing.T) {
		html := testRender(t,
			Dropdown(DropdownProps{ID: "account"}),
			DropdownTrigger(DropdownTriggerProps{}),
			testWithChildren(DropdownContent(DropdownContentProps{}),
				DropdownItem(DropdownItemProps{}),
			),
		)
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `id="account"`)
		assert.Contains(t, html, `id="account-trigger"`)
		assert.Contains(t, html, `aria-controls="account-content"`)
		assert.Contains(t, html, `<template x-teleport="#hui-portals"><div :data-open="open"`)
		assert.Contains(t, html, `id="account-content"`)
		assert.NotContains(t, html, "popover")
		assert.Contains(t, html, `aria-labelledby="account-trigger"`)
		assert.Contains(t, html, `role="menuitem"`)
	})

	t.Run("ships no script, the script is opt-in and rendered once", func(t *testing.T) {
		ctx := templ.InitializeContext(context.Background())

		html := testRenderIn(t, ctx, templ.Join(Dropdown(DropdownProps{}), DropdownScript(), DropdownScript()))

		assert.Equal(t, 1, testCount(html, "<script>"))
		assert.Contains(t, html, `"huiDropdown"`)
		assert.NotContains(t, testRender(t, Dropdown(DropdownProps{})), "<script>")
	})

	t.Run("generates an ID when none is given", func(t *testing.T) {
		ctx := WithIDScope(context.Background(), "page")

		html := testRenderIn(t, ctx, testWithChildren(Dropdown(DropdownProps{}), DropdownTrigger(DropdownTriggerProps{})))

		assert.Contains(t, html, `id="page-1"`)
		assert.Contains(t, html, `id="page-1-trigger"`)
	})

	t.Run("root holds the state", func(t *testing.T) {
		attrs := dropdownAttrs(DropdownProps{Base: Base{ID: "account"}})

		assert.Equal(t, "huiDropdown({ open: false })", attrs["x-data"])
		assert.Equal(t, "open", attrs["x-modelable"])
		assert.Equal(t, "open", attrs[":data-open"])
		assert.NotContains(t, attrs, "data-open")
		assert.NotContains(t, attrs, "@keydown")
	})

	t.Run("an open dropdown is visible before Alpine starts", func(t *testing.T) {
		ctx := dropdownIn("account", true)

		root := dropdownAttrs(DropdownProps{Base: Base{ID: "account"}, Open: true})
		trigger := dropdownTriggerAttrs(ctx, DropdownTriggerProps{})
		content, _ := dropdownContentAttrs(ctx, DropdownContentProps{})

		assert.Equal(t, "huiDropdown({ open: true })", root["x-data"])
		assert.Equal(t, true, root["data-open"])
		assert.Equal(t, "true", trigger["aria-expanded"])
		assert.Equal(t, true, trigger["data-open"])
		assert.Equal(t, true, content["data-open"])
		assert.NotContains(t, attrString(content["style"]), "display: none")
	})

	t.Run("open and close hooks listen for the dispatched events", func(t *testing.T) {
		attrs := dropdownAttrs(DropdownProps{
			OnOpen:  "log('open')",
			OnClose: "refresh()",
			Attrs:   templ.Attributes{"@hui-dropdown:closed.self": "track()"},
		})

		assert.Equal(t, "log('open')", attrs["@hui-dropdown:opened.self"])
		assert.Equal(t, "track(); refresh()", attrs["@hui-dropdown:closed.self"])
	})

	t.Run("trigger is a menu button", func(t *testing.T) {
		attrs := dropdownTriggerAttrs(dropdownIn("account", false), DropdownTriggerProps{})

		assert.Equal(t, "button", attrs["type"])
		assert.Equal(t, "menu", attrs["aria-haspopup"])
		assert.Equal(t, "account-content", attrs["aria-controls"])
		assert.Equal(t, "false", attrs["aria-expanded"])
		assert.Equal(t, "open", attrs[":aria-expanded"])
		assert.Equal(t, "open", attrs[":data-open"])
		assert.Equal(t, "anchor-name: --account", attrs["style"])
		assert.Equal(t, DropdownToggle, attrs["@click"])
		assert.Equal(t, DropdownKeydown, attrs["@keydown"])
	})

	t.Run("dropdown owns the part IDs", func(t *testing.T) {
		ctx := dropdownIn("account", false)

		trigger := dropdownTriggerAttrs(ctx, DropdownTriggerProps{Base: Base{ID: gofakeit.Word()}})
		content, _ := dropdownContentAttrs(ctx, DropdownContentProps{Base: Base{ID: gofakeit.Word()}})

		assert.Equal(t, "account-trigger", trigger["id"])
		assert.Equal(t, "account-content", content["id"])
	})

	t.Run("disabled trigger uses the native attribute", func(t *testing.T) {
		attrs := dropdownTriggerAttrs(dropdownIn("account", false), DropdownTriggerProps{Disabled: true})

		assert.Equal(t, true, attrs["disabled"])
		assert.Equal(t, true, attrs["data-disabled"])
		assert.NotContains(t, attrs, "data-enabled")
		assert.Equal(t, true, dropdownTriggerAttrs(dropdownIn("account", false), DropdownTriggerProps{})["data-enabled"])
	})

	t.Run("content is a teleported menu", func(t *testing.T) {
		attrs, teleport := dropdownContentAttrs(dropdownIn("account", false), DropdownContentProps{})

		assert.Equal(t, "#hui-portals", teleport)
		assert.Equal(t, "menu", attrs["role"])
		assert.Equal(t, "-1", attrs["tabindex"])
		assert.Equal(t, "open", attrs["x-show"])
		assert.Equal(t, "open", attrs[":data-open"])
		assert.Equal(t, "keydown($event)", attrs["@keydown"])
		assert.Equal(t, "select($event)", attrs["@click"])
		assert.Equal(t, "hover($event)", attrs["@pointermove"])
		assert.NotContains(t, attrs, "popover")
	})

	t.Run("content is anchored below the trigger by default", func(t *testing.T) {
		attrs, _ := dropdownContentAttrs(dropdownIn("account", false), DropdownContentProps{})

		assert.Equal(t, "position: fixed; position-anchor: --account; position-area: bottom span-right; position-try-fallbacks: flip-block; inset: auto; margin: 0; display: none", attrs["style"])
	})

	t.Run("side, align and offset", func(t *testing.T) {
		cases := map[string]struct {
			props DropdownContentProps
			want  string
		}{
			"top end":       {DropdownContentProps{Side: "top", Align: "end"}, "position-area: top span-left; position-try-fallbacks: flip-block"},
			"bottom center": {DropdownContentProps{Align: "center", SideOffset: 4}, "position-area: bottom center; position-try-fallbacks: flip-block; inset: auto; margin: 0; margin-top: 4px"},
			"right start":   {DropdownContentProps{Side: "right", SideOffset: 8}, "position-area: right span-bottom; position-try-fallbacks: flip-inline; inset: auto; margin: 0; margin-left: 8px"},
			"left end":      {DropdownContentProps{Side: "left", Align: "end"}, "position-area: left span-top; position-try-fallbacks: flip-inline"},
			"unknown side":  {DropdownContentProps{Side: "middle", Align: "middle"}, "position-area: bottom span-right"},
		}
		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				assert.Contains(t, anchorPosition("", c.props.Placement), c.want)
			})
		}
	})

	t.Run("a style from Attrs comes after the positioning", func(t *testing.T) {
		attrs, _ := dropdownContentAttrs(dropdownIn("account", false), DropdownContentProps{
			Base: Base{Attrs: templ.Attributes{"style": "margin-top: 1rem"}},
		})

		assert.True(t, strings.HasSuffix(attrString(attrs["style"]), "; margin: 0; display: none; margin-top: 1rem"))
	})

	t.Run("items are buttons that cannot submit by accident", func(t *testing.T) {
		text := gofakeit.Word()

		html := testRender(t, DropdownItem(DropdownItemProps{Class: "item"}), templ.Raw(text))

		assert.Equal(t, `<button class="item" data-enabled role="menuitem" tabindex="-1" type="button">`+text+`</button>`, html)
	})

	t.Run("a submit item keeps its form fields", func(t *testing.T) {
		name := gofakeit.Word()

		attrs, isLink := dropdownItemAttrs(DropdownItemProps{Type: "submit", Name: name, Form: "logout"})

		assert.False(t, isLink)
		assert.Equal(t, "submit", attrs["type"])
		assert.Equal(t, name, attrs["name"])
		assert.Equal(t, "logout", attrs["form"])
	})

	t.Run("an item with Href is a link", func(t *testing.T) {
		href := "/" + gofakeit.Word()

		html := testRender(t, DropdownItem(DropdownItemProps{Href: href}))

		assert.Equal(t, `<a data-enabled href="`+href+`" role="menuitem" tabindex="-1"></a>`, html)
	})

	t.Run("a disabled link item stays a menu item", func(t *testing.T) {
		attrs, isLink := dropdownItemAttrs(DropdownItemProps{Href: "/" + gofakeit.Word(), Disabled: true})

		assert.True(t, isLink)
		assert.Equal(t, "menuitem", attrs["role"])
		assert.Equal(t, "true", attrs["aria-disabled"])
		assert.Equal(t, true, attrs["data-disabled"])
		assert.NotContains(t, attrs, "href")
	})

	t.Run("item role cannot be overridden", func(t *testing.T) {
		attrs, _ := dropdownItemAttrs(DropdownItemProps{Base: Base{Attrs: templ.Attributes{"role": "option"}}})

		assert.Equal(t, "menuitem", attrs["role"])
	})

	t.Run("group is labelled by its label", func(t *testing.T) {
		html := testRenderIn(t, WithIDScope(context.Background(), "menu"),
			testWithChildren(DropdownGroup(DropdownGroupProps{}), DropdownLabel(DropdownLabelProps{}), DropdownItem(DropdownItemProps{})),
		)

		assert.Contains(t, html, `<div aria-labelledby="menu-1-label" id="menu-1" role="group"><div id="menu-1-label"></div>`)
	})

	t.Run("separator", func(t *testing.T) {
		html := testRender(t, DropdownSeparator(DropdownSeparatorProps{Class: "my-1 h-px"}))

		assert.Equal(t, `<div class="my-1 h-px" role="separator"></div>`, html)
	})

	t.Run("parts work outside a dropdown", func(t *testing.T) {
		trigger := dropdownTriggerAttrs(context.Background(), DropdownTriggerProps{Base: Base{ID: "menu-button"}})
		content, teleport := dropdownContentAttrs(context.Background(), DropdownContentProps{Base: Base{ID: "menu"}})

		assert.Equal(t, "menu-button", trigger["id"])
		assert.NotContains(t, trigger, "aria-controls")
		assert.NotContains(t, trigger, ":aria-expanded")
		assert.NotContains(t, trigger, "@click")
		assert.Equal(t, "menu", content["id"])
		assert.Empty(t, teleport)
		assert.NotContains(t, content, "aria-labelledby")
		assert.NotContains(t, content, "x-show")
		assert.NotContains(t, content, "@keydown")
		assert.NotContains(t, attrString(content["style"]), "position-anchor")
	})

	t.Run("anchor names are valid CSS identifiers", func(t *testing.T) {
		assert.Equal(t, "--user-menu-1", anchorName("user.menu:1"))
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		ctx := dropdownIn("account", false)

		root := dropdownAttrs(DropdownProps{
			Class:  "relative",
			Attrs:  templ.Attributes{"@pointermove": "track()", "hx-preserve": true},
			Events: Events{OnKeydown: "log($event)"},
		})
		trigger := dropdownTriggerAttrs(ctx, DropdownTriggerProps{
			Class:  "btn",
			Attrs:  templ.Attributes{"style": "color: red", "hx-get": "/menu"},
			Events: Events{OnClick: "count++"},
		})
		content, _ := dropdownContentAttrs(ctx, DropdownContentProps{
			Attrs:  templ.Attributes{"@pointermove": "track()", "role": "listbox"},
			Events: Events{OnKeydown: "log($event)"},
		})

		assert.Equal(t, "relative", root["class"])
		assert.Equal(t, "track()", root["@pointermove"])
		assert.Equal(t, "log($event)", root["@keydown"])
		assert.Equal(t, true, root["hx-preserve"])
		assert.Equal(t, "btn", trigger["class"])
		assert.Equal(t, "anchor-name: --account; color: red", trigger["style"])
		assert.Equal(t, "/menu", trigger["hx-get"])
		assert.Equal(t, "dropdownToggle($el, $event); count++", trigger["@click"])
		assert.Equal(t, "hover($event); track()", content["@pointermove"])
		assert.Equal(t, "keydown($event); log($event)", content["@keydown"])
		assert.Equal(t, "menu", content["role"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, dropdownJS)
		assert.Contains(t, dropdownJS, "window.Alpine.plugin")
		assert.Contains(t, dropdownJS, "alpine:init")
		assert.NotContains(t, dropdownJS, "export")
	})

	t.Run("an item has slots that stay out of its name", func(t *testing.T) {
		left := gofakeit.Word()
		text := gofakeit.Word()

		html := testRender(t, DropdownItem(DropdownItemProps{LeftSlot: templ.Raw(left)}), templ.Raw(text))

		assert.Contains(t, html, `<span aria-hidden="true" data-slot="left">`+left+`</span>`+text+`</button>`)
	})

	t.Run("checkbox item is a menu item that toggles", func(t *testing.T) {
		text := gofakeit.Word()

		html := testRender(t, DropdownCheckboxItem(DropdownCheckboxItemProps{Class: "item"}), templ.Raw(text))

		assert.Equal(t, `<button :aria-checked="ariaChecked()" :data-checked="checked" @click="dropdownCheck()" aria-checked="false" class="item" data-enabled role="menuitemcheckbox" tabindex="-1" type="button" x-data="huiDropdownCheckbox({ checked: false })" x-modelable="checked">`+text+`</button>`, html)
	})

	t.Run("a checked item carries its value for a form", func(t *testing.T) {
		name := gofakeit.Word()
		value := gofakeit.Word()

		html := testRender(t, DropdownCheckboxItem(DropdownCheckboxItemProps{Name: name, Value: value, Form: "columns", Checked: true}))

		assert.Contains(t, html, `aria-checked="true"`)
		assert.Contains(t, html, `data-checked`)
		assert.Contains(t, html, `form="columns" name="`+name+`"`)
		assert.Contains(t, html, `<input data-hui-checked form="columns" name="`+name+`" type="hidden" value="`+value+`">`)
	})

	t.Run("an unchecked item sends nothing", func(t *testing.T) {
		c := dropdownCheckboxItemAttrs(DropdownCheckboxItemProps{Name: gofakeit.Word()})

		assert.NotContains(t, c.attrs, "name")
		assert.Nil(t, c.input)
	})

	t.Run("a disabled checkable item is disabled", func(t *testing.T) {
		c := dropdownCheckboxItemAttrs(DropdownCheckboxItemProps{Disabled: true})

		assert.Equal(t, true, c.attrs["disabled"])
		assert.Equal(t, true, c.attrs["data-disabled"])
		assert.NotContains(t, c.attrs, "data-enabled")
	})

	t.Run("checkable item role cannot be overridden", func(t *testing.T) {
		c := dropdownCheckboxItemAttrs(DropdownCheckboxItemProps{Base: Base{Attrs: templ.Attributes{"role": "checkbox"}}})

		assert.Equal(t, "menuitemcheckbox", c.attrs["role"])
	})

	t.Run("checkable item chains a custom click", func(t *testing.T) {
		c := dropdownCheckboxItemAttrs(DropdownCheckboxItemProps{Base: Base{Events: Events{OnClick: DropdownHide}}})

		assert.Equal(t, "dropdownCheck(); dropdownHide()", c.attrs["@click"])
	})

	t.Run("radio group checks the item with its value", func(t *testing.T) {
		name := gofakeit.Word()

		html := testRenderIn(t, WithIDScope(context.Background(), "menu"),
			testWithChildren(DropdownRadioGroup(DropdownRadioGroupProps{Name: name, Value: "date", Form: "sort"}),
				DropdownLabel(DropdownLabelProps{}),
				DropdownRadioItem(DropdownRadioItemProps{Value: "name"}),
				DropdownRadioItem(DropdownRadioItemProps{Value: "date"}),
			),
		)

		assert.Contains(t, html, `<div aria-labelledby="menu-1-label" id="menu-1" role="group"><div id="menu-1-label"></div>`)
		assert.Equal(t, 2, testCount(html, `role="menuitemradio"`))
		assert.Equal(t, 1, testCount(html, `aria-checked="true"`))
		assert.Contains(t, html, `<input data-hui-checked form="sort" name="`+name+`" type="hidden" value="date">`)
		assert.Contains(t, html, `@click="dropdownPick()"`)
	})

	t.Run("a disabled radio group disables its items", func(t *testing.T) {
		ctx, _ := dropdownRadioGroupContext(context.Background(), DropdownRadioGroupProps{Disabled: true})

		c := dropdownRadioItemAttrs(ctx, DropdownRadioItemProps{Value: gofakeit.Word()})

		assert.Equal(t, true, c.attrs["disabled"])
		assert.Equal(t, true, dropdownRadioGroupAttrs(DropdownRadioGroupProps{Disabled: true})["data-disabled"])
	})

	t.Run("a radio item works outside a group", func(t *testing.T) {
		c := dropdownRadioItemAttrs(context.Background(), DropdownRadioItemProps{Checked: true})

		assert.Equal(t, "true", c.attrs["aria-checked"])
		assert.Equal(t, "huiDropdownRadio({ checked: true })", c.attrs["x-data"])
	})

	t.Run("indicators show the state before Alpine starts", func(t *testing.T) {
		off := testRender(t, DropdownCheckboxItem(DropdownCheckboxItemProps{LeftSlot: CheckboxIndicator(CheckboxIndicatorProps{})}))
		on := testRender(t, DropdownRadioItem(DropdownRadioItemProps{Checked: true, LeftSlot: RadioIndicator(RadioIndicatorProps{})}))

		assert.Contains(t, off, `style="display: none"`)
		assert.NotContains(t, on, `style="display: none"`)
		assert.Contains(t, on, `x-show="checked"`)
	})

	t.Run("checkable items ignore a surrounding field", func(t *testing.T) {
		html := testRender(t, Field(FieldProps{Name: "notify", Disabled: true, Error: gofakeit.Sentence()}), DropdownCheckboxItem(DropdownCheckboxItemProps{}))

		button := regexp.MustCompile(`<button[^>]*>`).FindString(html)

		assert.NotContains(t, button, `id=`)
		assert.NotContains(t, button, `aria-invalid`)
		assert.NotContains(t, button, `disabled`)
		assert.Contains(t, button, `data-enabled`)
	})
}
