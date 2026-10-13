package headlessui

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

type listKeyType struct{}

var listKey = listKeyType{}

type listItemKeyType struct{}

var listItemKey = listItemKeyType{}

type listState struct {
	id          string
	kind        string
	field       FieldProps
	inField     bool
	values      []string
	multiple    bool
	searchable  bool
	placeholder string
	open        bool
}

type listRoot struct {
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
}

func listFrom(ctx context.Context) (listState, bool) {
	s, ok := ctx.Value(listKey).(listState)
	return s, ok
}

func listContext(ctx context.Context, kind string, p listRoot) (context.Context, listState) {
	id := p.ID
	if id == "" {
		id = nextID(ctx, kind)
	}

	f, inField := controlFrom(ctx, FieldProps{
		Name:     p.Name,
		Value:    p.Value,
		Form:     p.Form,
		Disabled: p.Disabled,
		Required: p.Required,
		Invalid:  p.Invalid,
	})
	f.Readonly = false

	s := listState{id: id, kind: kind, field: f, inField: inField, multiple: p.Multiple, searchable: p.Searchable, placeholder: p.Placeholder, open: p.Open}
	switch {
	case p.Multiple:
		s.values = slices.DeleteFunc(slices.Clone(p.Values), func(v string) bool { return v == "" })
	case f.Value != "":
		s.values = []string{f.Value}
	}
	return context.WithValue(ctx, listKey, s), s
}

func (s listState) contentID() string {
	return s.id + "-content"
}

func (s listState) listID() string {
	return s.id + "-list"
}

func (s listState) searchID() string {
	if s.kind == "combobox" {
		return s.controlID()
	}
	return s.id + "-search"
}

func (s listState) labelledBy() string {
	if s.inField {
		return s.field.ID + "-label"
	}
	return s.controlID()
}

func (s listState) controlID() string {
	switch {
	case s.inField:
		return s.field.ID
	case s.kind == "combobox":
		return s.id + "-input"
	default:
		return s.id + "-trigger"
	}
}

func (s listState) optionID(value string) string {
	return s.id + "-option-" + idSafe(value)
}

func (s listState) selected(value string) bool {
	return slices.Contains(s.values, value)
}

func (s listState) data() string {
	config := make([]string, 0, 7)
	config = append(config, "open: "+strconv.FormatBool(s.open))
	if s.multiple {
		values := make([]string, 0, len(s.values))
		for _, v := range s.values {
			values = append(values, jsString(v))
		}
		config = append(config, "multiple: true", "value: ["+strings.Join(values, ", ")+"]")
	} else {
		value := ""
		if len(s.values) > 0 {
			value = s.values[0]
		}
		config = append(config, "value: "+jsString(value))
	}
	if s.field.Name != "" {
		config = append(config, "name: "+jsString(s.field.Name))
	}
	if s.field.Form != "" {
		config = append(config, "form: "+jsString(s.field.Form))
	}
	if s.field.Disabled {
		config = append(config, "disabled: true")
	}
	if s.placeholder != "" {
		config = append(config, "placeholder: "+jsString(s.placeholder))
	}
	if s.searchable {
		config = append(config, "searchable: true")
	}
	component := map[string]string{"listbox": "huiListbox", "combobox": "huiCombobox"}[s.kind]
	return component + "({ " + strings.Join(config, ", ") + " })"
}

func listRootAttrs(b Base, s listState, onOpen, onClose, onChange string) templ.Attributes {
	b.ID = s.id
	own := s.field.state(nil)
	own["x-data"] = s.data()
	own["x-modelable"] = "value"
	own[":data-open"] = "open"
	if s.open {
		own["data-open"] = true
	}
	if s.kind == "combobox" {
		own["style"] = "anchor-name: " + anchorName(s.id)
	}

	attrs := b.merge(own)
	setHandler(attrs, "@hui-"+s.kind+":opened.self", onOpen)
	setHandler(attrs, "@hui-"+s.kind+":closed.self", onClose)
	setHandler(attrs, "@hui-"+s.kind+":change.self", onChange)
	return attrs
}

func listInputAttrs(s listState) []templ.Attributes {
	if s.field.Name == "" {
		return nil
	}
	inputs := make([]templ.Attributes, 0, len(s.values)+1)
	for _, v := range s.values {
		input := templ.Attributes{"type": "hidden", "name": s.field.Name, "value": v, "data-hui-input": s.id}
		if s.field.Form != "" {
			input["form"] = s.field.Form
		}
		if s.field.Disabled {
			input["disabled"] = true
		}
		inputs = append(inputs, input)
	}
	return inputs
}

func listRequiredAttrs(s listState) templ.Attributes {
	if !s.field.Required {
		return nil
	}
	own := validationAttrs(s.id, len(s.values) > 0, s.field.Form, s.field.Disabled)
	own["@focus"] = "focusControl()"
	return own
}

func listTriggerAttrs(ctx context.Context, b Base, toggle string) templ.Attributes {
	own := make(templ.Attributes, 20)
	own["type"] = "button"
	own["aria-haspopup"] = "listbox"

	s, ok := listFrom(ctx)
	if !ok {
		return b.merge(own)
	}

	b.ID = s.controlID()
	own["style"] = "anchor-name: " + anchorName(s.id)
	if s.searchable {
		own["aria-controls"] = s.listID()
		own["aria-expanded"] = strconv.FormatBool(s.open)
		own[":aria-expanded"] = "open"
		listTriggerState(s, own)
	} else {
		listControlAttrs(s, own)
	}
	return b.mergeEvents(s.field.state(own), Events{OnClick: toggle, OnKeydown: "keydown($event)"})
}

