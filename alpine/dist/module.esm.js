// src/checkable.js
function checkable(config = {}) {
  return {
    checked: config.checked === true,
    indeterminate: config.indeterminate === true,
    initial: null,
    owner: null,
    reset: null,
    setup() {
      this.initial = { checked: this.checked, indeterminate: this.indeterminate };
      const required = this.validation();
      if (required) required.addEventListener("focus", () => this.$root.focus());
      this.owner = this.$root.form;
      if (this.owner) {
        this.reset = () => setTimeout(() => this.restore());
        this.owner.addEventListener("reset", this.reset);
      }
      this.$watch("checked", () => this.write());
      this.$watch("indeterminate", () => this.write());
    },
    teardown() {
      if (this.owner && this.reset) this.owner.removeEventListener("reset", this.reset);
    },
    restore() {
      this.checked = this.initial.checked;
      this.indeterminate = this.initial.indeterminate;
      this.write();
    },
    ariaChecked() {
      return this.indeterminate ? "mixed" : String(this.checked);
    },
    validation() {
      const next = this.$root.nextElementSibling;
      return next && next.matches("input[data-hui-require]") ? next : null;
    },
    input() {
      return this.$root.querySelector(":scope > input[data-hui-checked]");
    },
    write() {
      if (config.name && this.checked) {
        this.$root.setAttribute("name", config.name);
        this.$root.setAttribute("value", config.value ?? "on");
      } else {
        this.$root.removeAttribute("name");
        this.$root.removeAttribute("value");
      }
      const current = this.input();
      if (config.name && this.checked && !current) {
        const input = document.createElement("input");
        input.type = "hidden";
        input.name = config.name;
        input.value = config.value ?? "on";
        input.dataset.huiChecked = "";
        if (config.form) input.setAttribute("form", config.form);
        if (config.disabled) input.disabled = true;
        this.$root.prepend(input);
      } else if (!this.checked && current) {
        current.remove();
      }
      this.validate();
    },
    validate() {
      const required = this.validation();
      if (required) required.value = this.checked ? "1" : "";
    },
    emit() {
      this.$root.dispatchEvent(new Event("change", { bubbles: true }));
    }
  };
}

// src/checkbox.js
function checkbox(Alpine) {
  Alpine.data("huiCheckbox", (config = {}) => ({
    ...checkable(config),
    init() {
      this.setup();
    },
    destroy() {
      this.teardown();
    },
    toggle() {
      if (config.disabled) return;
      this.checked = this.indeterminate ? true : !this.checked;
      this.indeterminate = false;
      this.write();
      this.emit();
    },
    keydown(event) {
      if (event.key === "Enter") event.preventDefault();
    }
  }));
}

// node_modules/@alpinejs/collapse/dist/module.esm.js
function src_default(Alpine) {
  Alpine.directive("collapse", collapse);
  collapse.inline = (el, { modifiers }) => {
    if (!modifiers.includes("min"))
      return;
    el._x_doShow = () => {
    };
    el._x_doHide = () => {
    };
  };
  function collapse(el, { modifiers }) {
    let duration = modifierValue(modifiers, "duration", 250) / 1e3;
    let floor = modifierValue(modifiers, "min", 0);
    let fullyHide = !modifiers.includes("min");
    if (!el._x_isShown)
      el.style.height = `${floor}px`;
    if (!el._x_isShown && fullyHide)
      el.hidden = true;
    if (!el._x_isShown)
      el.style.overflow = "hidden";
    let setFunction = (el2, styles) => {
      let revertFunction = Alpine.setStyles(el2, styles);
      return styles.height ? () => {
      } : revertFunction;
    };
    let transitionStyles = {
      transitionProperty: "height",
      transitionDuration: `${duration}s`,
      transitionTimingFunction: "cubic-bezier(0.4, 0.0, 0.2, 1)"
    };
    el._x_transition = {
      in(before = () => {
      }, after = () => {
      }) {
        if (fullyHide)
          el.hidden = false;
        if (fullyHide)
          el.style.display = null;
        let current = el.getBoundingClientRect().height;
        el.style.height = "auto";
        let full = el.getBoundingClientRect().height;
        if (current === full) {
          current = floor;
        }
        Alpine.transition(el, Alpine.setStyles, {
          during: transitionStyles,
          start: { height: current + "px" },
          end: { height: full + "px" }
        }, () => el._x_isShown = true, () => {
          if (Math.abs(el.getBoundingClientRect().height - full) < 1) {
            el.style.overflow = null;
          }
        });
      },
      out(before = () => {
      }, after = () => {
      }) {
        let full = el.getBoundingClientRect().height;
        Alpine.transition(el, setFunction, {
          during: transitionStyles,
          start: { height: full + "px" },
          end: { height: floor + "px" }
        }, () => el.style.overflow = "hidden", () => {
          el._x_isShown = false;
          if (el.style.height == `${floor}px` && fullyHide) {
            el.style.display = "none";
            el.hidden = true;
          }
        });
      }
    };
  }
}
function modifierValue(modifiers, key2, fallback) {
  if (modifiers.indexOf(key2) === -1)
    return fallback;
  const rawValue = modifiers[modifiers.indexOf(key2) + 1];
  if (!rawValue)
    return fallback;
  if (key2 === "duration") {
    let match = rawValue.match(/([0-9]+)ms/);
    if (match)
      return match[1];
  }
  if (key2 === "min") {
    let match = rawValue.match(/([0-9]+)px/);
    if (match)
      return match[1];
  }
  return rawValue;
}
var module_default = src_default;

// src/collapse.js
var interactive = "a[href], button, input, select, textarea, label, summary, [tabindex], [contenteditable]";
function headerClicked(event) {
  const target = event.target instanceof Element ? event.target.closest(interactive) : null;
  if (target && event.currentTarget.contains(target)) return false;
  const selection = window.getSelection && window.getSelection();
  return !(selection && selection.toString());
}
function collapsing(Alpine) {
  if (Alpine.huiCollapse) return;
  Alpine.huiCollapse = true;
  Alpine.plugin(module_default);
}

