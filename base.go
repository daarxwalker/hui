package headlessui

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/a-h/templ"
)

type Base struct {
	Events
	ID    string
	Class string
	Attrs templ.Attributes
}

type Events struct {
	OnFocus   string
	OnBlur    string
	OnInput   string
	OnChange  string
	OnKeydown string
	OnClick   string
}

var MergeClass = func(existing, added string) string {
	switch {
	case existing == "":
		return added
	case added == "":
		return existing
	default:
		return existing + " " + added
	}
}

func (b Base) merge(own templ.Attributes) templ.Attributes {
	return b.mergeEvents(own, Events{})
}

func (b Base) mergeEvents(own templ.Attributes, internal Events) templ.Attributes {
	attrs := own
	if attrs == nil {
		attrs = make(templ.Attributes, len(b.Attrs)+2)
	}
	ownStyle := attrString(attrs["style"])
	for k, v := range b.Attrs {
		if _, taken := attrs[k]; !taken {
			attrs[k] = v
		}
	}
	if style := joinStyle(ownStyle, attrString(b.Attrs["style"])); style != "" {
		attrs["style"] = style
	}
	if b.ID != "" {
		attrs["id"] = b.ID
	}
	if class := MergeClass(attrString(attrs["class"]), b.Class); class != "" {
		attrs["class"] = class
	}
	setHandler(attrs, "@focus", chain(internal.OnFocus, b.OnFocus))
	setHandler(attrs, "@blur", chain(internal.OnBlur, b.OnBlur))
	setHandler(attrs, "@input", chain(internal.OnInput, b.OnInput))
	setHandler(attrs, "@change", chain(internal.OnChange, b.OnChange))
	setHandler(attrs, "@keydown", chain(internal.OnKeydown, b.OnKeydown))
	setHandler(attrs, "@click", chain(internal.OnClick, b.OnClick))
	return attrs
}

func joinStyle(own, added string) string {
	own = strings.TrimSuffix(strings.TrimSpace(own), ";")
	added = strings.TrimSpace(added)
	switch {
	case own == "":
		return added
	case added == "":
		return own
	default:
		return own + "; " + added
	}
}

func attrString(value any) string {
	text, _ := value.(string)
	return text
}

func setHandler(attrs templ.Attributes, name, expr string) {
	if expr == "" {
		return
	}
	if previous := attrString(attrs[name]); previous != "" {
		expr = chain(previous, expr)
	}
	attrs[name] = expr
}

func chain(internal, user string) string {
	internal = strings.TrimSpace(internal)
	user = strings.TrimSpace(user)
	switch {
	case internal == "":
		return user
	case user == "":
		return internal
	default:
		return strings.TrimSuffix(internal, ";") + "; " + user
	}
}

type idScopeKeyType struct{}

var idScopeKey = idScopeKeyType{}

type idScope struct {
	prefix string
	seq    atomic.Uint64
}

func WithIDScope(ctx context.Context, prefix string) context.Context {
	if prefix == "" {
		prefix = "id"
	}
	return context.WithValue(ctx, idScopeKey, &idScope{prefix: prefix})
}

var idSeq atomic.Uint64

func nextID(ctx context.Context, prefix string) string {
	if scope, ok := ctx.Value(idScopeKey).(*idScope); ok {
		return scope.prefix + "-" + strconv.FormatUint(scope.seq.Add(1), 36)
	}
	return prefix + "-" + strconv.FormatUint(idSeq.Add(1), 36)
}

func jsString(value string) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}
