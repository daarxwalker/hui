package headlessui

import (
	"context"
	_ "embed"
	"strconv"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/combobox.min.js
var comboboxJS string

var comboboxScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + comboboxJS + "</script>")))

func ComboboxScript() templ.Component {
	return comboboxScript.Once()
}

const (
	ComboboxShow   = "comboboxShow()"
	ComboboxHide   = "comboboxHide()"
	ComboboxToggle = "comboboxToggle()"
)

type ComboboxProps struct {
	Base
	Name     string
	Value    string
	Values   []string
	Multiple bool
	Form     string
	Disabled bool
	Required bool
	Invalid  bool
	Open     bool
	OnOpen   string
	OnClose  string
	OnChange string
}

type ComboboxTriggerProps struct {
	Base
}

type ComboboxContentProps struct {
	Base
	Placement
}

type ComboboxInputProps struct {
	Base
	QueryName   string
	Value       string
	Placeholder string
}

type ComboboxListProps struct {
	Base
}

type ComboboxItemProps struct {
	Base
	Value     string
	Label     string
	Disabled  bool
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type ComboboxItemIndicatorProps struct {
	Base
}

type ComboboxEmptyProps struct {
	Base
}

type ComboboxLoaderProps struct {
	Base
}

type ComboboxRemoveProps struct {
	Base
	Value string
}

type ComboboxGroupProps struct {
	Base
}

type ComboboxGroupLabelProps struct {
	Base
}

type ComboboxSeparatorProps struct {
	Base
}

func WithCombobox(ctx context.Context, p ComboboxProps) context.Context {
	ctx, _ = comboboxContext(ctx, p)
	return ctx
}

func comboboxContext(ctx context.Context, p ComboboxProps) (context.Context, listState) {
	return listContext(ctx, "combobox", listRoot{
		Base:       p.Base,
		Name:       p.Name,
		Value:      p.Value,
		Values:     p.Values,
		Multiple:   p.Multiple,
		Searchable: true,
		Form:       p.Form,
		Disabled:   p.Disabled,
		Required:   p.Required,
		Invalid:    p.Invalid,
		Open:       p.Open,
	})
}

func comboboxAttrs(s listState, p ComboboxProps) templ.Attributes {
	return listRootAttrs(p.Base, s, p.OnOpen, p.OnClose, p.OnChange)
}

func comboboxTriggerAttrs(ctx context.Context, p ComboboxTriggerProps) templ.Attributes {
	own := make(templ.Attributes, 16)
	own["type"] = "button"
	own["tabindex"] = "-1"
	own["aria-haspopup"] = "listbox"

	s, ok := listFrom(ctx)
	if !ok {
		buttonState(own, false)
		return p.merge(own)
	}

	p.ID = s.id + "-trigger"
	own["aria-controls"] = s.listID()
	own["aria-expanded"] = strconv.FormatBool(s.open)
	own[":aria-expanded"] = "open"
	own[":data-open"] = "open"
	if s.open {
		own["data-open"] = true
	}
	if s.inField {
		own["aria-labelledby"] = s.labelledBy()
	}
	if s.field.Disabled {
		own["disabled"] = true
	}
	own["@mousedown.prevent"] = true
	buttonState(own, s.field.Disabled)
	return p.mergeEvents(s.field.state(own), Events{OnClick: ComboboxToggle})
}

func comboboxContentAttrs(ctx context.Context, p ComboboxContentProps) (templ.Attributes, string) {
	return listContentAttrs(ctx, p.Base, p.Placement)
}

func comboboxInputAttrs(ctx context.Context, p ComboboxInputProps) templ.Attributes {
	own := make(templ.Attributes, 20)
	own["type"] = "text"
	own["autocomplete"] = "off"
	own["role"] = "combobox"
	own["aria-autocomplete"] = "list"
	if p.QueryName != "" {
		own["data-hui-query"] = p.QueryName
	}
	if p.Value != "" {
		own["value"] = p.Value
	}
	if p.Placeholder != "" {
		own["placeholder"] = p.Placeholder
	}

	s, ok := listFrom(ctx)
	if !ok {
		return p.merge(own)
	}

	p.ID = s.controlID()
	own["aria-controls"] = s.listID()
	own["aria-expanded"] = strconv.FormatBool(s.open)
	own[":aria-expanded"] = "open"
	own[":aria-activedescendant"] = "active"
	own[":data-open"] = "open"
	if s.open {
		own["data-open"] = true
	}
	if by := s.field.describedBy(); by != "" {
		own["aria-describedby"] = by
	}
	if s.field.Disabled {
		own["disabled"] = true
	}
	if s.field.Invalid {
		own["aria-invalid"] = "true"
	}
	if s.field.Required {
		own["aria-required"] = "true"
	}
	return p.mergeEvents(s.field.state(own), Events{OnInput: "input()", OnKeydown: "keydown($event)", OnBlur: "inputBlur($event)", OnClick: "inputClick()"})
}

func comboboxListAttrs(ctx context.Context, p ComboboxListProps) templ.Attributes {
	return listOptionsAttrs(ctx, p.Base)
}

func comboboxItemAttrs(ctx context.Context, p ComboboxItemProps) templ.Attributes {
	return listItemAttrs(ctx, p.Base, p.Value, p.Label, p.Disabled)
}

func comboboxRemoveAttrs(ctx context.Context, p ComboboxRemoveProps) templ.Attributes {
	own := templ.Attributes{"type": "button", "tabindex": "-1"}
	click := ""
	if _, ok := listFrom(ctx); ok {
		click = "remove(" + jsString(p.Value) + ")"
	}
	return p.mergeEvents(own, Events{OnClick: click})
}

func comboboxGroupContext(ctx context.Context, p ComboboxGroupProps) (context.Context, ComboboxGroupProps) {
	ctx, p.ID = groupContext(ctx, p.ID, "combobox-group")
	return ctx, p
}
