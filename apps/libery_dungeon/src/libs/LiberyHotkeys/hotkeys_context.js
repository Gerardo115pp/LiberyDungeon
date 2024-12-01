import { hasWindowContext } from "@libs/utils"
import { HotkeyData, default_hotkey_register_options  } from "./hotkeys"
import { HOTKEY_NULL_DESCRIPTION, HOTKEY_NULLISH_HANDLER } from "./hotkeys_consts"

/**
* CHORE: change all jsdoc for hotkey callbacks from function to HotkeyCallback
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
     * @param {"keydown"|"keyup"} mode
     * @returns {boolean}
     */
    hasHotkey(hotkey, mode="keydown") {
        let mode_hotkeys = this.#modeHotkeys(mode)

        return mode_hotkeys.has(hotkey);
    }

    /**
     * Returns if a hotkey trigger is registered on the context. different from hasHotkey, this handles strings and arrays of strings.
     * The method returns true if and only if all the hotkey combos are not registered on the context. 
     * @param {string | string[]} hotkey_trigger
     * @param {"keydown"|"keyup"} mode
     * @returns {boolean}
     */
    hasHotkeyTrigger(hotkey_trigger, mode="keydown") {
        if (!Array.isArray(hotkey_trigger)) {
            return this.hasHotkey(hotkey_trigger, mode);
        }

        let has_hotkey = false;

        for (let h = 0; h < hotkey_trigger.length && !has_hotkey; h++) {
            has_hotkey = this.hasHotkey(hotkey_trigger[h], mode);
        }

        return has_hotkey;
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
 * @property {boolean} [overwritten_trigger]
 * @property {boolean} [overwritten_callback]
 * @property {boolean} [overwritten_options]
 * @property {HotkeyRegisterParams} hotkey_register_params
*/

export class HotkeyAction {
    /**
     * Overwriting behavior for the hotkey action
     * @type {Symbol}
     */
    #overwrite_behavior

    /**
     * Whether the hotkey trigger has been overwritten
     * @type {boolean}
     * @default false
     */
    #overwritten_trigger

    /**
     * Whether the hotkey callback has been overwritten
     * @type {boolean}
     * @default false
     */
    #overwritten_callback


    /**
     * Whether the hotkey options have been overwritten
     * @type {boolean}
     * @default false
     */
    #overwritten_options

    /**
     * The hotkey register parameters
     * @type {HotkeyRegisterParams}
     */
    #hotkey_register_params
    
    /**
     * @param {HotkeyActionMeta} hotkey_action_meta
     */
    constructor(hotkey_action_meta) {
        this.#overwrite_behavior = hotkey_action_meta.overwrite_behavior;

        this.#overwritten_trigger = false;
        this.#overwritten_callback = false;
        this.#overwritten_options = false;

        this.#hotkey_register_params = hotkey_action_meta.hotkey_register_params
    }

    /**
     * The hotkey callback
     * @type {import('./hotkeys').HotkeyCallback}
     */
    get Callback() {
        return this.#hotkey_register_params.callback;
    }

    /**
     * The hotkey callback
     * @param {import('./hotkeys').HotkeyCallback} callback
     * @returns {void}
     */
    set Callback(callback) {
        this.#hotkey_register_params.callback = callback;
        this.#overwritten_callback = true;
    }

    /**
     * The hotkey triggers.
     * @type {string | string[]}
     */
    get HotkeyTriggers() {
        return this.#hotkey_register_params.hotkey_triggers;
    }

    /**
     * The hotkey triggers.
     * @param {string | string[]} hotkey_triggers
     * @returns {void}
     */
    set HotkeyTriggers(hotkey_triggers) {
        this.#hotkey_register_params.hotkey_triggers = hotkey_triggers;
        this.#overwritten_trigger = true;
    }

    /**
     * Whether the callback is nullish
     * @returns {boolean}
     */
    HasNullishCallback() {
        return this.#hotkey_register_params.callback === HOTKEY_NULLISH_HANDLER;
    }

    /**
     * Whether the description is nullish
     * @returns {boolean}
     */
    HasNullishDescription() {
        return this.#hotkey_register_params.options.description === HOTKEY_NULL_DESCRIPTION;
    }

    /**
     * The hotkey options
     * @type {import('./hotkeys').HotkeyRegisterOptions}
     */
    get Options() {
        return this.#hotkey_register_params.options;
    }

    /**
     * The hotkey options
     * @param {import('./hotkeys').HotkeyRegisterOptions} options
     * @returns {void}
     */
    set Options(options) {
        this.#hotkey_register_params.options = options;
        this.#overwritten_options = true;
    }

    /**
     * The overwriting behavior for the hotkey action
     * @type {Symbol}
     */
    get OverwriteBehavior() {
        return this.#overwrite_behavior;
    }

    /**
     * Whether the hotkey trigger has been overwritten
     * @type {boolean}
     */
    get OverwrittenTrigger() {
        return this.#overwritten_trigger;
    }

    /**
     * Whether the hotkey callback has been overwritten
     * @type {boolean}
     * @default false
     */
    get OverwrittenCallback() {
        return this.#overwritten_callback;
    }

    /**
     * Whether the hotkey options have been overwritten
     * @type {boolean}
     */
    get OverwrittenOptions() {
        return this.#overwritten_options;
    }

    /**
     * Overwrites the hotkey description. does not set the overwritten options flag.
     * @param {string} new_description
     */
    overwriteDescription(new_description) {
        this.#hotkey_register_params.options.description = new_description;
    }

    /**
     * Panics if the HotkeyRegisterParams.options have been overwritten. Description overwriting is ignored.
     * @param {string} [error_message]
     * @returns {void}
     */
    panicIfOptionsOverwritten(error_message) {
        if (!this.#overwritten_options) return;

        error_message = error_message ?? "Overwriting the options for this action is not acceptable by the owner component.";

        throw new Error(error_message);
    }

    /**
     *  Panics if the HotkeyRegisterParams.callback has been overwritten.
     * @param {string} [error_message]
     * @returns {void}
     */
    panicIfCallbackOverwritten(error_message) {
        if (!this.#overwritten_callback) return;

        error_message = error_message ?? "Overwriting the callback for this action is not acceptable by the owner component.";

        throw new Error(error_message);
    }

    /**
     * Panics if the HotkeyRegisterParams.hotkey_triggers has been overwritten.
     * @param {string} [error_message]
     * @returns {void}
     */
    panicIfTriggerOverwritten(error_message) {
        if (!this.#overwritten_trigger) return;

        error_message = error_message ?? "Overwriting the trigger for this action is not acceptable by the owner component.";

        throw new Error(error_message);
    }
}

/**
 * A callback for when the ComponentHotkeyContext.active property changes.
* @callback ActiveChangeCallback
 * @param {boolean} active
 * @returns {void}
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
     * An optional reference to the a parent's hotkeys context. This is useful to streamline the nesting of components that use their own hotkeys context.
     * @type {ComponentHotkeyContext | null}
     */
    #parent_hotkeys_context;

    /**
     * A map of child hotkey contexts.
     * @type {Map<string, ComponentHotkeyContext>}
     */
    #child_hotkeys_contexts;

    /** 
     * The hotkeys context for the component
     * @type {HotkeysContext | null} 
     */
    #hotkeys_context

    /** 
     * The actions that the component will handle.
     * @type {Map<Symbol, HotkeyAction>}
     */
    #actions

    /**
     * A low priority array of hotkeys a parent component wants the child component to register on their hotkeys context. This hotkeys are registered after the child component's hotkeys and if a conflict is 
     * found, the hotkey is ignored and an error is logged.
     * @type {HotkeyRegisterParams[]}
     */
    #extra_hotkeys;

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
     * Whether the hotkey context has been loaded, this shall be managed by the component itself and should be true if the component has loaded the hotkeys context onto the hotkey context manager, even if the hotkey context 
     * loaded in the hotkey binder is not the component's hotkey context but a child of it. it only returns to false if the component has dropped the hotkeys context and returned the hotkey control back to the parent.
     * @type {boolean}
     */
    #active;

    /**
     * A callback triggered whenever the #active property changes.
     * @type {ActiveChangeCallback[]}
     */
    #active_state_change_callbacks;

    /**
     * @param {string} hotkeys_context_name
     */
    constructor(hotkeys_context_name) {
        this.#parent_hotkeys_context = null;
        this.#child_hotkeys_contexts = new Map();
        this.#hotkeys_context = null;
        this.#actions = new Map();
        this.#extra_hotkeys = [];
        this.#hotkeys_context_name = hotkeys_context_name;
        this.#final = false;
        this.#active = false;
        this.#active_state_change_callbacks = [];
    }

    /**
     * Whether the hotkey context has been loaded, this shall be managed by the component itself and should be true if the component has loaded the hotkeys context onto the hotkey context manager, even if the hotkey context
     * loaded in the hotkey binder is not the component's hotkey context but a child of it. it only returns to false if the component has dropped the hotkeys context and returned the hotkey control back to the parent.
     * @type {boolean}
     */
    get Active() {
        return this.#active;
    }

    /**
     * Whether the hotkey context has been loaded, this shall be managed by the component itself and should be true if the component has loaded the hotkeys context onto the hotkey context manager, even if the hotkey context
     * loaded in the hotkey binder is not the component's hotkey context but a child of it. it only returns to false if the component has dropped the hotkeys context and returned the hotkey control back to the parent.
     * @param {boolean} active
     */
    set Active(active) {
        if (this.#active === active) return;

        this.#active = active;
        
        for (const callback of this.#active_state_change_callbacks) {
            callback(active);
        }
    }

    /**
     * Applies the extra hotkeys, if any, to the generated hotkeys context. which means this method has to be called after the hotkeys context has been generated otherwise it will panic.
     * @returns {void}
     */
    applyExtraHotkeys() {
        if (this.#hotkeys_context == null) {
            throw new Error(`The ComponentHotkeyContext '${this.#hotkeys_context_name}' has not generated its hotkeys context yet.`);
        }

        for (const hotkey_params of this.#extra_hotkeys) {
           if (this.#hotkeys_context.hasHotkeyTrigger(hotkey_params.hotkey_triggers, hotkey_params.options.mode)) {
                console.error(`The hotkey ${hotkey_params.hotkey_triggers} is already registered on the hotkeys context. Ignoring...`);
                continue;
           }

            this.#hotkeys_context.register(hotkey_params.hotkey_triggers, hotkey_params.callback, hotkey_params.options);
        }
    }

    /**
     * Adds a child hotkeys context to the ComponentHotkeyContext.
     * @param {ComponentHotkeyContext} child_hotkeys_context
     */
    addChildContext(child_hotkeys_context) {
        this.#child_hotkeys_contexts.set(child_hotkeys_context.HotkeysContextName, child_hotkeys_context);
    }

    /**
     * Copies all the extra hotkey registered on this component hotkey context to another.
     * @param {ComponentHotkeyContext} other
     */
    cloneExtraHotkeys(other) {
        for (let h = 0; h < this.#extra_hotkeys.length; h++) {
            other.registerExtraHotkey(this.#extra_hotkeys[h]);
        }
    }

    /**
     * Returns the overwriting behavior for a given action name. If the action does not exist, it returns undefined.
     * @param {Symbol} action_name
     * @returns {Symbol | undefined}
     */
    checkOverwriteBehavior(action_name) {
        let action = this.#actions.get(action_name);

        if (action == null) {
            return undefined;
        }

        return action.OverwriteBehavior;
    }

    /**
     * A map of child hotkey contexts. this is useful when a component has multiple child components that need to handle their own and we need to know about them from a non-direct parent hotkeys context.
     * @type {Map<string, ComponentHotkeyContext>}
     */
    get ChildHotkeysContexts() {
        return this.#child_hotkeys_contexts;
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
     * Returns a hotkey action for a given action name. If the action does not exist, it returns undefined.
     * @param {Symbol} action_name
     * @returns {HotkeyAction | undefined}
     */
    getHotkeyAction(action_name) {
        return this.#actions.get(action_name);
    }

    /**
     * Returns a hotkey action for a given action name. If the action does not exist, it panics.
     * @param {Symbol} action_name
     * @returns {HotkeyAction}
     */
    getHotkeyActionOrPanic(action_name) {
        let action = this.#actions.get(action_name);

        if (action == null) {
            throw new Error(`The action ${action_name} does not exist but it's required by the calling component.`);
        }

        return action;
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
            if (hotkey_action.HasNullishCallback() || hotkey_action.HasNullishDescription()) {
                console.warn(`The action ${action_name.toString()} has a nullish callback or description. Ignoring...`);
                continue;
            }

            this.#hotkeys_context.register(
                hotkey_action.HotkeyTriggers,
                hotkey_action.Callback,
                hotkey_action.Options
            );
        }

        return this.#hotkeys_context;        
    }

    /**
     * Returns the hotkey options for the given action. If the action does not exist, it returns undefined.
     * @param {Symbol} action_name
     * @returns {import('./hotkeys').HotkeyRegisterOptions | undefined}
     */
    getActionOptions(action_name) {
        let action = this.#actions.get(action_name);

        return action?.Options;    
    }

    /**
     * Returns the hotkey callback for the given action. If the action does not exist, it returns a nullish handler.
     * @param {Symbol} action_name
     * @returns {import('./hotkeys').HotkeyCallback}
     */
    getActionCallback(action_name) {
        let action = this.#actions.get(action_name);

        return action?.Callback ?? HOTKEY_NULLISH_HANDLER;
    }

    /**
     * Returns the hotkey trigger for the given action. If the action does not exist, it returns an empty array.
     * @param {Symbol} action_name
     * @returns {string | string[]}
     */
    getActionTrigger(action_name) {
        let action = this.#actions.get(action_name);

        return action?.HotkeyTriggers ?? [];
    }

    /**
     * Returns the hotkey context name of the last activated child hotkeys context in a Component hotkey context chain . e.g: consider a component hotkey context chain like A(active)->B(active)->C(active)->D(unactive).
     * The method will return the hotkey context name of C because it's active but non of it's children have activated their hotkeys context. If the component hotkey context is active but has not active childs, 
     * then returns it's own name. If the component hotkey context is not active then it returns an empty string.
     * @returns {string}
     */
    getLastActiveChildContextName() {
        if (!this.#active) return "";

        let active_child_context_name = this.#hotkeys_context_name;

        for (const [child_context_name, child_context] of this.#child_hotkeys_contexts) {
            if (child_context.Active) {
                active_child_context_name = child_context.getLastActiveChildContextName();
                break
            }
        }

        return active_child_context_name;
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
     * @param {Symbol} action_name
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
     * Whether the ComponentHotkeyContext has generated a hotkeys context.
     * @returns {boolean}
     */
    HasGeneratedHotkeysContext() {
        return this.#hotkeys_context != null;
    }

    /**
     * Passed down the extra hotkeys to all the child hotkey contexts.
     */
    inheritExtraHotkeys() {
        for (const [child_context_name, child_context] of this.#child_hotkeys_contexts) {
            this.cloneExtraHotkeys(child_context);
        }
    }

    /**
     * Overwrites the hotkey trigger for a given action. The action name must exist in the ComponentHotkeyContext already or the method panics
     * @param {Symbol} action_name
     * @param {string | string[]} hotkey_trigger
     */
    overwriteHotkeyTrigger(action_name, hotkey_trigger) {
        const hotkey_action = this.#actions.get(action_name);

        if (hotkey_action == null) {
            throw new Error(`The action ${action_name} does not exist in the ComponentHotkeyContext.`);
        }

        hotkey_action.HotkeyTriggers = hotkey_trigger;
    }

    /**
     * overwrites the hotkey callback for a given action. The action name must exist in the ComponentHotkeyContext already or the method panics.
     * @param {Symbol} action_name
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

        hotkey_data_params.Callback = callback;
    }

    /**
     * Overwrites the hotkey options for a given action. The action name must exist in the ComponentHotkeyContext already or the method panics.
     * @param {Symbol} action_name
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

        hotkey_data_params.Options = safe_options;
    }

    /**
     * Overwrites the hotkey description for a given action. Does not set the overwritten options flag.
     * @param {Symbol} action_name
     * @param {string} new_description
     * @returns {void}
     */
    overwriteHotkeyDescription(action_name, new_description) {
        const hotkey_action = this.#actions.get(action_name);

        if (hotkey_action == null) {
            throw new Error(`The action ${action_name} does not exist in the ComponentHotkeyContext.`);
        }

        hotkey_action.overwriteDescription(new_description);
    }

    /**
     * Set a callback to be triggered whenever the component's active state changes.
     * @param {(active: boolean) => void} callback
     */
    onActiveChange(callback) {
        for (const c of this.#active_state_change_callbacks) {
            if (c === callback) {
                console.warn("In ComponentHotkeyContext.onActiveChange: Attempted to add the same callback twice.");
                return;
            }
        }

        this.#active_state_change_callbacks.push(callback);
    }

    /**
     * An optional reference to the a parent's hotkeys context. This is useful to streamline the nesting of components that use their own hotkeys context.
     * @type {ComponentHotkeyContext | null}
     */
    get ParentHotkeysContext() {
        return this.#parent_hotkeys_context;
    }

    /**
     * An optional reference to the a parent's hotkeys context. This is useful to streamline the nesting of components that use their own hotkeys context.
     * @param {ComponentHotkeyContext} parent_hotkeys_context
     */
    set ParentHotkeysContext(parent_hotkeys_context) {
        this.#parent_hotkeys_context = parent_hotkeys_context;
    }
    
    /**
     * Panics if the HotkeyRegisterParams.options have been overwritten for a given action. description overwriting is ignored.
     * @param {Symbol} action_name
     * @param {string} [error_message]
     * @returns {void}
     */
    panicIfOptionsOverwritten(action_name, error_message) {
        const hotkey_action = this.getHotkeyActionOrPanic(action_name);

        hotkey_action.panicIfOptionsOverwritten(error_message);
    }

    /**
     * Panics if the HotkeyRegisterParams.callback has been overwritten for a given action.
     * @param {Symbol} action_name
     * @param {string} [error_message]
     * @returns {void}
     */
    panicIfCallbackOverwritten(action_name, error_message) {
        const hotkey_action = this.getHotkeyActionOrPanic(action_name);

        hotkey_action.panicIfCallbackOverwritten(error_message);
    }

    /**
     * Panics if the HotkeyRegisterParams.hotkey_triggers has been overwritten for a given action.
     * @param {Symbol} action_name
     * @param {string} [error_message]
     * @returns {void}
     */
    panicIfTriggerOverwritten(action_name, error_message) {
        const hotkey_action = this.getHotkeyActionOrPanic(action_name);

        hotkey_action.panicIfTriggerOverwritten(error_message);
    }

    /**
     * Registers hotkey data parameters for a given action.
     * @param {Symbol} action_name
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

        const new_hotkey_action = new HotkeyAction(safe_hotkey_action);

        this.#actions.set(action_name, new_hotkey_action);
    }

    /**
     * Registers an extra hotkey that a child component can register using methods and references from the parent component. This is ideal when you want to maintain certain functionality of the parent component in a child's
     * hotkeys context. This is completely safe to call after Final has been set.
     * @param {HotkeyRegisterParams} hotkey_params
     * @returns {void}
     */
    registerExtraHotkey(hotkey_params) {
        this.#extra_hotkeys.push(hotkey_params);
    }

    /**
     * Removes a on active change callback.
     * @param {ActiveChangeCallback} callback
     */
    removeOnActiveChange(callback) {
        this.#active_state_change_callbacks = this.#active_state_change_callbacks.filter(c => c !== callback);
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