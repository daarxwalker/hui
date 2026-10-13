package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	t.Run("renders native table parts", func(t *testing.T) {
		html := testRender(t, Table(TableProps{}),
			TableCaption(TableCaptionProps{}),
			testWithChildren(TableHeader(TableHeaderProps{}), testWithChildren(TableRow(TableRowProps{}), TableHead(TableHeadProps{}))),
			testWithChildren(TableBody(TableBodyProps{}), testWithChildren(TableRow(TableRowProps{Selected: true}), TableCell(TableCellProps{}))),
			TableFooter(TableFooterProps{}),
		)

		assert.Equal(t, `<table role="table"><caption></caption><thead role="rowgroup"><tr role="row"><th role="columnheader" scope="col"></th></tr></thead><tbody role="rowgroup"><tr data-selected role="row"><td role="cell"></td></tr></tbody><tfoot role="rowgroup"></tfoot></table>`, html)
	})

	t.Run("rows announce their place in a table that loads more", func(t *testing.T) {
		html := testRender(t, Table(TableProps{RowCount: 120}),
			testWithChildren(TableBody(TableBodyProps{}),
				testWithChildren(TableRow(TableRowProps{Index: 2}), TableCell(TableCellProps{})),
				testWithChildren(TableLoader(TableLoaderProps{Colspan: 3, Attrs: templ.Attributes{"hx-get": "/rows?page=2", "hx-trigger": "intersect once", "hx-swap": "outerHTML"}}), templ.Raw("Loading")),
			),
		)

		assert.Equal(t, `<table aria-rowcount="120" role="table"><tbody role="rowgroup"><tr aria-rowindex="2" role="row"><td role="cell"></td></tr><tr aria-hidden="true" data-hui-loader hx-get="/rows?page=2" hx-swap="outerHTML" hx-trigger="intersect once"><td colspan="3">Loading</td></tr></tbody></table>`, html)
	})

	t.Run("an unknown row count is announced as such", func(t *testing.T) {
		assert.Equal(t, "-1", tableAttrs(TableProps{RowCount: -1})["aria-rowcount"])
		assert.NotContains(t, tableAttrs(TableProps{}), "aria-rowcount")
		assert.NotContains(t, tableRowAttrs(TableRowProps{}), "aria-rowindex")
	})

	t.Run("sorting is announced on the header", func(t *testing.T) {
		attrs := tableHeadAttrs(TableHeadProps{Sort: "descending"})

		assert.Equal(t, "descending", attrs["aria-sort"])
		assert.Equal(t, "descending", attrs["data-sort"])
	})

	t.Run("unknown sort values are not rendered", func(t *testing.T) {
		attrs := tableHeadAttrs(TableHeadProps{Sort: "desc"})

		assert.NotContains(t, attrs, "aria-sort")
		assert.NotContains(t, attrs, "data-sort")
	})

	t.Run("row headers use their own scope", func(t *testing.T) {
		assert.Equal(t, "row", tableHeadAttrs(TableHeadProps{Scope: "row"})["scope"])
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		attrs := tableRowAttrs(TableRowProps{Selected: true, Base: Base{
			Class:  "hover:bg-gray-50",
			Attrs:  templ.Attributes{"hx-get": "/orders/1", "hx-trigger": "click"},
			Events: Events{OnClick: "track()"},
		}})

		assert.Equal(t, true, attrs["data-selected"])
		assert.Equal(t, "hover:bg-gray-50", attrs["class"])
		assert.Equal(t, "/orders/1", attrs["hx-get"])
		assert.Equal(t, "track()", attrs["@click"])
	})

	t.Run("cells carry their column label for a stacked layout", func(t *testing.T) {
		label := gofakeit.Word()

		assert.Equal(t, label, tableCellAttrs(TableCellProps{Label: label})["data-label"])
		assert.NotContains(t, tableCellAttrs(TableCellProps{}), "data-label")
	})

	t.Run("a header of a row is a rowheader", func(t *testing.T) {
		assert.Equal(t, "rowheader", tableHeadAttrs(TableHeadProps{Scope: "row"})["role"])
		assert.Equal(t, "columnheader", tableHeadAttrs(TableHeadProps{})["role"])
	})
}
