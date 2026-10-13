package headlessui

import (
	"context"
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/switch.min.js
var switchJS string

var switchScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + switchJS + "</script>")))

func SwitchScript() templ.Component {
	return switchScript.Once()
}

type SwitchProps struct {
	Base
	Name     string
	Value    string
	Form     string
	Checked  bool
	Disabled bool
	Required bool
	Invalid  bool
}

type SwitchThumbProps struct {
	Base
}

func switchAttrs(ctx context.Context, p SwitchProps) checkable {
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
	return checkableButton(p.Base, "switch", "huiSwitch", f, inField, false, Events{OnClick: "toggle()"})
}

func switchThumbAttrs(ctx context.Context, p SwitchThumbProps) templ.Attributes {
	own := templ.Attributes{"aria-hidden": "true", ":data-checked": "checked"}
	if s, ok := checkableFromCtx(ctx); ok && s.checked {
		own["data-checked"] = true
	}
	return p.merge(own)
}
