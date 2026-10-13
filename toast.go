package headlessui

import (
	"context"
	_ "embed"
	"strconv"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/toast.min.js
var toastJS string

var toastScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + toastJS + "</script>")))

func ToastScript() templ.Component {
	return toastScript.Once()
}

const toasterID = "hui-toaster"

type ToasterProps struct {
	Base
	Label string
}

type ToastProps struct {
	Base
	Duration int
	Urgent   bool
	OnClose  string
}

type ToastTitleProps struct {
	Base
}

type ToastDescriptionProps struct {
	Base
}

type ToastActionProps struct {
	Base
}

type ToastCloseProps struct {
	Base
}

type toastKeyType struct{}

var toastKey = toastKeyType{}

func toasterAttrs(p ToasterProps) templ.Attributes {
	if p.ID == "" {
		p.ID = toasterID
	}
	label := p.Label
	if label == "" {
		label = "Notifications"
	}
	own := make(templ.Attributes, 4)
	own["aria-label"] = label
	own["aria-live"] = "polite"
	own["data-hui-toaster"] = true
	return p.merge(own)
}

func toastContext(ctx context.Context, p ToastProps) (context.Context, ToastProps) {
	if p.ID == "" {
		p.ID = nextID(ctx, "toast")
	}
	return context.WithValue(ctx, toastKey, p.ID), p
}

func toastAttrs(p ToastProps) templ.Attributes {
	duration := p.Duration
	if duration == 0 {
		duration = 5000
	}
	if duration < 0 {
		duration = 0
	}

	own := make(templ.Attributes, 14)
	own["role"] = "status"
	if p.Urgent {
		own["role"] = "alert"
	}
	own["aria-atomic"] = "true"
	own["aria-labelledby"] = p.ID + "-title"
	own["aria-describedby"] = p.ID + "-description"
	own["tabindex"] = "0"
	own["x-data"] = "huiToast({ duration: " + strconv.Itoa(duration) + " })"
	if duration > 0 {
		own["style"] = "--hui-toast-duration: " + strconv.Itoa(duration) + "ms"
	}
	own[":data-paused"] = "paused"
	own["x-show"] = "open"
	own["data-open"] = true
	own[":data-open"] = "open"
	own["@pointerenter"] = chain("pause()", attrString(p.Attrs["@pointerenter"]))
	own["@pointerleave"] = chain("resume()", attrString(p.Attrs["@pointerleave"]))
	own["@focusin"] = chain("pause()", attrString(p.Attrs["@focusin"]))
	own["@focusout"] = chain("resume($event)", attrString(p.Attrs["@focusout"]))

	attrs := p.mergeEvents(own, Events{OnKeydown: "keydown($event)"})
	setHandler(attrs, "@hui-toast:closed.self", p.OnClose)
	return attrs
}

func toastPartAttrs(ctx context.Context, b Base, suffix string) templ.Attributes {
	if id, ok := ctx.Value(toastKey).(string); ok {
		b.ID = id + "-" + suffix
	}
	return b.merge(nil)
}

func toastActionAttrs(p ToastActionProps) templ.Attributes {
	return p.merge(templ.Attributes{"type": "button", "data-enabled": true})
}

func toastCloseAttrs(ctx context.Context, p ToastCloseProps) templ.Attributes {
	click := ""
	if _, ok := ctx.Value(toastKey).(string); ok {
		click = "dismiss()"
	}
	return p.mergeEvents(templ.Attributes{"type": "button", "data-enabled": true}, Events{OnClick: click})
}
