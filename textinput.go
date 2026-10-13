package headlessui

import (
	"context"

	"github.com/a-h/templ"
)

const (
	textInputStateEmpty  = "{ focused: false, empty: true }"
	textInputStateFilled = "{ focused: false, empty: false }"
)

var textInputEvents = Events{
	OnFocus: "focused = true",
	OnBlur:  "focused = false",
	OnInput: "empty = $event.target.value === ''",
}

type TextInputProps struct {
	Base
	Type         string
	Placeholder  string
	Autocomplete string
	Inputmode    string
	Name         string
	Value        string
	Form         string
	Disabled     bool
	Readonly     bool
	Required     bool
	Invalid      bool
	Alpine       bool
}

func textInputAttrs(ctx context.Context, p TextInputProps) templ.Attributes {
	typ := p.Type
	if typ == "" {
		typ = "text"
	}
	own := make(templ.Attributes, 16)
	own["type"] = typ

	f, inField := controlFrom(ctx, FieldProps{
		Name:     p.Name,
		Value:    p.Value,
		Form:     p.Form,
		Disabled: p.Disabled,
		Readonly: p.Readonly,
		Required: p.Required,
		Invalid:  p.Invalid,
	})
	if inField {
		p.ID = f.ID
		if by := f.describedBy(); by != "" {
			own["aria-describedby"] = by
		}
	}
	p.Value = f.Value

	if f.Name != "" {
		own["name"] = f.Name
	}
	if f.Value != "" {
		own["value"] = f.Value
	}
	if f.Form != "" {
		own["form"] = f.Form
	}
	if p.Placeholder != "" {
		own["placeholder"] = p.Placeholder
	}
	if p.Autocomplete != "" {
		own["autocomplete"] = p.Autocomplete
	}
	if p.Inputmode != "" {
		own["inputmode"] = p.Inputmode
	}
	if f.Disabled {
		own["disabled"] = true
	}
	if f.Readonly {
		own["readonly"] = true
	}
	if f.Required {
		own["required"] = true
	}
	if f.Invalid {
		own["aria-invalid"] = "true"
	}

	if !p.Alpine {
		return p.merge(f.state(own))
	}

	if p.Value == "" {
		own["x-data"] = textInputStateEmpty
		own["data-empty"] = true
	} else {
		own["x-data"] = textInputStateFilled
	}
	own[":data-focused"] = "focused"
	own[":data-empty"] = "empty"

	return p.mergeEvents(f.state(own), textInputEvents)
}
