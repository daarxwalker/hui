package headlessui

import (
	"context"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestField(t *testing.T) {
	t.Run("wires label, input, description and error together", func(t *testing.T) {
		name := gofakeit.Word()
		email := gofakeit.Email()
		description := gofakeit.Sentence(4)
		message := gofakeit.Sentence(3)

		html := testRender(t,
			Field(FieldProps{Name: name, Value: email, Description: description, Error: message}),
			Label(LabelProps{}),
			TextInput(TextInputProps{Type: "email"}),
			Description(DescriptionProps{}),
			Error(ErrorProps{}),
		)

		assert.Contains(t, html, `for="`+name+`"`)
		assert.Contains(t, html, `id="`+name+`"`)
		assert.Contains(t, html, `name="`+name+`"`)
		assert.Contains(t, html, `value="`+email+`"`)
		assert.Contains(t, html, `aria-describedby="`+name+`-description `+name+`-error"`)
		assert.Contains(t, html, `id="`+name+`-description">`+description+`</p>`)
		assert.Contains(t, html, `id="`+name+`-error"`)
		assert.Contains(t, html, `>`+message+`</p>`)
	})

	t.Run("an error always means invalid", func(t *testing.T) {
		props := FieldProps{Name: gofakeit.Word(), Error: gofakeit.Sentence(3)}

		require.False(t, props.Invalid)
		assert.True(t, props.normalize(context.Background()).Invalid)

		html := testRender(t, Field(props), TextInput(TextInputProps{}))

		assert.Contains(t, html, `aria-invalid="true"`)
		assert.Contains(t, html, `data-invalid`)
	})

	t.Run("state attributes land on wrapper, label and input", func(t *testing.T) {
		html := testRender(t,
			Field(FieldProps{Name: gofakeit.Word(), Disabled: true, Readonly: true, Required: true}),
			Label(LabelProps{}),
			TextInput(TextInputProps{}),
		)

		for _, state := range []string{"data-disabled", "data-readonly", "data-required"} {
			assert.Equal(t, 3, testCount(html, state), state)
		}
	})

	t.Run("no description or error means no aria-describedby", func(t *testing.T) {
		html := testRender(t,
			Field(FieldProps{Name: gofakeit.Word()}),
			Label(LabelProps{}),
			TextInput(TextInputProps{}),
			Description(DescriptionProps{}),
			Error(ErrorProps{}),
		)

		assert.NotContains(t, html, "aria-describedby")
		assert.NotContains(t, html, "<p")
	})

	t.Run("generates an ID when neither ID nor Name is set", func(t *testing.T) {
		html := testRender(t,
			Field(FieldProps{Description: gofakeit.Sentence(3)}),
			Label(LabelProps{}),
			TextInput(TextInputProps{}),
			Description(DescriptionProps{}),
		)

		match := regexp.MustCompile(`for="(field-[^"]+)"`).FindStringSubmatch(html)
		require.NotNil(t, match, html)

		id := match[1]
		assert.Contains(t, html, `id="`+id+`"`)
		assert.Contains(t, html, `aria-describedby="`+id+`-description"`)
		assert.Contains(t, html, `id="`+id+`-description"`)
	})

	t.Run("is the single source of truth for the input", func(t *testing.T) {
		html := testRender(t,
			Field(FieldProps{ID: "email", Name: "email", Disabled: true}),
			TextInput(TextInputProps{ID: "own", Name: "own", Value: "own", Required: true}),
		)

		assert.Contains(t, html, `id="email"`)
		assert.Contains(t, html, `name="email"`)
		assert.Contains(t, html, ` disabled`)
		assert.NotContains(t, html, `"own"`)
		assert.NotContains(t, html, ` required`)
	})

	t.Run("the wrapper is a swap target of its own", func(t *testing.T) {
		name := gofakeit.Word()

		html := testRender(t,
			Field(FieldProps{Name: name, Error: gofakeit.Sentence(3)}),
			TextInput(TextInputProps{}),
			Error(ErrorProps{}),
		)

		assert.Contains(t, html, `<div data-field="`+name+`" data-invalid id="`+name+`-field">`)
		assert.Equal(t, 1, testCount(html, `id="`+name+`"`))
	})

	t.Run("a scope makes generated IDs repeatable across renders", func(t *testing.T) {
		description := gofakeit.Sentence(3)
		render := func() string {
			return testRenderIn(t, WithIDScope(context.Background(), "signup"),
				Field(FieldProps{Description: description}),
				Label(LabelProps{}),
				TextInput(TextInputProps{}),
				Description(DescriptionProps{}),
			)
		}

		first, second := render(), render()

		assert.Equal(t, first, second)
		assert.Contains(t, first, `id="signup-1"`)
		assert.Contains(t, first, `id="signup-1-field"`)
		assert.Contains(t, first, `aria-describedby="signup-1-description"`)
	})
}
