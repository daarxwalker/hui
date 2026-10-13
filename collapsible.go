package headlessui

import (
	"context"
	_ "embed"
	"strconv"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/collapsible.min.js
var collapsibleJS string

var collapsibleScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + collapsibleJS + "</script>")))

func CollapsibleScript() templ.Component {
	return collapsibleScript.Once()
}

const (
	CollapsibleShow   = "collapsibleShow()"
	CollapsibleHide   = "collapsibleHide()"
	CollapsibleToggle = "collapsibleToggle()"
)

type CollapsibleProps struct {
	Base
	Open    bool
	OnOpen  string
	OnClose string
}

type CollapsibleHeaderProps struct {
	Base
	Clickable bool
}

type CollapsibleTitleProps struct {
	Base
}

type CollapsibleTriggerProps struct {
	Base
	Disabled bool
}

type CollapsibleContentProps struct {
	Base
}

type collapsibleKeyType struct{}

var collapsibleKey = collapsibleKeyType{}

type collapsibleState struct {
	overlayState
	inHeader bool
}

func collapsibleFrom(ctx context.Context) (collapsibleState, bool) {
	s, ok := ctx.Value(collapsibleKey).(collapsibleState)
	return s, ok
}

func collapsibleContext(ctx context.Context, p CollapsibleProps) (context.Context, CollapsibleProps) {
	if p.ID == "" {
		p.ID = nextID(ctx, "collapsible")
	}
	s := collapsibleState{overlayState: overlayState{id: p.ID, open: p.Open}}
	return context.WithValue(ctx, collapsibleKey, s), p
}

func collapsibleHeaderContext(ctx context.Context) context.Context {
	s, ok := collapsibleFrom(ctx)
	if !ok {
		return ctx
	}
	s.inHeader = true
	return context.WithValue(ctx, collapsibleKey, s)
}

func collapsibleAttrs(p CollapsibleProps) templ.Attributes {
	data := "huiCollapsible({ open: " + strconv.FormatBool(p.Open) + " })"
	return overlayRootAttrs(p.Base, "collapsible", data, p.Open, p.OnOpen, p.OnClose)
}

func collapsibleHeaderAttrs(ctx context.Context, p CollapsibleHeaderProps) templ.Attributes {
	own := make(templ.Attributes, 4)
	click := ""
	if s, ok := collapsibleFrom(ctx); ok {
		own[":data-open"] = "open"
		if s.open {
			own["data-open"] = true
		}
		if p.Clickable {
			click = "collapsibleHeaderClick($event)"
		}
	}
	if p.Clickable {
		own["data-clickable"] = true
	}
	return p.mergeEvents(own, Events{OnClick: click})
}

func collapsibleTitleAttrs(ctx context.Context, p CollapsibleTitleProps) templ.Attributes {
	own := make(templ.Attributes, 3)
	if s, ok := collapsibleFrom(ctx); ok {
		p.ID = s.id + "-title"
		own[":data-open"] = "open"
		if s.open {
			own["data-open"] = true
		}
	}
	return p.merge(own)
}

func collapsibleTriggerAttrs(ctx context.Context, p CollapsibleTriggerProps) templ.Attributes {
	own := make(templ.Attributes, 9)
	own["type"] = "button"

	click := ""
	if s, ok := collapsibleFrom(ctx); ok {
		p.ID = s.id + "-trigger"
		overlayTriggerAttrs(own, s.overlayState)
		delete(own, "style")
		if _, named := p.Attrs["aria-label"]; s.inHeader && !named {
			own["aria-labelledby"] = s.id + "-title " + s.id + "-trigger"
		}
		click = CollapsibleToggle
	}
	if p.Disabled {
		own["disabled"] = true
	}
	buttonState(own, p.Disabled)
	return p.mergeEvents(own, Events{OnClick: click})
}

func collapsibleContentAttrs(ctx context.Context, p CollapsibleContentProps) templ.Attributes {
	own := make(templ.Attributes, 6)
	if s, ok := collapsibleFrom(ctx); ok {
		p.ID = s.id + "-content"
		transition(own, p.Attrs, "x-collapse")
		overlayVisibility(own, p.Attrs, s.open)
	}
	return p.merge(own)
}
