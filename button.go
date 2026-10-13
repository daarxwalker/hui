package headlessui

import "github.com/a-h/templ"

type ButtonProps struct {
	Base
	Type      string
	Href      string
	Target    string
	Name      string
	Value     string
	Form      string
	Disabled  bool
	Pending   bool
	LeftSlot  templ.Component
	RightSlot templ.Component
}

func buttonAttrs(p ButtonProps) (templ.Attributes, bool) {
	if p.Pending {
		p.Disabled = true
	}

	own := make(templ.Attributes, 8)
	isLink := p.Href != ""

	switch {
	case isLink && p.Disabled:
		own["role"] = "link"
		own["aria-disabled"] = "true"
		own["tabindex"] = "-1"
	case isLink:
		own["href"] = string(templ.URL(p.Href))
		if p.Target != "" {
			own["target"] = p.Target
		}
		if p.Target == "_blank" {
			own["rel"] = "noopener noreferrer"
		}
	default:
		typ := p.Type
		if typ == "" {
			typ = "button"
		}
		own["type"] = typ
		if p.Name != "" {
			own["name"] = p.Name
		}
		if p.Value != "" {
			own["value"] = p.Value
		}
		if p.Form != "" {
			own["form"] = p.Form
		}
		if p.Disabled {
			own["disabled"] = true
		}
	}

	buttonState(own, p.Disabled)
	if p.Pending {
		own["aria-busy"] = "true"
		own["data-pending"] = true
	}

	return p.merge(own), isLink
}

func buttonState(own templ.Attributes, disabled bool) {
	if disabled {
		own["data-disabled"] = true
	} else {
		own["data-enabled"] = true
	}
}

func buttonSlotAttrs(side string) templ.Attributes {
	return templ.Attributes{"aria-hidden": "true", "data-slot": side}
}
