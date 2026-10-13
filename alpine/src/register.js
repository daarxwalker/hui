import pending from './pending.js'

export default function register(plugin) {
	pending()
	if (window.Alpine) window.Alpine.plugin(plugin)
	else document.addEventListener('alpine:init', () => window.Alpine.plugin(plugin))
}
