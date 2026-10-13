export default function checkable(config = {}) {
	return {
		checked: config.checked === true,
		indeterminate: config.indeterminate === true,
		initial: null,
		owner: null,
		reset: null,

		setup() {
			this.initial = { checked: this.checked, indeterminate: this.indeterminate }
			const required = this.validation()
			if (required) required.addEventListener('focus', () => this.$root.focus())
			this.owner = this.$root.form
			if (this.owner) {
				this.reset = () => setTimeout(() => this.restore())
				this.owner.addEventListener('reset', this.reset)
			}
			this.$watch('checked', () => this.write())
			this.$watch('indeterminate', () => this.write())
		},

		teardown() {
			if (this.owner && this.reset) this.owner.removeEventListener('reset', this.reset)
		},

		restore() {
			this.checked = this.initial.checked
			this.indeterminate = this.initial.indeterminate
			this.write()
		},

		ariaChecked() {
			return this.indeterminate ? 'mixed' : String(this.checked)
		},

		validation() {
			const next = this.$root.nextElementSibling
			return next && next.matches('input[data-hui-require]') ? next : null
		},

		input() {
			return this.$root.querySelector(':scope > input[data-hui-checked]')
		},

		write() {
			if (config.name && this.checked) {
				this.$root.setAttribute('name', config.name)
				this.$root.setAttribute('value', config.value ?? 'on')
			} else {
				this.$root.removeAttribute('name')
				this.$root.removeAttribute('value')
			}
			const current = this.input()
			if (config.name && this.checked && !current) {
				const input = document.createElement('input')
				input.type = 'hidden'
				input.name = config.name
				input.value = config.value ?? 'on'
				input.dataset.huiChecked = ''
				if (config.form) input.setAttribute('form', config.form)
				if (config.disabled) input.disabled = true
				this.$root.prepend(input)
			} else if (!this.checked && current) {
				current.remove()
			}
			this.validate()
		},

		validate() {
			const required = this.validation()
			if (required) required.value = this.checked ? '1' : ''
		},

		emit() {
			this.$root.dispatchEvent(new Event('change', { bubbles: true }))
		},
	}
}
