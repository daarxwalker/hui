const key = Symbol.for('headlessui.pending')

function mark(element) {
	const state = element[key] || (element[key] = { count: 0, kept: element.hasAttribute('data-pending') })
	state.count++
	element.setAttribute('data-pending', '')
	element.setAttribute('aria-busy', 'true')
}

function unmark(element) {
	const state = element[key]
	if (!state || --state.count > 0) return
	delete element[key]
	if (state.kept) return
	element.removeAttribute('data-pending')
	element.removeAttribute('aria-busy')
}

function leaving(ctx) {
	return !!ctx.hx && (!!ctx.hx.redirect || ctx.hx.refresh === 'true')
}

export function markPending(elements, ctx) {
	const marked = Array.from(new Set(elements)).filter(
		(element) => element instanceof Element && element !== document.body && element !== document.documentElement,
	)
	if (!marked.length) return
	marked.forEach(mark)
	const done = (event) => {
		if (!event.detail || event.detail.ctx !== ctx) return
		document.removeEventListener('htmx:finally:request', done)
		if (!leaving(ctx)) marked.forEach(unmark)
	}
	document.addEventListener('htmx:finally:request', done)
}

export default function pending() {
	if (document[key]) return
	document[key] = true
	document.addEventListener('htmx:before:request', (event) => {
		const ctx = event.detail && event.detail.ctx
		if (!ctx || !(event.target instanceof Element)) return
		markPending([ctx.sourceElement || event.target, ctx.request && ctx.request.submitter, ctx.target], ctx)
	})
}
