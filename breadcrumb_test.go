package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestBreadcrumb(t *testing.T) {
	t.Run("renders a labelled navigation list", func(t *testing.T) {
		html := testRender(t, Breadcrumb(BreadcrumbProps{}),
			testWithChildren(BreadcrumbList(BreadcrumbListProps{}),
				testWithChildren(BreadcrumbItem(BreadcrumbItemProps{}), BreadcrumbLink(BreadcrumbLinkProps{Href: "/"})),
				BreadcrumbSeparator(BreadcrumbSeparatorProps{}),
				testWithChildren(BreadcrumbItem(BreadcrumbItemProps{}), BreadcrumbLink(BreadcrumbLinkProps{Current: true})),
			),
		)

		assert.Equal(t, `<nav aria-label="Breadcrumb"><ol><li><a href="/"></a></li><li aria-hidden="true" role="presentation"></li><li><span aria-current="page" data-current></span></li></ol></nav>`, html)
	})

	t.Run("the label can be translated", func(t *testing.T) {
		label := gofakeit.Word()

		assert.Equal(t, label, breadcrumbAttrs(BreadcrumbProps{Label: label})["aria-label"])
	})

	t.Run("the current page can stay a link", func(t *testing.T) {
		href := "/" + gofakeit.Word()

		attrs, isLink := breadcrumbLinkAttrs(BreadcrumbLinkProps{Href: href, Current: true})

		assert.True(t, isLink)
		assert.Equal(t, href, attrs["href"])
		assert.Equal(t, "page", attrs["aria-current"])
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		attrs, _ := breadcrumbLinkAttrs(BreadcrumbLinkProps{
			Href:    "/docs",
			Current: true,
			Base: Base{
				Class:  "underline",
				Attrs:  templ.Attributes{"aria-current": "step", "hx-boost": "true"},
				Events: Events{OnClick: "track()"},
			},
		})

		assert.Equal(t, "page", attrs["aria-current"])
		assert.Equal(t, "underline", attrs["class"])
		assert.Equal(t, "true", attrs["hx-boost"])
		assert.Equal(t, "track()", attrs["@click"])
	})
}
