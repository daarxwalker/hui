package headlessui

import (
	"context"
	_ "embed"
	"strconv"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/dropdown.min.js
var dropdownJS string

var dropdownScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + dropdownJS + "</script>")))

func DropdownScript() templ.Component {
	return dropdownScript.Once()
}

const (
	DropdownShow    = "dropdownShow($el, $event)"
	DropdownHide    = "dropdownHide()"
	DropdownToggle  = "dropdownToggle($el, $event)"
	DropdownKeydown = "dropdownKeydown($event, $el)"
)

type DropdownProps struct {
	Base
	Open    bool
	OnOpen  string
	OnClose string
}

type DropdownTriggerProps struct {
	Base
	Disabled bool
}

type DropdownContentProps struct {
	Base
	Placement
}

type DropdownItemProps struct {
	Base
	Type      string
	Href      string
	Target    string
	Name      string
	Value     string
	Form      string
	Disabled  bool
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type DropdownCheckboxItemProps struct {
	Base
	Name      string
	Value     string
	Form      string
	Checked   bool
	Disabled  bool
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type DropdownRadioGroupProps struct {
	Base
	Name     string
	Value    string
	Form     string
	Disabled bool
}

type DropdownRadioItemProps struct {
	Base
	Value     string
	Checked   bool
	Disabled  bool
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type DropdownGroupProps struct {
	Base
}

type DropdownLabelProps struct {
	Base
}

type DropdownSeparatorProps struct {
	Base
}

type dropdownKeyType struct{}

var dropdownKey = dropdownKeyType{}

type dropdownRadioKeyType struct{}

var dropdownRadioKey = dropdownRadioKeyType{}

type dropdownState struct {
	id   string
	open bool
}

func dropdownFrom(ctx context.Context) (dropdownState, bool) {
	state, ok := ctx.Value(dropdownKey).(dropdownState)
	return state, ok
}

func dropdownContext(ctx context.Context, p DropdownProps) (context.Context, DropdownProps) {
	if p.ID == "" {
		p.ID = nextID(ctx, "dropdown")
	}
	return context.WithValue(ctx, dropdownKey, dropdownState{id: p.ID, open: p.Open}), p
}

func dropdownAttrs(p DropdownProps) templ.Attributes {
	own := make(templ.Attributes, 6)
	own["x-data"] = "huiDropdown({ open: " + strconv.FormatBool(p.Open) + " })"
	own["x-modelable"] = "open"
	own[":data-open"] = "open"
	if p.Open {
		own["data-open"] = true
	}

	attrs := p.merge(own)
	setHandler(attrs, "@hui-dropdown:opened.self", p.OnOpen)
	setHandler(attrs, "@hui-dropdown:closed.self", p.OnClose)
	return attrs
}

func dropdownTriggerAttrs(ctx context.Context, p DropdownTriggerProps) templ.Attributes {
	own := make(templ.Attributes, 10)
	own["type"] = "button"
	own["aria-haspopup"] = "menu"

	internal := Events{}
	if d, ok := dropdownFrom(ctx); ok {
		p.ID = d.id + "-trigger"
		own["aria-controls"] = d.id + "-content"
		own["aria-expanded"] = strconv.FormatBool(d.open)
		own[":aria-expanded"] = "open"
		own[":data-open"] = "open"
		if d.open {
			own["data-open"] = true
		}
		own["style"] = "anchor-name: " + anchorName(d.id)
		internal = Events{OnClick: DropdownToggle, OnKeydown: DropdownKeydown}
	}
	if p.Disabled {
		own["disabled"] = true
	}
	buttonState(own, p.Disabled)

	return p.mergeEvents(own, internal)
}

func dropdownContentAttrs(ctx context.Context, p DropdownContentProps) (templ.Attributes, string) {
	own := make(templ.Attributes, 12)
	own["role"] = "menu"
	own["tabindex"] = "-1"

	d, inDropdown := dropdownFrom(ctx)
	if !inDropdown {
		own["style"] = anchorPosition("", p.Placement)
		return p.merge(own), ""
	}

	p.ID = d.id + "-content"
	own["style"] = anchorPosition(anchorName(d.id), p.Placement)
	own["aria-labelledby"] = d.id + "-trigger"
	own["@pointermove"] = chain("hover($event)", attrString(p.Attrs["@pointermove"]))
	overlayVisibility(own, p.Attrs, d.open)

	return p.mergeEvents(own, Events{OnKeydown: "keydown($event)", OnClick: "select($event)"}), overlayTeleport(false)
}

func dropdownItemAttrs(p DropdownItemProps) (templ.Attributes, bool) {
	attrs, isLink := buttonAttrs(ButtonProps{
		Base:     p.Base,
		Type:     p.Type,
		Href:     p.Href,
		Target:   p.Target,
		Name:     p.Name,
		Value:    p.Value,
		Form:     p.Form,
		Disabled: p.Disabled,
	})
	attrs["role"] = "menuitem"
	attrs["tabindex"] = "-1"
	return attrs, isLink
}

func dropdownCheckable(b Base, role, component string, f FieldProps, onClick string) checkable {
	c := checkableButton(b, role, component, f, false, false, Events{OnClick: onClick})
	c.attrs["tabindex"] = "-1"
	buttonState(c.attrs, f.Disabled)
	return c
}

func dropdownCheckboxItemAttrs(p DropdownCheckboxItemProps) checkable {
	f := FieldProps{Name: p.Name, Value: p.Value, Form: p.Form, Disabled: p.Disabled, checked: p.Checked}
	return dropdownCheckable(p.Base, "menuitemcheckbox", "huiDropdownCheckbox", f, "dropdownCheck()")
}

func dropdownRadioGroupContext(ctx context.Context, p DropdownRadioGroupProps) (context.Context, DropdownRadioGroupProps) {
	ctx, p.ID = groupContext(ctx, p.ID, "dropdown-group")
	return context.WithValue(ctx, dropdownRadioKey, p), p
}

func dropdownRadioGroupAttrs(p DropdownRadioGroupProps) templ.Attributes {
	own := templ.Attributes{"role": "group", "aria-labelledby": p.ID + "-label"}
	if p.Disabled {
		own["data-disabled"] = true
	}
	return p.merge(own)
}

func dropdownRadioItemAttrs(ctx context.Context, p DropdownRadioItemProps) checkable {
	f := FieldProps{Value: p.Value, Disabled: p.Disabled, checked: p.Checked}
	if g, ok := ctx.Value(dropdownRadioKey).(DropdownRadioGroupProps); ok {
		f.Name = g.Name
		f.Form = g.Form
		f.Disabled = f.Disabled || g.Disabled
		f.checked = g.Value != "" && g.Value == p.Value
	}
	return dropdownCheckable(p.Base, "menuitemradio", "huiDropdownRadio", f, "dropdownPick()")
}

func dropdownGroupContext(ctx context.Context, p DropdownGroupProps) (context.Context, DropdownGroupProps) {
	ctx, p.ID = groupContext(ctx, p.ID, "dropdown-group")
	return ctx, p
}

func dropdownGroupAttrs(p DropdownGroupProps) templ.Attributes {
	return p.merge(templ.Attributes{"role": "group", "aria-labelledby": p.ID + "-label"})
}

func dropdownLabelAttrs(ctx context.Context, p DropdownLabelProps) templ.Attributes {
	if id, ok := groupFrom(ctx); ok {
		p.ID = id + "-label"
	}
	return p.merge(nil)
}

func dropdownSeparatorAttrs(p DropdownSeparatorProps) templ.Attributes {
	return p.merge(templ.Attributes{"role": "separator"})
}
