package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	t.Run("own attributes win over Attrs", func(t *testing.T) {
		base := Base{Attrs: templ.Attributes{"type": "number", "x-model": "email"}}

		attrs := base.merge(templ.Attributes{"type": "text"})

		assert.Equal(t, "text", attrs["type"])
		assert.Equal(t, "email", attrs["x-model"])
	})

	t.Run("ID is applied last", func(t *testing.T) {
		base := Base{ID: "email", Attrs: templ.Attributes{"id": "from-attrs"}}

		attrs := base.merge(templ.Attributes{"id": "own"})

		assert.Equal(t, "email", attrs["id"])
	})

	t.Run("class from Attrs is merged with Class", func(t *testing.T) {
		base := Base{Class: "rounded border", Attrs: templ.Attributes{"class": "w-full"}}

		attrs := base.merge(nil)

		assert.Equal(t, "w-full rounded border", attrs["class"])
	})

	t.Run("a custom MergeClass is used", func(t *testing.T) {
		original := MergeClass
		t.Cleanup(func() { MergeClass = original })
		MergeClass = func(existing, added string) string { return existing + "|" + added }

		attrs := Base{Class: "b", Attrs: templ.Attributes{"class": "a"}}.merge(nil)

		assert.Equal(t, "a|b", attrs["class"])
	})

	t.Run("a non string class in Attrs is replaced", func(t *testing.T) {
		base := Base{Class: "rounded", Attrs: templ.Attributes{"class": true}}

		attrs := base.merge(nil)

		assert.Equal(t, "rounded", attrs["class"])
	})

	t.Run("style from Attrs is appended to the own style", func(t *testing.T) {
		base := Base{Attrs: templ.Attributes{"style": "color: red"}}

		attrs := base.merge(templ.Attributes{"style": "anchor-name: --menu;"})

		assert.Equal(t, "anchor-name: --menu; color: red", attrs["style"])
	})

	t.Run("events are chained after the internal handler", func(t *testing.T) {
		base := Base{Events: Events{OnBlur: "$store.form.validate('email')"}}

		attrs := base.mergeEvents(nil, Events{OnBlur: "focused = false"})

		assert.Equal(t, "focused = false; $store.form.validate('email')", attrs["@blur"])
	})

	t.Run("handler from Attrs is chained, not dropped", func(t *testing.T) {
		base := Base{
			Attrs:  templ.Attributes{"@click": "track()"},
			Events: Events{OnClick: "open()"},
		}

		attrs := base.merge(nil)

		assert.Equal(t, "track(); open()", attrs["@click"])
	})

	t.Run("non string handler in Attrs is replaced", func(t *testing.T) {
		base := Base{
			Attrs:  templ.Attributes{"@focus": true},
			Events: Events{OnFocus: "log()"},
		}

		attrs := base.mergeEvents(nil, Events{OnFocus: "focused = true"})

		assert.Equal(t, "focused = true; log()", attrs["@focus"])
	})

	t.Run("every named event has its own attribute", func(t *testing.T) {
		base := Base{Events: Events{
			OnFocus:   "a()",
			OnBlur:    "b()",
			OnInput:   "c()",
			OnChange:  "d()",
			OnKeydown: "e()",
			OnClick:   "f()",
		}}

		attrs := base.merge(nil)

		assert.Equal(t, templ.Attributes{
			"@focus":   "a()",
			"@blur":    "b()",
			"@input":   "c()",
			"@change":  "d()",
			"@keydown": "e()",
			"@click":   "f()",
		}, attrs)
	})

	t.Run("chain", func(t *testing.T) {
		cases := map[string]struct {
			internal string
			user     string
			want     string
		}{
			"only user":          {"", "f()", "f()"},
			"only internal":      {"focused = true", "", "focused = true"},
			"both":               {"focused = true", "f()", "focused = true; f()"},
			"trailing semicolon": {"focused = true;", "f()", "focused = true; f()"},
			"surrounding spaces": {" focused = true ", " f() ", "focused = true; f()"},
		}
		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, c.want, chain(c.internal, c.user))
			})
		}
	})

	t.Run("generated IDs are unique", func(t *testing.T) {
		ctx := context.Background()
		assert.NotEqual(t, nextID(ctx, "field"), nextID(ctx, "field"))
	})
}
