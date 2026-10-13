import collapsing, { headerClicked } from './collapse.js'
import { commands, dispatch } from './layers.js'

export default function collapsible(Alpine) {
	collapsing(Alpine)

	Alpine.data('huiCollapsible', (config = {}) => ({
		open: config.open === true,
		cleanup: null,

		init() {
			this.cleanup = commands(this, 'hui-collapsible')
			this.$watch('open', (open) => dispatch(this.$root, open ? 'hui-collapsible:opened' : 'hui-collapsible:closed'))
		},

		destroy() {
			this.cleanup()
		},

		show() {
			this.open = true
		},

		hide() {
			this.open = false
		},

		toggle() {
			this.open = !this.open
		},

		collapsibleShow() {
			this.show()
		},

		collapsibleHide() {
			this.hide()
		},

		collapsibleToggle() {
			this.toggle()
		},

		collapsibleHeaderClick(event) {
			const trigger = document.getElementById(this.$root.id + '-trigger')
			if ((trigger && trigger.disabled) || !headerClicked(event)) return
			this.toggle()
		},
	}))
}
