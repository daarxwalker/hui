import layers, { commands, dispatch, idSafe, nextTabbable, process, restoreFocus, whenVisible } from './layers.js'
import { markPending } from './pending.js'

function textOf(option) {
	if (option.dataset.label) return option.dataset.label
	const clone = option.cloneNode(true)
	clone.querySelectorAll('[aria-hidden=true]').forEach((element) => element.remove())
	return clone.textContent.trim()
}

function fold(text) {
	return text.normalize('NFD').replace(/\p{M}/gu, '').toLowerCase()
}

function folded(text) {
	let value = ''
	const starts = []
	const ends = []
	for (let at = 0; at < text.length; ) {
		const char = String.fromCodePoint(text.codePointAt(at))
		const plain = fold(char)
		for (let i = 0; i < plain.length; i++) {
			starts.push(at)
			ends.push(at + char.length)
		}
		value += plain
		at += char.length
	}
	return { value, starts, ends }
}

function highlight(option, needle) {
	const hidden = (node) => {
		const element = node.parentElement && node.parentElement.closest('[aria-hidden=true]')
		return !!element && option.contains(element)
	}
	const walker = document.createTreeWalker(option, NodeFilter.SHOW_TEXT, {
		acceptNode: (node) => (hidden(node) ? NodeFilter.FILTER_REJECT : NodeFilter.FILTER_ACCEPT),
	})
	const nodes = []
	while (walker.nextNode()) nodes.push(walker.currentNode)
	nodes.forEach((node) => {
		const { value, starts, ends } = folded(node.data)
		const ranges = []
		for (let at = value.indexOf(needle); at !== -1; at = value.indexOf(needle, at + needle.length)) {
			ranges.push([starts[at], ends[at + needle.length - 1]])
		}
		ranges.reverse().forEach(([start, end]) => {
			const match = node.splitText(start)
			match.splitText(end - start)
			const mark = document.createElement('mark')
			mark.dataset.huiMatch = ''
			match.replaceWith(mark)
			mark.append(match)
		})
	})
}

