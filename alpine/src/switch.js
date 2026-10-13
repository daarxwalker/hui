import checkable from './checkable.js'

export default function toggleSwitch(Alpine) {
	Alpine.data('huiSwitch', (config = {}) => ({
		...checkable(config),

		init() {
			this.setup()
		},

		destroy() {
			this.teardown()
		},

		toggle() {
			if (config.disabled) return
			this.checked = !this.checked
			this.write()
			this.emit()
		},
	}))
}
