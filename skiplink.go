package headlessui

import "github.com/a-h/templ"

type SkipLinkProps struct {
	Base
	Href string
}

func skipLinkAttrs(p SkipLinkProps) templ.Attributes {
	href := p.Href
	if href == "" {
		href = "#main"
	}
	return p.merge(templ.Attributes{"href": string(templ.URL(href))})
}
