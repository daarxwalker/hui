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

func TestSwitch(t *testing.T) {
	ctx := context.Background()

	t.Run("renders a button with the switch role and a thumb", func(t *testing.T) {
		name := gofakeit.Word()

		html := testRender(t, Switch(SwitchProps{Name: name, Checked: true}), SwitchThumb(SwitchThumbProps{Class: "block size-4"}))
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `@click="toggle()" aria-checked="true" data-checked name="`+name+`" role="switch" type="button" value="on"`)
		assert.Contains(t, html, `<input data-hui-checked name="`+name+`" type="hidden" value="on">`)
		assert.Contains(t, html, `<span :data-checked="checked" aria-hidden="true" class="block size-4" data-checked></span></button>`)
		assert.NotContains(t, html, `@keydown`)
	})

	t.Run("an unchecked switch submits nothing and its thumb is unchecked", func(t *testing.T) {
		c := switchAttrs(ctx, SwitchProps{Name: gofakeit.Word()})
		thumb := switchThumbAttrs(checkableContext(ctx, c.state), SwitchThumbProps{})

		assert.Nil(t, c.input)
		assert.Equal(t, "false", c.attrs["aria-checked"])
		assert.NotContains(t, thumb, "data-checked")
	})

	t.Run("submits its own value when checked", func(t *testing.T) {
		name := gofakeit.Word()
		value := gofakeit.Word()

		c := switchAttrs(ctx, SwitchProps{Name: name, Value: value, Checked: true})

		assert.Equal(t, name, c.input["name"])
		assert.Equal(t, value, c.input["value"])
	})

	t.Run("carries its own state without a Field", func(t *testing.T) {
		c := switchAttrs(ctx, SwitchProps{
			ID:       "terms",
			Disabled: true,
			Required: true,
			Invalid:  true,
		})

		assert.Equal(t, "terms", c.attrs["id"])
		assert.Equal(t, true, c.attrs["disabled"])
		assert.Equal(t, "true", c.attrs["aria-required"])
		assert.Equal(t, "true", c.attrs["aria-invalid"])
		assert.Equal(t, true, c.attrs["data-disabled"])
		assert.Equal(t, true, c.attrs["data-required"])
		assert.Equal(t, true, c.attrs["data-invalid"])
		assert.NotNil(t, c.required)
	})

	t.Run("is wired into a Field", func(t *testing.T) {
		name := gofakeit.Word()

		html := testRender(t,
			Field(FieldProps{Name: name, Value: "yes", Description: gofakeit.Sentence(4), Error: gofakeit.Sentence(3)}),
			Switch(SwitchProps{ID: "own", Name: "own", Value: "own", Checked: true}),
			Label(LabelProps{}),
			Description(DescriptionProps{}),
			Error(ErrorProps{}),
		)

		assert.Contains(t, html, `for="`+name+`"`)
		assert.Contains(t, html, `id="`+name+`"`)
		assert.Contains(t, html, `aria-checked="true"`)
		assert.Contains(t, html, `name="`+name+`"`)
		assert.Contains(t, html, `value="yes"`)
		assert.Contains(t, html, `aria-describedby="`+name+`-description `+name+`-error"`)
		assert.Contains(t, html, `aria-invalid="true"`)
		assert.NotContains(t, html, `"own"`)
	})

	t.Run("ignores readonly because a switch has none", func(t *testing.T) {
		ctx := context.WithValue(ctx, fieldKey, FieldProps{Name: gofakeit.Word(), Readonly: true})

		c := switchAttrs(ctx, SwitchProps{})

		assert.NotContains(t, c.attrs, "readonly")
		assert.NotContains(t, c.attrs, "data-readonly")
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		c := switchAttrs(ctx, SwitchProps{Base: Base{
			Class: "peer rounded-full",
			Attrs: templ.Attributes{
				"class":      "data-checked:bg-blue-600",
				"hx-post":    "/settings/notifications",
				"type":       "submit",
				"role":       "button",
				"aria-label": "Notifications",
			},
			Events: Events{OnChange: "$store.settings.dirty = true"},
		}})

		assert.Equal(t, "data-checked:bg-blue-600 peer rounded-full", c.attrs["class"])
		assert.Equal(t, "/settings/notifications", c.attrs["hx-post"])
		assert.Equal(t, "Notifications", c.attrs["aria-label"])
		assert.Equal(t, "$store.settings.dirty = true", c.attrs["@change"])
		assert.Equal(t, "button", c.attrs["type"])
		assert.Equal(t, "switch", c.attrs["role"])
	})

	t.Run("belongs to a form outside of it", func(t *testing.T) {
		own := switchAttrs(ctx, SwitchProps{Form: "settings"})
		fallback := testRender(t, Field(FieldProps{Name: gofakeit.Word()}), Switch(SwitchProps{Form: "settings"}))

		assert.Equal(t, "settings", own.attrs["form"])
		assert.Contains(t, fallback, `form="settings"`)
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, switchJS)
		assert.Contains(t, switchJS, "window.Alpine.plugin")
	})
}
