package headlessui

import (
	"context"
	_ "embed"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/tooltip.min.js
var tooltipJS string

var tooltipScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + tooltipJS + "</script>")))

func TooltipScript() templ.Component {
	return tooltipScript.Once()
}

const (
	TooltipShow        = "tooltipShow($el)"
	TooltipHide        = "tooltipHide()"
	TooltipToggle      = "tooltipToggle($el, $event)"
	TooltipShowDelayed = "tooltipEnter($event, $el)"
	TooltipHideDelayed = "tooltipLeave($event)"
	TooltipFocus       = "tooltipFocus($event, $el)"
	TooltipBlur        = "tooltipBlur($event)"
	TooltipDismiss     = "tooltipDismiss()"
)

type TooltipProps struct {
	Base
	Open       bool
	Delay      int
	CloseDelay int
	OnOpen     string
	OnClose    string
}

type TooltipContentProps struct {
	Base
	Placement
}

type TooltipCloseProps struct {
	Base
}

type tooltipKeyType struct{}

var tooltipKey = tooltipKeyType{}

func tooltipFrom(ctx context.Context) (overlayState, bool) {
	s, ok := ctx.Value(tooltipKey).(overlayState)
	return s, ok
}

func tooltipContext(ctx context.Context, p TooltipProps) (context.Context, TooltipProps) {
	if p.ID == "" {
		p.ID = nextID(ctx, "tooltip")
	}
	return context.WithValue(ctx, tooltipKey, overlayState{id: p.ID, open: p.Open}), p
}

func tooltipData(p TooltipProps) string {
	config := make([]string, 0, 3)
	config = append(config, "open: "+strconv.FormatBool(p.Open))
	if p.Delay > 0 {
		config = append(config, "delay: "+strconv.Itoa(p.Delay))
	}
	if p.CloseDelay > 0 {
		config = append(config, "closeDelay: "+strconv.Itoa(p.CloseDelay))
	}
	return "huiTooltip({ " + strings.Join(config, ", ") + " })"
}

func tooltipAttrs(p TooltipProps) templ.Attributes {
	return overlayRootAttrs(p.Base, "tooltip", tooltipData(p), p.Open, p.OnOpen, p.OnClose)
}

func tooltipContentAttrs(ctx context.Context, p TooltipContentProps) (templ.Attributes, string) {
	own := make(templ.Attributes, 10)
	own["role"] = "tooltip"
	if p.Side == "" {
		p.Side = "top"
	}
	if p.Align == "" {
		p.Align = "center"
	}

	s, ok := tooltipFrom(ctx)
	if !ok {
		own["style"] = anchorPosition("", p.Placement)
		return p.merge(own), ""
	}

	p.ID = s.id + "-content"
	own["style"] = anchorPosition(anchorName(s.id), p.Placement)
	own["@pointerenter"] = chain(TooltipShowDelayed, attrString(p.Attrs["@pointerenter"]))
	own["@pointerleave"] = chain(TooltipHideDelayed, attrString(p.Attrs["@pointerleave"]))
	own["@focusout"] = chain(TooltipBlur, attrString(p.Attrs["@focusout"]))
	overlayVisibility(own, p.Attrs, s.open)
	return p.merge(own), overlayTeleport(false)
}

func tooltipCloseAttrs(ctx context.Context, p TooltipCloseProps) templ.Attributes {
	click := ""
	if _, ok := tooltipFrom(ctx); ok {
		click = TooltipDismiss
	}
	return p.mergeEvents(templ.Attributes{"type": "button", "data-enabled": true}, Events{OnClick: click})
}