// src/layers.js
var tabbable = "a[href], area[href], button, input, select, textarea, iframe, summary, [contenteditable]:not([contenteditable=false]), [tabindex]";
function tabbables(root) {
  if (!root) return [];
  return Array.from(root.querySelectorAll(tabbable)).filter(
    (element) => element.tabIndex >= 0 && !element.matches(":disabled") && !element.closest("[inert]") && element.getClientRects().length > 0
  );
}
function idSafe(value) {
  return Array.from(String(value), (char) => /[a-zA-Z0-9_-]/.test(char) ? char : "-").join("");
}
function anchorName(id) {
  return "--" + idSafe(id);
}
function moveAnchor(from, to, name) {
  if (from === to) return;
  if (from) from.style.removeProperty("anchor-name");
  if (to) to.style.setProperty("anchor-name", name);
}
function expand(element, controls, open) {
  if (!element) return;
  element.setAttribute("aria-controls", controls);
  element.setAttribute("aria-expanded", String(open));
  element.toggleAttribute("data-open", open);
}
function addToken(element, attribute, token) {
  const tokens = (element.getAttribute(attribute) || "").split(/\s+/).filter(Boolean);
  if (!tokens.includes(token)) element.setAttribute(attribute, tokens.concat(token).join(" "));
}
function dispatch(root, name, detail = {}) {
  root.dispatchEvent(new CustomEvent(name, { bubbles: true, detail: { id: root.id, ...detail } }));
}
function commands(component, prefix, names = ["show", "hide", "toggle"]) {
  const handler = (event) => {
    const target = event.detail && (event.detail.id || event.detail);
    if (target !== component.$root.id) return;
    component[event.type.slice(prefix.length + 1)]();
  };
  names.forEach((name) => window.addEventListener(prefix + ":" + name, handler));
  return () => names.forEach((name) => window.removeEventListener(prefix + ":" + name, handler));
}
function whenVisible(element, callback, attempt = 0) {
  if (!element) return;
  if (element.getClientRects().length) {
    callback();
    return;
  }
  if (attempt < 20) setTimeout(() => whenVisible(element, callback, attempt + 1), 16);
}
function process(elements) {
  if (!window.htmx) return;
  elements.forEach((element) => element && window.htmx.process(element));
}
function restoreFocus(content, candidates) {
  const current = document.activeElement;
  if (content && current && current !== document.body && !content.contains(current)) return;
  const target = candidates.find((element) => element && element !== document.body && element.isConnected);
  if (target) target.focus();
}
function nextTabbable(from) {
  const items = tabbables(document.body);
  const index = items.indexOf(from);
  return index === -1 ? null : items[index + 1] || null;
}
function layers(Alpine) {
  if (Alpine.huiLayers) return Alpine.huiLayers;
  const stack = [];
  let inerted = [];
  const inside = (layer, target) => layer.inside().some((element) => element && element.contains(target));
  const release = () => {
    inerted.forEach((element) => element.inert = false);
    inerted = [];
  };
  const isolate = () => {
    release();
    let index = stack.length - 1;
    while (index >= 0 && !stack[index].modal) index--;
    if (index === -1) return;
    const keep = /* @__PURE__ */ new Set();
    stack.slice(index).forEach((layer) => layer.parts().forEach((element) => element && keep.add(element)));
    document.querySelectorAll("[data-hui-toaster], [data-hui-announcer]").forEach((element) => keep.add(element));
    const path = /* @__PURE__ */ new Set();
    keep.forEach((element) => {
      let node = element;
      while (node && node !== document.body && !path.has(node)) {
        path.add(node);
        node = node.parentElement;
      }
    });
    path.forEach((node) => {
      if (!node.parentElement) return;
      for (const sibling of node.parentElement.children) {
        if (path.has(sibling) || sibling.inert || sibling.matches("script, style, template")) continue;
        sibling.inert = true;
        inerted.push(sibling);
      }
    });
  };
  document.addEventListener("keydown", (event) => {
    const top = stack[stack.length - 1];
    if (top && !event.defaultPrevented) top.keydown(event);
  });
  document.addEventListener(
    "click",
    (event) => {
      const target = event.target;
      if (!stack.length || !(target instanceof Element) || !target.isConnected) return;
      if (target.closest("[data-hui-toaster]")) return;
      for (let index = stack.length - 1; index >= 0; index--) {
        const layer = stack[index];
        if (inside(layer, target)) break;
        if (layer.dismiss(event) === false || layer.modal) break;
      }
    },
    true
  );
  Alpine.huiLayers = {
    open(layer) {
      if (stack.includes(layer)) return;
      stack.push(layer);
      isolate();
    },
    close(layer) {
      const index = stack.indexOf(layer);
      if (index === -1) return;
      stack.splice(index, 1);
      isolate();
    },
    top(layer) {
      return stack[stack.length - 1] === layer;
    }
  };
  return Alpine.huiLayers;
}

// src/collapsible.js
function collapsible(Alpine) {
  collapsing(Alpine);
  Alpine.data("huiCollapsible", (config = {}) => ({
    open: config.open === true,
    cleanup: null,
    init() {
      this.cleanup = commands(this, "hui-collapsible");
      this.$watch("open", (open) => dispatch(this.$root, open ? "hui-collapsible:opened" : "hui-collapsible:closed"));
    },
    destroy() {
      this.cleanup();
    },
    show() {
      this.open = true;
    },
    hide() {
      this.open = false;
    },
    toggle() {
      this.open = !this.open;
    },
    collapsibleShow() {
      this.show();
    },
    collapsibleHide() {
      this.hide();
    },
    collapsibleToggle() {
      this.toggle();
    },
    collapsibleHeaderClick(event) {
      const trigger = document.getElementById(this.$root.id + "-trigger");
      if (trigger && trigger.disabled || !headerClicked(event)) return;
      this.toggle();
    }
  }));
}

// src/pending.js
var key = /* @__PURE__ */ Symbol.for("headlessui.pending");
function mark(element) {
  const state = element[key] || (element[key] = { count: 0, kept: element.hasAttribute("data-pending") });
  state.count++;
  element.setAttribute("data-pending", "");
  element.setAttribute("aria-busy", "true");
}
function unmark(element) {
  const state = element[key];
  if (!state || --state.count > 0) return;
  delete element[key];
  if (state.kept) return;
  element.removeAttribute("data-pending");
  element.removeAttribute("aria-busy");
}
function leaving(ctx) {
  return !!ctx.hx && (!!ctx.hx.redirect || ctx.hx.refresh === "true");
}
function markPending(elements, ctx) {
  const marked = Array.from(new Set(elements)).filter(
    (element) => element instanceof Element && element !== document.body && element !== document.documentElement
  );
  if (!marked.length) return;
  marked.forEach(mark);
  const done = (event) => {
    if (!event.detail || event.detail.ctx !== ctx) return;
    document.removeEventListener("htmx:finally:request", done);
    if (!leaving(ctx)) marked.forEach(unmark);
  };
  document.addEventListener("htmx:finally:request", done);
}
function pending() {
  if (document[key]) return;
  document[key] = true;
  document.addEventListener("htmx:before:request", (event) => {
    const ctx = event.detail && event.detail.ctx;
    if (!ctx || !(event.target instanceof Element)) return;
    markPending([ctx.sourceElement || event.target, ctx.request && ctx.request.submitter, ctx.target], ctx);
  });
}

