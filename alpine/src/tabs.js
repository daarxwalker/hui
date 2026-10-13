import { dispatch } from './layers.js'

export default function tabs(Alpine) {
	Alpine.data('huiTabs', (config = {}) => ({
		value: config.value ?? '',

		init() {
			this.$watch('value', (value) => dispatch(this.$root, 'hui-tabs:change', { value }))
		},

		select(value) {
			this.value = value
		},

		triggers() {
			const prefix = this.$root.id + '-trigger-'
			return Array.from(this.$root.querySelectorAll('[role=tab]')).filter(
				(trigger) => trigger.id.startsWith(prefix) && !trigger.disabled,
			)
		},

		keydown(event) {
			const vertical = this.$root.dataset.orientation === 'vertical'
			const keys = vertical ? ['ArrowUp', 'ArrowDown'] : ['ArrowLeft', 'ArrowRight']
			const triggers = this.triggers()
			const index = triggers.indexOf(event.target)
			if (index === -1) return

			let next = null
			if (event.key === keys[0]) next = triggers[(index - 1 + triggers.length) % triggers.length]
			else if (event.key === keys[1]) next = triggers[(index + 1) % triggers.length]
			else if (event.key === 'Home') next = triggers[0]
			else if (event.key === 'End') next = triggers[triggers.length - 1]
			if (!next) return

			event.preventDefault()
			next.focus()
			if (!config.manual) this.select(next.dataset.value)
		},
	}))
}
