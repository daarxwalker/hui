package headlessui

import (
	"context"
	"slices"

	"github.com/a-h/templ"
)

type FieldsetProps struct {
	Base
	Name        string
	Value       string
	Values      []string
	Form        string
	Description string
	Error       string
	Disabled    bool
	Required    bool
	Invalid     bool
}

type LegendProps struct {
	Base
}

type fieldsetKeyType struct{}

var fieldsetKey = fieldsetKeyType{}

func fieldsetFrom(ctx context.Context) (FieldsetProps, bool) {
	s, ok := ctx.Value(fieldsetKey).(FieldsetProps)
	return s, ok
}

func (p FieldsetProps) normalize(ctx context.Context) FieldsetProps {
	if p.Error != "" {
		p.Invalid = true
	}
	switch {
	case p.ID != "":
	case p.Name != "":
		p.ID = p.Name
	default:
		p.ID = nextID(ctx, "fieldset")
	}
	return p
}

func (p FieldsetProps) field() FieldProps {
	return FieldProps{
		Base:        Base{ID: p.ID},
		Name:        p.Name,
		Form:        p.Form,
		Description: p.Description,
		Error:       p.Error,
		Disabled:    p.Disabled,
		Required:    p.Required,
		Invalid:     p.Invalid,
	}
}

func (p FieldsetProps) has(value string) bool {
	if value == "" {
		return false
	}
	return value == p.Value || slices.Contains(p.Values, value)
}

func (p FieldsetProps) option(f FieldProps) FieldProps {
	f.option = true
	f.Name = p.Name
	f.Invalid = f.Invalid || p.Invalid
	f.Required = p.Required
	f.Readonly = false
	f.checked = p.has(f.Value)
	return f
}

func fieldsetAttrs(p FieldsetProps) templ.Attributes {
	f := p.field()
	own := f.state(nil)
	if p.Name != "" {
		own["data-field"] = p.Name
	}
	if by := f.describedBy(); by != "" {
		own["aria-describedby"] = by
	}
	if p.Disabled {
		own["disabled"] = true
	}
	if p.Form != "" {
		own["form"] = p.Form
	}
	p.ID += "-field"
	return p.merge(own)
}

func legendAttrs(ctx context.Context, p LegendProps) templ.Attributes {
	s, ok := fieldsetFrom(ctx)
	if !ok {
		return p.merge(nil)
	}
	p.ID = s.ID + "-legend"
	return p.merge(s.field().state(nil))
}
