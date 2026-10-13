package headlessui

import (
	"context"
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/radio.min.js
var radioJS string

var radioScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + radioJS + "</script>")))

func RadioScript() templ.Component {
	return radioScript.Once()
}

type RadioProps struct {
	Base
	Name     string
	Value    string
	Form     string
	Checked  bool
	Disabled bool
	Required bool
	Invalid  bool
}

type RadioIndicatorProps struct {
	Base
}

func radioAttrs(ctx context.Context, p RadioProps) checkable {
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
	c := checkableButton(p.Base, "radio", "huiRadio", f, inField, false, Events{OnClick: "select()", OnKeydown: "keydown($event)"})
	if f.Name != "" {
		c.attrs["data-hui-radio"] = f.Name
	}
	if s, ok := fieldsetFrom(ctx); ok && f.option && !f.checked && (s.Value != "" || len(s.Values) > 0) {
		c.attrs["tabindex"] = "-1"
	}
	return c
}

func radioIndicatorAttrs(ctx context.Context, p RadioIndicatorProps) templ.Attributes {
	return checkableIndicatorAttrs(ctx, p.Base, "checked", func(s checkableState) bool { return s.checked })
}
