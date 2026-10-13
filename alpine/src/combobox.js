import options from './options.js'
import { dispatch } from './layers.js'

export default function combobox(Alpine) {
	Alpine.data('huiCombobox', (config = {}) => ({
		...options(Alpine, 'combobox', config),

		typed: false,

		init() {
			this.setup(() => [this.$root, this.content()])
		},

		destroy() {
			this.teardown()
		},

		comboboxToggle() {
			this.toggle()
			this.focusControl()
		},

		term() {
			const search = this.search()
			return this.typed && search ? search.value : ''
		},

		focusControl() {
			const search = this.search()
			if (search && document.activeElement !== search) search.focus()
		},

		seedLabel() {
			const search = this.search()
			if (config.multiple || !this.value || !search || !search.value) return
			this.labels[this.value] = search.value
		},

		input() {
			this.typed = true
			if (!this.open) this.show(null)
			this.active = null
			this.wantsNext = false
			this.ask()
		},

		ask() {
			const search = this.search()
			if (search && this.term() !== this.searched) dispatch(search, 'hui-combobox:query', { query: this.term() })
			this.updateSearching()
			this.mark()
		},

		inputClick() {
			if (this.open) return
			this.show('selected')
			const search = this.search()
			if (search && !this.typed) search.select()
		},

		inputBlur(event) {
			const content = this.content()
			if (!this.open || (content && content.contains(event.relatedTarget))) return
			this.hide()
		},

		keydown(event) {
			if (!this.open && event.key === 'ArrowDown' && !event.altKey) {
				event.preventDefault()
				this.show('selected')
				return
			}
			if (this.inputKeydown(event)) return
			if (event.key === 'Tab' && this.open) this.hide()
		},

		written(values) {
			const search = this.search()
			if (!search) return
			this.typed = false
			search.value = config.multiple || !values.length ? '' : this.labelOf(values[0]) ?? values[0]
			this.ask()
		},

		closed() {
			const search = this.search()
			if (search && this.typed && !config.multiple && search.value.trim() === '' && this.selected().length) {
				this.remove(this.value)
				return
			}
			this.written(this.selected())
		},
	}))
}
