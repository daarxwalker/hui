package headlessui

import (
	"context"
	_ "embed"
	"strings"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/tabs.min.js
var tabsJS string

var tabsScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + tabsJS + "</script>")))

func TabsScript() templ.Component {
	return tabsScript.Once()
}

type TabsProps struct {
	Base
	Value       string
	Orientation string
	Activation  string
	OnChange    string
}

type TabsListProps struct {
	Base
}

type TabsTriggerProps struct {
	Base
	Value    string
	Disabled bool
}

type TabsContentProps struct {
	Base
	Value string
}

type tabsKeyType struct{}

var tabsKey = tabsKeyType{}

type tabsState struct {
	id          string
	value       string
	orientation string
}

func tabsFrom(ctx context.Context) (tabsState, bool) {
	s, ok := ctx.Value(tabsKey).(tabsState)
	return s, ok
}

func tabsOrientation(orientation string) string {
	if orientation == "vertical" {
		return orientation
	}
	return "horizontal"
}

func tabsContext(ctx context.Context, p TabsProps) (context.Context, TabsProps) {
	if p.ID == "" {
		p.ID = nextID(ctx, "tabs")
	}
	s := tabsState{id: p.ID, value: p.Value, orientation: tabsOrientation(p.Orientation)}
	return context.WithValue(ctx, tabsKey, s), p
}

func (s tabsState) triggerID(value string) string {
	return s.id + "-trigger-" + idSafe(value)
}

func (s tabsState) contentID(value string) string {
	return s.id + "-content-" + idSafe(value)
}

func tabsAttrs(p TabsProps) templ.Attributes {
	config := []string{"value: " + jsString(p.Value)}
	if p.Activation == "manual" {
		config = append(config, "manual: true")
	}
	own := make(templ.Attributes, 4)
	own["x-data"] = "huiTabs({ " + strings.Join(config, ", ") + " })"
	own["x-modelable"] = "value"
	own["data-orientation"] = tabsOrientation(p.Orientation)

	attrs := p.merge(own)
	setHandler(attrs, "@hui-tabs:change.self", p.OnChange)
	return attrs
}

func tabsListAttrs(ctx context.Context, p TabsListProps) templ.Attributes {
	own := make(templ.Attributes, 4)
	own["role"] = "tablist"
	keydown := ""
	if s, ok := tabsFrom(ctx); ok {
		own["aria-orientation"] = s.orientation
		own["data-orientation"] = s.orientation
		keydown = "keydown($event)"
	}
	return p.mergeEvents(own, Events{OnKeydown: keydown})
}

func tabsTriggerAttrs(ctx context.Context, p TabsTriggerProps) templ.Attributes {
	own := make(templ.Attributes, 14)
	own["type"] = "button"
	own["role"] = "tab"
	if p.Disabled {
		own["disabled"] = true
	}
	buttonState(own, p.Disabled)

	s, ok := tabsFrom(ctx)
	if !ok {
		return p.merge(own)
	}

	p.ID = s.triggerID(p.Value)
	selected := s.value == p.Value
	test := "value === " + jsString(p.Value)
	own["data-value"] = p.Value
	own["aria-controls"] = s.contentID(p.Value)
	own["aria-selected"] = "false"
	own["tabindex"] = "-1"
	if selected {
		own["aria-selected"] = "true"
		own["tabindex"] = "0"
		own["data-selected"] = true
	}
	own[":aria-selected"] = test
	own[":tabindex"] = test + " ? 0 : -1"
	own[":data-selected"] = test
	return p.mergeEvents(own, Events{OnClick: "select(" + jsString(p.Value) + ")"})
}

func tabsContentAttrs(ctx context.Context, p TabsContentProps) templ.Attributes {
	own := make(templ.Attributes, 10)
	own["role"] = "tabpanel"
	own["tabindex"] = "0"

	s, ok := tabsFrom(ctx)
	if !ok {
		return p.merge(own)
	}

	p.ID = s.contentID(p.Value)
	test := "value === " + jsString(p.Value)
	own["aria-labelledby"] = s.triggerID(p.Value)
	own["x-show"] = test
	own[":data-selected"] = test
	if s.value == p.Value {
		own["data-selected"] = true
	} else {
		own["style"] = "display: none"
	}
	return p.merge(own)
}
