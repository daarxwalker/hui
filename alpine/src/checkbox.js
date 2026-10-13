import checkable from './checkable.js'

export default function checkbox(Alpine) {
	Alpine.data('huiCheckbox', (config = {}) => ({
		...checkable(config),

		init() {
			this.setup()
		},

		destroy() {
			this.teardown()
		},

		toggle() {
			if (config.disabled) return
			this.checked = this.indeterminate ? true : !this.checked
			this.indeterminate = false
			this.write()
			this.emit()
		},

		keydown(event) {
			if (event.key === 'Enter') event.preventDefault()
		},
	}))
}
