export default function range(Alpine) {
	Alpine.data('huiRange', () => ({
		values: [],
		form: null,
		reset: null,

		init() {
			const thumb = this.thumbs()[0]
			this.form = thumb ? thumb.form : null
			this.reset = () => setTimeout(() => this.sync())
			if (this.form) this.form.addEventListener('reset', this.reset)
			this.sync()
		},

		destroy() {
			if (this.form) this.form.removeEventListener('reset', this.reset)
		},

		thumbs() {
			return Array.from(this.$root.querySelectorAll('input[type="range"]'))
		},

		update(thumb) {
			const thumbs = this.thumbs()
			const index = thumbs.indexOf(thumb)
			if (index === -1) return
			const value = Number(thumb.value)
			const low = index > 0 ? Number(thumbs[index - 1].value) : -Infinity
			const high = index < thumbs.length - 1 ? Number(thumbs[index + 1].value) : Infinity
			const clamped = Math.min(Math.max(value, low), high)
			if (clamped !== value) thumb.value = clamped
			this.sync()
		},

		sync() {
			const thumbs = this.thumbs()
			const percents = thumbs.map((thumb) => this.percent(thumb))
			this.values = thumbs.map((thumb) => Number(thumb.value))

			const start = percents.length > 1 ? percents[0] : 0
			const end = percents.length > 0 ? percents[percents.length - 1] : 0
			this.$root.style.setProperty('--hui-range-start', start + '%')
			this.$root.style.setProperty('--hui-range-end', end + '%')

			if (thumbs.length < 2) return
			thumbs.forEach((thumb, index) => {
				thumb.style.zIndex = percents[index] > 50 ? thumbs.length - index : index + 1
			})
		},

		percent(thumb) {
			const min = Number(thumb.min)
			const max = Number(thumb.max)
			if (max <= min) return 0
			return Math.round(((Number(thumb.value) - min) / (max - min)) * 100000) / 1000
		},
	}))
}
