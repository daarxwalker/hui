package headlessui

import (
	"context"

	"github.com/a-h/templ"
)

type FieldProps struct {
	Base
	Name        string
	Value       string
	Form        string
	Description string
	Error       string
	Disabled    bool
	Readonly    bool
	Required    bool
	Invalid     bool
	option      bool
	checked     bool
}

type fieldKeyType struct{}

var fieldKey = fieldKeyType{}

func fieldFrom(ctx context.Context) (FieldProps, bool) {
	f, ok := ctx.Value(fieldKey).(FieldProps)
	return f, ok
}

func ownerFrom(ctx context.Context) (FieldProps, bool) {
	if f, ok := fieldFrom(ctx); ok {
		return f, true
	}
	if s, ok := fieldsetFrom(ctx); ok {
		return s.field(), true
	}
	return FieldProps{}, false
}

func controlFrom(ctx context.Context, own FieldProps) (FieldProps, bool) {
	if f, ok := fieldFrom(ctx); ok {
		if f.Form == "" {
			f.Form = own.Form
		}
		return f, true
	}
	if s, ok := fieldsetFrom(ctx); ok {
		own.Disabled = own.Disabled || s.Disabled
		if own.Form == "" {
			own.Form = s.Form
		}
	}
	return own, false
}

func checkableFrom(ctx context.Context, own FieldProps) (FieldProps, bool) {
	f, inField := controlFrom(ctx, own)
	if inField {
		if !f.option {
			f.checked = own.checked
		}
		return f, true
	}
	if s, ok := fieldsetFrom(ctx); ok && own.Name == "" {
		f = s.option(f)
	}
	return f, false
}

func (p FieldProps) normalize(ctx context.Context) FieldProps {
	if s, inFieldset := fieldsetFrom(ctx); inFieldset {
		p.Disabled = p.Disabled || s.Disabled
		if p.Form == "" {
			p.Form = s.Form
		}
		if p.Name == "" {
			p = s.option(p)
			if p.ID == "" && p.Value != "" {
				p.ID = s.ID + "-" + idSafe(p.Value)
			}
			if p.ID == "" {
				p.ID = nextID(ctx, s.ID)
			}
		}
	}
	if p.Error != "" {
		p.Invalid = true
	}
	switch {
	case p.ID != "":
	case p.Name != "":
		p.ID = p.Name
	default:
		p.ID = nextID(ctx, "field")
	}
	return p
}

func (p FieldProps) state(own templ.Attributes) templ.Attributes {
	if own == nil {
		own = make(templ.Attributes, 4)
	}
	if p.Disabled {
		own["data-disabled"] = true
	}
	if p.Invalid {
		own["data-invalid"] = true
	}
	if p.Readonly {
		own["data-readonly"] = true
	}
	if p.Required {
		own["data-required"] = true
	}
	return own
}

func (p FieldProps) descriptionID() string {
	if p.Description == "" {
		return ""
	}
	return p.ID + "-description"
}

func (p FieldProps) errorID() string {
	if p.Error == "" {
		return ""
	}
	return p.ID + "-error"
}

func (p FieldProps) describedBy() string {
	description, failure := p.descriptionID(), p.errorID()
	switch {
	case description != "" && failure != "":
		return description + " " + failure
	case description != "":
		return description
	default:
		return failure
	}
}

func fieldAttrs(p FieldProps) templ.Attributes {
	own := p.state(nil)
	switch {
	case p.option:
		own["data-option"] = p.Value
	case p.Name != "":
		own["data-field"] = p.Name
	}
	p.ID += "-field"
	return p.merge(own)
}

func validationAttrs(owner string, filled bool, form string, disabled bool) templ.Attributes {
	own := templ.Attributes{
		"required":         true,
		"tabindex":         "-1",
		"aria-hidden":      "true",
		"data-hui-require": owner,
		"value":            "",
		"style":            "position: absolute; width: 1px; height: 1px; padding: 0; border: 0; opacity: 0; pointer-events: none",
	}
	if filled {
		own["value"] = "1"
	}
	if form != "" {
		own["form"] = form
	}
	if disabled {
		own["disabled"] = true
	}
	return own
}
