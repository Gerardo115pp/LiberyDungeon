import HotkeysContext from "./hotkeys_context";
import  { default_hotkey_register_options } from "./hotkeys";
import { hotkeys_context_events, dispatchHotkeysContextEvent } from "./hotkeys_events";
import { Stack, canUseDOMEvents, hasWindowContext } from "@libs/utils";
import { destroyHotkeysBinder, getHotkeysBinder, setupHotkeysBinder } from "./hotkeys_binder";

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
    #last_context_load_requested = "";
    /**
     * The hotkeys key binder
     * @type {import('./hotkeys_binder').HotkeysController}
     */
    #hotkeys_controller
    /**
     * Whether the manager is destroyed and should no longer be used
     * @type {boolean}
     */
    #destroyed = false

    /**
     * 
     * @param {boolean} emit_events - if true the manager will emit the events defined on hotkeys_consts
     */
    constructor(emit_events) {
        console.log("Creating hotkeys manager!");
        let global_hotkeys_binder = null;
        if (getHotkeysBinder() !== null) {
            destroyHotkeysBinder();
        }

        this.#contexts = {};
        this.#current_context = null;
        this.#current_context_name = undefined;
        this.#context_stack = new Stack();
        this.#emit_events = emit_events && canUseDOMEvents(); // Only emit events in a browser environment.

        if (this.#emit_events) {
            setupHotkeysBinder();
            global_hotkeys_binder = getHotkeysBinder();
        }

        this.#hotkeys_controller = global_hotkeys_binder;
    }

    /**
     * the hotkeys binder.
     * @type {import('./hotkeys_binder').HotkeysController}
     */
    get Binder() {
        return this.#hotkeys_controller;
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
     * Destroys the hotkeys manager. Subsequent attempts to use the manager will panic
     * @returns {void}
     */
    destroy() {
        console.log("Destroying hotkeys manager");
        if (this.#emit_events) {
            destroyHotkeysBinder();
        }

        this.#contexts = {};
        this.#current_context = null;
        this.#current_context_name = undefined;
        this.#context_stack = null;
        this.#hotkeys_controller = null;
    }

    /**
     * Saves a context with a given name but doesn't activate it
     * @param {string} name
     * @param {HotkeysContext} context
     */
    declareContext(name, context) {
        this.#panicIfDestroyed();
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
     * Whether the manager is emitting events
     * @type {boolean}
     */
    get EmitEvents() {
        return this.#emit_events;
    }

    /**
     * Emits a context event
     * @private
     * @param {hotkeys_context_events} event
     * @param {import('./hotkeys_events').HotkeysContextEventDetail} detail
     * @returns {void}
     */
    #emitContextEvent(event, detail) {
        this.#panicIfDestroyed();
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
        this.#panicIfDestroyed();
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
        this.#panicIfDestroyed();
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
        this.#panicIfDestroyed();
        this.#context_locked = true;
    }

    /**
     * Panics if called and the manager is destroyed
     * @throws {Error} if the manager is destroyed
     */
    #panicIfDestroyed() {
        if (this.#destroyed) {
            throw new Error("The hotkeys manager is destroyed and should no longer be used");
        }
    }

    /**
     * Registers the current context's hotkeys
     * @throws {Error} if no context is loaded
     */
    #registerCurrentContext() {
        this.#panicIfDestroyed();
        if (!hasWindowContext()) return;

        if (this.#current_context === null) {
            throw new Error("No context loaded");
        }

        this.#hotkeys_controller.bindContext(this.#current_context);
    }

    /**
     * Registers a hotkey without persisting it on a context, which means that calling loadContext will overwrite it
     * @param {string|string[]} hotkey the hotkey's name or an array of hotkeys, if already registered it will be overwritten
     * @param {function} callback the callback to be called when the key is pressed
     * @param {import('./hotkeys').HotkeyRegisterOptions} options
     * @deprecated
     */
    registerHotkey(hotkey, callback, options) {
        this.#panicIfDestroyed();
        if (!hasWindowContext()) return;

        console.warn("In HotkeyContextManager.registerHotkey: This method is deprecated, use registerHotkeyOnContext instead.");
        
        options = {...default_hotkey_register_options, ...options};
        // Mousetrap.bind(hotkey, callback, options.mode)
    }

    /**
     * Registers a hotkey on the current context
     * @param {string|string[]} hotkey the hotkey's name or an array of hotkeys, if already registered it will be overwritten
     * @param {function} callback the callback to be called when the key is pressed
     * @param {import('./hotkeys').HotkeyRegisterOptions} options
     */
    registerHotkeyOnContext(hotkey, callback, options) {
        this.#panicIfDestroyed();
        if (!hasWindowContext() || !this.hasLoadedContext()) return;

        options = {...default_hotkey_register_options, ...options};
        this.#current_context.register(hotkey, callback, options);

        this.reloadCurrentContext();

        this.#emitContextEvent(hotkeys_context_events.CONTEXT_CHANGED, {context_name: this.#current_context_name});
    }

    /**
     * Reloads the current context. Useful when the context's hotkeys are changed
     * @throws {Error} if no context is loaded
     */
    reloadCurrentContext() {
        this.#panicIfDestroyed();
        if (!this.hasLoadedContext()) {
            throw new Error("No context loaded");
        }

        this.#hotkeys_controller.dropContext();
        this.#registerCurrentContext();
    }

    /**
     * Releases the context control, allowing contexts to be loaded again. If a context load was requested while the control was locked, 
     * the last requested context will be loaded.
     * @returns {void}
     */
    unlockContextControl() {
        this.#panicIfDestroyed();
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
        this.#panicIfDestroyed();
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
        this.#panicIfDestroyed();
        if (!hasWindowContext()) return;

        this.#current_context.unregister(hotkey, mode);

        this.#emitContextEvent(hotkeys_context_events.CONTEXT_CHANGED, {context_name: this.#current_context_name});
    }

    /**
     * Unregisters the current context's hotkeys
     */
    unregisterCurrentContext() {
        this.#panicIfDestroyed();
        if (!hasWindowContext() || !this.hasLoadedContext()) return;

        this.#hotkeys_controller.dropContext();

        this.#current_context = null;
        this.#current_context_name = undefined;

        this.#emitContextEvent(hotkeys_context_events.CONTEXT_CHANGED, {context_name: ""});
    }
}

/**
 * The global hotkeys manager. null if the environment doesn't support event listeners
 * @type {HotkeyContextManager | null}
 */
let global_hotkeys_manager = null;

export const setupHotkeysManager = () => {
    let emit_events = canUseDOMEvents();
    if (global_hotkeys_manager != null) {
        global_hotkeys_manager.destroy();
    }

    global_hotkeys_manager = new HotkeyContextManager(emit_events);

    globalThis.global_hotkeys_manager = global_hotkeys_manager;
}

export const getHotkeysManager = () => {
    return global_hotkeys_manager;
}