import Mousetrap from "mousetrap";
import HotkeysContext, { default_hotkey_register_options } from "./hotkeys_context";
import { hotkeys_context_events, dispatchHotkeysContextEvent } from "./hotkeys_events";
import { Stack, canUseDOMEvents, hasWindowContext } from "@libs/utils";

export class HotkeyContextManager {
    /** @type {Object<string,HotkeysContext>} */
    #contexts
    /** @type {HotkeysContext | null} */
    #current_context
    /** @type {string} current context name */
    #current_context_name
    /** @type {Stack<string>} */
    #context_stack
    /** @type {Boolean} - if true the manager will emit the events defined on hotkeys_consts */ 
    #emit_events
    /** @type {boolean} - true if the context control is locked */
    #context_locked = false
    /** 
     * The name of the last context load requested after the context was locked. as soon as lock is released this context name should be loaded
     * and the the variable should be set to an empty string
     * @type {string} 
     * @default ""
     */
    #last_context_load_requested = ""

    /**
     * 
     * @param {boolean} emit_events - if true the manager will emit the events defined on hotkeys_consts
     */
    constructor(emit_events) {
        this.#contexts = {};
        this.#current_context = null;
        this.#current_context_name = undefined;
        this.#context_stack = new Stack();
        this.#emit_events = emit_events && canUseDOMEvents(); // Only emit events in a browser environment.
    }

    /**
     * The current applied context
     * @type {HotkeysContext}
     * @readonly
    */
    get Context() {
        return this.#current_context;
    }

    /**
     * The current applied context's name
    */
    get ContextName() {
        return this.#current_context_name;
    }

    /**
     * Whether the manager is emitting events
     * @type {boolean}
     */
    get EmitEvents() {
        return this.#emit_events;
    }

    /**
     * Saves a context with a given name but doesn't activate it
     * @param {string} name
     * @param {HotkeysContext} context
    */
    declareContext(name, context) {
        this.#contexts[name] = context;
    }

    /**
     * Deletes a context
     * @param {string} name
    */
    dropContext(name) {
        delete this.#contexts[name];
    }

