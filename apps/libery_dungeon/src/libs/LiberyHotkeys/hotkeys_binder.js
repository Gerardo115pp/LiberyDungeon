import { StackBuffer } from "@libs/utils";
import { 
    MAX_PAST_EVENTS, 
    MAX_TIME_BETWEEN_SEQUENCE_KEYSTROKES
} from "./hotkeys_consts";


export class HotkeysController {

    /**
     * A stack with the last MAX_PAST_EVENTS keydown KeyboardEvents
     * @type {StackBuffer<KeyboardEvent>}
     */
    #keyboard_past_keydowns;

    /**
     * A stack with the last MAX_PAST_EVENTS keyup KeyboardEvents
     * @type {StackBuffer<KeyboardEvent>}
     */
    #keyboard_past_keyups;

    /**
     * keydown handler function bound to the HotkeyBinder instance so addEventListener wont override the this context with the EventTarget.
     * @type {function}
     */
    #bound_handleKeyDown;

    /**
     * keyup handler function bound to the HotkeyBinder instance so addEventListener wont override the this context with the EventTarget.
     * @type {function}
     */
    #bound_handleKeyUp;

    /**
     * A map of Hotkey triggers -> HotkeyData for keydown events. 
     * @type {Map<string, HotkeyData[]>}
     */
    #keydown_hotkey_triggers;

    /**
     * A map of Hotkey triggers -> HotkeyData[] for keyup events. 
     * @type {Map<string, HotkeyData[]>}
     */
    #keyup_hotkey_triggers;

    /**
     * The current hotkey context.
     * @type {import('./hotkeys_context').default | null}
     * @default null
     */
    #current_hotkey_context;


    constructor() {
        if (globalThis.addEventListener === undefined) {
            throw new Error("This environment does not support event listeners")
        }

        this.#keyboard_past_keydowns = new StackBuffer(MAX_PAST_EVENTS);
        this.#keyboard_past_keyups = new StackBuffer(MAX_PAST_EVENTS);

        this.#keydown_hotkey_triggers = new Map();
        this.#keyup_hotkey_triggers = new Map();
        
        this.#bound_handleKeyDown = this.#handleKeyDown.bind(this);
        this.#bound_handleKeyUp = this.#handleKeyUp.bind(this);

        this.#current_hotkey_context = null;

        this.#setup()
    }

    /**
     * Sets a new hotkey context. if there is already a context set, it will be disabled.
     * @param {import('./hotkeys_context').default} new_context
     */
    bindContext(new_context) {
        if (this.#current_hotkey_context != null) {
            this.dropContext();
        }

        this.#current_hotkey_context = new_context;

        this.#populateHotkeyTriggers(new_context);
    }

    /**
     * Drops the current hotkey context
     * @returns {void}
     */
    dropContext() {
        if (this.#current_hotkey_context == null) return;

        this.#keydown_hotkey_triggers.clear();
        this.#keyup_hotkey_triggers.clear();

        this.#keyboard_past_keydowns.Clear(); // Prevent key strokes issued in previous contexts from affecting behavior in the new context.
        this.#keyboard_past_keyups.Clear(); 

        this.#current_hotkey_context = null;
    }

    /**
     * Cleans up the binder
     */
    destroy() {
        globalThis.removeEventListener("keydown", this.#handleKeyDown);
        globalThis.removeEventListener("keyup", this.#handleKeyUp);
    }

    /**
     * Whether the binder has an active hotkey context or not.
     * @type {boolean}
     */
    get HasContext() {
        return this.#current_hotkey_context != null;
    }

    /**
     * Handles the keydown event
     * @param {KeyboardEvent} event
     */
    #handleKeyDown(event) {
        console.log("keydown", event) 

        this.#keyboard_past_keydowns.Add(event);
    }

    /**
     * Handles the keyup event
     * @param {KeyboardEvent} event
     */
    #handleKeyUp(event) {
        console.log("keyup", event)

        this.#keyboard_past_keyups.Add(event);
    }

    /**
     * Matches a Keyboard event with registered hotkeys
     * @param {KeyboardEvent} event
     */
    #matchHotkey(event) {
        let is_keydown = event.type === "keydown";
        let is_keyup = event.type === "keyup";

        if (!is_keydown && !is_keyup) {
            throw new Error("Event type not supported")
        }

        let past_events = is_keydown ? this.#keyboard_past_keydowns : this.#keyboard_past_keyups;
        let triggers = is_keydown ? this.#keydown_hotkey_triggers : this.#keyup_hotkey_triggers;


    }

    /**
     * The past keydown events stack.
     * @returns {StackBuffer<KeyboardEvent>}
     */
    get PastKeyDowns() {
        return this.#keyboard_past_keydowns;
    }

    /**
     * The past keyup events stack.
     * @returns {StackBuffer<KeyboardEvent>}
     */
    get PastKeyUps() {
        return this.#keyboard_past_keyups;
    }

    /**
     * Populates the hotkey triggers maps with the hotkeys from given context.
     * @param {import('./hotkeys_context').default} context
     * @returns {void}
     */
    #populateHotkeyTriggers(context) {
        let hotkeys = context.hotkeys;

        for (let hotkey of hotkeys) {
            if (hotkey.Mode !== "keydown" && hotkey.Mode !== "keyup") {
                console.error(`Hotkey mode ${hotkey.Mode} is not supported. Skipping hotkey ${hotkey.KeyCombo}`);
                continue;
            }

            let triggers = hotkey.Mode === "keydown" ? this.#keydown_hotkey_triggers : this.#keyup_hotkey_triggers;

            const hotkey_trigger = hotkey.Trigger;

            let trigger_hotkeys = [];

            if (triggers.has(hotkey_trigger)) {
                trigger_hotkeys = triggers.get(hotkey_trigger);
            }

            trigger_hotkeys.push(hotkey);

            triggers.set(hotkey_trigger, trigger_hotkeys);
        }
    }

    /**
     * Attaches the binder to the global this object
     */
    #setup() {
        globalThis.addEventListener("keydown", this.#bound_handleKeyDown);
        globalThis.addEventListener("keyup", this.#bound_handleKeyUp);
    }
}

/**
 * The global hotkeys binder
 * @type {HotkeysController | null}
 */
let GlobalHotkeysBinder = null;

export const setupHotkeysBinder = () => {
    if (GlobalHotkeysBinder != null) {
        throw new Error("Hotkeys binder already exists. Call destroyHotkeysBinder before setting up a new one")
    }

    GlobalHotkeysBinder = new HotkeysController();
}

export const destroyHotkeysBinder = () => {
    if (GlobalHotkeysBinder == null) {
        throw new Error("No hotkeys binder exists. Call setupHotkeysBinder before destroying it")
    }

    GlobalHotkeysBinder.destroy();

    GlobalHotkeysBinder = null;
}

/**
 * Returns the global hotkeys binder.
 * @returns {HotkeysController | null}
 */
export const getHotkeysBinder = () => {
    return GlobalHotkeysBinder;
}