// src/options.js
function textOf(option) {
  if (option.dataset.label) return option.dataset.label;
  const clone = option.cloneNode(true);
  clone.querySelectorAll("[aria-hidden=true]").forEach((element) => element.remove());
  return clone.textContent.trim();
}
function fold(text) {
  return text.normalize("NFD").replace(/\p{M}/gu, "").toLowerCase();
}
function folded(text) {
  let value = "";
  const starts = [];
  const ends = [];
  for (let at = 0; at < text.length; ) {
    const char = String.fromCodePoint(text.codePointAt(at));
    const plain = fold(char);
    for (let i = 0; i < plain.length; i++) {
      starts.push(at);
      ends.push(at + char.length);
    }
    value += plain;
    at += char.length;
  }
  return { value, starts, ends };
}
function highlight(option, needle) {
  const hidden = (node) => {
    const element = node.parentElement && node.parentElement.closest("[aria-hidden=true]");
    return !!element && option.contains(element);
  };
  const walker = document.createTreeWalker(option, NodeFilter.SHOW_TEXT, {
    acceptNode: (node) => hidden(node) ? NodeFilter.FILTER_REJECT : NodeFilter.FILTER_ACCEPT
  });
  const nodes = [];
  while (walker.nextNode()) nodes.push(walker.currentNode);
  nodes.forEach((node) => {
    const { value, starts, ends } = folded(node.data);
    const ranges = [];
    for (let at = value.indexOf(needle); at !== -1; at = value.indexOf(needle, at + needle.length)) {
      ranges.push([starts[at], ends[at + needle.length - 1]]);
    }
    ranges.reverse().forEach(([start, end]) => {
      const match = node.splitText(start);
      match.splitText(end - start);
      const mark2 = document.createElement("mark");
      mark2.dataset.huiMatch = "";
      match.replaceWith(mark2);
      mark2.append(match);
    });
  });
}
function options(Alpine, kind, config) {
  const stack = layers(Alpine);
  const prefix = "hui-" + kind;
  return {
    [kind + "Show"]() {
      this.show("selected");
      this.shown && this.focusControl();
    },
    [kind + "Hide"]() {
      this.hide();
    },
    [kind + "Toggle"]() {
      this.toggle();
    },
    open: config.open === true,
    value: config.multiple ? (config.value || []).slice() : config.value ?? "",
    active: null,
    shown: false,
    pending: null,
    layer: null,
    cleanup: null,
    observer: null,
    requests: null,
    labels: {},
    wantsNext: false,
    query: "",
    queryTimer: null,
    searching: false,
    searched: null,
    setup(inside) {
      this.layer = {
        parts: () => [this.content()],
        inside,
        dismiss: () => this.hide(),
        keydown: (event) => this.escape(event)
      };
      this.cleanup = commands(this, prefix);
      this.$watch("open", () => this.sync());
      this.$watch("value", () => this.write());
      this.$watch("active", () => this.paint());
      this.$watch("searching", () => this.paintSearching());
      this.watchRequests();
      this.$nextTick(() => {
        const content = this.content();
        if (content) {
          this.observer = new MutationObserver(() => this.refreshed());
          this.observer.observe(content, { childList: true, subtree: true });
        }
        this.searched = config.searchable ? this.term() : null;
        this.seedLabel();
        this.remember();
        this.paint();
        this.mark();
        this.sync();
      });
    },
    teardown() {
      this.cleanup();
      this.requests();
      clearTimeout(this.queryTimer);
      if (this.observer) this.observer.disconnect();
      this.deactivate();
    },
    watchRequests() {
      const mine = (event) => {
        const target = event.target;
        if (!(target instanceof Element)) return false;
        const content = this.content();
        return content && content.contains(target) || target === this.control();
      };
      const start = (event) => {
        if (!mine(event)) return;
        markPending([this.list()], event.detail.ctx);
        if (event.target === this.search()) this.trackSearch(event.detail.ctx);
      };
      const configure = (event) => {
        const search = this.search();
        const name = search && search.dataset.huiQuery;
        if (!name || event.target !== search) return;
        event.detail.ctx.request.body.set(name, this.term());
      };
      document.addEventListener("htmx:config:request", configure);
      document.addEventListener("htmx:before:request", start);
      this.requests = () => {
        document.removeEventListener("htmx:config:request", configure);
        document.removeEventListener("htmx:before:request", start);
      };
    },
    trackSearch(ctx) {
      const query = this.term();
      const done = (event) => {
        if (!event.detail || event.detail.ctx !== ctx) return;
        document.removeEventListener("htmx:finally:request", done);
        this.searched = query;
        this.updateSearching();
      };
      document.addEventListener("htmx:finally:request", done);
    },
    updateSearching() {
      const search = this.search();
      const remote = !!search && Array.from(search.attributes).some((attribute) => /^(data-)?hx-(get|post|put|patch|delete)$/.test(attribute.name));
      this.searching = remote && this.term() !== this.searched;
    },
    term() {
      const search = config.searchable && this.search();
      return search ? search.value : "";
    },
    mark() {
      const list = this.list();
      if (!list) return;
      list.querySelectorAll("mark[data-hui-match]").forEach((mark2) => {
        const parent = mark2.parentNode;
        mark2.replaceWith(...mark2.childNodes);
        parent.normalize();
      });
      const needle = fold(this.term().trim());
      if (needle) this.options(true).forEach((option) => highlight(option, needle));
      if (this.observer) this.observer.takeRecords();
    },
    paintSearching() {
      ;
      [this.search(), this.list()].forEach((element) => element && element.toggleAttribute("data-searching", this.searching));
    },
    listId() {
      return this.$root.id + "-list";
    },
    control() {
      if (config.searchable) return this.search();
      return document.querySelector('[role=combobox][aria-controls="' + this.listId() + '"]');
    },
    trigger() {
      return this.$root.querySelector('[aria-controls="' + this.listId() + '"]');
    },
    search() {
      if (kind === "combobox") return this.$root.querySelector("input[role=combobox]");
      return document.getElementById(this.$root.id + "-search");
    },
    content() {
      return document.getElementById(this.$root.id + "-content");
    },
    list() {
      return document.getElementById(this.listId()) || this.content();
    },
    loader() {
      const list = this.list();
      return list ? list.querySelector("[data-hui-loader]") : null;
    },
    focusControl() {
      const control = config.searchable ? this.trigger() : this.control();
      if (control) control.focus();
    },
    options(all = false) {
      const list = this.list();
      if (!list) return [];
      const options2 = Array.from(list.querySelectorAll("[role=option]"));
      return all ? options2 : options2.filter((option) => option.getAttribute("aria-disabled") !== "true");
    },
    selected() {
      if (config.multiple) return Array.isArray(this.value) ? this.value : [];
      return this.value === "" || this.value == null ? [] : [this.value];
    },
    isSelected(value) {
      return this.selected().includes(value);
    },
    show(which = null) {
      if (config.disabled) return;
      this.pending = which;
      this.open = true;
      this.sync();
    },
    hide() {
      this.open = false;
      this.sync();
    },
    toggle() {
      this.open ? this.hide() : this.show("selected");
    },
    sync() {
      if (this.open === this.shown) return;
      this.open ? this.activate() : this.deactivate(true);
    },
    activate() {
      this.shown = true;
      stack.open(this.layer);
      const content = this.content();
      process([content]);
      const which = this.pending;
      this.pending = null;
      whenVisible(content, () => {
        if (!this.shown) return;
        const search = config.searchable && this.search();
        if (search) search.focus();
        if (which) this.highlight(which);
      });
      dispatch(this.$root, prefix + ":opened");
    },
    deactivate(notify) {
      if (!this.shown) return;
      this.shown = false;
      stack.close(this.layer);
      this.active = null;
      this.wantsNext = false;
      if (!notify) return;
      restoreFocus(this.content(), [this.trigger()]);
      this.closed();
      dispatch(this.$root, prefix + ":closed");
    },
    closed() {
    },
    refreshed() {
      this.remember();
      this.paint();
      this.mark();
      this.paintSearching();
      if (!this.shown) return;
      const options2 = this.options();
      const current = options2.find((option) => option.id === this.active);
      if (current && this.wantsNext) {
        const next = options2[options2.indexOf(current) + 1];
        if (next) {
          this.wantsNext = false;
          this.highlight(next);
          return;
        }
      }
      if (current) return;
      this.wantsNext = false;
      this.highlight("first");
    },
    seedLabel() {
      if (config.multiple || !this.value) return;
      const label = this.$root.querySelector('[data-hui-value="' + this.$root.id + '"]');
      if (label && label.textContent.trim()) this.labels[this.value] = label.textContent.trim();
    },
    remember() {
      const size = this.loader() ? "-1" : null;
      this.options(true).forEach((option, index) => {
        if (!option.id) option.id = this.$root.id + "-option-" + idSafe(option.dataset.value);
        this.labels[option.dataset.value] = textOf(option);
        if (size) {
          option.setAttribute("aria-setsize", size);
          option.setAttribute("aria-posinset", String(index + 1));
        } else if (option.getAttribute("aria-setsize") === "-1") {
          option.removeAttribute("aria-setsize");
          option.removeAttribute("aria-posinset");
        }
      });
    },
    paint() {
      this.options(true).forEach((option) => {
        const selected = this.isSelected(option.dataset.value);
        option.setAttribute("aria-selected", String(selected));
        option.toggleAttribute("data-selected", selected);
        option.toggleAttribute("data-highlighted", option.id === this.active);
        option.querySelectorAll("[data-hui-indicator]").forEach((indicator) => {
          indicator.style.display = selected ? "" : "none";
        });
      });
    },
    highlight(which) {
      const options2 = this.options();
      let option = which;
      if (which === "first") option = options2[0];
      if (which === "last") option = options2[options2.length - 1];
      if (which === "selected") option = options2.find((item) => this.isSelected(item.dataset.value)) || options2[0];
      this.active = option ? option.id : null;
      if (option) option.scrollIntoView({ block: "nearest" });
    },
    move(delta) {
      const options2 = this.options();
      if (!options2.length) return;
      const index = options2.findIndex((option) => option.id === this.active);
      const loader = this.loader();
      if (delta > 0 && index === options2.length - 1 && loader) {
        this.wantsNext = true;
        loader.scrollIntoView({ block: "nearest" });
        return;
      }
      this.wantsNext = false;
      const next = index === -1 ? delta > 0 ? 0 : options2.length - 1 : (index + delta + options2.length) % options2.length;
      this.highlight(options2[next]);
    },
    current() {
      return this.options().find((option) => option.id === this.active) || null;
    },
    hover(event) {
      const option = event.target.closest && event.target.closest("[role=option]");
      if (!option || option.id === this.active || !this.options().includes(option)) return;
      this.active = option.id;
    },
    pick(event) {
      const option = event.target.closest("[role=option]");
      if (!option || this.searching || event.defaultPrevented || !this.options().includes(option)) return;
      this.choose(option);
    },
    choose(option) {
      const value = option.dataset.value;
      if (config.multiple) {
        this.value = this.isSelected(value) ? this.selected().filter((item) => item !== value) : this.selected().concat(value);
      } else {
        this.value = value;
      }
      this.active = option.id;
      this.write();
      this.emit();
      if (!config.multiple) this.hide();
    },
    remove(value) {
      if (config.disabled) return;
      this.value = config.multiple ? this.selected().filter((item) => item !== value) : "";
      this.write();
      this.emit();
    },
    emit() {
      const input = this.$root.querySelector('input[data-hui-input="' + this.$root.id + '"]');
      (input || this.$root).dispatchEvent(new Event("change", { bubbles: true }));
      dispatch(this.$root, prefix + ":change", { value: this.value });
    },
    labelOf(value) {
      const option = this.options(true).find((item) => item.dataset.value === value);
      if (option) return textOf(option);
      return Object.hasOwn(this.labels, value) ? this.labels[value] : null;
    },
    write() {
      const root = this.$root;
      const values = this.selected();
      root.querySelectorAll('input[data-hui-input="' + root.id + '"]').forEach((input) => input.remove());
      if (config.name) {
        const inputs = values.map((value) => {
          const input = document.createElement("input");
          input.type = "hidden";
          input.name = config.name;
          input.value = value;
          input.dataset.huiInput = root.id;
          if (config.form) input.setAttribute("form", config.form);
          if (config.disabled) input.disabled = true;
          return input;
        });
        root.prepend(...inputs);
      }
      const require2 = root.querySelector('input[data-hui-require="' + root.id + '"]');
      if (require2) require2.value = values.length ? "1" : "";
      this.paint();
      this.written(values);
    },
    written(values) {
      const label = this.$root.querySelector('[data-hui-value="' + this.$root.id + '"]');
      if (!label) return;
      const labels = values.map((value) => this.labelOf(value) ?? value);
      label.textContent = labels.length ? labels.join(", ") : config.placeholder || "";
    },
    searchFor(key2) {
      this.show("selected");
      whenVisible(this.search(), () => {
        const search = this.search();
        if (!search || !this.shown) return;
        search.focus();
        search.value += key2;
        search.dispatchEvent(new Event("input", { bubbles: true }));
      });
    },
    searchKeydown(event) {
      if (this.inputKeydown(event)) return;
      if (event.key !== "Tab") return;
      event.preventDefault();
      const trigger = this.trigger();
      this.hide();
      if (event.shiftKey || !trigger) return;
      const next = nextTabbable(trigger);
      if (next) next.focus();
    },
    triggerKeydown(event) {
      if (this.open) return;
      const printable = event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey;
      if (["ArrowDown", "ArrowUp", "Enter", " "].includes(event.key)) {
        event.preventDefault();
        this.show("selected");
      } else if (printable) {
        event.preventDefault();
        this.searchFor(event.key);
      }
    },
    typeahead(key2) {
      clearTimeout(this.queryTimer);
      this.query += key2.toLowerCase();
      this.queryTimer = setTimeout(() => this.query = "", 500);
      const options2 = this.options();
      const current = options2.findIndex((option) => option.id === this.active);
      const start = Math.max(this.query.length > 1 ? current : current + 1, 0);
      const ordered = options2.slice(start).concat(options2.slice(0, start));
      const match = ordered.find((option) => textOf(option).toLowerCase().startsWith(this.query));
      if (match) this.highlight(match);
    },
    input() {
      if (!this.open) this.show(null);
      this.active = null;
      this.wantsNext = false;
      this.updateSearching();
      this.mark();
    },
    inputKeydown(event) {
      switch (event.key) {
        case "ArrowDown":
          event.preventDefault();
          if (event.altKey) this.show(null);
          else if (!this.open) this.show("first");
          else this.move(1);
          return true;
        case "ArrowUp":
          event.preventDefault();
          if (!this.open) this.show("last");
          else this.move(-1);
          return true;
        case "Enter":
          if (!this.open || !this.current()) return false;
          event.preventDefault();
          if (!this.searching) this.choose(this.current());
          return true;
      }
      return false;
    },
    escape(event) {
      if (event.key !== "Escape") return;
      event.preventDefault();
      this.hide();
    }
  };
}