    /**
     * Emits a context event
     * @private
     * @param {hotkeys_context_events} event
     * @param {import('./hotkeys_events').HotkeysContextEventDetail} detail
     * @returns {void}
     */
    #emitContextEvent(event, detail) {
        if (this.#emit_events) {
            dispatchHotkeysContextEvent(event, detail);
        }
    }

    /**
     * Returns true if the provided context name exists
     * @param {string} name
     * @returns {boolean}
    */
    hasContext(name) {
        return this.#contexts[name] !== undefined;
    }

    /**
     * Returns true if the current context registers the provided hotkey
     * @param {string} hotkey
     * @param {"keypress"|"keydown"|"keyup"} mode
     * @returns {boolean}
    */
    hasHotkey(hotkey, mode="keydown") {
        if (!this.hasLoadedContext()) return false;

        return this.#current_context.hasHotkey(hotkey, mode);
    }

    /**
     * Returns whether there is a loaded context or not
     * @returns {boolean}
     */
    hasLoadedContext() {
        return this.#current_context !== null;
    }

    /**
     * Activates a context. Any previously registered hotkeys will be unregistered no matter if they were registered on the context or by calling registerHotkey
     * @param {string} name
     * @param {boolean} preserve_previous_context if true the previous context will be preserved and can be loaded again with loadPreviousContext
     * @throws {Error} if the context doesn't exist
    */
    loadContext(name, preserve_previous_context=true) {
        if (!hasWindowContext()) return;

        if (this.#context_locked) {
            
            if (this.hasContext(name)) {
                this.#last_context_load_requested = name;
            }

            return;
        }

        let saved_context = this.#contexts[name];

        if (saved_context === undefined) {
            throw new Error(`Context ${name} doesn't exist`);
        }

        // Mousetrap.reset();

        // If preserve_previous_context is true and the current context is not the one on the top of the stack, then add it to the stack. Prevents context duplication
        // also checks that there is a context loaded
        if (preserve_previous_context && this.hasLoadedContext() && this.#current_context_name !== this.#context_stack.Peek()) {
            this.#context_stack.Add(this.#current_context_name);
        }
        
        this.#current_context_name = name;
        this.#current_context = saved_context;

        // Notify the context change
        this.#emitContextEvent(hotkeys_context_events.CONTEXT_CHANGED, {context_name: name});

        this.#registerCurrentContext();
    }
    
    /**
     * Loads the previous context as the current if it exists or does nothing, returns true if the previous context was loaded
     * @returns {boolean} true if the previous context was loaded
    */
    loadPreviousContext() {
        if (!hasWindowContext()) return false;

        let previous_context_name = this.#context_stack.Pop();

        if (previous_context_name === null) return false;

        if (!this.hasContext(previous_context_name)) return false;

        this.loadContext(previous_context_name, false);
        return true;
    }

    /**
     * Locks the context control, preventing other contexts from being loaded. If a context is requested to be loaded while the control is locked, 
     * the context name will be saved and loaded as soon as the control is released.
     * @returns {void}
     */
    lockContextControl() {
        this.#context_locked = true;
    }

    /**
     * @typedef {Object} HotkeyRegisterOptions
     * @property {string} description - The hotkey's description
     * @property {"keypress"|"keydown"|"keyup"} mode - The mode of the keypress event. Default is "keydown"
     */

    /**
     * Registers the current context's hotkeys
     * @throws {Error} if no context is loaded
    */
    #registerCurrentContext() {
        if (!hasWindowContext()) return;

        if (this.#current_context === null) {
            throw new Error("No context loaded");
        }

        for (const hotkey of this.#current_context.hotkeys) {
            hotkey.key_bind()
        }
    }

    /**
     * Registers a hotkey without persisting it on a context, which means that calling loadContext will overwrite it
     * @param {string|string[]} hotkey the hotkey's name or an array of hotkeys, if already registered it will be overwritten
     * @param {function} callback the callback to be called when the key is pressed
     * @param {HotkeyRegisterOptions} options
     * @deprecated
     */
    registerHotkey(hotkey, callback, options) {
        if (!hasWindowContext()) return;

        console.warn("In HotkeyContextManager.registerHotkey: This method is deprecated, use registerHotkeyOnContext instead.");
        
        options = {...default_hotkey_register_options, ...options};
        // Mousetrap.bind(hotkey, callback, options.mode)
    }

    /**
     * Registers a hotkey on the current context
     * @param {string|string[]} hotkey the hotkey's name or an array of hotkeys, if already registered it will be overwritten
     * @param {function} callback the callback to be called when the key is pressed
     * @param {HotkeyRegisterOptions} options
     */
    registerHotkeyOnContext(hotkey, callback, options) {
        if (!hasWindowContext() || !this.hasLoadedContext()) return;

        options = {...default_hotkey_register_options, ...options};
        this.#current_context.register(hotkey, callback, {
            bind: true, 
            description: options.description, 
            mode: options.mode
        });

        this.#emitContextEvent(hotkeys_context_events.CONTEXT_CHANGED, {context_name: this.#current_context_name});
    }

    /**
     * Releases the context control, allowing contexts to be loaded again. If a context load was requested while the control was locked, 
     * the last requested context will be loaded.
     * @returns {void}
     */
    unlockContextControl() {
        this.#context_locked = false;
        let last_context = this.#last_context_load_requested;
        this.#last_context_load_requested = "";

        if (this.hasContext(last_context)) {
            this.loadContext(last_context, false);  
        }
    }

    /**
     * Unregisters a hotkey
     * @param {string|string[]} hotkey the hotkey's name or an array of hotkeys
     * @param {"keypress"|"keydown"|"keyup"} mode the mode of the keypress event
     * @returns {boolean} true if the hotkey was unregistered
     * @deprecated
    */
    unregisterHotkey(hotkey, mode="keydown") {
        if (!hasWindowContext()) return;

        console.warn("In HotkeyContextManager.unregisterHotkey: This method is deprecated, use unregisterHotkeyFromContext instead.");

        // return Mousetrap.unbind(hotkey, mode)
        return true;
    }

    /**
     * Unregisters a hotkey from the current context
     * @param {string|string[]} hotkey the hotkey's name or an array of hotkeys
     * @param {"keypress"|"keydown"|"keyup"} mode the mode of the keypress event
     */
    unregisterHotkeyFromContext(hotkey, mode="keydown") {
        if (!hasWindowContext()) return;

        this.#current_context.unregister(hotkey, mode);

        this.#emitContextEvent(hotkeys_context_events.CONTEXT_CHANGED, {context_name: this.#current_context_name});
    }

    /**
     * Unregisters the current context's hotkeys
     */
    unregisterCurrentContext() {
        if (!hasWindowContext() || !this.hasLoadedContext()) return;

        // Mousetrap.reset();

        this.#current_context = null;
        this.#current_context_name = undefined;

        this.#emitContextEvent(hotkeys_context_events.CONTEXT_CHANGED, {context_name: ""});
    }
}

const global_hotkeys_manager = new HotkeyContextManager(true);

if (globalThis.innerWidth != null) {
    window.global_hotkeys_manager = global_hotkeys_manager;
}

export default global_hotkeys_manager;
