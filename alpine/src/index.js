import checkbox from './checkbox.js'
import collapsible from './collapsible.js'
import combobox from './combobox.js'
import dropdown from './dropdown.js'
import fileUpload from './fileupload.js'
import listbox from './listbox.js'
import modal from './modal.js'
import pending from './pending.js'
import radio from './radio.js'
import range from './range.js'
import toggleSwitch from './switch.js'
import tabs from './tabs.js'
import toast from './toast.js'
import tooltip from './tooltip.js'

export {
	checkbox,
	collapsible,
	combobox,
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
	tooltip,
}

export default function headlessui(Alpine) {
	pending()
	Alpine.plugin(checkbox)
	Alpine.plugin(collapsible)
	Alpine.plugin(combobox)
	Alpine.plugin(dropdown)
	Alpine.plugin(fileUpload)
	Alpine.plugin(listbox)
	Alpine.plugin(modal)
	Alpine.plugin(radio)
	Alpine.plugin(range)
	Alpine.plugin(toggleSwitch)
	Alpine.plugin(tabs)
	Alpine.plugin(toast)
	Alpine.plugin(tooltip)
}
