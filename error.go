package headlessui

import (
	"context"

	"github.com/a-h/templ"
)

type ErrorProps struct {
	Base
}

func errorAttrs(ctx context.Context, p ErrorProps) (attrs templ.Attributes, text string, render bool) {
	own := templ.Attributes{"role": "alert"}
	f, inField := ownerFrom(ctx)
	if !inField {
		return p.merge(own), "", true
	}
	if f.Error == "" {
		return nil, "", false
	}
	p.ID = f.errorID()
	return p.merge(f.state(own)), f.Error, true
}
