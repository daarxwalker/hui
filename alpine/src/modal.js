import layers, { commands, dispatch, process, restoreFocus, tabbables, whenVisible } from './layers.js'

export default function modal(Alpine) {
	const stack = layers(Alpine)
	let locks = 0
	let overflow = ''

	const lock = () => {
		if (locks++ === 0) {
			overflow = document.documentElement.style.overflow
			document.documentElement.style.overflow = 'hidden'
		}
	}

	const unlock = () => {
		if (--locks === 0) document.documentElement.style.overflow = overflow
	}

	Alpine.data('huiModal', (config = {}) => ({
		open: config.open === true,
		active: false,
		returnTo: null,
		layer: null,
		cleanup: null,

		init() {
			const stored = this.restore()
			if (stored !== null) this.open = stored

			this.layer = {
				modal: true,
				parts: () => [this.overlay(), this.content()],
				inside: () => [this.content()],
				dismiss: () => this.dismiss(),
				keydown: (event) => this.keydown(event),
			}
			this.cleanup = commands(this, 'hui-modal')

			this.$watch('open', () => this.sync())
			this.$nextTick(() => this.sync())
		},

		destroy() {
			this.cleanup()
			this.deactivate(false)
		},

		show() {
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

		modalShow() {
			this.show()
		},

		modalHide() {
			this.hide()
		},

		modalToggle() {
			this.toggle()
		},

		dismiss() {
			if (config.static || config.alert) return false
			this.hide()
		},

		sync() {
			if (this.open === this.active) return
			this.open ? this.activate() : this.deactivate(true)
		},

		activate() {
			this.active = true
			this.returnTo = document.activeElement
			if (!config.allowScroll) lock()
			stack.open(this.layer)
			process(this.layer.parts())
			whenVisible(this.content(), () => this.focusIn())
			this.save(true)
			dispatch(this.$root, 'hui-modal:opened')
		},

		deactivate(notify) {
			if (!this.active) return
			this.active = false
			stack.close(this.layer)
			if (!config.allowScroll) unlock()
			if (!notify) return
			restoreFocus(this.content(), [this.returnTo, this.trigger()])
			this.save(false)
			dispatch(this.$root, 'hui-modal:closed')
		},

		content() {
			return document.getElementById(this.$root.id + '-content')
		},

		overlay() {
			return document.getElementById(this.$root.id + '-overlay')
		},

		trigger() {
			return this.$root.querySelector('[aria-controls="' + this.$root.id + '-content"]')
		},

		focusIn() {
			const content = this.content()
			if (!this.active || content.contains(document.activeElement)) return
			const target = content.querySelector('[autofocus]') || tabbables(content)[0] || content
			target.focus()
		},

		keydown(event) {
			const content = this.content()

			if (event.key === 'Escape') {
				if (config.static) return
				event.preventDefault()
				this.hide()
				return
			}

			if (event.key !== 'Tab' || !content) return
			const items = tabbables(content)
			const current = document.activeElement
			if (!items.length) {
				event.preventDefault()
				content.focus()
				return
			}
			const first = items[0]
			const last = items[items.length - 1]
			if (!content.contains(current)) {
				event.preventDefault()
				;(event.shiftKey ? last : first).focus()
			} else if (event.shiftKey && (current === first || current === content)) {
				event.preventDefault()
				last.focus()
			} else if (!event.shiftKey && current === last) {
				event.preventDefault()
				first.focus()
			}
		},

		restore() {
			if (!config.persist) return null
			try {
				const value = sessionStorage.getItem('hui-modal:' + config.persist)
				return value === null ? null : value === '1'
			} catch (error) {
				return null
			}
		},

		save(open) {
			if (!config.persist) return
			try {
				sessionStorage.setItem('hui-modal:' + config.persist, open ? '1' : '0')
			} catch (error) {}
		},
	}))
}
