import { dispatch } from './layers.js'

export default function toast(Alpine) {
	Alpine.data('huiToast', (config = {}) => ({
		open: true,
		timer: null,
		started: 0,
		remaining: config.duration ?? 5000,
		paused: false,
		visibility: null,

		init() {
			this.visibility = () => (document.hidden ? this.pause() : this.resume())
			document.addEventListener('visibilitychange', this.visibility)
			this.resume()
		},

		destroy() {
			clearTimeout(this.timer)
			document.removeEventListener('visibilitychange', this.visibility)
		},

		pause() {
			if (!this.remaining || this.paused) return
			this.paused = true
			clearTimeout(this.timer)
			this.remaining = Math.max(this.remaining - (Date.now() - this.started), 0)
		},

		resume(event) {
			if (event && event.relatedTarget && this.$root.contains(event.relatedTarget)) return
			if (!this.remaining || document.hidden || this.$root.matches(':hover')) return
			this.paused = false
			this.started = Date.now()
			clearTimeout(this.timer)
			this.timer = setTimeout(() => this.dismiss(), this.remaining)
		},

		dismiss() {
			if (!this.open) return
			clearTimeout(this.timer)
			const root = this.$root
			const toaster = root.closest('[data-hui-toaster]')
			const focused = root.contains(document.activeElement)
			this.open = false
			dispatch(root, 'hui-toast:closed')
			if (focused && toaster) {
				const next = Array.from(toaster.children).find((item) => item !== root && item.getClientRects().length)
				if (next) next.focus()
			}
			this.removeWhenHidden(root, 0)
		},

		removeWhenHidden(root, attempt) {
			if (getComputedStyle(root).display === 'none' || attempt > 40) {
				root.remove()
				return
			}
			setTimeout(() => this.removeWhenHidden(root, attempt + 1), 50)
		},

		keydown(event) {
			if (event.key !== 'Escape') return
			event.preventDefault()
			event.stopPropagation()
			this.dismiss()
		},
	}))
}
