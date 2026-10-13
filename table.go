package headlessui

import (
	"strconv"

	"github.com/a-h/templ"
)

type TableProps struct {
	Base
	RowCount int
}

type TableCaptionProps struct {
	Base
}

type TableHeaderProps struct {
	Base
}

type TableBodyProps struct {
	Base
}

type TableFooterProps struct {
	Base
}

type TableRowProps struct {
	Base
	Index    int
	Selected bool
}

type TableHeadProps struct {
	Base
	Scope string
	Sort  string
}

type TableCellProps struct {
	Base
	Label string
}

type TableLoaderProps struct {
	Base
	Colspan int
}

func tableAttrs(p TableProps) templ.Attributes {
	own := templ.Attributes{"role": "table"}
	if p.RowCount != 0 {
		own["aria-rowcount"] = strconv.Itoa(p.RowCount)
	}
	return p.merge(own)
}

func tableLoaderAttrs(p TableLoaderProps) (templ.Attributes, string) {
	colspan := ""
	if p.Colspan > 1 {
		colspan = strconv.Itoa(p.Colspan)
	}
	return p.merge(templ.Attributes{"aria-hidden": "true", "data-hui-loader": true}), colspan
}

func tableRowAttrs(p TableRowProps) templ.Attributes {
	own := templ.Attributes{"role": "row"}
	if p.Index > 0 {
		own["aria-rowindex"] = strconv.Itoa(p.Index)
	}
	if p.Selected {
		own["data-selected"] = true
	}
	return p.merge(own)
}

func tableHeadAttrs(p TableHeadProps) templ.Attributes {
	own := make(templ.Attributes, 3)
	scope := p.Scope
	if scope == "" {
		scope = "col"
	}
	own["scope"] = scope
	own["role"] = "columnheader"
	if scope == "row" || scope == "rowgroup" {
		own["role"] = "rowheader"
	}
	switch p.Sort {
	case "ascending", "descending", "none", "other":
		own["aria-sort"] = p.Sort
		own["data-sort"] = p.Sort
	}
	return p.merge(own)
}

func tableSectionAttrs(b Base) templ.Attributes {
	return b.merge(templ.Attributes{"role": "rowgroup"})
}

func tableCellAttrs(p TableCellProps) templ.Attributes {
	own := templ.Attributes{"role": "cell"}
	if p.Label != "" {
		own["data-label"] = p.Label
	}
	return p.merge(own)
}
