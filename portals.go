package headlessui

import (
	"context"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

const portalsID = "hui-portals"

type PortalsProps struct {
	Base
}

func portalsAttrs(p PortalsProps) templ.Attributes {
	p.ID = portalsID
	return p.merge(nil)
}

func animated(attrs templ.Attributes) bool {
	for name := range attrs {
		if strings.HasPrefix(name, "x-transition") || strings.HasPrefix(name, "x-collapse") {
			return true
		}
	}
	return false
}

func transition(own, attrs templ.Attributes, directive string) {
	if animated(own) || animated(attrs) {
		return
	}
	own[directive] = true
}

func overlayVisibility(own, attrs templ.Attributes, open bool) {
	own["x-show"] = "open"
	transition(own, attrs, "x-transition")
	own[":data-open"] = "open"
	if open {
		own["data-open"] = true
	} else {
		own["style"] = joinStyle(attrString(own["style"]), "display: none")
	}
}

func overlayTeleport(inOverlay bool) string {
	if inOverlay {
		return ""
	}
	return "#" + portalsID
}

type groupKeyType struct{}

var groupKey = groupKeyType{}

func groupContext(ctx context.Context, id, prefix string) (context.Context, string) {
	if id == "" {
		id = nextID(ctx, prefix)
	}
	return context.WithValue(ctx, groupKey, id), id
}

func groupFrom(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(groupKey).(string)
	return id, ok
}

type overlayState struct {
	id   string
	open bool
}

func overlayRootAttrs(b Base, component, data string, open bool, onOpen, onClose string) templ.Attributes {
	own := make(templ.Attributes, 6)
	own["x-data"] = data
	own["x-modelable"] = "open"
	own[":data-open"] = "open"
	if open {
		own["data-open"] = true
	}
	attrs := b.merge(own)
	setHandler(attrs, "@hui-"+component+":opened.self", onOpen)
	setHandler(attrs, "@hui-"+component+":closed.self", onClose)
	return attrs
}

func overlayTriggerAttrs(own templ.Attributes, s overlayState) {
	own["aria-controls"] = s.id + "-content"
	own["aria-expanded"] = strconv.FormatBool(s.open)
	own[":aria-expanded"] = "open"
	own[":data-open"] = "open"
	if s.open {
		own["data-open"] = true
	}
	own["style"] = "anchor-name: " + anchorName(s.id)
}
