package headlessui

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

func TestMeter(t *testing.T) {
	ctx := context.Background()

	t.Run("renders a meter with its value as a CSS variable", func(t *testing.T) {
		html := testRender(t, Meter(MeterProps{Value: 40, Label: "Storage"}))

		assert.Equal(t, `<div aria-label="Storage" aria-valuemax="100" aria-valuemin="0" aria-valuenow="40" data-level="medium" role="meter" style="--hui-meter: 40%"></div>`, html)
	})

	t.Run("value is clamped into its range", func(t *testing.T) {
		attrs := meterAttrs(ctx, MeterProps{Value: 5, Min: 10, Max: 20})

		assert.Equal(t, "10", attrs["aria-valuenow"])
		assert.Equal(t, "--hui-meter: 0%", attrs["style"])
	})

	t.Run("low and high split the range into levels", func(t *testing.T) {
		levels := map[float64]string{10: "low", 50: "medium", 90: "high"}
		for value, want := range levels {
			assert.Equal(t, want, meterAttrs(ctx, MeterProps{Value: value, Low: 25, High: 75})["data-level"])
		}
	})

	t.Run("optimum says which level is good", func(t *testing.T) {
		cases := []struct {
			optimum string
			value   float64
			want    string
		}{
			{"high", 90, "optimum"},
			{"high", 50, "suboptimal"},
			{"high", 10, "critical"},
			{"low", 10, "optimum"},
			{"low", 90, "critical"},
			{"medium", 50, "optimum"},
			{"medium", 90, "suboptimal"},
		}
		for _, c := range cases {
			attrs := meterAttrs(ctx, MeterProps{Value: c.value, Low: 25, High: 75, Optimum: c.optimum})
			assert.Equal(t, c.want, attrs["data-state"], c.optimum)
		}
		assert.NotContains(t, meterAttrs(ctx, MeterProps{Value: 50}), "data-state")
	})

	t.Run("is labelled by the Field label", func(t *testing.T) {
		html := testRender(t, Field(FieldProps{Name: "disk"}), Label(LabelProps{}), Meter(MeterProps{Value: 1}))

		assert.Contains(t, html, `aria-labelledby="disk-label"`)
		assert.Contains(t, html, `id="disk" role="meter"`)
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		attrs := meterAttrs(ctx, MeterProps{Value: 10, Base: Base{
			Class: "h-2",
			Attrs: templ.Attributes{"data-level": "custom", "title": "10%"},
		}})

		assert.Equal(t, "medium", attrs["data-level"])
		assert.Equal(t, "10%", attrs["title"])
		assert.Equal(t, "h-2", attrs["class"])
	})
}
