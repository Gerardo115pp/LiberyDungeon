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
     * @type {Map<string, import('./hotkeys').HotkeyData[]>}
     */
    #keydown_hotkey_triggers;

    /**
     * A map of Hotkey triggers -> HotkeyData[] for keyup events. 
     * @type {Map<string, import('./hotkeys').HotkeyData[]>}
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
     * Activates a given hotkey.
     * @param {import('./hotkeys').HotkeyData} hotkey
     * @param {KeyboardEvent} event
     */
    #activateHotkey(hotkey, event) {
        if (hotkey == null) {
            throw new Error("Hotkey is null")
        }
        console.log(`Activating hotkey ${hotkey.key_combo}`);

        hotkey.run(event);
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
        globalThis.removeEventListener("keydown", this.#bound_handleKeyDown);
        globalThis.removeEventListener("keyup", this.#bound_handleKeyUp);
    }

    /**
     * Returns the last N events of the given type.
     * @param {"keydown"|"keyup"} event_type
     * @param {number} n
     * @returns {KeyboardEvent[]}
     */
    #getLastEvents(event_type, n) {
        if (n < 0) {
            throw new Error("n must be greater than 0")
        }

        let all_past_events = event_type === "keydown" ? this.#keyboard_past_keydowns : this.#keyboard_past_keyups;

        if (all_past_events.Size < n) {
            throw new Error(`There are not enough past ${event_type} events to return ${n} of them`);
        }

        let last_n_events = [];

        for (let h = 0; h < n; h++) {
            let event =all_past_events.PeekN(h);

            if (event == null) break;
            
            last_n_events.push(event);
        }

        return last_n_events;
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
        if (this.#shouldIgnoreEvent(event)) return;
        console.log("keydown", event) 

        let matching_hotkey = this.#matchHotkey(event);

        if (matching_hotkey != null) {
            this.#activateHotkey(matching_hotkey, event);
        }

        this.#keyboard_past_keydowns.Add(event);
    }

    /**
     * Handles the keyup event
     * @param {KeyboardEvent} event
     */
    #handleKeyUp(event) {
        if (this.#shouldIgnoreEvent(event)) return;
        console.log("keyup", event);

        let matching_hotkey = this.#matchHotkey(event);

        if (matching_hotkey != null) {
            this.#activateHotkey(matching_hotkey, event);
        }

        this.#keyboard_past_keyups.Add(event);
    }

    /**
     * Matches a Keyboard event with registered hotkeys. if it finds a match, returns the the hotkey.
     * @param {KeyboardEvent} event
     * @returns {import('./hotkeys').HotkeyData | null}
     */
    #matchHotkey(event) {
        if (this.#current_hotkey_context == null) return null;

        let is_keydown = event.type === "keydown";
        let is_keyup = event.type === "keyup";

        if (!is_keydown && !is_keyup) {
            throw new Error("Event type not supported")
        }

        let past_events = is_keydown ? this.#keyboard_past_keydowns : this.#keyboard_past_keyups;

        let candidate_hotkeys = this.#getCandidateHotkeys(event.key, event.type);

        if (candidate_hotkeys.length === 0) return;

        /** @type {import('./hotkeys').HotkeyData | null} */
        let matching_hotkey = null;

        for (let hotkey of candidate_hotkeys) {
            let matching_events = [event];

            if (hotkey.Length > 1) {
                if (hotkey.Length > past_events.Size) {
                    continue;
                }

                matching_events = [
                    ...this.#getLastEvents(event.type, hotkey.Length - 1),
                    event
                ]
            }

            if (hotkey.match(matching_events)) {
                matching_hotkey = hotkey;
                break;
            }
        }

        console.log("Matched hotkey", matching_hotkey);

        return matching_hotkey;
    }

    /**
     * Returns the candidate hotkeys for a given key. Meaning the hotkeys that could potentially be triggered by the key.
     * @param {string} key
     * @param {string} mode
     * @returns {import('./hotkeys').HotkeyData[]}
     */
    #getCandidateHotkeys(key, mode) {
        let candidate_hotkeys = [];

        let triggers = mode === "keydown" ? this.#keydown_hotkey_triggers : this.#keyup_hotkey_triggers;

        let lower_case_matching = triggers.get(key.toLowerCase()) ?? [];
        let upper_case_matching = triggers.get(key.toUpperCase()) ?? [];

        candidate_hotkeys = [
            ...lower_case_matching,
            ...upper_case_matching
        ];

        return candidate_hotkeys;
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

            let all_triggers = hotkey.Mode === "keydown" ? this.#keydown_hotkey_triggers : this.#keyup_hotkey_triggers;

            const hotkey_triggers = hotkey.Triggers;

            for (let trigger of hotkey_triggers) {
                let hotkeys_on_trigger = [];

                if (all_triggers.has(trigger)) {
                    hotkeys_on_trigger = all_triggers.get(trigger);
                }

                hotkeys_on_trigger.push(hotkey);

                all_triggers.set(trigger, hotkeys_on_trigger);
            }
        }
    }

    /**
     * Attaches the binder to the global this object
     */
    #setup() {
        globalThis.addEventListener("keydown", this.#bound_handleKeyDown);
        globalThis.addEventListener("keyup", this.#bound_handleKeyUp);
    }

    /**
     * Whether a given KeyboardEvent was originated from a Input like element or not. In which case, the
     * event should not trigger any hotkeys.
     * @param {KeyboardEvent} event
     */
    #shouldIgnoreEvent(event) {
        let origin_element = event.target;

        if (origin_element == null) return false;

        let should_ignore = false; 

        if (origin_element instanceof HTMLInputElement) {
            should_ignore = true;
        }

        if (!should_ignore && origin_element instanceof HTMLTextAreaElement) {
            should_ignore = true;
        }

        return should_ignore;
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