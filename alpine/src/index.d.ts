export interface Alpine {
	plugin(callback: (Alpine: any) => void): void
	data(name: string, callback: (...args: any[]) => any): void
	directive(name: string, callback: (...args: any[]) => any): any
	magic(name: string, callback: (...args: any[]) => any): void
}

export type AlpinePlugin = (Alpine: Alpine) => void

declare const hui: AlpinePlugin
export default hui

export declare const checkbox: AlpinePlugin
export declare const collapsible: AlpinePlugin
export declare const combobox: AlpinePlugin
export declare const dropdown: AlpinePlugin
export declare const fileUpload: AlpinePlugin
export declare const listbox: AlpinePlugin
export declare const modal: AlpinePlugin
export declare const radio: AlpinePlugin
export declare const range: AlpinePlugin
export declare const tabs: AlpinePlugin
export declare const toast: AlpinePlugin
export declare const tooltip: AlpinePlugin
export declare function pending(): void

declare const toggleSwitch: AlpinePlugin
export { toggleSwitch as switch }
