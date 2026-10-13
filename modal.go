package headlessui

import (
	"context"
	_ "embed"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/modal.min.js
var modalJS string

var modalScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + modalJS + "</script>")))

func ModalScript() templ.Component {
	return modalScript.Once()
}

const (
	ModalShow   = "modalShow()"
	ModalHide   = "modalHide()"
	ModalToggle = "modalToggle()"
)

type ModalProps struct {
	Base
	Open        bool
	Alert       bool
	Persist     bool
	Static      bool
	AllowScroll bool
	OnOpen      string
	OnClose     string
}

type ModalTriggerProps struct {
	Base
	Disabled bool
}

type ModalOverlayProps struct {
	Base
}

type ModalContentProps struct {
	Base
}

type ModalTitleProps struct {
	Base
}

type ModalDescriptionProps struct {
	Base
}

type ModalCloseProps struct {
	Base
}

type modalKeyType struct{}

var modalKey = modalKeyType{}

type modalState struct {
	id        string
	open      bool
	alert     bool
	inOverlay bool
}

func modalFrom(ctx context.Context) (modalState, bool) {
	state, ok := ctx.Value(modalKey).(modalState)
	return state, ok
}

func modalContext(ctx context.Context, p ModalProps) (context.Context, ModalProps) {
	if p.ID == "" {
		p.ID = nextID(ctx, "modal")
	}
	return context.WithValue(ctx, modalKey, modalState{id: p.ID, open: p.Open, alert: p.Alert}), p
}

func modalOverlayContext(ctx context.Context) context.Context {
	m, ok := modalFrom(ctx)
	if !ok {
		return ctx
	}
	m.inOverlay = true
	return context.WithValue(ctx, modalKey, m)
}

func modalAttrs(p ModalProps) templ.Attributes {
	own := make(templ.Attributes, 6)
	own["x-data"] = modalData(p)
	own["x-modelable"] = "open"
	own[":data-open"] = "open"
	if p.Open {
		own["data-open"] = true
	}

	attrs := p.merge(own)
	setHandler(attrs, "@hui-modal:opened.self", p.OnOpen)
	setHandler(attrs, "@hui-modal:closed.self", p.OnClose)
	return attrs
}

func modalData(p ModalProps) string {
	config := make([]string, 0, 5)
	config = append(config, "open: "+strconv.FormatBool(p.Open))
	if p.Persist {
		config = append(config, "persist: '"+p.ID+"'")
	}
	if p.Static {
		config = append(config, "static: true")
	}
	if p.AllowScroll {
		config = append(config, "allowScroll: true")
	}
	if p.Alert {
		config = append(config, "alert: true")
	}
	return "huiModal({ " + strings.Join(config, ", ") + " })"
}

func modalTriggerAttrs(ctx context.Context, p ModalTriggerProps) templ.Attributes {
	own := make(templ.Attributes, 10)
	own["type"] = "button"
	own["aria-haspopup"] = "dialog"

	click := ""
	if m, ok := modalFrom(ctx); ok {
		click = ModalShow
		own["aria-controls"] = m.id + "-content"
		own["aria-expanded"] = strconv.FormatBool(m.open)
		own[":aria-expanded"] = "open"
		own[":data-open"] = "open"
		if m.open {
			own["data-open"] = true
		}
	}
	if p.Disabled {
		own["disabled"] = true
	}
	buttonState(own, p.Disabled)

	return p.mergeEvents(own, Events{OnClick: click})
}

func modalOverlayAttrs(ctx context.Context, p ModalOverlayProps) (templ.Attributes, string) {
	own := make(templ.Attributes, 6)
	own["aria-hidden"] = "true"

	teleport := ""
	if m, ok := modalFrom(ctx); ok {
		p.ID = m.id + "-overlay"
		teleport = overlayTeleport(false)
		transition(own, p.Attrs, "x-transition.opacity")
		overlayVisibility(own, p.Attrs, m.open)
	}

	return p.merge(own), teleport
}

func modalContentAttrs(ctx context.Context, p ModalContentProps) (templ.Attributes, string) {
	own := make(templ.Attributes, 10)
	own["role"] = "dialog"
	own["aria-modal"] = "true"
	own["tabindex"] = "-1"

	teleport := ""
	if m, ok := modalFrom(ctx); ok {
		p.ID = m.id + "-content"
		teleport = overlayTeleport(m.inOverlay)
		if m.alert {
			own["role"] = "alertdialog"
		}
		own["aria-labelledby"] = m.id + "-title"
		own["aria-describedby"] = m.id + "-description"
		overlayVisibility(own, p.Attrs, m.open)
	}

	return p.merge(own), teleport
}

func modalTitleAttrs(ctx context.Context, p ModalTitleProps) templ.Attributes {
	if m, ok := modalFrom(ctx); ok {
		p.ID = m.id + "-title"
	}
	return p.merge(nil)
}

func modalDescriptionAttrs(ctx context.Context, p ModalDescriptionProps) templ.Attributes {
	if m, ok := modalFrom(ctx); ok {
		p.ID = m.id + "-description"
	}
	return p.merge(nil)
}

func modalCloseAttrs(ctx context.Context, p ModalCloseProps) templ.Attributes {
	click := ""
	if _, ok := modalFrom(ctx); ok {
		click = ModalHide
	}
	return p.mergeEvents(templ.Attributes{"type": "button", "data-enabled": true}, Events{OnClick: click})
}
