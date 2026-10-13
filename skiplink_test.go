package headlessui

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

func TestSkipLink(t *testing.T) {
	t.Run("jumps to main by default", func(t *testing.T) {
		html := testRender(t, SkipLink(SkipLinkProps{Class: "sr-only focus:not-sr-only"}), templ.Raw("Skip to content"))

		assert.Equal(t, `<a class="sr-only focus:not-sr-only" href="#main">Skip to content</a>`, html)
	})

	t.Run("target can be changed", func(t *testing.T) {
		assert.Equal(t, "#content", skipLinkAttrs(SkipLinkProps{Href: "#content"})["href"])
	})
}
