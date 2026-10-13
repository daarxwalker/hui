# @daarxwalker/headlessui-alpine

The [Alpine.js](https://alpinejs.dev) plugins behind
[hui](https://github.com/daarxwalker/hui), a headless component library for
[templ](https://templ.guide). The Go components already embed these files. This package is
for when you bundle your own JavaScript or load Alpine from a CDN.

## Install

```bash
npm install alpinejs @daarxwalker/headlessui-alpine
```

```js
import Alpine from 'alpinejs'
import headlessui from '@daarxwalker/headlessui-alpine'

Alpine.plugin(headlessui)
Alpine.start()
```

Each component is also exported on its own:

```js
import { modal } from '@daarxwalker/headlessui-alpine'

Alpine.plugin(modal)
```

## CDN

Load the plugin before Alpine. It registers itself on `alpine:init`, and if Alpine has
already started (for example after an htmx swap), it registers straight away.

```html
<script defer src="https://cdn.jsdelivr.net/npm/@daarxwalker/headlessui-alpine@0.1.0/dist/cdn.min.js"></script>
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3/dist/cdn.min.js"></script>
```

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
registers it once, because its content uses `x-collapse` by default.

## Build

```bash
npm install
npm run build
```

| File | For |
|---|---|
| `dist/module.esm.js` | `import` from a bundler |
| `dist/cdn.min.js` | a `<script>` tag with every plugin |
| `dist/<component>.min.js` | a single self-registering plugin, embedded by the Go components |
