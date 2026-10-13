package headlessui

import (
	"context"
	_ "embed"
	"math"
	"slices"
	"strconv"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/range.min.js
var rangeJS string

var rangeScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + rangeJS + "</script>")))

func RangeScript() templ.Component {
	return rangeScript.Once()
}

var rangeThumbEvents = Events{
	OnInput: "update($el)",
}

type RangeProps struct {
	Base
	Name     string
	Form     string
	Value    []float64
	Min      float64
	Max      float64
	Step     float64
	Disabled bool
	Invalid  bool
}

type RangeThumbProps struct {
	Base
	Index int
}

type rangeKeyType struct{}

var rangeKey = rangeKeyType{}

type rangeState struct {
	field   FieldProps
	inField bool
	min     float64
	max     float64
	step    float64
	values  []float64
}

func newRangeState(ctx context.Context, p RangeProps) rangeState {
	f, inField := controlFrom(ctx, FieldProps{Name: p.Name, Form: p.Form, Disabled: p.Disabled, Invalid: p.Invalid})
	f.Readonly, f.Required = false, false

	s := rangeState{field: f, inField: inField, min: p.Min, max: p.Max, step: p.Step}
	if s.min == 0 && s.max == 0 {
		s.max = 100
	}
	if s.step <= 0 {
		s.step = 1
	}
	s.values = make([]float64, len(p.Value))
	for i, v := range p.Value {
		s.values[i] = min(max(v, s.min), s.max)
	}
	slices.Sort(s.values)
	return s
}

func (s rangeState) value(index int) float64 {
	if index < 0 || index >= len(s.values) {
		return s.min
	}
	return s.values[index]
}

func (s rangeState) percent(v float64) string {
	span := s.max - s.min
	if span <= 0 {
		return "0"
	}
	return formatNumber(math.Round((v-s.min)/span*100000) / 1000)
}

func (s rangeState) fill() (start, end string) {
	switch len(s.values) {
	case 0:
		return "0", "0"
	case 1:
		return "0", s.percent(s.values[0])
	default:
		return s.percent(s.values[0]), s.percent(s.values[len(s.values)-1])
	}
}

func formatNumber(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func rangeAttrs(ctx context.Context, p RangeProps) templ.Attributes {
	s := newRangeState(ctx, p)
	start, end := s.fill()

	own := make(templ.Attributes, 6)
	own["x-data"] = "huiRange"
	own["style"] = "--hui-range-start: " + start + "%; --hui-range-end: " + end + "%"

	return p.merge(s.field.state(own))
}

func rangeThumbAttrs(ctx context.Context, p RangeThumbProps) templ.Attributes {
	s, inRange := ctx.Value(rangeKey).(rangeState)
	if !inRange {
		s = newRangeState(ctx, RangeProps{})
	}

	own := make(templ.Attributes, 12)
	own["type"] = "range"
	own["min"] = formatNumber(s.min)
	own["max"] = formatNumber(s.max)
	own["step"] = formatNumber(s.step)
	own["value"] = formatNumber(s.value(p.Index))

	if s.inField {
		p.ID = s.field.ID
		if p.Index > 0 {
			p.ID += "-" + strconv.Itoa(p.Index)
		}
		if by := s.field.describedBy(); by != "" {
			own["aria-describedby"] = by
		}
	}
	if s.field.Name != "" {
		own["name"] = s.field.Name
	}
	if s.field.Form != "" {
		own["form"] = s.field.Form
	}
	if s.field.Disabled {
		own["disabled"] = true
	}
	if s.field.Invalid {
		own["aria-invalid"] = "true"
	}

	if !inRange {
		return p.merge(s.field.state(own))
	}
	return p.mergeEvents(s.field.state(own), rangeThumbEvents)
}
