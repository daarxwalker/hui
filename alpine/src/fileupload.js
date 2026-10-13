function accepts(file, accept) {
	if (!accept) return true
	const name = file.name.toLowerCase()
	const type = (file.type || '').toLowerCase()
	return accept
		.split(',')
		.map((rule) => rule.trim().toLowerCase())
		.filter(Boolean)
		.some((rule) => {
			if (rule.startsWith('.')) return name.endsWith(rule)
			if (rule.endsWith('/*')) return type.startsWith(rule.slice(0, -1))
			return type === rule
		})
}

export default function fileUpload(Alpine) {
	Alpine.data('huiFileUpload', () => ({
		dragging: false,
		depth: 0,
		files: [],

		input() {
			return this.$root.querySelector('input[type=file][data-hui-file="' + this.$root.id + '"]')
		},

		changed(event) {
			this.files = Array.from(event.target.files || []).map((file) => file.name)
		},

		browse(event) {
			const input = this.input()
			if (!input || input.disabled || event.target === input || event.target.closest('a, button, input, label')) return
			input.click()
		},

		enter(event) {
			if (!this.input() || this.input().disabled) return
			this.depth++
			this.dragging = true
		},

		over(event) {
			if (event.dataTransfer) event.dataTransfer.dropEffect = this.input() && !this.input().disabled ? 'copy' : 'none'
		},

		leave() {
			this.depth = Math.max(this.depth - 1, 0)
			if (this.depth === 0) this.dragging = false
		},

		drop(event) {
			this.depth = 0
			this.dragging = false
			const input = this.input()
			if (!input || input.disabled || !event.dataTransfer) return

			let files = Array.from(event.dataTransfer.files).filter((file) => accepts(file, input.accept))
			if (!input.multiple) files = files.slice(0, 1)
			if (!files.length) return

			const transfer = new DataTransfer()
			files.forEach((file) => transfer.items.add(file))
			input.files = transfer.files
			input.dispatchEvent(new Event('input', { bubbles: true }))
			input.dispatchEvent(new Event('change', { bubbles: true }))
		},
	}))
}
