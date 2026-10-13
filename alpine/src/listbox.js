import options from './options.js'

export default function listbox(Alpine) {
	Alpine.data('huiListbox', (config = {}) => ({
		...options(Alpine, 'listbox', config),

		init() {
			this.setup(() => [this.trigger(), this.content()])
		},

		destroy() {
			this.teardown()
		},

		keydown(event) {
			if (config.searchable) return this.triggerKeydown(event)
			const printable = event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey

			if (!this.open) {
				if (['ArrowDown', 'ArrowUp', 'Enter', ' '].includes(event.key)) {
					event.preventDefault()
					this.show('selected')
				} else if (event.key === 'Home' || event.key === 'End') {
					event.preventDefault()
					this.show(event.key === 'Home' ? 'first' : 'last')
				} else if (printable) {
					this.show(null)
					this.$nextTick(() => this.typeahead(event.key))
				}
				return
			}

			switch (event.key) {
				case 'ArrowDown':
					event.preventDefault()
					event.altKey ? this.hide() : this.move(1)
					break
				case 'ArrowUp':
					event.preventDefault()
					this.move(-1)
					break
				case 'Home':
					event.preventDefault()
					this.highlight('first')
					break
				case 'End':
					event.preventDefault()
					this.highlight('last')
					break
				case 'Enter':
					event.preventDefault()
					if (this.current()) this.choose(this.current())
					break
				case ' ':
					event.preventDefault()
					if (this.query) this.typeahead(event.key)
					else if (this.current()) this.choose(this.current())
					break
				case 'Tab':
					this.hide()
					break
				default:
					if (printable) this.typeahead(event.key)
			}
		},
	}))
}
