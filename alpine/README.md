# hui Alpine plugins

The [Alpine.js](https://alpinejs.dev) plugins behind the hui components. This directory is
only their source and build. It is not published anywhere: the Go module embeds `dist`,
and apps get the code through Go, either inline with `hui.Scripts()` or as a file in their
own sources with `go tool hui assets` (see
[Alpine](https://github.com/daarxwalker/hui#alpine) in the main README).

## Build

```bash
npm install
npm run build
```

| File | For |
|---|---|
| `dist/module.esm.js` | `go tool hui assets`, which writes it as `hui.js` with a version header |
| `dist/cdn.min.js` | `hui.Scripts()`, every plugin in one self-registering script |
| `dist/<component>.min.js` | `hui.<Component>Script()`, a single self-registering plugin |

`dist` is committed, because Go embeds it. Rebuild it after every change in `src`.

## Components

| Plugin | Registers | Docs |
|---|---|---|
| `checkbox` | `Alpine.data('huiCheckbox')` | [Checkbox](https://github.com/daarxwalker/hui#checkbox-switch-and-radio) |
| `collapsible` | `Alpine.data('huiCollapsible')` | [Collapsible](https://github.com/daarxwalker/hui#collapsible) |
| `combobox` | `Alpine.data('huiCombobox')` | [Combobox](https://github.com/daarxwalker/hui#combobox) |
| `dropdown` | `Alpine.data('huiDropdown')` | [Dropdown](https://github.com/daarxwalker/hui#dropdown) |
| `fileUpload` | `Alpine.data('huiFileUpload')` | [FileUpload](https://github.com/daarxwalker/hui#fileupload) |
| `listbox` | `Alpine.data('huiListbox')` | [Listbox](https://github.com/daarxwalker/hui#listbox) |
| `modal` | `Alpine.data('huiModal')` | [Modal](https://github.com/daarxwalker/hui#modal) |
| `radio` | `Alpine.data('huiRadio')` | [Radio](https://github.com/daarxwalker/hui#checkbox-switch-and-radio) |
| `range` | `Alpine.data('huiRange')` | [Range](https://github.com/daarxwalker/hui#range) |
| `switch` | `Alpine.data('huiSwitch')` | [Switch](https://github.com/daarxwalker/hui#checkbox-switch-and-radio) |
| `tabs` | `Alpine.data('huiTabs')` | [Tabs](https://github.com/daarxwalker/hui#tabs) |
| `toast` | `Alpine.data('huiToast')` | [Toast](https://github.com/daarxwalker/hui#toast) |
| `tooltip` | `Alpine.data('huiTooltip')` | [Tooltip](https://github.com/daarxwalker/hui#tooltip) |

The overlays (`dropdown`, `modal`, `tooltip`, `listbox`, `combobox`) share one
stack of open layers for Escape, clicks outside and `inert`. Each built file carries its own
copy of that code, and the first one to load keeps the stack on the `Alpine` object, so
loading several of them still gives one stack.

`collapsible` bundles [`@alpinejs/collapse`](https://alpinejs.dev/plugins/collapse) and
registers it once, because its content uses `x-collapse` by default. Alpine itself is never
bundled; the app loads it.
