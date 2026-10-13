package headlessui

import (
	"context"
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/listbox.min.js
var listboxJS string

var listboxScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + listboxJS + "</script>")))

func ListboxScript() templ.Component {
	return listboxScript.Once()
}

const (
	ListboxShow   = "listboxShow()"
	ListboxHide   = "listboxHide()"
	ListboxToggle = "listboxToggle()"
)

type ListboxProps struct {
	Base
	Name        string
	Value       string
	Values      []string
	Multiple    bool
	Searchable  bool
	Form        string
	Placeholder string
	Disabled    bool
	Required    bool
	Invalid     bool
	Open        bool
	OnOpen      string
	OnClose     string
	OnChange    string
}

type ListboxTriggerProps struct {
	Base
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type ListboxValueProps struct {
	Base
}

type ListboxContentProps struct {
	Base
	Placement
}

type ListboxSearchProps struct {
	Base
	QueryName   string
	Value       string
	Placeholder string
}

type ListboxListProps struct {
	Base
}

type ListboxItemProps struct {
	Base
	Value     string
	Label     string
	Disabled  bool
	LeftSlot  templ.Component
	RightSlot templ.Component
}

type ListboxItemIndicatorProps struct {
	Base
}

type ListboxEmptyProps struct {
	Base
}

type ListboxLoaderProps struct {
	Base
}

type ListboxGroupProps struct {
	Base
}

type ListboxGroupLabelProps struct {
	Base
}

type ListboxSeparatorProps struct {
	Base
}

func WithListbox(ctx context.Context, p ListboxProps) context.Context {
	ctx, _ = listboxContext(ctx, p)
	return ctx
}

func listboxContext(ctx context.Context, p ListboxProps) (context.Context, listState) {
	return listContext(ctx, "listbox", listRoot{
		Base:        p.Base,
		Name:        p.Name,
		Value:       p.Value,
		Values:      p.Values,
		Multiple:    p.Multiple,
		Searchable:  p.Searchable,
		Form:        p.Form,
		Placeholder: p.Placeholder,
		Disabled:    p.Disabled,
		Required:    p.Required,
		Invalid:     p.Invalid,
		Open:        p.Open,
	})
}

func listboxAttrs(s listState, p ListboxProps) templ.Attributes {
	return listRootAttrs(p.Base, s, p.OnOpen, p.OnClose, p.OnChange)
}

func listboxTriggerAttrs(ctx context.Context, p ListboxTriggerProps) templ.Attributes {
	return listTriggerAttrs(ctx, p.Base, ListboxToggle)
}

func listboxValueAttrs(ctx context.Context, p ListboxValueProps) (templ.Attributes, string) {
	return listValueAttrs(ctx, p.Base)
}

func listboxContentAttrs(ctx context.Context, p ListboxContentProps) (templ.Attributes, string) {
	return listContentAttrs(ctx, p.Base, p.Placement)
}

func listboxSearchAttrs(ctx context.Context, p ListboxSearchProps) templ.Attributes {
	return listSearchAttrs(ctx, p.Base, p.QueryName, p.Value, p.Placeholder)
}

func listboxListAttrs(ctx context.Context, p ListboxListProps) templ.Attributes {
	return listOptionsAttrs(ctx, p.Base)
}

func listboxItemAttrs(ctx context.Context, p ListboxItemProps) templ.Attributes {
	return listItemAttrs(ctx, p.Base, p.Value, p.Label, p.Disabled)
}

func listboxGroupContext(ctx context.Context, p ListboxGroupProps) (context.Context, ListboxGroupProps) {
	ctx, p.ID = groupContext(ctx, p.ID, "listbox-group")
	return ctx, p
}
