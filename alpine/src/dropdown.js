import checkable from './checkable.js'
import layers, { anchorName, commands, dispatch, expand, moveAnchor, process, restoreFocus, whenVisible } from './layers.js'

const menuItems = '[role=menuitem], [role=menuitemcheckbox], [role=menuitemradio]'

export default function dropdown(Alpine) {
	const stack = layers(Alpine)

	Alpine.data('huiDropdownCheckbox', (config = {}) => ({
		...checkable(config),

		init() {
			this.setup()
		},

		destroy() {
			this.teardown()
		},

		dropdownCheck() {
			if (config.disabled) return
			this.checked = !this.checked
			this.write()
			this.emit()
		},
	}))

	Alpine.data('huiDropdownRadio', (config = {}) => ({
		...checkable(config),

		init() {
			this.setup()
		},

		destroy() {
			this.teardown()
		},

		dropdownPick() {
			if (config.disabled || this.checked) return
			this.checked = true
			this.write()
			const group = this.$root.closest('[role=group], [role=menu]')
			if (group) {
				group.querySelectorAll('[role=menuitemradio]').forEach((item) => {
					if (item === this.$root) return
					const data = Alpine.$data(item)
					if (!data.checked) return
					data.checked = false
					data.write()
				})
			}
			this.emit()
		},
	}))

	Alpine.data('huiDropdown', (config = {}) => ({
		open: config.open === true,
		active: false,
		focusOnOpen: 'content',
		query: '',
		queryTimer: null,
		anchor: null,
		layer: null,
		cleanup: null,

		init() {
			this.anchor = this.trigger()
			this.layer = {
				parts: () => [this.content()],
				inside: () => [this.anchor, this.content()],
				dismiss: () => this.hide(),
				keydown: (event) => this.escape(event),
			}
			this.cleanup = commands(this, 'hui-dropdown')
			this.$watch('open', () => this.sync())
			this.$nextTick(() => this.sync())
		},

		destroy() {
			this.cleanup()
			clearTimeout(this.queryTimer)
			this.deactivate(false)
		},

		trigger() {
			return document.getElementById(this.$root.id + '-trigger')
		},

		content() {
			return document.getElementById(this.$root.id + '-content')
		},

		items() {
			const content = this.content()
			if (!content) return []
			return Array.from(content.querySelectorAll(menuItems)).filter(
				(item) => !item.matches(':disabled, [aria-disabled=true]'),
			)
		},

		show(focus) {
			if (focus) this.focusOnOpen = focus
			this.open = true
			this.sync()
		},

		hide() {
			this.open = false
			this.sync()
		},

		toggle() {
			this.open ? this.hide() : this.show()
		},

		attach(element) {
			const content = this.content()
			if (!element || element === this.anchor || (content && content.contains(element))) return
			moveAnchor(this.anchor, element, anchorName(this.$root.id))
			if (this.anchor) expand(this.anchor, this.$root.id + '-content', false)
			this.anchor = element
			expand(element, this.$root.id + '-content', this.open)
		},

		dropdownShow(element, event) {
			this.attach(element)
			if (event && event.type === 'click') this.focusOnOpen = event.detail === 0 ? 'first' : 'content'
			this.show()
		},

		dropdownHide() {
			this.hide()
		},

		dropdownToggle(element, event) {
			if (this.open && (!element || element === this.anchor)) this.hide()
			else this.dropdownShow(element, event)
		},

		sync() {
			if (this.open === this.active) return
			this.open ? this.activate() : this.deactivate(true)
		},

		activate() {
			this.active = true
			expand(this.anchor, this.$root.id + '-content', true)
			const focus = this.focusOnOpen
			this.focusOnOpen = 'content'
			stack.open(this.layer)
			process([this.content()])
			whenVisible(this.content(), () => this.active && this.focus(focus))
			dispatch(this.$root, 'hui-dropdown:opened')
		},

		deactivate(notify) {
			if (!this.active) return
			this.active = false
			stack.close(this.layer)
			expand(this.anchor, this.$root.id + '-content', false)
			if (!notify) return
			restoreFocus(this.content(), [this.anchor])
			dispatch(this.$root, 'hui-dropdown:closed')
		},

		focus(which) {
			const items = this.items()
			const item = which === 'first' ? items[0] : which === 'last' ? items[items.length - 1] : null
			const target = item || this.content()
			if (target) target.focus({ preventScroll: true })
		},

		move(delta) {
			const items = this.items()
			if (!items.length) return
			const index = items.indexOf(document.activeElement)
			const next = index === -1 ? (delta > 0 ? 0 : items.length - 1) : (index + delta + items.length) % items.length
			items[next].focus()
		},

		typeahead(key) {
			clearTimeout(this.queryTimer)
			this.query += key.toLowerCase()
			this.queryTimer = setTimeout(() => (this.query = ''), 500)

			const items = this.items()
			const current = items.indexOf(document.activeElement)
			const start = Math.max(this.query.length > 1 ? current : current + 1, 0)
			const ordered = items.slice(start).concat(items.slice(0, start))
			const match = ordered.find((item) => item.textContent.trim().toLowerCase().startsWith(this.query))
			if (match) match.focus()
		},

		dropdownKeydown(event, element) {
			if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp') return
			event.preventDefault()
			const focus = event.key === 'ArrowDown' ? 'first' : 'last'
			this.attach(element)
			this.open ? this.focus(focus) : this.show(focus)
		},

		escape(event) {
			if (event.key !== 'Escape') return
			event.preventDefault()
			this.hide()
		},

		keydown(event) {
			switch (event.key) {
				case 'ArrowDown':
					event.preventDefault()
					this.move(1)
					break
				case 'ArrowUp':
					event.preventDefault()
					this.move(-1)
					break
				case 'Home':
					event.preventDefault()
					this.focus('first')
					break
				case 'End':
					event.preventDefault()
					this.focus('last')
					break
				case 'Tab':
					this.hide()
					break
				case ' ':
					if (this.query) {
						event.preventDefault()
						this.typeahead(event.key)
					} else if (event.target.matches('a[role=menuitem]')) {
						event.preventDefault()
						event.target.click()
					}
					break
				default:
					if (event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
						this.typeahead(event.key)
					}
			}
		},

		select(event) {
			const item = event.target.closest('[role=menuitem], [role=menuitemradio]')
			if (!item || event.defaultPrevented || !this.items().includes(item)) return
			this.hide()
		},

		hover(event) {
			const item = event.target.closest && event.target.closest(menuItems)
			if (!item || item === document.activeElement || !this.items().includes(item)) return
			item.focus({ preventScroll: true })
		},
	}))
}
