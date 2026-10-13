package headlessui

import (
	"strconv"
	"strings"
)

type Placement struct {
	Side       string
	Align      string
	SideOffset int
}

func idSafe(value string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			return r
		default:
			return '-'
		}
	}, value)
}

func anchorName(id string) string {
	return "--" + idSafe(id)
}

func anchorSide(side string) string {
	switch side {
	case "top", "left", "right":
		return side
	default:
		return "bottom"
	}
}

func anchorAlign(align string) string {
	switch align {
	case "center", "end":
		return align
	default:
		return "start"
	}
}

func anchorPosition(anchor string, p Placement) string {
	side, align := anchorSide(p.Side), anchorAlign(p.Align)
	vertical := side == "top" || side == "bottom"

	area := map[string]string{"start": "span-right", "center": "center", "end": "span-left"}
	fallback := "flip-block"
	if !vertical {
		area = map[string]string{"start": "span-bottom", "center": "center", "end": "span-top"}
		fallback = "flip-inline"
	}
	facing := map[string]string{"bottom": "top", "top": "bottom", "left": "right", "right": "left"}

	rules := make([]string, 0, 7)
	rules = append(rules, "position: fixed")
	if anchor != "" {
		rules = append(rules, "position-anchor: "+anchor)
	}
	rules = append(rules,
		"position-area: "+side+" "+area[align],
		"position-try-fallbacks: "+fallback,
		"inset: auto",
		"margin: 0",
	)
	if p.SideOffset != 0 {
		rules = append(rules, "margin-"+facing[side]+": "+strconv.Itoa(p.SideOffset)+"px")
	}
	return strings.Join(rules, "; ")
}
