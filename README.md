# hui

Headless (unstyled) UI components built on [templ](https://templ.guide). The API shape
follows [bits-ui](https://bits-ui.com/). Accessibility comes first: the correct default
state is rendered from Go, so the components work with plain HTML, with
[htmx](https://htmx.org), or with [Alpine.js](https://alpinejs.dev) when you opt into it.

```bash
go get github.com/daarxwalker/hui
```

**This README is written for a coding agent** that builds a project's own UI library on top
of hui. Read it in this order:

1. [Requirements](#requirements): what the page has to load.
2. [Building a project UI library](#building-a-project-ui-library): the pattern every styled
   component follows, with complete code.
3. [Building a showcase app](#building-a-showcase-app): a small web app with every
   component in every state, to look at and try by hand.
4. [Styling reference](#styling-reference): every part, the element it renders and the
   state attributes to style.
5. The component sections from [What it looks like](#what-it-looks-like) on: props,
   behaviour and the edge cases of each component.

## Requirements

- **Go 1.27+.** Props embed `Base` and literals use embedded field selectors
  (`TextInputProps{ID: "email"}`).
- **templ v0.3.1020+.** The generated `*_templ.go` files are committed, so the library
  itself builds without the templ CLI. Your own `.templ` files still need it.
- **Tailwind v4.** Every state is a bare `data-*` attribute, which is what the
  `data-invalid:` style variants expect. The core itself ships no CSS.
- **htmx 4** for anything that talks to the server (`<hx-partial>`, `hx-include="this"`).
- **Alpine 3** for the interactive components (see [Alpine](#alpine)).
- In the layout, in this order: htmx, `@hui.Scripts()`, Alpine. Then `@hui.Portals()` once,
  directly in `<body>`, outside anything htmx swaps. Every overlay teleports into it:

```templ
<head>
    <script src="https://cdn.jsdelivr.net/npm/htmx.org@4.0.0/dist/htmx.min.js"></script>
    @hui.Scripts()
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3/dist/cdn.min.js"></script>
</head>
<body>
    ...
    @hui.Portals(hui.PortalsProps{}) {
        @ui.Toaster(ui.ToasterProps{Label: "Notifications"})
        @hui.Announcer(hui.AnnouncerProps{Class: "sr-only"})
    }
</body>
```

## Building a project UI library

hui renders structure, accessibility and state. Every project adds the look in its
own small package of styled wrappers. Do not fork the core and do not restyle it inline
at every call site. Write one wrapper per headless component, once, and use only the
wrappers in the app.

### Layout

```
ui/
  shared.go        Color, Variant, Size, shared class strings, layout helper
  button.go        ButtonProps, ButtonLook, class maps, createButtonProps
  button.templ     Button: calls the headless component with the built props
  modal.go
  modal.templ
  ...
```

- One flat package (`ui`), one `.go` + `.templ` pair per component, mirroring the core.
- The `.go` file holds all logic: props, class maps and a `create<Part>Props` function per
  part. The `.templ` file only calls the headless component with the built props and
  renders `{ children... }`.
- Import the core as `hui` so that `ui.Button` and `hui.Button` read apart.

### Shared tokens

```go
package ui

import (
	"maps"

	"github.com/a-h/templ"
	hui "github.com/daarxwalker/hui"
)

type Color string

type Variant string

type Size string

const (
	ColorDefault Color = "default"
	ColorNeutral Color = "neutral"
	ColorBrand   Color = "brand"
	ColorSuccess Color = "success"
	ColorWarning Color = "warning"
	ColorError   Color = "error"
)

const (
	VariantSolid   Variant = "solid"
	VariantOutline Variant = "outline"
	VariantGhost   Variant = "ghost"
)

const (
	SizeXS Size = "xs"
	SizeSM Size = "sm"
	SizeMD Size = "md"
	SizeLG Size = "lg"
)

var controlSizes = map[Size]string{
	SizeXS: "h-6 text-xs",
	SizeSM: "h-8 text-xs",
	SizeMD: "h-10 text-sm",
	SizeLG: "h-12 text-base",
}

const focusClass = "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-focus"

const inputClass = focusClass + " w-full rounded-md border border-default-divide bg-surface text-default-inverse transition placeholder:text-muted data-invalid:border-error-divide data-readonly:bg-default-subtle disabled:cursor-not-allowed disabled:opacity-50"

func createLayoutAttrs(class, own string, attrs templ.Attributes) templ.Attributes {
	merged := maps.Clone(attrs)
	if merged == nil {
		merged = templ.Attributes{}
	}
	if fromAttrs, ok := merged["class"].(string); ok {
		class = hui.MergeClass(class, fromAttrs)
	}
	merged["class"] = hui.MergeClass(class, own)
	return merged
}
```

`createLayoutAttrs` is for plain layout parts that have no headless counterpart (a modal
footer, a card body). They get the same `Class` + `Attrs` contract as everything else.

### A simple component

`button.go`:

```go
package ui

import (
	"cmp"
	"strings"

	hui "github.com/daarxwalker/hui"
)

type ButtonLook struct {
	Color     Color
	Variant   Variant
	Size      Size
	Square    bool
	Fullwidth bool
}

type ButtonProps struct {
	hui.ButtonProps
	ButtonLook
}

const buttonClass = focusClass + " inline-flex shrink-0 items-center justify-center gap-2 rounded-md font-semibold whitespace-nowrap transition data-enabled:active:scale-95 data-disabled:cursor-not-allowed data-disabled:opacity-50 data-pending:pointer-events-none"

var buttonPaddings = map[Size]string{
	SizeXS: "px-2",
	SizeSM: "px-3",
	SizeMD: "px-4",
	SizeLG: "px-5",
}

var buttonSquareWidths = map[Size]string{
	SizeXS: "w-6",
	SizeSM: "w-8",
	SizeMD: "w-10",
	SizeLG: "w-12",
}

var buttonColors = map[Variant]map[Color]string{
	VariantSolid: {
		ColorDefault: "bg-default text-default-contrast data-enabled:hover:bg-default-hover",
		ColorBrand:   "bg-brand text-brand-contrast data-enabled:hover:bg-brand-hover",
		ColorError:   "bg-error text-error-contrast data-enabled:hover:bg-error-hover",
	},
	VariantOutline: {
		ColorDefault: "border border-default-divide bg-surface text-default-inverse data-enabled:hover:bg-default-subtle",
		ColorBrand:   "border border-brand-divide bg-surface text-brand-inverse data-enabled:hover:bg-brand-subtle",
		ColorError:   "border border-error-divide bg-surface text-error-inverse data-enabled:hover:bg-error-subtle",
	},
	VariantGhost: {
		ColorDefault: "text-default-inverse data-enabled:hover:bg-default-subtle",
		ColorBrand:   "text-brand-inverse data-enabled:hover:bg-brand-subtle",
		ColorError:   "text-error-inverse data-enabled:hover:bg-error-subtle",
	},
}

func createButtonClass(l ButtonLook) string {
	color := cmp.Or(l.Color, ColorDefault)
	variant := cmp.Or(l.Variant, VariantSolid)
	size := cmp.Or(l.Size, SizeMD)

	widths := buttonPaddings
	if l.Square {
		widths = buttonSquareWidths
	}

	classes := []string{buttonClass, controlSizes[size], widths[size], buttonColors[variant][color]}
	if l.Fullwidth {
		classes = append(classes, "w-full")
	}
	return strings.Join(classes, " ")
}

func createButtonProps(p ButtonProps) hui.ButtonProps {
	props := p.ButtonProps
	props.Class = hui.MergeClass(createButtonClass(p.ButtonLook), p.Class)
	return props
}
```

`button.templ`:

```templ
package ui

import hui "github.com/daarxwalker/hui"

templ Button(props ButtonProps) {
	@hui.Button(createButtonProps(props)) {
		if props.Square {
			<span class="sr-only">
				{ children... }
			</span>
		} else {
			{ children... }
		}
	}
}
```

Use:

```templ
@ui.Button(ui.ButtonProps{Type: "submit", Color: ui.ColorBrand, Disabled: !valid}) {
    Save
}
```

What makes it work:

- **Embedding.** `ButtonProps` embeds `hui.ButtonProps`, so every headless prop (`ID`,
  `Name`, `Type`, `Href`, `Disabled`, `Pending`, `Class`, `Attrs`, `Events`, …) is set in
  the same literal and passes through untouched. The look lives in its own struct
  (`ButtonLook`) so that other parts can reuse it.
- **Look fields must not shadow headless ones.** A field of the outer struct hides a
  promoted field with the same name. Before you add `Size`, `Color` or `Label` to a look,
  check that the headless props have no field of that name.
- **Class builder.** `create<Part>Class` puts the look together from maps keyed by size,
  variant and color, with the defaults from `cmp.Or`. Keep each class string whole in
  the source, because Tailwind cannot see a class assembled from fragments.
- **Props builder.** `create<Part>Props` copies the headless props and sets
  `props.Class = hui.MergeClass(own, p.Class)`. The caller's `Class` comes last and wins.
  A `class` in `Attrs` is merged by the core too.
- **Thin templ.** The templ file never adds attributes of its own. `aria-*`, `role`,
  `x-*` and `@…` belong to the core, and anything extra goes through `Attrs` or `Events`,
  because the browser silently drops one of two duplicate attributes.

### A component with parts

Compound components get one wrapper per part. A styled part may render several headless
parts: `ModalContent` below renders the overlay and the dialog, and takes the overlay's
props as a field. Triggers and close buttons reuse `ButtonLook` and `createButtonClass`,
so every button-like part looks like a `Button`.

`modal.go`:

```go
package ui

import (
	"cmp"

	"github.com/a-h/templ"
	hui "github.com/daarxwalker/hui"
)

type ModalProps struct {
	hui.ModalProps
}

type ModalTriggerProps struct {
	hui.ModalTriggerProps
	ButtonLook
}

type ModalContentProps struct {
	hui.ModalContentProps
	Size    Size
	Overlay hui.ModalOverlayProps
}

type modalContentParts struct {
	Overlay hui.ModalOverlayProps
	Content hui.ModalContentProps
}

type ModalTitleProps struct {
	hui.ModalTitleProps
}

type ModalDescriptionProps struct {
	hui.ModalDescriptionProps
}

type ModalFooterProps struct {
	Class string
	Attrs templ.Attributes
}

type ModalCloseProps struct {
	hui.ModalCloseProps
	ButtonLook
}

const (
	modalOverlayClass     = "fixed inset-0 z-50 bg-overlay/80"
	modalContentClass     = "fixed inset-0 z-50 m-auto h-fit max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] overflow-y-auto rounded-lg border border-default-divide bg-surface p-6 text-default-inverse shadow-lg focus:outline-none md:max-h-[calc(100dvh-4rem)]"
	modalTitleClass       = "pr-8 text-base font-semibold text-default-inverse"
	modalDescriptionClass = "mt-1 text-sm text-body"
	modalFooterClass      = "mt-6 flex flex-wrap justify-end gap-2"
)

var modalWidths = map[Size]string{
	SizeSM: "max-w-sm",
	SizeMD: "max-w-lg",
	SizeLG: "max-w-2xl",
}

func createModalTriggerProps(p ModalTriggerProps) hui.ModalTriggerProps {
	props := p.ModalTriggerProps
	props.Class = hui.MergeClass(createButtonClass(p.ButtonLook), p.Class)
	return props
}

func createModalContentProps(p ModalContentProps) modalContentParts {
	parts := modalContentParts{Overlay: p.Overlay, Content: p.ModalContentProps}
	parts.Overlay.Class = hui.MergeClass(modalOverlayClass, p.Overlay.Class)
	parts.Content.Class = hui.MergeClass(modalContentClass+" "+modalWidths[cmp.Or(p.Size, SizeMD)], p.Class)
	return parts
}

func createModalTitleProps(p ModalTitleProps) hui.ModalTitleProps {
	props := p.ModalTitleProps
	props.Class = hui.MergeClass(modalTitleClass, p.Class)
	return props
}

func createModalDescriptionProps(p ModalDescriptionProps) hui.ModalDescriptionProps {
	props := p.ModalDescriptionProps
	props.Class = hui.MergeClass(modalDescriptionClass, p.Class)
	return props
}

func createModalCloseProps(p ModalCloseProps) hui.ModalCloseProps {
	props := p.ModalCloseProps
	props.Class = hui.MergeClass(createButtonClass(p.ButtonLook), p.Class)
	return props
}
```

`modal.templ`:

```templ
package ui

import hui "github.com/daarxwalker/hui"

templ Modal(props ModalProps) {
	@hui.Modal(props.ModalProps) {
		{ children... }
	}
}

templ ModalTrigger(props ModalTriggerProps) {
	@hui.ModalTrigger(createModalTriggerProps(props)) {
		{ children... }
	}
}

templ ModalContent(props ModalContentProps) {
	{{ parts := createModalContentProps(props) }}
	@hui.ModalOverlay(parts.Overlay)
	@hui.ModalContent(parts.Content) {
		{ children... }
	}
}

templ ModalTitle(props ModalTitleProps) {
	@hui.ModalTitle(createModalTitleProps(props)) {
		{ children... }
	}
}

templ ModalDescription(props ModalDescriptionProps) {
	@hui.ModalDescription(createModalDescriptionProps(props)) {
		{ children... }
	}
}

templ ModalFooter(props ModalFooterProps) {
	<div { createLayoutAttrs(modalFooterClass, props.Class, props.Attrs)... }>
		{ children... }
	</div>
}

templ ModalClose(props ModalCloseProps) {
	@hui.ModalClose(createModalCloseProps(props)) {
		{ children... }
	}
}
```

Use:

```templ
@ui.Modal(ui.ModalProps{ID: "delete-order"}) {
    @ui.ModalTrigger(ui.ModalTriggerProps{Color: ui.ColorError}) {
        Delete
    }
    @ui.ModalContent(ui.ModalContentProps{Size: ui.SizeSM}) {
        @ui.ModalTitle(ui.ModalTitleProps{}) {
            Delete the order?
        }
        @ui.ModalDescription(ui.ModalDescriptionProps{}) {
            This cannot be undone.
        }
        @ui.ModalFooter(ui.ModalFooterProps{}) {
            @ui.ModalClose(ui.ModalCloseProps{ButtonLook: ui.ButtonLook{Variant: ui.VariantOutline}}) {
                Cancel
            }
        }
    }
}
```

The same pattern covers every other compound component (`Dropdown`, `Listbox`, `Combobox`,
`Tabs`, `Collapsible`, `Table`, `Toast`, …): one wrapper per part, and a styled part may
render several headless parts.

### Design tokens

Components use semantic Tailwind tokens only, never a raw palette (`bg-red-600`). A
project then changes its look in `@theme` alone. Every color role has the same six shades:

| Token | Used for |
|---|---|
| `<role>` | solid background |
| `<role>-hover` | solid background on hover |
| `<role>-contrast` | text on the solid background |
| `<role>-subtle` | tinted background (ghost hover, alert, selected row) |
| `<role>-inverse` | text and icons on a light surface |
| `<role>-divide` | borders and dividers |

The roles are `default`, `neutral`, `brand`, `success`, `warning` and `error`, plus the
single tokens `surface`, `focus`, `muted`, `label`, `body` and `overlay`:

```css
@theme {
    --color-surface: var(--color-white);
    --color-focus: var(--color-brand-600);
    --color-muted: var(--color-zinc-500);
    --color-label: var(--color-zinc-600);
    --color-body: var(--color-zinc-700);
    --color-overlay: var(--color-zinc-950);

    --color-default: var(--color-zinc-100);
    --color-default-hover: var(--color-zinc-200);
    --color-default-contrast: var(--color-zinc-950);
    --color-default-subtle: var(--color-zinc-50);
    --color-default-inverse: var(--color-zinc-950);
    --color-default-divide: var(--color-zinc-300);

    --color-neutral: var(--color-zinc-950);
    --color-neutral-hover: var(--color-zinc-800);
    --color-neutral-contrast: var(--color-white);
    --color-neutral-subtle: var(--color-zinc-100);
    --color-neutral-inverse: var(--color-zinc-950);
    --color-neutral-divide: var(--color-zinc-950);

    --color-brand: var(--color-brand-300);
    --color-brand-hover: var(--color-brand-400);
    --color-brand-contrast: var(--color-brand-950);
    --color-brand-subtle: var(--color-brand-100);
    --color-brand-inverse: var(--color-brand-700);
    --color-brand-divide: var(--color-brand-500);

    --color-success: var(--color-green-700);
    --color-success-hover: var(--color-green-800);
    --color-success-contrast: var(--color-white);
    --color-success-subtle: var(--color-green-50);
    --color-success-inverse: var(--color-green-700);
    --color-success-divide: var(--color-green-700);

    --color-warning: var(--color-orange-700);
    --color-warning-hover: var(--color-orange-800);
    --color-warning-contrast: var(--color-white);
    --color-warning-subtle: var(--color-orange-50);
    --color-warning-inverse: var(--color-orange-700);
    --color-warning-divide: var(--color-orange-600);

    --color-error: var(--color-red-700);
    --color-error-hover: var(--color-red-600);
    --color-error-contrast: var(--color-white);
    --color-error-subtle: var(--color-red-50);
    --color-error-inverse: var(--color-red-700);
    --color-error-divide: var(--color-red-600);
}
```

`--color-brand-50` to `--color-brand-950` is the project's own brand scale.

### Styling state

State is never a class. Style it with Tailwind variants on the attributes the core renders
(full list in the [Styling reference](#styling-reference)):

| Variant | When |
|---|---|
| `data-invalid:` | field with an error, on the wrapper, label, control, description and error |
| `data-required:`, `data-readonly:` | the same parts of a required or read-only field |
| `data-disabled:` / `data-enabled:` | every core button carries exactly one of them; form parts carry `data-disabled` |
| `data-checked:` | `Checkbox`, `Switch`, `SwitchThumb`, `Radio`, checkbox and radio items of a `Dropdown` |
| `data-indeterminate:` | `Checkbox` in the mixed state |
| `data-open:` | open `Modal`, `Dropdown`, `Tooltip`, `Listbox`, `Combobox`, `Collapsible` and `Toast`, on the root, the trigger and the content |
| `data-highlighted:` | the active `ListboxItem` or `ComboboxItem` (dropdown items take real focus, use `focus:`) |
| `data-selected:` | selected listbox or combobox item, tab, tab panel and table row |
| `data-placeholder:` | `Listbox` trigger and value while nothing is selected |
| `data-pending:` | while an htmx request of this element runs, together with `aria-busy="true"` |
| `data-searching:` | from the first typed character until the server answers that query |
| `data-dragging:` | `FileUpload` while a file is dragged over it |
| `data-current:` | the current `BreadcrumbLink` |
| `data-sort:` / `aria-sort` | `TableHead` (`ascending`, `descending`, `none`, `other`) |
| `has-data-checked:` | a `Fieldset` option row whose control is checked |
| `in-data-pending:` | something without state of its own inside a pending element |

Every part carries its own state, so `group` is rarely needed.

### Class merging

`hui.MergeClass` joins strings by default. Plug in a Tailwind v4 aware merger once at
startup, so that `px-4` from the wrapper and `px-10` from the caller do not both survive:

```go
hui.MergeClass = func(existing, added string) string {
    return twmerge.Merge(existing, added)
}
```

`tailwind-merge-go` only covers Tailwind 3.0–3.4. Until there is a v4 merger, keep the
utilities of one property in one place.

### Overlay styling

The core positions anchored content (`Dropdown`, `Tooltip`, `Listbox`, `Combobox`) with
CSS anchor positioning and an inline `position: fixed`. The styled layer decides the
`z-index` and the width (`w-max` for a tooltip). `ModalOverlay` and `ModalContent` need
their own `fixed inset-0` and `z-index`. Use `Side`, `Align` and `SideOffset` for
placement, not margins.

### Checklist

For every headless component:

1. One wrapper per part, `<Part>Props` embedding `hui.<Part>Props` (plus a look struct
   where it makes sense), `create<Part>Props` in the `.go` file, thin `.templ`.
2. `Class`, `Attrs` and `Events` still work on every wrapper: the caller's class wins,
   `hx-*` and `x-*` in `Attrs` reach the element, an `Events` handler runs after the
   internal one.
3. No attribute the core sets for accessibility (`role`, `aria-*`, `id`, `for`, `x-*`,
   `@…`) is overwritten or duplicated.
4. Every state from the [Styling reference](#styling-reference) has a look: invalid,
   disabled, checked, open, highlighted, selected, pending.
5. Keyboard: a visible focus ring on everything focusable (`focus-visible:`).
6. It works in a `Field` and on its own, and inside an htmx partial.
7. It has a card in the [showcase app](#building-a-showcase-app).

## Building a showcase app

Next to the UI library, build a small web app that renders every component in every
state, so the library can be looked at and tried by hand. It is a plain `package main`
that imports `ui` and `hui` like any other code and holds no library logic.

```
showcase/
  main.go          ServeMux, the page and a few htmx endpoints
  showcase.templ   layout, showcaseCard, showcaseState and one card per component
```

Run it with `templ generate && go run ./showcase` and open http://localhost:5050.

**Server.** No framework, `net/http` is enough:

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		showcasePage().Render(r.Context(), w)
	})
	log.Println("listening on http://localhost:5050")
	log.Fatal(http.ListenAndServe(":5050", mux))
}
```

**Layout.** Tailwind from the browser CDN with the project's `@theme` inline, then
htmx, `@hui.Scripts()` and Alpine, and `Portals` with the `Toaster` and `Announcer` at the
end of `<body>`:

```templ
templ showcasePage() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>UI showcase</title>
			<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
			<style type="text/tailwindcss">
				@theme {
					/* the project's design tokens */
				}
			</style>
			<script src="https://cdn.jsdelivr.net/npm/htmx.org@4.0.0/dist/htmx.min.js"></script>
			@hui.Scripts()
			<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3/dist/cdn.min.js"></script>
		</head>
		<body class="bg-gray-50 p-6 text-gray-900">
			@ui.SkipLink(ui.SkipLinkProps{}) {
				Skip to content
			}
			<div id="main" tabindex="-1" class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
				@showcaseCard("Button") {
					@showcaseState("variants") {
						<div class="flex flex-wrap gap-2">
							@ui.Button(ui.ButtonProps{}) {
								Solid
							}
							@ui.Button(ui.ButtonProps{ButtonLook: ui.ButtonLook{Variant: ui.VariantOutline}}) {
								Outline
							}
						</div>
					}
					@showcaseState("disabled") {
						@ui.Button(ui.ButtonProps{Disabled: true}) {
							Disabled
						}
					}
				}
			</div>
			@hui.Portals(hui.PortalsProps{}) {
				@ui.Toaster(ui.ToasterProps{Label: "Notifications"})
				@hui.Announcer(hui.AnnouncerProps{Class: "sr-only"})
			}
		</body>
	</html>
}

templ showcaseCard(title string) {
	<section class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
		<h2 class="mb-4 text-lg font-semibold">{ title }</h2>
		<div class="grid gap-4">
			{ children... }
		</div>
	</section>
}

templ showcaseState(name string) {
	<div class="grid gap-1.5">
		<span class="font-mono text-xs text-gray-400">{ name }</span>
		<div>
			{ children... }
		</div>
	</div>
}
```

**Coverage.** One card per component, one `showcaseState` per state it has:

- default, with a value, disabled, read-only, invalid, every size, every variant × color,
- form components both on their own and inside a `Field`, choices inside a `Fieldset`,
- overlays with several `Side`/`Align` combinations, every modal size, a dropdown with
  checkbox and radio items, a listbox inside a modal (stacked layers),
- `Button` with `Pending`, an icon-only `Square` button, a link button (`Href`).

**htmx endpoints.** A few handlers over in-memory data make the server-driven parts
visible. Sleep a few hundred milliseconds in each, so `data-pending` and `data-searching`
are visible:

| Endpoint | Shows |
|---|---|
| `GET /countries` | `Combobox` search with paging and a `ComboboxLoader`, the partial rendered with `hui.WithCombobox` |
| `GET /listbox-countries` | searchable `Listbox` with paging, the partial rendered with `hui.WithListbox` |
| `GET /orders` | infinite `Table` with `TableLoader` and `aria-rowindex` |
| `POST /toast` | a `Toast` sent as `<hx-partial>` into `#hui-toaster` |
| `POST /announce`, `POST /alert` | `Announcer` and `Alert` |
| `POST /echo` | writes every submitted `name=value` back, to check what checkboxes, listboxes, file uploads and controls with `Form` really send |

## Styling reference

Every part, the element it renders and the attributes to style. **Go** means rendered on
the server, so the state is right before any JavaScript runs. **Alpine** means kept up to
date by the plugin. Go renders boolean attributes bare and Alpine may render them as
`="true"`, so always match on presence: the Tailwind variants `data-open:` and
`data-checked:` do exactly that. `aria-expanded`, `aria-selected` and `aria-checked` are
always `"true"` or `"false"`, never missing.

*Field state* below means `data-disabled`, `data-invalid`, `data-readonly` and
`data-required`, taken from the `Field` (or `Fieldset`) or from the control's own props.
*Button state* means exactly one of `data-enabled` and `data-disabled`.

There is no `data-side` or `data-align`: placement is inline CSS anchor positioning.
`data-state` exists on `Meter` only.

### Global

| Where | Attributes |
|---|---|
| element that sent an htmx request, its submitter and its `hx-target` | `data-pending`, `aria-busy="true"` (Alpine-free script). A server-rendered `data-pending` stays |
| listbox and combobox list during a request from inside | `data-pending`, `aria-busy="true"` |
| search input and list of `Listbox`/`Combobox` | `data-searching` while a typed query waits for the server |
| text inside an option that matches the query | `<mark data-hui-match>` |
| everything outside the top `Modal` | `inert` (not the toaster and the announcer) |

### Form

| Part | Element | Attributes |
|---|---|---|
| `Field` | `<div id="<id>-field">` | field state, `data-field="<name>"`, or `data-option="<value>"` for a choice inside a `Fieldset` |
| `Fieldset` | `<fieldset id="<id>-field">` | `data-disabled`, `data-invalid`, `data-required`, `data-field` |
| `Legend` | `<legend>` | the fieldset's state |
| `Label` | `<label for>` | field state |
| `Description` | `<p id="<id>-description">` | field state |
| `Error` | `<p role="alert" id="<id>-error">` | field state |
| `TextInput`, `Textarea` | `<input>`, `<textarea>` | field state, `aria-invalid`; with `Alpine: true` also `data-empty` and `data-focused` |
| `Checkbox` | `<button role="checkbox">` | `aria-checked` (`true`, `false`, `mixed`), `data-checked`, `data-indeterminate`, field state |
| `CheckboxIndicator` | `<span aria-hidden>` | shown only while checked (or indeterminate), style it through the checkbox |
| `Switch` | `<button role="switch">` | `aria-checked`, `data-checked`, field state |
| `SwitchThumb` | `<span aria-hidden>` | `data-checked` |
| `Radio` | `<button role="radio">` | `aria-checked`, `data-checked`, field state, roving `tabindex` |
| `RadioIndicator` | `<span aria-hidden>` | shown only while checked |
| `Range` | `<div>` | `data-disabled`, `data-invalid`, CSS variables `--hui-range-start` and `--hui-range-end` |
| `RangeThumb` | `<input type="range">` | `data-disabled`, `data-invalid`, `aria-invalid` |
| `Select` | `<select>` | field state, `aria-invalid`, native `:open` |
| `SelectTrigger` | `<button>` | button state, slots `data-slot="left\|right"` |
| `SelectItem`, `SelectGroup` | `<option>`, `<optgroup>` | `data-disabled`, native `:checked` |
| placeholder | `<option value="" data-placeholder>` | |
| `FileUpload` | `<div>` | `data-disabled`, `data-dragging` |
| `FileUploadInput` | `<input type="file">` | field state, `aria-invalid` |
| `FileUploadDropzone` | `<div>` | `data-dragging` |
| `FileUploadList`, `FileUploadFile` | `<ul>`, `<li>` | the list is hidden while it is empty |

### Listbox and Combobox

| Part | Element | Attributes |
|---|---|---|
| `Listbox`, `Combobox` | `<div>` | `data-open`, field state |
| `ListboxTrigger` | `<button role="combobox">` (no role when `Searchable`) | `aria-expanded`, `data-open`, `data-placeholder`, field state |
| `ListboxValue` | `<span>` | `data-placeholder` |
| `ListboxSearch` | `<input role="combobox">` | `aria-expanded`, `data-searching` |
| `ComboboxInput` | `<input role="combobox">` | `aria-expanded`, `data-open`, `data-searching`, field state |
| `ComboboxTrigger` | `<button tabindex="-1">` | `aria-expanded`, `data-open`, button state |
| `ListboxContent`, `ComboboxContent` | `<div>` in `#hui-portals` | `data-open` |
| `ListboxList`, `ComboboxList` | `<div role="listbox">` | `aria-multiselectable`, `data-searching`, `data-pending`, `aria-busy` |
| `ListboxItem`, `ComboboxItem` | `<div role="option">` | `aria-selected`, `data-selected`, `data-highlighted`, `data-disabled`, `aria-disabled`, slots `data-slot` |
| `…ItemIndicator` | `<span aria-hidden>` | shown only while selected, style it through `data-selected:` on the item |
| `…Group`, `…GroupLabel` | `<div role="group">`, `<div>` | |
| `…Separator` | `<div aria-hidden>` | |
| `…Empty` | `<div role="presentation">` | |
| `…Loader` | `<div aria-hidden data-hui-loader>` | |
| `ComboboxRemove` | `<button tabindex="-1">` | |

### Overlays

| Part | Element | Attributes |
|---|---|---|
| `Portals` | `<div id="hui-portals">` | |
| `Modal` | `<div>` | `data-open` |
| `ModalTrigger` | `<button aria-haspopup="dialog">` | `aria-expanded`, `data-open`, button state |
| `ModalOverlay` | `<div aria-hidden>` in `#hui-portals` | `data-open` |
| `ModalContent` | `<div role="dialog">` (`alertdialog`), `aria-modal` | `data-open` |
| `ModalTitle`, `ModalDescription` | `<h2>`, `<p>` | |
| `ModalClose` | `<button>` | `data-enabled` |
| `Dropdown` | `<div>` | `data-open` |
| `DropdownTrigger` | `<button aria-haspopup="menu">` | `aria-expanded`, `data-open`, button state |
| `DropdownContent` | `<div role="menu">` in `#hui-portals` | `data-open` |
| `DropdownItem` | `<button role="menuitem">` or `<a>` | button state, highlight is `focus:` |
| `DropdownCheckboxItem`, `DropdownRadioItem` | `<button role="menuitemcheckbox">`, `menuitemradio` | `aria-checked`, `data-checked`, button state |
| `DropdownRadioGroup`, `DropdownGroup` | `<div role="group">` | `data-disabled` on the radio group |
| `DropdownLabel`, `DropdownSeparator` | `<div>`, `<div role="separator">` | |
| `Tooltip` | `<div>` | `data-open` |
| `TooltipContent` | `<div role="tooltip">` in `#hui-portals` | `data-open` |
| `TooltipClose` | `<button>` | `data-enabled` |

Any element that opens an anchored overlay through a control constant (`DropdownToggle`,
`TooltipShow`, …) becomes its anchor, and a dropdown opener also gets `aria-expanded` and
`data-open`.

### Content

| Part | Element | Attributes |
|---|---|---|
| `Button` | `<button>` or `<a href>` | button state; `Pending` adds `data-pending` and `aria-busy`; a disabled link renders without `href`, with `aria-disabled` |
| slots | `<span aria-hidden data-slot="left\|right">` | on `Button`, `SelectTrigger`, `SelectItem`, `ListboxTrigger`, list items and dropdown items |
| `Tabs`, `TabsList` | `<div>`, `<div role="tablist">` | `data-orientation` |
| `TabsTrigger` | `<button role="tab">` | `aria-selected`, `data-selected`, button state |
| `TabsContent` | `<div role="tabpanel">` | `data-selected` |
| `Collapsible`, `CollapsibleHeader`, `CollapsibleTitle` | `<div>`, `<div>`, `<span>` | `data-open`; `data-clickable` on the header |
| `CollapsibleTrigger` | `<button>` | `aria-expanded`, `data-open`, button state |
| `CollapsibleContent` | `<div>` | `data-open` |
| `Breadcrumb`, `BreadcrumbList`, `BreadcrumbItem` | `<nav>`, `<ol>`, `<li>` | |
| `BreadcrumbLink` | `<a>`, or `<span>` without `Href` | `aria-current="page"`, `data-current` |
| `BreadcrumbSeparator` | `<li aria-hidden>` | |
| `Table` | `<table>` | `aria-rowcount` |
| `TableHeader`, `TableBody`, `TableFooter` | `<thead>`, `<tbody>`, `<tfoot>` | `data-pending` while they are an htmx target |
| `TableRow` | `<tr>` | `data-selected`, `aria-rowindex` |
| `TableHead` | `<th>` | `aria-sort`, `data-sort` |
| `TableCell` | `<td>` | `data-label` (for stacked mobile rows) |
| `TableLoader` | `<tr aria-hidden data-hui-loader>` | |
| `Meter` | `<div role="meter">` | `data-level` (`low`, `medium`, `high`), `data-state` (`optimum`, `suboptimal`, `critical`, only with `Optimum`), CSS variable `--hui-meter` |
| `Toaster` | `<ol id="hui-toaster">` | |
| `Toast` | `<li role="status">` (`alert` when urgent) | `data-open`, `data-paused`, CSS variable `--hui-toast-duration` |
| `ToastTitle`, `ToastDescription` | `<div>` | |
| `ToastAction`, `ToastClose` | `<button>` | `data-enabled` |
| `Announcer` | `<div role="status">` (`alert` when urgent) | |
| `Alert` | `<div role="alert">` (`status` when polite) | |
| `SkipLink` | `<a href="#main">` | |

## What it looks like

```templ
import hui "github.com/daarxwalker/hui"

@hui.Field(hui.FieldProps{
    Name:        "email",
    Value:       user.Email,
    Description: "We will send a confirmation there.",
    Error:       errs["email"],
    Required:    true,
    Class:       "grid gap-1.5",
}) {
    @hui.Label(hui.LabelProps{Class: "text-sm font-medium data-invalid:text-red-600"}) {
        Email
    }
    @hui.TextInput(hui.TextInputProps{
        Type:        "email",
        Class:       "rounded-md border px-3 py-2 data-invalid:border-red-500",
        Placeholder: "you@example.com",
        Attrs: templ.Attributes{
            "hx-post":   "/validate/email",
            "hx-target": "#email-field",
            "hx-swap":   "outerHTML",
        },
    })
    @hui.Description(hui.DescriptionProps{Class: "text-xs text-gray-500"})
    @hui.Error(hui.ErrorProps{Class: "text-xs text-red-600"})
}
```

You get a fully wired label, input, description and error. The `for`, `id`,
`aria-describedby` and `aria-invalid` attributes are worked out by the components, and
no JavaScript is involved:

```html
<div class="grid gap-1.5" data-field="email" data-invalid data-required id="email-field">
  <label class="text-sm font-medium data-invalid:text-red-600" data-invalid data-required
         for="email">Email</label>
  <input class="rounded-md border px-3 py-2 data-invalid:border-red-500"
         aria-describedby="email-description email-error" aria-invalid="true" required
         data-invalid data-required id="email" name="email" type="email" value="petr"
         placeholder="you@example.com"
         hx-post="/validate/email" hx-swap="outerHTML" hx-target="#email-field">
  <p class="text-xs text-gray-500" data-invalid data-required
     id="email-description">We will send a confirmation there.</p>
  <p class="text-xs text-red-600" data-invalid data-required id="email-error"
     role="alert">Not a valid email.</p>
</div>
```

State is exposed through `data-*` attributes, never through CSS classes of our own.

## Rules

- One flat package. Every component is a pair of files: `component.go` (props, attributes)
  and `component.templ` (thin markup that only calls the attrs function).
- Every component has its own `<Component>Props` struct embedding `Base`, so a literal
  reads `TextInputProps{ID: "email"}` instead of `TextInputProps{Base: Base{ID: "email"}}`
  (Go 1.27+).
- `Base` names the three things every component must support: `Class` (styling), `Attrs`
  (the escape hatch for anything at all — `hx-*`, `x-model`, `@keyup.enter`) and `Events`
  (named Alpine handlers).
- `Base.merge` assembles the final attributes: the caller's `Attrs` first, then the
  component's own attributes on top (those win, so the accessibility wiring cannot be
  overwritten by accident), then `ID`.
- Nothing the caller passes is silently dropped. A `class` in `Attrs` is **merged** with
  `Class`, a `style` in `Attrs` is appended after the component's own inline style (so
  yours wins in the cascade), and a handler in `Attrs` is chained with the matching
  `Events` entry. Named
  `Events` are appended to the internal handler of the same kind rather than replacing it —
  `@blur="focused = false; $store.form.validate('email')"`. The order matters: internal
  state has to settle before your expression runs.

## Field

Form components share their state through `Field` so that the label and the input can
never tell two different stories. `Field` puts `FieldProps` into the context and `Label`,
`Description`, `Error`, `TextInput`, `Textarea`, `Checkbox`, `Switch`, `Radio`, `Range`,
`Select`, `Listbox`, `Combobox`, `FileUploadInput` and `Meter` read it back.

- `Invalid` is true whenever `Error` is set — the two cannot drift apart.
- `ID` falls back to `Name`, and is generated when both are missing.
- The wrapper gets `id="<id>-field"` and `data-field="<name>"`, so one field is addressable
  on its own without hand-written attributes. The control keeps the plain `id`.
- `Description` and `Error` take their text from `FieldProps`, not from children. The input
  renders before the elements below it, so in a single pass it could not know whether a
  description or an error exists, and `aria-describedby` would point at IDs that are never
  rendered.
- Inside a `Field`, the `Field` is the single source of truth: `Name`, `Value`, `Disabled`,
  `Readonly`, `Required` and `Invalid` set on the input itself are ignored.
- The label gets `id="<id>-label"`, so controls that `for` cannot label (`Meter`, the
  list of a `Listbox`) point at it with `aria-labelledby`.
- `Form` is rendered as the native `form` attribute on the control (see
  [Outside a form](#outside-a-form)). The field's `Form` wins; without one the control
  keeps its own.

Every component also works on its own, without a `Field`. `Label` then takes `for` from its
own props, and `Description` and `Error` render their children. The price is that wiring up
`for`/`id`/`aria-describedby` is left to you (through `ID` and `Attrs`) — a deliberate
trade-off, not a bug.

## Fieldset and Radio

`Fieldset` is to a group what `Field` is to a single control. It renders a native
`<fieldset>`, holds the shared state and hands it to everything inside. It covers two
cases:

- **One name, several options** — a radio group, a group of checkboxes or switches
  (`plan=pro`, `tags=go&tags=htmx`).
- **Several fields that belong together** — an address with a street and a city, each with
  its own name.

```templ
@hui.Fieldset(hui.FieldsetProps{
    Name: "plan", Value: user.Plan, Required: true,
    Description: "You can change it later.", Error: errs["plan"],
    Class: "grid gap-2",
}) {
    @hui.Legend(hui.LegendProps{}) {
        Plan
    }
    @hui.Description(hui.DescriptionProps{})
    for _, plan := range plans {
        @hui.Field(hui.FieldProps{Value: plan.ID, Description: plan.Price, Disabled: plan.SoldOut,
            Class: "flex gap-2 has-data-checked:border-blue-600"}) {
            @hui.Radio(hui.RadioProps{Class: "size-4 rounded-full border data-checked:border-blue-600"}) {
                @hui.RadioIndicator(hui.RadioIndicatorProps{Class: "m-auto size-2 rounded-full bg-blue-600"})
            }
            @hui.Label(hui.LabelProps{}) {
                { plan.Name }
            }
            @hui.Description(hui.DescriptionProps{})
        }
    }
    @hui.Error(hui.ErrorProps{})
}
```

```html
<fieldset id="plan-field" data-field="plan" aria-describedby="plan-description plan-error" data-invalid data-required>
  <legend id="plan-legend" data-invalid data-required>Plan</legend>
  <p id="plan-description">You can change it later.</p>
  <div id="plan-pro-field" data-option="pro" data-invalid data-required>
    <button type="button" role="radio" id="plan-pro" name="plan" value="pro" aria-checked="true"
            data-checked data-hui-radio="plan" aria-invalid="true" aria-describedby="plan-pro-description"
            data-invalid data-required x-data="huiRadio({ … })">
      <input type="hidden" name="plan" value="pro" data-hui-checked>
      <span aria-hidden="true" x-show="checked" x-transition></span>
    </button>
    <input data-hui-require="plan-pro" required value="1" tabindex="-1" aria-hidden="true" style="…">
    <label for="plan-pro" id="plan-pro-label">Pro</label>
    <p id="plan-pro-description">$10 / month</p>
  </div>
  …
  <p id="plan-error" role="alert">Choose a plan</p>
</fieldset>
```

**IDs.** `ID` falls back to `Name` and is generated when both are missing. The fieldset gets
`<id>-field` (the same swap target as a `Field`), the legend `<id>-legend`, and its own
description and error `<id>-description` and `<id>-error`.

**A `Field` without a `Name` is an option of the group.** Its `Value` is the option's value.

| | Comes from |
|---|---|
| `Name` | the fieldset |
| `ID` | `<fieldset-id>-<value>` (characters outside `[A-Za-z0-9_-]` become `-`), unless you set one |
| `Disabled` | the fieldset **or** the option, so one option can be disabled on its own |
| `Invalid` | the fieldset or the option's own `Error` |
| `Required` | the fieldset. A radio renders a validation input for it (one checked radio satisfies it); a checkbox or switch does not, because natively every box would be required — check "at least one" on the server |
| checked | `Value == fieldset.Value` or `Value` in `fieldset.Values`; it wins over `Checked` on the input |
| wrapper | `data-option="<value>"` instead of `data-field`, so it is not mistaken for the group's swap target |
| `aria-describedby` | only the option's own description and error. The group's description is read with the `<fieldset>` and is not repeated for every option |

The option wrapper has no state of its own. Style the row from the control inside it with
`has-data-checked:`; the plugin keeps `data-checked` on the control current.

**A `Field` with a `Name` is an independent field inside the group.** It works exactly like
a `Field` anywhere else and only takes `Disabled` from the fieldset. The group's `Invalid`
and `Required` do not reach it — each field has its own `Error`.

```templ
@hui.Fieldset(hui.FieldsetProps{ID: "address", Disabled: !editable}) {
    @hui.Legend(hui.LegendProps{}) {
        Address
    }
    @hui.Field(hui.FieldProps{Name: "street", Value: a.Street, Error: errs["street"]}) {
        ...
    }
    @hui.Field(hui.FieldProps{Name: "city", Value: a.City, Error: errs["city"]}) {
        ...
    }
}
```

**Without a `Field`.** `Radio`, `Checkbox` and `Switch` placed straight in a fieldset (for
example wrapped in their own `<label>`) join the group by their `Value`, unless they have a
`Name` of their own. `TextInput` and `Textarea` take only `Disabled` and `Form` from it.

`Disabled` is also rendered as the native `disabled` attribute on the `<fieldset>`, which
disables every control inside even without JavaScript; the controls still get their own
`disabled` and `data-disabled` for styling. `Description` and `Error` belong to the closest
owner — a `Field` first, then a `Fieldset` — so the ones above the options describe the
group and the ones inside a `Field` describe the option. `Legend` outside a fieldset just
renders its children.

`Radio` works on its own too, with `Name`, `Value`, `Checked`, `Form`, `Disabled`,
`Required` and `Invalid`. It is described with `Checkbox` below.

## Textarea

`Textarea` is `TextInput` for multi-line text: the same `Field` wiring, the same props
(`Placeholder`, `Autocomplete`, `Disabled`, `Readonly`, `Required`, `Invalid`, `Form`) plus
`Rows`, and the same optional `Alpine: true` with `data-focused` and `data-empty`. The value
is rendered as the element's text. For a box that grows with its content, use CSS:
`field-sizing-content`.

## Checkbox, Switch and Radio

All three are a `<button type="button">` with the right role (`checkbox`, `switch`,
`radio`), so the indicator is your own markup instead of a styled native control. The
state lives in a small Alpine plugin, and each one renders its initial state from Go.

```templ
@hui.Field(hui.FieldProps{Name: "terms", Value: "yes", Required: true, Class: "flex gap-2"}) {
    @hui.Checkbox(hui.CheckboxProps{Checked: accepted,
        Class: "flex size-4 items-center justify-center rounded border data-checked:bg-blue-600 data-indeterminate:bg-blue-600"}) {
        @hui.CheckboxIndicator(hui.CheckboxIndicatorProps{}) {
            @icon.Check()
        }
        @hui.CheckboxIndicator(hui.CheckboxIndicatorProps{Indeterminate: true}) {
            @icon.Minus()
        }
    }
    @hui.Label(hui.LabelProps{}) {
        I accept the terms
    }
}

@hui.Switch(hui.SwitchProps{Name: "notifications", Checked: user.Notifications,
    Class: "inline-flex h-6 w-11 items-center rounded-full bg-gray-300 data-checked:bg-blue-600",
    Attrs: templ.Attributes{"hx-post": "/settings/notifications"}}) {
    @hui.SwitchThumb(hui.SwitchThumbProps{Class: "size-5 translate-x-0.5 rounded-full bg-white transition data-checked:translate-x-5.5"})
}
```

| Part | Renders | Notes |
|---|---|---|
| `Checkbox` | `<button role="checkbox">` | `aria-checked` (`true`, `false`, `mixed`), `data-checked`, `data-indeterminate`, `x-modelable="checked"` |
| `CheckboxIndicator` | `<span aria-hidden="true">` | shown while checked, or while indeterminate with `Indeterminate: true` |
| `Switch` | `<button role="switch">` | `aria-checked`, `data-checked`, Space and Enter toggle |
| `SwitchThumb` | `<span aria-hidden="true">` | always visible, carries `data-checked` for the slide |
| `Radio` | `<button role="radio">` | `aria-checked`, `data-checked`, `data-hui-radio="<name>"` |
| `RadioIndicator` | `<span aria-hidden="true">` | shown while checked |

- **Submitting.** A checked control holds `<input type="hidden" name value>` inside the
  button, so a `<form>` (or `hx-include` of a container) submits it exactly like a native
  checkbox. An unchecked one submits nothing. A checked button also carries `name` and
  `value` itself, because htmx sends a button's own value and never its children: an
  `hx-post` with `hx-trigger="change"` on the control sends `name=value` while it is
  checked. A `type="button"` is never submitted by a form, so nothing is sent twice.
  `Value` defaults to `on`.
- **Events.** Every change dispatches a bubbling `change` from the button, so
  `Events.OnChange`, `@change` and `hx-trigger="change"` work as with native inputs.
  `checked` is `x-modelable`.
- **Labels.** A `<label for>` clicks the button, so `Label` inside a `Field` works as
  usual, and so do `Description`, `Error` and `aria-describedby`.
- **Required.** HTML cannot validate a button, so `Required` renders a hidden validation
  input right after it (`data-hui-require`). The browser reports it when the box is
  unchecked, and focusing it moves focus to the button. A checkbox and a switch also get
  `aria-required`.
- **Keyboard.** Space toggles a checkbox (Enter does not, as for a native one). A switch
  toggles on Space and Enter. Radios with the same `Name` (and the same form) form a group:
  only the checked one is in the tab order, the arrow keys move to the next enabled radio
  and check it.
- **Reset.** A form `reset` restores the state rendered by the server.
- **Indicators** use `x-show` with the default `x-transition` (see [Overlays](#overlays));
  pass `Attrs: {"x-transition": false}` or your own `x-transition.*` to change it.
- Inside a `Field` the field's `Name`, `Value`, `Disabled`, `Required` and `Invalid` win,
  like for `TextInput`. `Checked` and `Indeterminate` stay on the control. HTML has no
  read-only checkbox, so a `Readonly` field never reaches it.
- **Without JavaScript** the controls render their state but cannot be toggled. This is
  the price of a custom indicator.

## Range

`Range` is a slider with one or more thumbs. Every thumb is a native
`<input type="range">`, so the arrow keys, Home/End, Page Up/Down, form submission and the
screen reader announcement all come from the browser. The thumbs are stacked over one
track, which turns two inputs into a single "from–to" slider.

```templ
@hui.Field(hui.FieldProps{Name: "price", Class: "grid gap-2"}) {
    @hui.Label(hui.LabelProps{}) {
        Price
    }
    @hui.Range(hui.RangeProps{
        Min: 0, Max: 1000, Step: 10,
        Value: []float64{filter.PriceFrom, filter.PriceTo},
        Class: "relative h-5",
    }) {
        <div class="absolute inset-x-0 top-2 h-1 rounded bg-gray-200"></div>
        <div class="absolute top-2 h-1 rounded bg-blue-600 left-(--hui-range-start) right-[calc(100%-var(--hui-range-end))]"></div>
        @hui.RangeThumb(hui.RangeThumbProps{Index: 0, Class: thumb, Attrs: templ.Attributes{"aria-label": "Minimum price"}})
        @hui.RangeThumb(hui.RangeThumbProps{Index: 1, Class: thumb, Attrs: templ.Attributes{"aria-label": "Maximum price"}})
    }
}
```

```go
const thumb = "pointer-events-none absolute inset-0 w-full appearance-none bg-transparent " +
    "[&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:size-5 " +
    "[&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full " +
    "[&::-webkit-slider-thumb]:bg-white [&::-webkit-slider-thumb]:shadow " +
    "[&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:size-5 " +
    "[&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-white"
```

- `Range` renders the root `<div>` and holds the shared state. `RangeThumb` renders one
  input and picks its value by `Index` from `Value`. Render one thumb per value; a thumb
  without a value starts at `Min`.
- `Min` and `Max` default to 0 and 100, `Step` to 1. Decimal steps work (`Step: 0.1`).
  Values from the server are clamped to `Min`–`Max` and sorted, so the thumbs never start
  crossed.
- The stacked inputs need `pointer-events-none` on the input and `pointer-events-auto` on
  the thumb pseudo-element, otherwise the top input swallows every click. This is the one
  piece of CSS a multi-thumb slider cannot do without, and it is left to you on purpose.
- The selected segment is exposed as two CSS variables on the root, `--hui-range-start`
  and `--hui-range-end`, both measured from the `Min` side in percent. With one thumb the
  segment runs from `Min` to the thumb, with more it runs from the first thumb to the last.
  The initial values are rendered from Go, so the fill is right before Alpine boots. A
  `style` you pass in `Attrs` is appended, not replaced.
- Every thumb submits under the same name, in order: `price=200&price=800`. In Go that is
  `r.Form["price"]`.
- Inside a `Field`, the field's `Name`, `Disabled` and `Invalid` win. The first thumb gets
  the field's id, so the `Label` points at it, and the others get `<id>-1`, `<id>-2`, …
  Every thumb carries the field's `aria-describedby`. `Required` and `Readonly` do not
  exist for a range input, so they stay on the wrapper. `Value` is read from `RangeProps`,
  because a field value is a single string.
- With more than one thumb, give each one its own `aria-label` ("Minimum price", "Maximum
  price"). A label can only point at one input, and a screen reader would otherwise
  announce two sliders with the same name.

**Alpine.** Range has a small Alpine plugin (`huiRange`, built from
[`alpine/src/range.js`](alpine/src/range.js)), loaded with the rest, see
[Alpine](#alpine). It does three things:

- A thumb cannot be moved past its neighbours, whether by mouse, touch or keyboard, so the
  submitted values stay in order.
- It keeps `--hui-range-start` and `--hui-range-end` in sync, also after a form reset.
- When thumbs sit on top of each other, it sets an inline `z-index` so the one that can
  still move is on top: at the upper end of the track that is the lower thumb.

The current values are reactive as `values`, so anything inside the root can show them or
derive from them:

```html
<output x-text="values.join(' – ')"></output>
```

Each thumb's own `@input` runs first, so anything you add through `Events.OnInput` sees the
value after it was clamped. Without Alpine the thumbs still work and submit, but they can
cross and the fill stays where the server left it.

**htmx.** The inputs fire `change` when you let go, and the event bubbles to the root, so a
filter that refreshes on release looks like this:

```go
Attrs: templ.Attributes{
    "hx-get":     "/products",
    "hx-trigger": "change",
    "hx-include": "this",
    "hx-target":  "#results",
}
```

## Select

`Select` is the native `<select>` with a custom trigger, in the shape of Headless UI's
Select: a [customizable select](https://developer.chrome.com/blog/a-customizable-select)
whose first child is a `<button>` with `<selectedcontent>`. The browser keeps the picker,
the keyboard, the form value and the validation; no JavaScript is involved.

```templ
@hui.Field(hui.FieldProps{Name: "fruit", Value: order.Fruit, Required: true}) {
    @hui.Label(hui.LabelProps{}) {
        Fruit
    }
    @hui.Select(hui.SelectProps{Placeholder: "Pick a fruit",
        Class: "[appearance:base-select] [&::picker(select)]:[appearance:base-select] rounded border"}) {
        @hui.SelectTrigger(hui.SelectTriggerProps{RightSlot: chevron()}) {
            @hui.SelectValue(hui.SelectValueProps{})
        }
        @hui.SelectGroup(hui.SelectGroupProps{Label: "Fruit"}) {
            for _, f := range fruits {
                @hui.SelectItem(hui.SelectItemProps{Value: f.Code, LeftSlot: f.Icon,
                    Class: "checked:font-medium hover:bg-gray-100"}) {
                    { f.Name }
                }
            }
        }
    }
}
```

| Part | Renders | Wiring |
|---|---|---|
| `Select` | `<select>` | `name`, `form`, `disabled`, `required`, `aria-invalid`, `aria-describedby` and the field's id and state |
| `SelectTrigger` | `<button>` | `data-enabled`/`data-disabled`, `LeftSlot`/`RightSlot`; followed by the placeholder option |
| `SelectValue` | `<selectedcontent>` | the browser copies the selected option into it |
| `SelectItem` | `<option>` | `value`, `selected` for the current value, `disabled`, `LeftSlot`/`RightSlot` |
| `SelectGroup` | `<optgroup>` | `label`, `disabled` |
| `SelectSeparator` | `<hr>` | |

- **The styling is yours.** Opting into the custom picker is CSS, so the core renders no
  style: `appearance: base-select` on the `<select>` and on `::picker(select)`. Style the
  picker through `[&::picker(select)]:…`, hide the native arrow with
  `[&::picker-icon]:hidden` and the checkmark of an option with `[&::checkmark]:hidden`.
  The open state is the `:open` pseudo-class (`open:`), the chosen option `:checked`
  (`checked:`).
- A browser without customizable selects ignores the `<button>` and shows a plain
  `<select>` with the same options, value and validation.
- `Placeholder` renders `<option value="" data-placeholder>` right after the trigger,
  chosen while there is no value. With `Required` the browser rejects it. Hide it from the
  picker with `[&_option[data-placeholder]]:hidden` and dim the trigger with
  `has-[option[data-placeholder]:checked]:text-gray-400`; no stale `data-placeholder`
  state is rendered on the select itself.
- `<selectedcontent>` copies the whole option, slots included: `*:data-[slot=right]:hidden`
  on `SelectValue` drops an icon that only belongs in the list.
- A select has one value. `Multiple`, a search or options loaded from the server are
  [`Listbox`](#listbox).

## Listbox

`Listbox` is a custom select in the shape of Headless UI's Listbox: a button that opens a
list of options. The value is kept in hidden inputs, so it submits with a form like a
native `<select>` would, and `Multiple: true` turns it into a multi-select. Unlike
`Select`, the list can be searched and loaded from the server.

```templ
@hui.Field(hui.FieldProps{Name: "country", Value: user.Country, Required: true}) {
    @hui.Label(hui.LabelProps{}) {
        Country
    }
    @hui.Listbox(hui.ListboxProps{Placeholder: "Pick a country"}) {
        @hui.ListboxTrigger(hui.ListboxTriggerProps{Class: "btn data-placeholder:text-gray-400"}) {
            @hui.ListboxValue(hui.ListboxValueProps{}) {
                { countryName(user.Country) }
            }
        }
        @hui.ListboxContent(hui.ListboxContentProps{Placement: hui.Placement{SideOffset: 4},
            Class: "z-50 rounded border bg-white p-1"}) {
            @hui.ListboxList(hui.ListboxListProps{Class: "max-h-60 overflow-auto"}) {
                for _, c := range countries {
                    @hui.ListboxItem(hui.ListboxItemProps{Value: c.Code, Label: c.Name,
                        Class: "data-highlighted:bg-gray-100 data-selected:font-medium"}) {
                        { c.Name }
                        @hui.ListboxItemIndicator(hui.ListboxItemIndicatorProps{}) {
                            ✓
                        }
                    }
                }
            }
        }
    }
}
```

| Part | Renders | Wiring |
|---|---|---|
| `Listbox` | `<div>` + hidden inputs | `x-data`, `x-modelable="value"`, `data-open`, one `<input type="hidden">` per value |
| `ListboxTrigger` | `<button role="combobox" aria-haspopup="listbox">` | `aria-controls` (the list), `aria-expanded`, `aria-activedescendant`, `data-placeholder`, the field's id and `aria-describedby`, anchor for the content. With `Searchable` a plain `<button aria-haspopup="listbox">` |
| `ListboxValue` | `<span>` | its children are the label of the current value, or `Placeholder` when there is none |
| `ListboxContent` | `<template x-teleport>` with `<div>` | the positioned popup (`<id>-content`), at least as wide as the trigger |
| `ListboxSearch` | `<input role="combobox">` | only with `Searchable`: `aria-controls` (the list), `aria-autocomplete="list"`, `aria-activedescendant`, `QueryName`, `Placeholder` |
| `ListboxList` | `<div role="listbox">` | `<id>-list`, `aria-labelledby` (the field label or the trigger), `aria-multiselectable`; the target the server redraws |
| `ListboxEmpty` | `<div role="presentation">` | for "nothing found" |
| `ListboxLoader` | `<div aria-hidden="true" data-hui-loader>` | the sentinel of an infinite list, see [Server-rendered lists](#server-rendered-lists) |
| `ListboxItem` | `<div role="option">` | `aria-selected`, `data-selected`, `data-highlighted`, `data-disabled` |
| `ListboxItemIndicator` | `<span aria-hidden="true">` | visible only while its item is selected |
| `ListboxGroup`, `ListboxGroupLabel` | `<div role="group">`, `<div>` | the group is labelled by its label |
| `ListboxSeparator` | `<div aria-hidden="true">` | |

- `Value` is a single value, `Values` the list for `Multiple`. Inside a `Field` the field's
  `Name` and `Value` win, `Values` stay on the listbox. The root carries the field's state
  (`data-invalid`, …), so does the trigger.
- `ListboxTrigger` and `ListboxItem` (and `ComboboxItem`) take `LeftSlot` and `RightSlot`
  like `Button`, rendered in `<span aria-hidden="true" data-slot>`. A flag on the left or a
  chevron on the right never reaches the label of the value, the typeahead or the
  accessible name: the plugin reads an item's text without its `aria-hidden` parts.
- Go cannot see the items while it renders the trigger, so the label of the current value
  is yours to pass as the children of `ListboxValue`. After a change the plugin replaces it
  with the item's `Label` (or the item's text without its `aria-hidden` parts).
- The trigger keeps focus the whole time; the highlighted option is announced through
  `aria-activedescendant`. Keyboard follows the
  [select-only combobox](https://www.w3.org/WAI/ARIA/apg/patterns/combobox/examples/combobox-select-only/)
  pattern: ArrowDown, ArrowUp, Enter and Space open it on the selected option, Home/End
  open it on the first/last one, typing jumps to a matching option, Enter and Space choose,
  Escape and Tab close. With `Multiple` the list stays open after a choice.
- Hidden inputs cannot be `required`, so `Required` adds a visually hidden validator input
  without a name. The browser blocks the submit and points at it, and the plugin hands the
  focus to the trigger.
- A change fires a bubbling `change` on the hidden input and `hui-listbox:change` with
  `{ id, value }` on the root (`OnChange`). `x-model` binds `value` — a string, or an array
  with `Multiple`. `OnOpen`/`OnClose` and the `hui-listbox:show|hide|toggle` events work like
  for `Modal`.
- Without JavaScript nothing opens; the hidden inputs still submit the server's value.

**Search.** `Searchable: true` puts a `ListboxSearch` in the popup, in the shape of a
popover with a command list. The server filters:

```templ
@hui.Listbox(hui.ListboxProps{ID: "country", Name: "country", Searchable: true}) {
    @hui.ListboxTrigger(hui.ListboxTriggerProps{}) { ... }
    @hui.ListboxContent(hui.ListboxContentProps{}) {
        @hui.ListboxSearch(hui.ListboxSearchProps{QueryName: "q", Placeholder: "Search…", Attrs: templ.Attributes{
            "hx-get": "/countries", "hx-trigger": "input changed delay:200ms",
            "hx-target": "#country-list", "hx-include": "#country"}})
        @hui.ListboxList(hui.ListboxListProps{}) {
            @countryItems(countries, 1)
        }
    }
}
```

- Opening moves focus into the search; it holds `aria-activedescendant` and takes the
  arrows, Enter and Tab. Escape, a choice or Tab close the popup and focus goes back to the
  trigger. A printable key on the closed trigger opens it and types into the search.
- The search is teleported out of any `<form>` and has no `form` attribute, so the query is
  never submitted with the value. `hx-include` of the root sends the current value along,
  so the server can render the selected item (see `WithListbox`).
- Only the list is swapped, so the search keeps focus and caret while the server redraws.
- What was typed is marked in the options, see [Match highlighting](#match-highlighting).

## Combobox

`Combobox` picks a value out of more options than fit in a `Listbox` — a user, a tag, a
city — in the shape of Headless UI's Combobox: the input is the field itself, you type
straight into it and the popup holds only the list. It shares the listbox with `Listbox`,
and the filtering is **the server's job**: the input asks htmx for new options, the plugin
only keeps the highlight, the keyboard, the query and the value.

```templ
@hui.Field(hui.FieldProps{Name: "country", Value: user.Country}) {
    @hui.Label(hui.LabelProps{}) {
        Country
    }
    @hui.Combobox(hui.ComboboxProps{ID: "country"}) {
        <div class="flex rounded border has-data-open:border-blue-500">
            @hui.ComboboxInput(hui.ComboboxInputProps{
                Value:       countryName(user.Country),
                QueryName:   "q",
                Placeholder: "Pick a country",
                Attrs: templ.Attributes{
                    "hx-get":     "/countries",
                    "hx-trigger": "hui-combobox:query delay:200ms",
                    "hx-target":  "#country-list",
                },
            })
            @hui.ComboboxTrigger(hui.ComboboxTriggerProps{}) {
                @chevron()
            }
        </div>
        @hui.ComboboxContent(hui.ComboboxContentProps{Class: "z-50 rounded border bg-white p-1"}) {
            @hui.ComboboxList(hui.ComboboxListProps{Class: "max-h-60 overflow-auto"}) {
                @countryItems(countries)
            }
        }
    }
}
```

```templ
templ countryItems(countries []Country) {
    for _, c := range countries {
        @hui.ComboboxItem(hui.ComboboxItemProps{Value: c.Code, Label: c.Name}) {
            <span>{ c.Name }</span>
        }
    }
    if len(countries) == 0 {
        @hui.ComboboxEmpty(hui.ComboboxEmptyProps{}) {
            Nothing found.
        }
    }
}
```

| Part | Renders | Wiring |
|---|---|---|
| `Combobox` | `<div>` + hidden inputs | `x-data`, `x-modelable="value"`, `data-open`, one `<input type="hidden">` per value, anchor for the content |
| `ComboboxInput` | `<input role="combobox">` | the field's id (`<id>-input` without one), `aria-controls` (the list), `aria-expanded`, `aria-autocomplete="list"`, `aria-activedescendant`, `aria-describedby`, `aria-invalid`, `aria-required`, `data-open`; `Value` (the label of the current value), `Placeholder`, `QueryName` |
| `ComboboxTrigger` | `<button type="button" tabindex="-1" aria-haspopup="listbox">` | optional: `<id>-trigger`, `aria-controls`, `aria-expanded`, `data-open`; toggles the list and keeps the focus in the input |
| `ComboboxContent` | `<template x-teleport>` with `<div>` | the positioned popup (`<id>-content`), at least as wide as the root |
| `ComboboxList` | `<div role="listbox">` | `<id>-list`, `aria-labelledby` (the field label or the input), `aria-multiselectable`; the target the server redraws |
| `ComboboxItem`, `ComboboxItemIndicator`, `ComboboxEmpty`, `ComboboxLoader`, `ComboboxGroup`, `ComboboxGroupLabel`, `ComboboxSeparator` | | as in `Listbox` |
| `ComboboxRemove` | `<button type="button" tabindex="-1">` | removes `Value`; for chips |

- **Focus stays in the input.** Typing opens the list and asks the server, ArrowDown opens
  it on the selected option, ArrowUp on the last one, the arrows move the highlight, Enter
  chooses, Escape and Tab close. A click into the input opens the list and selects its
  text, so typing replaces the current label instead of appending to it.
- **The query is not the label.** The input shows the label of the current value, which
  is never a query. The plugin tracks what was typed since the list opened or a choice was
  made, and dispatches `hui-combobox:query` on the input whenever that changes. Trigger the
  request on it (`hx-trigger="hui-combobox:query delay:200ms"`), not on `input`.
- The input has no `name`. The plugin adds the query to the requests of the input as
  `QueryName`, so the label is never sent as a query and nothing is submitted with a form.
  `hx-include` of the root sends the current value along (see `WithCombobox`).
- **Closing restores the label.** Escape, Tab or a click elsewhere puts the label of the
  current value back into the input, and when the list was filtered the plugin asks for
  the unfiltered one with an empty query. Clearing the input and closing it clears the
  value.
- Give the input an accessible name: a `Field` with a `Label` labels it natively,
  otherwise pass `aria-label` through `Attrs`. The trigger is out of the tab order.
- The `/countries` handler renders only `countryItems` — items do not need a `Combobox`
  around them, see [Server-rendered lists](#server-rendered-lists).
- From the first keystroke until the server answers that query, the input and the list
  carry `data-searching`. It starts before the debounced request is even sent and ends
  when the input holds the query that was answered last, so typing back to it clears it
  at once. Unlike `data-pending`, loading the next page never sets it, which makes it the
  hook for a spinner in the field (`peer-data-searching:`) and a list that cannot be used
  meanwhile (`data-searching:pointer-events-none`). Enter does not pick a stale option
  while searching. It needs an `hx-get` (or another `hx-*` verb) on the input; the same
  holds for the search of a `Searchable` `Listbox`.
- With `Multiple` the input only holds the query: after a choice it is emptied and the
  list stays open. Chips are yours to render inside the root, next to the input — for
  example `<template x-for="value in selected()">` with `labels[value]` as the text and
  `remove(value)` on a button, or server-rendered chips with `ComboboxRemove{Value}`.
- Hidden inputs, `Required`, `x-model` and the events work like in `Listbox`.

### Match highlighting

While a `Searchable` `Listbox` or a `Combobox` is searched, the plugin wraps every match of
the query in the options in `<mark data-hui-match>`, so `Ho` in the search turns
`Holení` into `<mark data-hui-match>Ho</mark>lení`.

- Matching ignores case and diacritics (`ces` marks `Česko`) and marks every occurrence.
  Parts with `aria-hidden` (slots, indicators) are skipped.
- The query of a `Combobox` is only what was typed. The label of the current value shown
  in the input is never marked, so opening the list on `Czechia` marks nothing until you
  type.
- The marks follow every redraw and appended page and disappear with an empty query; a
  `Combobox` drops its query when the list closes. The value, the label and the typeahead
  read the text without them.
- The core does not style them, `<mark>` keeps the browser's yellow. Style it through the
  item, for example `**:data-hui-match:bg-yellow-100 **:data-hui-match:font-semibold
  **:data-hui-match:text-inherit`. A mark splits the text, so in a `flex` item with a
  `gap` wrap the label in its own element (`<span>{ c.Name }</span>`), or the pieces of
  the word become separate flex items.

## Server-rendered lists

`Listbox`, `Combobox` and `Table` are built for a server that redraws and appends: a search
replaces the list, an infinite scroll appends the next page, a whole field comes back from
`hx-swap="outerHTML"` or `hx-swap-oob`. The result is the same as if everything had been
rendered at once.

**Infinite scroll** is the plain htmx sentinel. The last element of a page is a loader that
replaces itself with the next page:

```templ
templ countryItems(countries []Country, page int) {
    for _, c := range countries {
        @hui.ListboxItem(hui.ListboxItemProps{Value: c.Code, Label: c.Name}) {
            { c.Name }
        }
    }
    if hasMore {
        @hui.ListboxLoader(hui.ListboxLoaderProps{Attrs: templ.Attributes{
            "hx-get": fmt.Sprintf("/countries?page=%d", page+1),
            "hx-trigger": "intersect once", "hx-swap": "outerHTML", "hx-include": "#country"}}) {
            @spinner()
            Loading…
        }
    }
}
```

- **Partials render the same HTML.** `hui.WithListbox(ctx, props)` and
  `hui.WithCombobox(ctx, props)` return the context of a listbox or combobox, so items
  rendered by a handler get the same ids, `aria-selected`, `data-selected` and indicators
  as in the full page. Pass at least the `ID` and the current value:
  `exampleItems(...).Render(hui.WithListbox(r.Context(), hui.ListboxProps{ID: "country", Value: r.FormValue("country")}), w)`.
  Without it items still work: the plugin paints the same state as soon as they arrive.
- **Ids come from the value** (`<id>-option-<value>`, on the server and in the plugin), so
  the highlight and `aria-activedescendant` survive a redraw. When the highlighted item is
  gone, the first one is highlighted. Appending never moves the scroll position.
- **The keyboard crosses pages.** ArrowDown on the last loaded option scrolls the loader
  into view (which triggers `intersect`) and moves on as soon as the next page arrives.
- **A loader takes any markup.** `ListboxLoader`, `ComboboxLoader` and `TableLoader` render
  their children, so a spinner is whatever SVG or element you put inside. The loader is
  `aria-hidden` — the list announces the loading — and carries `data-pending` while its own
  request runs, so `group` on the loader plus `group-data-pending:` on the spinner shows it
  only while the next page is on its way.
- **A pending list is announced.** While a request from the loader or the search runs, the list
  has `aria-busy="true"` and `data-pending`. While a loader is in the list, options carry
  `aria-setsize="-1"` (unknown size) and `aria-posinset`.
- **State never lives in the items.** The value, the hidden inputs and the labels survive
  an option that is filtered out or not loaded yet. The plugin remembers the label of every
  option it has seen, and of the server-rendered `ListboxValue` (or the `Value` of
  `ComboboxInput`) for the initial value.
- **A whole component can be swapped.** Its plugin tears down the layer, listeners and
  observers, the teleported content leaves with it, and the new one starts from the state
  rendered by Go.
- Components inside a partial (a `Checkbox` in a table row, a `Tooltip` in a cell) ship
  their own script and initialise on their own.

**Table.** `TableProps.RowCount` renders `aria-rowcount` (`-1` when unknown) and
`TableRowProps.Index` `aria-rowindex` (1-based, header rows included), so a screen reader
knows where it is in a table that is only partly loaded. `TableLoader{Colspan}` is the
sentinel row: `<tr aria-hidden="true" data-hui-loader><td colspan>`.

## FileUpload

`FileUploadInput` is a native `<input type="file">` with `Name`, `Accept`, `Multiple`,
`Form`, `Disabled`, `Required` and `Invalid`, wired into a `Field` like any other control.
Inside a `FileUpload`, a `FileUploadDropzone` adds drag and drop:

```templ
@hui.FileUpload(hui.FileUploadProps{Class: "grid gap-2"}) {
    @hui.FileUploadDropzone(hui.FileUploadDropzoneProps{
        Class: "rounded border-2 border-dashed p-6 data-dragging:border-blue-600 has-focus-visible:border-blue-600"}) {
        @hui.FileUploadInput(hui.FileUploadInputProps{Name: "docs", Multiple: true, Accept: ".pdf,image/*", Class: "sr-only"})
        Drop files here, or click to choose
    }
    @hui.FileUploadList(hui.FileUploadListProps{}) {
        @hui.FileUploadFile(hui.FileUploadFileProps{}) {
            @hui.FileUploadFileName(hui.FileUploadFileNameProps{})
        }
    }
}
```

- Dropping puts the files into the input through `DataTransfer` and fires `input` and
  `change`, so a form submit or `hx-post` on the input sends them. Files that do not match
  `Accept` are dropped, and only the first one is kept without `Multiple`.
- A click anywhere in the dropzone opens the file picker. The dropzone is not focusable on
  its own — keyboard users reach the input, which can stay visually hidden with `sr-only`.
- The root and the dropzone carry `data-dragging` while files hover over them. Uploading
  is up to the form or htmx.
- `FileUploadList` shows the chosen files without writing any JavaScript: a `<ul>` that is
  hidden until something is chosen, `FileUploadFile` is the `<li>` repeated for every file
  and `FileUploadFileName` the `<span>` with its name. Put an icon or a size next to the
  name inside `FileUploadFile`.

```templ
@hui.FileUploadList(hui.FileUploadListProps{Class: "grid gap-1 text-sm"}) {
    @hui.FileUploadFile(hui.FileUploadFileProps{Class: "flex items-center gap-2"}) {
        @icon.File()
        @hui.FileUploadFileName(hui.FileUploadFileNameProps{})
    }
}
```

## Button

`Button` stands on its own and is not part of `Field`.

```templ
@hui.Button(hui.ButtonProps{Class: "btn", Attrs: templ.Attributes{"hx-post": "/cart"}}) {
    Add to cart
}
```

- `Type` defaults to `button`, so a button inside a form never submits it by accident.
  Pass `Type: "submit"` when you mean it, along with `Name`, `Value` or `Form` if needed.
- With `Href` the component renders an `<a>` instead, including
  `rel="noopener noreferrer"` when `Target` is `_blank`. The URL is sanitized, so a
  `javascript:` href never reaches the page.
- A disabled link stays an `<a>` but loses its `href` and gets `role="link"`,
  `aria-disabled="true"` and `tabindex="-1"`, because `<a>` has no `disabled` attribute.
  A disabled `<button>` uses the native attribute and no redundant `aria-disabled`.
- `LeftSlot` and `RightSlot` render around the children in
  `<span aria-hidden="true" data-slot="left|right">`, so an icon never leaks into the
  button's accessible name. A square icon button puts the icon in a slot and its name in
  children (`sr-only`). The slots carry no state: a spinner or an icon swap reacts to the
  button with `group` and `group-data-pending:`, because `data-pending` set during an htmx
  request lands on the button only.
- Every button carries exactly one of `data-enabled` and `data-disabled`, so hover and press
  styles read `data-enabled:hover:` and `data-enabled:active:` instead of
  `not-data-disabled:`. Tailwind's own `enabled:` matches `:enabled`, which an `<a>` never
  does. Every other button of the library follows the same rule — `DropdownItem`,
  `ModalTrigger`, `ModalClose`, `DropdownTrigger`, `CollapsibleTrigger`, `TooltipClose`,
  `ToastAction`, `ToastClose` and `TabsTrigger` — so one set of button classes styles them
  all.
- `Pending` implies `Disabled` and adds `aria-busy="true"` plus `data-pending`, so a
  spinner is pure CSS and a slow request cannot be fired twice.
- There is no Alpine state: `:hover`, `:active` and `:focus-visible` already exist in CSS,
  and a real `<button>` is keyboard operable on its own. Toggle buttons are one escape
  hatch away: `Attrs: templ.Attributes{":aria-pressed": "muted"}`.

## Overlays

`Modal`, `Dropdown`, `Tooltip`, `Listbox` and `Combobox` share one mechanism. None
of them uses the native `<dialog>` or `popover`: their content is teleported with Alpine's
`x-teleport` into a single container, so no ancestor with `transform`, `filter` or
`overflow` can clip it. Render that container once, directly in `<body>`, outside anything
htmx swaps:

```templ
<body>
    ...
    @hui.Portals(hui.PortalsProps{}) {
        @hui.Toaster(hui.ToasterProps{})
    }
</body>
```

It is required — every overlay teleports into `#hui-portals`.

**Positioning** is plain CSS [anchor
positioning](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_anchor_positioning), no
JavaScript measures anything. The trigger gets `anchor-name: --<id>` and the content an
inline `position: fixed` with `position-anchor`, `position-area` and
`position-try-fallbacks`, so it flips to the other side when it runs out of room. Every
content part embeds `Placement`: `Side` (`bottom`, `top`, `left`, `right`), `Align`
(`start`, `center`, `end`) and `SideOffset` in pixels. The inline style resets `inset` and
`margin`, so use `SideOffset` for the gap or pass your own `style` through `Attrs`, which
is appended last and wins. Teleported content is stacked in the order it was opened, but
it is still in the normal page, so give it a `z-index` that beats your layout. Content
positioned with `position-area` only gets the width of its area; `w-max` lets a short text
such as a tooltip take its natural width.

**One stack for all of them.** The plugins keep the open overlays in one stack:

- Escape only reaches the topmost one, so Escape in a `Listbox` inside a `Modal` closes the
  listbox, and the next Escape the modal.
- A click outside closes every overlay above the one you clicked into. A click inside a
  dropdown that belongs to a modal does not close the modal. A modal stops the walk, so a
  click on its overlay never closes what is below it.
- A modal makes everything except itself and the overlays opened from it `inert`. The
  `Toaster` and the `Announcer` are never made inert, so toasts stay usable and
  announcements are heard while a modal is open.

**Shared shape.** Every overlay root has `Open` (rendered from Go, so the right state is
there before Alpine starts), is `x-modelable` (`open` for the dialogs and menus, `value` for
`Listbox` and `Combobox`), dispatches `hui-<name>:opened` and `hui-<name>:closed` from the
root with `{ id }` (bound to `OnOpen` and `OnClose` with `.self`), and listens for
`hui-<name>:show`, `:hide` and `:toggle` on `window` with the id as the detail. The plugin
runs `htmx.process` on teleported content when it opens.

**Without JavaScript** nothing is teleported, so the content of an overlay is not rendered
at all. Triggers still report the right `aria-expanded`.

**Transitions are on by default.** Every part shown with `x-show` gets Alpine's
`x-transition`: the content of `Tooltip`, `Dropdown`, `Listbox`, `Combobox` and
`Modal`, every `Toast`, and the indicators of `Checkbox` and `Radio`. `ModalOverlay` gets
`x-transition.opacity`, `CollapsibleContent` gets `x-collapse` (the plugin is bundled,
nothing to install). Your own `Attrs` decide:

- any `x-transition…` or `x-collapse…` key replaces the default, for example
  `"x-transition.opacity.duration.100ms": true` or the `x-transition:enter`,
  `:enter-start`, … classes for Tailwind,
- `"x-transition": false` (or `"x-collapse": false`) turns it off — templ does not render
  a false attribute.

`Tabs` has no default transition: two panels animating at once would jump.

### Controlling state

Every component with an open state exports its controls as Alpine expressions. Put them on
any element inside the component, on any event, through `Attrs` or `Events` — this is how
a styled layer wires its own triggers:

```templ
@hui.Tooltip(hui.TooltipProps{ID: "vat"}) {
    <span tabindex="0" { templ.Attributes{
        "@pointerenter": hui.TooltipShowDelayed, "@pointerleave": hui.TooltipHideDelayed,
        "@focus": hui.TooltipFocus, "@blur": hui.TooltipBlur, "@click": hui.TooltipToggle}... }>
        VAT
    </span>
    @hui.TooltipContent(hui.TooltipContentProps{}) {
        Value added tax
    }
}
```

| Component | Controls |
|---|---|
| `Tooltip` | `TooltipShow`, `TooltipHide`, `TooltipToggle`, `TooltipShowDelayed`, `TooltipHideDelayed`, `TooltipFocus`, `TooltipBlur`, `TooltipDismiss` |
| `Dropdown` | `DropdownShow`, `DropdownHide`, `DropdownToggle`, `DropdownKeydown` |
| `Modal` | `ModalShow`, `ModalHide`, `ModalToggle` |
| `Collapsible` | `CollapsibleShow`, `CollapsibleHide`, `CollapsibleToggle` |
| `Listbox` | `ListboxShow`, `ListboxHide`, `ListboxToggle` |
| `Combobox` | `ComboboxShow`, `ComboboxHide`, `ComboboxToggle` |

- Each control calls a method named after its component (`tooltipToggle($el, $event)`, not
  `toggle()`). Alpine merges the scopes of nested components, so `hui.ModalHide` on a
  button inside a tooltip inside a modal still closes the modal.
- **The element becomes the anchor.** `Tooltip` and `Dropdown` move
  `anchor-name` to the element that opened them, so the content is positioned next to it.
  The tooltip adds its content to that element's `aria-describedby`; the
  dropdown gives it `aria-controls`, `aria-expanded` and `data-open`, and return focus to
  it on close.
- The built-in triggers (`DropdownTrigger`, `ModalTrigger`, …) are made of the same
  controls, so the internal wiring is the same code path. Adding a control to their
  `Events` chains it after the internal one. `Tooltip` has no built-in trigger: the controls
  are its only API.
- An element outside the component's `x-data` uses the window events instead:
  `$dispatch('hui-tooltip:toggle', 'vat')`, or `hx-on` with
  `window.dispatchEvent(new CustomEvent('hui-modal:show', { detail: 'confirm' }))`.

## Modal

`Modal` is a dialog whose state lives in Alpine. `open` is the single source of truth:
whatever sets it — a trigger, `x-model`, an event, a server-rendered `Open` — the plugin
reacts the same way, so the modal can be driven from anywhere without the DOM and the
state ever disagreeing. The markup follows bits-ui's Dialog.

The overlay and the content are teleported into `Portals` (see [Overlays](#overlays)).
Write the modal wherever it belongs:

```templ
@hui.Modal(hui.ModalProps{ID: "confirm", OnClose: "$store.cart.refresh()"}) {
    @hui.ModalTrigger(hui.ModalTriggerProps{Class: "btn"}) {
        Delete
    }
    @hui.ModalOverlay(hui.ModalOverlayProps{Class: "fixed inset-0 bg-black/50"})
    @hui.ModalContent(hui.ModalContentProps{Class: "fixed inset-0 m-auto size-fit rounded-xl p-6"}) {
        @hui.ModalTitle(hui.ModalTitleProps{}) {
            Delete this item?
        }
        @hui.ModalDescription(hui.ModalDescriptionProps{}) {
            This cannot be undone.
        }
        @hui.ModalClose(hui.ModalCloseProps{}) {
            Cancel
        }
    }
}
```

| Part | Renders | Wiring |
|---|---|---|
| `Modal` | `<div>` | `x-data`, `x-modelable="open"`, `data-open` |
| `ModalTrigger` | `<button type="button">` | opens on click, `aria-haspopup="dialog"`, `aria-controls`, `aria-expanded`, `data-open` |
| `ModalOverlay` | `<template x-teleport>` with `<div aria-hidden="true">` | `x-show`, `data-open`, id `<id>-overlay` |
| `ModalContent` | `<template x-teleport>` with `<div role="dialog" aria-modal="true" tabindex="-1">` | `x-show`, `data-open`, id `<id>-content`, `aria-labelledby`, `aria-describedby`, closes on a click outside; `role="alertdialog"` with `Alert` |
| `ModalTitle` | `<h2>` | id `<id>-title` |
| `ModalDescription` | `<p>` | id `<id>-description` |
| `ModalClose` | `<button type="button">` | closes on click |

The content always points `aria-labelledby` at the title, so render a `ModalTitle` — a
dialog without a name is an accessibility bug. `aria-describedby` points at the
description; without one the reference is simply ignored by assistive technology.

Overlay and content are siblings, each teleported on its own, or the content can sit
inside the overlay when you want to center it with flexbox — it then travels with the
overlay instead of being teleported separately. A click on the overlay closes the modal
either way, a click inside the content never does. The component adds no positioning, so
place both with `fixed` yourself.

**What teleporting means for you.** The root and the trigger stay where you wrote them,
and the Alpine scope comes along, so `open`, `hide()` and the bindings of the surrounding
`x-data` still work inside the content. The DOM, however, lives in `Portals`:

- A `<form>` around the `Modal` does not contain the teleported inputs — put the form
  inside `ModalContent`.
- Events from inside the content bubble to `Portals` and `<body>`, not to the
  ancestors of the modal on your page. `hui-modal:opened` and `hui-modal:closed` are
  dispatched from the root, so listening for them works as usual.
- htmx never sees the inside of a `<template>`, so the plugin runs `htmx.process` on the
  overlay and the content when the modal opens; `hx-*` attributes in there just work.
- When the modal itself is removed (an htmx swap, for example), Alpine removes its
  teleported parts from `Portals` too.

**What the plugin does while the modal is open:**

- Moves focus into the content: the first element with `autofocus`, otherwise the first
  tabbable element, otherwise the content itself. On close, focus returns to whatever had
  it before, falling back to the trigger.
- Keeps Tab and Shift+Tab cycling inside the content.
- Makes everything outside the content `inert`, so the rest of the page can be neither
  clicked, focused nor read by a screen reader — the same thing a native modal dialog does.
- Closes on Escape and on a click outside the content. If another overlay (a `Dropdown`,
  a `Listbox`, …) is open inside the modal, Escape and outside clicks close that one first.
- Locks page scroll.

Modals stack: open one from inside another and only the topmost one reacts to Escape and
outside clicks; closing it hands focus and interactivity back to the one below.

**Controlling it.** Inside the modal you are already in its Alpine scope, so `open`,
`show()`, `hide()` and `toggle()` are available directly. To bind it to your own state, pass
`x-model` — the modal is `x-modelable`, so the binding goes both ways:

```templ
<div x-data="{ confirming: false }">
    <button @click="confirming = true">Delete</button>
    @hui.Modal(hui.ModalProps{ID: "confirm", Attrs: templ.Attributes{"x-model": "confirming"}}) {
        ...
    }
</div>
```

From anywhere else on the page, dispatch an event with the modal's id — it works from a
button, from Alpine, from plain JavaScript and from htmx alike:

```html
<button @click="$dispatch('hui-modal:show', 'confirm')">Delete</button>
```

```js
window.dispatchEvent(new CustomEvent('hui-modal:hide', { detail: 'confirm' }))
```

`hui-modal:show`, `hui-modal:hide` and `hui-modal:toggle` all take the id as their detail,
so only the modal you name reacts.

**Reacting to it.** The modal dispatches `hui-modal:opened` and `hui-modal:closed` from its
root, with `{ id }` as the detail. `OnOpen` and `OnClose` are Alpine expressions bound to
exactly those events with the `.self` modifier, so a nested modal does not trigger its
parent's hooks. The events bubble, so you can also listen further up — or let htmx do it:

```html
hx-trigger="hui-modal:closed from:#confirm"
```

Anything you pass in `Attrs` under `@hui-modal:closed.self` is chained, not replaced, so
your own listener and `OnClose` both run.

**Transitions.** The overlay fades (`x-transition.opacity`) and the content uses the
default `x-transition`, see [Overlays](#overlays).

**State across htmx swaps.** `Open` is the server's opinion. It is rendered into the
markup — a closed modal hides its overlay and content with an inline `display: none`, an
open one carries `data-open` and `aria-expanded="true"` — and Alpine opens the modal as
soon as it starts. When a swap replaces the modal, the server decides again. Set
`Persist: true` if the client should decide instead — the open state is then mirrored in
`sessionStorage` under the modal's id and restored on init, so it survives swaps, back
navigation and reloads. It needs a stable `ID`; a generated one changes between renders.
An `x-model` binding always wins over both. If you would rather htmx did not touch the
modal at all, `hx-preserve` through `Attrs` keeps the live element with its state
untouched.

**The rest of the props.** `Static: true` stops Escape and outside clicks from closing the
modal (you close it yourself, for example after a request succeeds). `Alert: true` makes it
an [alert dialog](https://www.w3.org/WAI/ARIA/apg/patterns/alertdialog/): the content gets
`role="alertdialog"` and a click outside does nothing, Escape still closes it. Put
`autofocus` on the least destructive button. `AllowScroll: true`
leaves the background scrollable. `ModalTrigger` takes `Disabled`. Every part works on its
own outside a `Modal`, just without the wiring.

**The script.** Modal's logic is an Alpine plugin (built from
[`alpine/src/modal.js`](alpine/src/modal.js)), loaded with the rest, see [Alpine](#alpine).
Without Alpine nothing is teleported, so the overlay and the content are not rendered at
all — not even for `Open: true`. The trigger still reports the right `aria-expanded`.

## Dropdown

`Dropdown` is a menu button: a trigger that opens a list of actions or links. The list is
teleported and anchored like every other [overlay](#overlays), and the plugin adds what the
[menu pattern](https://www.w3.org/WAI/ARIA/apg/patterns/menu-button/) asks for: arrow keys,
Home/End, typeahead, Escape, closing on a click outside or after a choice, and returning
focus to the trigger.

```templ
@hui.Dropdown(hui.DropdownProps{ID: "account", Class: "inline-block"}) {
    @hui.DropdownTrigger(hui.DropdownTriggerProps{Class: "btn data-open:bg-gray-100"}) {
        Account
    }
    @hui.DropdownContent(hui.DropdownContentProps{SideOffset: 4, Class: "z-50 rounded-md border p-1 shadow"}) {
        @hui.DropdownItem(hui.DropdownItemProps{Href: "/profile", Class: "item focus:bg-gray-100"}) {
            Profile
        }
        @hui.DropdownItem(hui.DropdownItemProps{
            Class: "item focus:bg-gray-100",
            Attrs: templ.Attributes{"hx-post": "/archive"},
        }) {
            Archive
        }
        @hui.DropdownSeparator(hui.DropdownSeparatorProps{Class: "my-1 h-px bg-gray-200"})
        @hui.DropdownItem(hui.DropdownItemProps{Type: "submit", Form: "logout", Class: "item"}) {
            Sign out
        }
    }
}
```

**The parts.** `Dropdown` renders a `<div>` that holds the Alpine state, so style it like
any wrapper (`inline-block`, `relative`). `DropdownTrigger` is a `<button>` with
`aria-haspopup="menu"`, `aria-controls` and `aria-expanded`. `DropdownContent` is a
teleported `<div role="menu">` labelled by the trigger. `DropdownItem` is a `<button
role="menuitem" tabindex="-1">` — `type="button"` unless you pass `Type: "submit"` — or an
`<a>` when `Href` is set, with the same rules for `Target` and `Disabled` as `Button`.
`DropdownGroup` is a `role="group"` labelled by the `DropdownLabel` inside it, and
`DropdownSeparator` a `role="separator"` divider. `DropdownItem` takes `LeftSlot` and `RightSlot` like `Button`,
rendered in `<span aria-hidden="true" data-slot="left|right">`, so an icon never becomes
part of the item's name or of the typeahead match.

**IDs.** The dropdown owns them: the wrapper gets `ID` (or a generated one), the trigger
`<id>-trigger` and the content `<id>-content`. An `ID` set on the trigger or the content
inside a `Dropdown` is ignored, just like inside a `Field`.

**Positioning** comes from `Placement`, see [Overlays](#overlays).

**Styling the state.** Items get real DOM focus, both from the keyboard and on hover, so
`focus:` is the highlight — there is no separate `data-highlighted`. The wrapper, the
trigger and the content carry `data-open` while the menu is open. Disabled items keep
`data-disabled` and are skipped by the arrow keys.

**Keyboard.** Enter, Space or ArrowDown on the trigger open the menu on the first item,
ArrowUp on the last one; a mouse click focuses the menu itself. Inside, the arrow keys
wrap around, Home and End jump to the ends, typing a letter jumps to the next item that
starts with it, Escape closes the menu and returns focus to the trigger, and Tab closes it
and moves on.

**Choosing an item** closes the menu, unless your own click handler calls
`$event.preventDefault()`. The item still does its own job — a link navigates, a submit
button submits its form, `hx-*` attributes fire their request.

**Checkbox and radio items.** `DropdownCheckboxItem` is a `<button
role="menuitemcheckbox">` and `DropdownRadioItem` a `<button role="menuitemradio">`, both
with `aria-checked` and `data-checked` like `Checkbox` and `Radio`. Radio items belong to a
`DropdownRadioGroup` — a `role="group"` labelled by its `DropdownLabel` — which holds
`Name`, `Value`, `Form` and `Disabled`; the item whose `Value` matches is checked. A radio
item outside a group takes `Checked` itself.

```templ
@hui.DropdownContent(hui.DropdownContentProps{Class: "z-50 rounded-md border p-1 shadow"}) {
    @hui.DropdownCheckboxItem(hui.DropdownCheckboxItemProps{
        Name: "columns", Value: "total", Form: "view", Checked: true,
        LeftSlot: hui.CheckboxIndicator(hui.CheckboxIndicatorProps{}),
    }) {
        Total
    }
    @hui.DropdownSeparator(hui.DropdownSeparatorProps{})
    @hui.DropdownRadioGroup(hui.DropdownRadioGroupProps{Name: "sort", Value: "date", Form: "view"}) {
        @hui.DropdownLabel(hui.DropdownLabelProps{}) {
            Sort by
        }
        @hui.DropdownRadioItem(hui.DropdownRadioItemProps{
            Value: "date",
            LeftSlot: hui.RadioIndicator(hui.RadioIndicatorProps{}),
        }) {
            Date
        }
    }
}
```

- The indicator is the one from `Checkbox` and `Radio`: `CheckboxIndicator` or
  `RadioIndicator` anywhere inside the item, usually in `LeftSlot`. It is hidden from Go
  when the item is unchecked, so the state is right before Alpine starts.
- A checked item sends its value like `Checkbox`: `name`/`value` on the button and a hidden
  input inside it. The content is teleported out of your `<form>`, so set `Form`.
- Every change dispatches a bubbling `change` from the item. With htmx, react to it on the
  item itself — `hx-trigger="change"` fires after the state changed — or submit the form
  the items belong to.
- Choosing a radio item closes the menu, a checkbox item keeps it open so you can toggle
  several. Add `DropdownHide` to the item's `OnClick` to close it, or call
  `$event.preventDefault()` on a radio item to keep the menu open.
- Arrow keys, Home/End and typeahead move over all three kinds of items.

**Opening, closing and reacting** works like `Modal`: `Open`, `x-model`, the
`hui-dropdown:show|hide|toggle` events and `OnOpen`/`OnClose` for `hui-dropdown:opened`
and `hui-dropdown:closed`. The logic is the `dropdown` Alpine plugin, see [Alpine](#alpine).

## Tooltip

`Tooltip` shows a short description next to an element. It has no trigger part: you put
the [controls](#controlling-state) on your own element and choose the events yourself.

```templ
@hui.Tooltip(hui.TooltipProps{ID: "save-hint", Class: "inline-block"}) {
    <button type="button" aria-label="Save" aria-describedby="save-hint-content" { templ.Attributes{
        "@pointerenter": hui.TooltipShowDelayed, "@pointerleave": hui.TooltipHideDelayed,
        "@focus": hui.TooltipFocus, "@blur": hui.TooltipBlur}... }>
        @icon.Save()
    </button>
    @hui.TooltipContent(hui.TooltipContentProps{Class: "z-50 w-max rounded bg-gray-900 px-2 py-1 text-xs text-white"}) {
        Saved 2 minutes ago
    }
}
```

- `TooltipContent` is a teleported `<div role="tooltip">` with `id="<id>-content"`, placed
  above its anchor and centred unless you say otherwise.
- **Give the element `aria-describedby="<id>-content"` yourself**, so the text is read even
  before the first interaction. That element is also the default anchor, used when the
  tooltip opens from `Open`, `x-model` or a `hui-tooltip:*` event. Any other element that
  opens it through a control becomes the anchor and gets `aria-describedby` added.
- Hover: `TooltipShowDelayed` opens after `Delay` (700 ms), `TooltipHideDelayed` closes
  `CloseDelay` (300 ms) after the pointer leaves; moving the pointer onto the tooltip keeps
  it open. Both ignore touch. Focus: `TooltipFocus` opens right away on keyboard focus
  (`:focus-visible`), `TooltipBlur` closes it. Click: `TooltipToggle`, which is also the
  only way to open it on a touch screen. `TooltipShow` and `TooltipHide` open and close it
  immediately.
- Pick any combination: hover and focus is the classic tooltip, `TooltipToggle` alone opens
  it on click only, all of them together open it on hover and also on click.
- Escape, a click outside and a `TooltipClose` button inside the content close it
  (`TooltipDismiss`: it stays closed until the pointer leaves).
- Delays of 0 use the defaults; pass `Delay: 1` for an instant tooltip. `Open`, `x-model`
  and the `hui-tooltip:*` events work like for every overlay.
- Without JavaScript the tooltip is never shown, so do not put anything in it that the user
  cannot get elsewhere.

## Tabs

```templ
@hui.Tabs(hui.TabsProps{ID: "settings", Value: "profile"}) {
    @hui.TabsList(hui.TabsListProps{Attrs: templ.Attributes{"aria-label": "Settings"}}) {
        @hui.TabsTrigger(hui.TabsTriggerProps{Value: "profile", Class: "data-selected:border-b-2"}) {
            Profile
        }
        @hui.TabsTrigger(hui.TabsTriggerProps{Value: "billing", Class: "data-selected:border-b-2"}) {
            Billing
        }
    }
    @hui.TabsContent(hui.TabsContentProps{Value: "profile"}) {
        ...
    }
    @hui.TabsContent(hui.TabsContentProps{Value: "billing", Attrs: templ.Attributes{
        "hx-get": "/settings/billing", "hx-trigger": "intersect once"}}) {
    }
}
```

- `TabsList` is the `role="tablist"`, `TabsTrigger` a `<button role="tab">` and `TabsContent`
  a `role="tabpanel"` that is focusable and labelled by its tab. IDs are
  `<id>-trigger-<value>` and `<id>-content-<value>`.
- The selected state is rendered from Go: `aria-selected`, `data-selected`, the roving
  `tabindex` and `display: none` on the other panels, so the right tab is visible before
  Alpine starts.
- The arrow keys of the `Orientation` (`horizontal` by default, `vertical` gives
  `aria-orientation`) move between enabled tabs and wrap, Home and End jump to the ends.
  `Activation: "manual"` only moves focus and leaves the choice to Enter or Space.
- `value` is `x-modelable`, and every change dispatches `hui-tabs:change` with
  `{ id, value }` (`OnChange`). A panel loaded by htmx is one `hx-trigger` away, as above.

## Collapsible

`Collapsible` is a single show/hide: `Collapsible{Open}`, `CollapsibleTrigger` (with
`aria-expanded` and `aria-controls`) and `CollapsibleContent`, with the same `x-model`,
`OnOpen`/`OnClose` and `hui-collapsible:*` events as the overlays, and the
`CollapsibleShow|Hide|Toggle` [controls](#controlling-state).

**The row with the title.** It supports the two usual layouts:

- **The whole top is the button.** Put the title inside `CollapsibleTrigger` and make it
  full width.
- **A control button on the right.** `CollapsibleHeader` (a `<div>`) holds a
  `CollapsibleTitle` (`<span>`, `<id>-title`) and a small trigger. The trigger is labelled
  by the title and itself (`aria-labelledby="<title> <trigger>"`), so an icon-only button
  is announced with the title; an `aria-label` in its `Attrs` replaces that.

```templ
@hui.Collapsible(hui.CollapsibleProps{ID: "shipping"}) {
    @hui.CollapsibleHeader(hui.CollapsibleHeaderProps{Clickable: true, Class: "flex items-center justify-between"}) {
        @hui.CollapsibleTitle(hui.CollapsibleTitleProps{}) {
            Shipping address
        }
        @hui.CollapsibleTrigger(hui.CollapsibleTriggerProps{Class: "data-open:rotate-180"}) {
            @icon.ChevronDown()
        }
    }
    @hui.CollapsibleContent(hui.CollapsibleContentProps{}) {
        ...
    }
}
```

`Clickable: true` on the header makes a click anywhere in the row toggle too
(`data-clickable` is there for `cursor-pointer`). Clicks on links, buttons and inputs
inside the row keep their own meaning, a text selection does not toggle, and the keyboard
stays on the trigger. Header and title carry `data-open`.

## Breadcrumb

A navigation landmark without JavaScript.

```templ
@hui.Breadcrumb(hui.BreadcrumbProps{}) {
    @hui.BreadcrumbList(hui.BreadcrumbListProps{Class: "flex gap-2"}) {
        @hui.BreadcrumbItem(hui.BreadcrumbItemProps{}) {
            @hui.BreadcrumbLink(hui.BreadcrumbLinkProps{Href: "/"}) {
                Home
            }
        }
        @hui.BreadcrumbSeparator(hui.BreadcrumbSeparatorProps{}) {
            /
        }
        @hui.BreadcrumbItem(hui.BreadcrumbItemProps{}) {
            @hui.BreadcrumbLink(hui.BreadcrumbLinkProps{Current: true}) {
                Settings
            }
        }
    }
}
```

- `Breadcrumb` is a `<nav aria-label="Breadcrumb">` (`Label` translates it) around an `<ol>`.
  `BreadcrumbLink` is an `<a>`, or a `<span>` without `Href`; `Current` adds
  `aria-current="page"` and `data-current`. The separator is hidden from screen readers.

## Table

`Table`, `TableCaption`, `TableHeader`, `TableBody`, `TableFooter`, `TableRow`, `TableHead`
and `TableCell` render the native table elements. On top of that:

- `TableHead` gets `scope="col"` (or your `Scope`), and `Sort` (`ascending`, `descending`,
  `none`, `other`) renders `aria-sort` and `data-sort`.
- `TableRow{Selected}` renders `data-selected`, `TableRow{Index}` `aria-rowindex`.
- `Table{RowCount}` renders `aria-rowcount` and `TableLoader` is the sentinel row for an
  infinite scroll, see [Server-rendered lists](#server-rendered-lists).
- **Stacking on narrow screens.** `TableCell{Label}` renders `data-label`, so a stylesheet
  can turn every row into a card and show the column name before each value with
  `before:content-[attr(data-label)]` — no hidden duplicate of the header in the DOM. The
  label sits on the cell rather than being looked up from the header, so rows rendered on
  their own by htmx carry it too. Every part also renders its explicit role (`table`,
  `rowgroup`, `row`, `columnheader`/`rowheader`, `cell`): they are redundant on the native
  elements, but Chrome and Safari drop the table semantics for screen readers as soon as
  CSS changes `display`, and the roles bring them back.

There is no JavaScript: every interaction is the server's. Sorting is a `Button` with
`hx-get` in the header, selecting a row is a `Checkbox` with `hx-post`, and "select all" is
a `Checkbox` whose `Checked` and `Indeterminate` the server works out. The server renders
the table again with the new state.

## Meter

A `<div>` with the right role, so it can be styled freely; the value is exposed as a CSS
variable.

```templ
@hui.Meter(hui.MeterProps{Value: charge, Low: 30, High: 70, Optimum: "high", Label: "Battery", Class: "group h-2 rounded bg-gray-200"}) {
    <div class="h-full w-(--hui-meter) group-data-[state=critical]:bg-red-600"></div>
}
```

- `Meter` renders `role="meter"` for a value in a known range (`Min`, `Max`), with
  `--hui-meter`. `Low` and `High` split the range and `data-level` says where the value is
  (`low`, `medium`, `high`). `Optimum` (`low`, `medium`, `high`) names the good region and
  adds `data-state="optimum|suboptimal|critical"`, like the native `<meter>`.
- It takes `Label` (`aria-label`) and `ValueText` (`aria-valuetext`). Inside a `Field` it
  takes the field's id and is labelled by its `Label` through `aria-labelledby`.

## Toast

Toasts come from the server. Put a `Toaster` in your layout — inside `Portals` is a good
place — and send `Toast`s into it with an htmx partial from any response:

```templ
<hx-partial hx-target="#hui-toaster" hx-swap="beforeend">
    @hui.Toast(hui.ToastProps{Class: "rounded border bg-white p-4 shadow"}) {
        @hui.ToastTitle(hui.ToastTitleProps{}) {
            Saved
        }
        @hui.ToastDescription(hui.ToastDescriptionProps{}) {
            Your changes are live.
        }
        @hui.ToastAction(hui.ToastActionProps{Attrs: templ.Attributes{"hx-post": "/undo"}}) {
            Undo
        }
        @hui.ToastClose(hui.ToastCloseProps{}) {
            Close
        }
    }
</hx-partial>
```

- `Toaster` is an `<ol id="hui-toaster" aria-live="polite">` (`ID` and `Label` can be
  changed). A `Toast` is an `<li role="status">`, or
  `role="alert"` with `Urgent`, labelled by its title and described by its description.
- It closes itself after `Duration` (5000 ms by default, `-1` keeps it until closed). The
  timer pauses while the pointer is over the toast, while focus is inside it and while the
  page is hidden. Escape on a focused toast closes it and moves focus to the next one.
- A closed toast leaves the DOM once its `x-transition` has finished and dispatches
  `hui-toast:closed` (`OnClose`).
- For a progress bar, a toast with a timer carries `--hui-toast-duration` (for example
  `5000ms`) and `data-paused` while the timer is paused, so the bar is pure CSS:
  `animation-duration: var(--hui-toast-duration)` and
  `group-data-paused:[animation-play-state:paused]`. A toast that stays (`-1`) has no
  variable.

## Announcer and Alert

An htmx swap is silent for a screen reader: the table was filtered, twelve rows came back,
and nobody said so. `Announcer` is a live region the server writes into. Put one in your
layout — next to the `Toaster` inside `Portals` is a good place — and keep it visually
hidden:

```templ
@hui.Announcer(hui.AnnouncerProps{Class: "sr-only"})
```

Then send the message with any response:

```templ
<hx-partial hx-target="#hui-announcer" hx-swap="innerHTML">
    12 orders match "spring".
</hx-partial>
```

- `Announcer` is a `<div id="hui-announcer" role="status" aria-live="polite"
  aria-atomic="true" data-hui-announcer>`. With `Urgent` it is `role="alert"` and
  `aria-live="assertive"` with the id `hui-announcer-urgent`, so a page can have both.
  `ID` changes the id.
- The region has to be on the page before the message arrives — a live region that is
  swapped in together with its text is not announced. Replace its content, do not replace
  the region.
- Some screen readers skip a message identical to the previous one; say what changed
  ("12 orders match …") rather than a fixed "Updated".
- It needs no JavaScript and is never made `inert` by a `Modal`.

`Alert` is an inline message about the whole page or form — a failed submit, a lost
connection — a `<div role="alert">`, or `role="status"` with `Polite`. A message about one
field belongs in `Error` inside its `Field`. An alert is announced when it is inserted, so
swap it in with htmx; one that is already in the HTML on page load is not read out.

```templ
@hui.Alert(hui.AlertProps{ID: "checkout-error", Class: "rounded border border-red-300 p-3"}) {
    The card was declined. Try another one.
}
```

## SkipLink

`SkipLink` is an `<a href="#main">` (`Href` to change it) for keyboard users to jump past
the navigation. Hide it until it is focused: `sr-only focus:not-sr-only`. Give the target a
`tabindex="-1"` so it can take focus.

## Tailwind

The state attributes are rendered bare (`data-invalid`, not `data-invalid="true"`), which
is exactly what Tailwind v4 variants expect:

```html
class="border-gray-300 data-invalid:border-red-500 data-disabled:opacity-50"
```

Because every part of a field carries the same state, you rarely need `group`: the input,
the label and the error each react to `data-invalid:` on their own. Reach for
`group-data-invalid:` only when you style something that has no state of its own.

`Class` and a `class` passed through `Attrs` are merged, with `Class` last. Plain
concatenation is not enough for Tailwind — `px-4` and `px-10` would both survive and the
winner would be decided by the stylesheet, not by your call. Plug a merger in once:

```go
hui.MergeClass = func(existing, added string) string {
    return twmerge.Merge(existing, added)
}
```

There is no Tailwind v4 aware merger for Go yet — `tailwind-merge-go` covers v3.0–v3.4 —
so until one lands, keep the utilities of one property in one place and let the default
join do its job.

If you build a styled layer as a separate module, remember that Tailwind has to see the
class strings, and that module lives in the Go module cache. Either ship compiled CSS
from that repo, or vendor it and point Tailwind at the copy:

```css
@source "../vendor/github.com/you/ui/**/*.templ";
```

Keep class strings whole in the source. Tailwind cannot find a class you assemble from
fragments at runtime.

## htmx

`hx-*` attributes go through `Attrs` untouched — `merge` only ever rewrites keys that
start with `@`.

The wrapper id makes the usual server-side validation loop a one-liner: post the field,
render the same `Field` with `Error` filled in, and swap it in place.

```templ
@hui.TextInput(hui.TextInputProps{
    Attrs: templ.Attributes{
        "hx-post":   "/validate/email",
        "hx-target": "#email-field",
        "hx-swap":   "outerHTML",
    },
})
```

Keep ids stable across responses. htmx preserves the focus of an input between requests
only *"for inputs that have a defined id attribute"*, and `hx-target`, `hx-swap-oob` and
CSS transitions all rely on the same thing. In practice that means giving every field a
`Name` — it is what the id is derived from. For markup that has no name to lean on, scope
the generated ids per render so a fragment comes back with the ids it had in the full page:

```go
ctx = hui.WithIDScope(ctx, "signup")
```

Lists that the server filters or loads page by page (`Listbox`, `Combobox`, `Table`) are
covered in [Server-rendered lists](#server-rendered-lists).

### Pending requests

While an htmx request runs, three elements get `data-pending` and `aria-busy="true"`: the
element that sent it, the submit button that triggered it (for a form) and the element its
response will be swapped into (`hx-target`). They lose both as soon as the request is
finished, so the pending look is plain Tailwind on the component you already style:

```templ
@hui.Button(hui.ButtonProps{
    Class: "btn data-pending:pointer-events-none data-pending:opacity-60",
    Attrs: templ.Attributes{"hx-post": "/cart", "hx-target": "#cart"},
})

<div id="cart" class="transition-opacity data-pending:opacity-50">…</div>
```

- It works the same for every component: a `Switch` with `hx-post`, a `TextInput` that
  validates its `Field` (`hx-target="#email-field"` dims the whole row), a sort `Button` in
  a `TableHead` whose target is the `Table`, a `Listbox` root with `hx-include="this"`.
- Something without state of its own reacts to an ancestor with `in-data-pending:` (or
  `group` and `group-data-pending:`), for example a spinner inside a button.
- Requests that overlap are counted, so an element stays pending until the last one ends.
  A server-rendered `data-pending` (`Button{Pending}`) is never removed. A response that
  leaves the page (`HX-Redirect`, `HX-Refresh`) keeps the state until the page goes.
- The state is only visual. To really block a second submit, add htmx's own
  `"hx-disable": "this"`, which sets `disabled` for the duration of the request.
- The script is tiny and needs no Alpine. `@hui.Scripts()` and every plugin script bring it
  along. Without Alpine, render `@hui.PendingScript()` once in the layout.

## Outside a form

Every form control works without a `<form>` around it.

- **Somewhere else in the page.** `Form` on a control, a `Field` or a `Fieldset` renders the
  native `form="<id>"` attribute, so the value is submitted with that form wherever the
  control sits. `Listbox` and `Combobox` put it on their hidden inputs.
- **On its own with htmx.** A control with a `name` sends its own value, so `hx-post`
  with `hx-trigger="change"` on a `Checkbox`, `Switch`, `Radio`, `TextInput` or `Select` is all you
  need (the checkable buttons carry `name` and `value` while they are checked). `Listbox`, `Combobox` and `Range` keep their values in inputs inside a
  wrapper, so put the `hx-*` attributes on the root and add `hx-include="this"`, which
  includes every input inside it:

  ```go
  Attrs: templ.Attributes{"hx-post": "/filter", "hx-trigger": "change", "hx-include": "this"}
  ```

  A `FileUploadInput` needs `hx-encoding="multipart/form-data"` next to its `hx-post`.
- Every change fires a bubbling `change` (and `input` where it makes sense) on the element
  that carries the `name`, so `hx-trigger="change"` works the same with or without a form.
  No plugin looks for `closest('form')`; a form reset is found through `input.form`, which
  respects the `form` attribute.

## Alpine

`Modal`, `Dropdown`, `Tooltip`, `Listbox`, `Combobox`, `Tabs`, `Collapsible`, `Toast`, `FileUpload`, `Range`, `Checkbox`, `Switch` and `Radio` need Alpine. Alpine is used for
interaction only: focus, open state, highlight, keyboard and the value in a hidden input.
Filtering, sorting, selecting rows and rendering are the server's job.

**No component renders a `<script>`.** Nothing runs on its own: you load the plugins once,
in the way that suits you, and the components only render markup — also inside htmx
partials. Pick one:

- `@hui.Scripts()` once in the layout renders the whole bundle (every plugin plus the
  [pending request](#pending-requests) tracker), embedded in the Go module, so no Node or
  npm is needed. Put it before Alpine, or after it — it registers itself either way.
- `@hui.ModalScript()`, `@hui.ComboboxScript()`, … render only the plugins a page uses.
  Each one is rendered once per page however many times you call it.
- Bundle it yourself from [`alpine/`](alpine), the npm package
  `@daarxwalker/headlessui-alpine`:

```js
import Alpine from 'alpinejs'
import headlessui from '@daarxwalker/headlessui-alpine'

Alpine.plugin(headlessui)
Alpine.start()
```

Registering the same plugin twice is harmless — `Alpine.data` just replaces the previous
definition with an identical one. The default export also starts the
[pending request](#pending-requests) tracker; without Alpine, import `{ pending }` and call
`pending()` once.

`TextInput` can keep local state in Alpine, but it is off by default — with Tailwind,
`focus-visible:` and `placeholder-shown:` cover the same ground without any JavaScript.
Turn it on when you need the state as an attribute, for example to style a parent:

```templ
@hui.TextInput(hui.TextInputProps{Alpine: true, Class: "input"})
```

That renders `x-data="{ focused: false, empty: true }"`, keeps `data-focused` and
`data-empty` in sync through `@focus`, `@blur` and `@input`, and still chains anything you
pass in `Events` after the internal handler. With Alpine off, no state attribute is
rendered that could go stale.

## Development

```bash
make example
```

`make example` is the development loop: it builds the Alpine plugins, runs
`templ generate` and serves `example/` on http://localhost:5050. `example/` is the
showcase app of this repository, built exactly as described in
[Building a showcase app](#building-a-showcase-app), with a styled layer that follows
[Building a project UI library](#building-a-project-ui-library). For now it lives only
locally and is not part of the repository.

Before a change is done:

```bash
templ generate && go test ./...
templ generate -check && gofmt -l . && go vet ./...
cd alpine && npm install && npm run build
```

`alpine/dist` is committed, just like the generated `*_templ.go` files, because Go embeds
it. Rebuild it after every change in `alpine/src`.

## License

[MIT](LICENSE). The npm package [`@daarxwalker/headlessui-alpine`](alpine) is released under
the same license.
