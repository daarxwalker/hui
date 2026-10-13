package headlessui

import "github.com/a-h/templ"

const (
	announcerID       = "hui-announcer"
	announcerUrgentID = "hui-announcer-urgent"
)

type AnnouncerProps struct {
	Base
	Urgent bool
}

func announcerAttrs(p AnnouncerProps) templ.Attributes {
	own := make(templ.Attributes, 4)
	own["role"] = "status"
	own["aria-live"] = "polite"
	own["aria-atomic"] = "true"
	own["data-hui-announcer"] = true
	id := announcerID
	if p.Urgent {
		own["role"] = "alert"
		own["aria-live"] = "assertive"
		id = announcerUrgentID
	}
	if p.ID == "" {
		p.ID = id
	}
	return p.merge(own)
}
