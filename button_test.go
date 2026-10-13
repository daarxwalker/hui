package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestButton(t *testing.T) {
	t.Run("defaults to type button", func(t *testing.T) {
		text := gofakeit.Name()

		html := testRender(t, Button(ButtonProps{Class: "btn"}), templ.Raw(text))

		assert.Equal(t, `<button class="btn" data-enabled type="button">`+text+`</button>`, html)
	})

	t.Run("submits when asked to", func(t *testing.T) {
		name := gofakeit.Word()
		value := gofakeit.Word()

		html := testRender(t, Button(ButtonProps{Type: "submit", Name: name, Value: value, Form: "signup"}))

		assert.Contains(t, html, `type="submit"`)
		assert.Contains(t, html, `name="`+name+`"`)
		assert.Contains(t, html, `value="`+value+`"`)
		assert.Contains(t, html, `form="signup"`)
	})

	t.Run("disabled uses the native attribute", func(t *testing.T) {
		html := testRender(t, Button(ButtonProps{Disabled: true}))

		assert.Contains(t, html, ` disabled`)
		assert.Contains(t, html, `data-disabled`)
		assert.NotContains(t, html, `aria-disabled`)
		assert.NotContains(t, html, `data-enabled`)
	})

	t.Run("pending implies disabled", func(t *testing.T) {
		html := testRender(t, Button(ButtonProps{Pending: true}))

		assert.Contains(t, html, ` disabled`)
		assert.Contains(t, html, `data-disabled`)
		assert.Contains(t, html, `data-pending`)
		assert.Contains(t, html, `aria-busy="true"`)
		assert.NotContains(t, html, `data-enabled`)
	})

	t.Run("never ships a script, not even when it talks to the server", func(t *testing.T) {
		hx := templ.Attributes{"hx-post": "/" + gofakeit.Word()}

		assert.NotContains(t, testRender(t, Button(ButtonProps{Attrs: hx})), "<script>")
	})

	t.Run("the pending script is opt-in and rendered once", func(t *testing.T) {
		ctx := templ.InitializeContext(context.Background())

		html := testRenderIn(t, ctx, templ.Join(PendingScript(), PendingScript()))

		assert.Equal(t, 1, testCount(html, "<script>"))
		assert.Contains(t, html, "headlessui.pending")
	})

	t.Run("the pending script works without Alpine", func(t *testing.T) {
		require.NotEmpty(t, pendingJS)
		assert.Contains(t, pendingJS, "htmx:before:request")
		assert.Contains(t, pendingJS, "htmx:finally:request")
		assert.NotContains(t, pendingJS, "Alpine")
		assert.NotContains(t, pendingJS, "export")
	})

	t.Run("renders a link when Href is set", func(t *testing.T) {
		href := "/" + gofakeit.Word()
		text := gofakeit.Name()

		html := testRender(t, Button(ButtonProps{Href: href, Target: "_blank"}), templ.Raw(text))

		assert.Equal(t, `<a data-enabled href="`+href+`" rel="noopener noreferrer" target="_blank">`+text+`</a>`, html)
		assert.NotContains(t, html, `type="button"`)
	})

	t.Run("a disabled link loses its href", func(t *testing.T) {
		html := testRender(t, Button(ButtonProps{Href: "/" + gofakeit.Word(), Disabled: true}))

		assert.Contains(t, html, `role="link"`)
		assert.Contains(t, html, `aria-disabled="true"`)
		assert.Contains(t, html, `tabindex="-1"`)
		assert.Contains(t, html, `data-disabled`)
		assert.NotContains(t, html, `href`)
		assert.NotContains(t, html, `data-enabled`)
	})

	t.Run("slots wrap the children and stay out of the accessible name", func(t *testing.T) {
		text := gofakeit.Name()
		left := gofakeit.Word()
		right := gofakeit.Word()

		html := testRender(t, Button(ButtonProps{LeftSlot: templ.Raw(left), RightSlot: templ.Raw(right)}), templ.Raw(text))

		assert.Equal(t, `<button data-enabled type="button"><span aria-hidden="true" data-slot="left">`+left+`</span>`+text+`<span aria-hidden="true" data-slot="right">`+right+`</span></button>`, html)
	})

	t.Run("a link has slots too", func(t *testing.T) {
		left := gofakeit.Word()

		html := testRender(t, Button(ButtonProps{Href: "/" + gofakeit.Word(), LeftSlot: templ.Raw(left)}))

		assert.Contains(t, html, `<span aria-hidden="true" data-slot="left">`+left+`</span></a>`)
	})

	t.Run("sanitizes the href", func(t *testing.T) {
		html := testRender(t, Button(ButtonProps{Href: "javascript:alert(1)"}))

		assert.NotContains(t, html, "javascript:")
		assert.Contains(t, html, string(templ.FailedSanitizationURL))
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		html := testRender(t, Button(ButtonProps{
			Class:  "btn",
			Attrs:  templ.Attributes{"@click": "track()", "x-ref": "save", "type": "reset"},
			Events: Events{OnClick: "$store.cart.add(1)"},
		}))

		assert.Contains(t, html, `class="btn"`)
		assert.Contains(t, html, `x-ref="save"`)
		assert.Contains(t, html, `@click="track(); $store.cart.add(1)"`)
		assert.Contains(t, html, `type="button"`)
		assert.NotContains(t, html, `type="reset"`)
	})
}
