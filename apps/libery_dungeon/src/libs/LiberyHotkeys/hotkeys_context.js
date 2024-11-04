import { hasWindowContext } from "@libs/utils"
import { HotkeyData, default_hotkey_register_options  } from "./hotkeys"
import { HOTKEY_NULLISH_HANDLER } from "./hotkeys_consts"

/**
* CHORE: Rename all references to key_combo to hotkey_trigger
* CHORE: Rename change all jsdoc for hotkey callbacks from function to HotkeyCallback
*/

/**
 * Arguments to register a hotkey. each hotkey has one and only one trigger. But the HotkeyContext.register method accepts an  array of triggers but a different HotkeyData is created for each trigger.
 * This is the main difference between the HotkeyDataParams and HotkeyRegisterParams.
* @typedef {Object} HotkeyRegisterParams
 * @property {string | string[]} hotkey_triggers
 * @property {import('./hotkeys').HotkeyCallback} callback
 * @property {import('./hotkeys').HotkeyRegisterOptions} options
*/

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
     * @param {import('./hotkeys').HotkeyCallback} callback
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
     * @param {import('./hotkeys').HotkeyCallback} callback
     * @param {import('./hotkeys').HotkeyRegisterOptions} options
     */
    #saveHotkey(name, callback, options) {
        const mode_hotkeys = this.#modeHotkeys(options.mode ?? "keydown")
        const new_hotkey = new HotkeyData(name, callback, options)

        if (!new_hotkey.Valid) {
            return
        }

        mode_hotkeys.set(name, new_hotkey)
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

/**
 * Metadata about a hotkey action
* @typedef {Object} HotkeyActionMeta
 * @property {Symbol} overwrite_behavior - how the component declares that will handle hotkey callback overriding for this action.
 * @property {boolean} overwritten_trigger
 * @property {boolean} overwritten_callback
 * @property {boolean} overwritten_options
 * @property {HotkeyRegisterParams} hotkey_register_params
*/

/**
 * A wrapper class for the HotkeysContext. It's purpose is for components to expose the hotkeys it handles and allow a parent component to register over them or request a given action is taken when a hotkey triggers.
 * Whether an action handler gets replaced, wrapped or ignored is up to the actual component. For most cases, HotkeysContext is a better choice than ComponentHotkeyContext. ComponentHotkeyContext is a specialized version 
 * of HotkeysContext created for cases when you need to grant hotkey control to deeply nested child components and the hotkey action requires data from more than one component.
 */
export class ComponentHotkeyContext {

    /**
     * A component declares overwriting the hotkey callback for this action will replace the default behavior.
     * Typically this means that the component will not add a handler for the hotkey.  
     */
    static OVERRIDE_BEHAVIOR_REPLACE = Symbol("OVERRIDE_BEHAVIOR_REPLACE");

    /**
     * A component declares overwriting the hotkey callback for this action will cause the component's own handler to call the overwritten handler after it's done. 
     */
    static OVERRIDE_BEHAVIOR_WRAP = Symbol("OVERRIDE_BEHAVIOR_WRAP");

    /**
     * A component declares overwriting the hotkey callback for this action will be ignored. Meaning it will not accept changes for this action.
     */
    static OVERRIDE_BEHAVIOR_IGNORE = Symbol("OVERRIDE_BEHAVIOR_IGNORE");

    /** 
     * The hotkeys context for the component
     * @type {HotkeysContext | null} 
     */
    #hotkeys_context

    /** 
     * The actions that the component will handle.
     * @type {Map<string, HotkeyActionMeta>}
     */
    #actions

    /**
     * The name of the hotkey context
     * @type {string}
     */
    #hotkeys_context_name

    /**
     * Whether the ComponentHotkeyContext should accept any more actions. Even if set to true, it will accept overwriting of existing actions.
     * Just not new actions.
     * @type {boolean}
     */
    #final;

    /**
     * @param {string} hotkeys_context_name
     */
    constructor(hotkeys_context_name) {
        this.#hotkeys_context = null;
        this.#actions = new Map();
        this.#hotkeys_context_name = hotkeys_context_name;
        this.#final = false;
    }

    /**
     * Returns the overwriting behavior for a given action name. If the action does not exist, it returns undefined.
     * @param {string} action_name
     * @returns {Symbol | undefined}
     */
    checkOverwriteBehavior(action_name) {
        let action = this.#actions.get(action_name);

        if (action == null) {
            return undefined;
        }

        return action.overwrite_behavior;
    }

    /**
     * Drops the hotkeys context. If the hotkeys context has already been dropped, the method panics.
     * @returns {void}
     */
    dropHotkeysContext() {
        if (this.#hotkeys_context == null) {
            throw new Error(`The ComponentHotkeyContext '${this.#hotkeys_context_name}' has already dropped its hotkeys context.`);
        }

        this.#hotkeys_context = null;
    }

    /**
     * whether the ComponentHotkeyContext is final.
     * @type {boolean}
     */
    get Final() {
        return this.#final;
    }