// src/combobox.js
function combobox(Alpine) {
  Alpine.data("huiCombobox", (config = {}) => ({
    ...options(Alpine, "combobox", config),
    typed: false,
    init() {
      this.setup(() => [this.$root, this.content()]);
    },
    destroy() {
      this.teardown();
    },
    comboboxToggle() {
      this.toggle();
      this.focusControl();
    },
    term() {
      const search = this.search();
      return this.typed && search ? search.value : "";
    },
    focusControl() {
      const search = this.search();
      if (search && document.activeElement !== search) search.focus();
    },
    seedLabel() {
      const search = this.search();
      if (config.multiple || !this.value || !search || !search.value) return;
      this.labels[this.value] = search.value;
    },
    input() {
      this.typed = true;
      if (!this.open) this.show(null);
      this.active = null;
      this.wantsNext = false;
      this.ask();
    },
    ask() {
      const search = this.search();
      if (search && this.term() !== this.searched) dispatch(search, "hui-combobox:query", { query: this.term() });
      this.updateSearching();
      this.mark();
    },
    inputClick() {
      if (this.open) return;
      this.show("selected");
      const search = this.search();
      if (search && !this.typed) search.select();
    },
    inputBlur(event) {
      const content = this.content();
      if (!this.open || content && content.contains(event.relatedTarget)) return;
      this.hide();
    },
    keydown(event) {
      if (!this.open && event.key === "ArrowDown" && !event.altKey) {
        event.preventDefault();
        this.show("selected");
        return;
      }
      if (this.inputKeydown(event)) return;
      if (event.key === "Tab" && this.open) this.hide();
    },
    written(values) {
      const search = this.search();
      if (!search) return;
      this.typed = false;
      search.value = config.multiple || !values.length ? "" : this.labelOf(values[0]) ?? values[0];
      this.ask();
    },
    closed() {
      const search = this.search();
      if (search && this.typed && !config.multiple && search.value.trim() === "" && this.selected().length) {
        this.remove(this.value);
        return;
      }
      this.written(this.selected());
    }
  }));
}

