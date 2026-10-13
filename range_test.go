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

func rangeIn(ctx context.Context, props RangeProps) context.Context {
	return context.WithValue(ctx, rangeKey, newRangeState(ctx, props))
}

func rangeHTML(t *testing.T, ctx context.Context, c templ.Component, children ...templ.Component) string {
	t.Helper()
	return regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(testRenderIn(t, ctx, c, children...), "")
}

func TestRange(t *testing.T) {
	ctx := context.Background()

	t.Run("renders a native range input per thumb and ships its script once", func(t *testing.T) {
		ctx := templ.InitializeContext(context.Background())
		name := gofakeit.Word()

		html := testRenderIn(t, ctx, templ.Join(
			testWithChildren(Range(RangeProps{Name: name, Value: []float64{20, 80}}), RangeThumb(RangeThumbProps{Index: 0}), RangeThumb(RangeThumbProps{Index: 1})),
			testWithChildren(Range(RangeProps{}), RangeThumb(RangeThumbProps{})),
		))

		assert.Equal(t, 3, strings.Count(html, `<input `))
		assert.Equal(t, 2, strings.Count(html, `name="`+name+`"`))
		assert.NotContains(t, html, "<script>")
	})

	t.Run("each thumb takes its value from the root", func(t *testing.T) {
		ctx := rangeIn(ctx, RangeProps{Min: 10, Max: 500, Step: 5, Value: []float64{50, 250}})

		first := rangeThumbAttrs(ctx, RangeThumbProps{Index: 0})
		second := rangeThumbAttrs(ctx, RangeThumbProps{Index: 1})

		assert.Equal(t, "range", first["type"])
		assert.Equal(t, "10", first["min"])
		assert.Equal(t, "500", first["max"])
		assert.Equal(t, "5", first["step"])
		assert.Equal(t, "50", first["value"])
		assert.Equal(t, "250", second["value"])
	})

	t.Run("defaults to the native zero to a hundred in steps of one", func(t *testing.T) {
		attrs := rangeThumbAttrs(rangeIn(ctx, RangeProps{}), RangeThumbProps{})

		assert.Equal(t, "0", attrs["min"])
		assert.Equal(t, "100", attrs["max"])
		assert.Equal(t, "1", attrs["step"])
		assert.Equal(t, "0", attrs["value"])
	})

	t.Run("decimal steps are rendered as written", func(t *testing.T) {
		attrs := rangeThumbAttrs(rangeIn(ctx, RangeProps{Max: 1, Step: 0.05, Value: []float64{0.35}}), RangeThumbProps{})

		assert.Equal(t, "0.05", attrs["step"])
		assert.Equal(t, "0.35", attrs["value"])
	})

	t.Run("values are clamped and sorted so thumbs never start crossed", func(t *testing.T) {
		s := newRangeState(ctx, RangeProps{Value: []float64{120, -5, 40}})

		assert.Equal(t, []float64{0, 40, 100}, s.values)
	})

	t.Run("a thumb without a value starts at the minimum", func(t *testing.T) {
		attrs := rangeThumbAttrs(rangeIn(ctx, RangeProps{Min: 20, Max: 80, Value: []float64{50}}), RangeThumbProps{Index: 1})

		assert.Equal(t, "20", attrs["value"])
	})

	t.Run("the fill spans from the minimum to a single thumb", func(t *testing.T) {
		attrs := rangeAttrs(ctx, RangeProps{Value: []float64{30}})

		assert.Equal(t, "--hui-range-start: 0%; --hui-range-end: 30%", attrs["style"])
	})

	t.Run("the fill spans between the first and the last thumb", func(t *testing.T) {
		attrs := rangeAttrs(ctx, RangeProps{Min: 100, Max: 400, Value: []float64{200, 250, 300}})

		assert.Equal(t, "--hui-range-start: 33.333%; --hui-range-end: 66.667%", attrs["style"])
	})

	t.Run("the root hands interactivity to Alpine", func(t *testing.T) {
		root := rangeAttrs(ctx, RangeProps{})
		thumb := rangeThumbAttrs(rangeIn(ctx, RangeProps{}), RangeThumbProps{})

		assert.Equal(t, "huiRange", root["x-data"])
		assert.Equal(t, "update($el)", thumb["@input"])
	})

	t.Run("a thumb outside a Range is a plain range input", func(t *testing.T) {
		attrs := rangeThumbAttrs(ctx, RangeThumbProps{})

		assert.Equal(t, "range", attrs["type"])
		assert.NotContains(t, attrs, "@input")
	})

	t.Run("carries its own state without a Field", func(t *testing.T) {
		props := RangeProps{Name: gofakeit.Word(), Disabled: true, Invalid: true}

		root := rangeAttrs(ctx, props)
		thumb := rangeThumbAttrs(rangeIn(ctx, props), RangeThumbProps{ID: "min"})

		assert.Equal(t, true, root["data-disabled"])
		assert.Equal(t, true, root["data-invalid"])
		assert.NotContains(t, root, "name")
		assert.Equal(t, "min", thumb["id"])
		assert.Equal(t, props.Name, thumb["name"])
		assert.Equal(t, true, thumb["disabled"])
		assert.Equal(t, "true", thumb["aria-invalid"])
		assert.Equal(t, true, thumb["data-disabled"])
		assert.Equal(t, true, thumb["data-invalid"])
	})

	t.Run("Field is the single source of truth", func(t *testing.T) {
		name := gofakeit.Word()
		ctx := context.WithValue(ctx, fieldKey, FieldProps{
			Name:        name,
			Description: gofakeit.Sentence(3),
			Error:       gofakeit.Sentence(3),
			Disabled:    true,
		}.normalize(ctx))
		props := RangeProps{Name: "own", Value: []float64{10, 90}}
		ctx = rangeIn(ctx, props)

		first := rangeThumbAttrs(ctx, RangeThumbProps{ID: "own", Index: 0})
		second := rangeThumbAttrs(ctx, RangeThumbProps{Index: 1})

		assert.Equal(t, name, first["id"])
		assert.Equal(t, name+"-1", second["id"])
		for _, thumb := range []templ.Attributes{first, second} {
			assert.Equal(t, name, thumb["name"])
			assert.Equal(t, true, thumb["disabled"])
			assert.Equal(t, "true", thumb["aria-invalid"])
			assert.Equal(t, name+"-description "+name+"-error", thumb["aria-describedby"])
		}
	})

	t.Run("the label points at the first thumb", func(t *testing.T) {
		name := gofakeit.Word()

		html := rangeHTML(t, ctx,
			Field(FieldProps{Name: name}),
			Label(LabelProps{}),
			testWithChildren(Range(RangeProps{Value: []float64{10, 90}}), RangeThumb(RangeThumbProps{Index: 0}), RangeThumb(RangeThumbProps{Index: 1})),
		)

		assert.Contains(t, html, `for="`+name+`"`)
		assert.Contains(t, html, `id="`+name+`"`)
		assert.Contains(t, html, `id="`+name+`-1"`)
	})

	t.Run("required and readonly never reach the thumbs", func(t *testing.T) {
		html := rangeHTML(t, ctx,
			Field(FieldProps{Name: gofakeit.Word(), Required: true, Readonly: true}),
			testWithChildren(Range(RangeProps{}), RangeThumb(RangeThumbProps{})),
		)

		assert.NotContains(t, html, ` required`)
		assert.NotContains(t, html, ` readonly`)
		assert.Equal(t, 1, testCount(html, "data-required"))
		assert.Equal(t, 1, testCount(html, "data-readonly"))
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		root := rangeAttrs(ctx, RangeProps{
			Class: "relative h-2",
			Attrs: templ.Attributes{
				"class":       "w-full",
				"style":       "--track: gray",
				"x-data":      "{}",
				"hx-post":     "/filter",
				"hx-trigger":  "change",
				"hx-include":  "this",
				"data-testid": "price",
			},
		})
		thumb := rangeThumbAttrs(rangeIn(ctx, RangeProps{}), RangeThumbProps{
			Class: "pointer-events-none absolute inset-0",
			Attrs: templ.Attributes{
				"type":       "number",
				"aria-label": "Minimum price",
			},
			Events: Events{OnInput: "$dispatch('price', values)"},
		})

		assert.Equal(t, "w-full relative h-2", root["class"])
		assert.Equal(t, "--hui-range-start: 0%; --hui-range-end: 0%; --track: gray", root["style"])
		assert.Equal(t, "huiRange", root["x-data"])
		assert.Equal(t, "/filter", root["hx-post"])
		assert.Equal(t, "this", root["hx-include"])
		assert.Equal(t, "price", root["data-testid"])
		assert.Equal(t, "pointer-events-none absolute inset-0", thumb["class"])
		assert.Equal(t, "range", thumb["type"])
		assert.Equal(t, "Minimum price", thumb["aria-label"])
		assert.Equal(t, "update($el); $dispatch('price', values)", thumb["@input"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, rangeJS)
		assert.Contains(t, rangeJS, "window.Alpine.plugin")
		assert.Contains(t, rangeJS, "alpine:init")
		assert.NotContains(t, rangeJS, "export")
	})

	t.Run("every thumb belongs to the same form", func(t *testing.T) {
		html := testRender(t, Range(RangeProps{Name: "price", Form: "filters", Value: []float64{10, 20}}),
			RangeThumb(RangeThumbProps{Index: 0}),
			RangeThumb(RangeThumbProps{Index: 1}),
		)

		assert.Equal(t, 2, testCount(html, `form="filters"`))
	})
}
