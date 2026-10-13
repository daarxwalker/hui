package headlessui

import (
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/pending.min.js
var pendingJS string

var pendingScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + pendingJS + "</script>")))

func PendingScript() templ.Component {
	return pendingScript.Once()
}