// src/dropdown.js
var menuItems = "[role=menuitem], [role=menuitemcheckbox], [role=menuitemradio]";
function dropdown(Alpine) {
  const stack = layers(Alpine);
  Alpine.data("huiDropdownCheckbox", (config = {}) => ({
    ...checkable(config),
    init() {
      this.setup();
    },
    destroy() {
      this.teardown();
    },
    dropdownCheck() {
      if (config.disabled) return;
      this.checked = !this.checked;
      this.write();
      this.emit();
    }
  }));
  Alpine.data("huiDropdownRadio", (config = {}) => ({
    ...checkable(config),
    init() {
      this.setup();
    },
    destroy() {
      this.teardown();
    },
    dropdownPick() {
      if (config.disabled || this.checked) return;
      this.checked = true;
      this.write();
      const group = this.$root.closest("[role=group], [role=menu]");
      if (group) {
        group.querySelectorAll("[role=menuitemradio]").forEach((item) => {
          if (item === this.$root) return;
          const data = Alpine.$data(item);
          if (!data.checked) return;
          data.checked = false;
          data.write();
        });
      }
      this.emit();
    }
  }));
  Alpine.data("huiDropdown", (config = {}) => ({
    open: config.open === true,
    active: false,
    focusOnOpen: "content",
    query: "",
    queryTimer: null,
    anchor: null,
    layer: null,
    cleanup: null,
    init() {
      this.anchor = this.trigger();
      this.layer = {
        parts: () => [this.content()],
        inside: () => [this.anchor, this.content()],
        dismiss: () => this.hide(),
        keydown: (event) => this.escape(event)
      };
      this.cleanup = commands(this, "hui-dropdown");
      this.$watch("open", () => this.sync());
      this.$nextTick(() => this.sync());
    },
    destroy() {
      this.cleanup();
      clearTimeout(this.queryTimer);
      this.deactivate(false);
    },
    trigger() {
      return document.getElementById(this.$root.id + "-trigger");
    },
    content() {
      return document.getElementById(this.$root.id + "-content");
    },
    items() {
      const content = this.content();
      if (!content) return [];
      return Array.from(content.querySelectorAll(menuItems)).filter(
        (item) => !item.matches(":disabled, [aria-disabled=true]")
      );
    },
    show(focus) {
      if (focus) this.focusOnOpen = focus;
      this.open = true;
      this.sync();
    },
    hide() {
      this.open = false;
      this.sync();
    },
    toggle() {
      this.open ? this.hide() : this.show();
    },
    attach(element) {
      const content = this.content();
      if (!element || element === this.anchor || content && content.contains(element)) return;
      moveAnchor(this.anchor, element, anchorName(this.$root.id));
      if (this.anchor) expand(this.anchor, this.$root.id + "-content", false);
      this.anchor = element;
      expand(element, this.$root.id + "-content", this.open);
    },
    dropdownShow(element, event) {
      this.attach(element);
      if (event && event.type === "click") this.focusOnOpen = event.detail === 0 ? "first" : "content";
      this.show();
    },
    dropdownHide() {
      this.hide();
    },
    dropdownToggle(element, event) {
      if (this.open && (!element || element === this.anchor)) this.hide();
      else this.dropdownShow(element, event);
    },
    sync() {
      if (this.open === this.active) return;
      this.open ? this.activate() : this.deactivate(true);
    },
    activate() {
      this.active = true;
      expand(this.anchor, this.$root.id + "-content", true);
      const focus = this.focusOnOpen;
      this.focusOnOpen = "content";
      stack.open(this.layer);
      process([this.content()]);
      whenVisible(this.content(), () => this.active && this.focus(focus));
      dispatch(this.$root, "hui-dropdown:opened");
    },
    deactivate(notify) {
      if (!this.active) return;
      this.active = false;
      stack.close(this.layer);
      expand(this.anchor, this.$root.id + "-content", false);
      if (!notify) return;
      restoreFocus(this.content(), [this.anchor]);
      dispatch(this.$root, "hui-dropdown:closed");
    },
    focus(which) {
      const items = this.items();
      const item = which === "first" ? items[0] : which === "last" ? items[items.length - 1] : null;
      const target = item || this.content();
      if (target) target.focus({ preventScroll: true });
    },
    move(delta) {
      const items = this.items();
      if (!items.length) return;
      const index = items.indexOf(document.activeElement);
      const next = index === -1 ? delta > 0 ? 0 : items.length - 1 : (index + delta + items.length) % items.length;
      items[next].focus();
    },
    typeahead(key2) {
      clearTimeout(this.queryTimer);
      this.query += key2.toLowerCase();
      this.queryTimer = setTimeout(() => this.query = "", 500);
      const items = this.items();
      const current = items.indexOf(document.activeElement);
      const start = Math.max(this.query.length > 1 ? current : current + 1, 0);
      const ordered = items.slice(start).concat(items.slice(0, start));
      const match = ordered.find((item) => item.textContent.trim().toLowerCase().startsWith(this.query));
      if (match) match.focus();
    },
    dropdownKeydown(event, element) {
      if (event.key !== "ArrowDown" && event.key !== "ArrowUp") return;
      event.preventDefault();
      const focus = event.key === "ArrowDown" ? "first" : "last";
      this.attach(element);
      this.open ? this.focus(focus) : this.show(focus);
    },
    escape(event) {
      if (event.key !== "Escape") return;
      event.preventDefault();
      this.hide();
    },
    keydown(event) {
      switch (event.key) {
        case "ArrowDown":
          event.preventDefault();
          this.move(1);
          break;
        case "ArrowUp":
          event.preventDefault();
          this.move(-1);
          break;
        case "Home":
          event.preventDefault();
          this.focus("first");
          break;
        case "End":
          event.preventDefault();
          this.focus("last");
          break;
        case "Tab":
          this.hide();
          break;
        case " ":
          if (this.query) {
            event.preventDefault();
            this.typeahead(event.key);
          } else if (event.target.matches("a[role=menuitem]")) {
            event.preventDefault();
            event.target.click();
          }
          break;
        default:
          if (event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
            this.typeahead(event.key);
          }
      }
    },
    select(event) {
      const item = event.target.closest("[role=menuitem], [role=menuitemradio]");
      if (!item || event.defaultPrevented || !this.items().includes(item)) return;
      this.hide();
    },
    hover(event) {
      const item = event.target.closest && event.target.closest(menuItems);
      if (!item || item === document.activeElement || !this.items().includes(item)) return;
      item.focus({ preventScroll: true });
    }
  }));
}

