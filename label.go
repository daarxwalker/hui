package headlessui

import (
	"context"

	"github.com/a-h/templ"
)

type LabelProps struct {
	Base
	For string
}

func labelAttrs(ctx context.Context, p LabelProps) templ.Attributes {
	if f, inField := fieldFrom(ctx); inField {
		p.ID = f.ID + "-label"
		return p.merge(f.state(templ.Attributes{"for": f.ID}))
	}
	own := make(templ.Attributes, 1)
	if p.For != "" {
		own["for"] = p.For
	}
	return p.merge(own)
}
