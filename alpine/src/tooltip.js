import layers, { addToken, anchorName, commands, dispatch, moveAnchor } from './layers.js'

export default function tooltip(Alpine) {
	const stack = layers(Alpine)

	Alpine.data('huiTooltip', (config = {}) => ({
		open: config.open === true,
		active: false,
		blocked: false,
		timer: null,
		anchor: null,
		handled: null,
		layer: null,
		cleanup: null,

		init() {
			this.attach(this.$root.querySelector('[aria-describedby~="' + this.$root.id + '-content"]'))
			this.layer = {
				parts: () => [this.content()],
				inside: () => [this.anchor, this.content()],
				dismiss: () => this.hide(),
				keydown: (event) => this.escape(event),
			}
			this.cleanup = commands(this, 'hui-tooltip')
			this.$watch('open', () => this.sync())
			this.$nextTick(() => this.sync())
		},

		destroy() {
			this.cleanup()
			clearTimeout(this.timer)
			this.deactivate()
		},

		content() {
			return document.getElementById(this.$root.id + '-content')
		},

		owns(target) {
			return Boolean(target) && [this.anchor, this.content()].some((element) => element && element.contains(target))
		},

		attach(element) {
			const content = this.content()
			if (!element || element === this.anchor || (content && content.contains(element))) return
			moveAnchor(this.anchor, element, anchorName(this.$root.id))
			addToken(element, 'aria-describedby', this.$root.id + '-content')
			this.anchor = element
		},

		seen(event) {
			if (!event) return false
			if (this.handled === event) return true
			this.handled = event
			return false
		},

		show() {
			clearTimeout(this.timer)
			this.open = true
			this.sync()
		},

		hide() {
			clearTimeout(this.timer)
			this.open = false
			this.sync()
		},

		toggle() {
			this.open ? this.hide() : this.show()
		},

		dismiss() {
			this.blocked = true
			this.hide()
		},

		later(callback, delay) {
			clearTimeout(this.timer)
			this.timer = setTimeout(callback, delay)
		},

		sync() {
			if (this.open === this.active) return
			this.active = this.open
			if (this.open) stack.open(this.layer)
			else stack.close(this.layer)
			dispatch(this.$root, this.open ? 'hui-tooltip:opened' : 'hui-tooltip:closed')
		},

		deactivate() {
			if (this.active) stack.close(this.layer)
			this.active = false
		},

		tooltipShow(element) {
			this.attach(element)
			this.show()
		},

		tooltipHide() {
			this.hide()
		},

		tooltipToggle(element, event) {
			if (this.seen(event)) return
			if (this.open && (!element || element === this.anchor)) this.hide()
			else this.tooltipShow(element)
		},

		tooltipDismiss() {
			this.dismiss()
		},

		tooltipEnter(event, element) {
			if (event.pointerType === 'touch' || this.blocked) return
			this.attach(element)
			if (this.open) {
				clearTimeout(this.timer)
				return
			}
			this.later(() => this.show(), config.delay ?? 700)
		},

		tooltipLeave(event) {
			if (event.pointerType === 'touch' || this.owns(event.relatedTarget)) return
			this.blocked = false
			if (!this.open) {
				clearTimeout(this.timer)
				return
			}
			this.later(() => this.hide(), config.closeDelay ?? 300)
		},

		tooltipFocus(event, element) {
			if (this.blocked || !event.target.matches(':focus-visible')) return
			this.tooltipShow(element)
		},

		tooltipBlur(event) {
			if (this.owns(event.relatedTarget)) return
			this.blocked = false
			this.hide()
		},

		escape(event) {
			if (event.key !== 'Escape') return
			event.preventDefault()
			this.dismiss()
		},
	}))
}