// src/fileupload.js
function accepts(file, accept) {
  if (!accept) return true;
  const name = file.name.toLowerCase();
  const type = (file.type || "").toLowerCase();
  return accept.split(",").map((rule) => rule.trim().toLowerCase()).filter(Boolean).some((rule) => {
    if (rule.startsWith(".")) return name.endsWith(rule);
    if (rule.endsWith("/*")) return type.startsWith(rule.slice(0, -1));
    return type === rule;
  });
}
function fileUpload(Alpine) {
  Alpine.data("huiFileUpload", () => ({
    dragging: false,
    depth: 0,
    files: [],
    input() {
      return this.$root.querySelector('input[type=file][data-hui-file="' + this.$root.id + '"]');
    },
    changed(event) {
      this.files = Array.from(event.target.files || []).map((file) => file.name);
    },
    browse(event) {
      const input = this.input();
      if (!input || input.disabled || event.target === input || event.target.closest("a, button, input, label")) return;
      input.click();
    },
    enter(event) {
      if (!this.input() || this.input().disabled) return;
      this.depth++;
      this.dragging = true;
    },
    over(event) {
      if (event.dataTransfer) event.dataTransfer.dropEffect = this.input() && !this.input().disabled ? "copy" : "none";
    },
    leave() {
      this.depth = Math.max(this.depth - 1, 0);
      if (this.depth === 0) this.dragging = false;
    },
    drop(event) {
      this.depth = 0;
      this.dragging = false;
      const input = this.input();
      if (!input || input.disabled || !event.dataTransfer) return;
      let files = Array.from(event.dataTransfer.files).filter((file) => accepts(file, input.accept));
      if (!input.multiple) files = files.slice(0, 1);
      if (!files.length) return;
      const transfer = new DataTransfer();
      files.forEach((file) => transfer.items.add(file));
      input.files = transfer.files;
      input.dispatchEvent(new Event("input", { bubbles: true }));
      input.dispatchEvent(new Event("change", { bubbles: true }));
    }
  }));
}

// src/listbox.js
function listbox(Alpine) {
  Alpine.data("huiListbox", (config = {}) => ({
    ...options(Alpine, "listbox", config),
    init() {
      this.setup(() => [this.trigger(), this.content()]);
    },
    destroy() {
      this.teardown();
    },
    keydown(event) {
      if (config.searchable) return this.triggerKeydown(event);
      const printable = event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey;
      if (!this.open) {
        if (["ArrowDown", "ArrowUp", "Enter", " "].includes(event.key)) {
          event.preventDefault();
          this.show("selected");
        } else if (event.key === "Home" || event.key === "End") {
          event.preventDefault();
          this.show(event.key === "Home" ? "first" : "last");
        } else if (printable) {
          this.show(null);
          this.$nextTick(() => this.typeahead(event.key));
        }
        return;
      }
      switch (event.key) {
        case "ArrowDown":
          event.preventDefault();
          event.altKey ? this.hide() : this.move(1);
          break;
        case "ArrowUp":
          event.preventDefault();
          this.move(-1);
          break;
        case "Home":
          event.preventDefault();
          this.highlight("first");
          break;
        case "End":
          event.preventDefault();
          this.highlight("last");
          break;
        case "Enter":
          event.preventDefault();
          if (this.current()) this.choose(this.current());
          break;
        case " ":
          event.preventDefault();
          if (this.query) this.typeahead(event.key);
          else if (this.current()) this.choose(this.current());
          break;
        case "Tab":
          this.hide();
          break;
        default:
          if (printable) this.typeahead(event.key);
      }
    }
  }));
}

// src/modal.js
function modal(Alpine) {
  const stack = layers(Alpine);
  let locks = 0;
  let overflow = "";
  const lock = () => {
    if (locks++ === 0) {
      overflow = document.documentElement.style.overflow;
      document.documentElement.style.overflow = "hidden";
    }
  };
  const unlock = () => {
    if (--locks === 0) document.documentElement.style.overflow = overflow;
  };
  Alpine.data("huiModal", (config = {}) => ({
    open: config.open === true,
    active: false,
    returnTo: null,
    layer: null,
    cleanup: null,
    init() {
      const stored = this.restore();
      if (stored !== null) this.open = stored;
      this.layer = {
        modal: true,
        parts: () => [this.overlay(), this.content()],
        inside: () => [this.content()],
        dismiss: () => this.dismiss(),
        keydown: (event) => this.keydown(event)
      };
      this.cleanup = commands(this, "hui-modal");
      this.$watch("open", () => this.sync());
      this.$nextTick(() => this.sync());
    },
    destroy() {
      this.cleanup();
      this.deactivate(false);
    },
    show() {
      this.open = true;
      this.sync();
    },
    hide() {
      this.open = false;
      this.sync();
    },
    toggle() {
      this.open ? this.hide() : this.show();
    },
    modalShow() {
      this.show();
    },
    modalHide() {
      this.hide();
    },
    modalToggle() {
      this.toggle();
    },
    dismiss() {
      if (config.static || config.alert) return false;
      this.hide();
    },
    sync() {
      if (this.open === this.active) return;
      this.open ? this.activate() : this.deactivate(true);
    },
    activate() {
      this.active = true;
      this.returnTo = document.activeElement;
      if (!config.allowScroll) lock();
      stack.open(this.layer);
      process(this.layer.parts());
      whenVisible(this.content(), () => this.focusIn());
      this.save(true);
      dispatch(this.$root, "hui-modal:opened");
    },
    deactivate(notify) {
      if (!this.active) return;
      this.active = false;
      stack.close(this.layer);
      if (!config.allowScroll) unlock();
      if (!notify) return;
      restoreFocus(this.content(), [this.returnTo, this.trigger()]);
      this.save(false);
      dispatch(this.$root, "hui-modal:closed");
    },
    content() {
      return document.getElementById(this.$root.id + "-content");
    },
    overlay() {
      return document.getElementById(this.$root.id + "-overlay");
    },
    trigger() {
      return this.$root.querySelector('[aria-controls="' + this.$root.id + '-content"]');
    },
    focusIn() {
      const content = this.content();
      if (!this.active || content.contains(document.activeElement)) return;
      const target = content.querySelector("[autofocus]") || tabbables(content)[0] || content;
      target.focus();
    },
    keydown(event) {
      const content = this.content();
      if (event.key === "Escape") {
        if (config.static) return;
        event.preventDefault();
        this.hide();
        return;
      }
      if (event.key !== "Tab" || !content) return;
      const items = tabbables(content);
      const current = document.activeElement;
      if (!items.length) {
        event.preventDefault();
        content.focus();
        return;
      }
      const first = items[0];
      const last = items[items.length - 1];
      if (!content.contains(current)) {
        event.preventDefault();
        (event.shiftKey ? last : first).focus();
      } else if (event.shiftKey && (current === first || current === content)) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && current === last) {
        event.preventDefault();
        first.focus();
      }
    },
    restore() {
      if (!config.persist) return null;
      try {
        const value = sessionStorage.getItem("hui-modal:" + config.persist);
        return value === null ? null : value === "1";
      } catch (error) {
        return null;
      }
    },
    save(open) {
      if (!config.persist) return;
      try {
        sessionStorage.setItem("hui-modal:" + config.persist, open ? "1" : "0");
      } catch (error) {
      }
    }
  }));
}