func listValueAttrs(ctx context.Context, b Base) (templ.Attributes, string) {
	s, ok := listFrom(ctx)
	if !ok {
		return b.merge(nil), ""
	}
	own := templ.Attributes{"data-hui-value": s.id, ":data-placeholder": "!selected().length"}
	if len(s.values) > 0 {
		return b.merge(own), ""
	}
	own["data-placeholder"] = true
	return b.merge(own), s.placeholder
}

func listContentAttrs(ctx context.Context, b Base, placement Placement) (templ.Attributes, string) {
	own := make(templ.Attributes, 8)
	s, ok := listFrom(ctx)
	if !ok {
		own["style"] = anchorPosition("", placement)
		return b.merge(own), ""
	}
	b.ID = s.contentID()
	listPopupAttrs(s, own, b, placement)
	return b.merge(own), overlayTeleport(false)
}

func listSearchAttrs(ctx context.Context, b Base, queryName, value, placeholder string) templ.Attributes {
	own := make(templ.Attributes, 16)
	own["type"] = "text"
	own["autocomplete"] = "off"
	own["role"] = "combobox"
	own["aria-autocomplete"] = "list"
	if queryName != "" {
		own["name"] = queryName
	}
	if value != "" {
		own["value"] = value
	}
	if placeholder != "" {
		own["placeholder"] = placeholder
	}

	s, ok := listFrom(ctx)
	if !ok {
		return b.merge(own)
	}

	b.ID = s.searchID()
	own["aria-controls"] = s.listID()
	own["aria-expanded"] = strconv.FormatBool(s.open)
	own[":aria-expanded"] = "open"
	own[":aria-activedescendant"] = "active"
	if _, labelled := b.Attrs["aria-label"]; !labelled {
		own["aria-labelledby"] = s.labelledBy()
	}
	return b.mergeEvents(own, Events{OnInput: "input()", OnKeydown: "searchKeydown($event)"})
}

func listOptionsAttrs(ctx context.Context, b Base) templ.Attributes {
	own := make(templ.Attributes, 8)
	own["role"] = "listbox"
	own["tabindex"] = "-1"

	s, ok := listFrom(ctx)
	if !ok {
		return b.merge(own)
	}
	b.ID = s.listID()
	listListAttrs(s, own, b)
	return b.mergeEvents(own, Events{OnClick: "pick($event)"})
}

func listControlAttrs(s listState, own templ.Attributes) {
	own["role"] = "combobox"
	own["aria-controls"] = s.listID()
	own["aria-expanded"] = strconv.FormatBool(s.open)
	own[":aria-expanded"] = "open"
	own[":aria-activedescendant"] = "active"
	listTriggerState(s, own)
	if s.field.Required {
		own["aria-required"] = "true"
	}
}

func listTriggerState(s listState, own templ.Attributes) {
	own[":data-open"] = "open"
	if s.open {
		own["data-open"] = true
	}
	if len(s.values) == 0 {
		own["data-placeholder"] = true
	}
	own[":data-placeholder"] = "!selected().length"
	if by := s.field.describedBy(); by != "" {
		own["aria-describedby"] = by
	}
	if s.field.Disabled {
		own["disabled"] = true
	}
	if s.field.Invalid {
		own["aria-invalid"] = "true"
	}
}

func listPopupAttrs(s listState, own templ.Attributes, b Base, placement Placement) {
	own["style"] = anchorPosition(anchorName(s.id), placement) + "; min-width: anchor-size(width)"
	overlayVisibility(own, b.Attrs, s.open)
}

func listListAttrs(s listState, own templ.Attributes, b Base) {
	own["aria-labelledby"] = s.labelledBy()
	if s.multiple {
		own["aria-multiselectable"] = "true"
	}
	own["@mousedown.prevent"] = true
	own["@pointermove"] = chain("hover($event)", attrString(b.Attrs["@pointermove"]))
}

func listLoaderAttrs(b Base) templ.Attributes {
	return b.merge(templ.Attributes{"aria-hidden": "true", "data-hui-loader": true})
}

func listEmptyAttrs(b Base) templ.Attributes {
	return b.merge(templ.Attributes{"role": "presentation"})
}

func listItemContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, listItemKey, value)
}

func listItemAttrs(ctx context.Context, b Base, value, label string, disabled bool) templ.Attributes {
	own := make(templ.Attributes, 12)
	own["role"] = "option"
	own["data-value"] = value
	if label != "" {
		own["data-label"] = label
	}
	if disabled {
		own["aria-disabled"] = "true"
		own["data-disabled"] = true
	}

	s, ok := listFrom(ctx)
	if !ok {
		own["aria-selected"] = "false"
		return b.merge(own)
	}

	b.ID = s.optionID(value)
	own["aria-selected"] = strconv.FormatBool(s.selected(value))
	if s.selected(value) {
		own["data-selected"] = true
	}
	return b.merge(own)
}

func listIndicatorAttrs(ctx context.Context, b Base) templ.Attributes {
	own := templ.Attributes{"aria-hidden": "true", "data-hui-indicator": true}
	s, inList := listFrom(ctx)
	value, inItem := ctx.Value(listItemKey).(string)
	if inList && inItem && !s.selected(value) {
		own["style"] = "display: none"
	}
	return b.merge(own)
}

func listGroupAttrs(b Base) templ.Attributes {
	return b.merge(templ.Attributes{"role": "group", "aria-labelledby": b.ID + "-label"})
}

func listGroupLabelAttrs(ctx context.Context, b Base) templ.Attributes {
	if id, ok := groupFrom(ctx); ok {
		b.ID = id + "-label"
	}
	return b.merge(nil)
}

func listSeparatorAttrs(b Base) templ.Attributes {
	return b.merge(templ.Attributes{"aria-hidden": "true"})
}