    /**
     * Generates a HotkeysContext based on the ComponentHotkeyContext's actions. If the hotkeys context has already been generated, it has to be dropped first or the method panics.
     * @returns {HotkeysContext}
     */
    generateHotkeysContext() {
        if (this.#hotkeys_context instanceof HotkeysContext) {
            throw new Error(`The ComponentHotkeyContext '${this.#hotkeys_context_name}' has already generated its hotkeys context.`);
        }

        this.#hotkeys_context = new HotkeysContext();

        for (const [action_name, hotkey_action] of this.#actions) {
            this.#hotkeys_context.register(
                hotkey_action.hotkey_register_params.hotkey_triggers,
                hotkey_action.hotkey_register_params.callback,
                hotkey_action.hotkey_register_params.options
            );
        }

        return this.#hotkeys_context;        
    }
   
    /**
     * The hotkeys context name the component should use.
     *  @type {string}
     */
    get HotkeysContextName() {
        return this.#hotkeys_context_name;
    }

    /**
     * The hotkeys context the component should use.
     * @type {HotkeysContext | null}
     */
    get HotkeysContext() {
        return this.#hotkeys_context;
    }

    /**
     * Returns whether the ComponentHotkeyContext has a given action.
     * @param {string} action_name
     * @returns {boolean}
     */
    hasAction(action_name) {
        if (!this.#actions.has(action_name)) {
            return false;
        };

        let action = this.#actions.get(action_name);

        return action != null;
    }

    /**
     * Overwrites the hotkey trigger for a given action. The action name must exist in the ComponentHotkeyContext already or the method panics
     * @param {string} action_name
     * @param {string | string[]} hotkey_trigger
     */
    overwriteHotkeyTrigger(action_name, hotkey_trigger) {
        const hotkey_action = this.#actions.get(action_name);

        if (hotkey_action == null) {
            throw new Error(`The action ${action_name} does not exist in the ComponentHotkeyContext.`);
        }

        hotkey_action.hotkey_register_params.hotkey_triggers = hotkey_trigger;
    }

    /**
     * overwrites the hotkey callback for a given action. The action name must exist in the ComponentHotkeyContext already or the method panics.
     * @param {string} action_name
     * @param {import('./hotkeys').HotkeyCallback} callback
     */
    overwriteHotkeyCallback(action_name, callback) {
        if (typeof callback !== "function") {
            throw new Error(`In LiberyHotkeys/hotkeys_context ComponentHotkeyContext.overwriteHotkeyCallback: Attempted to set a non-function as a handler for the action ${action_name}`);
        }

        const hotkey_data_params = this.#actions.get(action_name);

        if (hotkey_data_params == null) {
            throw new Error(`The action ${action_name} does not exist in the ComponentHotkeyContext.`);
        }

        hotkey_data_params.hotkey_register_params = {
            ...hotkey_data_params.hotkey_register_params,
            callback: callback
        };
    }

    /**
     * Overwrites the hotkey options for a given action. The action name must exist in the ComponentHotkeyContext already or the method panics.
     * @param {string} action_name
     * @param {import('./hotkeys').HotkeyRegisterOptions} options
     */
    overwriteHotkeyOptions(action_name, options) {
        const hotkey_data_params = this.#actions.get(action_name);

        if (hotkey_data_params == null) {
            throw new Error(`The action ${action_name} does not exist in the ComponentHotkeyContext.`);
        }

        /**
         * @type {import('./hotkeys').HotkeyRegisterOptions}
         */
        let safe_options = { ...default_hotkey_register_options, ...options };

        hotkey_data_params.hotkey_register_params.options = safe_options;
    }

    /**
     * Registers hotkey data parameters for a given action.
     * @param {string} action_name
     * @param {HotkeyActionMeta} hotkey_action
     * @returns {void}
     */
    registerHotkeyAction(action_name, hotkey_action) {
        if (this.#final) {
            throw new Error(`The ComponentHotkeyContext '${this.#hotkeys_context_name}' is final and cannot accept new actions.`);
        }

        const hotkey_register_params = hotkey_action?.hotkey_register_params ?? {};
        const hotkey_register_options = hotkey_register_params.options ?? {};

        if (typeof hotkey_register_params.callback !== "function" || hotkey_register_params.hotkey_triggers == null) {
            throw new Error(`In LiberyHotkeys/hotkeys_context ComponentHotkeyContext.registerHotkeyAction: Attempted to register an action with missing parameters. The action name is ${action_name}`);
        }

        /**
         * @type {HotkeyActionMeta}
         */
        let safe_hotkey_action = { 
            ...default_hotkey_action_meta, 
            ...hotkey_action,
            hotkey_register_params: {
                ...hotkey_register_params,
                options: {
                    ...default_hotkey_register_options,
                    ...hotkey_register_options
                }
            }
        };

        this.#actions.set(action_name, hotkey_action);
    }

    /**
     * Sets the ComponentHotkeyContext to be final. 
     */
    SetFinal() {
        this.#final = true;
    }
}

/**
 * The default hotkey action meta
 * @type {HotkeyActionMeta} 
 */
const default_hotkey_action_meta = {
    overwrite_behavior: ComponentHotkeyContext.OVERRIDE_BEHAVIOR_IGNORE,
    overwritten_trigger: false,
    overwritten_callback: false,
    overwritten_options: false,
    hotkey_register_params: {
        hotkey_triggers: [],
        callback: HOTKEY_NULLISH_HANDLER,
        options: {
            ...default_hotkey_register_options
        }
    }
}