// src/radio.js
function radio(Alpine) {
  Alpine.data("huiRadio", (config = {}) => ({
    ...checkable(config),
    listener: null,
    init() {
      this.setup();
      this.listener = (event) => this.selected(event.detail);
      window.addEventListener("hui-radio:select", this.listener);
      this.$nextTick(() => this.settle());
    },
    destroy() {
      this.teardown();
      window.removeEventListener("hui-radio:select", this.listener);
    },
    restore() {
      this.checked = this.initial.checked;
      this.write();
      this.$nextTick(() => this.settle());
    },
    settle() {
      this.fill(this.group().some((element) => element.getAttribute("aria-checked") === "true"));
      this.rove();
    },
    group() {
      if (!config.name) return [this.$root];
      const form = this.$root.form;
      return Array.from(document.querySelectorAll("[role=radio][data-hui-radio]")).filter(
        (element) => element.dataset.huiRadio === config.name && element.form === form
      );
    },
    enabled() {
      return this.group().filter((element) => !element.disabled);
    },
    select() {
      if (config.disabled || this.checked) return;
      this.checked = true;
      this.write();
      this.fill(true);
      this.rove(this.$root);
      window.dispatchEvent(new CustomEvent("hui-radio:select", { detail: { name: config.name, root: this.$root } }));
      this.emit();
    },
    selected(detail) {
      if (!config.name || detail.name !== config.name || detail.root === this.$root) return;
      if (detail.root.form !== this.$root.form) return;
      if (this.checked) {
        this.checked = false;
        this.write();
      }
      this.fill(true);
      this.rove(detail.root);
    },
    validate() {
    },
    fill(filled) {
      const required = this.validation();
      if (required) required.value = filled ? "1" : "";
    },
    rove(checked = null) {
      if (!config.name) return;
      const group = this.enabled();
      const current = checked || group.find((element) => element.getAttribute("aria-checked") === "true");
      this.$root.tabIndex = this.$root === (current || group[0]) ? 0 : -1;
    },
    keydown(event) {
      const keys = { ArrowDown: 1, ArrowRight: 1, ArrowUp: -1, ArrowLeft: -1 };
      if (!(event.key in keys)) return;
      event.preventDefault();
      const group = this.enabled();
      const index = group.indexOf(this.$root);
      if (index === -1 || group.length < 2) return;
      const next = group[(index + keys[event.key] + group.length) % group.length];
      next.focus();
      next.click();
    }
  }));
}

// src/range.js
function range(Alpine) {
  Alpine.data("huiRange", () => ({
    values: [],
    form: null,
    reset: null,
    init() {
      const thumb = this.thumbs()[0];
      this.form = thumb ? thumb.form : null;
      this.reset = () => setTimeout(() => this.sync());
      if (this.form) this.form.addEventListener("reset", this.reset);
      this.sync();
    },
    destroy() {
      if (this.form) this.form.removeEventListener("reset", this.reset);
    },
    thumbs() {
      return Array.from(this.$root.querySelectorAll('input[type="range"]'));
    },
    update(thumb) {
      const thumbs = this.thumbs();
      const index = thumbs.indexOf(thumb);
      if (index === -1) return;
      const value = Number(thumb.value);
      const low = index > 0 ? Number(thumbs[index - 1].value) : -Infinity;
      const high = index < thumbs.length - 1 ? Number(thumbs[index + 1].value) : Infinity;
      const clamped = Math.min(Math.max(value, low), high);
      if (clamped !== value) thumb.value = clamped;
      this.sync();
    },
    sync() {
      const thumbs = this.thumbs();
      const percents = thumbs.map((thumb) => this.percent(thumb));
      this.values = thumbs.map((thumb) => Number(thumb.value));
      const start = percents.length > 1 ? percents[0] : 0;
      const end = percents.length > 0 ? percents[percents.length - 1] : 0;
      this.$root.style.setProperty("--hui-range-start", start + "%");
      this.$root.style.setProperty("--hui-range-end", end + "%");
      if (thumbs.length < 2) return;
      thumbs.forEach((thumb, index) => {
        thumb.style.zIndex = percents[index] > 50 ? thumbs.length - index : index + 1;
      });
    },
    percent(thumb) {
      const min = Number(thumb.min);
      const max = Number(thumb.max);
      if (max <= min) return 0;
      return Math.round((Number(thumb.value) - min) / (max - min) * 1e5) / 1e3;
    }
  }));
}

// src/switch.js
function toggleSwitch(Alpine) {
  Alpine.data("huiSwitch", (config = {}) => ({
    ...checkable(config),
    init() {
      this.setup();
    },
    destroy() {
      this.teardown();
    },
    toggle() {
      if (config.disabled) return;
      this.checked = !this.checked;
      this.write();
      this.emit();
    }
  }));
}

