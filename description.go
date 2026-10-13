package headlessui

import (
	"context"

	"github.com/a-h/templ"
)

type DescriptionProps struct {
	Base
}

func descriptionAttrs(ctx context.Context, p DescriptionProps) (attrs templ.Attributes, text string, render bool) {
	f, inField := ownerFrom(ctx)
	if !inField {
		return p.merge(nil), "", true
	}
	if f.Description == "" {
		return nil, "", false
	}
	p.ID = f.descriptionID()
	return p.merge(f.state(nil)), f.Description, true
}
