const tabbable =
	'a[href], area[href], button, input, select, textarea, iframe, summary, [contenteditable]:not([contenteditable=false]), [tabindex]'

export function tabbables(root) {
	if (!root) return []
	return Array.from(root.querySelectorAll(tabbable)).filter(
		(element) =>
			element.tabIndex >= 0 &&
			!element.matches(':disabled') &&
			!element.closest('[inert]') &&
			element.getClientRects().length > 0,
	)
}

export function idSafe(value) {
	return Array.from(String(value), (char) => (/[a-zA-Z0-9_-]/.test(char) ? char : '-')).join('')
}

export function anchorName(id) {
	return '--' + idSafe(id)
}

export function moveAnchor(from, to, name) {
	if (from === to) return
	if (from) from.style.removeProperty('anchor-name')
	if (to) to.style.setProperty('anchor-name', name)
}

export function expand(element, controls, open) {
	if (!element) return
	element.setAttribute('aria-controls', controls)
	element.setAttribute('aria-expanded', String(open))
	element.toggleAttribute('data-open', open)
}

export function addToken(element, attribute, token) {
	const tokens = (element.getAttribute(attribute) || '').split(/\s+/).filter(Boolean)
	if (!tokens.includes(token)) element.setAttribute(attribute, tokens.concat(token).join(' '))
}

export function removeToken(element, attribute, token) {
	const tokens = (element.getAttribute(attribute) || '').split(/\s+/).filter((item) => item && item !== token)
	if (tokens.length) element.setAttribute(attribute, tokens.join(' '))
	else element.removeAttribute(attribute)
}

export function dispatch(root, name, detail = {}) {
	root.dispatchEvent(new CustomEvent(name, { bubbles: true, detail: { id: root.id, ...detail } }))
}

export function commands(component, prefix, names = ['show', 'hide', 'toggle']) {
	const handler = (event) => {
		const target = event.detail && (event.detail.id || event.detail)
		if (target !== component.$root.id) return
		component[event.type.slice(prefix.length + 1)]()
	}
	names.forEach((name) => window.addEventListener(prefix + ':' + name, handler))
	return () => names.forEach((name) => window.removeEventListener(prefix + ':' + name, handler))
}

export function whenVisible(element, callback, attempt = 0) {
	if (!element) return
	if (element.getClientRects().length) {
		callback()
		return
	}
	if (attempt < 20) setTimeout(() => whenVisible(element, callback, attempt + 1), 16)
}

export function process(elements) {
	if (!window.htmx) return
	elements.forEach((element) => element && window.htmx.process(element))
}

export function restoreFocus(content, candidates) {
	const current = document.activeElement
	if (content && current && current !== document.body && !content.contains(current)) return
	const target = candidates.find((element) => element && element !== document.body && element.isConnected)
	if (target) target.focus()
}

export function nextTabbable(from) {
	const items = tabbables(document.body)
	const index = items.indexOf(from)
	return index === -1 ? null : items[index + 1] || null
}

export default function layers(Alpine) {
	if (Alpine.huiLayers) return Alpine.huiLayers

	const stack = []
	let inerted = []

	const inside = (layer, target) => layer.inside().some((element) => element && element.contains(target))

	const release = () => {
		inerted.forEach((element) => (element.inert = false))
		inerted = []
	}

	const isolate = () => {
		release()
		let index = stack.length - 1
		while (index >= 0 && !stack[index].modal) index--
		if (index === -1) return

		const keep = new Set()
		stack.slice(index).forEach((layer) => layer.parts().forEach((element) => element && keep.add(element)))
		document.querySelectorAll('[data-hui-toaster], [data-hui-announcer]').forEach((element) => keep.add(element))

		const path = new Set()
		keep.forEach((element) => {
			let node = element
			while (node && node !== document.body && !path.has(node)) {
				path.add(node)
				node = node.parentElement
			}
		})
		path.forEach((node) => {
			if (!node.parentElement) return
			for (const sibling of node.parentElement.children) {
				if (path.has(sibling) || sibling.inert || sibling.matches('script, style, template')) continue
				sibling.inert = true
				inerted.push(sibling)
			}
		})
	}

	document.addEventListener('keydown', (event) => {
		const top = stack[stack.length - 1]
		if (top && !event.defaultPrevented) top.keydown(event)
	})

	document.addEventListener(
		'click',
		(event) => {
			const target = event.target
			if (!stack.length || !(target instanceof Element) || !target.isConnected) return
			if (target.closest('[data-hui-toaster]')) return
			for (let index = stack.length - 1; index >= 0; index--) {
				const layer = stack[index]
				if (inside(layer, target)) break
				if (layer.dismiss(event) === false || layer.modal) break
			}
		},
		true,
	)

	Alpine.huiLayers = {
		open(layer) {
			if (stack.includes(layer)) return
			stack.push(layer)
			isolate()
		},

		close(layer) {
			const index = stack.indexOf(layer)
			if (index === -1) return
			stack.splice(index, 1)
			isolate()
		},

		top(layer) {
			return stack[stack.length - 1] === layer
		},
	}
	return Alpine.huiLayers
}
