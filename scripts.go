package headlessui

import (
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/cdn.min.js
var scriptsJS string

var scriptsHandle = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + scriptsJS + "</script>")))

func Scripts() templ.Component {
	return scriptsHandle.Once()
}