// src/tabs.js
function tabs(Alpine) {
  Alpine.data("huiTabs", (config = {}) => ({
    value: config.value ?? "",
    init() {
      this.$watch("value", (value) => dispatch(this.$root, "hui-tabs:change", { value }));
    },
    select(value) {
      this.value = value;
    },
    triggers() {
      const prefix = this.$root.id + "-trigger-";
      return Array.from(this.$root.querySelectorAll("[role=tab]")).filter(
        (trigger) => trigger.id.startsWith(prefix) && !trigger.disabled
      );
    },
    keydown(event) {
      const vertical = this.$root.dataset.orientation === "vertical";
      const keys = vertical ? ["ArrowUp", "ArrowDown"] : ["ArrowLeft", "ArrowRight"];
      const triggers = this.triggers();
      const index = triggers.indexOf(event.target);
      if (index === -1) return;
      let next = null;
      if (event.key === keys[0]) next = triggers[(index - 1 + triggers.length) % triggers.length];
      else if (event.key === keys[1]) next = triggers[(index + 1) % triggers.length];
      else if (event.key === "Home") next = triggers[0];
      else if (event.key === "End") next = triggers[triggers.length - 1];
      if (!next) return;
      event.preventDefault();
      next.focus();
      if (!config.manual) this.select(next.dataset.value);
    }
  }));
}

// src/toast.js
function toast(Alpine) {
  Alpine.data("huiToast", (config = {}) => ({
    open: true,
    timer: null,
    started: 0,
    remaining: config.duration ?? 5e3,
    paused: false,
    visibility: null,
    init() {
      this.visibility = () => document.hidden ? this.pause() : this.resume();
      document.addEventListener("visibilitychange", this.visibility);
      this.resume();
    },
    destroy() {
      clearTimeout(this.timer);
      document.removeEventListener("visibilitychange", this.visibility);
    },
    pause() {
      if (!this.remaining || this.paused) return;
      this.paused = true;
      clearTimeout(this.timer);
      this.remaining = Math.max(this.remaining - (Date.now() - this.started), 0);
    },
    resume(event) {
      if (event && event.relatedTarget && this.$root.contains(event.relatedTarget)) return;
      if (!this.remaining || document.hidden || this.$root.matches(":hover")) return;
      this.paused = false;
      this.started = Date.now();
      clearTimeout(this.timer);
      this.timer = setTimeout(() => this.dismiss(), this.remaining);
    },
    dismiss() {
      if (!this.open) return;
      clearTimeout(this.timer);
      const root = this.$root;
      const toaster = root.closest("[data-hui-toaster]");
      const focused = root.contains(document.activeElement);
      this.open = false;
      dispatch(root, "hui-toast:closed");
      if (focused && toaster) {
        const next = Array.from(toaster.children).find((item) => item !== root && item.getClientRects().length);
        if (next) next.focus();
      }
      this.removeWhenHidden(root, 0);
    },
    removeWhenHidden(root, attempt) {
      if (getComputedStyle(root).display === "none" || attempt > 40) {
        root.remove();
        return;
      }
      setTimeout(() => this.removeWhenHidden(root, attempt + 1), 50);
    },
    keydown(event) {
      if (event.key !== "Escape") return;
      event.preventDefault();
      event.stopPropagation();
      this.dismiss();
    }
  }));
}

// src/tooltip.js
function tooltip(Alpine) {
  const stack = layers(Alpine);
  Alpine.data("huiTooltip", (config = {}) => ({
    open: config.open === true,
    active: false,
    blocked: false,
    timer: null,
    anchor: null,
    handled: null,
    layer: null,
    cleanup: null,
    init() {
      this.attach(this.$root.querySelector('[aria-describedby~="' + this.$root.id + '-content"]'));
      this.layer = {
        parts: () => [this.content()],
        inside: () => [this.anchor, this.content()],
        dismiss: () => this.hide(),
        keydown: (event) => this.escape(event)
      };
      this.cleanup = commands(this, "hui-tooltip");
      this.$watch("open", () => this.sync());
      this.$nextTick(() => this.sync());
    },
    destroy() {
      this.cleanup();
      clearTimeout(this.timer);
      this.deactivate();
    },
    content() {
      return document.getElementById(this.$root.id + "-content");
    },
    owns(target) {
      return Boolean(target) && [this.anchor, this.content()].some((element) => element && element.contains(target));
    },
    attach(element) {
      const content = this.content();
      if (!element || element === this.anchor || content && content.contains(element)) return;
      moveAnchor(this.anchor, element, anchorName(this.$root.id));
      addToken(element, "aria-describedby", this.$root.id + "-content");
      this.anchor = element;
    },
    seen(event) {
      if (!event) return false;
      if (this.handled === event) return true;
      this.handled = event;
      return false;
    },
    show() {
      clearTimeout(this.timer);
      this.open = true;
      this.sync();
    },
    hide() {
      clearTimeout(this.timer);
      this.open = false;
      this.sync();
    },
    toggle() {
      this.open ? this.hide() : this.show();
    },
    dismiss() {
      this.blocked = true;
      this.hide();
    },
    later(callback, delay) {
      clearTimeout(this.timer);
      this.timer = setTimeout(callback, delay);
    },
    sync() {
      if (this.open === this.active) return;
      this.active = this.open;
      if (this.open) stack.open(this.layer);
      else stack.close(this.layer);
      dispatch(this.$root, this.open ? "hui-tooltip:opened" : "hui-tooltip:closed");
    },
    deactivate() {
      if (this.active) stack.close(this.layer);
      this.active = false;
    },
    tooltipShow(element) {
      this.attach(element);
      this.show();
    },
    tooltipHide() {
      this.hide();
    },
    tooltipToggle(element, event) {
      if (this.seen(event)) return;
      if (this.open && (!element || element === this.anchor)) this.hide();
      else this.tooltipShow(element);
    },
    tooltipDismiss() {
      this.dismiss();
    },
    tooltipEnter(event, element) {
      if (event.pointerType === "touch" || this.blocked) return;
      this.attach(element);
      if (this.open) {
        clearTimeout(this.timer);
        return;
      }
      this.later(() => this.show(), config.delay ?? 700);
    },
    tooltipLeave(event) {
      if (event.pointerType === "touch" || this.owns(event.relatedTarget)) return;
      this.blocked = false;
      if (!this.open) {
        clearTimeout(this.timer);
        return;
      }
      this.later(() => this.hide(), config.closeDelay ?? 300);
    },
    tooltipFocus(event, element) {
      if (this.blocked || !event.target.matches(":focus-visible")) return;
      this.tooltipShow(element);
    },
    tooltipBlur(event) {
      if (this.owns(event.relatedTarget)) return;
      this.blocked = false;
      this.hide();
    },
    escape(event) {
      if (event.key !== "Escape") return;
      event.preventDefault();
      this.dismiss();
    }
  }));
}

// src/index.js
function headlessui(Alpine) {
  pending();
  Alpine.plugin(checkbox);
  Alpine.plugin(collapsible);
  Alpine.plugin(combobox);
  Alpine.plugin(dropdown);
  Alpine.plugin(fileUpload);
  Alpine.plugin(listbox);
  Alpine.plugin(modal);
  Alpine.plugin(radio);
  Alpine.plugin(range);
  Alpine.plugin(toggleSwitch);
  Alpine.plugin(tabs);
  Alpine.plugin(toast);
  Alpine.plugin(tooltip);
}
export {
  checkbox,
  collapsible,
  combobox,
  headlessui as default,
  dropdown,
  fileUpload,
  listbox,
  modal,
  pending,
  radio,
  range,
  toggleSwitch as switch,
  tabs,
  toast,
  tooltip
};
