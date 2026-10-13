package headlessui

import "github.com/a-h/templ"

type AlertProps struct {
	Base
	Polite bool
}

func alertAttrs(p AlertProps) templ.Attributes {
	role := "alert"
	if p.Polite {
		role = "status"
	}
	return p.merge(templ.Attributes{"role": role})
}
