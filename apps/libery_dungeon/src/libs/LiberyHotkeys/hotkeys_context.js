import { hasWindowContext } from "@libs/utils"
import { HotkeyData, default_hotkey_register_options  } from "./hotkeys"

export default class HotkeysContext {
    /** @type {Map<string, HotkeyData>} */
    #keydown_hotkeys
    /** @type {Map<string, HotkeyData>} */
    #keypress_hotkeys
    /** @type {Map<string, HotkeyData>} */
    #keyup_hotkeys

    constructor() {
        this.#keydown_hotkeys = new Map();
        this.#keypress_hotkeys = new Map();
        this.#keyup_hotkeys = new Map();
    }

    /**
     * The Context's hotkeys
     * @type {HotkeyData[]}
     */
    get hotkeys() {
        let all_hotkeys = Array.from(this.#keydown_hotkeys.values());

        all_hotkeys = all_hotkeys.concat(Array.from(this.#keypress_hotkeys.values()));
        all_hotkeys = all_hotkeys.concat(Array.from(this.#keyup_hotkeys.values()));

        return all_hotkeys
    }

    /**
     * returns true if the hotkey is registered on the context
     * @param {string} hotkey
     * @param {"keypress"|"keydown"|"keyup"} mode
     * @returns {boolean}
     */
    hasHotkey(hotkey, mode="keydown") {
        let mode_hotkeys = this.#modeHotkeys(mode)

        return mode_hotkeys.has(hotkey);
    }

    /**
     * Returns the mode's hotkeys
     * @param {"keypress"|"keydown"|"keyup"} mode
     * @returns {Map<string, HotkeyData>}
     */
    #modeHotkeys(mode) {
        switch (mode) {
            case "keydown":
                return this.#keydown_hotkeys
            case "keypress": // TODO: Remove support for keypress as it is deprecated and will slowly be removed from browsers.
                console.warn("keypress is deprecated and will be removed from browsers, use keydown instead")
                return this.#keypress_hotkeys
            case "keyup":
                return this.#keyup_hotkeys
        }
    }

    /**
     * Register a hotkey on the context, if the hotkey already exists it will be overwritten.
     * If a array of hotkeys is passed, all of them will be registered with the same callback
     * @param {string|string[]} name
     * @param {function} callback
     * @param {import('./hotkeys').HotkeyRegisterOptions} options
     */
    register(name, callback, options) {
        if (!hasWindowContext()) return;

        options = {...default_hotkey_register_options, ...options}; // Merge the options with the default options. the 'options' param takes precedence.

        const names = Array.isArray(name) ? name : [name]

        for (const n of names) {
            this.#saveHotkey(n, callback, options)
        }
        return
    }

    /**
     * Saves a given hotkey name as a HotkeyData in the appropriate mode.
     * @param {string} name
     * @param {function} callback
     * @param {import('./hotkeys').HotkeyRegisterOptions} options
     */
    #saveHotkey(name, callback, options) {
        const mode_hotkeys = this.#modeHotkeys(options.mode)
        const new_hotkey = new HotkeyData(name, callback, options)
        mode_hotkeys.set(name, new_hotkey)
        if (options.bind) {
            console.warn("The bind option is deprecated and will be removed soon.");
            new_hotkey.key_bind()
        }
    }

    /**
     * Unregister a hotkey on the context
     * @param {string|string[]} name could be a key name or an array of key names(keyboard keys)
     * @param {"keypress"|"keydown"|"keyup"} mode
     */
    unregister(name, mode) {
        if (!hasWindowContext()) return;

        const mode_hotkeys = this.#modeHotkeys(mode) 

        if (Array.isArray(name)) {
            for (const key_name of name) {
                mode_hotkeys.delete(key_name)
            }
            return
        }

        mode_hotkeys.delete(name)
    }
}
