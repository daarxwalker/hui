import collapse from '@alpinejs/collapse'

export const interactive = 'a[href], button, input, select, textarea, label, summary, [tabindex], [contenteditable]'

export function headerClicked(event) {
	const target = event.target instanceof Element ? event.target.closest(interactive) : null
	if (target && event.currentTarget.contains(target)) return false
	const selection = window.getSelection && window.getSelection()
	return !(selection && selection.toString())
}

export default function collapsing(Alpine) {
	if (Alpine.huiCollapse) return
	Alpine.huiCollapse = true
	Alpine.plugin(collapse)
}
