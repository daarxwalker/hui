package headlessui

import (
	"context"
	"strconv"

	"github.com/a-h/templ"
)

type TextareaProps struct {
	Base
	Placeholder  string
	Autocomplete string
	Name         string
	Value        string
	Form         string
	Rows         int
	Disabled     bool
	Readonly     bool
	Required     bool
	Invalid      bool
	Alpine       bool
}

func textareaAttrs(ctx context.Context, p TextareaProps) (templ.Attributes, string) {
	own := make(templ.Attributes, 16)

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

	if f.Name != "" {
		own["name"] = f.Name
	}
	if f.Form != "" {
		own["form"] = f.Form
	}
	if p.Rows > 0 {
		own["rows"] = strconv.Itoa(p.Rows)
	}
	if p.Placeholder != "" {
		own["placeholder"] = p.Placeholder
	}
	if p.Autocomplete != "" {
		own["autocomplete"] = p.Autocomplete
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
		return p.merge(f.state(own)), f.Value
	}

	if f.Value == "" {
		own["x-data"] = textInputStateEmpty
		own["data-empty"] = true
	} else {
		own["x-data"] = textInputStateFilled
	}
	own[":data-focused"] = "focused"
	own[":data-empty"] = "empty"

	return p.mergeEvents(f.state(own), textInputEvents), f.Value
}
