package headlessui

import "github.com/a-h/templ"

type BreadcrumbProps struct {
	Base
	Label string
}

type BreadcrumbListProps struct {
	Base
}

type BreadcrumbItemProps struct {
	Base
}

type BreadcrumbLinkProps struct {
	Base
	Href    string
	Current bool
}

type BreadcrumbSeparatorProps struct {
	Base
}

func breadcrumbAttrs(p BreadcrumbProps) templ.Attributes {
	label := p.Label
	if label == "" {
		label = "Breadcrumb"
	}
	return p.merge(templ.Attributes{"aria-label": label})
}

func breadcrumbLinkAttrs(p BreadcrumbLinkProps) (templ.Attributes, bool) {
	own := make(templ.Attributes, 3)
	if p.Href != "" {
		own["href"] = string(templ.URL(p.Href))
	}
	if p.Current {
		own["aria-current"] = "page"
		own["data-current"] = true
	}
	return p.merge(own), p.Href != ""
}

func breadcrumbSeparatorAttrs(p BreadcrumbSeparatorProps) templ.Attributes {
	return p.merge(templ.Attributes{"role": "presentation", "aria-hidden": "true"})
}
