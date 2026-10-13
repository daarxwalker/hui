package headlessui

import (
	"context"

	"github.com/a-h/templ"
)

type SelectProps struct {
	Base
	Name         string
	Value        string
	Form         string
	Placeholder  string
	Autocomplete string
	Disabled     bool
	Required     bool
	Invalid      bool
}

type SelectTriggerProps struct {
	Base
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type SelectValueProps struct {
	Base
}

type SelectItemProps struct {
	Base
	Value     string
	Disabled  bool
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type SelectGroupProps struct {
	Base
	Label    string
	Disabled bool
}

type SelectSeparatorProps struct {
	Base
}

type selectState struct {
	value       string
	placeholder string
	disabled    bool
}

type selectKeyType struct{}

var selectKey = selectKeyType{}

func selectFrom(ctx context.Context) (selectState, bool) {
	s, ok := ctx.Value(selectKey).(selectState)
	return s, ok
}

func selectContext(ctx context.Context, p SelectProps) (context.Context, templ.Attributes) {
	own := make(templ.Attributes, 12)
	f, inField := controlFrom(ctx, FieldProps{
		Name:     p.Name,
		Value:    p.Value,
		Form:     p.Form,
		Disabled: p.Disabled,
		Required: p.Required,
		Invalid:  p.Invalid,
	})
	f.Readonly = false
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
	if p.Autocomplete != "" {
		own["autocomplete"] = p.Autocomplete
	}
	if f.Disabled {
		own["disabled"] = true
	}
	if f.Required {
		own["required"] = true
	}
	if f.Invalid {
		own["aria-invalid"] = "true"
	}

	s := selectState{value: f.Value, placeholder: p.Placeholder, disabled: f.Disabled}
	return context.WithValue(ctx, selectKey, s), p.merge(f.state(own))
}

func selectTriggerAttrs(ctx context.Context, p SelectTriggerProps) (templ.Attributes, templ.Attributes, string) {
	own := make(templ.Attributes, 4)
	s, ok := selectFrom(ctx)
	buttonState(own, s.disabled)
	if !ok || s.placeholder == "" {
		return p.merge(own), nil, ""
	}
	placeholder := templ.Attributes{"value": "", "data-placeholder": true}
	if s.value == "" {
		placeholder["selected"] = true
	}
	return p.merge(own), placeholder, s.placeholder
}

func selectValueAttrs(p SelectValueProps) templ.Attributes {
	return p.merge(nil)
}

func selectItemAttrs(ctx context.Context, p SelectItemProps) templ.Attributes {
	own := templ.Attributes{"value": p.Value}
	if s, ok := selectFrom(ctx); ok && s.value != "" && s.value == p.Value {
		own["selected"] = true
	}
	if p.Disabled {
		own["disabled"] = true
		own["data-disabled"] = true
	}
	return p.merge(own)
}

func selectGroupAttrs(p SelectGroupProps) templ.Attributes {
	own := make(templ.Attributes, 3)
	if p.Label != "" {
		own["label"] = p.Label
	}
	if p.Disabled {
		own["disabled"] = true
		own["data-disabled"] = true
	}
	return p.merge(own)
}

func selectSeparatorAttrs(p SelectSeparatorProps) templ.Attributes {
	return p.merge(nil)
}
