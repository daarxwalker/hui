import checkable from './checkable.js'

export default function radio(Alpine) {
	Alpine.data('huiRadio', (config = {}) => ({
		...checkable(config),
		listener: null,

		init() {
			this.setup()
			this.listener = (event) => this.selected(event.detail)
			window.addEventListener('hui-radio:select', this.listener)
			this.$nextTick(() => this.settle())
		},

		destroy() {
			this.teardown()
			window.removeEventListener('hui-radio:select', this.listener)
		},

		restore() {
			this.checked = this.initial.checked
			this.write()
			this.$nextTick(() => this.settle())
		},

		settle() {
			this.fill(this.group().some((element) => element.getAttribute('aria-checked') === 'true'))
			this.rove()
		},

		group() {
			if (!config.name) return [this.$root]
			const form = this.$root.form
			return Array.from(document.querySelectorAll('[role=radio][data-hui-radio]')).filter(
				(element) => element.dataset.huiRadio === config.name && element.form === form,
			)
		},

		enabled() {
			return this.group().filter((element) => !element.disabled)
		},

		select() {
			if (config.disabled || this.checked) return
			this.checked = true
			this.write()
			this.fill(true)
			this.rove(this.$root)
			window.dispatchEvent(new CustomEvent('hui-radio:select', { detail: { name: config.name, root: this.$root } }))
			this.emit()
		},

		selected(detail) {
			if (!config.name || detail.name !== config.name || detail.root === this.$root) return
			if (detail.root.form !== this.$root.form) return
			if (this.checked) {
				this.checked = false
				this.write()
			}
			this.fill(true)
			this.rove(detail.root)
		},

		validate() {},

		fill(filled) {
			const required = this.validation()
			if (required) required.value = filled ? '1' : ''
		},

		rove(checked = null) {
			if (!config.name) return
			const group = this.enabled()
			const current = checked || group.find((element) => element.getAttribute('aria-checked') === 'true')
			this.$root.tabIndex = this.$root === (current || group[0]) ? 0 : -1
		},

		keydown(event) {
			const keys = { ArrowDown: 1, ArrowRight: 1, ArrowUp: -1, ArrowLeft: -1 }
			if (!(event.key in keys)) return
			event.preventDefault()
			const group = this.enabled()
			const index = group.indexOf(this.$root)
			if (index === -1 || group.length < 2) return
			const next = group[(index + keys[event.key] + group.length) % group.length]
			next.focus()
			next.click()
		},
	}))
}
