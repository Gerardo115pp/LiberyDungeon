import { hasWindowContext } from "@libs/utils"
import { HOTKEYS_GENERAL_GROUP } from "./hotkeys_consts"
import Mousetrap from "mousetrap"

/**
* @typedef {Object} HotkeyDataParams
 * @property {string} key_combo
 * @property {function} handler 
 * @property {"keypress"|"keydown"|"keyup"} mode
 * @property {?string} description
*/

export class HotkeyData {
    /**
     * @type {string} the key's name e.g: 'a', 'esc', '-', etc
     */
    name

    /**
     * @type {function} the callback to be called when the key is pressed
     */
    callback

    /**
     * @type {"keypress"|"keydown"|"keyup"} the mode of the keypress event
     */
    mode

    /**
     * @type {string} the hotkey's description
     * @default "<General>No information available"
     */
    #description

    /**
     * @param {string} name the key's name e.g: 'a', 'esc', '-', etc
     * @param {function} callback the callback to be called when the key is pressed
     * @param {"keypress"|"keydown"|"keyup"} mode the mode of the keypress event
     * @param {?string} description The hotkey's description. if the string is prefized with "<group_name>", `group_name` will be the hotkey's group if not it will be "General"
     * @constructor
     */
    constructor(name, callback, mode, description) {
        /** @type {string} the key's name e.g: 'a', 'esc', '-', etc */
        this.name = name
        /** @type {function} the callback to be called when the key is pressed */   
        this.callback = callback
        /** @type {"keypress"|"keydown"|"keyup"} the mode of the keypress event */
        this.mode = mode
        /** @type {string} the hotkey's description */
        this.#description = description ?? "<General>No information available";
    }

    /**
     * The clean hotkey's description string
     * @returns {string}
     */
    get Description() {
        return this.#description.replace(/<.*>/, "");
    }

    /**
     * The hotkey's group
     * @returns {string}
     */
    get Group() {
        let hotkey_group = HOTKEYS_GENERAL_GROUP;
        console.log(this)
        
        /** @type {RegExpMatchArray} */
        let group_matches = this.#description.match(/<(.*)>/);

        if (group_matches != null && group_matches[1] !== undefined) {
            hotkey_group = group_matches[1];
        }

        return hotkey_group;
    }

    /**
     * Registers the hotkey
    */
    key_bind() {
        if (this.mode === "") {
            Mousetrap.bind(this.name, this.callback)
            return
        }
        Mousetrap.bind(this.name, this.callback, this.mode)
    }
    
}

/**
* @typedef {Object} HotkeyRegisterOptions
 * @property {boolean} bind - If true the hotkey will be binded immediately. default is false
 * @property {?string} description - The hotkey's description   
 * @property {"keypress"|"keydown"|"keyup"} mode - The mode of the keypress event. Default is "keydown"
*/

/**
* The default hotkey register options
 * @type {HotkeyRegisterOptions}
 */
export const default_hotkey_register_options = {
    bind: false,
    description: null,
    mode: "keydown"
}

export default class HotkeysContext {
    /** @type{Object<string,HotkeyData>} */
    #keydown_hotkeys
    /** @type {Object<string,HotkeyData>} */
    #keypress_hotkeys
    /** @type {Object<string,HotkeyData>} */
    #keyup_hotkeys
    constructor() {
        this.#keydown_hotkeys = {}
        this.#keypress_hotkeys = {}
        this.#keyup_hotkeys = {}
    }

    /**
     * The Context's hotkeys
     * @type {HotkeyData[]}
    */
    get hotkeys() {
        let all_hotkeys = [];

        all_hotkeys = all_hotkeys.concat(Object.values(this.#keydown_hotkeys))
        all_hotkeys = all_hotkeys.concat(Object.values(this.#keypress_hotkeys))
        all_hotkeys = all_hotkeys.concat(Object.values(this.#keyup_hotkeys))

        return all_hotkeys
    }

    /**
     * returns true if the hotkey is registered on the context
     * @param {string} hotkey
     * @param {"keypress"|"keydown"|"keyup"} mode
     * @returns {boolean}
    */
    hasHotkey(hotkey, mode="keydown") {
        return this.#keydown_hotkeys[hotkey] !== undefined;
    }

    /**
     * Returns the mode's hotkeys
     * @param {"keypress"|"keydown"|"keyup"} mode
    */
    #modeHotkeys(mode) {
        switch (mode) {
            case "keydown":
                return this.#keydown_hotkeys
            case "keypress":
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
     * @param {HotkeyRegisterOptions} options
    */
    register(name, callback, options) {
        if (!hasWindowContext()) return;

        options = {...default_hotkey_register_options, ...options}; // Merge the options with the default options. the 'options' param takes precedence.

        const mode_hotkeys = this.#modeHotkeys(options.mode)

        if (Array.isArray(name)) {
            for (const n of name) {
                mode_hotkeys[n] = new HotkeyData(n, callback, options.mode, options.description)
                if (options.bind) {
                    mode_hotkeys[n].key_bind()
                }
            }
            return
        }

        mode_hotkeys[name] = new HotkeyData(name, callback, options.mode, options.description)
        if (options.bind) {
            mode_hotkeys[name].key_bind();
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
                delete mode_hotkeys[key_name]
            }
            return
        }
        delete mode_hotkeys[name]
    }
}