export default function options(Alpine, kind, config) {
	const stack = layers(Alpine)
	const prefix = 'hui-' + kind

	return {
		[kind + 'Show']() {
			this.show('selected')
			this.shown && this.focusControl()
		},

		[kind + 'Hide']() {
			this.hide()
		},

		[kind + 'Toggle']() {
			this.toggle()
		},

		open: config.open === true,
		value: config.multiple ? (config.value || []).slice() : config.value ?? '',
		active: null,
		shown: false,
		pending: null,
		layer: null,
		cleanup: null,
		observer: null,
		requests: null,
		labels: {},
		wantsNext: false,
		query: '',
		queryTimer: null,
		searching: false,
		searched: null,

		setup(inside) {
			this.layer = {
				parts: () => [this.content()],
				inside,
				dismiss: () => this.hide(),
				keydown: (event) => this.escape(event),
			}
			this.cleanup = commands(this, prefix)
			this.$watch('open', () => this.sync())
			this.$watch('value', () => this.write())
			this.$watch('active', () => this.paint())
			this.$watch('searching', () => this.paintSearching())
			this.watchRequests()
			this.$nextTick(() => {
				const content = this.content()
				if (content) {
					this.observer = new MutationObserver(() => this.refreshed())
					this.observer.observe(content, { childList: true, subtree: true })
				}
				this.searched = config.searchable ? this.term() : null
				this.seedLabel()
				this.remember()
				this.paint()
				this.mark()
				this.sync()
			})
		},

		teardown() {
			this.cleanup()
			this.requests()
			clearTimeout(this.queryTimer)
			if (this.observer) this.observer.disconnect()
			this.deactivate()
		},

		watchRequests() {
			const mine = (event) => {
				const target = event.target
				if (!(target instanceof Element)) return false
				const content = this.content()
				return (content && content.contains(target)) || target === this.control()
			}
			const start = (event) => {
				if (!mine(event)) return
				markPending([this.list()], event.detail.ctx)
				if (event.target === this.search()) this.trackSearch(event.detail.ctx)
			}
			const configure = (event) => {
				const search = this.search()
				const name = search && search.dataset.huiQuery
				if (!name || event.target !== search) return
				event.detail.ctx.request.body.set(name, this.term())
			}
			document.addEventListener('htmx:config:request', configure)
			document.addEventListener('htmx:before:request', start)
			this.requests = () => {
				document.removeEventListener('htmx:config:request', configure)
				document.removeEventListener('htmx:before:request', start)
			}
		},

		trackSearch(ctx) {
			const query = this.term()
			const done = (event) => {
				if (!event.detail || event.detail.ctx !== ctx) return
				document.removeEventListener('htmx:finally:request', done)
				this.searched = query
				this.updateSearching()
			}
			document.addEventListener('htmx:finally:request', done)
		},

		updateSearching() {
			const search = this.search()
			const remote = !!search && Array.from(search.attributes).some((attribute) => /^(data-)?hx-(get|post|put|patch|delete)$/.test(attribute.name))
			this.searching = remote && this.term() !== this.searched
		},

		term() {
			const search = config.searchable && this.search()
			return search ? search.value : ''
		},

		mark() {
			const list = this.list()
			if (!list) return
			list.querySelectorAll('mark[data-hui-match]').forEach((mark) => {
				const parent = mark.parentNode
				mark.replaceWith(...mark.childNodes)
				parent.normalize()
			})
			const needle = fold(this.term().trim())
			if (needle) this.options(true).forEach((option) => highlight(option, needle))
			if (this.observer) this.observer.takeRecords()
		},

		paintSearching() {
			;[this.search(), this.list()].forEach((element) => element && element.toggleAttribute('data-searching', this.searching))
		},

		listId() {
			return this.$root.id + '-list'
		},

		control() {
			if (config.searchable) return this.search()
			return document.querySelector('[role=combobox][aria-controls="' + this.listId() + '"]')
		},

		trigger() {
			return this.$root.querySelector('[aria-controls="' + this.listId() + '"]')
		},

		search() {
			if (kind === 'combobox') return this.$root.querySelector('input[role=combobox]')
			return document.getElementById(this.$root.id + '-search')
		},

		content() {
			return document.getElementById(this.$root.id + '-content')
		},

		list() {
			return document.getElementById(this.listId()) || this.content()
		},

		loader() {
			const list = this.list()
			return list ? list.querySelector('[data-hui-loader]') : null
		},

		focusControl() {
			const control = config.searchable ? this.trigger() : this.control()
			if (control) control.focus()
		},

		options(all = false) {
			const list = this.list()
			if (!list) return []
			const options = Array.from(list.querySelectorAll('[role=option]'))
			return all ? options : options.filter((option) => option.getAttribute('aria-disabled') !== 'true')
		},

		selected() {
			if (config.multiple) return Array.isArray(this.value) ? this.value : []
			return this.value === '' || this.value == null ? [] : [this.value]
		},

		isSelected(value) {
			return this.selected().includes(value)
		},

		show(which = null) {
			if (config.disabled) return
			this.pending = which
			this.open = true
			this.sync()
		},

		hide() {
			this.open = false
			this.sync()
		},

		toggle() {
			this.open ? this.hide() : this.show('selected')
		},

		sync() {
			if (this.open === this.shown) return
			this.open ? this.activate() : this.deactivate(true)
		},

		activate() {
			this.shown = true
			stack.open(this.layer)
			const content = this.content()
			process([content])
			const which = this.pending
			this.pending = null
			whenVisible(content, () => {
				if (!this.shown) return
				const search = config.searchable && this.search()
				if (search) search.focus()
				if (which) this.highlight(which)
			})
			dispatch(this.$root, prefix + ':opened')
		},

		deactivate(notify) {
			if (!this.shown) return
			this.shown = false
			stack.close(this.layer)
			this.active = null
			this.wantsNext = false
			if (!notify) return
			restoreFocus(this.content(), [this.trigger()])
			this.closed()
			dispatch(this.$root, prefix + ':closed')
		},

		closed() {},

		refreshed() {
			this.remember()
			this.paint()
			this.mark()
			this.paintSearching()
			if (!this.shown) return
			const options = this.options()
			const current = options.find((option) => option.id === this.active)
			if (current && this.wantsNext) {
				const next = options[options.indexOf(current) + 1]
				if (next) {
					this.wantsNext = false
					this.highlight(next)
					return
				}
			}
			if (current) return
			this.wantsNext = false
			this.highlight('first')
		},

		seedLabel() {
			if (config.multiple || !this.value) return
			const label = this.$root.querySelector('[data-hui-value="' + this.$root.id + '"]')
			if (label && label.textContent.trim()) this.labels[this.value] = label.textContent.trim()
		},

		remember() {
			const size = this.loader() ? '-1' : null
			this.options(true).forEach((option, index) => {
				if (!option.id) option.id = this.$root.id + '-option-' + idSafe(option.dataset.value)
				this.labels[option.dataset.value] = textOf(option)
				if (size) {
					option.setAttribute('aria-setsize', size)
					option.setAttribute('aria-posinset', String(index + 1))
				} else if (option.getAttribute('aria-setsize') === '-1') {
					option.removeAttribute('aria-setsize')
					option.removeAttribute('aria-posinset')
				}
			})
		},

		paint() {
			this.options(true).forEach((option) => {
				const selected = this.isSelected(option.dataset.value)
				option.setAttribute('aria-selected', String(selected))
				option.toggleAttribute('data-selected', selected)
				option.toggleAttribute('data-highlighted', option.id === this.active)
				option.querySelectorAll('[data-hui-indicator]').forEach((indicator) => {
					indicator.style.display = selected ? '' : 'none'
				})
			})
		},

		highlight(which) {
			const options = this.options()
			let option = which
			if (which === 'first') option = options[0]
			if (which === 'last') option = options[options.length - 1]
			if (which === 'selected') option = options.find((item) => this.isSelected(item.dataset.value)) || options[0]
			this.active = option ? option.id : null
			if (option) option.scrollIntoView({ block: 'nearest' })
		},

		move(delta) {
			const options = this.options()
			if (!options.length) return
			const index = options.findIndex((option) => option.id === this.active)
			const loader = this.loader()
			if (delta > 0 && index === options.length - 1 && loader) {
				this.wantsNext = true
				loader.scrollIntoView({ block: 'nearest' })
				return
			}
			this.wantsNext = false
			const next = index === -1 ? (delta > 0 ? 0 : options.length - 1) : (index + delta + options.length) % options.length
			this.highlight(options[next])
		},

		current() {
			return this.options().find((option) => option.id === this.active) || null
		},

		hover(event) {
			const option = event.target.closest && event.target.closest('[role=option]')
			if (!option || option.id === this.active || !this.options().includes(option)) return
			this.active = option.id
		},

		pick(event) {
			const option = event.target.closest('[role=option]')
			if (!option || this.searching || event.defaultPrevented || !this.options().includes(option)) return
			this.choose(option)
		},

		choose(option) {
			const value = option.dataset.value
			if (config.multiple) {
				this.value = this.isSelected(value) ? this.selected().filter((item) => item !== value) : this.selected().concat(value)
			} else {
				this.value = value
			}
			this.active = option.id
			this.write()
			this.emit()
			if (!config.multiple) this.hide()
		},

		remove(value) {
			if (config.disabled) return
			this.value = config.multiple ? this.selected().filter((item) => item !== value) : ''
			this.write()
			this.emit()
		},

		emit() {
			const input = this.$root.querySelector('input[data-hui-input="' + this.$root.id + '"]')
			;(input || this.$root).dispatchEvent(new Event('change', { bubbles: true }))
			dispatch(this.$root, prefix + ':change', { value: this.value })
		},

		labelOf(value) {
			const option = this.options(true).find((item) => item.dataset.value === value)
			if (option) return textOf(option)
			return Object.hasOwn(this.labels, value) ? this.labels[value] : null
		},

		write() {
			const root = this.$root
			const values = this.selected()
			root.querySelectorAll('input[data-hui-input="' + root.id + '"]').forEach((input) => input.remove())
			if (config.name) {
				const inputs = values.map((value) => {
					const input = document.createElement('input')
					input.type = 'hidden'
					input.name = config.name
					input.value = value
					input.dataset.huiInput = root.id
					if (config.form) input.setAttribute('form', config.form)
					if (config.disabled) input.disabled = true
					return input
				})
				root.prepend(...inputs)
			}
			const require = root.querySelector('input[data-hui-require="' + root.id + '"]')
			if (require) require.value = values.length ? '1' : ''
			this.paint()
			this.written(values)
		},

		written(values) {
			const label = this.$root.querySelector('[data-hui-value="' + this.$root.id + '"]')
			if (!label) return
			const labels = values.map((value) => this.labelOf(value) ?? value)
			label.textContent = labels.length ? labels.join(', ') : config.placeholder || ''
		},

		searchFor(key) {
			this.show('selected')
			whenVisible(this.search(), () => {
				const search = this.search()
				if (!search || !this.shown) return
				search.focus()
				search.value += key
				search.dispatchEvent(new Event('input', { bubbles: true }))
			})
		},

		searchKeydown(event) {
			if (this.inputKeydown(event)) return
			if (event.key !== 'Tab') return
			event.preventDefault()
			const trigger = this.trigger()
			this.hide()
			if (event.shiftKey || !trigger) return
			const next = nextTabbable(trigger)
			if (next) next.focus()
		},

		triggerKeydown(event) {
			if (this.open) return
			const printable = event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey
			if (['ArrowDown', 'ArrowUp', 'Enter', ' '].includes(event.key)) {
				event.preventDefault()
				this.show('selected')
			} else if (printable) {
				event.preventDefault()
				this.searchFor(event.key)
			}
		},

		typeahead(key) {
			clearTimeout(this.queryTimer)
			this.query += key.toLowerCase()
			this.queryTimer = setTimeout(() => (this.query = ''), 500)

			const options = this.options()
			const current = options.findIndex((option) => option.id === this.active)
			const start = Math.max(this.query.length > 1 ? current : current + 1, 0)
			const ordered = options.slice(start).concat(options.slice(0, start))
			const match = ordered.find((option) => textOf(option).toLowerCase().startsWith(this.query))
			if (match) this.highlight(match)
		},

		input() {
			if (!this.open) this.show(null)
			this.active = null
			this.wantsNext = false
			this.updateSearching()
			this.mark()
		},

		inputKeydown(event) {
			switch (event.key) {
				case 'ArrowDown':
					event.preventDefault()
					if (event.altKey) this.show(null)
					else if (!this.open) this.show('first')
					else this.move(1)
					return true
				case 'ArrowUp':
					event.preventDefault()
					if (!this.open) this.show('last')
					else this.move(-1)
					return true
				case 'Enter':
					if (!this.open || !this.current()) return false
					event.preventDefault()
					if (!this.searching) this.choose(this.current())
					return true
			}
			return false
		},

		escape(event) {
			if (event.key !== 'Escape') return
			event.preventDefault()
			this.hide()
		},
	}
}
