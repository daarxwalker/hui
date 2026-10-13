package headlessui

import (
	"context"
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/checkbox.min.js
var checkboxJS string

var checkboxScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + checkboxJS + "</script>")))

func CheckboxScript() templ.Component {
	return checkboxScript.Once()
}

type CheckboxProps struct {
	Base
	Name          string
	Value         string
	Form          string
	Checked       bool
	Indeterminate bool
	Disabled      bool
	Required      bool
	Invalid       bool
}

type CheckboxIndicatorProps struct {
	Base
	Indeterminate bool
}

func checkboxAttrs(ctx context.Context, p CheckboxProps) checkable {
	f, inField := checkableFrom(ctx, FieldProps{
		Name:     p.Name,
		Value:    p.Value,
		Form:     p.Form,
		Disabled: p.Disabled,
		Required: p.Required,
		Invalid:  p.Invalid,
		checked:  p.Checked,
	})
	f.Readonly = false
	return checkableButton(p.Base, "checkbox", "huiCheckbox", f, inField, p.Indeterminate, Events{OnClick: "toggle()", OnKeydown: "keydown($event)"})
}

func checkboxIndicatorAttrs(ctx context.Context, p CheckboxIndicatorProps) templ.Attributes {
	if p.Indeterminate {
		return checkableIndicatorAttrs(ctx, p.Base, "indeterminate", func(s checkableState) bool { return s.indeterminate })
	}
	return checkableIndicatorAttrs(ctx, p.Base, "checked && !indeterminate", func(s checkableState) bool { return s.checked && !s.indeterminate })
}
